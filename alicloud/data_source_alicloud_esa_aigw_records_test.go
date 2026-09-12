package alicloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func dataSourceEsaAigwRecordsConfig(name string) string {
	instanceId := os.Getenv("ALICLOUD_ESA_AIGW_INSTANCE_ID")
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

variable "aigw_instance_id" {
  default = "%s"
}

resource "alicloud_esa_rate_plan_instance" "default" {
  type         = "NS"
  auto_renew   = "false"
  period       = "1"
  payment_type = "Subscription"
  coverage     = "overseas"
  auto_pay     = "true"
  plan_name    = "high"
}

resource "alicloud_esa_site" "default" {
  site_name   = "tf-testacc-aigw-ds-%[1]s.com"
  instance_id = alicloud_esa_rate_plan_instance.default.id
  coverage    = "overseas"
  access_type = "NS"
}

resource "alicloud_esa_aigw_record" "default" {
  site_id      = alicloud_esa_site.default.id
  instance_id  = var.aigw_instance_id
  record_name  = "ds-aigw-%[1]s.example.com"
}
`, name, instanceId)
}

func TestAccAliCloudEsaAigwRecordsDataSource_basic(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%d", rand)
	resourceId := "data.alicloud_esa_aigw_records.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceEsaAigwRecordsConfig)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccEsaAigwRecordPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id": "${var.aigw_instance_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "records.#", "1"),
				),
			},
		},
	})
}
