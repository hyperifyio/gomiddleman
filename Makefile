.PHONY: build clean tidy

GOMIDDLEMAN_SOURCES := \
    ./internal/gomiddleman/connectionhandlers/connectionhandler.go \
    ./internal/gomiddleman/connectionhandlers/tcp.go \
    ./internal/gomiddleman/connectionhandlers/tls.go \
    ./internal/gomiddleman/connectors/connector.go \
    ./internal/gomiddleman/connectors/tcp.go \
    ./internal/gomiddleman/connectors/tls.go \
    ./internal/gomiddleman/listeners/listener.go \
    ./internal/gomiddleman/listeners/tcp.go \
    ./internal/gomiddleman/listeners/tls.go \
    ./internal/gomiddleman/tlsutils/loadtlsconfig.go \
    ./internal/gomiddleman/proxyutils/forwardtraffic.go \
    ./internal/gomiddleman/proxyutils/handleconnection.go \
    ./internal/gomiddleman/proxy.go \
    ./cmd/gomiddleman/main.go

all: build

build: gomiddleman

tidy:
	go mod tidy

gomiddleman: $(GOMIDDLEMAN_SOURCES)
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o gomiddleman ./cmd/gomiddleman
	chmod 700 ./gomiddleman

ca-key.pem:
	openssl genrsa -out ca-key.pem 4096

ca.pem: ca-key.pem
	openssl req -x509 -new -nodes -key ca-key.pem -sha256 -days 3650 -out ca.pem -subj "/C=FI/ST=Tampere/L=Tampere/O=HyperifyIO/OU=Developers/CN=HyperifyCA"

client-key.pem:
	openssl genrsa -out client-key.pem 4096

client-csr.pem: client-key.pem
	openssl req -new -key client-key.pem -out client-csr.pem -config client-cert.conf -extensions req_ext

client-cert.pem: ca.pem ca-key.pem client-csr.pem
	openssl x509 -req -in client-csr.pem -CA ca.pem -CAkey ca-key.pem -CAcreateserial -out client-cert.pem -days 3650 -sha256

key.pem:
	openssl genpkey -algorithm RSA -out key.pem -pkeyopt rsa_keygen_bits:2048

cert-csr.pem: cert.conf key.pem
	openssl req -new -key key.pem -out cert-csr.pem -config cert.conf -extensions req_ext

cert.pem: key.pem cert.conf cert-csr.pem ca.pem ca-key.pem
	openssl x509 -req -in cert-csr.pem -CA ca.pem -CAkey ca-key.pem -CAcreateserial -out cert.pem -days 3650 -sha256 -extfile cert.conf -extensions req_ext

test: cert.pem client-cert.pem ca.pem
	go test -v ./...

clean:
	rm -f gomiddleman cert.pem key.pem cert-csr.pem client-csr.pem client-key.pem client-cert.pem ca.pem ca-key.pem
