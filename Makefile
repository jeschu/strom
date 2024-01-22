default: show-config

show-config:
	@echo ">> show-config <<"
	@echo "=== ~/Library/Application\ Support/strom.yml ==="
	@cat ~/Library/Application\ Support/strom.yml
	@echo "--- ~/Library/Application\ Support/strom.yml ---"

build:
	@echo ">> build <<"
	@go build -a -o "build/strom" ./...

install: build
	@echo ">> install <<"
	@go install ./...