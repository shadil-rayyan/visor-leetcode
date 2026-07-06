DATA_REPO = leetcode-companywise-interview-questions
DATA_URL  = https://github.com/snehasishroy/leetcode-companywise-interview-questions.git
SKIP_TAGS ?= true

data:
	@if [ ! -d "$(DATA_REPO)/.git" ]; then rm -rf "$(DATA_REPO)" && git clone --depth 1 $(DATA_URL) $(DATA_REPO); fi
	cd merger && ROOT_DIR=../$(DATA_REPO) SKIP_TAGS=$(SKIP_TAGS) go run . --json

dev:
	cd web && npm run dev

build:
	cd web && npm run build

.PHONY: data dev build
