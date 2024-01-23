default: show-config

show-config:
	@echo ">> show-config <<"
	@echo "=== ~/Library/Application\ Support/strom.yml ==="
	@cat ~/Library/Application\ Support/strom.yml
	@echo "--- ~/Library/Application\ Support/strom.yml ---"

clean:
	@echo ">> clean <<"
	@rm -rfv version.go build

generate:
	@echo ">> generate <<"
	@go generate

build: generate
	@echo ">> build <<"
	@go build -o "build/strom" ./...

install: clean build
	@echo ">> install <<"
	@go install ./...