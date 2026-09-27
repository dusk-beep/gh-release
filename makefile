build:
	mkdir -p bin/
	go build -o bin/gh-release cmd/gh-release/main.go 

install: build
	mv -v bin/gh-release ~/.local/bin/

clean:
	rm -rv bin/
