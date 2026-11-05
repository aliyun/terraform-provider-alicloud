package dts_test

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	product "github.com/aliyun/terraform-provider-alicloud/alicloud/service/dts"
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
		{product.ResourceInstance, "destination_endpoint_engine_name", "MySQL", "__audit_invalid_enum__"},
		{product.ResourceInstance, "du", 1, 101},
		{product.ResourceInstance, "fee_type", "ONLY_CONFIGURATION_FEE", "__audit_invalid_enum__"},
		{product.ResourceInstance, "instance_class", "4xlarge", "__audit_invalid_enum__"},
		{product.ResourceInstance, "payment_type", "PayAsYouGo", "__audit_invalid_enum__"},
		{product.ResourceInstance, "period", "Month", "__audit_invalid_enum__"},
		{product.ResourceInstance, "source_endpoint_engine_name", "MySQL", "__audit_invalid_enum__"},
		{product.ResourceInstance, "sync_architecture", "bidirectional", "__audit_invalid_enum__"},
		{product.ResourceInstance, "synchronization_direction", "Forward", "__audit_invalid_enum__"},
		{product.ResourceInstance, "type", "migration", "__audit_invalid_enum__"},
		{product.ResourceInstance, "used_time", 1, 10},
		{product.ResourceJobMonitorRule, "state", "Y", "__audit_invalid_enum__"},
		{product.ResourceJobMonitorRule, "type", "delay", "__audit_invalid_enum__"},
		{product.ResourceSubscriptionJob, "destination_endpoint_engine_name", "ADS", "__audit_invalid_enum__"},
		{product.ResourceSubscriptionJob, "instance_class", "large", "__audit_invalid_enum__"},
		{product.ResourceSubscriptionJob, "payment_duration_unit", "Month", "__audit_invalid_enum__"},
		{product.ResourceSubscriptionJob, "payment_type", "PayAsYouGo", "__audit_invalid_enum__"},
		{product.ResourceSubscriptionJob, "source_endpoint_engine_name", "MySQL", "__audit_invalid_enum__"},
		{product.ResourceSubscriptionJob, "source_endpoint_instance_type", "CEN", "__audit_invalid_enum__"},
		{product.ResourceSubscriptionJob, "source_endpoint_ssl", "0", "__audit_invalid_enum__"},
		{product.ResourceSubscriptionJob, "status", "Abnormal", "__audit_invalid_enum__"},
		{product.ResourceSubscriptionJob, "subscription_instance_network_type", "classic", "__audit_invalid_enum__"},
		{product.ResourceSubscriptionJob, "sync_architecture", "bidirectional", "__audit_invalid_enum__"},
		{product.ResourceSynchronizationJob, "destination_endpoint_engine_name", "ADB20", "__audit_invalid_enum__"},
		{product.ResourceSynchronizationJob, "destination_endpoint_instance_type", "ads", "__audit_invalid_enum__"},
		{product.ResourceSynchronizationJob, "destination_endpoint_ssl", "0", "__audit_invalid_enum__"},
		{product.ResourceSynchronizationJob, "instance_class", "xmicro", "__audit_invalid_enum__"},
		{product.ResourceSynchronizationJob, "source_endpoint_engine_name", "AS400", "__audit_invalid_enum__"},
		{product.ResourceSynchronizationJob, "source_endpoint_instance_type", "CEN", "__audit_invalid_enum__"},
		{product.ResourceSynchronizationJob, "source_endpoint_ssl", "0", "__audit_invalid_enum__"},
		{product.ResourceSynchronizationJob, "status", "Synchronizing", "__audit_invalid_enum__"},
		{product.ResourceSynchronizationJob, "synchronization_direction", "Forward", "__audit_invalid_enum__"},
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
