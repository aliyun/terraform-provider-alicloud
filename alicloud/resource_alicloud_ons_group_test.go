package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_ons_group", &resource.Sweeper{
		Name: "alicloud_ons_group",
		F:    testSweepOnsGroup,
	})
}

func testSweepOnsGroup(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"GID-tf-testAcc",
		"GID_tf-testacc",
		"CID-tf-testAcc",
		"CID_tf-testacc",
	}

	action := "OnsInstanceInServiceList"
	request := make(map[string]interface{})
	var response map[string]interface{}
	conn, err := client.NewOnsClient()
	if err != nil {
		return WrapError(err)
	}
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-02-14"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
	if err != nil {
		log.Printf("[ERROR] Failed to retrieve ons instance in service list: %s", err)
	}
	resp, err := jsonpath.Get("$.Data.InstanceVO", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Data.InstanceVO", response)
	}

	var instanceIds []string
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		instanceIds = append(instanceIds, item["InstanceId"].(string))
	}

	for _, instanceId := range instanceIds {

		action := "OnsGroupList"
		request := make(map[string]interface{})
		var response map[string]interface{}
		conn, err := client.NewOnsClient()
		if err != nil {
			return WrapError(err)
		}
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-02-14"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ons_groups", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Data.SubscribeInfoDo", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Data.SubscribeInfoDo", response)
		}

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			name := item["GroupId"].(string)
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
					skip = false
					break
				}
			}
			if skip {
				log.Printf("[INFO] Skipping ons group: %s ", name)
				continue
			}
			log.Printf("[INFO] delete ons group: %s ", name)

			action := "OnsGroupDelete"
			request := map[string]interface{}{
				"GroupId":    name,
				"InstanceId": instanceId,
			}
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-02-14"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				log.Printf("[ERROR] Failed to delete ons group (%s): %s", name, err)

			}
		}
	}

	return nil
}

func TestAccAlicloudOnsGroup_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ons_group.default"
	ra := resourceAttrInit(resourceId, onsGroupBasicMap)
	serviceFunc := func() interface{} {
		return &OnsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandIntRange(1000000, 9999999)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("GID-tf-testacconsgroupbasic%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceOnsGroupConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id": "${alicloud_ons_instance.default.id}",
					"group_id":    "${var.group_id}",
					"remark":      "alicloud_ons_group_remark",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_id": fmt.Sprintf("GID-tf-testacconsgroupbasic%v", rand),
						"remark":   "alicloud_ons_group_remark",
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
					"tags": map[string]string{
						"Created": "TFM",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "1",
						"tags.Created": "TFM",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"group_id": "${var.group_id}_change",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_id": fmt.Sprintf("GID-tf-testacconsgroupbasic%v_change", rand)}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"group_id": "${var.group_id}",
					"remark":   "alicloud_ons_group_remark",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_id": fmt.Sprintf("GID-tf-testacconsgroupbasic%v", rand),
						"remark":   "alicloud_ons_group_remark",
					}),
				),
			},
		},
	})

}

func resourceOnsGroupConfigDependence(name string) string {
	return fmt.Sprintf(`
resource "alicloud_ons_instance" "default" {
  name = "%s"
}

variable "group_id" {
 default = "%s"
}
`, name, name)
}

var onsGroupBasicMap = map[string]string{
	"group_id": "${var.group_id}",
	"remark":   "alicloud_ons_group_remark",
}
