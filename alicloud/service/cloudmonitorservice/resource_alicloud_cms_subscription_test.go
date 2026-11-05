package cloudmonitorservice_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/alibabacloud-go/tea/tea"
	credentials "github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/cloudmonitorservice"
	"github.com/hashicorp/terraform-plugin-framework/path"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
	fwschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccAliCloudCmsSubscription_basic(t *testing.T) {
	resourceId := "alicloud_cms_subscription.default"
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc-cmssubscription-%d", rand)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckCmsSubscriptionDependenciesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCmsSubscriptionConfigMinimal(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "subscription_name", name),
					resource.TestCheckResourceAttr(resourceId, "workspace", strings.ToLower(name)),
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
					resource.TestCheckResourceAttrSet(resourceId, "id"),
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
				// Two distinguishable conditions referenced by position: the swap
				// changes the expression semantics, so a set-based waiter would
				// accept the pre-update order.
				Config: testAccCmsSubscriptionConfigTwoConditions(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "filter_setting.0.expression", "1 and 2"),
					resource.TestCheckResourceAttr(resourceId, "filter_setting.0.conditions.#", "2"),
					resource.TestCheckResourceAttr(resourceId, "filter_setting.0.conditions.0.field", "labels.alertname"),
					resource.TestCheckResourceAttr(resourceId, "filter_setting.0.conditions.0.value", "HighCPUUsage"),
					resource.TestCheckResourceAttr(resourceId, "filter_setting.0.conditions.1.field", "labels.severity"),
					resource.TestCheckResourceAttr(resourceId, "filter_setting.0.conditions.1.value", "P1,P2"),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateIdFunc: testAccCmsSubscriptionImportStateIdFunc(resourceId),
				ImportStateVerify: true,
			},
			{
				Config: testAccCmsSubscriptionConfigTwoConditionsSwapped(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "filter_setting.0.expression", "1 and 2"),
					resource.TestCheckResourceAttr(resourceId, "filter_setting.0.conditions.0.field", "labels.severity"),
					resource.TestCheckResourceAttr(resourceId, "filter_setting.0.conditions.0.value", "P1,P2"),
					resource.TestCheckResourceAttr(resourceId, "filter_setting.0.conditions.1.field", "labels.alertname"),
					resource.TestCheckResourceAttr(resourceId, "filter_setting.0.conditions.1.value", "HighCPUUsage"),
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
				// The cleared state must survive import as well.
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
	name := fmt.Sprintf("tf-testAcc-cmssubscription-full-%d", rand)
	var boundId string

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckCmsSubscriptionDependenciesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCmsSubscriptionConfigCompleteUnbound(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCmsSubscriptionCaptureId(resourceId, &boundId),
					resource.TestCheckResourceAttr(resourceId, "description", "complete subscription"),
					resource.TestCheckResourceAttr(resourceId, "enable", "true"),
					resource.TestCheckResourceAttr(resourceId, "subscribe_legacy_event", "true"),
					resource.TestCheckNoResourceAttr(resourceId, "notify_strategy_id"),
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
			{
				// Binding the notification strategy is an in-place update: the
				// subscription must keep the identity captured while it was unbound.
				Config: testAccCmsSubscriptionConfigComplete(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCmsSubscriptionCheckIdUnchanged(resourceId, &boundId),
					resource.TestCheckResourceAttrPair(resourceId, "notify_strategy_id", "alicloud_cms_event_notify_policy.default", "uuid"),
					resource.TestCheckResourceAttr(resourceId, "filter_setting.0.conditions.#", "2"),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateIdFunc: testAccCmsSubscriptionImportStateIdFunc(resourceId),
				ImportStateVerify: true,
			},
			{
				// A name/description/filter update on the bound subscription: the
				// binding stays, the identity stays and the changed members must be
				// visible remotely.
				Config: testAccCmsSubscriptionConfigCompleteUpdated(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCmsSubscriptionCheckIdUnchanged(resourceId, &boundId),
					resource.TestCheckResourceAttr(resourceId, "subscription_name", name+"-updated"),
					resource.TestCheckResourceAttr(resourceId, "description", "complete subscription updated"),
					resource.TestCheckResourceAttrPair(resourceId, "notify_strategy_id", "alicloud_cms_event_notify_policy.default", "uuid"),
					resource.TestCheckResourceAttr(resourceId, "filter_setting.0.conditions.1.value", "P1"),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateIdFunc: testAccCmsSubscriptionImportStateIdFunc(resourceId),
				ImportStateVerify: true,
			},
			{
				// Omitting subscribe_legacy_event falls back to the provider default
				// and the service normalizes the cleared member to false.
				Config: testAccCmsSubscriptionConfigCompleteLegacyOmitted(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCmsSubscriptionCheckIdUnchanged(resourceId, &boundId),
					resource.TestCheckResourceAttr(resourceId, "subscribe_legacy_event", "false"),
					resource.TestCheckResourceAttr(resourceId, "description", "complete subscription updated"),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateIdFunc: testAccCmsSubscriptionImportStateIdFunc(resourceId),
				ImportStateVerify: true,
			},
			{
				// Moving the subscription to another workspace replaces it: the
				// replacement must return a different identity.
				Config: testAccCmsSubscriptionConfigCompleteOtherWorkspace(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCmsSubscriptionCheckIdChanged(resourceId, &boundId),
					resource.TestCheckResourceAttr(resourceId, "workspace", strings.ToLower(name)+"-second"),
					resource.TestCheckNoResourceAttr(resourceId, "notify_strategy_id"),
					resource.TestCheckResourceAttr(resourceId, "description", "complete subscription updated"),
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

func testAccCmsSubscriptionCaptureId(resourceId string, captured *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceId]
		if !ok {
			return fmt.Errorf("resource %s not found in state", resourceId)
		}
		*captured = rs.Primary.ID
		return nil
	}
}

func testAccCmsSubscriptionCheckIdUnchanged(resourceId string, captured *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceId]
		if !ok {
			return fmt.Errorf("resource %s not found in state", resourceId)
		}
		if rs.Primary.ID != *captured {
			return fmt.Errorf("subscription id changed from %q to %q; this update must stay in place", *captured, rs.Primary.ID)
		}
		return nil
	}
}

func testAccCmsSubscriptionCheckIdChanged(resourceId string, captured *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceId]
		if !ok {
			return fmt.Errorf("resource %s not found in state", resourceId)
		}
		if rs.Primary.ID == *captured {
			return fmt.Errorf("subscription id %q did not change; a workspace change must replace the subscription", *captured)
		}
		return nil
	}
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

// testAccCheckCmsSubscriptionDependenciesDestroy verifies the subscription and
// every object this test creates: the workspaces, the event notification policy
// and the SLS project they depend on. Each dependency is looked up directly by
// its own id and scope so a not-found answer for another object, or a text that
// merely mentions not-found, cannot be mistaken for a successful deletion.
func testAccCheckCmsSubscriptionDependenciesDestroy(s *terraform.State) error {
	client := acctest.Provider.Meta().(*connectivity.AliyunClient)
	for _, rs := range s.RootModule().Resources {
		var check func() (bool, error)
		switch rs.Type {
		case "alicloud_cms_subscription":
			id, workspace := rs.Primary.ID, rs.Primary.Attributes["workspace"]
			check = func() (bool, error) {
				_, err := cloudmonitorservice.NewCmsService(client).DescribeSubscription(context.Background(), id, workspace)
				return subscriptionDestroyed(err)
			}
		case "alicloud_cms_workspace":
			action := fmt.Sprintf("/workspace/%s", rs.Primary.ID)
			check = func() (bool, error) {
				_, err := client.RoaGet("Cms", "2024-03-30", action, nil, nil, nil)
				return cmsDependencyDestroyed(err, "WorkspaceNotExist")
			}
		case "alicloud_cms_event_notify_policy":
			uuid, workspace := rs.Primary.Attributes["uuid"], rs.Primary.Attributes["workspace"]
			query := map[string]*string{
				"uuid":      helper.StringPointer(uuid),
				"workspace": helper.StringPointer(workspace),
			}
			check = func() (bool, error) {
				_, err := client.RoaGet("Cms", "2024-03-30", "/api/eventbase/notify-policy", query, nil, nil)
				return cmsDependencyDestroyed(err, "NotFound")
			}
		case "alicloud_log_project":
			hostMap := map[string]*string{"project": helper.StringPointer(rs.Primary.ID)}
			check = func() (bool, error) {
				_, err := client.Do("Sls", helper.RoaParam("GET", "2020-12-30", "GetProject", "/"), nil, nil, nil, hostMap, false)
				return cmsDependencyDestroyed(err, "ProjectNotExist")
			}
		default:
			continue
		}
		if err := waitForCmsDependencyDestroy(check); err != nil {
			return fmt.Errorf("%s %s is not confirmed deleted: %s", rs.Type, rs.Primary.ID, err)
		}
	}
	return nil
}

// subscriptionDestroyed accepts only the contract-level ResourceNotFound the
// service documents for GetSubscription; any other failure, including a
// permission or workspace error whose text mentions not-found, propagates.
func subscriptionDestroyed(err error) (bool, error) {
	if err == nil {
		return false, nil
	}
	if cloudmonitorservice.IsSubscriptionNotFound(err) {
		return true, nil
	}
	return false, err
}

// cmsDependencyDestroyed accepts only the documented absence code returned by
// the dependency lookup itself. An authorization or server error is never a
// deletion, even when its message or data happens to contain "NotFound", and
// the code must match the contract for the object that was queried.
func cmsDependencyDestroyed(err error, codes ...string) (bool, error) {
	if err == nil {
		return false, nil
	}
	sdkErr, ok := err.(*tea.SDKError)
	if !ok {
		return false, err
	}
	status := tea.IntValue(sdkErr.StatusCode)
	code := tea.StringValue(sdkErr.Code)
	if status == http.StatusUnauthorized || status == http.StatusForbidden || status >= http.StatusInternalServerError {
		return false, fmt.Errorf("the dependency lookup was rejected with HTTP %d and code %q", status, code)
	}
	for _, expected := range codes {
		if code == expected {
			return true, nil
		}
	}
	return false, fmt.Errorf("the dependency lookup returned code %q, which is not the documented absence code", code)
}

// waitForCmsDependencyDestroy polls because Cms and SLS delete asynchronously:
// a workspace that is still being deleted reports its own transient error
// instead of the documented not-found signal.
func waitForCmsDependencyDestroy(check func() (bool, error)) error {
	deadline := time.Now().Add(5 * time.Minute)
	var lastErr error
	for {
		gone, err := check()
		if gone {
			return nil
		}
		lastErr = err
		if time.Now().After(deadline) {
			if lastErr != nil {
				return lastErr
			}
			return fmt.Errorf("the object still exists after the destroy wait")
		}
		time.Sleep(5 * time.Second)
	}
}

func testAccCmsSubscriptionDependence(name string) string {
	// SLS project names only accept lowercase letters, digits and hyphens, so the
	// dependency name cannot carry the mixed-case tf-testAcc prefix verbatim.
	baseName := strings.ToLower(name)
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

variable "base_name" {
  default = "%s"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
  project_name = var.base_name
}

resource "alicloud_cms_workspace" "default" {
  workspace_name = var.base_name
  sls_project    = alicloud_log_project.default.project_name
}
`, name, baseName)
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

func testAccCmsSubscriptionConfigTwoConditions(name string) string {
	return fmt.Sprintf(`%s
resource "alicloud_cms_subscription" "default" {
  subscription_name = var.name
  workspace         = alicloud_cms_workspace.default.workspace_name
  description       = "ordered conditions"

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
`, testAccCmsSubscriptionDependence(name))
}

func testAccCmsSubscriptionConfigTwoConditionsSwapped(name string) string {
	return fmt.Sprintf(`%s
resource "alicloud_cms_subscription" "default" {
  subscription_name = var.name
  workspace         = alicloud_cms_workspace.default.workspace_name
  description       = "ordered conditions"

  filter_setting {
    expression = "1 and 2"

    conditions {
      field = "labels.severity"
      op    = "IN"
      value = "P1,P2"
    }

    conditions {
      field = "labels.alertname"
      op    = "EQ"
      value = "HighCPUUsage"
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

func testAccCmsSubscriptionSecondWorkspaceDependence(name string) string {
	baseName := strings.ToLower(name)
	return fmt.Sprintf(`%s
resource "alicloud_log_project" "second" {
  project_name = "%s-second"
}

resource "alicloud_cms_workspace" "second" {
  workspace_name = alicloud_log_project.second.project_name
  sls_project    = alicloud_log_project.second.project_name
}
`, testAccCmsSubscriptionCompleteDependence(name), baseName)
}

func testAccCmsSubscriptionConfigCompleteUnbound(name string) string {
	return fmt.Sprintf(`%s
resource "alicloud_cms_subscription" "default" {
  subscription_name      = var.name
  workspace              = alicloud_cms_workspace.default.workspace_name
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

func testAccCmsSubscriptionConfigCompleteOtherWorkspace(name string) string {
	return fmt.Sprintf(`%s
resource "alicloud_cms_subscription" "default" {
  subscription_name = "${var.name}-updated"
  workspace         = alicloud_cms_workspace.second.workspace_name
  description       = "complete subscription updated"
  enable            = true

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
      value = "P1"
    }
  }
}
`, testAccCmsSubscriptionSecondWorkspaceDependence(name))
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

func testAccCmsSubscriptionConfigCompleteUpdated(name string) string {
	return fmt.Sprintf(`%s
resource "alicloud_cms_subscription" "default" {
  subscription_name      = "${var.name}-updated"
  workspace              = alicloud_cms_workspace.default.workspace_name
  notify_strategy_id     = alicloud_cms_event_notify_policy.default.uuid
  description            = "complete subscription updated"
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
      value = "P1"
    }
  }
}
`, testAccCmsSubscriptionCompleteDependence(name))
}

func testAccCmsSubscriptionConfigCompleteLegacyOmitted(name string) string {
	return fmt.Sprintf(`%s
resource "alicloud_cms_subscription" "default" {
  subscription_name  = "${var.name}-updated"
  workspace          = alicloud_cms_workspace.default.workspace_name
  notify_strategy_id = alicloud_cms_event_notify_policy.default.uuid
  description        = "complete subscription updated"
  enable             = true

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
      value = "P1"
    }
  }
}
`, testAccCmsSubscriptionCompleteDependence(name))
}

// The tests below drive the resource through its public factory, schema and
// lifecycle entry points against a local Cms endpoint, so the assertions cover
// the request the resource would send and the state it produces.

const subscriptionUnitServerId = "sub-unit-123"

type subscriptionUnitRecorder struct {
	mu       sync.Mutex
	methods  []string
	paths    []string
	bodies   []map[string]interface{}
	response func(method, path string, body map[string]interface{}) (int, interface{})
}

func (r *subscriptionUnitRecorder) handle(w http.ResponseWriter, req *http.Request) {
	var body map[string]interface{}
	if req.Body != nil {
		_ = json.NewDecoder(req.Body).Decode(&body)
	}
	r.mu.Lock()
	r.methods = append(r.methods, req.Method)
	r.paths = append(r.paths, req.URL.Path)
	r.bodies = append(r.bodies, body)
	handler := r.response
	r.mu.Unlock()

	status, payload := handler(req.Method, req.URL.Path, body)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if payload != nil {
		_ = json.NewEncoder(w).Encode(payload)
	}
}

func (r *subscriptionUnitRecorder) requests(method string) []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	paths := make([]string, 0)
	for i, recorded := range r.methods {
		if recorded == method {
			paths = append(paths, r.paths[i])
		}
	}
	return paths
}

func (r *subscriptionUnitRecorder) bodiesFor(method string) []map[string]interface{} {
	r.mu.Lock()
	defer r.mu.Unlock()
	bodies := make([]map[string]interface{}, 0)
	for i, recorded := range r.methods {
		if recorded == method {
			bodies = append(bodies, r.bodies[i])
		}
	}
	return bodies
}

func subscriptionUnitClient(t *testing.T, response func(method, path string, body map[string]interface{}) (int, interface{})) (*connectivity.AliyunClient, *subscriptionUnitRecorder) {
	t.Helper()
	recorder := &subscriptionUnitRecorder{response: response}
	server := httptest.NewServer(http.HandlerFunc(recorder.handle))
	t.Cleanup(server.Close)

	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").SetAccessKeyId("test-key").SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatal(err)
	}
	endpoints := new(sync.Map)
	endpoint := strings.TrimPrefix(server.URL, "http://")
	t.Setenv("NO_PROXY", endpoint)
	config := &connectivity.Config{
		AccessKey: "test-key", SecretKey: "test-secret", Credential: credential,
		RegionId: "cn-hangzhou", AccountType: "test", Protocol: "http",
		Endpoints: endpoints, SignVersion: new(sync.Map), SkipRegionValidation: true,
	}
	client, err := config.Client()
	if err != nil {
		t.Fatal(err)
	}
	endpoints.Store("cms", endpoint)
	return client, recorder
}

func subscriptionUnitResource(t *testing.T, apiClient *connectivity.AliyunClient) fwresource.Resource {
	t.Helper()
	r := cloudmonitorservice.NewSubscriptionResource()
	configurable, ok := r.(fwresource.ResourceWithConfigure)
	if !ok {
		t.Fatalf("%T does not implement resource.ResourceWithConfigure", r)
	}
	var resp fwresource.ConfigureResponse
	configurable.Configure(context.Background(), fwresource.ConfigureRequest{ProviderData: apiClient}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("configuring the resource: %v", resp.Diagnostics)
	}
	return r
}

func subscriptionUnitSchema(t *testing.T, r fwresource.Resource) (fwschema.Schema, tftypes.Object) {
	t.Helper()
	ctx := context.Background()
	var schemaResp fwresource.SchemaResponse
	r.Schema(ctx, fwresource.SchemaRequest{}, &schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("building schema: %v", schemaResp.Diagnostics)
	}
	objType, ok := schemaResp.Schema.Type().TerraformType(ctx).(tftypes.Object)
	if !ok {
		t.Fatalf("schema type is %T, want tftypes.Object", schemaResp.Schema.Type().TerraformType(ctx))
	}
	return schemaResp.Schema, objType
}

type subscriptionUnitUnknown struct{}

func subscriptionUnitValue(t *testing.T, valueType tftypes.Type, raw interface{}) tftypes.Value {
	t.Helper()
	if prepared, ok := raw.(tftypes.Value); ok {
		return prepared
	}
	if _, ok := raw.(subscriptionUnitUnknown); ok {
		return tftypes.NewValue(valueType, tftypes.UnknownValue)
	}
	if raw == nil {
		return tftypes.NewValue(valueType, nil)
	}
	switch {
	case valueType.Is(tftypes.String):
		return tftypes.NewValue(valueType, raw.(string))
	case valueType.Is(tftypes.Bool):
		return tftypes.NewValue(valueType, raw.(bool))
	case valueType.Is(tftypes.Number):
		return tftypes.NewValue(valueType, raw.(float64))
	}
	switch typed := valueType.(type) {
	case tftypes.List:
		items := raw.([]interface{})
		values := make([]tftypes.Value, 0, len(items))
		for _, item := range items {
			values = append(values, subscriptionUnitValue(t, typed.ElementType, item))
		}
		return tftypes.NewValue(valueType, values)
	case tftypes.Object:
		attributes := raw.(map[string]interface{})
		values := make(map[string]tftypes.Value, len(typed.AttributeTypes))
		for name, attributeType := range typed.AttributeTypes {
			attributeRaw, present := attributes[name]
			if !present {
				values[name] = tftypes.NewValue(attributeType, nil)
				continue
			}
			values[name] = subscriptionUnitValue(t, attributeType, attributeRaw)
		}
		return tftypes.NewValue(valueType, values)
	}
	t.Fatalf("unsupported tftypes type %T", valueType)
	return tftypes.Value{}
}

func subscriptionUnitObject(t *testing.T, objType tftypes.Object, attributes map[string]interface{}) tftypes.Value {
	t.Helper()
	values := make(map[string]tftypes.Value, len(objType.AttributeTypes))
	for name, attributeType := range objType.AttributeTypes {
		raw, present := attributes[name]
		if !present {
			values[name] = tftypes.NewValue(attributeType, nil)
			continue
		}
		values[name] = subscriptionUnitValue(t, attributeType, raw)
	}
	return tftypes.NewValue(objType, values)
}

func subscriptionUnitObjectType(t *testing.T, objType tftypes.Object, name string) tftypes.Object {
	t.Helper()
	listType, ok := objType.AttributeTypes[name].(tftypes.List)
	if !ok {
		t.Fatalf("%s is %T, want a list", name, objType.AttributeTypes[name])
	}
	elementType, ok := listType.ElementType.(tftypes.Object)
	if !ok {
		t.Fatalf("%s element is %T, want an object", name, listType.ElementType)
	}
	return elementType
}

func subscriptionUnitStringAttribute(t *testing.T, state tfsdk.State, name string) types.String {
	t.Helper()
	var value types.String
	if diags := state.GetAttribute(context.Background(), path.Root(name), &value); diags.HasError() {
		t.Fatalf("reading state attribute %q: %v", name, diags)
	}
	return value
}

// TestUnitCmsSubscriptionCreateRequestShape proves that a legal zero value is
// sent explicitly and that unset values stay absent on the wire. An empty
// string is not a legal value: it is rejected by validation and reads back as
// null, so it is never part of a request.
func TestUnitCmsSubscriptionCreateRequestShape(t *testing.T) {
	created := map[string]interface{}{}
	apiClient, recorder := subscriptionUnitClient(t, func(method, path string, body map[string]interface{}) (int, interface{}) {
		if method == http.MethodPost {
			for key, value := range body {
				created[key] = value
			}
			return http.StatusOK, map[string]interface{}{"requestId": "req-1", "subscriptionId": subscriptionUnitServerId}
		}
		subscription := map[string]interface{}{
			"subscriptionId":   subscriptionUnitServerId,
			"workspace":        "ws-1",
			"enable":           false,
			"userId":           "1",
			"subscriptionType": "NORMAL",
			"createTime":       "1",
			"updateTime":       "1",
		}
		for key, value := range created {
			subscription[key] = value
		}
		return http.StatusOK, map[string]interface{}{"requestId": "req-2", "subscription": subscription}
	})
	r := subscriptionUnitResource(t, apiClient)
	resourceSchema, objType := subscriptionUnitSchema(t, r)

	planValue := subscriptionUnitObject(t, objType, map[string]interface{}{
		"id":                     subscriptionUnitUnknown{},
		"subscription_name":      "zero-values",
		"workspace":              "ws-1",
		"enable":                 false,
		"subscribe_legacy_event": false,
		"filter_setting": []interface{}{
			map[string]interface{}{
				"conditions": []interface{}{
					map[string]interface{}{"field": "labels.alertname", "op": "EQ", "value": "CPU"},
				},
			},
		},
	})
	req := fwresource.CreateRequest{Plan: tfsdk.Plan{Schema: resourceSchema, Raw: planValue}}
	resp := &fwresource.CreateResponse{State: tfsdk.State{Schema: resourceSchema, Raw: planValue}}

	r.Create(context.Background(), req, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Create returned diagnostics: %v", resp.Diagnostics)
	}

	bodies := recorder.bodiesFor(http.MethodPost)
	if len(bodies) != 1 {
		t.Fatalf("expected exactly one CreateSubscription call, got %v", recorder.requests(http.MethodPost))
	}
	body := bodies[0]
	if got, ok := body["description"]; ok {
		t.Fatalf("an unset description must stay absent from the request, got %v", got)
	}
	if got, ok := body["subscribeLegacyEvent"]; !ok || got != false {
		t.Fatalf("subscribeLegacyEvent must be sent explicitly as false, got %v (present=%t)", got, ok)
	}
	filter, ok := body["filterSetting"].(map[string]interface{})
	if !ok {
		t.Fatalf("filterSetting missing from the request: %v", body)
	}
	if got, ok := filter["expression"]; ok {
		t.Fatalf("an unset expression must stay absent from the request, got %v", got)
	}
	conditions, ok := filter["conditions"].([]interface{})
	if !ok || len(conditions) != 1 {
		t.Fatalf("conditions must carry the configured member, got %v", filter["conditions"])
	}
	condition := conditions[0].(map[string]interface{})
	if condition["field"] != "labels.alertname" || condition["op"] != "EQ" || condition["value"] != "CPU" {
		t.Fatalf("unexpected condition payload: %v", condition)
	}

	if got := subscriptionUnitStringAttribute(t, resp.State, "id").ValueString(); got != subscriptionUnitServerId {
		t.Fatalf("id = %q, want %q", got, subscriptionUnitServerId)
	}
	if got := subscriptionUnitStringAttribute(t, resp.State, "region_id").ValueString(); got != "cn-hangzhou" {
		t.Fatalf("region_id = %q, want the provider region", got)
	}
}

// TestUnitCmsSubscriptionCreateKeepsIdentityOnWaitFailure proves that a failure
// after CreateSubscription returned an id still leaves a destroyable state.
func TestUnitCmsSubscriptionCreateKeepsIdentityOnWaitFailure(t *testing.T) {
	apiClient, recorder := subscriptionUnitClient(t, func(method, path string, body map[string]interface{}) (int, interface{}) {
		if method == http.MethodPost {
			return http.StatusOK, map[string]interface{}{"requestId": "req-1", "subscriptionId": subscriptionUnitServerId}
		}
		return http.StatusForbidden, map[string]interface{}{"requestId": "req-2", "Code": "Forbidden", "Message": "denied"}
	})
	r := subscriptionUnitResource(t, apiClient)
	resourceSchema, objType := subscriptionUnitSchema(t, r)
	planValue := subscriptionUnitObject(t, objType, map[string]interface{}{
		"id":                subscriptionUnitUnknown{},
		"subscription_name": "identity",
		"workspace":         "ws-1",
		"enable":            true,
	})
	req := fwresource.CreateRequest{Plan: tfsdk.Plan{Schema: resourceSchema, Raw: planValue}}
	resp := &fwresource.CreateResponse{State: tfsdk.State{Schema: resourceSchema, Raw: planValue}}

	r.Create(context.Background(), req, resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected Create to fail when the first read never succeeds")
	}
	if got := subscriptionUnitStringAttribute(t, resp.State, "id").ValueString(); got != subscriptionUnitServerId {
		t.Fatalf("state id = %q, want the created id to be persisted before the failing read", got)
	}
	if got := subscriptionUnitStringAttribute(t, resp.State, "workspace").ValueString(); got != "ws-1" {
		t.Fatalf("state workspace = %q, want ws-1", got)
	}
	if len(recorder.requests(http.MethodGet)) == 0 {
		t.Fatal("no read was recorded")
	}
}

// TestUnitCmsSubscriptionCreateRejectsIdentityConflict proves a response that
// reports another subscription id is an error instead of overwriting the
// identity returned by CreateSubscription.
func TestUnitCmsSubscriptionCreateRejectsIdentityConflict(t *testing.T) {
	apiClient, _ := subscriptionUnitClient(t, func(method, path string, body map[string]interface{}) (int, interface{}) {
		if method == http.MethodPost {
			return http.StatusOK, map[string]interface{}{"requestId": "req-1", "subscriptionId": subscriptionUnitServerId}
		}
		return http.StatusOK, map[string]interface{}{"requestId": "req-2", "subscription": map[string]interface{}{
			"subscriptionId":   "another-id",
			"subscriptionName": "identity",
			"workspace":        "ws-1",
		}}
	})
	r := subscriptionUnitResource(t, apiClient)
	resourceSchema, objType := subscriptionUnitSchema(t, r)
	planValue := subscriptionUnitObject(t, objType, map[string]interface{}{
		"id":                subscriptionUnitUnknown{},
		"subscription_name": "identity",
		"workspace":         "ws-1",
	})
	req := fwresource.CreateRequest{Plan: tfsdk.Plan{Schema: resourceSchema, Raw: planValue}}
	resp := &fwresource.CreateResponse{State: tfsdk.State{Schema: resourceSchema, Raw: planValue}}

	r.Create(context.Background(), req, resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected Create to fail when the read reports a different subscription id")
	}
	if got := subscriptionUnitStringAttribute(t, resp.State, "id").ValueString(); got != subscriptionUnitServerId {
		t.Fatalf("state id = %q, want the confirmed id %q", got, subscriptionUnitServerId)
	}
}

// TestUnitCmsSubscriptionReadFlattensLists proves list attributes rebuild with
// the proper scalar element type instead of nested lists.
func TestUnitCmsSubscriptionReadFlattensLists(t *testing.T) {
	apiClient, _ := subscriptionUnitClient(t, func(method, path string, body map[string]interface{}) (int, interface{}) {
		return http.StatusOK, map[string]interface{}{"requestId": "req", "subscription": map[string]interface{}{
			"subscriptionId":   "sub-1",
			"subscriptionName": "name",
			"workspace":        "ws-1",
			"pushingSetting": map[string]interface{}{
				"alertActionIds":   []interface{}{"a1"},
				"restoreActionIds": []interface{}{"a2"},
			},
			"workspaceFilterSetting": map[string]interface{}{
				"workspaceUuids": []interface{}{"w1"},
			},
			"agentConfig": map[string]interface{}{
				"agentUuid": "agent-1",
				"routes": []interface{}{
					map[string]interface{}{
						"channels": []interface{}{
							map[string]interface{}{
								"channelType":        "WEBHOOK",
								"receivers":          []interface{}{"r1"},
								"enabledSubChannels": []interface{}{"EMAIL"},
							},
						},
						"effectTimeRange": map[string]interface{}{
							"dayInWeek": []interface{}{1, 2},
						},
					},
				},
			},
		}}
	})
	r := subscriptionUnitResource(t, apiClient)
	resourceSchema, objType := subscriptionUnitSchema(t, r)
	stateValue := subscriptionUnitObject(t, objType, map[string]interface{}{
		"id":        "sub-1",
		"workspace": "ws-1",
	})
	req := fwresource.ReadRequest{State: tfsdk.State{Schema: resourceSchema, Raw: stateValue}}
	resp := &fwresource.ReadResponse{State: tfsdk.State{Schema: resourceSchema, Raw: stateValue}}

	r.Read(context.Background(), req, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Read returned diagnostics for valid list values: %v", resp.Diagnostics)
	}

	assertStringList := func(name string, attributePath path.Path, want []string) {
		t.Helper()
		var list types.List
		if diags := resp.State.GetAttribute(context.Background(), attributePath, &list); diags.HasError() {
			t.Fatalf("reading %s: %v", name, diags)
		}
		if list.IsNull() {
			t.Fatalf("%s is null", name)
		}
		elements := list.Elements()
		if len(elements) != len(want) {
			t.Fatalf("%s has %d elements, want %d", name, len(elements), len(want))
		}
		for i, element := range elements {
			value, ok := element.(types.String)
			if !ok {
				t.Fatalf("%s[%d] is %T, want types.String (nested list element type)", name, i, element)
			}
			if value.ValueString() != want[i] {
				t.Fatalf("%s[%d] = %q, want %q", name, i, value.ValueString(), want[i])
			}
		}
	}

	assertStringList("pushing_setting.alert_action_ids",
		path.Root("pushing_setting").AtListIndex(0).AtName("alert_action_ids"), []string{"a1"})
	assertStringList("pushing_setting.restore_action_ids",
		path.Root("pushing_setting").AtListIndex(0).AtName("restore_action_ids"), []string{"a2"})
	assertStringList("workspace_filter_setting.workspace_uuids",
		path.Root("workspace_filter_setting").AtListIndex(0).AtName("workspace_uuids"), []string{"w1"})
	channelPath := path.Root("agent_config").AtListIndex(0).AtName("routes").AtListIndex(0).
		AtName("channels").AtListIndex(0)
	assertStringList("agent_config.routes.channels.receivers", channelPath.AtName("receivers"), []string{"r1"})
	assertStringList("agent_config.routes.channels.enabled_sub_channels", channelPath.AtName("enabled_sub_channels"), []string{"EMAIL"})
}

// TestUnitCmsSubscriptionReadHandlesMissingLists proves an omitted list reads
// back as null without a state diagnostic; this is the empty-list path that
// previously built a list of the wrong element type.
func TestUnitCmsSubscriptionReadHandlesMissingLists(t *testing.T) {
	apiClient, _ := subscriptionUnitClient(t, func(method, path string, body map[string]interface{}) (int, interface{}) {
		return http.StatusOK, map[string]interface{}{"requestId": "req", "subscription": map[string]interface{}{
			"subscriptionId":   "sub-1",
			"subscriptionName": "name",
			"workspace":        "ws-1",
			"agentConfig": map[string]interface{}{
				"agentUuid": "agent-1",
				"routes": []interface{}{
					map[string]interface{}{
						"channels": []interface{}{
							map[string]interface{}{"channelType": "WEBHOOK"},
						},
					},
				},
			},
		}}
	})
	r := subscriptionUnitResource(t, apiClient)
	resourceSchema, objType := subscriptionUnitSchema(t, r)
	stateValue := subscriptionUnitObject(t, objType, map[string]interface{}{
		"id":        "sub-1",
		"workspace": "ws-1",
	})
	req := fwresource.ReadRequest{State: tfsdk.State{Schema: resourceSchema, Raw: stateValue}}
	resp := &fwresource.ReadResponse{State: tfsdk.State{Schema: resourceSchema, Raw: stateValue}}

	r.Read(context.Background(), req, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Read returned diagnostics for missing lists: %v", resp.Diagnostics)
	}

	var receivers types.List
	channelPath := path.Root("agent_config").AtListIndex(0).AtName("routes").AtListIndex(0).
		AtName("channels").AtListIndex(0)
	if diags := resp.State.GetAttribute(context.Background(), channelPath.AtName("receivers"), &receivers); diags.HasError() {
		t.Fatalf("reading receivers: %v", diags)
	}
	if !receivers.IsNull() {
		t.Fatalf("receivers = %v, want null", receivers)
	}
}

// TestUnitCmsSubscriptionReadRejectsNullSubscription proves an off-contract
// null payload keeps state and reports an error instead of deleting the
// resource.
func TestUnitCmsSubscriptionReadRejectsNullSubscription(t *testing.T) {
	apiClient, _ := subscriptionUnitClient(t, func(method, path string, body map[string]interface{}) (int, interface{}) {
		return http.StatusOK, map[string]interface{}{"requestId": "req", "subscription": nil}
	})
	r := subscriptionUnitResource(t, apiClient)
	resourceSchema, objType := subscriptionUnitSchema(t, r)
	stateValue := subscriptionUnitObject(t, objType, map[string]interface{}{
		"id":        "sub-1",
		"workspace": "ws-1",
	})
	req := fwresource.ReadRequest{State: tfsdk.State{Schema: resourceSchema, Raw: stateValue}}
	resp := &fwresource.ReadResponse{State: tfsdk.State{Schema: resourceSchema, Raw: stateValue}}

	r.Read(context.Background(), req, resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected Read to fail for a null subscription payload")
	}
	if got := subscriptionUnitStringAttribute(t, resp.State, "id").ValueString(); got != "sub-1" {
		t.Fatalf("state id = %q, want the resource to keep its state", got)
	}
}

// TestUnitCmsSubscriptionReadDropsRemovedBlocks proves Refresh no longer hides
// a remote block removal behind the previous state, so Refresh and a bare
// import see the same configuration.
func TestUnitCmsSubscriptionReadDropsRemovedBlocks(t *testing.T) {
	apiClient, _ := subscriptionUnitClient(t, func(method, path string, body map[string]interface{}) (int, interface{}) {
		return http.StatusOK, map[string]interface{}{"requestId": "req", "subscription": map[string]interface{}{
			"subscriptionId":   "sub-1",
			"subscriptionName": "name",
			"workspace":        "ws-1",
		}}
	})
	r := subscriptionUnitResource(t, apiClient)
	resourceSchema, objType := subscriptionUnitSchema(t, r)
	filterType := subscriptionUnitObjectType(t, objType, "filter_setting")
	conditionType := subscriptionUnitObjectType(t, filterType, "conditions")
	stateValue := subscriptionUnitObject(t, objType, map[string]interface{}{
		"id":        "sub-1",
		"workspace": "ws-1",
		"filter_setting": []interface{}{
			map[string]interface{}{
				"relation": "AND",
				"conditions": []interface{}{
					subscriptionUnitObject(t, conditionType, map[string]interface{}{
						"field": "labels.alertname", "op": "EQ", "value": "CPU",
					}),
				},
			},
		},
	})
	req := fwresource.ReadRequest{State: tfsdk.State{Schema: resourceSchema, Raw: stateValue}}
	resp := &fwresource.ReadResponse{State: tfsdk.State{Schema: resourceSchema, Raw: stateValue}}

	r.Read(context.Background(), req, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Read returned diagnostics: %v", resp.Diagnostics)
	}

	var filterSetting types.List
	if diags := resp.State.GetAttribute(context.Background(), path.Root("filter_setting"), &filterSetting); diags.HasError() {
		t.Fatalf("reading filter_setting: %v", diags)
	}
	if !filterSetting.IsNull() && len(filterSetting.Elements()) != 0 {
		t.Fatalf("filter_setting = %v, want an empty state after the service dropped the block", filterSetting)
	}
}

// TestUnitCmsSubscriptionUpdateEnableOnlySkipsBody proves a toggle-only update
// does not rewrite the subscription body.
func TestUnitCmsSubscriptionUpdateEnableOnlySkipsBody(t *testing.T) {
	disabled := false
	apiClient, recorder := subscriptionUnitClient(t, func(method, path string, body map[string]interface{}) (int, interface{}) {
		if method == http.MethodPut {
			disabled = true
			return http.StatusOK, map[string]interface{}{"requestId": "req-1"}
		}
		return http.StatusOK, map[string]interface{}{"requestId": "req-2", "subscription": map[string]interface{}{
			"subscriptionId":   "sub-1",
			"subscriptionName": "name",
			"workspace":        "ws-1",
			"enable":           !disabled,
		}}
	})
	r := subscriptionUnitResource(t, apiClient)
	resourceSchema, objType := subscriptionUnitSchema(t, r)
	stateValue := subscriptionUnitObject(t, objType, map[string]interface{}{
		"id":                "sub-1",
		"subscription_name": "name",
		"workspace":         "ws-1",
		"enable":            true,
	})
	planValue := subscriptionUnitObject(t, objType, map[string]interface{}{
		"id":                "sub-1",
		"subscription_name": "name",
		"workspace":         "ws-1",
		"enable":            false,
	})

	req := fwresource.UpdateRequest{
		Plan:  tfsdk.Plan{Schema: resourceSchema, Raw: planValue},
		State: tfsdk.State{Schema: resourceSchema, Raw: stateValue},
	}
	resp := &fwresource.UpdateResponse{State: tfsdk.State{Schema: resourceSchema, Raw: stateValue}}
	r.Update(context.Background(), req, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Update returned diagnostics: %v", resp.Diagnostics)
	}

	puts := recorder.requests(http.MethodPut)
	if len(puts) != 1 || !strings.HasSuffix(puts[0], "/disable") {
		t.Fatalf("put requests = %v, want only the disable action", puts)
	}
}

// TestUnitCmsSubscriptionUpdateRequiresOrderedConditions proves the update
// waiter compares conditions by position: a stale response with the previous
// order must not converge.
func TestUnitCmsSubscriptionUpdateRequiresOrderedConditions(t *testing.T) {
	apiClient, _ := subscriptionUnitClient(t, func(method, path string, body map[string]interface{}) (int, interface{}) {
		if method == http.MethodPut {
			return http.StatusOK, map[string]interface{}{"requestId": "req-1"}
		}
		return http.StatusOK, map[string]interface{}{"requestId": "req-2", "subscription": map[string]interface{}{
			"subscriptionId":   "sub-1",
			"subscriptionName": "name",
			"workspace":        "ws-1",
			"filterSetting": map[string]interface{}{
				"expression": "1 and 2",
				"conditions": []interface{}{
					map[string]interface{}{"field": "labels.alertname", "op": "EQ", "value": "CPU"},
					map[string]interface{}{"field": "labels.severity", "op": "IN", "value": "P1,P2"},
				},
			},
		}}
	})
	r := subscriptionUnitResource(t, apiClient)
	resourceSchema, objType := subscriptionUnitSchema(t, r)
	filterType := subscriptionUnitObjectType(t, objType, "filter_setting")
	conditionType := subscriptionUnitObjectType(t, filterType, "conditions")

	stateValue := subscriptionUnitObject(t, objType, map[string]interface{}{
		"id":                "sub-1",
		"subscription_name": "name",
		"workspace":         "ws-1",
		"filter_setting": []interface{}{
			map[string]interface{}{
				"expression": "1 and 2",
				"conditions": []interface{}{
					subscriptionUnitObject(t, conditionType, map[string]interface{}{"field": "labels.alertname", "op": "EQ", "value": "CPU"}),
					subscriptionUnitObject(t, conditionType, map[string]interface{}{"field": "labels.severity", "op": "IN", "value": "P1,P2"}),
				},
			},
		},
	})
	planValue := subscriptionUnitObject(t, objType, map[string]interface{}{
		"id":                "sub-1",
		"subscription_name": "name",
		"workspace":         "ws-1",
		"filter_setting": []interface{}{
			map[string]interface{}{
				"expression": "1 and 2",
				"conditions": []interface{}{
					subscriptionUnitObject(t, conditionType, map[string]interface{}{"field": "labels.severity", "op": "IN", "value": "P1,P2"}),
					subscriptionUnitObject(t, conditionType, map[string]interface{}{"field": "labels.alertname", "op": "EQ", "value": "CPU"}),
				},
			},
		},
	})

	ctx, cancel := context.WithTimeout(context.Background(), 400*time.Millisecond)
	defer cancel()
	req := fwresource.UpdateRequest{
		Plan:  tfsdk.Plan{Schema: resourceSchema, Raw: planValue},
		State: tfsdk.State{Schema: resourceSchema, Raw: stateValue},
	}
	resp := &fwresource.UpdateResponse{State: tfsdk.State{Schema: resourceSchema, Raw: stateValue}}
	r.Update(ctx, req, resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected Update to fail while the service still returns the previous condition order")
	}
}

// TestUnitCmsSubscriptionUpdateRequiresMemberRemoval proves a same-length
// channel whose removed member is still returned does not converge.
func TestUnitCmsSubscriptionUpdateRequiresMemberRemoval(t *testing.T) {
	apiClient, _ := subscriptionUnitClient(t, func(method, path string, body map[string]interface{}) (int, interface{}) {
		if method == http.MethodPut {
			return http.StatusOK, map[string]interface{}{"requestId": "req-1"}
		}
		return http.StatusOK, map[string]interface{}{"requestId": "req-2", "subscription": map[string]interface{}{
			"subscriptionId":   "sub-1",
			"subscriptionName": "name",
			"workspace":        "ws-1",
			"agentConfig": map[string]interface{}{
				"agentUuid": "agent-1",
				"routes": []interface{}{
					map[string]interface{}{
						"channels": []interface{}{
							map[string]interface{}{
								"channelType": "WEBHOOK",
								"receivers":   []interface{}{"r1"},
							},
						},
					},
				},
			},
		}}
	})
	r := subscriptionUnitResource(t, apiClient)
	resourceSchema, objType := subscriptionUnitSchema(t, r)
	agentType := subscriptionUnitObjectType(t, objType, "agent_config")
	routeType := subscriptionUnitObjectType(t, agentType, "routes")
	channelType := subscriptionUnitObjectType(t, routeType, "channels")

	buildAgent := func(withReceivers bool) []interface{} {
		channel := map[string]interface{}{"channel_type": "WEBHOOK"}
		if withReceivers {
			channel["receivers"] = []interface{}{"r1"}
		}
		route := subscriptionUnitObject(t, routeType, map[string]interface{}{
			"channels": []interface{}{subscriptionUnitObject(t, channelType, channel)},
		})
		return []interface{}{
			subscriptionUnitObject(t, agentType, map[string]interface{}{
				"agent_uuid": "agent-1",
				"routes":     []interface{}{route},
			}),
		}
	}

	stateValue := subscriptionUnitObject(t, objType, map[string]interface{}{
		"id":                "sub-1",
		"subscription_name": "name",
		"workspace":         "ws-1",
		"agent_config":      buildAgent(true),
	})
	planValue := subscriptionUnitObject(t, objType, map[string]interface{}{
		"id":                "sub-1",
		"subscription_name": "name",
		"workspace":         "ws-1",
		"agent_config":      buildAgent(false),
	})

	ctx, cancel := context.WithTimeout(context.Background(), 400*time.Millisecond)
	defer cancel()
	req := fwresource.UpdateRequest{
		Plan:  tfsdk.Plan{Schema: resourceSchema, Raw: planValue},
		State: tfsdk.State{Schema: resourceSchema, Raw: stateValue},
	}
	resp := &fwresource.UpdateResponse{State: tfsdk.State{Schema: resourceSchema, Raw: stateValue}}
	r.Update(ctx, req, resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected Update to fail while the service still returns the removed receivers")
	}
}

// TestUnitCmsSubscriptionChannelValidation proves the public config validator
// distinguishes unknown, null and empty enabled_sub_channels.
func TestUnitCmsSubscriptionChannelValidation(t *testing.T) {
	r := cloudmonitorservice.NewSubscriptionResource()
	configValidators, ok := r.(fwresource.ResourceWithConfigValidators)
	if !ok {
		t.Fatalf("%T does not implement resource.ResourceWithConfigValidators", r)
	}
	resourceSchema, objType := subscriptionUnitSchema(t, r)
	agentType := subscriptionUnitObjectType(t, objType, "agent_config")
	routeType := subscriptionUnitObjectType(t, agentType, "routes")
	channelType := subscriptionUnitObjectType(t, routeType, "channels")

	cases := map[string]struct {
		channel   map[string]interface{}
		wantError bool
	}{
		"contactWithoutSubChannels": {
			channel:   map[string]interface{}{"channel_type": "CONTACT", "receivers": []interface{}{"c1"}},
			wantError: true,
		},
		"contactWithNullSubChannels": {
			channel:   map[string]interface{}{"channel_type": "CONTACT", "enabled_sub_channels": nil},
			wantError: true,
		},
		"groupWithEmptySubChannels": {
			channel:   map[string]interface{}{"channel_type": "GROUP", "enabled_sub_channels": []interface{}{}},
			wantError: true,
		},
		"contactWithUnknownSubChannels": {
			channel:   map[string]interface{}{"channel_type": "CONTACT", "enabled_sub_channels": subscriptionUnitUnknown{}},
			wantError: false,
		},
		"dutyWithSubChannels": {
			channel:   map[string]interface{}{"channel_type": "DUTY", "enabled_sub_channels": []interface{}{"SMS"}},
			wantError: false,
		},
		"webhookWithoutSubChannels": {
			channel:   map[string]interface{}{"channel_type": "WEBHOOK", "receivers": []interface{}{"w1"}},
			wantError: false,
		},
	}
	for name, testCase := range cases {
		t.Run(name, func(t *testing.T) {
			channelValue := subscriptionUnitObject(t, channelType, testCase.channel)
			routeValue := subscriptionUnitObject(t, routeType, map[string]interface{}{
				"channels": []interface{}{channelValue},
			})
			configValue := subscriptionUnitObject(t, objType, map[string]interface{}{
				"subscription_name": "name",
				"workspace":         "ws-1",
				"agent_config": []interface{}{
					subscriptionUnitObject(t, agentType, map[string]interface{}{
						"agent_uuid": "agent-1",
						"routes":     []interface{}{routeValue},
					}),
				},
			})
			req := fwresource.ValidateConfigRequest{Config: tfsdk.Config{Schema: resourceSchema, Raw: configValue}}
			resp := &fwresource.ValidateConfigResponse{}
			for _, validator := range configValidators.ConfigValidators(context.Background()) {
				validator.ValidateResource(context.Background(), req, resp)
			}
			if testCase.wantError {
				if !resp.Diagnostics.HasError() {
					t.Fatal("expected a local validation error for a conditional channel type without enabled_sub_channels")
				}
				if !strings.Contains(fmt.Sprint(resp.Diagnostics), "enabled_sub_channels") {
					t.Fatalf("diagnostics %v do not point at enabled_sub_channels", resp.Diagnostics)
				}
				return
			}
			if resp.Diagnostics.HasError() {
				t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
			}
		})
	}
}

// TestUnitCmsSubscriptionDescribeRejectsOffContractPayload proves a payload
// that is neither the documented object nor the ResourceNotFound code stays an
// error instead of being read as a deletion.
func TestUnitCmsSubscriptionDescribeRejectsOffContractPayload(t *testing.T) {
	cases := map[string]interface{}{
		"nullPayload":   nil,
		"stringPayload": "not-an-object",
		"numberPayload": 42.0,
	}
	for name, subscription := range cases {
		t.Run(name, func(t *testing.T) {
			apiClient, _ := subscriptionUnitClient(t, func(method, path string, body map[string]interface{}) (int, interface{}) {
				return http.StatusOK, map[string]interface{}{"requestId": "req", "subscription": subscription}
			})
			_, err := cloudmonitorservice.NewCmsService(apiClient).DescribeSubscription(context.Background(), "sub-1", "ws-1")
			if err == nil {
				t.Fatal("expected an error for an off-contract subscription payload")
			}
			if cloudmonitorservice.IsSubscriptionNotFound(err) {
				t.Fatal("an off-contract payload must not be reported as a missing resource")
			}
		})
	}
}

// TestUnitCmsSubscriptionDescribeRejectsObjectWithoutName proves a successful
// payload that is empty or has no usable name is rejected instead of being
// flattened into state, while the contract ResourceNotFound still reads as a
// deletion and a normal object still succeeds.
func TestUnitCmsSubscriptionDescribeRejectsObjectWithoutName(t *testing.T) {
	cases := map[string]struct {
		statusCode   int
		subscription interface{}
		wantErr      bool
		wantNotFound bool
	}{
		"emptyObject": {
			statusCode:   http.StatusOK,
			subscription: map[string]interface{}{},
			wantErr:      true,
		},
		"missingName": {
			statusCode:   http.StatusOK,
			subscription: map[string]interface{}{"subscriptionId": "sub-1", "workspace": "ws-1"},
			wantErr:      true,
		},
		"emptyName": {
			statusCode:   http.StatusOK,
			subscription: map[string]interface{}{"subscriptionId": "sub-1", "subscriptionName": ""},
			wantErr:      true,
		},
		"resourceNotFound": {
			statusCode:   http.StatusNotFound,
			subscription: map[string]interface{}{"Code": "ResourceNotFound", "Message": "not found"},
			wantErr:      true,
			wantNotFound: true,
		},
		"normalObject": {
			statusCode:   http.StatusOK,
			subscription: map[string]interface{}{"subscriptionId": "sub-1", "subscriptionName": "name", "workspace": "ws-1"},
			wantErr:      false,
		},
	}
	for name, testCase := range cases {
		t.Run(name, func(t *testing.T) {
			apiClient, _ := subscriptionUnitClient(t, func(method, path string, body map[string]interface{}) (int, interface{}) {
				if testCase.statusCode != http.StatusOK {
					return testCase.statusCode, testCase.subscription
				}
				return http.StatusOK, map[string]interface{}{"requestId": "req", "subscription": testCase.subscription}
			})
			_, err := cloudmonitorservice.NewCmsService(apiClient).DescribeSubscription(context.Background(), "sub-1", "ws-1")
			if !testCase.wantErr {
				if err != nil {
					t.Fatalf("unexpected error for a usable object: %v", err)
				}
				return
			}
			if err == nil {
				t.Fatal("expected DescribeSubscription to reject the payload")
			}
			if cloudmonitorservice.IsSubscriptionNotFound(err) != testCase.wantNotFound {
				t.Fatalf("IsSubscriptionNotFound = %t, want %t (err: %v)", cloudmonitorservice.IsSubscriptionNotFound(err), testCase.wantNotFound, err)
			}
		})
	}
}

// TestUnitCmsSubscriptionEmptyConfigurationRejected proves an empty block and
// an explicit empty optional string are rejected before a request is built;
// an omitted member is the only way to express "unset".
func TestUnitCmsSubscriptionEmptyConfigurationRejected(t *testing.T) {
	r := cloudmonitorservice.NewSubscriptionResource()
	configValidators, ok := r.(fwresource.ResourceWithConfigValidators)
	if !ok {
		t.Fatalf("%T does not implement resource.ResourceWithConfigValidators", r)
	}
	resourceSchema, objType := subscriptionUnitSchema(t, r)
	agentType := subscriptionUnitObjectType(t, objType, "agent_config")
	routeType := subscriptionUnitObjectType(t, agentType, "routes")
	channelType := subscriptionUnitObjectType(t, routeType, "channels")
	windowType := subscriptionUnitObjectType(t, routeType, "effect_time_range")

	base := func(extra map[string]interface{}) map[string]interface{} {
		attributes := map[string]interface{}{
			"subscription_name": "name",
			"workspace":         "ws-1",
		}
		for key, value := range extra {
			attributes[key] = value
		}
		return attributes
	}
	agentWith := func(attributes map[string]interface{}) []interface{} {
		return []interface{}{subscriptionUnitObject(t, agentType, attributes)}
	}

	cases := map[string]struct {
		attributes map[string]interface{}
		wantError  bool
	}{
		"emptyFilterSetting": {
			attributes: base(map[string]interface{}{
				"filter_setting": []interface{}{subscriptionUnitObject(t, subscriptionUnitObjectType(t, objType, "filter_setting"), map[string]interface{}{})},
			}),
			wantError: true,
		},
		"emptyPushingSetting": {
			attributes: base(map[string]interface{}{
				"pushing_setting": []interface{}{subscriptionUnitObject(t, subscriptionUnitObjectType(t, objType, "pushing_setting"), map[string]interface{}{})},
			}),
			wantError: true,
		},
		"emptyWorkspaceFilterSetting": {
			attributes: base(map[string]interface{}{
				"workspace_filter_setting": []interface{}{subscriptionUnitObject(t, subscriptionUnitObjectType(t, objType, "workspace_filter_setting"), map[string]interface{}{})},
			}),
			wantError: true,
		},
		"emptyAgentConfig": {
			attributes: base(map[string]interface{}{
				"agent_config": agentWith(map[string]interface{}{}),
			}),
			wantError: true,
		},
		"emptyRoute": {
			attributes: base(map[string]interface{}{
				"agent_config": agentWith(map[string]interface{}{
					"agent_uuid": "agent-1",
					"routes":     []interface{}{subscriptionUnitObject(t, routeType, map[string]interface{}{})},
				}),
			}),
			wantError: true,
		},
		"emptyChannel": {
			attributes: base(map[string]interface{}{
				"agent_config": agentWith(map[string]interface{}{
					"routes": []interface{}{subscriptionUnitObject(t, routeType, map[string]interface{}{
						"channels": []interface{}{subscriptionUnitObject(t, channelType, map[string]interface{}{})},
					})},
				}),
			}),
			wantError: true,
		},
		"emptyEffectTimeRange": {
			attributes: base(map[string]interface{}{
				"agent_config": agentWith(map[string]interface{}{
					"routes": []interface{}{subscriptionUnitObject(t, routeType, map[string]interface{}{
						"effect_time_range": []interface{}{subscriptionUnitObject(t, windowType, map[string]interface{}{})},
					})},
				}),
			}),
			wantError: true,
		},
		"omittedValues": {
			attributes: base(map[string]interface{}{
				"filter_setting": []interface{}{subscriptionUnitObject(t, subscriptionUnitObjectType(t, objType, "filter_setting"), map[string]interface{}{"relation": "AND"})},
			}),
			wantError: false,
		},
		// A block whose only member is resolved at apply is not empty: a dynamic
		// value must not be mistaken for an omitted one.
		"dynamicFilterSettingExpression": {
			attributes: base(map[string]interface{}{
				"filter_setting": []interface{}{subscriptionUnitObject(t, subscriptionUnitObjectType(t, objType, "filter_setting"), map[string]interface{}{"expression": subscriptionUnitUnknown{}})},
			}),
			wantError: false,
		},
		"dynamicPushingTemplateUuid": {
			attributes: base(map[string]interface{}{
				"pushing_setting": []interface{}{subscriptionUnitObject(t, subscriptionUnitObjectType(t, objType, "pushing_setting"), map[string]interface{}{"template_uuid": subscriptionUnitUnknown{}})},
			}),
			wantError: false,
		},
		"dynamicEffectTimeRangeMember": {
			attributes: base(map[string]interface{}{
				"agent_config": agentWith(map[string]interface{}{
					"routes": []interface{}{subscriptionUnitObject(t, routeType, map[string]interface{}{
						"effect_time_range": []interface{}{subscriptionUnitObject(t, windowType, map[string]interface{}{"start_time_in_minute": subscriptionUnitUnknown{}})},
					})},
				}),
			}),
			wantError: false,
		},
		// A nested collection that is entirely unknown cannot be judged at plan
		// time and must not be rejected.
		"unknownAgentConfig": {
			attributes: base(map[string]interface{}{
				"agent_config": subscriptionUnitUnknown{},
			}),
			wantError: false,
		},
		"unknownRoutes": {
			attributes: base(map[string]interface{}{
				"agent_config": agentWith(map[string]interface{}{
					"agent_uuid": "agent-1",
					"routes":     subscriptionUnitUnknown{},
				}),
			}),
			wantError: false,
		},
		"unknownConditions": {
			attributes: base(map[string]interface{}{
				"filter_setting": []interface{}{subscriptionUnitObject(t, subscriptionUnitObjectType(t, objType, "filter_setting"), map[string]interface{}{
					"conditions": subscriptionUnitUnknown{},
				})},
			}),
			wantError: false,
		},
	}
	for name, testCase := range cases {
		t.Run(name, func(t *testing.T) {
			configValue := subscriptionUnitObject(t, objType, testCase.attributes)
			req := fwresource.ValidateConfigRequest{Config: tfsdk.Config{Schema: resourceSchema, Raw: configValue}}
			resp := &fwresource.ValidateConfigResponse{}
			for _, validator := range configValidators.ConfigValidators(context.Background()) {
				validator.ValidateResource(context.Background(), req, resp)
			}
			if testCase.wantError {
				if !resp.Diagnostics.HasError() {
					t.Fatal("expected the configuration to be rejected")
				}
				return
			}
			if resp.Diagnostics.HasError() {
				t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
			}
		})
	}
}

// TestUnitCmsSubscriptionEmptyOptionalStringRejected proves every optional
// string attribute rejects an explicit empty value and accepts both a non-empty
// value and omission.
func TestUnitCmsSubscriptionEmptyOptionalStringRejected(t *testing.T) {
	r := cloudmonitorservice.NewSubscriptionResource()
	resourceSchema, _ := subscriptionUnitSchema(t, r)

	nestedBlock := func(blocks map[string]fwschema.Block, name string) fwschema.ListNestedBlock {
		t.Helper()
		block, ok := blocks[name].(fwschema.ListNestedBlock)
		if !ok {
			t.Fatalf("%s is %T, want a list nested block", name, blocks[name])
		}
		return block
	}
	filterBlock := nestedBlock(resourceSchema.Blocks, "filter_setting")
	pushingBlock := nestedBlock(resourceSchema.Blocks, "pushing_setting")
	agentBlock := nestedBlock(resourceSchema.Blocks, "agent_config")
	routeBlock := nestedBlock(agentBlock.NestedObject.Blocks, "routes")
	channelBlock := nestedBlock(routeBlock.NestedObject.Blocks, "channels")
	windowBlock := nestedBlock(routeBlock.NestedObject.Blocks, "effect_time_range")

	stringValidators := func(attributes map[string]fwschema.Attribute, name string) []validator.String {
		t.Helper()
		attribute, ok := attributes[name].(fwschema.StringAttribute)
		if !ok {
			t.Fatalf("%s is %T, want a string attribute", name, attributes[name])
		}
		return attribute.Validators
	}
	channelsPath := path.Root("agent_config").AtListIndex(0).AtName("routes").AtListIndex(0).
		AtName("channels").AtListIndex(0)
	windowsPath := path.Root("agent_config").AtListIndex(0).AtName("routes").AtListIndex(0).
		AtName("effect_time_range").AtListIndex(0)

	cases := map[string]struct {
		path       path.Path
		validators []validator.String
	}{
		"description":      {path.Root("description"), stringValidators(resourceSchema.Attributes, "description")},
		"notifyStrategyId": {path.Root("notify_strategy_id"), stringValidators(resourceSchema.Attributes, "notify_strategy_id")},
		"expression":       {path.Root("filter_setting").AtListIndex(0).AtName("expression"), stringValidators(filterBlock.NestedObject.Attributes, "expression")},
		"templateUuid":     {path.Root("pushing_setting").AtListIndex(0).AtName("template_uuid"), stringValidators(pushingBlock.NestedObject.Attributes, "template_uuid")},
		"responsePlanId":   {path.Root("pushing_setting").AtListIndex(0).AtName("response_plan_id"), stringValidators(pushingBlock.NestedObject.Attributes, "response_plan_id")},
		"agentUuid":        {path.Root("agent_config").AtListIndex(0).AtName("agent_uuid"), stringValidators(agentBlock.NestedObject.Attributes, "agent_uuid")},
		"channelType":      {channelsPath.AtName("channel_type"), stringValidators(channelBlock.NestedObject.Attributes, "channel_type")},
		"timeZone":         {windowsPath.AtName("time_zone"), stringValidators(windowBlock.NestedObject.Attributes, "time_zone")},
	}
	for name, testCase := range cases {
		t.Run(name, func(t *testing.T) {
			if len(testCase.validators) == 0 {
				t.Fatal("the attribute carries no validator")
			}
			rejected := func(value types.String) bool {
				resp := &validator.StringResponse{}
				for _, attributeValidator := range testCase.validators {
					attributeValidator.ValidateString(context.Background(), validator.StringRequest{ConfigValue: value, Path: testCase.path}, resp)
				}
				return resp.Diagnostics.HasError()
			}
			if !rejected(types.StringValue("")) {
				t.Fatal("an explicit empty string must be rejected")
			}
			if rejected(types.StringNull()) {
				t.Fatal("an omitted value must be accepted")
			}
			if rejected(types.StringValue("value")) {
				t.Fatal("a non-empty value must be accepted")
			}
		})
	}
}

// TestUnitCmsSubscriptionUpdateRequiresConvergedEnable proves a combined update
// whose final read confirms the new body but still reports the previous enabled
// switch does not converge.
func TestUnitCmsSubscriptionUpdateRequiresConvergedEnable(t *testing.T) {
	gets := 0
	apiClient, _ := subscriptionUnitClient(t, func(method, path string, body map[string]interface{}) (int, interface{}) {
		if method == http.MethodPut {
			return http.StatusOK, map[string]interface{}{"requestId": "req"}
		}
		gets++
		enabled := true
		if gets == 1 {
			enabled = false
		}
		return http.StatusOK, map[string]interface{}{"requestId": "req", "subscription": map[string]interface{}{
			"subscriptionId":   "sub-1",
			"subscriptionName": "name",
			"workspace":        "ws-1",
			"description":      "new description",
			"enable":           enabled,
		}}
	})
	r := subscriptionUnitResource(t, apiClient)
	resourceSchema, objType := subscriptionUnitSchema(t, r)

	stateValue := subscriptionUnitObject(t, objType, map[string]interface{}{
		"id":                "sub-1",
		"subscription_name": "name",
		"workspace":         "ws-1",
		"description":       "old description",
		"enable":            true,
	})
	planValue := subscriptionUnitObject(t, objType, map[string]interface{}{
		"id":                "sub-1",
		"subscription_name": "name",
		"workspace":         "ws-1",
		"description":       "new description",
		"enable":            false,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 400*time.Millisecond)
	defer cancel()
	req := fwresource.UpdateRequest{
		Plan:  tfsdk.Plan{Schema: resourceSchema, Raw: planValue},
		State: tfsdk.State{Schema: resourceSchema, Raw: stateValue},
	}
	resp := &fwresource.UpdateResponse{State: tfsdk.State{Schema: resourceSchema, Raw: stateValue}}
	r.Update(ctx, req, resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected Update to fail while the read still reports the previous enabled switch")
	}
}

// TestUnitCmsSubscriptionCreateCancellationKeepsLateIdentity proves a create
// cancelled while the request is on the wire keeps the id returned by the late
// response, issues no second attempt and stops issuing follow-up reads.
func TestUnitCmsSubscriptionCreateCancellationKeepsLateIdentity(t *testing.T) {
	started := make(chan struct{}, 1)
	release := make(chan struct{})
	apiClient, recorder := subscriptionUnitClient(t, func(method, path string, body map[string]interface{}) (int, interface{}) {
		if method == http.MethodPost {
			started <- struct{}{}
			<-release
			return http.StatusOK, map[string]interface{}{"requestId": "req-1", "subscriptionId": subscriptionUnitServerId}
		}
		return http.StatusOK, map[string]interface{}{"requestId": "req-2", "subscription": map[string]interface{}{
			"subscriptionId":   subscriptionUnitServerId,
			"subscriptionName": "identity",
			"workspace":        "ws-1",
		}}
	})
	r := subscriptionUnitResource(t, apiClient)
	resourceSchema, objType := subscriptionUnitSchema(t, r)
	planValue := subscriptionUnitObject(t, objType, map[string]interface{}{
		"id":                subscriptionUnitUnknown{},
		"subscription_name": "identity",
		"workspace":         "ws-1",
		"enable":            true,
	})
	req := fwresource.CreateRequest{Plan: tfsdk.Plan{Schema: resourceSchema, Raw: planValue}}
	resp := &fwresource.CreateResponse{State: tfsdk.State{Schema: resourceSchema, Raw: planValue}}

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		defer close(done)
		r.Create(ctx, req, resp)
	}()
	<-started
	cancel()
	close(release)
	<-done

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected Create to stop once the operation context is cancelled")
	}
	if posts := len(recorder.requests(http.MethodPost)); posts != 1 {
		t.Fatalf("create requests = %d, want exactly one attempt", posts)
	}
	if gets := len(recorder.requests(http.MethodGet)); gets != 0 {
		t.Fatalf("follow-up reads = %d, want none after cancellation", gets)
	}
	if got := subscriptionUnitStringAttribute(t, resp.State, "id").ValueString(); got != subscriptionUnitServerId {
		t.Fatalf("state id = %q, want the id of the late successful response", got)
	}
}

// TestUnitCmsSubscriptionCreateKeepsIdentityWithUnknownRelation proves a failed
// create of an expression configuration whose computed relation is not planned
// yet returns a state with the confirmed identity and no unknown values.
func TestUnitCmsSubscriptionCreateKeepsIdentityWithUnknownRelation(t *testing.T) {
	apiClient, _ := subscriptionUnitClient(t, func(method, path string, body map[string]interface{}) (int, interface{}) {
		if method == http.MethodPost {
			return http.StatusOK, map[string]interface{}{"requestId": "req-1", "subscriptionId": subscriptionUnitServerId}
		}
		return http.StatusForbidden, map[string]interface{}{"requestId": "req-2", "Code": "Forbidden", "Message": "denied"}
	})
	r := subscriptionUnitResource(t, apiClient)
	resourceSchema, objType := subscriptionUnitSchema(t, r)
	filterType := subscriptionUnitObjectType(t, objType, "filter_setting")
	conditionType := subscriptionUnitObjectType(t, filterType, "conditions")

	planValue := subscriptionUnitObject(t, objType, map[string]interface{}{
		"id":                subscriptionUnitUnknown{},
		"subscription_name": "identity",
		"workspace":         "ws-1",
		"filter_setting": []interface{}{
			subscriptionUnitObject(t, filterType, map[string]interface{}{
				"expression": "1",
				"relation":   subscriptionUnitUnknown{},
				"conditions": []interface{}{
					subscriptionUnitObject(t, conditionType, map[string]interface{}{
						"field": "labels.alertname", "op": "EQ", "value": "CPU",
					}),
				},
			}),
		},
	})
	req := fwresource.CreateRequest{Plan: tfsdk.Plan{Schema: resourceSchema, Raw: planValue}}
	resp := &fwresource.CreateResponse{State: tfsdk.State{Schema: resourceSchema, Raw: planValue}}

	r.Create(context.Background(), req, resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected Create to fail when the first read never succeeds")
	}
	if got := subscriptionUnitStringAttribute(t, resp.State, "id").ValueString(); got != subscriptionUnitServerId {
		t.Fatalf("state id = %q, want the confirmed id", got)
	}
	if got := subscriptionUnitStringAttribute(t, resp.State, "workspace").ValueString(); got != "ws-1" {
		t.Fatalf("state workspace = %q, want ws-1", got)
	}
	subscriptionUnitAssertKnown(t, resp.State.Raw, "state")
}

func subscriptionUnitAssertKnown(t *testing.T, value tftypes.Value, at string) {
	t.Helper()
	if value.IsNull() {
		return
	}
	if !value.IsKnown() {
		t.Fatalf("%s is unknown; the failed create must not return unknown nested values", at)
	}
	switch value.Type().(type) {
	case tftypes.Object:
		var attributes map[string]tftypes.Value
		if err := value.As(&attributes); err != nil {
			t.Fatalf("%s: %v", at, err)
		}
		for name, attribute := range attributes {
			subscriptionUnitAssertKnown(t, attribute, at+"."+name)
		}
	case tftypes.List:
		var items []tftypes.Value
		if err := value.As(&items); err != nil {
			t.Fatalf("%s: %v", at, err)
		}
		for index, item := range items {
			subscriptionUnitAssertKnown(t, item, fmt.Sprintf("%s[%d]", at, index))
		}
	}
}

// TestUnitCmsSubscriptionReadUnifiesOmittedListState proves the same response
// produces the same state through Refresh and a bare Import: an omitted or
// empty list reads back as null instead of reusing the previous state.
func TestUnitCmsSubscriptionReadUnifiesOmittedListState(t *testing.T) {
	apiClient, _ := subscriptionUnitClient(t, func(method, path string, body map[string]interface{}) (int, interface{}) {
		return http.StatusOK, map[string]interface{}{"requestId": "req", "subscription": map[string]interface{}{
			"subscriptionId":   "sub-1",
			"subscriptionName": "name",
			"workspace":        "ws-1",
			"pushingSetting": map[string]interface{}{
				"templateUuid": "tpl-1",
			},
			"workspaceFilterSetting": map[string]interface{}{
				"workspaceUuids": []interface{}{},
			},
		}}
	})
	r := subscriptionUnitResource(t, apiClient)
	resourceSchema, objType := subscriptionUnitSchema(t, r)
	pushingType := subscriptionUnitObjectType(t, objType, "pushing_setting")
	workspaceFilterType := subscriptionUnitObjectType(t, objType, "workspace_filter_setting")

	refreshed := subscriptionUnitObject(t, objType, map[string]interface{}{
		"id":                "sub-1",
		"subscription_name": "name",
		"workspace":         "ws-1",
		"pushing_setting": []interface{}{
			subscriptionUnitObject(t, pushingType, map[string]interface{}{
				"alert_action_ids":   []interface{}{"previous"},
				"restore_action_ids": []interface{}{"previous"},
				"template_uuid":      "tpl-1",
			}),
		},
		"workspace_filter_setting": []interface{}{
			subscriptionUnitObject(t, workspaceFilterType, map[string]interface{}{
				"workspace_uuids": []interface{}{"previous"},
			}),
		},
	})
	imported := subscriptionUnitObject(t, objType, map[string]interface{}{
		"id":        "sub-1",
		"workspace": "ws-1",
	})

	read := func(stateValue tftypes.Value) *fwresource.ReadResponse {
		req := fwresource.ReadRequest{State: tfsdk.State{Schema: resourceSchema, Raw: stateValue}}
		resp := &fwresource.ReadResponse{State: tfsdk.State{Schema: resourceSchema, Raw: stateValue}}
		r.Read(context.Background(), req, resp)
		if resp.Diagnostics.HasError() {
			t.Fatalf("Read returned diagnostics: %v", resp.Diagnostics)
		}
		return resp
	}

	states := map[string]*fwresource.ReadResponse{"refresh": read(refreshed), "import": read(imported)}
	for name, resp := range states {
		for label, attributePath := range map[string]path.Path{
			"pushing_setting.alert_action_ids":         path.Root("pushing_setting").AtListIndex(0).AtName("alert_action_ids"),
			"pushing_setting.restore_action_ids":       path.Root("pushing_setting").AtListIndex(0).AtName("restore_action_ids"),
			"workspace_filter_setting.workspace_uuids": path.Root("workspace_filter_setting").AtListIndex(0).AtName("workspace_uuids"),
		} {
			var list types.List
			if diags := resp.State.GetAttribute(context.Background(), attributePath, &list); diags.HasError() {
				t.Fatalf("%s: reading %s: %v", name, label, diags)
			}
			if !list.IsNull() {
				t.Fatalf("%s: %s = %v, want null for an omitted member", name, label, list)
			}
		}
	}
}

// TestUnitCmsSubscriptionCreateVerifiesConfigurationWithDisable proves the
// initial disable wait requires the requested configuration and the switch
// together: a response that only reports the switch does not converge.
func TestUnitCmsSubscriptionCreateVerifiesConfigurationWithDisable(t *testing.T) {
	cases := map[string]struct {
		keepConfiguration bool
		wantError         bool
	}{
		"switchOnly":  {keepConfiguration: false, wantError: true},
		"fullPayload": {keepConfiguration: true, wantError: false},
	}
	for name, testCase := range cases {
		t.Run(name, func(t *testing.T) {
			var mu sync.Mutex
			created := map[string]interface{}{}
			disabled := false
			apiClient, recorder := subscriptionUnitClient(t, func(method, path string, body map[string]interface{}) (int, interface{}) {
				mu.Lock()
				defer mu.Unlock()
				if method == http.MethodPost {
					for key, value := range body {
						created[key] = value
					}
					return http.StatusOK, map[string]interface{}{"requestId": "req-1", "subscriptionId": subscriptionUnitServerId}
				}
				if method == http.MethodPut && strings.HasSuffix(path, "/disable") {
					disabled = true
					return http.StatusOK, map[string]interface{}{"requestId": "req-2"}
				}
				subscription := map[string]interface{}{
					"subscriptionId":   subscriptionUnitServerId,
					"subscriptionName": "disable",
					"workspace":        "ws-1",
					"enable":           !disabled,
				}
				// Before the disable the configuration is complete; after it a
				// switch-only response must not be accepted.
				if !disabled || testCase.keepConfiguration {
					for key, value := range created {
						subscription[key] = value
					}
				}
				return http.StatusOK, map[string]interface{}{"requestId": "req-3", "subscription": subscription}
			})
			r := subscriptionUnitResource(t, apiClient)
			resourceSchema, objType := subscriptionUnitSchema(t, r)
			planValue := subscriptionUnitObject(t, objType, map[string]interface{}{
				"id":                subscriptionUnitUnknown{},
				"subscription_name": "disable",
				"workspace":         "ws-1",
				"description":       "keep me",
				"enable":            false,
			})
			ctx, cancel := context.WithTimeout(context.Background(), 400*time.Millisecond)
			defer cancel()
			req := fwresource.CreateRequest{Plan: tfsdk.Plan{Schema: resourceSchema, Raw: planValue}}
			resp := &fwresource.CreateResponse{State: tfsdk.State{Schema: resourceSchema, Raw: planValue}}
			r.Create(ctx, req, resp)

			if len(recorder.requests(http.MethodPut)) == 0 {
				t.Fatal("the disable action was never called")
			}
			if testCase.wantError {
				if !resp.Diagnostics.HasError() {
					t.Fatal("expected Create to fail while the read never confirms the requested configuration")
				}
				return
			}
			if resp.Diagnostics.HasError() {
				t.Fatalf("Create returned diagnostics: %v", resp.Diagnostics)
			}
			if got := subscriptionUnitStringAttribute(t, resp.State, "description").ValueString(); got != "keep me" {
				t.Fatalf("state description = %q, want the confirmed configuration", got)
			}
		})
	}
}

// TestUnitCmsSubscriptionCreateStopsAfterCancellationOnRetryableFailure proves
// a retryable failure does not produce another attempt once the operation is
// cancelled: the SDK-level retry is disabled and the provider-level loop
// observes the same context.
func TestUnitCmsSubscriptionCreateStopsAfterCancellationOnRetryableFailure(t *testing.T) {
	attempted := make(chan struct{}, 1)
	apiClient, recorder := subscriptionUnitClient(t, func(method, path string, body map[string]interface{}) (int, interface{}) {
		if method == http.MethodPost {
			attempted <- struct{}{}
			return http.StatusInternalServerError, map[string]interface{}{"requestId": "req-1", "Code": "InternalError", "Message": "code: 500, internal error"}
		}
		return http.StatusOK, map[string]interface{}{"requestId": "req-2"}
	})
	r := subscriptionUnitResource(t, apiClient)
	resourceSchema, objType := subscriptionUnitSchema(t, r)
	planValue := subscriptionUnitObject(t, objType, map[string]interface{}{
		"id":                subscriptionUnitUnknown{},
		"subscription_name": "retryable",
		"workspace":         "ws-1",
	})
	req := fwresource.CreateRequest{Plan: tfsdk.Plan{Schema: resourceSchema, Raw: planValue}}
	resp := &fwresource.CreateResponse{State: tfsdk.State{Schema: resourceSchema, Raw: planValue}}

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		defer close(done)
		r.Create(ctx, req, resp)
	}()
	<-attempted
	cancel()
	<-done

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected Create to stop once the operation context is cancelled")
	}
	if posts := len(recorder.requests(http.MethodPost)); posts != 1 {
		t.Fatalf("create requests = %d, want exactly one attempt after cancellation", posts)
	}
	if gets := len(recorder.requests(http.MethodGet)); gets != 0 {
		t.Fatalf("follow-up reads = %d, want none after cancellation", gets)
	}
}

// TestUnitCmsSubscriptionRequestCoversExternalFixtures records the request and
// the read-back state for the pushing, workspace-filter and agent blocks. Their
// referenced action, template, response plan and agent objects have no
// Terraform resource, so the wire contract is the evidence for these members.
func TestUnitCmsSubscriptionRequestCoversExternalFixtures(t *testing.T) {
	var mu sync.Mutex
	created := map[string]interface{}{}
	apiClient, recorder := subscriptionUnitClient(t, func(method, path string, body map[string]interface{}) (int, interface{}) {
		mu.Lock()
		defer mu.Unlock()
		if method == http.MethodPost {
			for key, value := range body {
				created[key] = value
			}
			return http.StatusOK, map[string]interface{}{"requestId": "req-1", "subscriptionId": subscriptionUnitServerId}
		}
		subscription := map[string]interface{}{
			"subscriptionId":   subscriptionUnitServerId,
			"subscriptionName": "fixtures",
			"workspace":        "ws-1",
			"enable":           true,
			"userId":           "1",
			"subscriptionType": "NORMAL",
			"createTime":       "1",
			"updateTime":       "1",
		}
		for key, value := range created {
			subscription[key] = value
		}
		return http.StatusOK, map[string]interface{}{"requestId": "req-2", "subscription": subscription}
	})
	r := subscriptionUnitResource(t, apiClient)
	resourceSchema, objType := subscriptionUnitSchema(t, r)
	pushingType := subscriptionUnitObjectType(t, objType, "pushing_setting")
	workspaceFilterType := subscriptionUnitObjectType(t, objType, "workspace_filter_setting")
	agentType := subscriptionUnitObjectType(t, objType, "agent_config")
	routeType := subscriptionUnitObjectType(t, agentType, "routes")
	channelType := subscriptionUnitObjectType(t, routeType, "channels")
	windowType := subscriptionUnitObjectType(t, routeType, "effect_time_range")

	planValue := subscriptionUnitObject(t, objType, map[string]interface{}{
		"id":                subscriptionUnitUnknown{},
		"subscription_name": "fixtures",
		"workspace":         "ws-1",
		"pushing_setting": []interface{}{
			subscriptionUnitObject(t, pushingType, map[string]interface{}{
				"alert_action_ids":   []interface{}{"action-1"},
				"restore_action_ids": []interface{}{"action-2"},
				"template_uuid":      "template-1",
				"response_plan_id":   "plan-1",
			}),
		},
		"workspace_filter_setting": []interface{}{
			subscriptionUnitObject(t, workspaceFilterType, map[string]interface{}{
				"workspace_uuids": []interface{}{"workspace-uuid-1"},
			}),
		},
		"agent_config": []interface{}{
			subscriptionUnitObject(t, agentType, map[string]interface{}{
				"agent_uuid": "agent-1",
				"routes": []interface{}{
					subscriptionUnitObject(t, routeType, map[string]interface{}{
						"channels": []interface{}{
							subscriptionUnitObject(t, channelType, map[string]interface{}{
								"channel_type":         "WEBHOOK",
								"receivers":            []interface{}{"receiver-1"},
								"enabled_sub_channels": []interface{}{"EMAIL"},
							}),
						},
						"effect_time_range": []interface{}{
							subscriptionUnitObject(t, windowType, map[string]interface{}{
								"time_zone":            "Asia/Shanghai",
								"day_in_week":          []interface{}{float64(1), float64(2)},
								"start_time_in_minute": float64(0),
								"end_time_in_minute":   float64(1439),
							}),
						},
					}),
				},
			}),
		},
	})
	req := fwresource.CreateRequest{Plan: tfsdk.Plan{Schema: resourceSchema, Raw: planValue}}
	resp := &fwresource.CreateResponse{State: tfsdk.State{Schema: resourceSchema, Raw: planValue}}
	r.Create(context.Background(), req, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Create returned diagnostics: %v", resp.Diagnostics)
	}

	bodies := recorder.bodiesFor(http.MethodPost)
	if len(bodies) != 1 {
		t.Fatalf("expected one CreateSubscription call, got %v", recorder.requests(http.MethodPost))
	}
	body := bodies[0]
	pushing, ok := body["pushingSetting"].(map[string]interface{})
	if !ok || pushing["templateUuid"] != "template-1" || pushing["responsePlanId"] != "plan-1" {
		t.Fatalf("unexpected pushingSetting payload: %v", body["pushingSetting"])
	}
	workspaceFilter, ok := body["workspaceFilterSetting"].(map[string]interface{})
	if !ok || fmt.Sprint(workspaceFilter["workspaceUuids"]) != "[workspace-uuid-1]" {
		t.Fatalf("unexpected workspaceFilterSetting payload: %v", body["workspaceFilterSetting"])
	}
	agent, ok := body["agentConfig"].(map[string]interface{})
	if !ok || agent["agentUuid"] != "agent-1" {
		t.Fatalf("unexpected agentConfig payload: %v", body["agentConfig"])
	}

	stringList := func(attributePath path.Path) []string {
		t.Helper()
		var list types.List
		if diags := resp.State.GetAttribute(context.Background(), attributePath, &list); diags.HasError() {
			t.Fatalf("reading %s: %v", attributePath, diags)
		}
		result := make([]string, 0)
		for _, element := range list.Elements() {
			result = append(result, element.(types.String).ValueString())
		}
		return result
	}
	if got := fmt.Sprint(stringList(path.Root("pushing_setting").AtListIndex(0).AtName("alert_action_ids"))); got != "[action-1]" {
		t.Fatalf("state alert_action_ids = %s, want [action-1]", got)
	}
	if got := fmt.Sprint(stringList(path.Root("workspace_filter_setting").AtListIndex(0).AtName("workspace_uuids"))); got != "[workspace-uuid-1]" {
		t.Fatalf("state workspace_uuids = %s, want [workspace-uuid-1]", got)
	}
	if got := subscriptionUnitStringAttribute(t, resp.State, "subscription_name").ValueString(); got != "fixtures" {
		t.Fatalf("state subscription_name = %q", got)
	}
	channelPath := path.Root("agent_config").AtListIndex(0).AtName("routes").AtListIndex(0).AtName("channels").AtListIndex(0)
	var channelTypeValue types.String
	if diags := resp.State.GetAttribute(context.Background(), channelPath.AtName("channel_type"), &channelTypeValue); diags.HasError() {
		t.Fatalf("reading channel_type: %v", diags)
	}
	if channelTypeValue.ValueString() != "WEBHOOK" {
		t.Fatalf("state channel_type = %q, want WEBHOOK", channelTypeValue.ValueString())
	}
	if got := fmt.Sprint(stringList(channelPath.AtName("enabled_sub_channels"))); got != "[EMAIL]" {
		t.Fatalf("state enabled_sub_channels = %s, want [EMAIL]", got)
	}
}

// TestUnitCmsSubscriptionEmptyOptionalListRejected proves every optional list
// rejects an explicitly empty value and accepts omission, matching the
// documented "omit the argument to unset it" rule.
func TestUnitCmsSubscriptionEmptyOptionalListRejected(t *testing.T) {
	r := cloudmonitorservice.NewSubscriptionResource()
	resourceSchema, _ := subscriptionUnitSchema(t, r)

	nestedBlock := func(blocks map[string]fwschema.Block, name string) fwschema.ListNestedBlock {
		t.Helper()
		block, ok := blocks[name].(fwschema.ListNestedBlock)
		if !ok {
			t.Fatalf("%s is %T, want a list nested block", name, blocks[name])
		}
		return block
	}
	pushingBlock := nestedBlock(resourceSchema.Blocks, "pushing_setting")
	workspaceFilterBlock := nestedBlock(resourceSchema.Blocks, "workspace_filter_setting")
	agentBlock := nestedBlock(resourceSchema.Blocks, "agent_config")
	routeBlock := nestedBlock(agentBlock.NestedObject.Blocks, "routes")
	channelBlock := nestedBlock(routeBlock.NestedObject.Blocks, "channels")
	windowBlock := nestedBlock(routeBlock.NestedObject.Blocks, "effect_time_range")

	listValidators := func(attributes map[string]fwschema.Attribute, name string) []validator.List {
		t.Helper()
		attribute, ok := attributes[name].(fwschema.ListAttribute)
		if !ok {
			t.Fatalf("%s is %T, want a list attribute", name, attributes[name])
		}
		return attribute.Validators
	}
	channelsPath := path.Root("agent_config").AtListIndex(0).AtName("routes").AtListIndex(0).AtName("channels").AtListIndex(0)
	cases := map[string][]validator.List{
		"alertActionIds":     listValidators(pushingBlock.NestedObject.Attributes, "alert_action_ids"),
		"restoreActionIds":   listValidators(pushingBlock.NestedObject.Attributes, "restore_action_ids"),
		"workspaceUuids":     listValidators(workspaceFilterBlock.NestedObject.Attributes, "workspace_uuids"),
		"receivers":          listValidators(channelBlock.NestedObject.Attributes, "receivers"),
		"enabledSubChannels": listValidators(channelBlock.NestedObject.Attributes, "enabled_sub_channels"),
		"dayInWeek":          listValidators(windowBlock.NestedObject.Attributes, "day_in_week"),
	}
	for name, attributeValidators := range cases {
		t.Run(name, func(t *testing.T) {
			if len(attributeValidators) == 0 {
				t.Fatal("the attribute carries no validator")
			}
			rejected := func(value types.List) bool {
				resp := &validator.ListResponse{}
				for _, attributeValidator := range attributeValidators {
					attributeValidator.ValidateList(context.Background(), validator.ListRequest{ConfigValue: value, Path: channelsPath}, resp)
				}
				return resp.Diagnostics.HasError()
			}
			if !rejected(types.ListValueMust(types.StringType, []attr.Value{})) {
				t.Fatal("an explicit empty list must be rejected")
			}
			if rejected(types.ListNull(types.StringType)) {
				t.Fatal("an omitted list must be accepted")
			}
		})
	}
}

// TestUnitCmsSubscriptionDestroyCheckUsesStructuredSignal proves the destroy
// check only accepts the documented absence code of the object it queried: a
// permission error whose text mentions not-found, an answer for another object
// or free text never counts as a successful deletion.
func TestUnitCmsSubscriptionDestroyCheckUsesStructuredSignal(t *testing.T) {
	if gone, err := cmsDependencyDestroyed(nil, "WorkspaceNotExist"); gone || err != nil {
		t.Fatalf("existing object = (%t, %v), want (false, nil)", gone, err)
	}
	missing := &tea.SDKError{
		StatusCode: tea.Int(http.StatusNotFound),
		Code:       tea.String("WorkspaceNotExist"),
		Message:    tea.String("code: 404, the workspace does not exist"),
	}
	if gone, err := cmsDependencyDestroyed(missing, "WorkspaceNotExist"); !gone || err != nil {
		t.Fatalf("contract absence = (%t, %v), want (true, nil)", gone, err)
	}
	forbidden := &tea.SDKError{
		StatusCode: tea.Int(http.StatusForbidden),
		Code:       tea.String("NotFound"),
		Message:    tea.String("code: 403, NotFound: the caller is not authorized"),
	}
	if gone, err := cmsDependencyDestroyed(forbidden, "NotFound"); gone || err == nil {
		t.Fatalf("403 mentioning NotFound = (%t, %v), want (false, err)", gone, err)
	}
	otherObject := &tea.SDKError{
		StatusCode: tea.Int(http.StatusNotFound),
		Code:       tea.String("ProjectNotExist"),
		Message:    tea.String("code: 404, the project does not exist"),
	}
	if gone, err := cmsDependencyDestroyed(otherObject, "WorkspaceNotExist"); gone || err == nil {
		t.Fatalf("absence of another object = (%t, %v), want (false, err)", gone, err)
	}
	textual := fmt.Errorf("the workspace not found because the request was denied")
	if gone, err := cmsDependencyDestroyed(textual, "WorkspaceNotExist"); gone || err == nil {
		t.Fatalf("textual not-found = (%t, %v), want (false, err)", gone, err)
	}

	subscriptionTextual := fmt.Errorf("subscription not found")
	if gone, err := subscriptionDestroyed(subscriptionTextual); gone || err == nil {
		t.Fatalf("textual subscription not-found = (%t, %v), want (false, err)", gone, err)
	}
}
