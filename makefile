default: darwin
	GOOS=linux GOARCH=amd64 go build -o bin/skeef *.go

darwin: windows
	GOOS=darwin GOARCH=amd64 go build -o bin/skeef-mac *.go

windows: npm-install
	GOOS=windows GOARCH=amd64 go build -o bin/skeef.exe *.go

npm-install: css
	npm install

css: build
	npm run dev-css	

build:
	npm run build	
