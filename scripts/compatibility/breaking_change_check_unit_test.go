package main

import (
	"fmt"
	"testing"
)

// TestNewRequiredFieldDetection tests if we can detect new required fields
func TestNewRequiredFieldDetection(t *testing.T) {
	oldAttrs := map[string]map[string]interface{}{
		"field1": {
			"Name":     "field1",
			"Type":     "TypeString",
			"Optional": true,
		},
	}

	newAttrs := map[string]map[string]interface{}{
		"field1": {
			"Name":     "field1",
			"Type":     "TypeString",
			"Optional": true,
		},
		"new_required_field": {
			"Name":     "new_required_field",
			"Type":     "TypeBool",
			"Required": true,
		},
	}

	hasBreaking := IsBreakingChange(oldAttrs, newAttrs)
	if !hasBreaking {
		t.Error("Expected breaking change for new required field, but got none")
	} else {
		fmt.Println("✓ Test passed: New required field detected")
	}
}

// TestAttributeDeletion tests if we can detect attribute deletion
func TestAttributeDeletion(t *testing.T) {
	oldAttrs := map[string]map[string]interface{}{
		"field1": {
			"Name":     "field1",
			"Type":     "TypeString",
			"Optional": true,
		},
		"field2": {
			"Name":     "field2",
			"Type":     "TypeString",
			"Optional": true,
		},
	}

	newAttrs := map[string]map[string]interface{}{
		"field1": {
			"Name":     "field1",
			"Type":     "TypeString",
			"Optional": true,
		},
	}

	hasBreaking := IsBreakingChange(oldAttrs, newAttrs)
	if !hasBreaking {
		t.Error("Expected breaking change for attribute deletion, but got none")
	} else {
		fmt.Println("✓ Test passed: Attribute deletion detected")
	}
}

// TestOptionalToRequired tests if we can detect optional->required change
func TestOptionalToRequired(t *testing.T) {
	oldAttrs := map[string]map[string]interface{}{
		"field1": {
			"Name":     "field1",
			"Type":     "TypeString",
			"Optional": true,
		},
	}

	newAttrs := map[string]map[string]interface{}{
		"field1": {
			"Name":     "field1",
			"Type":     "TypeString",
			"Required": true,
		},
	}

	hasBreaking := IsBreakingChange(oldAttrs, newAttrs)
	if !hasBreaking {
		t.Error("Expected breaking change for optional->required, but got none")
	} else {
		fmt.Println("✓ Test passed: Optional->Required detected")
	}
}

// TestTypeChange tests if we can detect type changes
func TestTypeChange(t *testing.T) {
	oldAttrs := map[string]map[string]interface{}{
		"count_field": {
			"Name":     "count_field",
			"Type":     "TypeInt",
			"Optional": true,
		},
	}

	newAttrs := map[string]map[string]interface{}{
		"count_field": {
			"Name":     "count_field",
			"Type":     "TypeString",
			"Optional": true,
		},
	}

	hasBreaking := IsBreakingChange(oldAttrs, newAttrs)
	if !hasBreaking {
		t.Error("Expected breaking change for type change, but got none")
	} else {
		fmt.Println("✓ Test passed: Type change detected")
	}
}

// TestForceNewAdded tests if we can detect ForceNew being added
func TestForceNewAdded(t *testing.T) {
	oldAttrs := map[string]map[string]interface{}{
		"field1": {
			"Name":     "field1",
			"Type":     "TypeString",
			"Optional": true,
		},
	}

	newAttrs := map[string]map[string]interface{}{
		"field1": {
			"Name":     "field1",
			"Type":     "TypeString",
			"Optional": true,
			"ForceNew": true,
		},
	}

	hasBreaking := IsBreakingChange(oldAttrs, newAttrs)
	if !hasBreaking {
		t.Error("Expected breaking change for ForceNew added, but got none")
	} else {
		fmt.Println("✓ Test passed: ForceNew addition detected")
	}
}

// TestSafeChange tests if safe changes are allowed
func TestSafeChange(t *testing.T) {
	oldAttrs := map[string]map[string]interface{}{
		"field1": {
			"Name":     "field1",
			"Type":     "TypeString",
			"Required": true,
		},
	}

	newAttrs := map[string]map[string]interface{}{
		"field1": {
			"Name":     "field1",
			"Type":     "TypeString",
			"Required": true,
		},
		"new_optional_field": {
			"Name":     "new_optional_field",
			"Type":     "TypeString",
			"Optional": true,
		},
	}

	hasBreaking := IsBreakingChange(oldAttrs, newAttrs)
	if hasBreaking {
		t.Error("Expected no breaking change for new optional field, but got one")
	} else {
		fmt.Println("✓ Test passed: New optional field allowed")
	}
}

// TestRetryCodeRemoval tests if retry error code removal is detected
func TestRetryCodeRemoval(t *testing.T) {
	oldCodes := map[string]map[string]struct{}{
		"CreateInstance": {
			"Throttling":         {},
			"ServiceUnavailable": {},
			"OperationConflict":  {},
		},
	}

	newCodes := map[string]map[string]struct{}{
		"CreateInstance": {
			"Throttling":         {},
			"ServiceUnavailable": {},
		},
	}

	hasBreaking := IsRetryCodeBreaking(oldCodes, newCodes)
	if !hasBreaking {
		t.Error("Expected breaking change for retry code removal, but got none")
	} else {
		fmt.Println("✓ Test passed: Retry code removal detected")
	}
}

