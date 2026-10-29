// Package mns_test.
package mns_test

import (
	"fmt"
	"log"
	"strings"
	"testing"

	ali_mns "github.com/aliyun/aliyun-mns-go-sdk"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/mns"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_mns_topic", &resource.Sweeper{
		Name: "alicloud_mns_topic",
		F:    testSweepMnsTopics,
	})
}

func testSweepMnsTopics(region string) error {
	rawClient, err := sweep.SharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf-testacc",
	}

	var topicAttrs []ali_mns.TopicAttribute
	for _, namePrefix := range prefixes {
		for {
			var nextMaker string
			raw, err := client.WithMnsTopicManager(func(topicManager ali_mns.AliTopicManager) (interface{}, error) {
				return topicManager.ListTopicDetail(nextMaker, 1000, namePrefix)
			})
			if err != nil {
				return fmt.Errorf("get topicDetails  error: %#v", err)
			}
			topicDetails, _ := raw.(ali_mns.TopicDetails)
			for _, attr := range topicDetails.Attrs {
				topicAttrs = append(topicAttrs, attr)
			}
			nextMaker = topicDetails.NextMarker
			if nextMaker == "" {
				break
			}
		}
	}
	for _, topicAttr := range topicAttrs {
		name := topicAttr.TopicName
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping mns topic : %s ", name)
			continue
		}
		log.Printf("[INFO] delete  mns topic : %s ", name)
		_, err := client.WithMnsTopicManager(func(topicManager ali_mns.AliTopicManager) (interface{}, error) {
			return nil, topicManager.DeleteTopic(name)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete mns topic (%s (%s)): %s", topicAttr.TopicName, topicAttr.TopicName, err)
		}
	}

	return nil
}

func TestAccAliCloudMnsTopic_basic(t *testing.T) {
	var v *ali_mns.TopicAttribute
	resourceId := "alicloud_mns_topic.default"
	ra := acctest.ResourceAttrInit(resourceId, mnsTopicMap)
	serviceFunc := func() interface{} {
		return mns.NewMnsService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccMNSTopicConfig-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceMnsTopicConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name,
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
					"maximum_message_size": "12357",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"maximum_message_size": "12357",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"logging_enabled": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"logging_enabled": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"maximum_message_size": "65536",
					"logging_enabled":      "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"maximum_message_size": "65536",
						"logging_enabled":      "false",
					}),
				),
			},
		},
	})
}

func TestAccAliCloudMnsTopic_multi(t *testing.T) {
	var v *ali_mns.TopicAttribute
	resourceId := "alicloud_mns_topic.default.4"
	ra := acctest.ResourceAttrInit(resourceId, mnsTopicMap)
	serviceFunc := func() interface{} {
		return mns.NewMnsService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccMNSTopicConfig-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceMnsTopicConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":  name + "${count.index}",
					"count": "5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func resourceMnsTopicConfigDependence(name string) string {
	return ""
}

var mnsTopicMap = map[string]string{
	"name":                 acctest.CHECKSET,
	"maximum_message_size": "65536",
	"logging_enabled":      "false",
}
