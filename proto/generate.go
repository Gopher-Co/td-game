package proto

//go:generate protoc --go_out=../models/coopstate --go-grpc_out=../models/coopstate --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative ./client.proto ./common.proto ./server.proto
