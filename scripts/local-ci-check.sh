#!/bin/bash

# Local CI Checks Script
# This script runs the same checks as GitHub Actions CI locally before pushing
# Based on: .github/workflows/pull_requests.yml, document-check.yml, acctest-terraform-lint.yml
#
# Usage:
#   make ci-check              # Run from project root
#   ./scripts/local-ci-check.sh  # Run directly
#   ./scripts/local-ci-check.sh --skip-build  # Skip build check
#   ./scripts/local-ci-check.sh --skip-test   # Skip unit tests
#   ./scripts/local-ci-check.sh --quick       # Skip build and tests
#
# Features:
#   - Auto-installs missing tools (terrafmt, misspell, markdownlint, markdown-link-check)
#   - Runs code quality checks (gofmt, goimports, go vet, errcheck)
#   - Checks documentation (links, formatting, spelling, content, consistency)
#   - Builds and tests the project
#   - Beautiful colored output with summary

set -e

# Parse command line arguments
SKIP_BUILD=false
SKIP_TEST=true  # Default: skip tests locally (CI still runs them)
SKIP_ERRCHECK=true  # Default: skip errcheck locally (CI still runs it)
SKIP_MARKDOWN_LINT=true  # Default: skip markdown lint locally (CI still runs it)
SKIP_MARKDOWN_LINK_CHECK=true  # Default: skip markdown link check locally (CI still runs it)
SKIP_EXAMPLE_TEST=false  # Default: run example tests (use --skip-example-test to disable)
STRICT_MODE=false  # Strict mode checks ALL docs like CI does

while [[ $# -gt 0 ]]; do
  case $1 in
    --skip-build)
      SKIP_BUILD=true
      shift
      ;;
    --skip-test)
      SKIP_TEST=true
      shift
      ;;
    --run-test)
      SKIP_TEST=false
      shift
      ;;
    --run-markdown-lint)
      SKIP_MARKDOWN_LINT=false
      shift
      ;;
    --run-markdown-link-check)
      SKIP_MARKDOWN_LINK_CHECK=false
      shift
      ;;
    --run-errcheck)
      SKIP_ERRCHECK=false
      shift
      ;;
    --skip-example-test)
      SKIP_EXAMPLE_TEST=true
      shift
      ;;
    --quick)
      SKIP_BUILD=true
      SKIP_TEST=true
      SKIP_ERRCHECK=true
      SKIP_EXAMPLE_TEST=true
      shift
      ;;
    --strict)
      STRICT_MODE=true
      shift
      ;;
    -h|--help)
      echo "Usage: $0 [OPTIONS]"
      echo ""
      echo "Options:"
      echo "  --skip-build              Skip build check"
      echo "  --skip-test               Skip unit tests (default: skipped locally)"
      echo "  --run-test                Run unit tests"
      echo "  --run-errcheck            Run errcheck (default: skipped locally)"
      echo "  --run-markdown-lint       Run markdown lint (default: skipped locally)"
      echo "  --run-markdown-link-check Run markdown link check (default: skipped locally)"
      echo "  --skip-example-test       Skip example tests (default: enabled)"
      echo "  --quick                   Skip build, tests, errcheck, and example tests (faster checks)"
      echo "  --strict                  Check ALL docs (like CI), not just changed files"
      echo "  -h, --help                Show this help message"
      echo ""
      echo "Note: By default, example tests are ENABLED and will create real resources."
      echo "      Use --skip-example-test to disable example tests."
      echo "      Unit tests, errcheck, markdown lint, and markdown link checks are skipped locally."
      echo "      Use --run-test, --run-errcheck, --run-markdown-lint, or --run-markdown-link-check to enable them."
      exit 0
      ;;
    *)
      echo "Unknown option: $1"
      echo "Use --help to see available options"
      exit 1
      ;;
  esac
done

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

echo -e "${GREEN}╔════════════════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║  Terraform Provider AliCloud - Local CI Checks                ║${NC}"
echo -e "${GREEN}╚════════════════════════════════════════════════════════════════╝${NC}"
echo

# Track overall status
FAILED_CHECKS=()
PASSED_CHECKS=()

