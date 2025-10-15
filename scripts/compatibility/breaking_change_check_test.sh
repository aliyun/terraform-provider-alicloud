#!/bin/bash

# Breaking Change Check Test Suite
# This script tests all scenarios of breaking changes

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TEST_DIR="/tmp/breaking_change_test_$$"
PASS_COUNT=0
FAIL_COUNT=0

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

mkdir -p "$TEST_DIR"

# Cleanup function
cleanup() {
    rm -rf "$TEST_DIR"
}
trap cleanup EXIT

# Test result tracking
test_passed() {
    echo -e "${GREEN}✓ PASS${NC}: $1"
    ((PASS_COUNT++))
}

test_failed() {
    echo -e "${RED}✗ FAIL${NC}: $1"
    echo -e "${YELLOW}  Expected: Breaking Change${NC}"
    echo -e "${YELLOW}  Got: No Breaking Change${NC}"
    ((FAIL_COUNT++))
}

test_expected_pass() {
    echo -e "${GREEN}✓ PASS${NC}: $1 (correctly allowed)"
    ((PASS_COUNT++))
}

test_unexpected_fail() {
    echo -e "${RED}✗ FAIL${NC}: $1"
    echo -e "${YELLOW}  Expected: No Breaking Change${NC}"
    echo -e "${YELLOW}  Got: Breaking Change${NC}"
    ((FAIL_COUNT++))
}

# Run a test
run_test() {
    local test_name="$1"
    local old_content="$2"
    local new_content="$3"
    local should_break="$4"  # "yes" or "no"
    
    echo
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${BLUE}Testing: $test_name${NC}"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    
    # Create test files
    local old_file="$TEST_DIR/resource_test_old.go"
    local new_file="$TEST_DIR/resource_test_new.go"
    
    echo "$old_content" > "$old_file"
    echo "$new_content" > "$new_file"
    
    # Create a mock git repository for testing
    (
        cd "$TEST_DIR"
        git init -q
        git config user.email "test@test.com"
        git config user.name "Test"
        cp "$old_file" alicloud/
        mkdir -p alicloud
        cp "$old_file" alicloud/resource_test.go
        git add .
        git commit -q -m "Initial"
        cp "$new_file" alicloud/resource_test.go
        git add .
        git commit -q -m "Update"
        git diff HEAD^ HEAD > diff.out
    )
    
    # Run breaking change check
    local output
    output=$(cd "$TEST_DIR" && go run "$SCRIPT_DIR/breaking_change_check.go" -fileNames=diff.out 2>&1 || true)
    
    # Check result
    if echo "$output" | grep -q "Breaking Change"; then
        if [ "$should_break" = "yes" ]; then
            test_passed "$test_name"
            echo "$output" | grep "Breaking Change" | head -3
        else
            test_unexpected_fail "$test_name"
            echo "$output" | grep "Breaking Change"
        fi
    else
        if [ "$should_break" = "no" ]; then
            test_expected_pass "$test_name"
        else
            test_failed "$test_name"
        fi
    fi
}

# Test 1: Attribute Deletion
echo
echo -e "${YELLOW}═══════════════════════════════════════════${NC}"
echo -e "${YELLOW}  Test Suite: Schema Breaking Changes${NC}"
echo -e "${YELLOW}═══════════════════════════════════════════${NC}"

run_test "Attribute Deletion" \
'package alicloud
func resource() *schema.Resource {
    return &schema.Resource{
        Schema: map[string]*schema.Schema{
            "field1": {
                Type:     schema.TypeString,
                Optional: true,
            },
            "field2": {
                Type:     schema.TypeString,
                Optional: true,
            },
        },
    }
}' \
'package alicloud
func resource() *schema.Resource {
    return &schema.Resource{
        Schema: map[string]*schema.Schema{
            "field1": {
                Type:     schema.TypeString,
                Optional: true,
            },
        },
    }
}' \
"yes"

