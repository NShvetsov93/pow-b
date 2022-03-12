include ./scripts/Makefile

UNAME := $(shell uname)
PROJECT_NAME=pow-b

.PHONY: lintfix
lintfix:
    ifeq ($(UNAME), Linux)
		find . \( -path './cmd/*' -or -path './internal/*' \) \
		-type f -name '*.go' -print0 | \
		xargs -0  sed -i '/import (/,/)/{/^\s*$$/d;}'
    endif
    ifeq ($(UNAME), Darwin)
		find . \( -path './cmd/*' -or -path './internal/*' \) \
		-type f -name '*.go' -print0 | \
		xargs -0  sed -i '' '/import (/,/)/{/^\s*$$/d;}'
    endif
	goimports -local=$(PROJECT_NAME) -w ./cmd/ ./internal

.PHONY: run
run:
	go run cmd/main.go


# Cosmetics
YELLOW := "\e[1;33m"
NC := "\e[0m"
# Shell Functions
INFO := @bash -c '\
  printf $(YELLOW); \
  echo "=> $$1"; \
  printf $(NC)' VALUE