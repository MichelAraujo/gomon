
build:
	go build -o Gomon ./main.go

buildToTest: build
	mv ./Gomon ./tests