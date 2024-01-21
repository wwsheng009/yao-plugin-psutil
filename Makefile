
build:
	CGO_ENABLED=0 go build -o psutil.so

windows:
	GOOS=windows CGO_ENABLED=0 GOARCH=amd64 go build -o psutil.dll

.PHONY:	clean
clean:
	rm -f "../psutil.so"
	rm -f "../psutil.dll"