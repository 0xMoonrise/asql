ENTRY=./cmd/cli
TARGET=asql

all:
	go run $(ENTRY) -f code.asql

build:
	go build -o $(TARGET) $(ENTRY)

clean:
	rm $(TARGET)

test:
	go test ./...

server:
	go run ./cmd/server
