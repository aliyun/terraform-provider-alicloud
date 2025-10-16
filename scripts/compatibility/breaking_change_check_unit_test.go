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

// TestRetryCodeAddition tests if adding retry error codes is allowed (safe change)
func TestRetryCodeAddition(t *testing.T) {
	oldCodes := map[string]map[string]struct{}{
		"CreateInstance": {
			"Throttling": {},
		},
	}

	newCodes := map[string]map[string]struct{}{
		"CreateInstance": {
			"Throttling":         {},
			"ServiceUnavailable": {},
			"OperationConflict":  {},
		},
	}

	hasBreaking := IsRetryCodeBreaking(oldCodes, newCodes)
	if hasBreaking {
		t.Error("Expected no breaking change for retry code addition, but got one")
	} else {
		fmt.Println("✓ Test passed: Retry code addition allowed")
	}
}

// TestRetryCodeCompleteRemoval tests if complete removal of IsExpectedErrors is detected
func TestRetryCodeCompleteRemoval(t *testing.T) {
	oldCodes := map[string]map[string]struct{}{
		"CreateInstance": {
			"Throttling":         {},
			"ServiceUnavailable": {},
		},
	}

	newCodes := map[string]map[string]struct{}{
		// CreateInstance completely removed
	}

	hasBreaking := IsRetryCodeBreaking(oldCodes, newCodes)
	if !hasBreaking {
		t.Error("Expected breaking change for complete retry code removal, but got none")
	} else {
		fmt.Println("✓ Test passed: Complete retry code removal detected")
	}
}

// TestRetryCodeParsingFromContent tests parsing retry codes from actual Go code
func TestRetryCodeParsingFromContent(t *testing.T) {
	content := `
package alicloud

func resourceCreate(d *schema.ResourceData) error {
	action := "CreateInstance"
	if err := client.DoAction(action); err != nil {
		if IsExpectedErrors(err, []string{"Throttling", "ServiceUnavailable"}) {
			return resource.RetryableError(err)
		}
		return err
	}
	
	action = "DescribeInstance"
	if err := client.DoAction(action); err != nil {
		if IsExpectedErrors(err, []string{"NotFound", "InvalidId"}) {
			return resource.RetryableError(err)
		}
		return err
	}
	return nil
}
`

	codes := ParseRetryErrorCodesFromContent(content)

	// Should have 2 APIs
	if len(codes) != 2 {
		t.Errorf("Expected 2 APIs, got %d", len(codes))
		return
	}

	// Check CreateInstance codes
	if createCodes, ok := codes["CreateInstance"]; ok {
		if len(createCodes) != 2 {
			t.Errorf("Expected 2 codes for CreateInstance, got %d", len(createCodes))
		}
		if _, ok := createCodes["Throttling"]; !ok {
			t.Error("Expected Throttling code for CreateInstance")
		}
		if _, ok := createCodes["ServiceUnavailable"]; !ok {
			t.Error("Expected ServiceUnavailable code for CreateInstance")
		}
	} else {
		t.Error("CreateInstance not found in parsed codes")
	}

	// Check DescribeInstance codes
	if describeCodes, ok := codes["DescribeInstance"]; ok {
		if len(describeCodes) != 2 {
			t.Errorf("Expected 2 codes for DescribeInstance, got %d", len(describeCodes))
		}
		if _, ok := describeCodes["NotFound"]; !ok {
			t.Error("Expected NotFound code for DescribeInstance")
		}
		if _, ok := describeCodes["InvalidId"]; !ok {
			t.Error("Expected InvalidId code for DescribeInstance")
		}
	} else {
		t.Error("DescribeInstance not found in parsed codes")
	}

	fmt.Println("✓ Test passed: Retry code parsing from content")
}
