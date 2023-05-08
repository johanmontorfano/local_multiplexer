go build -o ./build/epr.exe
cmd /C "SET GOOS=linux&& SET GOARCH=amd64&& go build -o ./build/epr-linux"