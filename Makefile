BINARY="mc.wasm"

all: 
	build

build: 
	GOOS=js GOARCH=wasm  go build -o ${BINARY} ./cmd/MedianCut

