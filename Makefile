.PHONY: deps
deps:
		go mod download


.PHONY: fmt
fmt:
		go fmt  ./...


# Runs tests in short mode, without database adapters
.PHONY: quicktest
quicktest:
		go test -short all


# Runs tests in short mode, without database adapters
.PHONY: test
test:
		go test all
