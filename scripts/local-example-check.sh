#!/usr/bin/env bash

# Local Docs Example Test Script
# This script runs the same example tests as GitHub Actions CI locally
# Based on: ci/tasks/docs-example.sh and .github/workflows/acctest-terraform-integration.yml
#
# Usage:
#   ./scripts/local-example-check.sh              # Check changed docs
#   ./scripts/local-example-check.sh --all        # Check all docs
#   ./scripts/local-example-check.sh --file <path> # Check specific doc file
#   ./scripts/local-example-check.sh --dry-run    # Show what would be tested without running
#
# Requirements:
#   - Terraform must be installed
#   - Provider binary must be built (run 'make build' first)
#   - Valid AliCloud credentials (ALICLOUD_ACCESS_KEY, ALICLOUD_SECRET_KEY, ALICLOUD_REGION)
#
# Note: This runs actual terraform operations and may create real resources!
#       Make sure you understand what resources will be created.

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Global variable to track if we created terraformrc backup
TERRAFORMRC_BACKUP=""

# Cleanup function to restore .terraformrc
cleanup() {
  if [ "$TERRAFORMRC_BACKUP" = "CREATED" ]; then
    # We created a new .terraformrc, remove it
    echo -e "${BLUE}Removing temporary ~/.terraformrc...${NC}"
    rm -f "$HOME/.terraformrc"
  elif [ -n "$TERRAFORMRC_BACKUP" ] && [ -f "$TERRAFORMRC_BACKUP" ]; then
    # Restore from backup
    echo -e "${BLUE}Restoring ~/.terraformrc from backup...${NC}"
    mv "$TERRAFORMRC_BACKUP" "$HOME/.terraformrc"
  fi
}

# Register cleanup function to run on exit
trap cleanup EXIT INT TERM

# Parse command line arguments
CHECK_ALL=false
DRY_RUN=false
SPECIFIC_FILE=""
SKIP_BUILD=false

while [[ $# -gt 0 ]]; do
  case $1 in
    --all)
      CHECK_ALL=true
      shift
      ;;
    --file)
      SPECIFIC_FILE="$2"
      shift 2
      ;;
    --dry-run)
      DRY_RUN=true
      shift
      ;;
    --skip-build)
      SKIP_BUILD=true
      shift
      ;;
    -h|--help)
      echo "Usage: $0 [OPTIONS]"
      echo ""
      echo "Options:"
      echo "  --all                 Check all documentation files (not just changed)"
      echo "  --file <path>         Check a specific documentation file"
      echo "  --dry-run             Show what would be tested without running"
      echo "  --skip-build          Skip building provider binary"
      echo "  -h, --help            Show this help message"
      echo ""
      echo "Environment Variables (required for actual testing):"
      echo "  ALICLOUD_ACCESS_KEY   AliCloud Access Key"
      echo "  ALICLOUD_SECRET_KEY   AliCloud Secret Key"
      echo "  ALICLOUD_REGION       AliCloud Region (default: cn-hangzhou)"
      echo ""
      echo "Note: This script will create real resources in your AliCloud account!"
      exit 0
      ;;
    *)
      echo "Unknown option: $1"
      echo "Use --help to see available options"
      exit 1
      ;;
  esac
done

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

echo -e "${GREEN}╔════════════════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║  Terraform Provider AliCloud - Local Example Tests            ║${NC}"
echo -e "${GREEN}╚════════════════════════════════════════════════════════════════╝${NC}"
echo

cd "$PROJECT_ROOT"

# Check if terraform is installed
if ! command -v terraform &> /dev/null; then
  echo -e "${RED}ERROR: Terraform is not installed.${NC}"
  echo -e "${YELLOW}Please install Terraform first: https://www.terraform.io/downloads${NC}"
  exit 1
fi

echo -e "${GREEN}✓ Terraform is installed: $(terraform version | head -1)${NC}"
echo

