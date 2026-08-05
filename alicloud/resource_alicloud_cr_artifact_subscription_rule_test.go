package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudCrArtifactSubscriptionRule_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cr_artifact_subscription_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudCrArtifactSubscriptionRuleMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CrServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCrArtifactSubscriptionRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccr%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCrArtifactSubscriptionRuleBasicDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":           "${alicloud_cr_ee_instance.default.id}",
					"source_provider":       "DOCKER_HUB",
					"source_namespace_name": "library",
					"source_repo_name":      "alpine",
					"namespace_name":        "${alicloud_cr_ee_namespace.default.name}",
					"repo_name":             "${alicloud_cr_ee_repo.default.name}",
					"tag_regexp":            ".*",
					"tag_count":             "1",
					"override":              "true",
					"accelerate":            "false",
					"platform":              []string{"linux/amd64", "linux/arm64"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_provider":       "DOCKER_HUB",
						"source_namespace_name": "library",
						"source_repo_name":      "alpine",
						"tag_regexp":            ".*",
						"tag_count":             "1",
						"override":              "true",
						"accelerate":            "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"source_provider":       "QUAY",
					"source_namespace_name": "library",
					"source_repo_name":      "busybox",
					"namespace_name":        "${alicloud_cr_ee_namespace.second.name}",
					"repo_name":             "${alicloud_cr_ee_repo.second.name}",
					"tag_regexp":            "v.*",
					"tag_count":             "2",
					"override":              "false",
					"accelerate":            "true",
					"platform":              []string{"linux/amd64"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_provider":  "QUAY",
						"source_repo_name": "busybox",
						"tag_regexp":       "v.*",
						"tag_count":        "2",
						"override":         "false",
						"accelerate":       "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"source_provider":       "DOCKER_HUB",
					"source_namespace_name": "",
					"source_repo_name":      "alpine",
					"namespace_name":        "${alicloud_cr_ee_namespace.default.name}",
					"repo_name":             "${alicloud_cr_ee_repo.default.name}",
					"tag_regexp":            ".*",
					"tag_count":             "3",
					"override":              "true",
					"accelerate":            "false",
					"platform":              []string{"*/*"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tag_count":  "3",
						"override":   "true",
						"accelerate": "false",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudCrArtifactSubscriptionRuleMap = map[string]string{
	"create_time":                   CHECKSET,
	"modified_time":                 CHECKSET,
	"artifact_subscription_rule_id": CHECKSET,
	"source_domain":                 CHECKSET,
}

func AlicloudCrArtifactSubscriptionRuleBasicDependence(name string) string {
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

resource "alicloud_cr_ee_namespace" "second" {
  instance_id        = alicloud_cr_ee_instance.default.id
  name               = "${var.name}-second"
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

resource "alicloud_cr_ee_repo" "second" {
  instance_id = alicloud_cr_ee_instance.default.id
  namespace   = alicloud_cr_ee_namespace.second.name
  name        = "${var.name}-second"
  repo_type   = "PRIVATE"
  summary     = "test repository second for artifact subscription rule"
}
`, name)
}
