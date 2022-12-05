.PHONY: audit
audit:
	go list -m all | nancy sleuth

.PHONY: test
test:
	go test --cover -v ./...

.PHONY: local
local:
	docker compose up

.PHONY: down
down:
	docker compose down