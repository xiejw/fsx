# ------------------------------------------------------------------------------
# compiler tooling.
# ------------------------------------------------------------------------------

GO = go build
LD = go build -o
EX =
TS = go test
FM = go fmt

CL = @echo ${CCYAN}'!!! cleaning...'${CEND} &&

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
# common actions.
# ------------------------------------------------------------------------------

BUILD_DIR = .build

${BUILD_DIR}:
	@mkdir -p ${BUILD_DIR}

