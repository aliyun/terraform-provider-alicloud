package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCREndpointAclPoliciesDataSource(t *testing.T) {

	rand := acctest.RandIntRange(10000, 99999)
	resourceId := "data.alicloud_cr_endpoint_acl_policies.default"
	name := fmt.Sprintf("tf-testacc-cr-aclp-%v", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceCrEndpointAclPoliciesConfigDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id":   "${alicloud_cr_endpoint_acl_policy.default.instance_id}",
			"endpoint_type": "internet",
			"module_name":   "Registry",
			"ids":           []string{"${alicloud_cr_endpoint_acl_policy.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id":   "${alicloud_cr_endpoint_acl_policy.default.instance_id}",
			"endpoint_type": "internet",
			"ids":           []string{"${alicloud_cr_endpoint_acl_policy.default.id}_fake"},
		}),
	}

	var existCrEndpointAclPoliciesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                    "1",
			"policies.#":               "1",
			"policies.0.endpoint_type": "internet",
		}
	}

	var fakeCrEndpointAclPoliciesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"policies.#": "0",
		}
	}

	var CrEndpointAclPoliciesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existCrEndpointAclPoliciesMapFunc,
		fakeMapFunc:  fakeCrEndpointAclPoliciesMapFunc,
	}
	preCheck := func() {
		testAccPreCheck(t)
	}
	CrEndpointAclPoliciesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf)
}

// The dependence must create a dedicated instance: picking an arbitrary
// account instance (data.alicloud_cr_ee_instances ids.0) can select one whose
// status does not support endpoint operations, and enabling the endpoint then
// fails with INSTANCE_STATUS_NOT_SUPPORT.
func dataSourceCrEndpointAclPoliciesConfigDependence(name string) string {
	return fmt.Sprintf(`
		variable "name" {
		  default = "%v"
		}
		resource "alicloud_cr_ee_instance" "default" {
		  payment_type   = "Subscription"
		  period         = 1
		  renewal_status = "ManualRenewal"
		  instance_type  = "Economy"
		  instance_name  = var.name
		  image_scanner  = "DISABLE"
		}
		data "alicloud_cr_endpoint_acl_service" "default" {
		  endpoint_type = "internet"
		  enable        = true
		  instance_id   = alicloud_cr_ee_instance.default.id
		  module_name   = "Registry"
		}
		resource "alicloud_cr_endpoint_acl_policy" "default" {
		  instance_id   = alicloud_cr_ee_instance.default.id
		  entry         = "192.168.1.0/24"
		  description   = var.name
		  module_name   = "Registry"
		  endpoint_type = "internet"
		}
		`, name)
}
