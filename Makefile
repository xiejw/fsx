BUILD=.build
PACKAGES=github.com/xiejw/fsx/go/...

compile:
	go build ${PACKAGES}

fmt:
	go fmt ${PACKAGES}

test:
	go test ${PACKAGES}

bench:
	go test -bench=. ${PACKAGES}

# use go get package/path to update
tidy:
	go mod tidy

binary:
	mkdir -p ${BUILD} && \
	  go build -o ${BUILD}/snapshot cmd/snapshot/main.go
