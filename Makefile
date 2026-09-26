
.PHONY: build

build: 
	rm -rf build
	mkdir -p build
	go get 
	go tool templ generate
	go build -o build/app
	cp -r migrations build
	cp -r statics build


