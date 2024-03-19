.PHONY: clean test all

include .env
export

clean:
	@rm internal/pb/*.go

gen:
	@mkdir -p internal/pb
	@protoc -I proto --go_out=internal/pb/ --go_opt=paths=source_relative \
	--go-grpc_out=internal/pb/ --go-grpc_opt=paths=source_relative proto/*.proto

build:
	@go build -o bin/cleanbg cmd/cleanbg/main.go
	@go build -o bin/tgbot cmd/tgbot/main.go

run-http:build
	@bin/cleanbg

run-tg:build
	@bin/tgbot

run-service:
	@cd src && python3 -m main.py