# Build provider if not skipped
if [ "$SKIP_BUILD" = false ] && [ "$DRY_RUN" = false ]; then
  echo -e "${BLUE}Building provider binary...${NC}"
  if ! make build; then
    echo -e "${RED}ERROR: Failed to build provider${NC}"
    exit 1
  fi

  # Setup provider override in ~/.terraformrc
  echo -e "${BLUE}Setting up provider override in ~/.terraformrc...${NC}"

  PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
  PROVIDER_DIR="$PROJECT_ROOT/bin"
  TERRAFORMRC="$HOME/.terraformrc"

  # Extract provider binary from tar.gz if needed
  PLATFORM="$(go env GOOS)_$(go env GOARCH)"
  PROVIDER_ARCHIVE="$PROVIDER_DIR/terraform-provider-alicloud_${PLATFORM}.tgz"

  if [ -f "$PROVIDER_ARCHIVE" ]; then
    echo -e "${BLUE}  Extracting provider binary from archive...${NC}"
    cd "$PROVIDER_DIR"
    tar -xzf "terraform-provider-alicloud_${PLATFORM}.tgz" 2>/dev/null || true
    cd "$PROJECT_ROOT"
  fi

  # Verify provider binary exists
  if [ ! -f "$PROVIDER_DIR/terraform-provider-alicloud" ]; then
    echo -e "${RED}ERROR: Provider binary not found at $PROVIDER_DIR/terraform-provider-alicloud${NC}"
    exit 1
  fi

  # Make sure provider binary is executable
  chmod +x "$PROVIDER_DIR/terraform-provider-alicloud"

  # Backup existing .terraformrc if it exists
  if [ -f "$TERRAFORMRC" ]; then
    TERRAFORMRC_BACKUP="${TERRAFORMRC}.backup.$(date +%s)"
    cp "$TERRAFORMRC" "$TERRAFORMRC_BACKUP"
    echo -e "${YELLOW}  Backed up existing ~/.terraformrc to ${TERRAFORMRC_BACKUP}${NC}"
  else
    # Mark that we need to clean up (remove the file we're about to create)
    TERRAFORMRC_BACKUP="CREATED"
  fi

  # Create or update .terraformrc with dev_overrides
  cat > "$TERRAFORMRC" <<EOF
provider_installation {
  dev_overrides {
    "registry.terraform.io/aliyun/alicloud" = "$PROVIDER_DIR"
  }
  direct {}
}
EOF

  echo -e "${GREEN}✓ Provider override configured in ~/.terraformrc${NC}"
  echo -e "  Provider directory: ${PROVIDER_DIR}"
  echo -e "  Provider binary: terraform-provider-alicloud"
  echo
fi

# Check credentials (skip in dry-run mode)
if [ "$DRY_RUN" = false ]; then
  if [ -z "$ALICLOUD_ACCESS_KEY" ] || [ -z "$ALICLOUD_SECRET_KEY" ]; then
    echo -e "${YELLOW}WARNING: AliCloud credentials not set!${NC}"
    echo -e "${YELLOW}Please set the following environment variables:${NC}"
    echo -e "  export ALICLOUD_ACCESS_KEY=your_access_key"
    echo -e "  export ALICLOUD_SECRET_KEY=your_secret_key"
    echo -e "  export ALICLOUD_REGION=cn-hangzhou  # optional, defaults to cn-hangzhou"
    echo
    echo -e "${RED}Credentials are required to run actual tests.${NC}"
    echo -e "${YELLOW}Use --dry-run to see what would be tested without credentials.${NC}"
    exit 1
  fi

  ALICLOUD_REGION=${ALICLOUD_REGION:-cn-hangzhou}
  echo -e "${GREEN}✓ AliCloud credentials configured${NC}"
  echo -e "  Region: ${ALICLOUD_REGION}"
  echo
fi

# Get files to check
FILES_TO_CHECK=""

