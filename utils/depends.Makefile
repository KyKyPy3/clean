define github_url
    https://github.com/$(GITHUB)/releases/download/v$(VERSION)/$(ARCHIVE)
endef

# creates a directory bin.
bin:
	@ mkdir -p $@

# ~~ [migrate] ~~~ https://github.com/golang-migrate/migrate ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

MIGRATE := $(shell command -v migrate || echo "bin/migrate")
migrate: bin/migrate ## Install migrate (database migration)

bin/migrate: VERSION := 4.18.1
bin/migrate: GITHUB  := golang-migrate/migrate
bin/migrate: ARCHIVE := migrate.$(OSTYPE)-${ARCH}.tar.gz
bin/migrate: bin
	@ printf "Install migrate... "
	@ curl -Ls $(call github_url) | tar -zOxf - ./migrate > $@ && chmod +x $@
	@ echo "done."

# ~~ [ gotestsum ] ~~~ https://github.com/gotestyourself/gotestsum ~~~~~~~~~~~~~~~~~~~~~~~

GOTESTSUM := $(shell command -v gotestsum || echo "bin/gotestsum")
gotestsum: bin/gotestsum ## Installs gotestsum (testing go code)

bin/gotestsum: VERSION := 1.12.0
bin/gotestsum: GITHUB  := gotestyourself/gotestsum
bin/gotestsum: ARCHIVE := gotestsum_$(VERSION)_$(OSTYPE)_${ARCH}.tar.gz
bin/gotestsum: bin
	@ printf "Install gotestsum... "
	@ curl -Ls $(call github_url) | tar -zOxvf - gotestsum > $@ && chmod +x $@
	@ echo "done."

# ~~ [ tparse ] ~~~ https://github.com/mfridman/tparse ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

TPARSE := $(shell command -v tparse || echo "bin/tparse")
tparse: bin/tparse ## Installs tparse (testing go code)

bin/tparse: VERSION := 0.16.0
bin/tparse: GITHUB  := mfridman/tparse
bin/tparse: ARCHIVE := tparse_$(OSTYPE)_${ARCH}
bin/tparse: bin
	@ printf "Install tparse... "
	@ printf "Download from https://github.com/$(GITHUB)/releases/download/v$(VERSION)/$(ARCHIVE)... "
	@ curl -Ls $(shell echo $(call github_url) | tr A-Z a-z) > $@ && chmod +x $@
	@ echo "done."

# ~~ [ mockery ] ~~~ https://github.com/vektra/mockery ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

MOCKERY := $(shell command -v mockery || echo "bin/mockery")
mockery: bin/mockery ## Installs mockery (mocks generation)

bin/mockery: VERSION := 2.49.0
bin/mockery: GITHUB  := vektra/mockery
bin/mockery: ARCHIVE := mockery_$(VERSION)_$(OSTYPE)_${ARCH}.tar.gz
bin/mockery: bin
	@ printf "Install mockery... "
	@ curl -Ls $(call github_url) | tar -zOxf - mockery > $@ && chmod +x $@
	@ echo "done."

# ~~ [ golangci-lint ] ~~~ https://github.com/golangci/golangci-lint ~~~~~~~~~~~~~~~~~~~~~
GOLANGCI := $(shell command -v golangci-lint || echo "bin/golangci-lint")
golangci-lint: bin/golangci-lint ## Installs golangci-lint (linter)

bin/golangci-lint: VERSION := 1.62.0
bin/golangci-lint: GITHUB  := golangci/golangci-lint
bin/golangci-lint: ARCHIVE := golangci-lint-$(VERSION)-$(OSTYPE)-${ARCH}.tar.gz
bin/golangci-lint: bin
	@ printf "Install golangci-linter... "
	@ curl -Ls $(call github_url) | tar -zOxvf - $(shell printf golangci-lint-$(VERSION)-$(OSTYPE)-${ARCH}/golangci-lint) > $@ && chmod +x $@
	@ echo "done."

# ~~ [ go-arch-lint ] ~~~ https://github.com/fe3dback/go-arch-lint/ ~~~~~~~~~~~~~~~~~~~~~
GOARCH := $(shell command -v go-arch-lint || echo "bin/go-arch-lint")
go-arch-lint: bin/go-arch-lint ## Installs go-arch-lint (linter)

bin/go-arch-lint: VERSION := 1.11.6
bin/go-arch-lint: GITHUB  := fe3dback/go-arch-lint
bin/go-arch-lint: ARCHIVE := go-arch-lint_$(VERSION)_$(OSTYPE)_${ARCH}.tar.gz
bin/go-arch-lint: bin
	@ printf "Install go-arch-lint... "
	@ curl -Ls $(call github_url) | tar -zOxf - ./go-arch-lint > $@ && chmod +x $@
	@ echo "done."

# ~~ [ swag ] ~~~ https://github.com/swaggo/swag ~~~~~~~~~~~~~~~~~~~~~
SWAG := $(shell command -v swag || echo "bin/swag")
swag: bin/swag ## Installs swag (doc generator)

bin/swag: VERSION := 1.16.4
bin/swag: GITHUB  := swaggo/swag
bin/swag: ARCHIVE := swag_$(VERSION)_$(OSTYPE)_${ARCH}.tar.gz
bin/swag: bin
	@ printf "Install swag... "
	@ curl -Ls $(call github_url) | tar -zOxf - swag > $@ && chmod +x $@
	@ echo "done."
