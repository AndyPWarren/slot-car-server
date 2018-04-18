# Pi environment varibles
_PI_USER := $(if $(PI_USER),$(PI_USER),$(USER))
_PI_SSH_KEY := $(if $(PI_SSH_KEY),$(PI_SSH_KEY),'~/.ssh/id_rsa')
# Compilation Flags
GOOS            ?=
GOARCH          ?=
# Flags
FLAGS           ?=
# LDFlags
BUILD_TIME_LDFLAG       ?= -X main.timestamp=$(shell date +%s)
BUILD_VERSION_LDFLAG    ?= -X main.version=$(shell git rev-parse HEAD)
LDFLAGS                 ?= "$(BUILD_TIME_LDFLAG) $(BUILD_VERSION_LDFLAG)"
# Binary Name
BIN_NAME        ?= slot-car-server
BIN_SUFFIX      ?=
ifneq ($(GOOS),)
ifneq ($(GOARCH),)
BIN_SUFFIX      = .$(GOOS)-$(GOARCH)
endif
endif

# Packages to test
TEST_PKGS=$(shell go list ./... | grep -v ./vendor/)

.PHONY: build testdata

# Run the go application
run:
	go run -ldflags $(LDFLAGS) main.go $(FLAGS)

# Build a binary
build:
	GOOS=$(GOOS) \
	GOARCH=$(GOARCH) \
	go build \
		-v \
		-ldflags $(LDFLAGS) \
		-o "bin/$(BIN_NAME)$(BIN_SUFFIX)"

# Build a binary for linux arm
linuxarm:
	GOOS=linux GOARCH=arm make build

# Build docker image
image:
	docker build --force-rm -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

# Run tests
test:
	go test -v -cover $(TEST_PKGS)

# Generate test data
testdata:
	go test -v $(TEST_PKGS) -update

deploy:
	scp -i $(_PI_SSH_KEY) ./bin/$(BIN_NAME).linux-arm $(_PI_USER)@raspberrypi:~

run:
	./bin/$(BIN_NAME)$(BIN_SUFFIX)

buildrun:
	make build
	make run

buildpi:
	make linuxarm
	make deploy