if [ -n "$SPECIFIC_FILE" ]; then
  if [ ! -f "$SPECIFIC_FILE" ]; then
    echo -e "${RED}ERROR: File not found: $SPECIFIC_FILE${NC}"
    exit 1
  fi
  FILES_TO_CHECK="$SPECIFIC_FILE"
  echo -e "${BLUE}Checking specific file: ${SPECIFIC_FILE}${NC}"
elif [ "$CHECK_ALL" = true ]; then
  FILES_TO_CHECK=$(find website/docs/r website/docs/d -name "*.html.markdown" 2>/dev/null || true)
  echo -e "${BLUE}Checking all documentation files${NC}"
else
  # Get changed files
  CHANGED_FILES=$(git diff --name-only origin/master...HEAD 2>/dev/null || git diff --name-only HEAD~1 HEAD 2>/dev/null || true)

  if [ -z "$CHANGED_FILES" ]; then
    echo -e "${YELLOW}No changed files detected. Use --all to check all files.${NC}"
    exit 0
  fi

  # Filter for documentation files
  FILES_TO_CHECK=$(echo "$CHANGED_FILES" | grep -E "^website/docs/[rd]/.*\.html\.markdown$" || true)

  if [ -z "$FILES_TO_CHECK" ]; then
    echo -e "${YELLOW}No documentation files changed. Skipping example tests.${NC}"
    exit 0
  fi

  echo -e "${BLUE}Checking changed documentation files:${NC}"
  echo "$FILES_TO_CHECK" | sed 's/^/  - /'
fi

echo
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo -e "${BLUE}  Analyzing Documentation Files${NC}"
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo

TOTAL_FILES=0
TOTAL_EXAMPLES=0
FILES_WITH_EXAMPLES=""

# Count examples in each file
for doc_file in $FILES_TO_CHECK; do
  if [ ! -f "$doc_file" ]; then
    continue
  fi

  example_count=$(grep -c '```terraform' "$doc_file" || echo "0")

  if [ "$example_count" -gt 0 ]; then
    TOTAL_FILES=$((TOTAL_FILES + 1))
    TOTAL_EXAMPLES=$((TOTAL_EXAMPLES + example_count))
    FILES_WITH_EXAMPLES="${FILES_WITH_EXAMPLES}${doc_file}:${example_count} "

    echo -e "${GREEN}✓${NC} $doc_file: ${example_count} example(s)"
  else
    echo -e "${YELLOW}⚠${NC} $doc_file: No examples found"
  fi
done

echo
echo -e "${BLUE}Summary:${NC}"
echo -e "  Files with examples: ${TOTAL_FILES}"
echo -e "  Total examples: ${TOTAL_EXAMPLES}"
echo

if [ "$TOTAL_EXAMPLES" -eq 0 ]; then
  echo -e "${YELLOW}No examples found to test.${NC}"
  exit 0
fi

if [ "$DRY_RUN" = true ]; then
  echo -e "${YELLOW}Dry run mode: Not executing tests${NC}"
  exit 0
fi

echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo -e "${BLUE}  Running Example Tests${NC}"
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo
echo -e "${RED}WARNING: This will create real resources in your AliCloud account!${NC}"
echo -e "${YELLOW}Press Ctrl+C within 5 seconds to cancel...${NC}"
sleep 5
echo

# Create temporary directory for test runs
TEST_DIR=$(mktemp -d)
echo -e "${BLUE}Test directory: ${TEST_DIR}${NC}"
echo

# Track test results
PASSED_EXAMPLES=0
FAILED_EXAMPLES=0
FAILED_EXAMPLE_LIST=""

