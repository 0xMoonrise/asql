ENTRY=./cmd/cli
TARGET=asql
RUN ?= .

all:
	go run $(ENTRY) -f code.asql

build:
	go build -o $(TARGET) $(ENTRY)

clean:
	rm $(TARGET)

server:
	go run ./cmd/server

cli:
	go run $(ENTRY)

debug:
	dlv debug $(ENTRY) -- -f code.asql

test:
	gotestsum --format testname -- -run $(RUN) ./...
