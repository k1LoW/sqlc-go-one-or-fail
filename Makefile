PKG = github.com/k1LoW/sqlc-go-one-or-fail
COMMIT = $$(git describe --tags --always)
OSNAME=${shell uname -s}
ifeq ($(OSNAME),Darwin)
	DATE = $$(gdate --utc '+%Y-%m-%d_%H:%M:%S')
else
	DATE = $$(date --utc '+%Y-%m-%d_%H:%M:%S')
endif

export GO111MODULE=on

BUILD_LDFLAGS = -X $(PKG).commit=$(COMMIT) -X $(PKG).date=$(DATE)

default: test

ci: depsdev test test-integration

build:
	go build -ldflags="$(BUILD_LDFLAGS)" -o sqlc-go-one-or-fail main.go

install: build
	go install

lint:
	golangci-lint run ./...
	govulncheck ./...

test:
	go test ./... -coverprofile=coverage.out -covermode=count

test-integration:
	rm -rf internal/rewriter/testdata/default/gen
	cd internal/rewriter/testdata/default && sqlc generate
	cd internal/rewriter/testdata/default && go run main.go
	go run main.go internal/rewriter/testdata/default/gen
	(cd internal/rewriter/testdata/default && go run main.go) || if [ $$? -ne 1 ]; then echo "should be failed with status 1" && exit 1 ; fi;
	rm -rf internal/rewriter/testdata/emit_methods_with_db_argument/gen
	cd internal/rewriter/testdata/emit_methods_with_db_argument && sqlc generate
	cd internal/rewriter/testdata/emit_methods_with_db_argument && go run main.go
	go run main.go internal/rewriter/testdata/emit_methods_with_db_argument/gen
	(cd internal/rewriter/testdata/emit_methods_with_db_argument && go run main.go) || if [ $$? -ne 1 ]; then echo "should be failed with status 1" && exit 1 ; fi;
	rm -rf internal/rewriter/testdata/emit_result_struct_pointers/gen
	cd internal/rewriter/testdata/emit_result_struct_pointers && sqlc generate
	cd internal/rewriter/testdata/emit_result_struct_pointers && go run main.go
	go run main.go internal/rewriter/testdata/emit_result_struct_pointers/gen
	(cd internal/rewriter/testdata/emit_result_struct_pointers && go run main.go) || if [ $$? -ne 1 ]; then echo "should be failed with status 1" && exit 1 ; fi;

depsdev:
	go install github.com/Songmu/ghch/cmd/ghch@latest
	go install github.com/Songmu/gocredits/cmd/gocredits@latest
	go install golang.org/x/vuln/cmd/govulncheck@latest
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

prerelease_for_tagpr:
	gocredits -w .
	git add CHANGELOG.md CREDITS go.mod go.sum

release:
	git push origin main --tag
	goreleaser --clean

.PHONY: default test
