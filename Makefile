build-packets:
	protoc --go_out=./packets packets/proto/*.proto packets/proto/components/*.proto