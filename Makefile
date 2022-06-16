GO ?= go
SOURCE_CMD_SERVER = gocalcharger/cmd/server
SOURCE_CMD_CLIENT = gocalcharger/cmd/client
SOURCE_CMD_WEB    = gocalcharger/cmd/web
TARGET_CMD_SERVER = GocalChargerServer
TARGET_CMD_CLIENT = GocalChargerClient
TARGET_CMD_WEB    = GocalChargerWeb

GO_FLAGS=-O2 -fPIC -fstack-protector-all -D_FORTIFY_SOURCE=2
GO_LD_FLAGS="--extldflags=-Wl,-z,now,-z,relro,-z,noexecstack -s"
GO_BUILD = $(GO) build --buildmode=exe --buildmode=pie --trimpath \
			--ldflags $(GO_LD_FLAGS)

export GO111MODULE=on
export CGO_ENABLED=1
export CGO_CFLAGS=$(GO_FLAGS)
export CGO_CXXFLAGS=$(GO_FLAGS)

.PHONY: all
all: cmd


.PHONY: cmd
cmd: cmd/server cmd/client cmd/web

.PHONY: cmd/server
cmd/server:
	$(GO_BUILD) -o $(TARGET_CMD_SERVER) $(SOURCE_CMD_SERVER)

.PHONY: cmd/client
cmd/client:
	$(GO_BUILD) -o $(TARGET_CMD_CLIENT) $(SOURCE_CMD_CLIENT)

.PHONY: cmd/web
cmd/web:
	$(GO_BUILD) -o $(TARGET_CMD_WEB) $(SOURCE_CMD_WEB)

.PHONY: clean
clean:
	$(RM) $(TARGET_CMD_SERVER) $(TARGET_CMD_CLIENT)
