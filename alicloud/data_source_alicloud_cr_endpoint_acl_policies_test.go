package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCREndpointAclPoliciesDataSource(t *testing.T) {

	rand := acctest.RandInt()
	resourceId := "data.alicloud_cr_endpoint_acl_policies.default"
	name := fmt.Sprintf("tf-testacc-CrEndpointAclPolicies%v", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceCrEndpointAclPoliciesConfigDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id":   "${alicloud_cr_endpoint_acl_policy.default.instance_id}",
			"endpoint_type": "internet",
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

func dataSourceCrEndpointAclPoliciesConfigDependence(name string) string {
	return fmt.Sprintf(`
		variable "name" {
		  default = "%v"
		}
		data "alicloud_cr_ee_instances" "default" {}
		data "alicloud_cr_endpoint_acl_service" "default" {
		  endpoint_type = "internet"
		  enable        = true
		  instance_id   = data.alicloud_cr_ee_instances.default.ids.0
		  module_name   = "Registry"
		}
		resource "alicloud_cr_endpoint_acl_policy" "default" {
		  instance_id   = data.alicloud_cr_ee_instances.default.ids.0
		  entry         = "192.168.1.0/24"
		  description   = var.name
		  module_name   = "Registry"
		  endpoint_type = "internet"
		}
		`, name)
}
