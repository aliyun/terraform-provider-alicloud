package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceAlicloudCrArtifactSubscriptionRules_basic(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccr%d", rand)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAlicloudCrArtifactSubscriptionRulesConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.alicloud_cr_artifact_subscription_rules.default", "rules.#"),
				),
			},
		},
	})
}

func testAccDataSourceAlicloudCrArtifactSubscriptionRulesConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_cr_ee_instance" "default" {
  default_oss_bucket = "true"
  instance_name      = var.name
  renewal_status     = "ManualRenewal"
  image_scanner      = "ACR"
  period             = "1"
  payment_type       = "Subscription"
  instance_type      = "Standard"

  timeouts {
    create = "30m"
  }
}

resource "alicloud_cr_ee_namespace" "default" {
  instance_id        = alicloud_cr_ee_instance.default.id
  name               = var.name
  auto_create        = false
  default_visibility = "PRIVATE"
}

resource "alicloud_cr_ee_repo" "default" {
  instance_id = alicloud_cr_ee_instance.default.id
  namespace   = alicloud_cr_ee_namespace.default.name
  name        = var.name
  repo_type   = "PRIVATE"
  summary     = "test repository for artifact subscription rule"
}

resource "alicloud_cr_artifact_subscription_rule" "default" {
  instance_id           = alicloud_cr_ee_instance.default.id
  source_provider       = "DOCKER_HUB"
  source_namespace_name = "library"
  source_repo_name      = "alpine"
  namespace_name        = alicloud_cr_ee_namespace.default.name
  repo_name             = alicloud_cr_ee_repo.default.name
  tag_regexp            = ".*"
  tag_count             = 1
  override              = true
  accelerate            = false
  platform              = ["linux/amd64"]
}

data "alicloud_cr_artifact_subscription_rules" "default" {
  instance_id    = alicloud_cr_ee_instance.default.id
  enable_details  = true
  ids             = [alicloud_cr_artifact_subscription_rule.default.id]
}
`, name)
}
