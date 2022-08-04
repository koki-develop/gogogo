build:
	mkdir -p ./dist
	cp ./static/* ./dist
	GOOS=js GOARCH=wasm go build -o ./dist/main.wasm
