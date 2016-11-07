CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/badger-go .
docker build -t badger-go .
