build:
	@go build -o bin/song-library cmd/main.go
run: build
	@./bin/song-library