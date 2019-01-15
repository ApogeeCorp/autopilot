DOCKER_HUB_AUTOPILOT_TAG ?= latest
AUTOPILOT_IMG=$(DOCKER_HUB_REPO)/autopilot:$(DOCKER_HUB_AUTOPILOT_TAG)

ifndef PKGS
PKGS := $(shell go list ./... 2>&1 | grep -v 'github.com/libopenstorage/autopilot/vendor' | grep -v 'pkg/client/informers/externalversions' | grep -v versioned | grep -v 'pkg/apis/autopilot')
endif

ifeq ($(BUILD_TYPE),debug)
BUILDFLAGS += -gcflags "-N -l"
endif

BASE_DIR    := $(shell git rev-parse --show-toplevel)
BIN         :=$(BASE_DIR)/bin

LDFLAGS += "-s -w"
BUILD_OPTIONS := -ldflags=$(LDFLAGS)

.DEFAULT_GOAL=all
.PHONY: clean vendor vendor-update container deploy

all: autopilot vet lint staticcheck

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

$(GOPATH)/bin/staticcheck:
	go get -u honnef.co/go/tools/cmd/staticcheck

staticcheck: $(GOPATH)/bin/staticcheck
	$(GOPATH)/bin/staticcheck $(PKGS)

errcheck:
	go get -v github.com/kisielk/errcheck
	errcheck -verbose -blank $(PKGS)

pretest: lint vet errcheck staticcheck

codegen:
	@echo "Generating CRD"
	@hack/update-codegen.sh

autopilot:
	@echo "Building the autopilot binary"
	@cd cmd/autopilot && CGO_ENABLED=0 go build $(BUILD_OPTIONS) -o $(BIN)/autopilot

container:
	@echo "Building container: docker build --tag $(AUTOPILOT_IMG) -f Dockerfile ."
	sudo docker build --tag $(AUTOPILOT_IMG) -f Dockerfile .

deploy: container
	sudo docker push $(AUTOPILOT_IMG)

clean:
	go clean -i $(PKGS)

