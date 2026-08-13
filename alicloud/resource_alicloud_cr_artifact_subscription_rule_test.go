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

	// baselineConfig returns the step1 (create) values for every attribute.
	baselineConfig := func() map[string]interface{} {
		return map[string]interface{}{
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
		}
	}

	// secondValue is the alternate value each mutable attribute is flipped
	// to, one attribute per step, so the single-attribute Update path
	// (d.HasChange -> request -> RPC -> field-match waiter -> Read) is
	// exercised in isolation for every attribute.
	secondValue := map[string]interface{}{
		"source_provider":       "QUAY",
		"source_repo_name":      "busybox",
		"source_namespace_name": "public",
		"tag_regexp":            "v.*",
		"tag_count":             "2",
		"override":              "false",
		"accelerate":            "true",
		"platform":              []string{"linux/amd64"},
		"namespace_name":        "${alicloud_cr_ee_namespace.second.name}",
		"repo_name":             "${alicloud_cr_ee_repo.second.name}",
	}

	// singleAttrOrder lists the attributes flipped one per step in the
	// accumulating config: stepConfig(k) is the baseline with the first k
	// attributes already flipped, so step k+1 differs from step k in exactly
	// one attribute and d.HasChange fires for that single attribute only.
	singleAttrOrder := []string{
		"source_provider",
		"source_repo_name",
		"source_namespace_name",
		"tag_regexp",
		"tag_count",
		"override",
		"accelerate",
		"platform",
	}
	stepConfig := func(k int) map[string]interface{} {
		cfg := baselineConfig()
		for i := 0; i < k; i++ {
			cfg[singleAttrOrder[i]] = secondValue[singleAttrOrder[i]]
		}
		return cfg
	}
	// pairConfig flips namespace_name and repo_name together as the final
	// update step: the target repo must exist in the target namespace, so the
	// two are advanced to the "second" namespace/repo pair in one step rather
	// than producing an invalid intermediate state.
	pairConfig := func() map[string]interface{} {
		cfg := stepConfig(len(singleAttrOrder))
		cfg["namespace_name"] = secondValue["namespace_name"]
		cfg["repo_name"] = secondValue["repo_name"]
		return cfg
	}

	// secondName is the resolved name of the "second" namespace and repo
	// resources, both named "<name>-second" in the dependence.
	secondName := fmt.Sprintf("%s-second", name)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			// step1: create the rule with the baseline configuration.
			{
				Config: testAccConfig(stepConfig(0)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":           CHECKSET,
						"source_provider":       "DOCKER_HUB",
						"source_namespace_name": "library",
						"source_repo_name":      "alpine",
						"tag_regexp":            ".*",
						"tag_count":             "1",
						"override":              "true",
						"accelerate":            "false",
						"platform.#":            "2",
					}),
				),
			},
			// step2: update source_provider only.
			{
				Config: testAccConfig(stepConfig(1)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_provider": "QUAY",
					}),
				),
			},
			// step3: update source_repo_name only.
			{
				Config: testAccConfig(stepConfig(2)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_repo_name": "busybox",
					}),
				),
			},
			// step4: update source_namespace_name only.
			{
				Config: testAccConfig(stepConfig(3)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_namespace_name": "public",
					}),
				),
			},
			// step5: update tag_regexp only.
			{
				Config: testAccConfig(stepConfig(4)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tag_regexp": "v.*",
					}),
				),
			},
			// step6: update tag_count only.
			{
				Config: testAccConfig(stepConfig(5)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tag_count": "2",
					}),
				),
			},
			// step7: update override only.
			{
				Config: testAccConfig(stepConfig(6)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"override": "false",
					}),
				),
			},
			// step8: update accelerate only.
			{
				Config: testAccConfig(stepConfig(7)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"accelerate": "true",
					}),
				),
			},
			// step9: update platform only (list value).
			{
				Config: testAccConfig(stepConfig(8)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"platform.#": "1",
					}),
					resource.TestCheckResourceAttr(resourceId, "platform.0", "linux/amd64"),
				),
			},
			// step10: update namespace_name and repo_name together to the
			// "second" namespace/repo pair (the target repo must live in the
			// target namespace, so the two are flipped in one step).
			{
				Config: testAccConfig(pairConfig()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"namespace_name": secondName,
						"repo_name":      secondName,
					}),
				),
			},
			// step11: import and verify the final state.
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
