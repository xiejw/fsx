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
LD    = ${QUIET_LD} go build
EX    = ${QUIET_EX}
TS    = go test
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
QUIET_LD  = @printf '    %b %b\n' $(LINKCOLOR)LD$(ENDCOLOR) \
          $(BINCOLOR)`basename $@`$(ENDCOLOR) 1>&2;
QUIET_EX  = @printf '    %b %b\n' $(LINKCOLOR)EX$(ENDCOLOR) \
          $(BINCOLOR)$@$(ENDCOLOR) 1>&2;

TS        := @sh -c 'printf "    %b %b\n" $(LINKCOLOR)TS $(ENDCOLOR)$(BINCOLOR)"$$1"$(ENDCOLOR) 1>&2; ${TS} "$$1"' sh
FM        := @sh -c 'printf "    %b %b\n" $(LINKCOLOR)FM $(ENDCOLOR)$(BINCOLOR)"$$1"$(ENDCOLOR) 1>&2; ${FM} "$$1"' sh
endif

LIBS      = github.com/xiejw/${REPO}/${LIB_DIR}/...
CMD_LIBS  = github.com/xiejw/${REPO}/${CMD_DIR}/...
CMD_CANDS = $(patsubst cmd/%,${BUILD_DIR}/%,$(wildcard cmd/*))

compile: ${BUILD_DIR} compile_lib compile_cmd

${BUILD_DIR}:
	@mkdir -p ${BUILD_DIR}

compile_lib: ${LIBS}

${LIBS}:
	${GO} ${LIBS}

compile_cmd: ${BUILD_DIR} ${CMD_CANDS}

# convention is cmd/<binary>/main.go
${BUILD_DIR}/%: cmd/%/main.go
	${LD} -o $@ $<;

fmt:
	${FM} ${LIBS}
	${FM} ${CMD_LIBS}

test:
	${TS} ${LIBS}

clean:
	#go mod tidy
	rm -rf ${BUILD_DIR}