# Test 2: New Required Field
run_test "New Required Field Added" \
'package alicloud
func resource() *schema.Resource {
    return &schema.Resource{
        Schema: map[string]*schema.Schema{
            "field1": {
                Type:     schema.TypeString,
                Optional: true,
            },
        },
    }
}' \
'package alicloud
func resource() *schema.Resource {
    return &schema.Resource{
        Schema: map[string]*schema.Schema{
            "field1": {
                Type:     schema.TypeString,
                Optional: true,
            },
            "new_required_field": {
                Type:     schema.TypeString,
                Required: true,
            },
        },
    }
}' \
"yes"

# Test 3: New Optional Field (Should Pass)
run_test "New Optional Field Added (Safe)" \
'package alicloud
func resource() *schema.Resource {
    return &schema.Resource{
        Schema: map[string]*schema.Schema{
            "field1": {
                Type:     schema.TypeString,
                Optional: true,
            },
        },
    }
}' \
'package alicloud
func resource() *schema.Resource {
    return &schema.Resource{
        Schema: map[string]*schema.Schema{
            "field1": {
                Type:     schema.TypeString,
                Optional: true,
            },
            "new_optional_field": {
                Type:     schema.TypeString,
                Optional: true,
            },
        },
    }
}' \
"no"

# Test 4: Optional -> Required
run_test "Optional Changed to Required" \
'package alicloud
func resource() *schema.Resource {
    return &schema.Resource{
        Schema: map[string]*schema.Schema{
            "field1": {
                Type:     schema.TypeString,
                Optional: true,
            },
        },
    }
}' \
'package alicloud
func resource() *schema.Resource {
    return &schema.Resource{
        Schema: map[string]*schema.Schema{
            "field1": {
                Type:     schema.TypeString,
                Required: true,
            },
        },
    }
}' \
"yes"

# Test 5: Type Change
run_test "Type Changed (TypeInt -> TypeString)" \
'package alicloud
func resource() *schema.Resource {
    return &schema.Resource{
        Schema: map[string]*schema.Schema{
            "count_field": {
                Type:     schema.TypeInt,
                Optional: true,
            },
        },
    }
}' \
'package alicloud
func resource() *schema.Resource {
    return &schema.Resource{
        Schema: map[string]*schema.Schema{
            "count_field": {
                Type:     schema.TypeString,
                Optional: true,
            },
        },
    }
}' \
"yes"

# Test 6: Adding ForceNew
run_test "Non-ForceNew -> ForceNew" \
'package alicloud
func resource() *schema.Resource {
    return &schema.Resource{
        Schema: map[string]*schema.Schema{
            "field1": {
                Type:     schema.TypeString,
                Optional: true,
            },
        },
    }
}' \
'package alicloud
func resource() *schema.Resource {
    return &schema.Resource{
        Schema: map[string]*schema.Schema{
            "field1": {
                Type:     schema.TypeString,
                Optional: true,
                ForceNew: true,
            },
        },
    }
}' \
"yes"

# Test 7: Multiple Safe Changes (Should Pass)
run_test "Multiple Safe Changes" \
'package alicloud
func resource() *schema.Resource {
    return &schema.Resource{
        Schema: map[string]*schema.Schema{
            "field1": {
                Type:     schema.TypeString,
                Required: true,
            },
        },
    }
}' \
'package alicloud
func resource() *schema.Resource {
    return &schema.Resource{
        Schema: map[string]*schema.Schema{
            "field1": {
                Type:     schema.TypeString,
                Required: true,
            },
            "field2": {
                Type:     schema.TypeString,
                Optional: true,
            },
            "field3": {
                Type:     schema.TypeInt,
                Optional: true,
                Computed: true,
            },
        },
    }
}' \
"no"

# Print Summary
echo
echo -e "${YELLOW}═══════════════════════════════════════════${NC}"
echo -e "${YELLOW}  Test Summary${NC}"
echo -e "${YELLOW}═══════════════════════════════════════════${NC}"
echo -e "${GREEN}Passed: $PASS_COUNT${NC}"
echo -e "${RED}Failed: $FAIL_COUNT${NC}"
echo -e "${BLUE}Total:  $((PASS_COUNT + FAIL_COUNT))${NC}"
echo

if [ $FAIL_COUNT -eq 0 ]; then
    echo -e "${GREEN}✅ All tests passed!${NC}"
    exit 0
else
    echo -e "${RED}❌ Some tests failed!${NC}"
    exit 1
fi

