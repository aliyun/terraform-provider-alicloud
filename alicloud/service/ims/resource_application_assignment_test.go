package ims_test

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// The IMS application itself is not managed by a Terraform resource, so each
// acceptance test builds its own application through CreateApplication. A
// sweeper is registered so applications leaked by failed or interrupted runs
// (named with the tf-testAccImsApp prefix) are reclaimed on the next sweep.
func init() {
	resource.AddTestSweepers("alicloud_ims_application", &resource.Sweeper{
		Name: "alicloud_ims_application",
		F:    testSweepImsApplications,
	})
}

// testSweepImsApplications deletes IMS applications whose AppName carries the
// tf-testAccImsApp prefix. ListApplications returns the full set in one call,
// so a single traversal is enough; transient delete failures are logged and
// retried by the next sweep rather than aborting the pass.
func testSweepImsApplications(region string) error {
	rawClient, err := sweep.SharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("getting Alicloud client for IMS application sweeper: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	response, err := client.RpcPost("Ims", "2019-08-15", "ListApplications", nil, map[string]interface{}{}, true)
	if err != nil {
		log.Printf("[ERROR] Failed to list IMS applications during sweep: %s", err)
		return nil
	}
	container, _ := response["Applications"].(map[string]interface{})
	items, _ := container["Application"].([]interface{})
	for _, raw := range items {
		item, _ := raw.(map[string]interface{})
		if item == nil {
			continue
		}
		appName, _ := item["AppName"].(string)
		if !strings.HasPrefix(appName, "tf-testAccImsApp") {
			continue
		}
		appId, _ := item["AppId"].(string)
		log.Printf("[INFO] Sweeping IMS application %q (AppId=%s)", appName, appId)
		if _, err := client.RpcPost("Ims", "2019-08-15", "DeleteApplication", nil, map[string]interface{}{"AppId": appId}, true); err != nil {
			log.Printf("[ERROR] Failed to delete IMS application %s during sweep: %s", appId, err)
		}
	}
	return nil
}

func TestAccAliCloudImsApplicationAssignment_basic(t *testing.T) {
	// The IMS application itself is not managed by a Terraform resource, so the
	// AppId it provides must be known before the step configs are formatted.
	// It is therefore built before resource.Test, which means the TF_ACC Skip
	// inside resource.Test can no longer protect the prerequisite CreateApplication
	// call. Replicate the acceptance-test gate here, ahead of any shared client
	// construction or cloud request, so the prerequisite only fires under the
	// same opt-in that governs resource.Test (R7.17).
	if os.Getenv("TF_ACC") == "" {
		t.Skip("TF_ACC is not set")
	}
	// The shared sweeper client reads the same credentials and region as the
	// provider under test, so the application lives in the account and region
	// that the assignment will target.
	region := os.Getenv("ALICLOUD_REGION")
	if region == "" {
		region = "cn-beijing"
	}
	r := tfacctest.RandIntRange(10000, 99999)
	appName := fmt.Sprintf("tf-testAccImsApp%d", r)

	rawClient, err := sweep.SharedClientForRegion(region)
	if err != nil {
		t.Fatalf("building shared client for IMS application prereq: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	appId := testAccCreateImsApplication(t, client, appName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccCheckImsApplicationAssignmentDestroy(client, appId),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccImsApplicationAssignmentConfig(r, appId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_ims_application_assignment.default", "id"),
					resource.TestCheckResourceAttr("alicloud_ims_application_assignment.default", "app_id", appId),
					resource.TestCheckResourceAttr("alicloud_ims_application_assignment.default", "principal_type", "user"),
					resource.TestCheckResourceAttrSet("alicloud_ims_application_assignment.default", "create_time"),
					resource.TestCheckResourceAttrPair("alicloud_ims_application_assignment.default", "principal_name", "alicloud_ims_user.default", "user_principal_name"),
					testAccCheckImsApplicationAssignmentID("alicloud_ims_application_assignment.default", appId),
				),
			},
			{
				Config: testAccImsApplicationAssignmentConfigUser2(r, appId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_ims_application_assignment.default", "id"),
					resource.TestCheckResourceAttr("alicloud_ims_application_assignment.default", "app_id", appId),
					resource.TestCheckResourceAttr("alicloud_ims_application_assignment.default", "principal_type", "user"),
					resource.TestCheckResourceAttrPair("alicloud_ims_application_assignment.default", "principal_name", "alicloud_ims_user.second", "user_principal_name"),
					testAccCheckImsApplicationAssignmentID("alicloud_ims_application_assignment.default", appId),
				),
			},
			{
				ResourceName:      "alicloud_ims_application_assignment.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// testAccCreateImsApplication provisions an OnDemand WebApp and returns its
// AppId. The application is removed by a t.Cleanup so it cannot outlive the
// test even when the assignment steps fail partway through.
func testAccCreateImsApplication(t *testing.T, client *connectivity.AliyunClient, appName string) string {
	t.Helper()
	request := map[string]interface{}{
		"AppName":        appName,
		"DisplayName":    appName,
		"AppType":        "WebApp",
		"AssignmentType": "OnDemand",
	}
	response, err := client.RpcPost("Ims", "2019-08-15", "CreateApplication", nil, request, true)
	if err != nil {
		t.Fatalf("creating IMS application %q: %s", appName, err)
	}
	app, _ := response["Application"].(map[string]interface{})
	appId, _ := app["AppId"].(string)
	if appId == "" {
		t.Fatalf("CreateApplication returned no AppId: %v", response)
	}
	t.Cleanup(func() {
		if _, err := client.RpcPost("Ims", "2019-08-15", "DeleteApplication", nil, map[string]interface{}{"AppId": appId}, true); err != nil {
			log.Printf("[ERROR] Failed to delete IMS application %s during cleanup: %s", appId, err)
		}
	})
	return appId
}

// testAccCheckImsApplicationAssignmentID asserts that the stored id is the
// composite <app_id>:<principal_type>:<principal_name> and that app_id matches
// the application created for this run.
func testAccCheckImsApplicationAssignmentID(name, appId string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set for %s", name)
		}
		gotAppId := rs.Primary.Attributes["app_id"]
		if gotAppId != appId {
			return fmt.Errorf("app_id = %q, want %q", gotAppId, appId)
		}
		want := fmt.Sprintf("%s:%s:%s",
			rs.Primary.Attributes["app_id"],
			rs.Primary.Attributes["principal_type"],
			rs.Primary.Attributes["principal_name"])
		if rs.Primary.ID != want {
			return fmt.Errorf("id = %q, want composite %q", rs.Primary.ID, want)
		}
		return nil
	}
}

