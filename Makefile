format:
	go fmt ./...
.PHONY: format

test: format
	go test -v -count=1 ./...
.PHONY: test

server:
	oapi-codegen -config configs/server.cfg.yaml api/openapi.yaml
.PHONY: server

service:
	protoc --go_out=./internal/generated/clients \
	--go-grpc_out=./internal/generated/clients \
	./api/proto/geo_service.proto
.PHONY: service

service2:
	protoc --go_out=./internal/generated ./api/proto/basket_confirmed.proto
.PHONY: service2