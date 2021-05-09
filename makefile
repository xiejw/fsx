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
CBLUE   = "\033[34m"
CHBLUE  = "\033[34;1m"
CYELL   = "\033[33m"
CCYAN   = "\033[36;1m"
CEND    = "\033[0m"

GO      := @sh -c 'printf "    %b %b\n"   $(CBLUE)GO$(CEND)  $(CYELL)"$$1"$(CEND)            1>&2; ${GO} "$$1"' sh
LD      := @sh -c 'printf "    %b %b\n"   $(CHBLUE)LD$(CEND) $(CCYAN)`basename $$1`$(CEND)   1>&2; ${LD}  $$* ' sh
EX      := @sh -c 'printf "    %b %b\n\n" $(CHBLUE)EX$(CEND) $(CCYAN)"`basename $$1`"$(CEND) 1>&2; ${EX} "$$1"' sh
TS      := @sh -c 'printf "    %b %b\n\n" $(CHBLUE)TS$(CEND) $(CCYAN)"$$1"$(CEND)            1>&2; ${TS} "$$1"' sh
FM      := @sh -c 'printf "    %b %b\n"   $(CHBLUE)FM$(CEND) $(CCYAN)"$$1"$(CEND)            1>&2; ${FM} "$$1"' sh
endif

# ------------------------------------------------------------------------------
# actions.
# ------------------------------------------------------------------------------

compile: ${BUILD_DIR} compile_lib compile_cmd

compile_lib: ${BUILD_DIR} ${LIBS}

compile_cmd: ${BUILD_DIR} ${CMD_CANDS}

${BUILD_DIR}:
	@mkdir -p ${BUILD_DIR}

.PHONY: ${LIBS}

${LIBS}:
	${GO} $@

${BUILD_DIR}/%: cmd/%/main.go  # convention is cmd/<binary>/main.go
	${LD} $@ $<;

fmt:
	${FM} ${LIBS}
	${FM} ${CMD_LIBS}

test:
	${TS} ${LIBS}

clean:
	@echo 'cleaning...' && rm -rf ${BUILD_DIR} && go mod tidy

snapshot: compile_cmd
	${EX} ${BUILD_DIR}/snapshot
