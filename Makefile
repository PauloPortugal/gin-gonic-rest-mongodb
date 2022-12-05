.PHONY: audit
audit:
	go list -m all | nancy sleuth

.PHONY: test
test:
	go test --cover -v ./...