build:
	mkdir -p bin/
	go build -o bin/gh-release cmd/gh-release/main.go 
clean:
	rm -rv bin/
