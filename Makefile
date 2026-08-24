.PHONY: test lint wasm gen serve fmt

test:
	go test ./... -race -count=1

lint:
	go vet ./...

wasm:
	GOOS=js GOARCH=wasm go build -o web/static/saint-hubbins.wasm ./cmd/saint-hubbins-wasm && cp $$(go env GOROOT)/lib/wasm/wasm_exec.js web/static/ 2>/dev/null || cp $$(go env GOROOT)/misc/wasm/wasm_exec.js web/static/ 2>/dev/null || echo "wasm build skipped (no wasm target)"

gen:
	go generate ./...

serve:
	go run ./cmd/saint-hubbins serve

fmt:
	gofmt -w .
