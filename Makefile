.PHONY: clean test all

clean:
	@rm internal/pb/*.go

gen:
	@mkdir -p internal/pb
	@protoc -I proto --go_out=internal/pb/ --go_opt=paths=source_relative \
	--go-grpc_out=internal/pb/ --go-grpc_opt=paths=source_relative proto/*.proto

build:
	@go build -o bin/cleanbg cmd/cleanbg/main.go

run:build
	@bin/cleanbg


