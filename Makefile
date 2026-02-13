ENTRY=./cmd/asql
TARGET=asql

all:
	go run $(ENTRY) -f code.txt

build:
	go build -o $(TARGET) $(ENTRY)

clean:
	rm $(TARGET)

test:
	go test ./...

