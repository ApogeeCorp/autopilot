ifndef PKGS
PKGS := $(shell go list ./... 2>&1 | grep -v 'github.com/libopenstorage/autopilot/vendor' | grep -v 'pkg/client/informers/externalversions' | grep -v versioned | grep -v 'pkg/apis/autopilot')
endif

ifeq ($(BUILD_TYPE),debug)
BUILDFLAGS += -gcflags "-N -l"
endif

BASE_DIR    := $(shell git rev-parse --show-toplevel)

LDFLAGS += "-s -w"
BUILD_OPTIONS := -ldflags=$(LDFLAGS)

.DEFAULT_GOAL=all
.PHONY: clean vendor vendor-update

all: vet lint simple

vendor-update:
	dep ensure -update

vendor:
	dep ensure

lint:
	go get -v golang.org/x/lint/golint
	for file in $$(find . -name '*.go' | grep -v vendor | \
                                       grep -v '\.pb\.go' | \
                                       grep -v '\.pb\.gw\.go' | \
                                       grep -v 'externalversions' | \
                                       grep -v 'versioned' | \
                                       grep -v 'generated'); do \
		golint $${file}; \
		if [ -n "$$(golint $${file})" ]; then \
			exit 1; \
		fi; \
	done

vet:
	go vet $(PKGS)

$(GOPATH)/bin/gosimple:
	go get -u honnef.co/go/tools/cmd/gosimple

simple: $(GOPATH)/bin/gosimple
	$(GOPATH)/bin/gosimple $(PKGS)

errcheck:
	go get -v github.com/kisielk/errcheck
	errcheck -verbose -blank $(PKGS)

pretest: lint vet errcheck simple

codegen:
	@echo "Generating CRD"
	@hack/update-codegen.sh

clean:
	go clean -i $(PKGS)

