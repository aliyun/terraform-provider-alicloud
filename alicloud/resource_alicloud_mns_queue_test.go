package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	ali_mns "github.com/aliyun/aliyun-mns-go-sdk"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func init() {
	resource.AddTestSweepers("alicloud_mns_queue", &resource.Sweeper{
		Name: "alicloud_mns_queue",
		F:    testSweepMnsQueues,
	})
}

func testSweepMnsQueues(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf-testacc",
	}

	var queueAttrs []ali_mns.QueueAttribute
	for _, namePrefix := range prefixes {
		for {
			var nextMaker string
			raw, err := client.WithMnsQueueManager(func(queueManager ali_mns.AliQueueManager) (interface{}, error) {
				return queueManager.ListQueueDetail(nextMaker, 1000, namePrefix)
			})
			if err != nil {
				return fmt.Errorf("get queueDetails  error: %#v", err)
			}
			queueDetails, _ := raw.(ali_mns.QueueDetails)
			for _, attr := range queueDetails.Attrs {
				queueAttrs = append(queueAttrs, attr)
			}
			nextMaker = queueDetails.NextMarker
			if nextMaker == "" {
				break
			}
		}
	}
	for _, queueAttr := range queueAttrs {
		name := queueAttr.QueueName
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping mns queque: %s ", name)
			continue
		}
		log.Printf("[INFO] delete  mns queque: %s ", name)
		_, err := client.WithMnsQueueManager(func(queueManager ali_mns.AliQueueManager) (interface{}, error) {
			return nil, queueManager.DeleteQueue(queueAttr.QueueName)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete mnsQueue (%s (%s)): %s", queueAttr.QueueName, queueAttr.QueueName, err)
		}
	}

	return nil
}

func TestAccAlicloudMnsQueue_basic(t *testing.T) {
	var v ali_mns.QueueAttribute
	resourceId := "alicloud_mns_queue.default"
	ra := resourceAttrInit(resourceId, mnsQueueMap)
	serviceFunc := func() interface{} {
		return &MnsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccMNSQueueConfig-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceMnsQueueConfigDependence)

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
					"delay_seconds": "60478",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delay_seconds": "60478",
					}),
				),
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
					"message_retention_period": "256000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"message_retention_period": "256000",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"visibility_timeout": "40",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"visibility_timeout": "40",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"polling_wait_seconds": "3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"polling_wait_seconds": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"delay_seconds":            "0",
					"maximum_message_size":     "65536",
					"message_retention_period": "345600",
					"visibility_timeout":       "30",
					"polling_wait_seconds":     "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delay_seconds":            "0",
						"maximum_message_size":     "65536",
						"message_retention_period": "345600",
						"visibility_timeout":       "30",
						"polling_wait_seconds":     "0",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudMnsQueue_multi(t *testing.T) {
	var v ali_mns.QueueAttribute
	resourceId := "alicloud_mns_queue.default.4"
	ra := resourceAttrInit(resourceId, mnsQueueMap)
	serviceFunc := func() interface{} {
		return &MnsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccMNSQueueConfig-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceMnsQueueConfigDependence)

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

func resourceMnsQueueConfigDependence(name string) string {
	return ""
}

var mnsQueueMap = map[string]string{
	"name":                     CHECKSET,
	"delay_seconds":            "0",
	"maximum_message_size":     "65536",
	"message_retention_period": "345600",
	"visibility_timeout":       "30",
	"polling_wait_seconds":     "0",
}
