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
	resource.AddTestSweepers("alicloud_iot_device_group", &resource.Sweeper{
		Name: "alicloud_iot_device_group",
		F:    testSweepDeviceGroup,
	})
}

func testSweepDeviceGroup(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "Error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	request := map[string]interface{}{
		"CurrentPage": 1,
		"PageSize":    PageSizeSmall,
	}
	var response map[string]interface{}
	action := "QueryDeviceGroupList"
	conn, err := client.NewIotClient()
	if err != nil {
		return WrapError(err)
	}

	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-01-20"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			log.Printf("[ERROR] %s got an error: %#v", action, err)
			return nil
		}
		addDebug(action, response, request)

		resp, err := jsonpath.Get("$.Data.GroupInfo", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Data.GroupInfo", response)
		}

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			skip := true
			if _, ok := item["GroupName"]; !ok {
				skip = false
			} else {
				for _, prefix := range prefixes {
					if strings.HasPrefix(strings.ToLower(fmt.Sprintf("%v", item["GroupName"])), strings.ToLower(prefix)) {
						skip = false
						break
					}
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Iot Device Group: %v", item["GroupName"])
				continue
			}
			action := "DeleteDeviceGroup"
			request := map[string]interface{}{
				"GroupId": item["GroupId"],
			}
			_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-01-20"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				log.Printf("[ERROR] Failed to delete Iot Device Group (%v): %s", item["GroupName"], err)
				continue
			}

			log.Printf("[INFO] Delete Iot Device Group Success: %v ", item["GroupName"])
		}
		if len(result) < PageSizeSmall {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}

func TestAccAlicloudIotDeviceGroup_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_iot_device_group.default"
	ra := resourceAttrInit(resourceId, AlicloudIotDeviceGroupMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &IotService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeIotDeviceGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 1000)
	name := fmt.Sprintf("tf_testacciot%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudIotDeviceGroupBasicDependence0)
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
					"group_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"group_desc": name + "group_desc",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_desc": name + "group_desc",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"group_desc": name + "group_desc",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_desc": name + "group_desc",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"iot_instance_id", "super_group_id"},
			},
		},
	})
}

var AlicloudIotDeviceGroupMap0 = map[string]string{}

func AlicloudIotDeviceGroupBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
`, name)
}
