build:
	mkdir -p ./dist
	cp ./static/* ./dist
	cp $$(go env GOROOT)/misc/wasm/wasm_exec.js ./dist/wasm_exec.js
	GOOS=js GOARCH=wasm go build -o ./dist/main.wasm

deploy: build
	aws s3 sync --delete ./dist/ s3://gogogo-frontend
