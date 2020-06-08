BUILD=.build
PACKAGES=github.com/xiejw/fsx/go/...
CMD_PACKAGES=github.com/xiejw/fsx/cmd/...

compile: compile_pkg
	mkdir -p ${BUILD} && \
	  go build -o ${BUILD}/snapshot cmd/snapshot/main.go

compile_pkg:
	go build ${PACKAGES}

fmt:
	go fmt ${PACKAGES}
	go fmt ${CMD_PACKAGES}

test:
	go test ${PACKAGES}

bench:
	go test -bench=. ${PACKAGES}

# use go get package/path to update
tidy:
	go mod tidy

