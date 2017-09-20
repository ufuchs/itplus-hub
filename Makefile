GO ?= go
GIT ?= git

export GOOS=linux
#export GOARCH=arm
export GOARCH=amd64

#export GOARM=7


ifeq ($(OS),Windows_NT)
    CCFLAGS += -D WIN32
    GO = /d/Go/bin/go
    ifeq ($(PROCESSOR_ARCHITECTURE),AMD64)
        CCFLAGS += -D AMD64
    endif
    ifeq ($(PROCESSOR_ARCHITECTURE),x86)
        CCFLAGS += -D IA32
    endif
else
    UNAME_S := $(shell uname -s)
    ifeq ($(UNAME_S),Linux)
        CCFLAGS += -D LINUX
    endif
    ifeq ($(UNAME_S),Darwin)
		SHA1 = shasum
    endif
    UNAME_P := $(shell uname -p)
    ifeq ($(UNAME_P),x86_64)
        CCFLAGS += -D AMD64
    endif
    ifneq ($(filter %86,$(UNAME_P)),)
        CCFLAGS += -D IA32
    endif
    ifneq ($(filter arm%,$(UNAME_P)),)
        CCFLAGS += -D ARM
    endif
endif

#FTP_HOST := ftp://ftp.hidrive.strato.com/public
FTP_HOST := ftp://localhost/ftp

#-------------------------------------------------------------------------------
APP := itplus-hub
VERSION := 0.10
REVL := $(shell git rev-parse HEAD | tail -c 8)

REVH := $(shell git rev-parse HEAD | head -c 7)
REV := $(REVH)..$(REVL)
ZIPFILE = $(APP).zip
VER_STRING := $(VERSION) (build $(REV))
#-------------------------------------------------------------------------------
#GOPATH := $(CURDIR)/_vendor:$(GOPATH)
#-------------------------------------------------------------------------------
all: build

version:
	@echo 'package main' > version.go
	@echo 'var (' >> version.go
	@echo '    version = "$(VER_STRING)"' >> version.go
	@echo '    buildDate = "$(shell date)";' >> version.go
	@echo '    builder = "$(LOGNAME)@$(shell hostname)"' >> version.go
	@echo ')' >> version.go

upload: build
	@echo 'version: $(APP) version $(VER_STRING)' > INFO
	@echo 'sha1: $(shell $(SHA1) $(APP) | cut -d" " -f 1)' >> INFO
	zip $(ZIPFILE) $(APP)
	curl -T $(ZIPFILE) -u ufuchs:$(PASSWORD) $(FTP_HOST)/$(APP)/$(ZIPFILE)
	curl -T INFO -u ufuchs:$(PASSWORD) $(FTP_HOST)/$(APP)/INFO

build: version
#	golint ./
	@$(GO) build -ldflags "-s -w" -o $(APP)

