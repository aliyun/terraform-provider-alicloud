#!/bin/bash

# Count auto-generated resource files
# Criteria:
# 1. In alicloud directory
# 2. Starts with "resource"
# 3. Does not end with "_test.go"
# 4. Contains "This file is generated automatically"

count=0
echo "Counting auto-generated resource files..."
echo

# Iterate through all matching files
for file in alicloud/resource*.go; do
    # Exclude files ending with _test.go
    if [[ "$file" == *_test.go ]]; then
        continue
    fi
    
    # Check if file contains auto-generation marker
    if grep -q "This file is generated automatically" "$file"; then
        ((count++))
        echo "✓ $file"
    fi
done

echo
echo "==============================================="
echo "Total auto-generated resource files: $count"
echo "==============================================="

