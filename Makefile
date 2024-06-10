build:
	export CGO_ENABLED=0 && go build -o bin/main .

test:
	go test -v -cover ./server/handle/

run:
	go run main.go