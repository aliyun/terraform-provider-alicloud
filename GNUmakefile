TEST?=$$(go list ./... |grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
WEBSITE_REPO=github.com/hashicorp/terraform-website
PKG_NAME=alicloud
RELEASE_ALPHA_VERSION=$(VERSION)-alpha$(shell date +'%Y%m%d')
RELEASE_ALPHA_NAME=terraform-provider-alicloud_v$(RELEASE_ALPHA_VERSION)
PARALLEL ?= $(shell nproc 2>/dev/null || sysctl -n hw.ncpu 2>/dev/null || echo 4)

tools:
	@echo "==> installing required tooling..."
	go install golang.org/x/tools/cmd/goimports@latest
	GOTOOLCHAIN=go1.26.10 go install github.com/bflad/tfproviderlint/cmd/tfproviderlint@v0.31.0
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $$(go env GOPATH || $$GOPATH)/bin v2.12.2

default: build

build: fmtcheck	all

test: fmtcheck
	go test $(TEST) -timeout=30s -parallel=4

testacc: fmtcheck
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m

tflint:
	@echo "Run tfproviderlint ..."
	./scripts/run-tflint.sh

goimports: tools
	@echo "Fixing imports in all Go files..."
	@find . -name '*.go' -print0 \
	| xargs -0 -P "$(PARALLEL)" -I {} ./scripts/goimport-file.sh {}
	@echo "Done. Processed with up to $(PARALLEL) parallel jobs."

# Test a specific resource with debug logs
# Usage: make test-resource-debug RESOURCE=alicloud_vpc TESTCASE=basic LOGLEVEL=TRACE LOGFILE=vpc-test.log
# Flat-layout resources run against ./alicloud; split-layout resources
# (alicloud/service/<product>/) are discovered from the service package and
# run in their own subpackage.
test-resource-debug:
	@if [ -z "$(RESOURCE)" ]; then \
		echo "Error: RESOURCE is required. Usage: make test-resource-debug RESOURCE=alicloud_vpc"; \
		exit 1; \
	fi
	@RESOURCE_NAME=$$(echo "$(RESOURCE)" | sed 's/^alicloud_//'); \
	TEST_FILE="alicloud/resource_alicloud_$${RESOURCE_NAME}_test.go"; \
	TEST_TYPE="resource"; \
	TEST_PKG="./alicloud"; \
	if [ ! -f "$$TEST_FILE" ]; then \
		TEST_FILE="alicloud/data_source_alicloud_$${RESOURCE_NAME}_test.go"; \
		TEST_TYPE="data_source"; \
	fi; \
	if [ ! -f "$$TEST_FILE" ]; then \
		for f in alicloud/service/*/resource_*_test.go alicloud/service/*/data_source_*_test.go; do \
			[ -f "$$f" ] || continue; \
			d=$$(basename $$(dirname "$$f")); \
			b=$$(basename "$$f" _test.go); b=$${b#resource_}; b=$${b#data_source_}; \
			if [ "alicloud_$${d}_$${b}" = "$(RESOURCE)" ]; then \
				TEST_FILE="$$f"; \
				TEST_TYPE="split"; \
				TEST_PKG="./alicloud/service/$${d}"; \
				break; \
			fi; \
		done; \
	fi; \
	if [ ! -f "$$TEST_FILE" ]; then \
		echo "Error: Test file not found for resource $(RESOURCE)"; \
		echo "Tried alicloud/resource_alicloud_$${RESOURCE_NAME}_test.go, alicloud/data_source_alicloud_$${RESOURCE_NAME}_test.go and alicloud/service/*/{resource,data_source}_$${RESOURCE_NAME}_test.go"; \
		exit 1; \
	fi; \
	echo "Found test file: $$TEST_FILE ($$TEST_TYPE, package $$TEST_PKG)"; \
	LOGLEVEL=$${LOGLEVEL:-DEBUG}; \
	LOGFILE_BASE="$(LOGFILE)"; \
	PROJECT_ROOT=$$(pwd); \
	if [ -n "$$LOGFILE_BASE" ]; then \
		LOGFILE_BASE=$$(echo "$$LOGFILE_BASE" | sed 's/\.log$$//'); \
		API_LOG="$${PROJECT_ROOT}/$${LOGFILE_BASE}-api.log"; \
		CONSOLE_LOG="$${PROJECT_ROOT}/$${LOGFILE_BASE}-console.log"; \
		echo "API logs will be saved to: $$API_LOG"; \
		echo "Console output will be saved to: $$CONSOLE_LOG"; \
	fi; \
	echo "Log level: $$LOGLEVEL"; \
	echo ""; \
	ALL_TESTS=$$(grep -E "^func (TestAccAli[Cc]loud[A-Za-z0-9_]+)" "$$TEST_FILE" | sed -E 's/func ([A-Za-z0-9_]+).*/\1/'); \
	if [ -z "$$ALL_TESTS" ]; then \
		echo "Error: No test functions found in $$TEST_FILE"; \
		exit 1; \
	fi; \
	echo "=== Test cases to be executed ==="; \
	if [ -n "$(TESTCASE)" ]; then \
		SELECTED_TESTS=$$(echo "$$ALL_TESTS" | grep -i "$(TESTCASE)"); \
		if [ -z "$$SELECTED_TESTS" ]; then \
			echo "  No matching test cases found"; \
			exit 1; \
		fi; \
		echo "$$SELECTED_TESTS" | sed 's/^/  - /'; \
	else \
		SELECTED_TESTS="$$ALL_TESTS"; \
		echo "$$SELECTED_TESTS" | sed 's/^/  - /'; \
	fi; \
	echo "=================================="; \
	echo ""; \
	TEST_COUNT=$$(echo "$$SELECTED_TESTS" | wc -l | tr -d ' '); \
	if [ -n "$(TESTCASE)" ]; then \
		echo "Running $$TEST_COUNT test(s) for resource: $(RESOURCE), testcase: $(TESTCASE) with log level: $$LOGLEVEL"; \
	else \
		echo "Running $$TEST_COUNT test(s) for resource: $(RESOURCE) with log level: $$LOGLEVEL"; \
	fi; \
	echo "Test file: $$TEST_FILE"; \
	echo ""; \
	TEST_NUM=0; \
	FAILED_TESTS=""; \
	echo "$$SELECTED_TESTS" | while IFS= read -r TEST_NAME; do \
		TEST_NUM=$$((TEST_NUM + 1)); \
		echo "=================================================="; \
		echo "[$$TEST_NUM/$$TEST_COUNT] Running: $$TEST_NAME"; \
		echo "=================================================="; \
		if [ -n "$$LOGFILE_BASE" ]; then \
			TEST_API_LOG="$${PROJECT_ROOT}/$${LOGFILE_BASE}-$${TEST_NAME}-api.log"; \
			TEST_CONSOLE_LOG="$${PROJECT_ROOT}/$${LOGFILE_BASE}-$${TEST_NAME}-console.log"; \
			if TF_LOG=$$LOGLEVEL TF_LOG_PATH=$$TEST_API_LOG TF_ACC=1 go test -v $$TEST_PKG -run="^$$TEST_NAME\$$" -timeout 360m 2>&1 | tee $$TEST_CONSOLE_LOG; then \
				echo "✓ PASSED: $$TEST_NAME"; \
				echo "  API logs: $$TEST_API_LOG"; \
				echo "  Console logs: $$TEST_CONSOLE_LOG"; \
			else \
				echo "✗ FAILED: $$TEST_NAME"; \
				echo "  API logs: $$TEST_API_LOG"; \
				echo "  Console logs: $$TEST_CONSOLE_LOG"; \
				FAILED_TESTS="$$FAILED_TESTS$$TEST_NAME\n"; \
			fi; \
		else \
			if TF_LOG=$$LOGLEVEL TF_ACC=1 go test -v $$TEST_PKG -run="^$$TEST_NAME\$$" -timeout 360m; then \
				echo "✓ PASSED: $$TEST_NAME"; \
			else \
				echo "✗ FAILED: $$TEST_NAME"; \
				FAILED_TESTS="$$FAILED_TESTS$$TEST_NAME\n"; \
			fi; \
		fi; \
		echo ""; \
	done; \
	if [ -n "$$FAILED_TESTS" ]; then \
		echo "=================================================="; \
		echo "Failed tests:"; \
		echo "$$FAILED_TESTS" | sed 's/^/  - /'; \
		exit 1; \
	fi

vet:
	"$(CURDIR)/scripts/vetcheck.sh"

fmt:
	gofmt -w $(GOFMT_FILES)
	goimports -w $(GOFMT_FILES)

# Real production build of the root provider main, covering the root main
# package and every transitive production dependency it links — including the
# split alicloud/service/* subpackages. This is a genuine `go build`/link, not a
# renamed `vet`. It is the entry that proves the root main links; the cloudauth
# init probe (probe_init) compiles+inits a single service package via a
# no-such-test-case run and does NOT compile or link the root main.
provider-build:
	@echo "==> building root provider main (production)..."
	go build -o /dev/null .

# Assert the root alicloud production dependency graph (root main, non-test)
# does not depend — directly or transitively — on
# github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource. Unlike
# test-framework-guard.sh (a diff-level scan of newly added import lines),
# this queries the COMPLETE production dependency graph with `go list -deps`,
# so already-migrated and not-yet-migrated files are both covered.
check-production-deps:
	@"$(CURDIR)/scripts/ci/check-production-deps.sh"

fmtcheck:
	"$(CURDIR)/scripts/gofmtcheck.sh"

importscheck:
	"$(CURDIR)/scripts/goimportscheck.sh"

errcheck:
	@sh -c "'$(CURDIR)/scripts/errcheck.sh'"

test-compile:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./$(PKG_NAME)"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)

website:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

website-test:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), getting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	ln -sf ../../../../ext/providers/alicloud/website/docs $(GOPATH)/src/github.com/hashicorp/terraform-website/content/source/docs/providers/alicloud
	ln -sf ../../../ext/providers/alicloud/website/alicloud.erb $(GOPATH)/src/github.com/hashicorp/terraform-website/content/source/layouts/alicloud.erb
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider-test PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

commit:
	@bash "$(CURDIR)/scripts/generate-commit.sh"

# Local CI checks
# Usage:
#   make ci-check                    # Run all checks including example tests and resource integration tests (default)
#   make ci-check SKIP_EXAMPLE=1     # Run all checks except example tests
#   make ci-check SKIP_TEST=1        # Skip resource integration tests
#   make ci-check SKIP_BUILD=1       # Skip build check
#   make ci-check-quick              # Quick check (run changed-file unit tests; skip build, errcheck, and example tests)
#   CI_CHECK_BASE=<commit> make ci-check-quick  # Compare that base to the current working tree
# An explicit base must exist locally. Unset it to retain the default change detection.
ci-check: fmtcheck
	@echo "Building for current OS only..."
	@OS=$$(uname -s | tr '[:upper:]' '[:lower:]'); \
	if [ "$(SKIP_BUILD)" != "1" ]; then \
	    if [ "$$OS" = "darwin" ]; then \
	        $(MAKE) native; \
		elif [ "$$OS" = "linux" ]; then \
	        $(MAKE) linux; \
	    elif echo "$$OS" | grep -q "mingw\|msys\|cygwin"; then \
	        $(MAKE) windows; \
	    else \
	        echo "Unknown OS: $$OS, building for linux as fallback"; \
	        $(MAKE) linux; \
	    fi; \
	fi
	@CI_CHECK_FLAGS="--skip-build"; \
	if [ "$(SKIP_EXAMPLE)" = "1" ]; then \
		CI_CHECK_FLAGS="$$CI_CHECK_FLAGS --skip-example-test"; \
	fi; \
	if [ "$(SKIP_TEST)" = "1" ]; then \
		CI_CHECK_FLAGS="$$CI_CHECK_FLAGS --skip-resource-test"; \
	fi; \
	bash "$(CURDIR)/scripts/local-ci-check.sh" $$CI_CHECK_FLAGS

# Quick CI check (quality checks and changed-file unit tests in one run)
ci-check-quick:
	@bash "$(CURDIR)/scripts/local-ci-check.sh" --quick

# Calculate minimal test set for a resource (100% coverage)
# Usage: make minimal-test-set RESOURCE=alicloud_drds_polardbx_instance
# Usage: make minimal-test-set RESOURCE=alicloud_drds_polardbx_instance FORMAT=json
minimal-test-set:
	@if [ -z "$(RESOURCE)" ]; then \
		echo "Error: RESOURCE is required. Usage: make minimal-test-set RESOURCE=alicloud_drds_polardbx_instance"; \
		exit 1; \
	fi
	@FORMAT=$${FORMAT:-summary}; \
	if [ "$$FORMAT" != "summary" ] && [ "$$FORMAT" != "json" ]; then \
		echo "Error: FORMAT must be 'summary' or 'json' (default: summary)"; \
		exit 1; \
	fi; \
	go run scripts/testing/minimal_test_set_calculator.go -resource $(RESOURCE) -format $$FORMAT

# Sweep test resources in a specific region
# Usage:
#   make sweep REGION=cn-hangzhou RESOURCE=alicloud_vpc
#   make sweep REGION=cn-hangzhou RESOURCE=alicloud_vpc OLDER_THAN=12h              # only delete VPCs created more than 12h ago
#   make sweep REGION=cn-hangzhou RESOURCE=alicloud_vpc SWEEP_ALL=1                 # delete all resource names; VPC cleanup protects explicit do-not-delete VPCs
#   make sweep REGION=cn-hangzhou RESOURCE=alicloud_vpc SWEEP_ALL=1 OLDER_THAN=24h  # delete all, but protect resources newer than 24h
# Note: For resources with dependencies (specified in Dependencies field), include them with pipe separator.
#       Example: RESOURCE="alicloud_vpc_ipam_ipam_pool|alicloud_vpc_ipam_ipam"
# OLDER_THAN accepts a Go duration (12h, 24h, 30m, 1h30m). When set, resources created within that
# window are spared even if they match the name prefix — useful when several people run AccTest on
# the same account and you only want to clean up stale leftovers.
# SWEEP_ALL traverses a whole dependency graph and may exceed Go's default 10m.
# SWEEP_TIMEOUT overrides this inner Go deadline; remote runner deadlines still apply.
SWEEP_TIMEOUT ?= $(if $(filter 1,$(SWEEP_ALL)),3h,10m)
sweep:
	@if [ -z "$(REGION)" ]; then \
		echo "Error: REGION is required. Usage: make sweep REGION=cn-hangzhou RESOURCE=alicloud_instance"; \
		exit 1; \
	fi
	@if [ -z "$(RESOURCE)" ]; then \
		echo "Error: RESOURCE is required. Usage: make sweep REGION=cn-hangzhou RESOURCE=alicloud_instance"; \
		exit 1; \
	fi
	@if [ "$(SWEEP_ALL)" = "1" ]; then \
		echo "Sweeping ALL $(RESOURCE) resources in region $(REGION) (skipping prefix check)$(if $(OLDER_THAN), protecting resources newer than $(OLDER_THAN))..."; \
		ALICLOUD_SWEEP_ALL_RESOURCES=true ALICLOUD_SWEEP_OLDER_THAN=$(OLDER_THAN) TF_ACC=1 go test ./alicloud -v -timeout=$(SWEEP_TIMEOUT) -sweep=$(REGION) -sweep-run=$(RESOURCE); \
	else \
		echo "Sweeping $(RESOURCE) in region $(REGION)$(if $(OLDER_THAN), protecting resources newer than $(OLDER_THAN))..."; \
		ALICLOUD_SWEEP_OLDER_THAN=$(OLDER_THAN) TF_ACC=1 go test ./alicloud -v -timeout=$(SWEEP_TIMEOUT) -sweep=$(REGION) -sweep-run=$(RESOURCE); \
	fi

.PHONY: build test testacc test-resource-debug vet fmt fmtcheck importscheck goimports errcheck test-compile website website-test commit ci-check ci-check-quick minimal-test-set sweep tools native provider-build check-production-deps

all: mac windows linux

dev: clean mac copy

devlinux: clean fmt linux linuxcopy

devwin: clean fmt windows windowscopy

copy:
	tar -xvf bin/terraform-provider-alicloud_darwin-amd64.tgz && mv bin/terraform-provider-alicloud $(shell dirname `which terraform`)

clean:
	rm -rf bin/*

mac:
	GOOS=darwin GOARCH=amd64 go build -o bin/terraform-provider-alicloud
	tar czvf bin/terraform-provider-alicloud_darwin-amd64.tgz bin/terraform-provider-alicloud
	rm -rf bin/terraform-provider-alicloud

# Build for the host platform and architecture, named the way
# scripts/local-example-check.sh discovers binaries
# (terraform-provider-alicloud_<os>-<arch>), for local checks
native:
	GOOS= GOARCH= go build -o bin/terraform-provider-alicloud_$$(go env GOOS)-$$(go env GOARCH)

windowscopy:
	tar -xvf bin/terraform-provider-alicloud_windows-amd64.tgz && mv bin/terraform-provider-alicloud $(shell dirname `which terraform`)
    
windows:
	GOOS=windows GOARCH=amd64 go build -o bin/terraform-provider-alicloud.exe
	tar czvf bin/terraform-provider-alicloud_windows-amd64.tgz bin/terraform-provider-alicloud.exe
	rm -rf bin/terraform-provider-alicloud.exe

linuxcopy:
	tar -xvf bin/terraform-provider-alicloud_linux-amd64.tgz && mv bin/terraform-provider-alicloud $(shell dirname `which terraform`)

linux:
	GOOS=linux GOARCH=amd64 go build -o bin/terraform-provider-alicloud
	tar czvf bin/terraform-provider-alicloud_linux-amd64.tgz bin/terraform-provider-alicloud
	rm -rf bin/terraform-provider-alicloud

alpha:
	GOOS=linux GOARCH=amd64 go build -o bin/$(RELEASE_ALPHA_NAME)
	aliyun --profile terraformer --region cn-hangzhou oss cp bin/$(RELEASE_ALPHA_NAME) oss://iac-service-terraform/terraform/alphaplugins/registry.terraform.io/aliyun/alicloud/$(RELEASE_ALPHA_VERSION)/linux_amd64/$(RELEASE_ALPHA_NAME)
	#aliyun oss cp bin/$(RELEASE_ALPHA_NAME) oss://iac-service-terraform/terraform/alphaplugins/registry.terraform.io/hashicorp/alicloud/$(RELEASE_ALPHA_VERSION)/linux_amd64/$(RELEASE_ALPHA_NAME)  --profile terraformer --region cn-hangzhou
	rm -rf bin/$(RELEASE_ALPHA_NAME)

macarm:
	GOOS=darwin GOARCH=arm64 go build -o bin/terraform-provider-alicloud_v1.0.0
	cp bin/terraform-provider-alicloud_v1.0.0 ~/.terraform.d/plugins/registry.terraform.io/aliyun/alicloud/1.0.0/darwin_arm64/
	mv bin/terraform-provider-alicloud_v1.0.0 ~/.terraform.d/plugins/registry.terraform.io/hashicorp/alicloud/1.0.0/darwin_arm64/
