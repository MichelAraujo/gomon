
build:
	go build -o Gomon ./main.go

build_test: build
	mv ./Gomon ./tests