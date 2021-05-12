include go.mk

# ------------------------------------------------------------------------------
# global configuraiton.
# ------------------------------------------------------------------------------
REPO      = fsx
LIB_DIR   = src
CMD_DIR   = cmd

LIBS      = github.com/xiejw/${REPO}/${LIB_DIR}/...
CMDP      = github.com/xiejw/${REPO}/${CMD_DIR}/...
CMDS      = $(patsubst ${CMD_DIR}/%,%,$(wildcard ${CMD_DIR}/*))
CMDB      = $(patsubst ${CMD_DIR}/%,${BUILD_DIR}/%,$(wildcard ${CMD_DIR}/*))

# ------------------------------------------------------------------------------
# actions.
# ------------------------------------------------------------------------------
#
.PHONY: compile ${LIBS}

compile: ${BUILD_DIR} ${LIBS} ${CMDB}

${LIBS}:
	${GO} $@

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

$(foreach cmd,$(CMDS),$(eval $(call objs,$(cmd),$(BUILD_DIR))))

