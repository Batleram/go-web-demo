
.PHONY: build

build: 
	@rm -rf build
	@mkdir -p build
	@echo "\033[34m==\033[37m Downloading dependencies"
	@go get 
	@echo "\033[34m==\033[37m Building templates"
	@go tool templ generate
	@echo "\033[34m==\033[37m Building app"
	@go build -o build/app
	@cp -r migrations build
	@cp -r statics build
	@echo  "[\033[32mdone\033[0m]"


