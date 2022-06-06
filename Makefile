.PHONY: all clean build build-clean
.PHONY: fmt lint fieldalign fieldalign-fix

OUT_DIR = ./bin
CMD_DIR = ./cmd/stuffnotifier

OUTFILE = $(OUT_DIR)/main

GO = $(shell which go)
TEST_FLAGS := -count 1

define go_test
$(GO) test $(MODULE_DIRS) $(strip $(TEST_FLAGS) $1)
endef

define _buildcmd
@rm -rf $1
@mkdir -p ./bin
$(GO) build -o $1 $2
endef

define _fieldalign
@fieldalignment $(strip $2) $1
endef

all : clean build

build :
	$(call _buildcmd,$(OUTFILE),$(CMD_DIR))

build-clean : clean build

test-ci :
	$(call go_test, -coverprofile=profile.cov)

lint :
	@golangci-lint run --config=.golangci.yml

fieldalign :
	$(call _fieldalign,$(FILES))

fieldalign-fix :
	$(call _fieldalign,$(FILES),-fix)

fmt :
	gofmt -s -w ./

clean :
	@rm -rf ./bin