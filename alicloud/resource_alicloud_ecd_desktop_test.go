package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers(
		"alicloud_ecd_desktop",
		&resource.Sweeper{
			Name: "alicloud_ecd_desktop",
			F:    testSweepEcdDesktop,
		})
}

func testSweepEcdDesktop(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	action := "DescribeDesktops"
	request := map[string]interface{}{}

	request["MaxResults"] = PageSizeLarge
	var response map[string]interface{}
	conn, err := client.NewGwsecdClient()
	if err != nil {
		log.Printf("[ERROR] %s get an error: %#v", action, err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-09-30"), StringPointer("AK"), nil, request, &runtime)
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
		resp, err := jsonpath.Get("$.Desktops", response)
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})

			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["DesktopName"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping EcdDesktop: %s", item["DesktopName"].(string))
				continue
			}

			action := "DeleteDesktops"
			request := map[string]interface{}{
				"DesktopId": []string{item["DesktopId"].(string)},
			}

			_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-09-30"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				log.Printf("[ERROR] Failed to delete EcdDesktop (%s): %s", item["DesktopName"].(string), err)
			}
			log.Printf("[INFO] Delete EcdDesktop success: %s ", item["DesktopName"].(string))
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	return nil
}

func TestAccAlicloudECDDesktop_basic0(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.EcdUserSupportRegions)
	resourceId := "alicloud_ecd_desktop.default"
	ra := resourceAttrInit(resourceId, AlicloudECDDesktopMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcdService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcdDesktop")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testaccdesktop%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudECDDesktopBasicDependence0)
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
					"office_site_id":  "${alicloud_ecd_simple_office_site.default.id}",
					"policy_group_id": "${alicloud_ecd_policy_group.default.id}",
					"bundle_id":       "${data.alicloud_ecd_bundles.default.bundles.1.id}",
					"desktop_name":    name,
					"amount":          "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"office_site_id":  CHECKSET,
						"policy_group_id": CHECKSET,
						"bundle_id":       CHECKSET,
						"desktop_name":    name,
						"amount":          "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_group_id": "${alicloud_ecd_policy_group.default0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"desktop_name": fmt.Sprintf("tf-testaccdesknewname%d", rand),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"desktop_name": fmt.Sprintf("tf-testaccdesknewname%d", rand),
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
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF1",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Running",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Running",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status":             "Stopped",
					"stopped_mode":       "KeepCharging",
					"root_disk_size_gib": "200",
					"user_disk_size_gib": "200",
					"policy_group_id":    "${alicloud_ecd_policy_group.default.id}",
					"desktop_name":       fmt.Sprintf("tf-testaccdesknewname%d", rand),
					"end_user_ids":       []string{"${alicloud_ecd_user.default.id}", "${alicloud_ecd_user.default1.id}"},
					"desktop_type":       "eds.graphics.24c1t4",
					"tags": map[string]string{
						"Created": "TF1",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":             "Stopped",
						"stopped_mode":       "KeepCharging",
						"root_disk_size_gib": "200",
						"user_disk_size_gib": "200",
						"policy_group_id":    CHECKSET,
						"desktop_name":       fmt.Sprintf("tf-testaccdesknewname%d", rand),
						"end_user_ids.#":     "2",
						"desktop_type":       "eds.graphics.24c1t4",
						"tags.%":             "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"root_disk_size_gib", "auto_renew", "period", "bundle_id", "user_assign_mode", "user_disk_size_gib", "host_name", "period_unit", "stopped_mode", "amount", "auto_pay", "desktop_id"},
			},
		},
	})
}

func TestAccAlicloudECDDesktop_basic1(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.EcdUserSupportRegions)
	resourceId := "alicloud_ecd_desktop.default"
	ra := resourceAttrInit(resourceId, AlicloudECDDesktopMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcdService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcdDesktop")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testaccdesktop%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudECDDesktopBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithTime(t, []int{1})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"office_site_id":  "${data.alicloud_ecd_simple_office_sites.default.sites.0.id}",
					"policy_group_id": "${data.alicloud_ecd_policy_groups.ids.groups.2.id}",
					"bundle_id":       "${data.alicloud_ecd_bundles.default.bundles.1.id}",
					"desktop_name":    name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"office_site_id":  CHECKSET,
						"policy_group_id": CHECKSET,
						"bundle_id":       CHECKSET,
						"desktop_name":    name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type": "Subscription",
					"period":       "1",
					"period_unit":  "Week",
					"auto_pay":     "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type": "Subscription",
						"period":       "1",
						"period_unit":  "Week",
						"auto_pay":     "true",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"root_disk_size_gib", "auto_renew", "period", "bundle_id", "user_assign_mode", "user_disk_size_gib", "host_name", "period_unit", "stopped_mode", "amount", "auto_pay", "desktop_id"},
			},
		},
	})
}

