# A simple makefile template for use with Webstorm Templates on golang projects
NAME=$(notdir $(shell pwd))

VERSION=0.1.0
BUILD_TIME=$(shell date +%FT%T%z)

GITHUB_NAME=$(shell git config --list|grep -e "user.name="|cut -c 11-|tr -d ' ')

BUILD_DEBUG=debug
BUILD_RELEASE=release
BUILD_TEST=test

ifndef $(build)
build=$(BUILD_DEBUG)
endif

MAIN=main

GOMAIN=$(MAIN).go


LDFLAGS=-ldflags '-X $(MAIN).BUILD_VERSION=$(VERSION)\
 -X $(MAIN).BUILD_TIME=$(BUILD_TIME)\
 -X $(MAIN).BUILD_TYPE=$(build)\
 -X $(MAIN).BUILD_NAME=$(NAME)'

GO=go
GODBG=dlv

GOFMT=gofmt

GOBUILD=$(GO) build
GORUN=$(GO) run

GOINSTALL=$(GO) install
GOCLEAN=$(GO) clean

GOTEST=$(GO) test
GOCOVER=$(GO) test -cover -html=cover_$(BUILD_TIME).html

.PHONY: $(NAME)
$(NAME): dependencies
	@printf "\n\e[32mBuilding\e[0m '\e[34m%s\e[0m'...\n" $@
	@printf "\tExported to '\e[34m%s\e[0m':\n" $(GOMAIN)
	@printf "\t\t\e[35mBUILD_NAME: \e[31m%s\n" $(NAME)
	@printf "\t\t\e[35mBUILD_VERSION: \e[31m%s\n" $(VERSION)
	@printf "\t\t\e[35mBUILD_TYPE: \e[31m%s\n" $(build)
	@printf "\t\t\e[35mBUILD_TIME: \e[31m%s\e[0m\n\n" $(BUILD_TIME)
	@$(GOBUILD) $(LDFLAGS) -o $(NAME) $(GOMAIN)
	@printf "\e[32mDone!\e[0m\n"

.PHONY: dependencies
dependencies:
	@printf "\n\e[32mInstalling\e[0m '\e[34m%s\e[0m'...\n" $@
	@$(GO) get
	@printf "\e[32mDone!\e[0m\n"

.PHONY: debug
debug: $(NAME)
	@printf "\n\e[32mDebugging\e[0m '\e[34m$<\e[0m' with\e[35m%s\e[0m...\n" $(GODBG)
	@printf "\t\t\e[35mBUILD_TYPE: \e[31m%s\e[0m\n\n" $(build)
	@$(GODBG) exec $<
	@printf "\e[32mDone!\e[0m\n"

.PHONY: install
install: $(GOMAIN) $(NAME) test
	@printf "\n\e[32mInstalling\e[0m '\e[34m%s\e[0m'...\n\t" $(NAME)
	$(GOINSTALL) $(LDFLAGS)

.PHONY: run
run: $(GOMAIN) $(NAME)
	@printf "\n\e[32mRunning\e[0m '\e[34m%s\e[0m'...\n\n" $(NAME)
	@$(GORUN) $(LDFLAGS) $(GOMAIN)

.PHONY: test
test: test/ | $(NAME)
	@printf "\n\e[32mTesting\e[0m '\e[34m%s\e[0m'...\n" $(NAME)
	@build=$(BUILD_TEST)
	@printf "\t\t\e[35mBUILD_TYPE: \e[31m%s\e[0m\n\n" $(build)
	@$(GOTEST) $(LDFLAGS) github.com/$(GITHUB_NAME)/$(notdir $(shell pwd))

.PHONY: cover
cover: test/ | $(NAME)
	@printf "\n\e[32mChecking Coverage\e[0m of '\e[34m%s\e[0m'...\n" $(NAME)
	@build=$(BUILD_TEST)
	@printf "\t\t\e[35mBUILD_TYPE: \e[31m%s\e[0m\n\n" $(build)
	@$(GOCOVER) $(LDFLAGS)

.PHONY: format
format: $(shell find . -type f -name '*.go')
	@printf "\nFormatting '%s'...\n" $^
	@$(GOFMT) -w $^

.PHONY: clean
clean:
	$(GOCLEAN)
