GO ?= go
SOURCE_GUI = gui/main.go
SOURCE_CMD_SERVER = cmd/server/server.go
SOURCE_CMD_CLIENT = cmd/client/client.go
TARGET_GUI = GocalChargerGui
TARGET_CMD_SERVER = GocalChargerServer
TARGET_CMD_CLIENT = GocalChargerClient

GO_FLAGS=-O2 -fPIC -fstack-protector-all -D_FORTIFY_SOURCE=2
GO_LD_FLAGS="--extldflags=-Wl,-z,now,-z,relro,-z,noexecstack -s"
GO_BUILD = $(GO) build --buildmode=exe --buildmode=pie --trimpath \
			--ldflags $(GO_LD_FLAGS)

export GO111MODULE=on
export CGO_ENABLED=1
export CGO_CFLAGS=$(GO_FLAGS)
export CGO_CXXFLAGS=$(GO_FLAGS)

.PHONY: all
all: gui cmd

.PHONY: gui
gui:
	$(GO_BUILD) -o $(TARGET_GUI) $(SOURCE_GUI)

.PHONY: cmd
cmd: cmd/server cmd/client

.PHONY: cmd/server
cmd/server:
	$(GO_BUILD) -o $(TARGET_CMD_SERVER) $(SOURCE_CMD_SERVER)

.PHONY: cmd/client
cmd/client:
	$(GO_BUILD) -o $(TARGET_CMD_CLIENT) $(SOURCE_CMD_CLIENT)

.PHONY: clean
clean:
	$(RM) $(TARGET_GUI) $(TARGET_CMD_SERVER) $(TARGET_CMD_CLIENT)
