args = `arg="$(filter-out $@,$(MAKECMDGOALS))" && echo $${arg:-${1}}`

VERSION=$(shell git describe --match 'v[0-9]*' --dirty='.m' --always)
REVISION=$(shell git rev-parse HEAD)$(shell if ! git diff --no-ext-diff --quiet --exit-code; then echo .m; fi)
DATETIME:=$(shell date)


PKG=local.pkg
EXTRA_LDFLAGS= -s -w
GO_LDFLAGS=-ldflags '-X $(PKG)/main.Version=$(VERSION) -X $(PKG)/main.Revision=$(REVISION) -X $(PKG)/main.Package=$(PACKAGE) -X $(PKG)/main.Date=$(DATE) $(EXTRA_LDFLAGS)'

all: main

.PHONY: main
main: 
	TEMPL_EXPERIMENT=rawgo templ generate
	go build -buildmode=plugin -o components.so components.go ExampleComponents_templ.go 

.PHONY: clean 
clean:
	/bin/rm main ./templcomponents/components.so

