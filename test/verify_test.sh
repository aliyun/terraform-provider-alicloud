#!/bin/bash

# Test harness for verify.sh credential checks
# Stubs gh to fail, verifies output and exit code

test_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
proj_root="$(cd "$test_dir/.." && pwd)"

# Create a temporary bin directory with a stub gh that fails
tmpbin=$(mktemp -d)

# Create a temporary file for output
verify_output=$(mktemp)
trap "rm -rf $tmpbin; rm -f $verify_output" EXIT

# Create stub gh that exits with non-zero
cat > "$tmpbin/gh" << 'EOF'
#!/bin/bash
exit 1
EOF
chmod +x "$tmpbin/gh"

# Run verify.sh with stubbed PATH
export PATH="$tmpbin:$PATH"

# Capture output and exit code
bash "$proj_root/bootstrap/verify.sh" > "$verify_output" 2>&1
exit_code=$?
output=$(cat "$verify_output")

echo "=== Test Output ==="
echo "$output"
echo ""
echo "=== Exit Code ==="
echo "$exit_code"
echo ""

# Check assertions
echo "=== Assertions ==="
if echo "$output" | grep -q "gh-cred"; then
    echo "FAIL: Output should not depend on ambient gh auth"
    echo "Output was: $output"
    exit 1
else
    echo "PASS: Output does not contain ambient gh auth check"
fi

if echo "$output" | grep -q "FAIL jarvis-github-token"; then
    echo "PASS: Output contains 'FAIL jarvis-github-token'"
else
    echo "FAIL: Output does NOT contain 'FAIL jarvis-github-token'"
    echo "Output was: $output"
    exit 1
fi

if [ "$exit_code" -ne 0 ]; then
    echo "PASS: Exit code is non-zero ($exit_code)"
else
    echo "FAIL: Exit code is zero (should be non-zero)"
    exit 1
fi

echo ""
echo "✓ All test assertions passed"
exit 0
