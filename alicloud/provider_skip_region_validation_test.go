package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// skipRegionValidationTestSchema returns a single-field schema map containing
// only the provider-level skip_region_validation attribute, so the SDK diff
// machinery (which applies DefaultFunc) can be exercised in isolation.
func skipRegionValidationTestSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"skip_region_validation": Provider().(*schema.Provider).Schema["skip_region_validation"],
	}
}

func TestProviderSkipRegionValidation(t *testing.T) {
	t.Run("env unset defaults to false", func(t *testing.T) {
		t.Setenv("ALIBABA_CLOUD_SKIP_REGION_VALIDATION", "")
		t.Setenv("ALICLOUD_SKIP_REGION_VALIDATION", "")
		d := schema.TestResourceDataRaw(t, skipRegionValidationTestSchema(), map[string]interface{}{})
		if got := d.Get("skip_region_validation").(bool); got != false {
			t.Fatalf("expected false when no env is set, got %v", got)
		}
	})

	t.Run("ALICLOUD env true enables skip", func(t *testing.T) {
		t.Setenv("ALIBABA_CLOUD_SKIP_REGION_VALIDATION", "")
		t.Setenv("ALICLOUD_SKIP_REGION_VALIDATION", "true")
		d := schema.TestResourceDataRaw(t, skipRegionValidationTestSchema(), map[string]interface{}{})
		if got := d.Get("skip_region_validation").(bool); got != true {
			t.Fatalf("expected true from ALICLOUD_SKIP_REGION_VALIDATION, got %v", got)
		}
	})

	t.Run("ALIBABA_CLOUD takes priority over ALICLOUD", func(t *testing.T) {
		t.Setenv("ALIBABA_CLOUD_SKIP_REGION_VALIDATION", "true")
		t.Setenv("ALICLOUD_SKIP_REGION_VALIDATION", "false")
		d := schema.TestResourceDataRaw(t, skipRegionValidationTestSchema(), map[string]interface{}{})
		if got := d.Get("skip_region_validation").(bool); got != true {
			t.Fatalf("expected ALIBABA_CLOUD_SKIP_REGION_VALIDATION to win (true), got %v", got)
		}
	})

	t.Run("explicit config value overrides env", func(t *testing.T) {
		t.Setenv("ALIBABA_CLOUD_SKIP_REGION_VALIDATION", "true")
		t.Setenv("ALICLOUD_SKIP_REGION_VALIDATION", "true")
		d := schema.TestResourceDataRaw(t, skipRegionValidationTestSchema(), map[string]interface{}{"skip_region_validation": false})
		if got := d.Get("skip_region_validation").(bool); got != false {
			t.Fatalf("expected explicit config value false to override env, got %v", got)
		}
	})

	t.Run("invalid boolean env value returns an error", func(t *testing.T) {
		t.Setenv("ALIBABA_CLOUD_SKIP_REGION_VALIDATION", "")
		t.Setenv("ALICLOUD_SKIP_REGION_VALIDATION", "notabool")
		defaultFunc := Provider().(*schema.Provider).Schema["skip_region_validation"].DefaultFunc
		if defaultFunc == nil {
			t.Fatal("expected skip_region_validation to declare a DefaultFunc")
		}
		if _, err := defaultFunc(); err == nil {
			t.Fatal("expected an error for an invalid boolean env value, got nil")
		}
	})
}
