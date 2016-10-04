SET GOOS=windows
SET GOARCH=amd64


go install github.com/pieterlouw/go-jsonapigateway_tmpl/gateway
go install github.com/pieterlouw/go-jsonapigateway_tmpl/boltdb
go build -o bin/authswitch.exe cmd/authswitch/main.go
