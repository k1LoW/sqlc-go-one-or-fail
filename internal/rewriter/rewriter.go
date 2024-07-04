package rewriter

import (
	"context"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io"

	"golang.org/x/tools/go/ast/astutil"
)

func Rewrite(ctx context.Context, p string, w io.Writer) error {
	fset := token.NewFileSet()
	n, err := parser.ParseFile(fset, p, nil, parser.ParseComments|parser.AllErrors)
	if err != nil {
		return err
	}
	rewrited := false
	astutil.Apply(n, nil, func(c *astutil.Cursor) bool {
		n := c.Node()
		switch x := n.(type) {
		case *ast.FuncDecl:
			// Find the following function with 4 statements
			//
			// ```go
			//func (q *Queries) GetAuthorByName(ctx context.Context, name string) (Author, error) {
			//	row := q.db.QueryRowContext(ctx, getAuthorByName, name)
			//	var i Author
			//	err := row.Scan(&i.ID, &i.Name, &i.Bio)
			//	return i, err
			//}
			// ```
			if hasQueryRowContext(x) && hasScanOne(x) && len(x.Body.List) == 4 {
				// Rewrite
				qa := getArgs(x.Body.List[0].(*ast.AssignStmt).Rhs[0].(*ast.CallExpr))
				sa := getArgs(x.Body.List[2].(*ast.AssignStmt).Rhs[0].(*ast.CallExpr))
				vardef := x.Body.List[1]
				x.Body.List = generateStmts(x.Type.Params, x.Type.Results, vardef, qa, sa)
				rewrited = true
				return true
			}
		}
		return true
	})
	if rewrited {
		appendPackage(n, "database/sql")
		appendPackage(n, "fmt")
	}
	if err := format.Node(w, fset, n); err != nil {
		return err
	}
	return nil
}

func appendPackage(n ast.Node, pkg string) {
	spkg := fmt.Sprintf("\"%s\"", pkg)
	astutil.Apply(n, nil, func(c *astutil.Cursor) bool {
		n := c.Node()
		switch x := n.(type) {
		case *ast.GenDecl:
			if x.Tok != token.IMPORT {
				return true
			}
			for _, s := range x.Specs {
				is := s.(*ast.ImportSpec)
				if is.Path.Value == spkg {
					return true
				}
			}
			x.Specs = append(x.Specs, &ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: spkg,
				},
			})
		}
		return true
	})
}

func hasQueryRowContext(n ast.Node) bool {
	var use bool
	astutil.Apply(n, nil, func(c *astutil.Cursor) bool {
		n := c.Node()
		switch x := n.(type) {
		case *ast.CallExpr:
			se, ok := x.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			if se.Sel.Name == "QueryRowContext" {
				use = true
			}
		}
		return true
	})
	return use
}

func getArgs(expr *ast.CallExpr) []ast.Expr {
	a := make([]ast.Expr, len(expr.Args))
	copy(a, expr.Args)
	return a
}

func hasScanOne(n ast.Node) bool {
	var cnt int
	astutil.Apply(n, nil, func(c *astutil.Cursor) bool {
		n := c.Node()
		switch x := n.(type) {
		case *ast.CallExpr:
			se, ok := x.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			if se.Sel.Name == "Scan" {
				cnt++
			}
		}
		return true
	})
	return cnt == 1
}

