
.PHONY: build

build: 
	@rm -rf build
	@mkdir -p build
	@echo -e "\033[34m==\033[37m Downloading dependencies"
	@go get 
	@echo -e "\033[34m==\033[37m Building templates"
	@go tool templ generate
	@echo -e "\033[34m==\033[37m Building app"
	@go build -o build/app
	@cp -r migrations build
	@cp -r statics build
	@echo -e "[\033[32mdone\033[0m]"


