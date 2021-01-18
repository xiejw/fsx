# ------------------------------------------------------------------------------
# global configuraiton.
# ------------------------------------------------------------------------------
REPO      = fsx
LIB_DIR   = src
CMD_DIR   = cmd
BUILD_DIR = .build

LIBS      = github.com/xiejw/${REPO}/${LIB_DIR}/...
CMD_LIBS  = github.com/xiejw/${REPO}/${CMD_DIR}/...
CMD_CANDS = $(patsubst cmd/%,${BUILD_DIR}/%,$(wildcard cmd/*))

GO = go build
LD = go build -o
EX =
TS = go test
FM = go fmt

# ------------------------------------------------------------------------------
# color printing.
# ------------------------------------------------------------------------------

# enable verbose cmd by `make V=1`
ifndef V

CCCOLOR   = "\033[34m"
LINKCOLOR = "\033[34;1m"
SRCCOLOR  = "\033[33m"
BINCOLOR  = "\033[36;1m"
ENDCOLOR  = "\033[0m"

GO        := @sh -c 'printf "    %b %b\n" $(CCCOLOR)GO$(ENDCOLOR) $(SRCCOLOR)"$$1"$(ENDCOLOR) 1>&2; ${GO} "$$1"' sh
LD        := @sh -c 'printf "    %b %b\n" $(LINKCOLOR)LD$(ENDCOLOR) $(BINCOLOR)`basename $$1`$(ENDCOLOR) 1>&2; ${LD} $$*' sh
EX        := @sh -c 'printf "    %b %b\n\n" $(LINKCOLOR)EX$(ENDCOLOR) $(BINCOLOR)"`basename $$1`"$(ENDCOLOR) 1>&2; ${EX} "$$1"' sh
TS        := @sh -c 'printf "    %b %b\n\n" $(LINKCOLOR)TS $(ENDCOLOR)$(BINCOLOR)"$$1"$(ENDCOLOR) 1>&2; ${TS} "$$1"' sh
FM        := @sh -c 'printf "    %b %b\n" $(LINKCOLOR)FM $(ENDCOLOR)$(BINCOLOR)"$$1"$(ENDCOLOR) 1>&2; ${FM} "$$1"' sh
endif

# ------------------------------------------------------------------------------
# actions.
# ------------------------------------------------------------------------------

compile: ${BUILD_DIR} compile_lib compile_cmd

${BUILD_DIR}:
	@mkdir -p ${BUILD_DIR}

compile_lib: ${LIBS}

compile_cmd: ${BUILD_DIR} ${CMD_CANDS}

${LIBS}:
	${GO} ${LIBS}

${BUILD_DIR}/%: cmd/%/main.go  # convention is cmd/<binary>/main.go
	${LD} $@ $<;

fmt:
	${FM} ${LIBS}
	${FM} ${CMD_LIBS}

test:
	${TS} ${LIBS}

clean:
	@echo 'cleaning...' && rm -rf ${BUILD_DIR}

snapshot: compile_cmd
	${EX} ${BUILD_DIR}/snapshot
