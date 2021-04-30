update-pkg-cache: go.sum go.mod prototools.go prototools_test.go
	GOPROXY=https://proxy.golang.org GO111MODULE=on \
	go get github.com/johnsiilver/prototools
