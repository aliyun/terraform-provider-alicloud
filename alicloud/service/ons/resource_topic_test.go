package ons_test

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	ons "github.com/aliyun/terraform-provider-alicloud/alicloud/service/ons"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_ons_topic", &resource.Sweeper{
		Name: "alicloud_ons_topic",
		F:    testSweepOnsTopic,
	})
}

func testSweepOnsTopic(region string) error {
	rawClient, err := sweep.SharedClientForRegion(region)
	if err != nil {
		return errs.WrapErrorf(err, "error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf-testacc",
	}

	action := "OnsInstanceInServiceList"
	request := make(map[string]interface{})
	var response map[string]interface{}
	response, err = client.RpcPost("Ons", "2019-02-14", action, nil, request, false)
	if err != nil {
		log.Printf("[ERROR] Failed to retrieve ons instance in service list: %s", err)
	}
	resp, err := jsonpath.Get("$.Data.InstanceVO", response)
	if err != nil {
		return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.Data.InstanceVO", response)
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
		response, err = client.RpcPost("Ons", "2019-02-14", action, nil, request, false)
		if err != nil {
			log.Printf("[ERROR] Failed to retrieve ons instance in service list: %s", err)
		}
		resp, err := jsonpath.Get("$.Data.PublishInfoDo", response)
		if err != nil {
			return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.Data.PublishInfoDo", response)
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
			response, err = client.RpcPost("Ons", "2019-02-14", action, nil, request, false)
			if err != nil {
				log.Printf("[ERROR] Failed to delete ons topic (%s): %s", name, err)
			}
		}
	}

	return nil
}

func TestAccAliCloudOnsTopic_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ons_topic.default"
	ra := resourceAttrInit(resourceId, onsTopicBasicMap)
	serviceFunc := func() interface{} {
		return ons.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := tfacctest.RandIntRange(1000000, 9999999)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testacc%sonstopicbasic%v", acctest.DefaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceOnsTopicConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		// module name
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"topic_name":   name,
					"instance_id":  "${alicloud_ons_instance.default.id}",
					"message_type": "1",
					"remark":       "alicloud_ons_topic_remark",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"topic":        fmt.Sprintf("tf-testacc%sonstopicbasic%v", acctest.DefaultRegionToTest, rand),
						"remark":       "alicloud_ons_topic_remark",
						"message_type": "1",
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
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
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

resource "alicloud_ons_topic" "mutilple" {
  count        = 10
  topic_name   = format("%%s%%s", var.topic, count.index)
  instance_id  = alicloud_ons_instance.default.id
  message_type = 0
  remark       = "dafault_ons_topic_remark"
}
`, name, name)
}

var onsTopicBasicMap = map[string]string{
	"message_type": "1",
	"remark":       "alicloud_ons_topic_remark",
}
