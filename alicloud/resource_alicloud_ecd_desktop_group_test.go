package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers(
		"alicloud_ecd_desktop_group",
		&resource.Sweeper{
			Name: "alicloud_ecd_desktop_group",
			F:    testSweepEcdDesktopGroup,
		})
}

func testSweepEcdDesktopGroup(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	action := "DescribeDesktopGroups"
	request := map[string]interface{}{
		"RegionId":   region,
		"MaxResults": PageSizeLarge,
	}

	var response map[string]interface{}
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("ecd", "2020-09-30", action, nil, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			log.Printf("[ERROR] %s get an error: %#v", action, err)
			return nil
		}
		resp, err := jsonpath.Get("$.DesktopGroups", response)
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})

			skip := true
			if !sweepAll() {
				for _, prefix := range prefixes {
					if strings.HasPrefix(strings.ToLower(fmt.Sprint(item["DesktopGroupName"])), strings.ToLower(prefix)) {
						skip = false
					}
				}
				if skip {
					log.Printf("[INFO] Skipping EcdDesktopGroup: %s", fmt.Sprint(item["DesktopGroupName"]))
					continue
				}
			}
			action := "DeleteDesktopGroup"
			request := map[string]interface{}{
				"RegionId":       region,
				"DesktopGroupId": fmt.Sprint(item["DesktopGroupId"]),
			}

			_, err = client.RpcPost("ecd", "2020-09-30", action, nil, request, false)
			if err != nil {
				log.Printf("[ERROR] Failed to delete EcdDesktopGroup (%s): %s", fmt.Sprint(item["DesktopGroupId"]), err)
			}
			log.Printf("[INFO] Delete EcdDesktopGroup success: %s ", fmt.Sprint(item["DesktopGroupId"]))
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	return nil
}