# Process each file
for file_info in $FILES_WITH_EXAMPLES; do
  doc_file="${file_info%:*}"
  example_count="${file_info##*:}"

  echo -e "${GREEN}═══════════════════════════════════════════════════════════════${NC}"
  echo -e "${GREEN}  Testing: ${doc_file} (${example_count} example(s))${NC}"
  echo -e "${GREEN}═══════════════════════════════════════════════════════════════${NC}"
  echo

  # Determine resource name from file path
  resource_name=$(basename "$doc_file" .html.markdown)

  # Extract and run each example
  example_num=0
  in_example=false
  example_content=""

  while IFS= read -r line; do
    if [[ "$line" == '```terraform' ]]; then
      in_example=true
      example_content=""
      continue
    fi

    if [[ "$line" == '```' ]] && [ "$in_example" = true ]; then
      in_example=false

      # Create example directory
      if [[ "$doc_file" == *"/d/"* ]]; then
        example_dir="${TEST_DIR}/data_source_alicloud_${resource_name}_example_${example_num}"
      else
        example_dir="${TEST_DIR}/resource_alicloud_${resource_name}_example_${example_num}"
      fi

      mkdir -p "$example_dir"

      # Create provider configuration
      cat > "$example_dir/terraform.tf" <<EOF
terraform {
  required_providers {
    alicloud = {
      source = "registry.terraform.io/aliyun/alicloud"
    }
  }
}

EOF

      # Write example content to main.tf
      echo "$example_content" > "$example_dir/main.tf"

      # Run terraform operations
      example_name=$(basename "$example_dir")
      echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
      echo -e "${BLUE}▶ Testing: ${example_name}${NC}"
      echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"

      test_failed=false
      error_log="$example_dir/error.log"

      # APPLY phase
      echo -e "${YELLOW}=== RUN   ${example_name} APPLY${NC}"

      # Run terraform init
      if ! (cd "$example_dir" && \
          TF_VAR_access_key="$ALICLOUD_ACCESS_KEY" \
          TF_VAR_secret_key="$ALICLOUD_SECRET_KEY" \
          TF_VAR_region="$ALICLOUD_REGION" \
          terraform init > "$error_log" 2>&1); then
        echo -e "${RED} - init failed.${NC}"
        echo -e "${YELLOW}Error details:${NC}"
        cat "$error_log" | tail -20
        test_failed=true
      # Run terraform plan
      elif ! (cd "$example_dir" && \
          TF_VAR_access_key="$ALICLOUD_ACCESS_KEY" \
          TF_VAR_secret_key="$ALICLOUD_SECRET_KEY" \
          TF_VAR_region="$ALICLOUD_REGION" \
          terraform plan > "$error_log" 2>&1); then
        echo -e "${RED} - plan failed.${NC}"
        echo -e "${YELLOW}Error details:${NC}"
        cat "$error_log" | tail -20
        test_failed=true
      # Run terraform apply
      elif ! (cd "$example_dir" && \
          TF_VAR_access_key="$ALICLOUD_ACCESS_KEY" \
          TF_VAR_secret_key="$ALICLOUD_SECRET_KEY" \
          TF_VAR_region="$ALICLOUD_REGION" \
          terraform apply -auto-approve > "$error_log" 2>&1); then
        echo -e "${RED} - apply failed.${NC}"
        echo -e "${YELLOW}Error details:${NC}"
        cat "$error_log" | tail -20
        test_failed=true
      else
        echo -e "${GREEN} - apply check: success.${NC}"

        # Double check: verify no drift and no deprecated attributes
        plan_output="$example_dir/plan_output.log"
        if ! (cd "$example_dir" && \
            TF_VAR_access_key="$ALICLOUD_ACCESS_KEY" \
            TF_VAR_secret_key="$ALICLOUD_SECRET_KEY" \
            TF_VAR_region="$ALICLOUD_REGION" \
            terraform plan > "$plan_output" 2>&1); then
          echo -e "${RED} - apply diff check: fail (plan failed after apply).${NC}"
          echo -e "${YELLOW}Plan output:${NC}"
          cat "$plan_output" | tail -20
          test_failed=true
        else
          # Check for deprecated attributes
          if grep -qi "deprecated" "$plan_output"; then
            echo -e "${RED} - deprecated attributes check: fail.${NC}"
            echo -e "${YELLOW}Deprecated warnings:${NC}"
            grep -i "deprecated" "$plan_output"
            test_failed=true
          fi

          # Check for no changes (No changes. Your infrastructure matches the configuration.)
          if ! grep -q "No changes" "$plan_output"; then
            echo -e "${RED} - apply diff check: fail.${NC}"
            echo -e "${YELLOW}Unexpected changes detected:${NC}"
            cat "$plan_output" | tail -30
            test_failed=true
          else
            echo -e "${GREEN} - apply diff check: success.${NC}"
          fi
        fi
      fi

      # DESTROY phase
      echo -e "${YELLOW}=== RUN   ${example_name} DESTROY${NC}"

      if (cd "$example_dir" && \
          TF_VAR_access_key="$ALICLOUD_ACCESS_KEY" \
          TF_VAR_secret_key="$ALICLOUD_SECRET_KEY" \
          TF_VAR_region="$ALICLOUD_REGION" \
          terraform destroy -auto-approve > "$error_log" 2>&1); then
        echo -e "${GREEN} - destroy check: success.${NC}"

        # Verify complete destruction
        if (cd "$example_dir" && \
            TF_VAR_access_key="$ALICLOUD_ACCESS_KEY" \
            TF_VAR_secret_key="$ALICLOUD_SECRET_KEY" \
            TF_VAR_region="$ALICLOUD_REGION" \
            terraform plan -destroy -detailed-exitcode > "$error_log" 2>&1); then
          echo -e "${GREEN} - destroy diff check: success.${NC}"
        else
          echo -e "${RED} - destroy diff check: fail (resources not fully destroyed).${NC}"
          echo -e "${YELLOW}Remaining resources:${NC}"
          cat "$error_log" | tail -20
          test_failed=true
        fi
      else
        echo -e "${RED} - destroy check: fail.${NC}"
        echo -e "${YELLOW}Error details:${NC}"
        cat "$error_log" | tail -20
        test_failed=true
      fi

      # Mark result
      if [ "$test_failed" = true ]; then
        echo -e "${RED}--- FAIL: ${example_name}${NC}"
        FAILED_EXAMPLES=$((FAILED_EXAMPLES + 1))
        FAILED_EXAMPLE_LIST="${FAILED_EXAMPLE_LIST}${example_name} "
      else
        echo -e "${GREEN}--- PASS: ${example_name}${NC}"
        PASSED_EXAMPLES=$((PASSED_EXAMPLES + 1))
      fi

      echo
      example_num=$((example_num + 1))
      continue
    fi

    if [ "$in_example" = true ]; then
      example_content="${example_content}${line}"$'\n'
    fi
  done < "$doc_file"
