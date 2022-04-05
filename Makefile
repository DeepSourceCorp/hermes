NAME=hermes
BIN := hermes

.PHONY: build
## build: Build hermes
build: $(BIN)
	go build -o ${BIN} cmd/server/server.go

.PHONY: clean
## clean: Clean projects and previous builds
clean:
	@rm -rf $(NAME)

.PHONY: deps
## deps: Download modules
deps:
	@go mod download

.PHONY: major
## major: Create a major release
major:
	@git pull --tags; \
	IFS='.' read -ra tag <<< "$$(git describe --tags `git rev-list --tags --max-count=1`)"; \
	bump=$$(($${tag[0]:1} + 1)); \
	ver=v$$bump.0.0; \
	rem=$$(git remote); \
	git tag $$ver; \
	git push $$rem $$ver

.PHONY: minor
## minor: Create a minor release
minor:
	@git pull --tags; \
	IFS='.' read -ra tag <<< "$$(git describe --tags `git rev-list --tags --max-count=1`)"; \
	bump=$$(($${tag[1]} + 1)); \
	ver=$${tag[0]}.$$bump.0; \
	rem=$$(git remote); \
	git tag $$ver; \
	git push $$rem $$ver

.PHONY: patch
## patch: Create a patch
patch:
	@git pull --tags; \
	IFS='.' read -ra tag <<< "$$(git describe --tags `git rev-list --tags --max-count=1`)"; \
	bump=$$(($${tag[2]} + 1)); \
	ver=$${tag[0]}.$${tag[1]}.$$bump; \
	rem=$$(git remote); \
	git tag $$ver; \
	git push $$rem $$ver

.PHONY: help
all: help
## help: show this help message
help: makefile
	@echo
	@echo " Choose a command to run:"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' | sed -e 's/^/ /'
	@echo


.DEFAULT_GOAL := help