// testAccCleanupEcdSimpleOfficeSite deletes the simple office site named name once the test finishes.
// Destroying alicloud_ecd_simple_office_site does not call DeleteOfficeSites when the office site has
// EnableInternetAccess set (see resourceAlicloudEcdSimpleOfficeSiteDelete), so every run would otherwise
// leave an office site behind and exhaust the office site quota of the test account.
func testAccCleanupEcdSimpleOfficeSite(t *testing.T, name string) {
	t.Cleanup(func() {
		client, ok := testAccProvider.Meta().(*connectivity.AliyunClient)
		if !ok || client == nil {
			return
		}
		request := map[string]interface{}{
			"RegionId":       client.RegionId,
			"OfficeSiteType": "SIMPLE",
			"MaxResults":     PageSizeLarge,
		}
		officeSiteIds := make([]string, 0)
		for {
			response, err := client.RpcPost("ecd", "2020-09-30", "DescribeOfficeSites", nil, request, true)
			if err != nil {
				t.Logf("[WARN] Failed to list the office sites to clean up %s: %s", name, err)
				return
			}
			officeSites, _ := response["OfficeSites"].([]interface{})
			for _, raw := range officeSites {
				if officeSite, ok := raw.(map[string]interface{}); ok && fmt.Sprint(officeSite["Name"]) == name {
					officeSiteIds = append(officeSiteIds, fmt.Sprint(officeSite["OfficeSiteId"]))
				}
			}
			if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
				request["NextToken"] = nextToken
			} else {
				break
			}
		}
		for _, officeSiteId := range officeSiteIds {
			deleteRequest := map[string]interface{}{
				"RegionId":     client.RegionId,
				"OfficeSiteId": []string{officeSiteId},
			}
			err := resource.Retry(5*time.Minute, func() *resource.RetryError {
				_, err := client.RpcPost("ecd", "2020-09-30", "DeleteOfficeSites", nil, deleteRequest, false)
				if err != nil {
					if NeedRetry(err) {
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			if err != nil {
				t.Logf("[WARN] Failed to delete the office site %s (%s): %s", officeSiteId, name, err)
				continue
			}
			t.Logf("[INFO] Deleted the office site %s (%s)", officeSiteId, name)
		}
	})
}

func TestAccAliCloudEcdDesktopGroup_basic0(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.EcdUserSupportRegions)
	resourceId := "alicloud_ecd_desktop_group.default"
	ra := resourceAttrInit(resourceId, AlicloudECDDesktopGroupMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcdServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcdDesktopGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := 10000 + acctest.RandIntRange(0, 89999)
	name := fmt.Sprintf("tf-testaccdesktopgroup%d", rand)
	testAccCleanupEcdSimpleOfficeSite(t, name)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudECDDesktopGroupBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"office_site_id":     "${alicloud_ecd_simple_office_site.default.id}",
					"policy_group_id":    "${alicloud_ecd_policy_group.default.id}",
					"bundle_id":          "${data.alicloud_ecd_bundles.default.bundles.0.id}",
					"end_user_ids":       []string{"${alicloud_ecd_user.default.id}"},
					"desktop_group_name": name,
					"comments":           "test-comments",
					"keep_duration":      "180000",
					"min_desktops_count": "0",
					"max_desktops_count": "1",
					"allow_auto_setup":   "0",
					"allow_buffer_count": "0",
					"directory_id":       "",
					"pay_type":           "PostPaid",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"office_site_id":     CHECKSET,
						"policy_group_id":    CHECKSET,
						"bundle_id":          CHECKSET,
						"end_user_ids.#":     "1",
						"desktop_group_name": name,
						"comments":           "test-comments",
						"keep_duration":      "180000",
						"min_desktops_count": "0",
						"max_desktops_count": "1",
						"allow_auto_setup":   "0",
						"allow_buffer_count": "0",
						"pay_type":           "PostPaid",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"desktop_group_name": fmt.Sprintf("tf-testaccdesktopgroupnew%d", rand),
					"comments":           "test-comments-update",
					"keep_duration":      "360000",
					"max_desktops_count": "2",
					"allow_buffer_count": "1",
					"scale_strategy_id":  "",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"desktop_group_name": fmt.Sprintf("tf-testaccdesktopgroupnew%d", rand),
						"comments":           "test-comments-update",
						"keep_duration":      "360000",
						"max_desktops_count": "2",
						"allow_buffer_count": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_group_id":    "${alicloud_ecd_policy_group.default0.id}",
					"bundle_id":          "${data.alicloud_ecd_bundles.default.bundles.1.id}",
					"allow_auto_setup":   "1",
					"min_desktops_count": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_group_id":    CHECKSET,
						"bundle_id":          CHECKSET,
						"allow_auto_setup":   "1",
						"min_desktops_count": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"end_user_ids": []string{"${alicloud_ecd_user.default.id}", "${alicloud_ecd_user.default1.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"end_user_ids.#": "2",
					}),
				),
			},
			{
				// Reordering the authorized users is reported as a change, because end_user_ids is a TypeList.
				Config: testAccConfig(map[string]interface{}{
					"end_user_ids": []string{"${alicloud_ecd_user.default1.id}", "${alicloud_ecd_user.default.id}"},
				}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				// Applying that order converges: Read keeps the configured order, so the post-apply plan is empty.
				Config: testAccConfig(map[string]interface{}{
					"end_user_ids": []string{"${alicloud_ecd_user.default1.id}", "${alicloud_ecd_user.default.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"end_user_ids.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"end_user_ids": []string{"${alicloud_ecd_user.default1.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"end_user_ids.#": "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"scale_strategy_id"},
			},
		},
	})
}

var AlicloudECDDesktopGroupMap0 = map[string]string{
	"office_site_id":     CHECKSET,
	"policy_group_id":    CHECKSET,
	"bundle_id":          CHECKSET,
	"end_user_ids.#":     CHECKSET,
	"desktop_group_name": CHECKSET,
	"cpu":                CHECKSET,
	"create_time":        CHECKSET,
	"creator":            CHECKSET,
	"memory":             CHECKSET,
	"office_site_name":   CHECKSET,
	"office_site_type":   CHECKSET,
	"pay_type":           CHECKSET,
	"policy_group_name":  CHECKSET,
}

func AlicloudECDDesktopGroupBasicDependence0(name string) string {
	rand := 10000 + acctest.RandIntRange(0, 89999)
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_ecd_bundles" "default" {
  bundle_type = "SYSTEM"
}

resource "alicloud_ecd_simple_office_site" "default" {
  cidr_block          = "172.16.0.0/12"
  desktop_access_type = "Internet"
  office_site_name    = var.name
}

resource "alicloud_ecd_policy_group" "default" {
  policy_group_name = var.name
  clipboard         = "readwrite"
  local_drive       = "read"
  authorize_access_policy_rules {
    description = var.name
    cidr_ip     = "1.2.3.4/24"
  }
  authorize_security_policy_rules {
    type        = "inflow"
    policy      = "accept"
    description = var.name
    port_range  = "80/80"
    ip_protocol = "TCP"
    priority    = "1"
    cidr_ip     = "0.0.0.0/0"
  }
}

resource "alicloud_ecd_policy_group" "default0" {
  policy_group_name = var.name
  clipboard         = "readwrite"
  local_drive       = "read"
  authorize_access_policy_rules {
    description = var.name
    cidr_ip     = "1.2.3.4/24"
  }
  authorize_security_policy_rules {
    type        = "inflow"
    policy      = "accept"
    description = var.name
    port_range  = "80/80"
    ip_protocol = "TCP"
    priority    = "1"
    cidr_ip     = "0.0.0.0/0"
  }
}

resource "alicloud_ecd_user" "default" {
  end_user_id = "tf_testacc-dg-u1-%d"
  email       = "hello.dg1.%d@aaa.com"
  phone       = "158016%d"
  password    = "%d"
}

resource "alicloud_ecd_user" "default1" {
  end_user_id = "tf_testacc-dg-u2-%d"
  email       = "hello.dg2.%d@aaa.com"
  phone       = "139017%d"
  password    = "%d"
}
`, name, rand, rand, rand, rand, rand, rand, rand, rand)
}
