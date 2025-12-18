release:
	GOOS=linux GOARCH=amd64; go build -ldflags="-s -w"
	tar -cJf simon.tar.xz simon
	rm -f simon