# Function to install a tool if not present
install_tool() {
  local tool_name="$1"
  local install_command="$2"
  local check_command="${3:-command -v $tool_name}"
  
  if eval "$check_command" &> /dev/null; then
    return 0
  fi
  
  echo -e "${YELLOW}⚠ Tool '${tool_name}' not found. Attempting to install...${NC}"
  echo -e "${BLUE}   Installing: ${install_command}${NC}"
  
  # Create a temporary file for the installation output
  local temp_output=$(mktemp)
  
  if eval "$install_command" > "$temp_output" 2>&1; then
    echo -e "${GREEN}✓ Successfully installed ${tool_name}${NC}"
    rm -f "$temp_output"
    echo
    return 0
  else
    echo -e "${RED}✗ Failed to install ${tool_name}${NC}"
    echo -e "${YELLOW}   Error output (last 5 lines):${NC}"
    tail -5 "$temp_output" | sed 's/^/   /'
    echo -e "${YELLOW}   You can manually install with: ${install_command}${NC}"
    rm -f "$temp_output"
    echo
    return 1
  fi
}

# Function to run a check
run_check() {
  local check_name="$1"
  local check_command="$2"
  
  echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
  echo -e "${BLUE}▶ Running: ${check_name}${NC}"
  echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
  
  if eval "$check_command"; then
    echo -e "${GREEN}✓ PASSED: ${check_name}${NC}"
    PASSED_CHECKS+=("$check_name")
    echo
    return 0
  else
    echo -e "${RED}✗ FAILED: ${check_name}${NC}"
    FAILED_CHECKS+=("$check_name")
    echo
    return 1
  fi
}

# Check if we're in a git repository
if ! git rev-parse --git-dir > /dev/null 2>&1; then
  echo -e "${RED}ERROR: Not in a git repository.${NC}"
  exit 1
fi

cd "$PROJECT_ROOT"

# Check Go version
if [ -f "$PROJECT_ROOT/.go-version" ]; then
  EXPECTED_GO_VERSION=$(cat "$PROJECT_ROOT/.go-version")
  CURRENT_GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
  
  echo -e "${BLUE}Go Version Check:${NC}"
  echo -e "  Expected: ${GREEN}${EXPECTED_GO_VERSION}${NC}"
  echo -e "  Current:  ${GREEN}${CURRENT_GO_VERSION}${NC}"
  
  # Compare major.minor version
  EXPECTED_MAJOR_MINOR=$(echo "$EXPECTED_GO_VERSION" | cut -d. -f1,2)
  CURRENT_MAJOR_MINOR=$(echo "$CURRENT_GO_VERSION" | cut -d. -f1,2)
  
  if [ "$EXPECTED_MAJOR_MINOR" != "$CURRENT_MAJOR_MINOR" ]; then
    echo -e "${YELLOW}⚠ Warning: Go version mismatch (expected ${EXPECTED_GO_VERSION}, got ${CURRENT_GO_VERSION})${NC}"
    echo -e "${YELLOW}  This may cause issues. Consider using the expected version.${NC}"
    echo
  else
    echo -e "${GREEN}✓ Go version is compatible${NC}"
    echo
  fi
fi

# Get changed files from the latest commit
# Priority:
# 1. Uncommitted changes (if any)
# 2. Latest commit changes
# 3. All changes from master branch
CHANGED_FILES=$(git diff --name-only HEAD 2>/dev/null)

if [ -z "$CHANGED_FILES" ]; then
  # No uncommitted changes, check the latest commit
  echo -e "${BLUE}Checking latest commit changes...${NC}"
  CHANGED_FILES=$(git diff --name-only HEAD~1 HEAD 2>/dev/null)
fi

if [ -z "$CHANGED_FILES" ]; then
  # No changes in latest commit, fall back to comparing with master
  echo -e "${YELLOW}No changes in latest commit. Comparing with master branch...${NC}"
  CHANGED_FILES=$(git diff --name-only origin/master...HEAD 2>/dev/null || git diff --name-only master...HEAD 2>/dev/null)
fi

if [ -z "$CHANGED_FILES" ]; then
  echo -e "${YELLOW}No changed files detected. Running full checks...${NC}"
