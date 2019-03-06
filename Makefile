all:
	GOOS=windows    GOARCH=386      go build -o ./build/Release/01-get-balance.exe      cmd/01-get-balance/main.go
	GOOS=windows    GOARCH=amd64    go build -o ./build/Release/01-get-balance_64.exe   cmd/01-get-balance/main.go
	GOOS=linux      GOARCH=amd64    go build -o ./build/Release/01-get-balance          cmd/01-get-balance/main.go
	#
	GOOS=windows    GOARCH=386      go build -o ./build/Release/02-get-balance.exe      cmd/02-get-balance/main.go
	GOOS=windows    GOARCH=amd64    go build -o ./build/Release/02-get-balance_64.exe   cmd/02-get-balance/main.go
	GOOS=linux      GOARCH=amd64    go build -o ./build/Release/02-get-balance          cmd/02-get-balance/main.go