func TestAccAlicloudECDDesktop_basic2(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.EcdUserSupportRegions)
	resourceId := "alicloud_ecd_desktop.default"
	ra := resourceAttrInit(resourceId, AlicloudECDDesktopMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcdService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcdDesktop")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testaccdesktop%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudECDDesktopBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithTime(t, []int{1})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"office_site_id":   "${data.alicloud_ecd_simple_office_sites.default.sites.0.id}",
					"policy_group_id":  "${data.alicloud_ecd_policy_groups.ids.groups.2.id}",
					"bundle_id":        "${data.alicloud_ecd_bundles.default.bundles.1.id}",
					"desktop_name":     name,
					"payment_type":     "Subscription",
					"auto_renew":       "true",
					"user_assign_mode": "ALL",
					"period":           "1",
					"period_unit":      "Week",
					"auto_pay":         "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"office_site_id":   CHECKSET,
						"policy_group_id":  CHECKSET,
						"bundle_id":        CHECKSET,
						"desktop_name":     name,
						"payment_type":     "Subscription",
						"auto_renew":       "true",
						"user_assign_mode": "ALL",
						"period":           "1",
						"period_unit":      "Week",
						"auto_pay":         "true",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"root_disk_size_gib", "auto_renew", "period", "bundle_id", "user_assign_mode", "user_disk_size_gib", "host_name", "period_unit", "stopped_mode", "amount", "auto_pay", "desktop_id"},
			},
		},
	})
}

var AlicloudECDDesktopMap0 = map[string]string{
	"user_disk_size_gib": NOSET,
	"host_name":          NOSET,
	"period_unit":        NOSET,
	"desktop_id":         NOSET,
	"desktop_name":       CHECKSET,
	"stopped_mode":       NOSET,
	"office_site_id":     CHECKSET,
	"policy_group_id":    CHECKSET,
	"desktop_type":       CHECKSET,
	"tags.%":             NOSET,
	"amount":             NOSET,
	"end_user_ids.#":     NOSET,
	"payment_type":       CHECKSET,
	"auto_pay":           NOSET,
	"root_disk_size_gib": NOSET,
	"auto_renew":         NOSET,
	"period":             NOSET,
	"status":             CHECKSET,
	"bundle_id":          NOSET,
	"user_assign_mode":   NOSET,
}

func AlicloudECDDesktopBasicDependence0(name string) string {
	rand := acctest.RandIntRange(10000, 99999)
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

data "alicloud_ecd_bundles" "default"{
  bundle_type = "SYSTEM"
}

resource "alicloud_ecd_simple_office_site" "default" {
  cidr_block = "172.16.0.0/12"
  desktop_access_type = "Internet"
  office_site_name    = var.name
}

resource "alicloud_ecd_policy_group" "default" {
  policy_group_name = var.name
  clipboard = "readwrite"
  local_drive = "read"
  authorize_access_policy_rules{
    description= var.name
    cidr_ip=     "1.2.3.4/24"
  }
  authorize_security_policy_rules  {
    type=        "inflow"
    policy=      "accept"
    description=  var.name
    port_range= "80/80"
    ip_protocol= "TCP"
    priority=    "1"
    cidr_ip=     "0.0.0.0/0"
  }
}

resource "alicloud_ecd_policy_group" "default0" {
  policy_group_name = var.name
  clipboard = "readwrite"
  local_drive = "read"
  authorize_access_policy_rules{
    description= var.name
    cidr_ip=     "1.2.3.4/24"
  }
  authorize_security_policy_rules  {
    type=        "inflow"
    policy=      "accept"
    description=  var.name
    port_range= "80/80"
    ip_protocol= "TCP"
    priority=    "1"
    cidr_ip=     "0.0.0.0/0"
  }
}

resource "alicloud_ecd_user" "default" {
	end_user_id = "tf_testaccecduser%d"
	email       = "hello.%d@aaa.com"
	phone       = "158016%d"
	password    = "%d"
}

resource "alicloud_ecd_user" "default1" {
	end_user_id = "tf_testaccecduser%d"
	email       = "hello.%d@aaa.com"
	phone       = "158016%d"
	password    = "%d"
}

`, name, rand, rand, rand, rand, rand, rand, rand, rand)
}

func AlicloudECDDesktopBasicDependence1(name string) string {

	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

data "alicloud_ecd_bundles" "default"{
  bundle_type = "SYSTEM"
}

data "alicloud_ecd_simple_office_sites" "default" {}

data "alicloud_ecd_policy_groups" "ids" {}



`, name)
}
