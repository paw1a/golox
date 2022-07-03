build:
	go build -o golox ./cmd/interpreter/main.go

run: build
	./golox

clean:
	rm -rf ./golox

.DEFAULT_GOAL := build
