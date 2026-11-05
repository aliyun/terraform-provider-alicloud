package cloudmonitorservice_test

import (
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
	name := fmt.Sprintf("tfacccmssubscription%d", rand)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckCmsSubscriptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCmsSubscriptionConfigBasic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "subscription_name", name),
					resource.TestCheckResourceAttr(resourceId, "description", "created by terraform"),
					resource.TestCheckResourceAttr(resourceId, "enable", "false"),
					resource.TestCheckResourceAttr(resourceId, "filter_setting.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "filter_setting.0.relation", "AND"),
					resource.TestCheckResourceAttr(resourceId, "filter_setting.0.conditions.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "filter_setting.0.conditions.0.field", "labels.alertname"),
					resource.TestCheckResourceAttr(resourceId, "filter_setting.0.conditions.0.op", "EQ"),
					resource.TestCheckResourceAttr(resourceId, "filter_setting.0.conditions.0.value", "HighCPUUsage"),
					resource.TestCheckResourceAttrSet(resourceId, "create_time"),
				),
			},
			{
				Config: testAccCmsSubscriptionConfigUpdate(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "description", "updated by terraform"),
					resource.TestCheckResourceAttr(resourceId, "enable", "true"),
					resource.TestCheckResourceAttr(resourceId, "subscribe_legacy_event", "true"),
					resource.TestCheckResourceAttr(resourceId, "filter_setting.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "filter_setting.0.expression", "1 and 2"),
					resource.TestCheckResourceAttr(resourceId, "filter_setting.0.conditions.#", "2"),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckCmsSubscriptionDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_cms_subscription" {
			continue
		}
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)
		_, err := cloudmonitorservice.NewCmsService(client).DescribeCmsSubscription(rs.Primary.ID, rs.Primary.Attributes["workspace"])
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

func testAccCmsSubscriptionConfigBasic(name string) string {
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
`, name)
}

func testAccCmsSubscriptionConfigUpdate(name string) string {
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

resource "alicloud_cms_subscription" "default" {
  subscription_name      = var.name
  workspace              = alicloud_cms_workspace.default.workspace_name
  description            = "updated by terraform"
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
`, name)
}
