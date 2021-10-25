SHELL=/bin/bash -o pipefail

VERSION=0.10.0-RC-1
TARGETS=linux/amd64 windows/amd64 darwin/amd64 linux/s390x linux/ppc64le  linux/arm64


GO ?= go

GIT_REV := $(shell git rev-parse HEAD 2> /dev/null || true)
GIT_COMMIT := $(if $(shell git status --porcelain --untracked-files=no),${GIT_REV}-dirty,${GIT_REV})
GIT_BRANCH ?= $(shell git rev-parse --abbrev-ref HEAD 2>/dev/null)

export GO111MODULE=on
PWD=$(shell pwd)
LDFLAGS := -s -X github.com/codenotary/cas/pkg/meta.version=v${VERSION} \
			  -X github.com/codenotary/cas/pkg/meta.gitCommit=${GIT_COMMIT} \
			  -X github.com/codenotary/cas/pkg/meta.gitBranch=${GIT_BRANCH}
LDFLAGS_STATIC := ${LDFLAGS} \
				  -X github.com/codenotary/cas/pkg/meta.static=static \
				  -extldflags "-static"
TEST_FLAGS ?= -v -race
CASEXE=cas-v${VERSION}-windows-amd64.exe
SETUPEXE=codenotary_cas_v${VERSION}_setup.exe

.PHONY: cas
cas:
	$(GO) build -ldflags '${LDFLAGS} -X github.com/codenotary/cas/pkg/meta.version=v${VERSION}-dev' ./cmd/cas

.PHONY: vendor
vendor:
	$(GO) mod vendor

.PHONY: test
test:
	$(GO) vet ./...
	$(GO) test ${TEST_FLAGS} ./...

.PHONY: install
install: TEST_FLAGS=-v
install: vendor test
	$(GO) install -ldflags '${LDFLAGS}' ./cmd/cas

.PHONY: static
static:
	$(GO) build -a -tags netgo -ldflags '${LDFLAGS_STATIC}' ./cmd/cas

.PHONY: docs/cmd
docs/cmd:
	rm -rf docs/cmd/*.md
	$(GO) run docs/cmd/main.go


.PHONY: clean/dist
clean/dist:
	rm -Rf ./dist

.PHONY: clean
clean: clean/dist
	rm -f ./cas

.PHONY: CHANGELOG.md
CHANGELOG.md:
	git-chglog -o CHANGELOG.md

.PHONY: CHANGELOG.md.next-tag
CHANGELOG.md.next-tag:
	git-chglog -o CHANGELOG.md --next-tag v${VERSION}

.PHONY: dist
dist: clean/dist
	@$(GO) build -a -ldflags "${LDFLAGS_STATIC}" \
    -o ./dist/cas-v${VERSION}-linux-amd64-static \
    ./cmd/cas ;\
	for p in ${TARGETS}; do \
		platform_split=($${p//\// }) ; \
		GOOS=$${platform_split[0]} ; \
		GOARCH=$${platform_split[1]} ; \
		output_name='cas-v'${VERSION}'-'$$GOOS'-'$$GOARCH ; \
		if [ "$$GOOS" = "windows" ]; then \
			output_name+='.exe' ; \
		fi ;\
		GOOS="$$GOOS" GOARCH="$$GOARCH" $(GO) build -a -ldflags "${LDFLAGS}" -o ./dist/$$output_name ./cmd/cas || { echo "failed compiling for $$GOOS/$$GOARCH"; exit 2; } ; \
	done

.PHONY: dist/all
dist/all: dist

.PHONY: dist/binary.md
dist/binary.md:
	@for f in ./dist/*; do \
		ff=$$(basename $$f); \
		shm_id=$$(sha256sum $$f | awk '{print $$1}'); \
		printf "[$$ff](https://github.com/codenotary/cas/releases/download/v${VERSION}/$$ff) | $$shm_id \n" ; \
	done

.PHONY: build/codegen
build/codegen:
	protoc -I pkg/lcgrpc pkg/lcgrpc/lc.proto  \
	-I${GOPATH}/pkg/mod \
	-I${GOPATH}/pkg/mod/github.com/codenotary/immudb@v0.8.0/pkg/api/schema \
	-I${GOPATH}/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.14.4/third_party/googleapis \
	-I${GOPATH}/pkg/mod/github.com/dgraph-io/badger/v2@v2.0.0-20200408100755-2e708d968e94 \
	-I${GOPATH}/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.14.4 \
	--go_out=pkg/lcgrpc --go-grpc_out=pkg/lcgrpc \
	--go_opt=paths=source_relative \
	--go-grpc_opt=paths=source_relative
