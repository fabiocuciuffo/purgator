Arguments GOOS et GOARCH pour build pour macos sillicon
``` zsh
GOOS=darwin GOARCH=arm64 go build -o bin/app-arm64-darwin app.go
```

Compiler le code go dans un container Docker :
``` zsh
docker run --rm -v "$PWD":/usr/src/myapp -w /usr/src/myapp -e GOOS=darwin -e GOARCH=arm64 golang go build -v
ou
docker run --rm -v "$PWD":/usr/src/myapp -w /usr/src/myapp -e GOOS=darwin -e GOARCH=arm64 golang go build -o purgator.bin main.go
```

Run sans compiler :
``` zsh
docker run --rm -v "$PWD":/usr/src/myapp -w /usr/src/myapp  golang go run main.go
```