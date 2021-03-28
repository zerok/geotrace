all: bin/geotrace

bin:
	mkdir -p bin

bin/geotrace: $(shell find . -name '*.go') go.mod bin
	cd cmd/geotrace && go build -mod=mod -o ../../$@

test:
	go test ./... -v

clean:
	rm -rf bin
