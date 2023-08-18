package rewriter

// func run() error {
// 	var req pb.CodeGenRequest
// 	reqBlob, err := io.ReadAll(os.Stdin)
// 	if err != nil {
// 		return err
// 	}
// 	if err := proto.Unmarshal(reqBlob, &req); err != nil {
// 		return err
// 	}
// 	resp, err := gen(context.Background(), &req)
// 	if err != nil {
// 		return err
// 	}
// 	respBlob, err := proto.Marshal(resp)
// 	if err != nil {
// 		return err
// 	}
// 	w := bufio.NewWriter(os.Stdout)
// 	if _, err := w.Write(respBlob); err != nil {
// 		return err
// 	}
// 	if err := w.Flush(); err != nil {
// 		return err
// 	}
// 	return nil
// }

// func gen(ctx context.Context, req *pb.CodeGenRequest) (*pb.CodeGenResponse, error) {
// 	tmpdir, err := os.MkdirTemp(os.TempDir(), "sqlc-go-one-or-fail-")
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer os.RemoveAll(tmpdir)
// 	b := req.GetPluginOptions()
// 	g := map[string]any{}
// 	if err := json.Unmarshal(b, &g); err != nil {
// 		return nil, err
// 	}
// 	g["out"] = tmpdir
// 	c := map[string]any{
// 		"version": "2",
// 		"sql": []map[string]any{
// 			{
// 				"schema":  req.GetSettings().GetSchema()[0],
// 				"queries": req.GetSettings().GetQueries()[0],
// 				"engine":  req.GetSettings().GetEngine(),
// 				"gen": map[string]any{
// 					"go": g,
// 				},
// 			},
// 		},
// 	}
// 	cb, err := json.Marshal(c)
// 	if err != nil {
// 		return nil, err
// 	}
// 	tmpc := filepath.Join(tmpdir, "sqlc.yaml")
// 	if err := os.WriteFile(tmpc, cb, os.ModePerm); err != nil {
// 		return nil, err
// 	}
// 	_, _ = fmt.Fprintf(os.Stderr, "%#v\n", os.Environ())
// 	cmd := exec.CommandContext(ctx, "sqlc", "generate", "-f", tmpc)
// 	cmd.Stderr = os.Stderr
// 	o, err := cmd.Output()
// 	if err != nil {
// 		return nil, err
// 	}
// 	_, _ = fmt.Fprintf(os.Stderr, "%#v\n", string(o))

// 	resp := &pb.CodeGenResponse{}
// 	out := req.GetSettings().GetCodegen().GetOut()
// 	var files []string
// 	for _, q := range req.Queries {
// 		if q.Filename == "" {
// 			continue
// 		}
// 		p := fmt.Sprintf("%s/%s.go", out, q.Filename)
// 		files = append(files, p)
// 	}
// 	files = unique(files)
// 	for _, p := range files {
// 		_, err := replace(ctx, p)
// 		if err != nil {
// 			return nil, err
// 		}
// 		// resp.Files = append(resp.Files, f)
// 	}
// 	return resp, nil
// }
