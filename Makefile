# http://alexyu.se/content/2014/04/simple-makefile-go
#
#  Makefile for Go
#
GO_CMD=go
GO_BUILD=$(GO_CMD) build -o $(OUTPUT_NAME)
GO_BUILD_RACE=$(GO_CMD) build -race
GO_TEST=$(GO_CMD) test
GO_TEST_VERBOSE=$(GO_CMD) test -v
GO_INSTALL=$(GO_CMD) install -v
GO_CLEAN=$(GO_CMD) clean
GO_DEPS=$(GO_CMD) get -d -v
GO_DEPS_UPDATE=$(GO_CMD) get -d -v -u
GO_VET=$(GO_CMD) vet
GO_FMT=$(GO_CMD) fmt
GO_LINT=golint

# Color output
CRESET=\x1b[0m
CRED=\x1b[31;01m
CGREEN=\x1b[32;01m
CYELLOW=\x1b[33;01m
CBLUE=\x1b[34;01m
CMAGENTA=\x1b[35;01m
CCYAN=\x1b[36;01m

# Packages
GITHUB_USER=zchee
TOP_PACKAGE_DIR := github.com/$(GITHUB_USER)
# Current 
PACKAGE_LIST := `basename $(PWD)`
# Parse "func main()" word only .go in current dir. Required ag
# FIXME: Not support main.go
# OUTPUT_NAME := `ag --go -l --depth 1 "func main\(\)" | sed -e 's/\.go//g'`
# Get output name no dependency
OUTPUT_NAME := `go list | awk -F "/" '{print $NF}'`

# go get -u -v github.com/jstemmer/gotags
CTAGS_CMD=gotags
# gotags options
#		-f string
#			Output specified file name
#			If specified "-", output to stdout
#		-R
#			Recurse into directories in the file list
#		-fields string
#			Include selected extension fields (only +l)
#		-sort
#			Sort tags (default true)
#		-tag-relative
#			File path s should be relative to the directory containing the tag file
CTAGS_OPTIONS=-f tags -R -fields=+l -sort -tag-relative .

.PHONY: all build build-race test test-verbose deps update-deps install clean fmt vet lint

all: build

debug:
		# defined command list
		@echo GO_CMD=$(GO_CMD)
		@echo GO_BUILD=$(GO_BUILD)
		@echo GO_BUILD_RACE=$(GO_BUILD_RACE)
		@echo GO_TEST=$(GO_TEST)
		@echo GO_TEST_VERBOSE=$(GO_TEST_VERBOSE)
		@echo GO_INSTALL=$(GO_INSTALL)
		@echo GO_CLEAN=$(GO_CLEAN)
		@echo GO_DEPS=$(GO_DEPS)
		@echo GO_DEPS_UPDATE=$(GO_DEPS_UPDATE)
		@echo GO_VET=$(GO_VET)
		@echo GO_FMT=$(GO_FMT)
		@echo GO_LINT=$(GO_LINT)

build: vet
	@for p in $(PACKAGE_LIST); do \
		echo "$(CBLUE)==>$(CRESET) Build $(CGREEN)$$p$(CRESET) ..."; \
		$(GO_BUILD) $(TOP_PACKAGE_DIR)/$$p || exit 1; \
	done

build-race: vet
	@for p in $(PACKAGE_LIST); do \
		echo "$(CBLUE)==>$(CRESET) Build $(CGREEN)$$p$(CRESET) ...h -race flag..."; \
		time $(GO_BUILD_RACE) $(TOP_PACKAGE_DIR)/$$p || exit 1; \
	done

build-force: vet
	@for p in $(PACKAGE_LIST); do \
		echo "$(CBLUE)==>$(CRESET) Build $(CGREEN)$$p$(CRESET) ..."; \
		time $(GO_BUILD) -a $(TOP_PACKAGE_DIR)/$$p || exit 1; \
	done

build-verbose: vet
	@for p in $(PACKAGE_LIST); do \
		echo "$(CBLUE)==>$(CRESET) Build $(CGREEN)$$p$(CRESET) ..."; \
		time $(GO_BUILD) -v -x $(TOP_PACKAGE_DIR)/$$p || exit 1; \
	done

ctags:
	@for p in $(PACKAGE_LIST); do \
		echo "$(CBLUE)==>$(CRESET) Create ctags file..."; \
		$(CTAGS_CMD) $(CTAGS_OPTIONS) || exit 1; \
	done

test: deps
	@for p in $(PACKAGE_LIST); do \
		echo "$(CBLUE)==>$(CRESET) Unit Testing $(CGREEN)$$p$(CRESET) ..."; \
		$(GO_TEST) $(TOP_PACKAGE_DIR)/$$p || exit 1; \
	done

test-verbose: deps
	@for p in $(PACKAGE_LIST); do \
		echo "$(CBLUE)==>$(CRESET) Unit Testing $(CGREEN)$$p$(CRESET) ..."; \
		$(GO_TEST_VERBOSE) $(TOP_PACKAGE_DIR)/$$p || exit 1; \
	done

deps:
	@for p in $(PACKAGE_LIST); do \
		echo "$(CBLUE)==>$(CRESET) Install dependencies for $(CGREEN)$$p$(CRESET) ..."; \
		$(GO_DEPS) $(TOP_PACKAGE_DIR)/$$p || exit 1; \
	done

update-deps:
	@for p in $(PACKAGE_LIST); do \
		echo "$(CBLUE)==>$(CRESET) Update dependencies for $(CGREEN)$$p$(CRESET) ..."; \
		$(GO_DEPS_UPDATE) $(TOP_PACKAGE_DIR)/$$p || exit 1; \
	done

install:
	@for p in $(PACKAGE_LIST); do \
		echo "$(CBLUE)==>$(CRESET) Install $(CGREEN)$$p$(CRESET) ..."; \
		$(GO_INSTALL) $(TOP_PACKAGE_DIR)/$$p || exit 1; \
	done

clean:
	@for p in $(PACKAGE_LIST); do \
		echo "$(CBLUE)==>$(CRESET) Clean $(CGREEN)$$p$(CRESET) ..."; \
		$(GO_CLEAN) $(TOP_PACKAGE_DIR)/$$p; \
	done

fmt:
	@for p in $(PACKAGE_LIST); do \
		echo "$(CBLUE)==>$(CRESET) Formatting $(CGREEN)$$p$(CRESET) ..."; \
		$(GO_FMT) $(TOP_PACKAGE_DIR)/$$p || exit 1; \
	done

vet:
	@for p in $(PACKAGE_LIST); do \
		echo "$(CBLUE)==>$(CRESET) Vet $(CGREEN)$$p$(CRESET) ..."; \
		$(GO_VET) $(TOP_PACKAGE_DIR)/$$p; \
	done

lint:
	@for p in $(PACKAGE_LIST); do \
		echo "$(CBLUE)==>$(CRESET) Lint $(CGREEN)$$p$(CRESET) ..."; \
		$(GO_LINT) src/$(TOP_PACKAGE_DIR)/$$p; \
	done
