# ------------------------------------------------------------------------------
# global configuraiton.
# ------------------------------------------------------------------------------
REPO      = fsx
LIB_DIR   = src
CMD_DIR   = cmd
BUILD_DIR = .build

# ------------------------------------------------------------------------------
# color printing.
# ------------------------------------------------------------------------------

GO    = ${QUIET_GO} go build
EX    = ${QUIET_EX}
FM    = go fmt

CCCOLOR   = "\033[34m"
LINKCOLOR = "\033[34;1m"
SRCCOLOR  = "\033[33m"
BINCOLOR  = "\033[36;1m"
ENDCOLOR  = "\033[0m"

# enable verbose cmd by `make V=1`
ifndef V
QUIET_GO  = @printf '    %b %b\n' $(CCCOLOR)GO$(ENDCOLOR) \
          $(SRCCOLOR)$@$(ENDCOLOR) 1>&2;
QUIET_EX  = @printf '    %b %b\n' $(LINKCOLOR)EX$(ENDCOLOR) \
          $(BINCOLOR)$@$(ENDCOLOR) 1>&2;
FM        := @sh -c 'printf "    %b %b\n" $(LINKCOLOR)FM $(ENDCOLOR)$(BINCOLOR)"$$1"$(ENDCOLOR) 1>&2; ${FM} "$$1"' sh
endif

LIBS      = github.com/xiejw/${REPO}/${LIB_DIR}/...
CMD_LIBS  = github.com/xiejw/${REPO}/${CMD_DIR}/...
# convention is cmd/<binary>/main.go
CMD_CANDS = $(patsubst cmd/%,%,$(wildcard cmd/*))

compile: ${BUILD_DIR} compile_lib compile_cmd

${BUILD_DIR}:
	@mkdir -p ${BUILD_DIR}

compile_lib: ${LIBS}

${LIBS}:
	${GO} ${LIBS}

compile_cmd: ${CMD_CANDS}

	#@for cmd in ${CMD_CANDS}; do \
	#	echo 'compile cmd/'$$cmd && \
	#  ${GO} -o ${BUILD_DIR}/$$cmd cmd/$$cmd/main.go; \
	#done

fmt:
	${FM} ${LIBS}
	${FM} ${CMD_LIBS}

test:
	go test ${LIBS}

clean:
	go mod tidy
	@echo "clean '"${BUILD_DIR}"'" && rm -rf ${BUILD_DIR}

