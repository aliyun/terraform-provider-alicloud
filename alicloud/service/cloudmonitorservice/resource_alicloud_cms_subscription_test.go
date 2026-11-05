package cloudmonitorservice_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/cloudmonitorservice"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccAliCloudCmsSubscription_basic(t *testing.T) {
	resourceId := "alicloud_cms_subscription.default"
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-cmssubscription-%d", rand)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckCmsSubscriptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCmsSubscriptionConfigMinimal(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "subscription_name", name),
					resource.TestCheckResourceAttr(resourceId, "workspace", name),
					resource.TestCheckResourceAttr(resourceId, "enable", "true"),
					resource.TestCheckResourceAttr(resourceId, "region_id", "cn-hangzhou"),
					resource.TestCheckResourceAttrSet(resourceId, "user_id"),
					resource.TestCheckResourceAttrSet(resourceId, "subscription_type"),
					resource.TestCheckResourceAttrSet(resourceId, "create_time"),
					resource.TestCheckResourceAttrSet(resourceId, "update_time"),
					resource.TestCheckNoResourceAttr(resourceId, "description"),
					resource.TestCheckNoResourceAttr(resourceId, "notify_strategy_id"),
					resource.TestCheckResourceAttr(resourceId, "filter_setting.#", "0"),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateIdFunc: testAccCmsSubscriptionImportStateIdFunc(resourceId),
				ImportStateVerify: true,
			},
			{
				Config: testAccCmsSubscriptionConfigFiltered(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "description", "created by terraform"),
					resource.TestCheckResourceAttr(resourceId, "enable", "false"),
					resource.TestCheckResourceAttr(resourceId, "filter_setting.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "filter_setting.0.relation", "AND"),
					resource.TestCheckResourceAttr(resourceId, "filter_setting.0.conditions.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "filter_setting.0.conditions.0.field", "labels.alertname"),
					resource.TestCheckResourceAttr(resourceId, "filter_setting.0.conditions.0.op", "EQ"),
					resource.TestCheckResourceAttr(resourceId, "filter_setting.0.conditions.0.value", "HighCPUUsage"),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateIdFunc: testAccCmsSubscriptionImportStateIdFunc(resourceId),
				ImportStateVerify: true,
			},
			{
				// Same condition count with different content: a count-only waiter
				// would accept the stale response here.
				Config: testAccCmsSubscriptionConfigFilteredUpdated(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "description", "updated by terraform"),
					resource.TestCheckResourceAttr(resourceId, "filter_setting.0.conditions.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "filter_setting.0.conditions.0.value", "CriticalCPUUsage"),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateIdFunc: testAccCmsSubscriptionImportStateIdFunc(resourceId),
				ImportStateVerify: true,
			},
			{
				// Removing every optional value must clear it remotely, and the
				// explicit false must survive the round trip.
				Config: testAccCmsSubscriptionConfigCleared(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "enable", "true"),
					resource.TestCheckResourceAttr(resourceId, "subscribe_legacy_event", "false"),
					resource.TestCheckNoResourceAttr(resourceId, "description"),
					resource.TestCheckResourceAttr(resourceId, "filter_setting.#", "0"),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateIdFunc: testAccCmsSubscriptionImportStateIdFunc(resourceId),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAliCloudCmsSubscription_complete(t *testing.T) {
	resourceId := "alicloud_cms_subscription.default"
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-cmssubscription-full-%d", rand)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckCmsSubscriptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCmsSubscriptionConfigComplete(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "description", "complete subscription"),
					resource.TestCheckResourceAttr(resourceId, "enable", "true"),
					resource.TestCheckResourceAttr(resourceId, "subscribe_legacy_event", "true"),
					resource.TestCheckResourceAttrPair(resourceId, "notify_strategy_id", "alicloud_cms_event_notify_policy.default", "uuid"),
					resource.TestCheckResourceAttr(resourceId, "filter_setting.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "filter_setting.0.expression", "1 and 2"),
					resource.TestCheckResourceAttr(resourceId, "filter_setting.0.conditions.#", "2"),
					resource.TestCheckResourceAttr(resourceId, "filter_setting.0.conditions.0.field", "labels.alertname"),
					resource.TestCheckResourceAttr(resourceId, "filter_setting.0.conditions.0.op", "EQ"),
					resource.TestCheckResourceAttr(resourceId, "filter_setting.0.conditions.0.value", "HighCPUUsage"),
					resource.TestCheckResourceAttr(resourceId, "filter_setting.0.conditions.1.field", "labels.severity"),
					resource.TestCheckResourceAttr(resourceId, "filter_setting.0.conditions.1.op", "IN"),
					resource.TestCheckResourceAttr(resourceId, "filter_setting.0.conditions.1.value", "P1,P2"),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateIdFunc: testAccCmsSubscriptionImportStateIdFunc(resourceId),
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCmsSubscriptionImportStateIdFunc(resourceId string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceId]
		if !ok {
			return "", fmt.Errorf("resource %s not found in state", resourceId)
		}
		return fmt.Sprintf("%s:%s", rs.Primary.Attributes["workspace"], rs.Primary.ID), nil
	}
}

func testAccCheckCmsSubscriptionDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_cms_subscription" {
			continue
		}
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)
		_, err := cloudmonitorservice.NewCmsService(client).DescribeCmsSubscription(context.Background(), rs.Primary.ID, rs.Primary.Attributes["workspace"])
		if err != nil {
			if errs.NotFoundError(err) {
				continue
			}
			return fmt.Errorf("describing cms subscription %s failed: %s", rs.Primary.ID, err)
		}
		return fmt.Errorf("cms subscription %s still exists", rs.Primary.ID)
	}
	return nil
}

func testAccCmsSubscriptionDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
  project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
  workspace_name = var.name
  sls_project    = alicloud_log_project.default.project_name
}
`, name)
}

func testAccCmsSubscriptionConfigMinimal(name string) string {
	return fmt.Sprintf(`%s
resource "alicloud_cms_subscription" "default" {
  subscription_name = var.name
  workspace         = alicloud_cms_workspace.default.workspace_name
}
`, testAccCmsSubscriptionDependence(name))
}

func testAccCmsSubscriptionConfigFiltered(name string) string {
	return fmt.Sprintf(`%s
resource "alicloud_cms_subscription" "default" {
  subscription_name = var.name
  workspace         = alicloud_cms_workspace.default.workspace_name
  description       = "created by terraform"
  enable            = false

  filter_setting {
    relation = "AND"

    conditions {
      field = "labels.alertname"
      op    = "EQ"
      value = "HighCPUUsage"
    }
  }
}
`, testAccCmsSubscriptionDependence(name))
}

func testAccCmsSubscriptionConfigFilteredUpdated(name string) string {
	return fmt.Sprintf(`%s
resource "alicloud_cms_subscription" "default" {
  subscription_name = var.name
  workspace         = alicloud_cms_workspace.default.workspace_name
  description       = "updated by terraform"
  enable            = false

  filter_setting {
    relation = "AND"

    conditions {
      field = "labels.alertname"
      op    = "EQ"
      value = "CriticalCPUUsage"
    }
  }
}
`, testAccCmsSubscriptionDependence(name))
}

func testAccCmsSubscriptionConfigCleared(name string) string {
	return fmt.Sprintf(`%s
resource "alicloud_cms_subscription" "default" {
  subscription_name      = var.name
  workspace              = alicloud_cms_workspace.default.workspace_name
  enable                 = true
  subscribe_legacy_event = false
}
`, testAccCmsSubscriptionDependence(name))
}

func testAccCmsSubscriptionCompleteDependence(name string) string {
	return fmt.Sprintf(`%s
resource "alicloud_cms_event_notify_policy" "default" {
  workspace   = alicloud_cms_workspace.default.workspace_name
  name        = var.name
  description = "Event notification policy managed by Terraform"

  response_plan {
    repeat_notify_setting {
      end_incident_state = "resolved"
      repeat_interval    = "30"
    }
    auto_recover_seconds = "600"
  }

  notify_strategy {
    description                  = "Notify strategy for subscription test"
    ignore_restored_notification = false

    grouping_setting {
      grouping_keys = ["severity"]
      period_min    = "5"
      times         = "1"
      silence_sec   = "300"
    }

    routes {
      effect_time_range {
        time_zone            = "Asia/Shanghai"
        start_time_in_minute = "0"
        end_time_in_minute   = "1439"
        day_in_week          = ["1", "2", "3", "4", "5"]
      }

      channels {
        receivers            = ["tf-test-group"]
        channel_type         = "DING"
        enabled_sub_channels = []
      }
    }
  }
}
`, testAccCmsSubscriptionDependence(name))
}

func testAccCmsSubscriptionConfigComplete(name string) string {
	return fmt.Sprintf(`%s
resource "alicloud_cms_subscription" "default" {
  subscription_name      = var.name
  workspace              = alicloud_cms_workspace.default.workspace_name
  notify_strategy_id     = alicloud_cms_event_notify_policy.default.uuid
  description            = "complete subscription"
  enable                 = true
  subscribe_legacy_event = true

  filter_setting {
    expression = "1 and 2"

    conditions {
      field = "labels.alertname"
      op    = "EQ"
      value = "HighCPUUsage"
    }

    conditions {
      field = "labels.severity"
      op    = "IN"
      value = "P1,P2"
    }
  }
}
`, testAccCmsSubscriptionCompleteDependence(name))
}
