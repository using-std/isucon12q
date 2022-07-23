isuports: test go.mod go.sum *.go cmd/isuports/*
	go build -o isuports ./cmd/isuports

migrate: cmd/migrate/*
	go build -o migrate ./cmd/migrate

test:
	go test -v ./...
