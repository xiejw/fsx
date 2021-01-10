REPO      = fsx
LIB_DIR   = src
CMD_DIR   = cmd
BUILD_DIR = .build

LIBS      = github.com/xiejw/${REPO}/${LIB_DIR}/...
CMD_LIBS  = github.com/xiejw/${REPO}/${CMD_DIR}/...
# convention is cmd/<binary>/main.go
CMD_CANDS = $(patsubst cmd/%,%,$(wildcard cmd/*))

compile: compile_lib compile_cmd

compile_lib:
	go build ${LIBS}

compile_cmd:
	@mkdir -p ${BUILD_DIR}
	@for cmd in ${CMD_CANDS}; do \
		echo 'compile cmd/'$$cmd && \
	  go build -o ${BUILD_DIR}/$$cmd cmd/$$cmd/main.go; \
	done

fmt:
	go fmt ${LIBS}
	go fmt ${CMD_LIBS}

test:
	go test ${LIBS}

clean:
	go mod tidy
	@echo "clean '"${BUILD_DIR}"'" && rm -rf ${BUILD_DIR}