done

# Cleanup
echo -e "${BLUE}Cleaning up test directory...${NC}"
rm -rf "$TEST_DIR"
echo

# Print summary
echo -e "${GREEN}╔════════════════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║  Example Test Summary                                          ║${NC}"
echo -e "${GREEN}╚════════════════════════════════════════════════════════════════╝${NC}"
echo
echo -e "${BLUE}Total Examples: ${TOTAL_EXAMPLES}${NC}"
echo -e "${GREEN}Passed: ${PASSED_EXAMPLES}${NC}"
echo -e "${RED}Failed: ${FAILED_EXAMPLES}${NC}"
echo

if [ "$FAILED_EXAMPLES" -gt 0 ]; then
  echo -e "${RED}Failed Examples:${NC}"
  for example in $FAILED_EXAMPLE_LIST; do
    echo -e "  ${RED}✗${NC} $example"
  done
  echo
  echo -e "${RED}════════════════════════════════════════════════════════════════${NC}"
  echo -e "${RED}  Some example tests failed!${NC}"
  echo -e "${RED}════════════════════════════════════════════════════════════════${NC}"
  exit 1
else
  echo -e "${GREEN}════════════════════════════════════════════════════════════════${NC}"
  echo -e "${GREEN}  ✓ All example tests passed!${NC}"
  echo -e "${GREEN}════════════════════════════════════════════════════════════════${NC}"
  exit 0
fi
