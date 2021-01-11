default: build

build: proto
		go build ./cmd/lored
		go build ./cmd/lore-cli
		go build ./cmd/loreid


install: proto
		go install ./cmd/lored
		go install ./cmd/lore-cli
		go install ./cmd/loreid

proto: grpc/grpc.pb.go

grpc/grpc.pb.go: grpc/grpc.proto
		protoc --go_out=. --go_opt=paths=source_relative grpc/grpc.proto

clean:
		rm lored lore-cli
