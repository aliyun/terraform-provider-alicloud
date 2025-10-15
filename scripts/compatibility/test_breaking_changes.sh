#!/bin/bash

# Simplified Breaking Change Test Suite
# Tests all breaking change scenarios using pre-generated diff files

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TEST_DIR="/tmp/bc_test_$$"
PASS_COUNT=0
FAIL_COUNT=0

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

mkdir -p "$TEST_DIR"

cleanup() {
    rm -rf "$TEST_DIR"
}
trap cleanup EXIT

test_passed() {
    echo -e "${GREEN}✓ PASS${NC}: $1"
    ((PASS_COUNT++))
}

test_failed() {
    echo -e "${RED}✗ FAIL${NC}: $1"
    ((FAIL_COUNT++))
}

echo -e "${YELLOW}═══════════════════════════════════════════${NC}"
echo -e "${YELLOW}  Breaking Change Detection Test Suite${NC}"
echo -e "${YELLOW}═══════════════════════════════════════════${NC}"

# Test 1: New Required Field
echo
echo -e "${BLUE}Test 1: New Required Field Detection${NC}"
cat > "$TEST_DIR/test1.diff" << 'EOF'
diff --git a/alicloud/resource_alicloud_test.go b/alicloud/resource_alicloud_test.go
index abc123..def456 100644
--- a/alicloud/resource_alicloud_test.go
+++ b/alicloud/resource_alicloud_test.go
@@ -10,6 +10,10 @@ func resourceTest() *schema.Resource {
 				Type:     schema.TypeString,
 				Optional: true,
 			},
+			"new_required_field": {
+				Type:     schema.TypeBool,
+				Required: true,
+			},
 		},
 	}
 }
EOF

if go run "$SCRIPT_DIR/breaking_change_check.go" -fileNames="$TEST_DIR/test1.diff" 2>&1 | grep -q "New required attribute.*should not been added"; then
    test_passed "New Required Field Detection"
else
    test_failed "New Required Field Detection - should detect new required field as breaking change"
fi

# Test 2: Attribute Deletion
echo
echo -e "${BLUE}Test 2: Attribute Deletion Detection${NC}"
cat > "$TEST_DIR/test2.diff" << 'EOF'
diff --git a/alicloud/resource_alicloud_test.go b/alicloud/resource_alicloud_test.go
index abc123..def456 100644
--- a/alicloud/resource_alicloud_test.go
+++ b/alicloud/resource_alicloud_test.go
@@ -10,10 +10,6 @@ func resourceTest() *schema.Resource {
 				Type:     schema.TypeString,
 				Optional: true,
 			},
-			"removed_field": {
-				Type:     schema.TypeString,
-				Optional: true,
-			},
 		},
 	}
 }
EOF

if go run "$SCRIPT_DIR/breaking_change_check.go" -fileNames="$TEST_DIR/test2.diff" 2>&1 | grep -q "Attribute.*should not been removed"; then
    test_passed "Attribute Deletion Detection"
else
    test_failed "Attribute Deletion Detection"
fi

