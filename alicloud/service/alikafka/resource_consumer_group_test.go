package alikafka_test

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	alikafkasdk "github.com/aliyun/alibaba-cloud-sdk-go/services/alikafka"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/alikafka"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_alikafka_consumer_group", &resource.Sweeper{
		Name: "alicloud_alikafka_consumer_group",
		F:    testSweepAlikafkaConsumerGroup,
	})
}

func testSweepAlikafkaConsumerGroup(region string) error {
	rawClient, err := sweep.SharedClientForRegion(region)
	if err != nil {
		return errs.WrapErrorf(err, "error getting AliCloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testacc",
	}

	instanceListReq := alikafkasdk.CreateGetInstanceListRequest()
	instanceListReq.RegionId = acctest.DefaultRegionToTest

	raw, err := client.WithAlikafkaClient(func(alikafkaClient *alikafkasdk.Client) (interface{}, error) {
		return alikafkaClient.GetInstanceList(instanceListReq)
	})
	if err != nil {
		log.Printf("[ERROR] Failed to retrieve alikafka instance in service list: %s", err)
	}

	instanceListResp, _ := raw.(*alikafkasdk.GetInstanceListResponse)

	var instanceIds []string
	for _, v := range instanceListResp.InstanceList.InstanceVO {
		instanceIds = append(instanceIds, v.InstanceId)
	}

	for _, instanceId := range instanceIds {

		// Control the consumer group list request rate.
		time.Sleep(time.Duration(400) * time.Millisecond)

		request := alikafkasdk.CreateGetConsumerListRequest()
		request.InstanceId = instanceId
		request.RegionId = acctest.DefaultRegionToTest

		raw, err := client.WithAlikafkaClient(func(alikafkaClient *alikafkasdk.Client) (interface{}, error) {
			return alikafkaClient.GetConsumerList(request)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to retrieve alikafka consumer groups on instance (%s): %s", instanceId, err)
			continue
		}

		consumerListResp, _ := raw.(*alikafkasdk.GetConsumerListResponse)
		consumers := consumerListResp.ConsumerList.ConsumerVO
		for _, v := range consumers {
			name := v.ConsumerId
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
					skip = false
					break
				}
			}
			if skip {
				log.Printf("[INFO] Skipping alikafka consumer id: %s ", name)
				continue
			}
			log.Printf("[INFO] delete alikafka consumer group: %s ", name)

			// Control the consumer group delete rate
			time.Sleep(time.Duration(400) * time.Millisecond)

			request := alikafkasdk.CreateDeleteConsumerGroupRequest()
			request.InstanceId = instanceId
			request.ConsumerId = v.ConsumerId
			request.RegionId = acctest.DefaultRegionToTest

			_, err := client.WithAlikafkaClient(func(alikafkaClient *alikafkasdk.Client) (interface{}, error) {
				return alikafkaClient.DeleteConsumerGroup(request)
			})
			if err != nil {
				log.Printf("[ERROR] Failed to delete alikafka consumer group (%s): %s", name, err)
			}
		}
	}

	return nil
}

// Test AliKafka ConsumerGroup. >>> Resource test cases, automatically generated.
// Case group全生命周期 9851
func TestAccAliCloudAliKafkaConsumerGroup_basic9851(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alikafka_consumer_group.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudAliKafkaConsumerGroupMap9851)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return alikafka.NewAlikafkaServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeAliKafkaConsumerGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccalikafka%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudAliKafkaConsumerGroupBasicDependence9851)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id": "${data.alicloud_alikafka_instances.default.instances.0.id}",
					"consumer_id": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id": acctest.CHECKSET,
						"consumer_id": acctest.CHECKSET,
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
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": acctest.REMOVEKEY,
						"tags.For":     acctest.REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func TestAccAliCloudAliKafkaConsumerGroup_basic9851_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alikafka_consumer_group.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudAliKafkaConsumerGroupMap9851)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return alikafka.NewAlikafkaServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeAliKafkaConsumerGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccalikafka%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudAliKafkaConsumerGroupBasicDependence9851)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id": "${data.alicloud_alikafka_instances.default.instances.0.id}",
					"consumer_id": name,
					"remark":      name,
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":  acctest.CHECKSET,
						"consumer_id":  acctest.CHECKSET,
						"remark":       name,
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

// Case group全生命周期, 适配废弃字段description 9852
func TestAccAliCloudAliKafkaConsumerGroup_basic9852_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alikafka_consumer_group.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudAliKafkaConsumerGroupMap9851)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return alikafka.NewAlikafkaServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeAliKafkaConsumerGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccalikafka%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudAliKafkaConsumerGroupBasicDependence9851)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id": "${data.alicloud_alikafka_instances.default.instances.0.id}",
					"consumer_id": name,
					"description": name,
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":  acctest.CHECKSET,
						"consumer_id":  acctest.CHECKSET,
						"description":  name,
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func TestAccAliCloudAliKafkaConsumerGroup_multi(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alikafka_consumer_group.default.5"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudAliKafkaConsumerGroupMap9851)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return alikafka.NewAlikafkaServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeAliKafkaConsumerGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccalikafka%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudAliKafkaConsumerGroupBasicDependence9851)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"count":       "6",
					"instance_id": "${data.alicloud_alikafka_instances.default.instances.0.id}",
					"consumer_id": name + "-${count.index}",
					"remark":      name + "-${count.index}",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":  acctest.CHECKSET,
						"consumer_id":  acctest.CHECKSET,
						"remark":       name + fmt.Sprint(-5),
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
		},
	})
}

var AliCloudAliKafkaConsumerGroupMap9851 = map[string]string{
	"create_time": acctest.CHECKSET,
	"region_id":   acctest.CHECKSET,
}

func AliCloudAliKafkaConsumerGroupBasicDependence9851(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_alikafka_instances" "default" {
	}
`, name)
}

// Test AliKafka ConsumerGroup. <<< Resource test cases, automatically generated.
