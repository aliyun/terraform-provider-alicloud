package fwutil

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestUnitTimeoutBlockMatchesOperations(t *testing.T) {
	block := TimeoutBlock("create", "delete")
	if len(block.Attributes) != 2 {
		t.Fatalf("TimeoutBlock returned %d attributes, want 2", len(block.Attributes))
	}
	for _, operation := range []string{"create", "delete"} {
		attribute, ok := block.Attributes[operation]
		if !ok {
			t.Fatalf("TimeoutBlock is missing %q", operation)
		}
		stringAttribute, ok := attribute.(schema.StringAttribute)
		if !ok {
			t.Fatalf("%q is not a string attribute", operation)
		}
		if !stringAttribute.Optional || stringAttribute.Required {
			t.Fatalf("%q must be optional, got Optional=%v Required=%v", operation, stringAttribute.Optional, stringAttribute.Required)
		}
		if len(stringAttribute.Validators) != 1 {
			t.Fatalf("%q must carry the duration validator, got %d", operation, len(stringAttribute.Validators))
		}
	}

	attributeTypes := TimeoutTypes("create", "delete")
	if len(attributeTypes) != 2 || attributeTypes["create"] != types.StringType || attributeTypes["delete"] != types.StringType {
		t.Fatalf("TimeoutTypes = %v, want create/delete string types", attributeTypes)
	}
}

func TestUnitTimeoutDuration(t *testing.T) {
	attributeTypes := TimeoutTypes("create", "delete")
	object := func(create, delete attr.Value) types.Object {
		value, diagnostics := types.ObjectValue(attributeTypes, map[string]attr.Value{"create": create, "delete": delete})
		if diagnostics.HasError() {
			t.Fatalf("building timeouts object: %v", diagnostics)
		}
		return value
	}

	if duration, diagnostics := TimeoutDuration(types.ObjectNull(attributeTypes), "create", 20*time.Minute); diagnostics.HasError() || duration != 20*time.Minute {
		t.Fatalf("null block: duration=%v diagnostics=%v", duration, diagnostics)
	}
	if duration, diagnostics := TimeoutDuration(object(types.StringNull(), types.StringNull()), "create", 20*time.Minute); diagnostics.HasError() || duration != 20*time.Minute {
		t.Fatalf("unset operation: duration=%v diagnostics=%v", duration, diagnostics)
	}
	if duration, diagnostics := TimeoutDuration(object(types.StringValue("10m"), types.StringNull()), "create", 20*time.Minute); diagnostics.HasError() || duration != 10*time.Minute {
		t.Fatalf("configured operation: duration=%v diagnostics=%v", duration, diagnostics)
	}
	if _, diagnostics := TimeoutDuration(object(types.StringValue("soon"), types.StringNull()), "create", 20*time.Minute); !diagnostics.HasError() {
		t.Fatal("an unparsable duration must be rejected")
	}
	if _, diagnostics := TimeoutDuration(object(types.StringUnknown(), types.StringNull()), "create", 20*time.Minute); !diagnostics.HasError() {
		t.Fatal("an unknown duration must be rejected")
	}
	if _, diagnostics := TimeoutDuration(types.ObjectUnknown(attributeTypes), "create", 20*time.Minute); !diagnostics.HasError() {
		t.Fatal("an unknown timeouts block must be rejected")
	}
}
