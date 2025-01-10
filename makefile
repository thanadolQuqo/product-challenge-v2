run:
	go run cmd/api/main.go

db:
	cockroach start-single-node --advertise-addr 'localhost' --insecure

install-db:
	brew install cockroachdb/tap/cockroach