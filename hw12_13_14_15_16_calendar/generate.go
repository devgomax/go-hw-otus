package calendar

//go:generate protoc -I ./api --go_out ./internal/pb --go_opt paths=source_relative --go-grpc_out ./internal/pb --go-grpc_opt paths=source_relative --grpc-gateway_out ./internal/pb --grpc-gateway_opt paths=source_relative ./api/events/events.proto
