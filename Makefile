.PHONY: test
test:
	go test `go list ./... | grep -v '/mocks'` -race -cover -count=1

.PHONY: deps
deps:
	GO111MODULE=on go mod tidy
	GO111MODULE=on go mod download

.PHONY: bench
bench:
	go test -bench=. -cpu=1,2,4,8,16 -benchmem
