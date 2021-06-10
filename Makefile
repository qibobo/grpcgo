gen:
	# protoc --proto_path=models --go_out=models --go_opt=paths=source_relative     --go-grpc_out=models --go-grpc_opt=paths=source_relative models/*.proto
	protoc --proto_path=models --go_out=. --go-grpc_out=models    --go-grpc_opt=paths=source_relative --grpc-gateway_out=:. --openapiv2_out=:swagger models/*.proto  
clean:
	rm -f models/*.pb.go && rm -f models/*.gw.go
test:
	go test -cover -race ./service
server:
	go run cmd/server/main.go -port 8080
rest:
	go run cmd/server/main.go -port 8080 -serverType rest
clients:
	go run cmd/client/main.go -address 127.0.0.1:8080
gencert:
	cd cert && rm *.pem &&./gen.sh&&cd ..