else
  echo -e "${BLUE}Changed files ($(echo "$CHANGED_FILES" | wc -l | tr -d ' ')):${NC}"
  echo "$CHANGED_FILES" | sed 's/^/  - /'
fi
echo

# ============================================================================
# 0. Check and Install Required Tools
# ============================================================================

echo -e "${GREEN}═══════════════════════════════════════════════════════════════${NC}"
echo -e "${GREEN}  Checking Required Tools${NC}"
echo -e "${GREEN}═══════════════════════════════════════════════════════════════${NC}"
echo

# Check for goimports (needed for code quality checks)
if ! command -v goimports &> /dev/null; then
  install_tool "goimports" "go install golang.org/x/tools/cmd/goimports@latest"
fi

if command -v goimports &> /dev/null; then
  echo -e "${GREEN}✓ goimports is available${NC}"
else
  echo -e "${YELLOW}⚠ goimports not available, some checks may fail${NC}"
fi
echo

# ============================================================================
# 1. Code Quality Checks
# ============================================================================

echo -e "${GREEN}═══════════════════════════════════════════════════════════════${NC}"
echo -e "${GREEN}  Part 1: Code Quality Checks${NC}"
echo -e "${GREEN}═══════════════════════════════════════════════════════════════${NC}"
echo

# 1.1 Go Format Check
run_check "Go Format Check (gofmt)" \
  "\"$SCRIPT_DIR/gofmtcheck.sh\"" \
  || true

# 1.2 Go Imports Check
run_check "Go Imports Check (goimports)" \
  "\"$SCRIPT_DIR/goimportscheck.sh\"" \
  || true

# 1.3 Go Vet
run_check "Go Vet" \
  "make -C \"$PROJECT_ROOT\" vet" \
  || true

# 1.4 Error Check
if [ "$SKIP_ERRCHECK" = true ]; then
  echo -e "${YELLOW}⚠ Skipping Error Check (errcheck) (disabled locally, use --run-errcheck to enable)${NC}"
  echo
else
  run_check "Error Check (errcheck)" \
    "\"$SCRIPT_DIR/errcheck.sh\"" \
    || true
fi

# 1.5 Basic Check (fmt.Println, etc.)
if echo "$CHANGED_FILES" | grep -q "\.go$\|website/docs"; then
  run_check "Basic Check (fmt.Println, doc links)" \
    "\"$SCRIPT_DIR/basic-check.sh\"" \
    || true
fi

# ============================================================================
# 2. Testing Coverage Rate Check
# ============================================================================

echo -e "${GREEN}═══════════════════════════════════════════════════════════════${NC}"
echo -e "${GREEN}  Part 2: Testing Coverage Rate Check${NC}"
echo -e "${GREEN}═══════════════════════════════════════════════════════════════${NC}"
echo

# Check if there are any resource or data source changes
RESOURCE_CHANGES=$(echo "$CHANGED_FILES" | grep -E "alicloud/(resource|data_source).*\.go$" | grep -v "_test\.go$" || true)

if [ -n "$RESOURCE_CHANGES" ]; then
  echo -e "${BLUE}▶ Running testing coverage rate check...${NC}"
  
  # Create a temporary diff file
  TEMP_DIFF=$(mktemp)
  
  # Generate diff: prioritize uncommitted changes, then latest commit, then compare with master
  if git diff --name-only HEAD 2>/dev/null | grep -q .; then
    # There are uncommitted changes
    git diff HEAD > "$TEMP_DIFF"
  else
    # No uncommitted changes, check the latest commit
    git diff HEAD~1 HEAD > "$TEMP_DIFF"
  fi
  
  # If diff is empty, try comparing with master
  if [ ! -s "$TEMP_DIFF" ]; then
    git diff origin/master...HEAD > "$TEMP_DIFF" 2>/dev/null || git diff master...HEAD > "$TEMP_DIFF" 2>/dev/null || true
  fi
  
  if [ -s "$TEMP_DIFF" ]; then
    if go run "$SCRIPT_DIR/testing/testing_coverage_rate_check.go" -fileNames="$TEMP_DIFF"; then
      echo -e "${GREEN}✓ PASSED: Testing coverage rate check${NC}"
      PASSED_CHECKS+=("Testing coverage rate check")
    else
      echo -e "${RED}✗ FAILED: Testing coverage rate check${NC}"
      FAILED_CHECKS+=("Testing coverage rate check")
    fi
  else
    echo -e "${YELLOW}⚠ No changes to check${NC}"
  fi
  
  rm -f "$TEMP_DIFF"
  echo
