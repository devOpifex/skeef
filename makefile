build:
				npm install
				npm run dev-css	
				npm run build	
				GOOS=linux GOARCH=amd64 go build -o bin/skeef *.go
				GOOS=windows GOARCH=amd64 go build -o bin/skeef.exe *.go
				GOOS=darwin GOARCH=amd64 go build -o bin/skeef-mac *.go
run:
				go run -race .
