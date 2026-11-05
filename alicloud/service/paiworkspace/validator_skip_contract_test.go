package paiworkspace_test

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	product "github.com/aliyun/terraform-provider-alicloud/alicloud/service/paiworkspace"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// The shared helper caches skip=true for the lifetime of the process. Each
// mode must run in a fresh process so other tests cannot mask a regression.
func TestUnitValidatorSkipContract(t *testing.T) {
	if os.Getenv("ALICLOUD_VALIDATOR_CONTRACT_CHILD") != "1" {
		for _, mode := range []string{"false", "true"} {
			t.Run(mode, func(t *testing.T) {
				cmd := exec.Command(os.Args[0], "-test.run=^TestUnitValidatorSkipContract$", "-test.v")
				for _, v := range os.Environ() {
					if !strings.HasPrefix(v, "TF_SKIP_RESOURCE_SCHEMA_VALIDATION=") && !strings.HasPrefix(v, "ALICLOUD_VALIDATOR_CONTRACT_CHILD=") {
						cmd.Env = append(cmd.Env, v)
					}
				}
				cmd.Env = append(cmd.Env, "TF_SKIP_RESOURCE_SCHEMA_VALIDATION="+mode, "ALICLOUD_VALIDATOR_CONTRACT_CHILD=1")
				if output, err := cmd.CombinedOutput(); err != nil {
					t.Fatalf("skip=%s: %v\n%s", mode, err, output)
				}
			})
		}
		return
	}
	cases := []struct {
		factory        func() *schema.Resource
		name           string
		valid, invalid interface{}
	}{
		{product.ResourcePaiWorkspaceCodeSource, "accessibility", "PRIVATE", "__audit_invalid_enum__"},
		{product.ResourcePaiWorkspaceDataset, "accessibility", "PRIVATE", "__audit_invalid_enum__"},
		{product.ResourcePaiWorkspaceDataset, "data_source_type", "OSS", "__audit_invalid_enum__"},
		{product.ResourcePaiWorkspaceDataset, "data_type", "COMMON", "__audit_invalid_enum__"},
		{product.ResourcePaiWorkspaceDataset, "property", "FILE", "__audit_invalid_enum__"},
		{product.ResourcePaiWorkspaceDataset, "source_type", "USER", "__audit_invalid_enum__"},
		{product.ResourcePaiWorkspaceDataset, "user_id", "12345", "invalid-owner"},
		{product.ResourcePaiWorkspaceDatasetversion, "property", "FILE", "__audit_invalid_enum__"},
		{product.ResourcePaiWorkspaceExperiment, "accessibility", "PRIVATE", "__audit_invalid_enum__"},
		{product.ResourcePaiWorkspaceExperiment, "workspace_id", "12345", "invalid-owner"},
		{product.ResourcePaiWorkspaceModel, "accessibility", "PRIVATE", "__audit_invalid_enum__"},
	}
	skip := os.Getenv("TF_SKIP_RESOURCE_SCHEMA_VALIDATION") == "true"
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			validate := tc.factory().Schema[tc.name].ValidateFunc
			if validate == nil {
				t.Fatal("missing validator")
			}
			warnings, errors := validate(tc.valid, tc.name)
			if len(warnings) != 0 || len(errors) != 0 {
				t.Fatalf("valid input: warnings=%v errors=%v", warnings, errors)
			}
			warnings, errors = validate(tc.invalid, tc.name)
			if skip {
				if len(errors) != 0 || len(warnings) == 0 {
					t.Fatalf("skip=true must warn without errors: warnings=%v errors=%v", warnings, errors)
				}
			} else if len(errors) == 0 || len(warnings) != 0 {
				t.Fatalf("skip=false must reject: warnings=%v errors=%v", warnings, errors)
			}
		})
	}
}
