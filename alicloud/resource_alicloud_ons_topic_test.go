package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ons"

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
	onsService := OnsService{client}

	prefixes := []string{
		"tf-testAcc",
		"tf-testacc",
	}

	instanceListReq := ons.CreateOnsInstanceInServiceListRequest()

	raw, err := onsService.client.WithOnsClient(func(onsClient *ons.Client) (interface{}, error) {
		return onsClient.OnsInstanceInServiceList(instanceListReq)
	})
	if err != nil {
		log.Printf("[ERROR] Failed to retrieve ons instance in service list: %s", err)
	}

	instanceListResp, _ := raw.(*ons.OnsInstanceInServiceListResponse)

	var instanceIds []string
	for _, v := range instanceListResp.Data.InstanceVO {
		instanceIds = append(instanceIds, v.InstanceId)
	}

	for _, instanceId := range instanceIds {
		request := ons.CreateOnsTopicListRequest()
		request.InstanceId = instanceId

		raw, err := onsService.client.WithOnsClient(func(onsClient *ons.Client) (interface{}, error) {
			return onsClient.OnsTopicList(request)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to retrieve ons topics on instance (%s): %s", instanceId, err)
			continue
		}

		topicListResp, _ := raw.(*ons.OnsTopicListResponse)
		topics := topicListResp.Data.PublishInfoDo

		for _, v := range topics {
			name := v.Topic
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

			request := ons.CreateOnsTopicDeleteRequest()
			request.InstanceId = instanceId
			request.Topic = v.Topic

			_, err := onsService.client.WithOnsClient(func(onsClient *ons.Client) (interface{}, error) {
				return onsClient.OnsTopicDelete(request)
			})
			if err != nil {
				log.Printf("[ERROR] Failed to delete ons topic (%s): %s", name, err)
			}
		}
	}

	return nil
}

func TestAccAlicloudOnsTopic_basic(t *testing.T) {
	var v ons.PublishInfoDo
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

			{
				Config: testAccConfig(map[string]interface{}{
					"message_type": "5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"message_type": "5"}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"perm": "4",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"perm": "4"}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"topic":  "tf-testacc-alicloud_ons_default_topic_change",
					"remark": "default remark",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"topic":  "tf-testacc-alicloud_ons_default_topic_change",
						"remark": "default remark"}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"topic":        "${var.topic}",
					"message_type": "0",
					"remark":       "alicloud_ons_topic_remark",
					"perm":         "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"topic":        fmt.Sprintf("tf-testacc%sonstopicbasic%v", defaultRegionToTest, rand),
						"message_type": "0",
						"remark":       "alicloud_ons_topic_remark",
						"perm":         "2",
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
	"topic":        "${var.topic}",
	"message_type": "1",
	"remark":       "alicloud_ons_topic_remark",
	"perm":         "6",
}
