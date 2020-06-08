REPO=fsx

BUILD=.build
PACKAGES=github.com/xiejw/${REPO}/go/...
CMD_PACKAGES=github.com/xiejw/${REPO}/cmd/...
CMDS=$(patsubst cmd/%,%,$(wildcard cmd/*))

compile: compile_pkg compile_cmd

compile_pkg:
	go build ${PACKAGES}

compile_cmd:
	@mkdir -p ${BUILD}
	@for cmd in ${CMDS}; do \
		echo 'compile cmd/'$$cmd && go build -o ${BUILD}/$$cmd cmd/$$cmd/main.go; \
	done

fmt:
	go fmt ${PACKAGES}
	go fmt ${CMD_PACKAGES}

test:
	go test ${PACKAGES}

bench:
	go test -bench=. ${PACKAGES}

clean:
	go mod tidy
	@echo 'clean' ${BUILD} && rm -rf ${BUILD}

