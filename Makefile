gen:
	protoc --go_out=. --go_opt=paths=source_relative     --go-grpc_out=. --go-grpc_opt=paths=source_relative models/*.proto
clean:
	rm -f models/*.pb.go
test:
	go test -cover -race ./service
server:
	go run cmd/server/main.go -port 8080
clients:
	go run cmd/client/main.go -address 127.0.0.1:8080
gencert:
	cd cert && rm *.pem &&./gen.sh&&cd ..