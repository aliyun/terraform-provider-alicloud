package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/alikafka"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func init() {
	resource.AddTestSweepers("alicloud_alikafka_instance", &resource.Sweeper{
		Name: "alicloud_alikafka_instance",
		F:    testSweepAlikafkaInstance,
	})
}

func testSweepAlikafkaInstance(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)
	alikafkaService := AlikafkaService{client}

	prefixes := []string{
		"tf-testAcc",
		"tf_testacc",
	}

	instanceListReq := alikafka.CreateGetInstanceListRequest()
	instanceListReq.RegionId = region

	raw, err := alikafkaService.client.WithAlikafkaClient(func(alikafkaClient *alikafka.Client) (interface{}, error) {
		return alikafkaClient.GetInstanceList(instanceListReq)
	})
	if err != nil {
		log.Printf("[ERROR] Failed to retrieve alikafka instance in service list: %s", err)
	}

	instanceListResp, _ := raw.(*alikafka.GetInstanceListResponse)

	for _, v := range instanceListResp.InstanceList.InstanceVO {

		name := v.Name
		skip := true
		for _, prefix := range prefixes {

			// ServiceStatus equals 5 means the instance is in running status.
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) && v.ServiceStatus == 5 {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping alikafka instance: %s ", name)
			continue
		}
		log.Printf("[INFO] delete alikafka instance: %s ", name)

		request := alikafka.CreateReleaseInstanceRequest()
		request.InstanceId = v.InstanceId

		_, err := alikafkaService.client.WithAlikafkaClient(func(alikafkaClient *alikafka.Client) (interface{}, error) {
			return alikafkaClient.ReleaseInstance(request)
		})

		if err != nil {
			log.Printf("[ERROR] Failed to delete alikafka instance (%s): %s", name, err)
		}
	}

	return nil
}

func TestAccAlicloudAlikafkaInstance_basic(t *testing.T) {

	var v *alikafka.InstanceVO
	resourceId := "alicloud_alikafka_instance.default"
	ra := resourceAttrInit(resourceId, alikafkaInstanceBasicMap)
	serviceFunc := func() interface{} {
		return &AlikafkaService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandInt()
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testacc-alikafkainstancebasic%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAlikafkaInstanceConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, connectivity.AlikafkaSupportedRegions)
			testAccPreCheck(t)
			testAccPreCheckWithNoDefaultVswitch(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckAlikafkaInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":        var.name,
					"topic_quota": "50",
					"disk_type":   "1",
					"disk_size":   "500",
					"deploy_type": "5",
					"io_max":      "20",
					"vswitch_id":  data.alicloud_vswitches.default.ids.0,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": fmt.Sprintf("tf-testacc-alikafkainstancebasic%v", rand),
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
					"name": fmt.Sprintf("tf-testacc-alikafkainstancechange%v", rand),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": fmt.Sprintf("tf-testacc-alikafkainstancechange%v", rand)}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"topic_quota": "51",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"topic_quota": "51"}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"disk_type": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disk_type": "0"}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"disk_size": "600",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disk_size": "600"}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"deploy_type": "4",
					"eip_max":     "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deploy_type": "4",
						"eip_max":     "1"}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"io_max": "30",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"io_max": "30",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "acceptance test",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test",
						"Updated": "TF",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "3",
						"tags.Created": "TF",
						"tags.For":     "acceptance test",
						"tags.Updated": "TF",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"name":        var.name,
					"topic_quota": "50",
					"disk_type":   "1",
					"disk_size":   "500",
					"deploy_type": "5",
					"io_max":      "20",
					"eip_max":     "0",
					"paid_type":   "PrePaid",
					"spec_type":   "professional",
					"tags":        REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":         fmt.Sprintf("tf-testacc-alikafkainstancebasic%v", rand),
						"topic_quota":  "50",
						"disk_type":    "1",
						"disk_size":    "500",
						"deploy_type":  "5",
						"io_max":       "20",
						"eip_max":      "0",
						"paid_type":    "PrePaid",
						"spec_type":    "professional",
						"tags.%":       REMOVEKEY,
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
						"tags.Updated": REMOVEKEY,
					}),
				),
			},
		},
	})
}

func TestAccAlicloudAlikafkaInstance_multi(t *testing.T) {

	var v *alikafka.InstanceVO
	resourceId := "alicloud_alikafka_instance.default.1"
	ra := resourceAttrInit(resourceId, alikafkaInstanceBasicMap)
	serviceFunc := func() interface{} {
		return &AlikafkaService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandInt()
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testacc-alikafkainstancebasic%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAlikafkaInstanceConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, connectivity.AlikafkaSupportedRegions)
			testAccPreCheck(t)
			testAccPreCheckWithNoDefaultVswitch(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"count":       "2",
					"name":        "${var.name}-${count.index}",
					"topic_quota": "50",
					"disk_type":   "1",
					"disk_size":   "500",
					"deploy_type": "5",
					"io_max":      "20",
					"vswitch_id":  data.alicloud_vswitches.default.ids.0,
					"paid_type":   "PostPaid",
					"spec_type":   "normal",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})

}

func testAccCheckAlikafkaInstanceDestroy(s *terraform.State) error {
	// We expect instance to still exist
	return nil
}

func resourceAlikafkaInstanceConfigDependence(name string) string {
	return fmt.Sprintf(`
		variable "name" {
 			default = "%v"
		}

		data "alicloud_vswitches" "default" {
		  is_default = true
		}
		`, name)
}

var alikafkaInstanceBasicMap = map[string]string{
	"topic_quota": CHECKSET,
	"disk_type":   CHECKSET,
	"disk_size":   CHECKSET,
	"deploy_type": CHECKSET,
	"io_max":      CHECKSET,
	"vswitch_id":  CHECKSET,
	"paid_type":   CHECKSET,
	"spec_type":   CHECKSET,
}
