# proto
# Используем bin в текущей директории для установки плагинов protoc
LOCAL_BIN:=$(CURDIR)/bin

# Добавляем bin в текущей директории в PATH при запуске protoc
#PROTOC = PATH="$$PATH:$(LOCAL_BIN)" protoc

ORDER_PROTO_PATH:=api/proto/checker/v1
ORDER_PROTO_PATH_OUT:=api
ORDER_DOCS_PATH:=docs

# Установка всех необходимых зависимостей
.PHONY: .bin-deps
bin-deps:
	$(info Installing binary dependencies...)

	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	GOBIN=$(LOCAL_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	GOBIN=$(LOCAL_BIN) go install github.com/envoyproxy/protoc-gen-validate@latest

# Вендоринг внешних proto файлов
vendor-proto: vendor-proto-rm vendor-proto/google/protobuf vendor-proto/google/api vendor-proto/protoc-gen-openapiv2/options vendor-proto/validate

vendor-proto-rm:
	rm -fdr 'vendor.proto' || true

# Устанавливаем proto описания protoc-gen-openapiv2/options
.PHONY: vendor-proto/protoc-gen-openapiv2/options
vendor-proto/protoc-gen-openapiv2/options:
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
 		https://github.com/grpc-ecosystem/grpc-gateway vendor.proto/grpc-ecosystem && \
 	cd vendor.proto/grpc-ecosystem && \
	git sparse-checkout set --no-cone protoc-gen-openapiv2/options && \
	git checkout
	mkdir -p vendor.proto/protoc-gen-openapiv2
	mv vendor.proto/grpc-ecosystem/protoc-gen-openapiv2/options vendor.proto/protoc-gen-openapiv2
	rm -rf vendor.proto/grpc-ecosystem

# Устанавливаем proto описания google/protobuf
.PHONY: vendor-proto/google/protobuf
vendor-proto/google/protobuf:
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
		https://github.com/protocolbuffers/protobuf vendor.proto/protobuf &&\
	cd vendor.proto/protobuf &&\
	git sparse-checkout set --no-cone src/google/protobuf &&\
	git checkout
	mkdir -p vendor.proto/google
	mv vendor.proto/protobuf/src/google/protobuf vendor.proto/google
	rm -rf vendor.proto/protobuf

.PHONY: vendor-proto/google/api
vendor-proto/google/api:
	git clone -b master --single-branch -n --depth=1 --filter=tree:0 \
 		https://github.com/googleapis/googleapis vendor.proto/googleapis && \
 	cd vendor.proto/googleapis && \
	git sparse-checkout set --no-cone google/api && \
	git checkout
	mkdir -p  vendor.proto/google
	mv vendor.proto/googleapis/google/api vendor.proto/google
	rm -rf vendor.proto/googleapis

.PHONY: vendor-proto/validate
vendor-proto/validate:
	git clone -b main --single-branch --depth=2 --filter=tree:0 \
		https://github.com/bufbuild/protoc-gen-validate vendor.proto/tmp && \
		cd vendor.proto/tmp && \
		git sparse-checkout set --no-cone validate &&\
		git checkout
		mkdir -p vendor.proto/validate
		mv vendor.proto/tmp/validate vendor.proto/
		rm -rf vendor.proto/tmp

.PHONY: generate-proto
generate-proto:
	mkdir -p pkg/$(ORDER_PROTO_PATH_OUT)
	mkdir -p $(ORDER_DOCS_PATH)
	protoc -I api/proto \
		${ORDER_PROTO_PATH}/checker.proto \
		--go_out=./pkg/$(ORDER_PROTO_PATH_OUT) --go_opt=paths=source_relative\
		--go-grpc_out=./pkg/$(ORDER_PROTO_PATH_OUT) --go-grpc_opt=paths=source_relative \

do-all: bin-deps vendor-proto generate-proto
