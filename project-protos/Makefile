OUT_DIR := .

generate:
	protoc \
	  --go_out=$(OUT_DIR) --go_opt=paths=source_relative \
	  --go-grpc_out=$(OUT_DIR) --go-grpc_opt=paths=source_relative \
	  proto/user/user.proto proto/task/task.proto

clean:
	find . -name "*.pb.go" -delete
