build:
				npm install
				npm run build-css	
				npm run build	
				GOOS=linux GOARCH=amd64 go build -o skeef *.go