// testAccCheckImsApplicationAssignmentDestroy confirms that no assignment to a
// test principal remains on the application after the resources are destroyed.
func testAccCheckImsApplicationAssignmentDestroy(client *connectivity.AliyunClient, appId string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		request := map[string]interface{}{
			"AppId":      appId,
			"MaxResults": 100,
		}
		for {
			response, err := client.RpcPost("Ims", "2019-08-15", "ListApplicationAssignments", nil, request, true)
			if err != nil {
				return fmt.Errorf("ListApplicationAssignments failed during destroy check: %s", err)
			}
			container, _ := response["ApplicationAssignment"].(map[string]interface{})
			items, _ := container["ApplicationAssignmentInfo"].([]interface{})
			for _, raw := range items {
				item, _ := raw.(map[string]interface{})
				if item == nil {
					continue
				}
				pn, _ := item["PrincipalName"].(string)
				if strings.Contains(pn, "tf-testAccImsAssign") {
					return fmt.Errorf("IMS application assignment still exists: app_id=%s principal_name=%s", appId, pn)
				}
			}
			next, _ := response["NextToken"].(string)
			if next == "" {
				break
			}
			request["NextToken"] = next
		}
		return nil
	}
}

func testAccImsApplicationAssignmentConfig(r int, appId string) string {
	return fmt.Sprintf(`
data "alicloud_ims_default_domain" "default" {}

resource "alicloud_ims_user" "default" {
  user_principal_name = "tf-testAccImsAssign%[1]d@${data.alicloud_ims_default_domain.default.default_domain}"
  display_name        = "tf-%[1]d"
}

resource "alicloud_ims_application_assignment" "default" {
  app_id         = "%[2]s"
  principal_type = "user"
  principal_name = alicloud_ims_user.default.user_principal_name
}
`, r, appId)
}

// testAccImsApplicationAssignmentConfigUser2 swaps the assigned principal to a
// second RAM user, exercising the ForceNew replacement of principal_name.
func testAccImsApplicationAssignmentConfigUser2(r int, appId string) string {
	return fmt.Sprintf(`
data "alicloud_ims_default_domain" "default" {}

resource "alicloud_ims_user" "default" {
  user_principal_name = "tf-testAccImsAssign%[1]d@${data.alicloud_ims_default_domain.default.default_domain}"
  display_name        = "tf-%[1]d"
}

resource "alicloud_ims_user" "second" {
  user_principal_name = "tf-testAccImsAssign%[1]db@${data.alicloud_ims_default_domain.default.default_domain}"
  display_name        = "tf-%[1]db"
}

resource "alicloud_ims_application_assignment" "default" {
  app_id         = "%[2]s"
  principal_type = "user"
  principal_name = alicloud_ims_user.second.user_principal_name
}
`, r, appId)
}
