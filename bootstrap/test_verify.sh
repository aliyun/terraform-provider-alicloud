#!/bin/bash

# Test harness for verify.sh credential checks
# Stubs gh to fail, verifies output and exit code

set -e

test_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$test_dir"

# Create a temporary bin directory with a stub gh that fails
tmpbin=$(mktemp -d)
trap "rm -rf $tmpbin" EXIT

# Create stub gh that exits with non-zero
cat > "$tmpbin/gh" << 'EOF'
#!/bin/bash
exit 1
EOF
chmod +x "$tmpbin/gh"

# Run verify.sh with stubbed PATH
export PATH="$tmpbin:$PATH"

# Capture output and exit code
bash verify.sh > /tmp/verify_output.txt 2>&1
exit_code=$?
output=$(cat /tmp/verify_output.txt)

echo "=== Test Output ==="
echo "$output"
echo ""
echo "=== Exit Code ==="
echo "$exit_code"
echo ""

# Check assertions
echo "=== Assertions ==="
if echo "$output" | grep -q "FAIL gh-cred"; then
    echo "PASS: Output contains 'FAIL gh-cred'"
else
    echo "FAIL: Output does NOT contain 'FAIL gh-cred'"
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
