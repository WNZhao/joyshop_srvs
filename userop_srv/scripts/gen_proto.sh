#!/bin/bash

# 生成 proto 文件
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/goods.proto 