func generateStmts(params, results *ast.FieldList, vardef ast.Stmt, qa, sa []ast.Expr) []ast.Stmt {
	var zerov, resv, qv string
	var first ast.Expr
	varname := vardef.(*ast.DeclStmt).Decl.(*ast.GenDecl).Specs[0].(*ast.ValueSpec).Names[0].Name
	t, ok := results.List[0].Type.(*ast.Ident)
	if ok {
		switch t.Name {
		case "string":
			zerov = "\"\""
		case "int", "int32", "int64", "uint", "uint32", "uint64":
			zerov = "0"
		default:
			zerov = t.Name + "{}"
		}
		resv = t.Name
		first = &ast.Ident{
			Name: varname,
		}
	} else {
		// emit_result_struct_pointers = true
		switch e := results.List[0].Type.(type) {
		case *ast.StarExpr:
			zerov = "nil"
			resv = e.X.(*ast.Ident).Name
			first = &ast.UnaryExpr{
				Op: token.AND,
				X: &ast.Ident{
					Name: varname,
				},
			}
		case *ast.SelectorExpr:
			// return sql.*, nil
			zerov = varname
			resv = fmt.Sprintf("%s.%s", e.X.(*ast.Ident).Name, e.Sel.Name)
			first = &ast.Ident{
				Name: varname,
			}
		}
	}
	// Get query variable name
	for _, a := range qa {
		if a.(*ast.Ident).Name != "ctx" {
			qv = a.(*ast.Ident).Name
			break
		}
	}
	// Set db instance
	var db ast.Expr
	db = &ast.SelectorExpr{
		X: &ast.Ident{
			Name: "q",
		},
		Sel: &ast.Ident{
			Name: "db",
		},
	}
	for _, p := range params.List {
		if p.Names[0].Name == "db" {
			// emit_methods_with_db_argument = true
			db = &ast.Ident{
				Name: "db",
			}
			break
		}
	}

	stmts := []ast.Stmt{
		&ast.DeclStmt{
			Decl: &ast.GenDecl{
				Tok: token.VAR,
				Specs: []ast.Spec{
					&ast.ValueSpec{
						Names: []*ast.Ident{
							&ast.Ident{
								Name: varname,
							},
						},
						Type: &ast.Ident{
							Name: resv,
						},
					},
				},
			},
		},
		&ast.AssignStmt{
			Lhs: []ast.Expr{
				&ast.Ident{
					Name: "rows",
				},
				&ast.Ident{
					Name: "err",
				},
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X: db,
						Sel: &ast.Ident{
							Name: "QueryContext",
						},
					},
					Args: qa,
				},
			},
		},
		&ast.IfStmt{
			Cond: &ast.BinaryExpr{
				X: &ast.Ident{
					Name: "err",
				},
				Op: token.NEQ,
				Y: &ast.Ident{
					Name: "nil",
				},
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.ReturnStmt{
						Results: []ast.Expr{
							&ast.Ident{
								Name: zerov,
							},
							&ast.Ident{
								Name: "err",
							},
						},
					},
				},
			},
		},
		&ast.DeferStmt{
			Call: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X: &ast.Ident{
						Name: "rows",
					},
					Sel: &ast.Ident{
						Name: "Close",
					},
				},
			},
		},
		&ast.IfStmt{
			Cond: &ast.UnaryExpr{
				Op: token.NOT,
				X: &ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X: &ast.Ident{
							Name: "rows",
						},
						Sel: &ast.Ident{
							Name: "Next",
						},
					},
				},
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.IfStmt{
						Init: &ast.AssignStmt{
							Lhs: []ast.Expr{
								&ast.Ident{
									Name: "err",
								},
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "rows",
										},
										Sel: &ast.Ident{
											Name: "Err",
										},
									},
								},
							},
						},
						Cond: &ast.BinaryExpr{
							X: &ast.Ident{
								Name: "err",
							},
							Op: token.NEQ,
							Y: &ast.Ident{
								Name: "nil",
							},
						},
						Body: &ast.BlockStmt{
							List: []ast.Stmt{
								&ast.ReturnStmt{
									Results: []ast.Expr{
										&ast.Ident{
											Name: zerov,
										},
										&ast.Ident{
											Name: "err",
										},
									},
								},
							},
						},
					},
					&ast.ReturnStmt{
						Results: []ast.Expr{
							&ast.Ident{
								Name: zerov,
							},
							&ast.SelectorExpr{
								X: &ast.Ident{
									Name: "sql",
								},
								Sel: &ast.Ident{
									Name: "ErrNoRows",
								},
							},
						},
					},
				},
			},
		},
		&ast.IfStmt{
			Init: &ast.AssignStmt{
				Lhs: []ast.Expr{
					&ast.Ident{
						Name: "err",
					},
				},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{
					&ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X: &ast.Ident{
								Name: "rows",
							},
							Sel: &ast.Ident{
								Name: "Scan",
							},
						},
						Args: sa,
					},
				},
			},
			Cond: &ast.BinaryExpr{
				X: &ast.Ident{
					Name: "err",
				},
				Op: token.NEQ,
				Y: &ast.Ident{
					Name: "nil",
				},
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.ReturnStmt{
						Results: []ast.Expr{
							&ast.Ident{
								Name: zerov,
							},
							&ast.Ident{
								Name: "err",
							},
						},
					},
				},
			},
		},
		&ast.IfStmt{
			Cond: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X: &ast.Ident{
						Name: "rows",
					},
					Sel: &ast.Ident{
						Name: "Next",
					},
				},
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.ReturnStmt{
						Results: []ast.Expr{
							&ast.Ident{
								Name: zerov,
							},
							&ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "fmt",
									},
									Sel: &ast.Ident{
										Name: "Errorf",
									},
								},
								Args: []ast.Expr{
									&ast.BasicLit{
										Kind:  token.STRING,
										Value: "\"multiple records were retrieved when the following query was executed: %q\"",
									},
									&ast.Ident{
										Name: qv,
									},
								},
							},
						},
					},
				},
			},
		},
		&ast.IfStmt{
			Init: &ast.AssignStmt{
				Lhs: []ast.Expr{
					&ast.Ident{
						Name: "err",
					},
				},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{
					&ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X: &ast.Ident{
								Name: "rows",
							},
							Sel: &ast.Ident{
								Name: "Close",
							},
						},
					},
				},
			},
			Cond: &ast.BinaryExpr{
				X: &ast.Ident{
					Name: "err",
				},
				Op: token.NEQ,
				Y: &ast.Ident{
					Name: "nil",
				},
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.ReturnStmt{
						Results: []ast.Expr{
							&ast.Ident{
								Name: zerov,
							},
							&ast.Ident{
								Name: "err",
							},
						},
					},
				},
			},
		},
		&ast.IfStmt{
			Init: &ast.AssignStmt{
				Lhs: []ast.Expr{
					&ast.Ident{
						Name: "err",
					},
				},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{
					&ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X: &ast.Ident{
								Name: "rows",
							},
							Sel: &ast.Ident{
								Name: "Err",
							},
						},
					},
				},
			},
			Cond: &ast.BinaryExpr{
				X: &ast.Ident{
					Name: "err",
				},
				Op: token.NEQ,
				Y: &ast.Ident{
					Name: "nil",
				},
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.ReturnStmt{
						Results: []ast.Expr{
							&ast.Ident{
								Name: zerov,
							},
							&ast.Ident{
								Name: "err",
							},
						},
					},
				},
			},
		},
		&ast.ReturnStmt{
			Results: []ast.Expr{
				first,
				&ast.Ident{
					Name: "err",
				},
			},
		},
	}
	return stmts
}
