
test: gojsx
	./gojsx

static/react-tools.js: index.js
	./node_modules/.bin/browserify ./node_modules/react-tools/main.js > static/react-tools.js

bindata.go: static/react-tools.js
	go generate

gojsx: main.go bindata.go static/react-tools.js
	go build -v