else
  echo -e "${BLUE}▶ No resource or data source changes detected${NC}"
  echo -e "${GREEN}✓ SKIPPED: Testing coverage rate check (no resource changes)${NC}"
  echo
fi

# ============================================================================
# 3. Documentation Checks
# ============================================================================

if echo "$CHANGED_FILES" | grep -q "website/docs"; then
  echo -e "${GREEN}═══════════════════════════════════════════════════════════════${NC}"
  echo -e "${GREEN}  Part 2: Documentation Checks${NC}"
  echo -e "${GREEN}═══════════════════════════════════════════════════════════════${NC}"
  echo

  # 2.1 Check for exact help.aliyun.com links
  echo -e "${BLUE}▶ Checking documentation links...${NC}"
  error=false
  for doc in $CHANGED_FILES; do
    if [[ "$doc" == website/docs/d/* ]] || [[ "$doc" == website/docs/r/* ]]; then
      if [ -f "$doc" ]; then
        if grep -q "https://help.aliyun.com/)\.$" "$doc" 2>/dev/null; then
          echo -e "${RED}  ✗ Doc :${doc}: Please input the exact link, not just https://help.aliyun.com/${NC}"
          error=true
        fi
      fi
    fi
  done
  
  if [ "$error" = false ]; then
    echo -e "${GREEN}✓ PASSED: Documentation link check${NC}"
    PASSED_CHECKS+=("Documentation link check")
  else
    echo -e "${RED}✗ FAILED: Documentation link check${NC}"
    FAILED_CHECKS+=("Documentation link check")
  fi
  echo

  # 2.2 Terrafmt check (Terraform code formatting in docs)
  if ! command -v terrafmt &> /dev/null; then
    install_tool "terrafmt" "go install github.com/katbyte/terrafmt@latest"
  fi

  if command -v terrafmt &> /dev/null; then
    if [ "$STRICT_MODE" = true ]; then
      # CI mode: check all docs
      run_check "Terrafmt Check (Terraform code in docs)" \
        "terrafmt diff ./website --check --pattern '*.markdown'" \
        || true
    else
      # Smart mode: check only changed docs
      CHANGED_DOCS=$(echo "$CHANGED_FILES" | grep "website/docs.*\.markdown$" || true)
      if [[ -n "$CHANGED_DOCS" ]]; then
        run_check "Terrafmt Check (changed docs only)" \
          "echo '$CHANGED_DOCS' | xargs -I {} terrafmt diff {} --check" \
          || true
      else
        echo -e "${BLUE}▶ Terrafmt Check: No changed markdown files${NC}"
        echo
      fi
    fi
  else
    echo -e "${YELLOW}⚠ Skipping Terrafmt check (installation failed)${NC}"
    echo
  fi

  # 2.3 Misspell check
  if ! command -v misspell &> /dev/null; then
    install_tool "misspell" "go install github.com/client9/misspell/cmd/misspell@latest"
  fi

  if command -v misspell &> /dev/null; then
    if [ "$STRICT_MODE" = true ]; then
      # CI mode: check all docs
      run_check "Spell Check (misspell)" \
        "misspell -error -locale=US ./website/docs" \
        || true
    else
      # Smart mode: check only changed docs
      CHANGED_DOCS=$(echo "$CHANGED_FILES" | grep "website/docs" || true)
      if [[ -n "$CHANGED_DOCS" ]]; then
        run_check "Spell Check (changed docs only)" \
          "echo '$CHANGED_DOCS' | xargs misspell -error -locale=US" \
          || true
      else
        echo -e "${BLUE}▶ Misspell Check: No changed doc files${NC}"
        echo
      fi
    fi
  else
    echo -e "${YELLOW}⚠ Skipping Misspell check (installation failed)${NC}"
    echo
  fi

  # 2.4 Markdown lint
  if [ "$SKIP_MARKDOWN_LINT" = true ]; then
    echo -e "${YELLOW}⚠ Skipping Markdown lint (disabled locally, use --run-markdown-lint to enable)${NC}"
    echo
  else
    if ! command -v markdownlint &> /dev/null && ! command -v markdownlint-cli &> /dev/null; then
      # Check if npm is available
      if command -v npm &> /dev/null; then
        install_tool "markdownlint-cli" "npm install -g markdownlint-cli" "command -v markdownlint"
      else
        echo -e "${YELLOW}⚠ npm not found, cannot install markdownlint-cli${NC}"
        echo -e "${YELLOW}  Please install Node.js and npm first${NC}"
        echo
      fi
    fi

    if command -v markdownlint &> /dev/null || command -v markdownlint-cli &> /dev/null; then
      if [ "$STRICT_MODE" = true ]; then
        # CI mode: check all docs
        run_check "Markdown Lint" \
          "markdownlint website/docs 2>&1 || markdownlint-cli website/docs 2>&1" \
          || true
      else
        # Smart mode: check only changed docs
        CHANGED_DOCS=$(echo "$CHANGED_FILES" | grep "website/docs.*\.\(md\|markdown\)$" || true)
        if [[ -n "$CHANGED_DOCS" ]]; then
          run_check "Markdown Lint (changed docs only)" \
            "echo '$CHANGED_DOCS' | xargs markdownlint 2>&1 || echo '$CHANGED_DOCS' | xargs markdownlint-cli 2>&1" \
            || true
        else
          echo -e "${BLUE}▶ Markdown Lint: No changed markdown files${NC}"
          echo
        fi
      fi
    else
      echo -e "${YELLOW}⚠ Skipping Markdown lint (installation failed)${NC}"
      echo
    fi
  fi

  # 2.5 Markdown link check
  if [ "$SKIP_MARKDOWN_LINK_CHECK" = true ]; then
    echo -e "${YELLOW}⚠ Skipping Markdown link check (disabled locally, use --run-markdown-link-check to enable)${NC}"
    echo
  else
    if ! command -v markdown-link-check &> /dev/null; then
      # Check if npm is available
      if command -v npm &> /dev/null; then
        install_tool "markdown-link-check" "npm install -g markdown-link-check" "command -v markdown-link-check"
      else
        echo -e "${YELLOW}⚠ npm not found, cannot install markdown-link-check${NC}"
        echo -e "${YELLOW}  Please install Node.js and npm first${NC}"
        echo
      fi
    fi

    if command -v markdown-link-check &> /dev/null; then
      if [ "$STRICT_MODE" = true ]; then
        # CI mode: check all docs
        run_check "Markdown Link Check" \
          "find website/docs -name '*.md' -o -name '*.markdown' | xargs -I {} markdown-link-check {} --config .markdown-link-check.json --quiet" \
          || true
      else
        # Smart mode: check only changed docs
        CHANGED_DOCS=$(echo "$CHANGED_FILES" | grep "website/docs.*\.\(md\|markdown\)$" || true)
        if [[ -n "$CHANGED_DOCS" ]]; then
          run_check "Markdown Link Check (changed docs only)" \
            "echo '$CHANGED_DOCS' | xargs -I {} markdown-link-check {} --config .markdown-link-check.json --quiet" \
            || true
        else
          echo -e "${BLUE}▶ Markdown Link Check: No changed markdown files${NC}"
          echo
        fi
      fi
    else
      echo -e "${YELLOW}⚠ Skipping Markdown link check (installation failed)${NC}"
      echo
    fi
  fi

  # 2.6 Documentation content check
  echo -e "${BLUE}▶ Running documentation content checks...${NC}"
  doc_content_failed=false
  for doc in $CHANGED_FILES; do
    if [[ "$doc" == website/docs/d/* ]] || [[ "$doc" == website/docs/r/* ]]; then
      if [ -f "$doc" ]; then
        echo -e "  Checking: $doc"
        if ! go run "$SCRIPT_DIR/document/document_check.go" "$doc"; then
          doc_content_failed=true
        fi
      fi
    fi
  done
  
  if [ "$doc_content_failed" = false ]; then
    echo -e "${GREEN}✓ PASSED: Documentation content check${NC}"
    PASSED_CHECKS+=("Documentation content check")
  else
    echo -e "${RED}✗ FAILED: Documentation content check${NC}"
    FAILED_CHECKS+=("Documentation content check")
  fi
  echo

  # 2.7 Consistency check
  echo -e "${BLUE}▶ Running consistency check...${NC}"
  # Create a temporary diff file
  TEMP_DIFF=$(mktemp)
  git diff origin/master...HEAD > "$TEMP_DIFF" 2>/dev/null || git diff HEAD~1 HEAD > "$TEMP_DIFF"
  
  if [ -s "$TEMP_DIFF" ]; then
    if go run "$SCRIPT_DIR/consistency/consistency_check.go" -fileNames="$TEMP_DIFF"; then
      echo -e "${GREEN}✓ PASSED: Consistency check${NC}"
      PASSED_CHECKS+=("Consistency check")
    else
      echo -e "${RED}✗ FAILED: Consistency check${NC}"
      FAILED_CHECKS+=("Consistency check")
    fi
  else
    echo -e "${YELLOW}⚠ No changes to check${NC}"
  fi
  rm -f "$TEMP_DIFF"
  echo
fi

# ============================================================================
# 4. Build Check
# ============================================================================

if [ "$SKIP_BUILD" = false ]; then
  echo -e "${GREEN}═══════════════════════════════════════════════════════════════${NC}"
  echo -e "${GREEN}  Part 4: Build Check${NC}"
  echo -e "${GREEN}═══════════════════════════════════════════════════════════════${NC}"
  echo

  run_check "Build Check" \
    "make -C \"$PROJECT_ROOT\" build" \
    || true
else
  echo -e "${YELLOW}⚠ Skipping build check (--skip-build flag)${NC}"
  echo
fi

# ============================================================================
# 5. Unit Tests
# ============================================================================

if [ "$SKIP_TEST" = false ]; then
  echo -e "${GREEN}═══════════════════════════════════════════════════════════════${NC}"
  echo -e "${GREEN}  Part 5: Unit Tests${NC}"
  echo -e "${GREEN}═══════════════════════════════════════════════════════════════${NC}"
  echo

  run_check "Unit Tests" \
    "make -C \"$PROJECT_ROOT\" test" \
    || true
else
  echo -e "${YELLOW}⚠ Skipping unit tests (--skip-test flag)${NC}"
  echo
fi

# ============================================================================
# 6. Example Tests (Documentation Examples)
# ============================================================================

echo -e "${GREEN}═══════════════════════════════════════════════════════════════${NC}"
echo -e "${GREEN}  Part 6: Example Tests (Documentation Examples)${NC}"
echo -e "${GREEN}═══════════════════════════════════════════════════════════════${NC}"
echo

# Detect changed documentation files
CHANGED_DOCS=""
CHANGED_RESOURCES=""

# Check for changed documentation files
for file in $CHANGED_FILES; do
  if [[ "$file" == website/docs/r/*.html.markdown ]] || [[ "$file" == website/docs/d/*.html.markdown ]]; then
    CHANGED_DOCS="${CHANGED_DOCS}${file} "
  elif [[ "$file" == alicloud/resource_alicloud*.go ]] && [[ "$file" != *_test.go ]]; then
    # Convert resource file to doc file
    resource_name=$(echo "$file" | sed 's|alicloud/resource_alicloud_||' | sed 's|\.go||')
    doc_file="website/docs/r/${resource_name}.html.markdown"
    if [ -f "$doc_file" ]; then
      CHANGED_RESOURCES="${CHANGED_RESOURCES}${doc_file} "
    fi
  elif [[ "$file" == alicloud/data_source_alicloud*.go ]] && [[ "$file" != *_test.go ]]; then
    # Convert data source file to doc file
    resource_name=$(echo "$file" | sed 's|alicloud/data_source_alicloud_||' | sed 's|\.go||')
    doc_file="website/docs/d/${resource_name}.html.markdown"
    if [ -f "$doc_file" ]; then
      CHANGED_RESOURCES="${CHANGED_RESOURCES}${doc_file} "
    fi
  fi
done

# Combine and deduplicate
ALL_AFFECTED_DOCS=$(echo "$CHANGED_DOCS $CHANGED_RESOURCES" | tr ' ' '\n' | sort -u | tr '\n' ' ')

if [ -z "$ALL_AFFECTED_DOCS" ]; then
  echo -e "${BLUE}▶ No documentation changes detected${NC}"
  echo -e "${GREEN}✓ SKIPPED: Example Tests (no docs changed)${NC}"
  echo
else
  # Count examples in affected docs
  TOTAL_EXAMPLES=0
  DOCS_WITH_EXAMPLES=""

  echo -e "${BLUE}▶ Analyzing affected documentation files:${NC}"
  DOC_COUNT=$(echo "$ALL_AFFECTED_DOCS" | wc -w | tr -d ' ')
  echo -e "${BLUE}  Total documentation files to check: ${DOC_COUNT}${NC}"
  for doc_file in $ALL_AFFECTED_DOCS; do
    if [ -f "$doc_file" ]; then
      example_count=$(grep -c '```terraform' "$doc_file" 2>/dev/null || echo "0")
      if [ "$example_count" -gt 0 ]; then
        TOTAL_EXAMPLES=$((TOTAL_EXAMPLES + example_count))
        DOCS_WITH_EXAMPLES="${DOCS_WITH_EXAMPLES}${doc_file} "
        echo -e "  ${GREEN}✓${NC} $doc_file (${example_count} example(s))"
      else
        echo -e "  ${YELLOW}⚠${NC} $doc_file (no examples)"
      fi
    fi
  done
  echo

  if [ "$TOTAL_EXAMPLES" -eq 0 ]; then
    echo -e "${YELLOW}⚠ No terraform examples found in changed documentation${NC}"
    echo -e "${GREEN}✓ SKIPPED: Example Tests (no examples)${NC}"
    echo
  else
    echo -e "${BLUE}Summary:${NC}"
    echo -e "  Documentation files changed: $(echo "$ALL_AFFECTED_DOCS" | wc -w | tr -d ' ')"
    echo -e "  Files with examples: $(echo "$DOCS_WITH_EXAMPLES" | wc -w | tr -d ' ')"
    echo -e "  Total examples to test: ${TOTAL_EXAMPLES}"
    echo

    # Check if example tests should run
    if [ "$SKIP_EXAMPLE_TEST" = true ]; then
      echo -e "${YELLOW}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
      echo -e "${YELLOW}  Example tests are skipped (use default to enable)${NC}"
      echo -e "${YELLOW}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
      echo
      echo -e "${BLUE}💡 To enable example tests, run:${NC}"
      echo -e "  ${GREEN}make ci-check${NC}  (default behavior, requires credentials)"
      echo
    else
      # Check credentials
      if [ -z "$ALICLOUD_ACCESS_KEY" ] || [ -z "$ALICLOUD_SECRET_KEY" ]; then
        echo -e "${RED}✗ Cannot run example tests: AliCloud credentials not set${NC}"
        echo
        echo -e "${YELLOW}Please set the following environment variables:${NC}"
        echo -e "  ${GREEN}export ALICLOUD_ACCESS_KEY=your_access_key${NC}"
        echo -e "  ${GREEN}export ALICLOUD_SECRET_KEY=your_secret_key${NC}"
        echo -e "  ${GREEN}export ALICLOUD_REGION=cn-hangzhou${NC}  ${YELLOW}# optional${NC}"
        echo
        echo -e "${BLUE}💡 To skip example tests temporarily:${NC}"
        echo -e "  ${GREEN}make ci-check SKIP_EXAMPLE=1${NC}"
        echo
        FAILED_CHECKS+=("Example Tests (credentials not set)")
      else
        echo -e "${YELLOW}⚠ WARNING: Example tests will create REAL resources in your AliCloud account!${NC}"
        echo

        # Check if example test script exists
        if [ -f "$SCRIPT_DIR/local-example-check.sh" ]; then
          run_check "Example Tests (Documentation Examples)" \
            "\"$SCRIPT_DIR/local-example-check.sh\" --skip-build" \
            || true
        else
          echo -e "${RED}✗ Example test script not found at: $SCRIPT_DIR/local-example-check.sh${NC}"
          FAILED_CHECKS+=("Example Tests (script not found)")
          echo
        fi
      fi
    fi
  fi
fi

# ============================================================================
# Summary
# ============================================================================

echo -e "${GREEN}╔════════════════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║  Check Summary                                                 ║${NC}"
echo -e "${GREEN}╚════════════════════════════════════════════════════════════════╝${NC}"
echo

TOTAL_CHECKS=$((${#PASSED_CHECKS[@]} + ${#FAILED_CHECKS[@]}))

echo -e "${BLUE}Total Checks Run: ${TOTAL_CHECKS}${NC}"
echo -e "${GREEN}Passed: ${#PASSED_CHECKS[@]}${NC}"
echo -e "${RED}Failed: ${#FAILED_CHECKS[@]}${NC}"
echo

if [ ${#PASSED_CHECKS[@]} -gt 0 ]; then
  echo -e "${GREEN}✓ Passed Checks:${NC}"
  for check in "${PASSED_CHECKS[@]}"; do
    echo -e "  ${GREEN}✓${NC} $check"
  done
  echo
fi

if [ ${#FAILED_CHECKS[@]} -gt 0 ]; then
  echo -e "${RED}✗ Failed Checks:${NC}"
  for check in "${FAILED_CHECKS[@]}"; do
    echo -e "  ${RED}✗${NC} $check"
  done
  echo
  echo -e "${RED}════════════════════════════════════════════════════════════════${NC}"
  echo -e "${RED}  Some checks failed. Please fix the issues before pushing.${NC}"
  echo -e "${RED}════════════════════════════════════════════════════════════════${NC}"
  echo
  echo -e "${YELLOW}Tips:${NC}"
  echo -e "  • Run ${BLUE}make fmt${NC} to auto-fix formatting issues"
  echo -e "  • Run ${BLUE}make vet${NC} to check suspicious constructs"
  echo -e "  • Check the error messages above for specific issues"
  exit 1
else
  echo -e "${GREEN}════════════════════════════════════════════════════════════════${NC}"
  echo -e "${GREEN}  ✓ All checks passed! Ready to push.${NC}"
  echo -e "${GREEN}════════════════════════════════════════════════════════════════${NC}"
  echo
  
  # Check git status
  if ! git diff-index --quiet HEAD -- 2>/dev/null; then
    echo -e "${BLUE}📝 Git Status:${NC}"
    echo -e "${YELLOW}  You have uncommitted changes.${NC}"
    echo
    echo -e "${BLUE}Next steps:${NC}"
    echo -e "  1. Review your changes: ${GREEN}git status${NC}"
    echo -e "  2. Stage your changes:  ${GREEN}git add .${NC}"
    echo -e "  3. Commit your changes: ${GREEN}git commit -m 'your message'${NC}"
    echo -e "  4. Push to remote:      ${GREEN}git push${NC}"
    echo
  else
    CURRENT_BRANCH=$(git branch --show-current 2>/dev/null || echo "unknown")
    echo -e "${BLUE}📝 Git Status:${NC}"
    echo -e "  Current branch: ${GREEN}${CURRENT_BRANCH}${NC}"
    echo -e "  ${GREEN}✓ No uncommitted changes${NC}"
    echo
    if [ "$CURRENT_BRANCH" != "master" ] && [ "$CURRENT_BRANCH" != "main" ]; then
      echo -e "${BLUE}Ready to push:${NC}"
      echo -e "  ${GREEN}git push origin ${CURRENT_BRANCH}${NC}"
      echo
    fi
  fi
  
  echo -e "${BLUE}💡 Note:${NC}"
  echo -e "  • GitHub CI will run additional checks (PR title, commit count, etc.)"
  echo -e "  • Make sure your PR follows the contribution guidelines"
  echo -e "  • Acceptance tests will run automatically in CI"
  
  exit 0
fi
