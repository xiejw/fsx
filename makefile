include go.mk

# ------------------------------------------------------------------------------
# global configuraiton.
# ------------------------------------------------------------------------------
REPO      = fsx
LIB_DIR   = src
CMD_DIR   = cmd

LIBS      = github.com/xiejw/${REPO}/${LIB_DIR}/...
CMDP      = github.com/xiejw/${REPO}/${CMD_DIR}/...
CMDS      = $(patsubst ${CMD_DIR}/%,${BUILD_DIR}/%,$(wildcard ${CMD_DIR}/*))

# ------------------------------------------------------------------------------
# actions.
# ------------------------------------------------------------------------------
#
.PHONY: ${LIBS} compile compile_lib compile_cmd

compile: ${BUILD_DIR} compile_lib compile_cmd

compile_lib: ${BUILD_DIR} ${LIBS}

compile_cmd: ${BUILD_DIR} ${CMDS}

${LIBS}:
	${GO} $@

${BUILD_DIR}/%: cmd/%/main.go  # convention is cmd/<binary>/main.go
	${LD} $@ $<;

fmt:
	${FM} ${LIBS}
	${FM} ${CMDP}

test:
	${TS} ${LIBS}

clean:
	${CL} rm -rf ${BUILD_DIR} && go mod tidy

# ------------------------------------------------------------------------------
# binaries.
# ------------------------------------------------------------------------------

snapshot:
	@touch cmd/$@/main.go  # force a rebuild
	@make compile
	${EX} ${BUILD_DIR}/snapshot
