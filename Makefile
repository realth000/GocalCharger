GO ?= go
SOURCE_GUI = gui/main.go
TARGET_GUI = GocalCharger

GO_FLAGS=-O2 -fPIC -fstack-protector-all -D_FORTIFY_SOURCE=2
GO_LD_FLAGS="--extldflags=-Wl,-z,now,-z,relro,-z,noexecstack -s"
GO_BUILD = $(GO) build --buildmode=exe --trimpath \
			--ldflags $(GO_LD_FLAGS)

export GO111MODULE=on
export CGO_ENABLED=1
export CGO_CFLAGS=$(GO_FLAGS)
export CGO_CXXFLAGS=$(GO_FLAGS)

.PHONY: all
all: gui

.PHONY: gui
gui:
	$(GO_BUILD) --buildmode=pie -o $(TARGET_GUI) $(SOURCE_GUI)

.PHONY: clean
clean:
	$(RM) $(TARGET_GUI)
