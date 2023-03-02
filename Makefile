tidy:
	go mod tidy

build:
	go build -o bin/blocker

run: build
	./bin/blocker

test:
	go test -v ./...


# ==========================
# Proto

.PHONY: proto
proto:
	protoc --go_out=. --go_opt=paths=source_relative \
	--go-grpc_out=. --go-grpc_opt=paths=source_relative \
	proto/*.proto


# will format the proto file to look like go code, using the .clang-format file.
proto/format:
	@clang-format -i ./*/*.proto