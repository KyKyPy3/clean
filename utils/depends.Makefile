define github_url
    https://github.com/$(GITHUB)/releases/download/v$(VERSION)/$(ARCHIVE)
endef

# creates a directory bin.
bin:
	@ mkdir -p $@

# ~~ [migrate] ~~~ https://github.com/golang-migrate/migrate ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

MIGRATE := $(shell command -v migrate || echo "bin/migrate")
migrate: bin/migrate ## Install migrate (database migration)

bin/migrate: VERSION := 4.16.2
bin/migrate: GITHUB  := golang-migrate/migrate
bin/migrate: ARCHIVE := migrate.$(OSTYPE)-amd64.tar.gz
bin/migrate: bin
	@ printf "Install migrate... "
	@ curl -Ls $(call github_url) | tar -zOxf - ./migrate > $@ && chmod +x $@
	@ echo "done."

# ~~ [ gotestsum ] ~~~ https://github.com/gotestyourself/gotestsum ~~~~~~~~~~~~~~~~~~~~~~~

GOTESTSUM := $(shell command -v gotestsum || echo "bin/gotestsum")
gotestsum: bin/gotestsum ## Installs gotestsum (testing go code)

bin/gotestsum: VERSION := 1.11.0
bin/gotestsum: GITHUB  := gotestyourself/gotestsum
bin/gotestsum: ARCHIVE := gotestsum_$(VERSION)_$(OSTYPE)_amd64.tar.gz
bin/gotestsum: bin
	@ printf "Install gotestsum... "
	@ curl -Ls $(call github_url) | tar -zOxvf - gotestsum > $@ && chmod +x $@
	@ echo "done."

# ~~ [ tparse ] ~~~ https://github.com/mfridman/tparse ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

TPARSE := $(shell command -v tparse || echo "bin/tparse")
tparse: bin/tparse ## Installs tparse (testing go code)

bin/tparse: VERSION := 0.8.3
bin/tparse: GITHUB  := mfridman/tparse
bin/tparse: ARCHIVE := tparse_$(VERSION)_$(OSTYPE)_x86_64.tar.gz
bin/tparse: bin
	@ printf "Install tparse... "
	@ curl -Ls $(shell echo $(call github_url) | tr A-Z a-z) | tar -zOxf - tparse > $@ && chmod +x $@
	@ echo "done."

# ~~ [ mockery ] ~~~ https://github.com/vektra/mockery ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

MOCKERY := $(shell command -v mockery || echo "bin/mockery")
mockery: bin/mockery ## Installs mockery (mocks generation)

bin/mockery: VERSION := 2.36.0
bin/mockery: GITHUB  := vektra/mockery
bin/mockery: ARCHIVE := mockery_$(VERSION)_$(OSTYPE)_x86_64.tar.gz
bin/mockery: bin
	@ printf "Install mockery... "
	@ curl -Ls $(call github_url) | tar -zOxf - mockery > $@ && chmod +x $@
	@ echo "done."

# ~~ [ golangci-lint ] ~~~ https://github.com/golangci/golangci-lint ~~~~~~~~~~~~~~~~~~~~~
GOLANGCI := $(shell command -v golangci-lint || echo "bin/golangci-lint")
golangci-lint: bin/golangci-lint ## Installs golangci-lint (linter)

bin/golangci-lint: VERSION := 1.55.2
bin/golangci-lint: GITHUB  := golangci/golangci-lint
bin/golangci-lint: ARCHIVE := golangci-lint-$(VERSION)-$(OSTYPE)-amd64.tar.gz
bin/golangci-lint: bin
	@ printf "Install golangci-linter... "
	@ curl -Ls $(call github_url) | tar -zOxvf - $(shell printf golangci-lint-$(VERSION)-$(OSTYPE)-amd64/golangci-lint) > $@ && chmod +x $@
	@ echo "done."

# ~~ [ swag ] ~~~ https://github.com/swaggo/swag ~~~~~~~~~~~~~~~~~~~~~
SWAG := $(shell command -v swag || echo "bin/swag")
swag: bin/swag ## Installs swag (doc generator)

bin/swag: VERSION := 1.16.2
bin/swag: GITHUB  := swaggo/swag
bin/swag: ARCHIVE := swag_$(VERSION)_$(OSTYPE)_x86_64.tar.gz
bin/swag: bin
	@ printf "Install swag... "
	@ curl -Ls $(call github_url) | tar -zOxf - swag > $@ && chmod +x $@
	@ echo "done."