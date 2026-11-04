package cs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestUnitAliCloudCSServerlessKubernetesEncryptionSchema(t *testing.T) {
	r := ResourceCSServerlessKubernetes()

	keySchema, ok := r.Schema["encryption_provider_key"]
	if !ok {
		t.Fatal("expected schema to contain 'encryption_provider_key'")
	}
	if keySchema.Type != schema.TypeString {
		t.Errorf("encryption_provider_key type = %v, want TypeString", keySchema.Type)
	}
	if !keySchema.Optional {
		t.Error("encryption_provider_key should be Optional")
	}
	if keySchema.DiffSuppressFunc == nil {
		t.Error("encryption_provider_key should set DiffSuppressFunc")
	}

	disableSchema, ok := r.Schema["disable_encryption"]
	if !ok {
		t.Fatal("expected schema to contain 'disable_encryption'")
	}
	if disableSchema.Type != schema.TypeBool {
		t.Errorf("disable_encryption type = %v, want TypeBool", disableSchema.Type)
	}
	if !disableSchema.Optional || !disableSchema.Computed {
		t.Error("disable_encryption should be Optional and Computed")
	}
}

func TestUnitAliCloudCSServerlessKubernetesKmsEncryptionDiffSuppress(t *testing.T) {
	r := ResourceCSServerlessKubernetes()

	// When disable_encryption is true, a diff on encryption_provider_key is suppressed.
	dDisabled := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{
		"disable_encryption": true,
	})
	if !kmsEncryptionDiffSuppressFunc("encryption_provider_key", "key-old", "", dDisabled) {
		t.Error("expected diff to be suppressed when disable_encryption=true")
	}

	// When disable_encryption is false, the diff is not suppressed.
	dEnabled := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{
		"disable_encryption": false,
	})
	if kmsEncryptionDiffSuppressFunc("encryption_provider_key", "", "key-new", dEnabled) {
		t.Error("expected diff not to be suppressed when disable_encryption=false")
	}
}
