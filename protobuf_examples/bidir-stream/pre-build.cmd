
setlocal 

:: https://github.com/protocolbuffers/protobuf/releases взят отсюда

go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

set PATH=C:\Programs\protoc\bin;%PATH%
protoc --go_out=. --go-grpc_out=. chat.proto

go mod init main & go mod tidy