package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_cms_monitor_group", &resource.Sweeper{
		Name: "alicloud_cms_monitor_group",
		F:    testSweepCmsMonitorgroup,
	})
}

func testSweepCmsMonitorgroup(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "DescribeMonitorGroups"
	request := make(map[string]interface{})
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	request["Type"] = "custom"
	var response map[string]interface{}
	conn, err := client.NewCmsClient()
	if err != nil {
		return WrapError(err)
	}

	for {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cms_monitor_groups", action, AlibabaCloudSdkGoERROR)
		}

		resp, err := jsonpath.Get("$.Resources.Resource", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Resources.Resource", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			name := item["GroupName"].(string)
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
					skip = false
					break
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Cms Monitor Group: %s ", name)
				continue
			}
			log.Printf("[INFO] Delete Cms Monitor Group: %s ", name)

			delAction := "DeleteMonitorGroup"
			conn, err := client.NewCmsClient()
			if err != nil {
				return WrapError(err)
			}
			delRequest := map[string]interface{}{
				"GroupId": fmt.Sprint(formatInt(item["GroupId"])),
			}
			_, err = conn.DoRequest(StringPointer(delAction), nil, StringPointer("POST"), StringPointer("2019-01-01"), StringPointer("AK"), nil, delRequest, &util.RuntimeOptions{})
			if err != nil {
				log.Printf("[ERROR] Failed to delete Cms Monitor Group (%s): %s", name, err)
			}
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}

func TestAccAlicloudCmsMonitorGroup_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_monitor_group.default"
	ra := resourceAttrInit(resourceId, AlicloudCmsMonitorGroupMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsMonitorGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAlicloudCmsMonitorGroup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCmsMonitorGroupBasicDependence)
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
					"monitor_group_name": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"monitor_group_name": name,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"contact_groups": []string{"${alicloud_cms_alarm_contact_group.default.alarm_contact_group_name}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"contact_groups.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"monitor_group_name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"monitor_group_name": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance-test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "acceptance-test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"contact_groups":     []string{"${alicloud_cms_alarm_contact_group.default.alarm_contact_group_name}", "${alicloud_cms_alarm_contact_group.default1.alarm_contact_group_name}"},
					"monitor_group_name": "${var.name}",
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "acceptance-test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"contact_groups.#":   "2",
						"monitor_group_name": name,
						"tags.%":             "2",
						"tags.Created":       "TF-update",
						"tags.For":           "acceptance-test-update",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudCmsMonitorGroup_ByResourceGroupId(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_monitor_group.default"
	ra := resourceAttrInit(resourceId, AlicloudCmsMonitorGroupMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsMonitorGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAlicloudCmsMonitorGroup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCmsMonitorGroupBasicDependence1)
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
					"contact_groups":      []string{"${alicloud_cms_alarm_contact_group.default.alarm_contact_group_name}"},
					"resource_group_id":   "${alicloud_resource_manager_resource_group.default.id}",
					"resource_group_name": "${alicloud_resource_manager_resource_group.default.resource_group_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"contact_groups.#":    "1",
						"resource_group_id":   CHECKSET,
						"resource_group_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"monitor_group_name": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"monitor_group_name": name,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"resource_group_id", "resource_group_name"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"contact_groups": []string{"${alicloud_cms_alarm_contact_group.default.alarm_contact_group_name}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"contact_groups.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"monitor_group_name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"monitor_group_name": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance-test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "acceptance-test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"contact_groups":     []string{"${alicloud_cms_alarm_contact_group.default.alarm_contact_group_name}", "${alicloud_cms_alarm_contact_group.default1.alarm_contact_group_name}"},
					"monitor_group_name": "${var.name}",
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "acceptance-test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"contact_groups.#":   "2",
						"monitor_group_name": name,
						"tags.%":             "2",
						"tags.Created":       "TF-update",
						"tags.For":           "acceptance-test-update",
					}),
				),
			},
		},
	})
}

var AlicloudCmsMonitorGroupMap = map[string]string{}

func AlicloudCmsMonitorGroupBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

resource "alicloud_cms_alarm_contact_group" "default" {
alarm_contact_group_name = var.name
}

resource "alicloud_cms_alarm_contact_group" "default1" {
alarm_contact_group_name = "${var.name}_update"
}
`, name)
}

func AlicloudCmsMonitorGroupBasicDependence1(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

resource "alicloud_cms_alarm_contact_group" "default" {
	alarm_contact_group_name = var.name
}

resource "alicloud_cms_alarm_contact_group" "default1" {
	alarm_contact_group_name = "${var.name}_update"
}

resource "alicloud_resource_manager_resource_group" "default" {
	resource_group_name = var.name
	display_name        = var.name
}
`, name)
}
