gen:
	protoc --go_out=. --go_opt=paths=source_relative     --go-grpc_out=. --go-grpc_opt=paths=source_relative models/person.proto models/person_service.proto
clean:
	rm -f models/*.pb.go
test:
	go test -cover -race ./service
server:
	go build -o server cmd/server/main.go && ./server -port 8080
client:
	go run cmd/client/main.go -address 127.0.0.1:8080