# Test 3: Type Change
echo
echo -e "${BLUE}Test 3: Type Change Detection (TypeInt -> TypeString)${NC}"
cat > "$TEST_DIR/test3.diff" << 'EOF'
diff --git a/alicloud/resource_alicloud_test.go b/alicloud/resource_alicloud_test.go
index abc123..def456 100644
--- a/alicloud/resource_alicloud_test.go
+++ b/alicloud/resource_alicloud_test.go
@@ -8,7 +8,7 @@ func resourceTest() *schema.Resource {
 		Schema: map[string]*schema.Schema{
 			"count_field": {
-				Type:     schema.TypeInt,
+				Type:     schema.TypeString,
 				Optional: true,
 			},
 		},
EOF

if go run "$SCRIPT_DIR/breaking_change_check.go" -fileNames="$TEST_DIR/test3.diff" 2>&1 | grep -q "type should not been changed"; then
    test_passed "Type Change Detection"
else
    test_failed "Type Change Detection"
fi

# Test 4: Optional -> Required
echo
echo -e "${BLUE}Test 4: Optional -> Required Detection${NC}"
cat > "$TEST_DIR/test4.diff" << 'EOF'
diff --git a/alicloud/resource_alicloud_test.go b/alicloud/resource_alicloud_test.go
index abc123..def456 100644
--- a/alicloud/resource_alicloud_test.go
+++ b/alicloud/resource_alicloud_test.go
@@ -8,7 +8,7 @@ func resourceTest() *schema.Resource {
 		Schema: map[string]*schema.Schema{
 			"field1": {
 				Type:     schema.TypeString,
-				Optional: true,
+				Required: true,
 			},
 		},
 	}
EOF

if go run "$SCRIPT_DIR/breaking_change_check.go" -fileNames="$TEST_DIR/test4.diff" 2>&1 | grep -q "should not been changed.*to required"; then
    test_passed "Optional -> Required Detection"
else
    test_failed "Optional -> Required Detection"
fi

# Test 5: ForceNew Added
echo
echo -e "${BLUE}Test 5: ForceNew Added Detection${NC}"
cat > "$TEST_DIR/test5.diff" << 'EOF'
diff --git a/alicloud/resource_alicloud_test.go b/alicloud/resource_alicloud_test.go
index abc123..def456 100644
--- a/alicloud/resource_alicloud_test.go
+++ b/alicloud/resource_alicloud_test.go
@@ -9,6 +9,7 @@ func resourceTest() *schema.Resource {
 			"field1": {
 				Type:     schema.TypeString,
 				Optional: true,
+				ForceNew: true,
 			},
 		},
 	}
EOF

if go run "$SCRIPT_DIR/breaking_change_check.go" -fileNames="$TEST_DIR/test5.diff" 2>&1 | grep -q "should not been changed to ForceNew"; then
    test_passed "ForceNew Added Detection"
else
    test_failed "ForceNew Added Detection"
fi

# Test 6: Retry Error Code Removed
echo
echo -e "${BLUE}Test 6: Retry Error Code Removal Detection${NC}"
cat > "$TEST_DIR/test6.diff" << 'EOF'
diff --git a/alicloud/resource_alicloud_test.go b/alicloud/resource_alicloud_test.go
index abc123..def456 100644
--- a/alicloud/resource_alicloud_test.go
+++ b/alicloud/resource_alicloud_test.go
@@ -10,7 +10,7 @@ func resourceTest() {
 	action := "CreateInstance"
 	response, err := client.DoRequest(action, request)
 	if err != nil {
-		if IsExpectedErrors(err, []string{"Throttling", "ServiceUnavailable", "OperationConflict"}) {
+		if IsExpectedErrors(err, []string{"Throttling", "ServiceUnavailable"}) {
 			return resource.RetryableError(err)
 		}
 		return resource.NonRetryableError(err)
EOF

if go run "$SCRIPT_DIR/breaking_change_check.go" -fileNames="$TEST_DIR/test6.diff" 2>&1 | grep -q "Retry error code.*should not been removed"; then
    test_passed "Retry Error Code Removal Detection"
else
    test_failed "Retry Error Code Removal Detection"
fi

# Test 7: Safe Change - New Optional Field
echo
echo -e "${BLUE}Test 7: Safe Change - New Optional Field (Should Pass)${NC}"
cat > "$TEST_DIR/test7.diff" << 'EOF'
diff --git a/alicloud/resource_alicloud_test.go b/alicloud/resource_alicloud_test.go
index abc123..def456 100644
--- a/alicloud/resource_alicloud_test.go
+++ b/alicloud/resource_alicloud_test.go
@@ -10,6 +10,10 @@ func resourceTest() *schema.Resource {
 				Type:     schema.TypeString,
 				Optional: true,
 			},
+			"new_optional_field": {
+				Type:     schema.TypeString,
+				Optional: true,
+			},
 		},
 	}
 }
EOF

if go run "$SCRIPT_DIR/breaking_change_check.go" -fileNames="$TEST_DIR/test7.diff" 2>&1 | grep -q "PASS"; then
    test_passed "Safe Change - New Optional Field"
else
    test_failed "Safe Change - should allow new optional field"
fi

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
    echo -e "${RED}❌ $FAIL_COUNT test(s) failed!${NC}"
    exit 1
fi

