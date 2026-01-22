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
	go install github.com/bflad/tfproviderlint/cmd/tfproviderlint@latest
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $$(go env GOPATH || $$GOPATH)/bin v2.7.2

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
test-resource-debug:
	@if [ -z "$(RESOURCE)" ]; then \
		echo "Error: RESOURCE is required. Usage: make test-resource-debug RESOURCE=alicloud_vpc"; \
		exit 1; \
	fi
	@RESOURCE_NAME=$$(echo "$(RESOURCE)" | sed 's/^alicloud_//'); \
	TEST_FILE="alicloud/resource_alicloud_$${RESOURCE_NAME}_test.go"; \
	TEST_TYPE="resource"; \
	if [ ! -f "$$TEST_FILE" ]; then \
		TEST_FILE="alicloud/data_source_alicloud_$${RESOURCE_NAME}_test.go"; \
		TEST_TYPE="data_source"; \
		if [ ! -f "$$TEST_FILE" ]; then \
			echo "Error: Test file not found for resource $(RESOURCE)"; \
			echo "Tried: resource_alicloud_$${RESOURCE_NAME}_test.go and data_source_alicloud_$${RESOURCE_NAME}_test.go"; \
			exit 1; \
		fi; \
	fi; \
	echo "Found test file: $$TEST_FILE ($$TEST_TYPE)"; \
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
			if TF_LOG=$$LOGLEVEL TF_LOG_PATH=$$TEST_API_LOG TF_ACC=1 go test -v ./alicloud -run="^$$TEST_NAME\$$" -timeout 360m 2>&1 | tee $$TEST_CONSOLE_LOG; then \
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
			if TF_LOG=$$LOGLEVEL TF_ACC=1 go test -v ./alicloud -run="^$$TEST_NAME\$$" -timeout 360m; then \
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
#   make ci-check                    # Run all checks including example tests (default)
#   make ci-check SKIP_EXAMPLE=1     # Run all checks except example tests
#   make ci-check SKIP_BUILD=1       # Skip build check
#   make ci-check-quick              # Quick check (skip build, tests, errcheck, and example tests)
ci-check:
	@if [ "$(SKIP_EXAMPLE)" = "1" ]; then \
		bash "$(CURDIR)/scripts/local-ci-check.sh" --skip-example-test; \
	elif [ "$(SKIP_BUILD)" = "1" ]; then \
		bash "$(CURDIR)/scripts/local-ci-check.sh" --skip-build; \
	else \
		bash "$(CURDIR)/scripts/local-ci-check.sh"; \
	fi

# Quick CI check (skip build and tests)
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

.PHONY: build test testacc test-resource test-resource-debug vet fmt fmtcheck errcheck test-compile website website-test commit ci-check ci-check-quick minimal-test-set

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
