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
	resource.AddTestSweepers("alicloud_ons_topic", &resource.Sweeper{
		Name: "alicloud_ons_topic",
		F:    testSweepOnsTopic,
	})
}

func testSweepOnsTopic(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf-testacc",
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
		action := "OnsTopicList"
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
		resp, err := jsonpath.Get("$.Data.PublishInfoDo", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Data.PublishInfoDo", response)
		}

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			name := item["Topic"].(string)
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
					skip = false
					break
				}
			}
			if skip {
				log.Printf("[INFO] Skipping ons topic: %s ", name)
				continue
			}
			log.Printf("[INFO] delete ons topic: %s ", name)

			action := "OnsTopicDelete"
			request := map[string]interface{}{
				"InstanceId": instanceId,
				"Topic":      item["Topic"],
			}
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-02-14"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				log.Printf("[ERROR] Failed to delete ons topic (%s): %s", name, err)
			}
		}
	}

	return nil
}

func TestAccAlicloudOnsTopic_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ons_topic.default"
	ra := resourceAttrInit(resourceId, onsTopicBasicMap)
	serviceFunc := func() interface{} {
		return &OnsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandIntRange(1000000, 9999999)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testacc%sonstopicbasic%v", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceOnsTopicConfigDependence)

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
					"instance_id":  "${alicloud_ons_instance.default.id}",
					"topic":        "${var.topic}",
					"message_type": "1",
					"remark":       "alicloud_ons_topic_remark",
					"perm":         "6",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"topic":  fmt.Sprintf("tf-testacc%sonstopicbasic%v", defaultRegionToTest, rand),
						"remark": "alicloud_ons_topic_remark",
						"perm":   "6",
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
						"Created": "TF",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "1",
						"tags.Created": "TF",
					}),
				),
			},
			// TODO: there is an openapi bug that OnsTopicUpdate updating perm does not work
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"perm": "4",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{"perm": "4"}),
			//	),
			//},

			{
				Config: testAccConfig(map[string]interface{}{
					//"perm": "2",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						//"perm":         "2",
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
		},
	})

}

func resourceOnsTopicConfigDependence(name string) string {
	return fmt.Sprintf(`
resource "alicloud_ons_instance" "default" {
  name = "%s"
}

variable "topic" {
 default = "%s"
}
`, name, name)
}

var onsTopicBasicMap = map[string]string{
	"message_type": "1",
	"remark":       "alicloud_ons_topic_remark",
	"perm":         "6",
}
