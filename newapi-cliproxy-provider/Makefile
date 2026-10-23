# 默认构建当前平台；可通过命令行覆盖 Go 工具、目标平台和输出目录。
.DEFAULT_GOAL := build
GO ?= go
# Docker 目标无需本机安装 Go；其他目标继续自动识别本机平台。
ifneq ($(filter linux-amd64 linux-amd64-debug,$(MAKECMDGOALS)),)
GOOS ?= linux
GOARCH ?= amd64
endif
ifndef GOOS
GOOS := $(shell $(GO) env GOOS)
endif
ifndef GOARCH
GOARCH := $(shell $(GO) env GOARCH)
endif

# 插件通过 CGO 导出原生接口，跨平台构建时需自行配置对应的 C 编译器。
export CGO_ENABLED := 1
export GOOS
export GOARCH

# 保持与 README 中的插件名称及平台目录一致。
PLUGIN_ID := newapi-provider
LIB_EXT := so
ifeq ($(GOOS),windows)
LIB_EXT := dll
else ifeq ($(GOOS),darwin)
LIB_EXT := dylib
endif
OUTPUT_DIR ?= plugins/$(GOOS)/$(GOARCH)
OUTPUT := $(OUTPUT_DIR)/$(PLUGIN_ID).$(LIB_EXT)

# Docker 构建固定使用本地引擎，并从 go.mod 读取镜像的 Go 版本。
DOCKER ?= docker
PROJECT_DIR := $(CURDIR)
# MSYS2 的 /d/... 路径需显式转换，供 Windows Docker 客户端挂载宿主目录。
ifeq ($(OS),Windows_NT)
DOCKER_PROJECT_DIR = $(shell cygpath -m "$(PROJECT_DIR)")
else
DOCKER_PROJECT_DIR = $(PROJECT_DIR)
endif
GO_VERSION = $(shell awk '$$1 == "go" { print $$2; exit }' "$(PROJECT_DIR)/go.mod")
DOCKER_IMAGE = golang:$(GO_VERSION)-bookworm

.PHONY: build debug linux-amd64 linux-amd64-debug

# 普通构建使用空的调试日志实现。
build:
	$(GO) build -trimpath -buildmode=c-shared -o "$(OUTPUT)" .

# 启用现有 debug 日志；输出路径与普通构建相同，产物会覆盖。
debug:
	$(GO) build -trimpath -tags debug -buildmode=c-shared -o "$(OUTPUT)" .

# Linux/amd64 普通构建和 debug 构建复用同一容器流程。
linux-amd64: DOCKER_BUILD_TARGET := build
linux-amd64-debug: DOCKER_BUILD_TARGET := debug
# 禁用 MSYS 路径转换；工作目录使用双斜杠，兼容混用 Git Bash 与 MSYS2 make 的环境。
linux-amd64 linux-amd64-debug: export MSYS_NO_PATHCONV := 1
linux-amd64 linux-amd64-debug: export MSYS2_ARG_CONV_EXCL := *
linux-amd64 linux-amd64-debug:
	MSYS_NO_PATHCONV=1 MSYS2_ARG_CONV_EXCL='*' $(DOCKER) --context default run --rm --platform linux/amd64 \
		--mount "type=bind,source=$(DOCKER_PROJECT_DIR),target=/src" \
		--mount type=volume,source=newapi-provider-go-mod-cache,target=//go/pkg/mod \
		--mount type=volume,source=newapi-provider-go-build-cache,target=//root/.cache/go-build \
		--workdir //src \
		--env CGO_ENABLED=1 \
		--env GOOS=linux \
		--env GOARCH=amd64 \
		"$(DOCKER_IMAGE)" \
		make $(DOCKER_BUILD_TARGET) GOOS=linux GOARCH=amd64
