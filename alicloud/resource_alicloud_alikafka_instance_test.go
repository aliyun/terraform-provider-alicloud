package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alikafka"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_alikafka_instance", &resource.Sweeper{
		Name: "alicloud_alikafka_instance",
		F:    testSweepAlikafkaInstance,
		// When implemented, these should be removed firstly
		Dependencies: []string{
			"alicloud_alikafka_consumer_group",
			"alicloud_alikafka_sasl_acl",
			"alicloud_alikafka_topic",
			"alicloud_alikafka_sasl_user",
		},
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
	service := VpcService{client}
	for _, v := range instanceListResp.InstanceList.InstanceVO {

		name := v.Name
		skip := true
		for _, prefix := range prefixes {

			// ServiceStatus equals 5 means the instance is in running status.
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		// If a ES description is not set successfully, it should be fetched by vswitch name and deleted.
		if skip {
			if need, err := service.needSweepVpc(v.VpcId, v.VSwitchId); err == nil {
				skip = !need
			}
		}
		if skip {
			log.Printf("[INFO] Skipping alikafka instance: %s ", name)
			continue
		}
		if v.ServiceStatus != 10 {
			log.Printf("[INFO] release alikafka instance: %s ", name)

			request := alikafka.CreateReleaseInstanceRequest()
			request.InstanceId = v.InstanceId
			request.ForceDeleteInstance = requests.NewBoolean(true)
			_, err := alikafkaService.client.WithAlikafkaClient(func(alikafkaClient *alikafka.Client) (interface{}, error) {
				return alikafkaClient.ReleaseInstance(request)
			})

			if err != nil {
				log.Printf("[ERROR] Failed to release alikafka instance (%s): %s", name, err)
			}
		}

		log.Printf("[INFO] Delete alikafka instance: %s ", name)
		request2 := alikafka.CreateDeleteInstanceRequest()
		request2.InstanceId = v.InstanceId
		_, err = alikafkaService.client.WithAlikafkaClient(func(alikafkaClient *alikafka.Client) (interface{}, error) {
			return alikafkaClient.DeleteInstance(request2)
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
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rc.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":           "${var.name}",
					"topic_quota":    "50",
					"disk_type":      "1",
					"disk_size":      "500",
					"deploy_type":    "5",
					"io_max":         "20",
					"vswitch_id":     "${data.alicloud_vswitches.default.ids.0}",
					"security_group": "${alicloud_security_group.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":           fmt.Sprintf("tf-testacc-alikafkainstancebasic%v", rand),
						"security_group": CHECKSET,
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
						"topic_quota": "51",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"disk_type": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disk_type": "0",
					}),
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
					"eip_max":     "3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deploy_type": "4",
						"eip_max":     "3",
					}),
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
					"spec_type":       "professional",
					"service_version": "2.2.0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"spec_type":       "professional",
						"service_version": "2.2.0",
					}),
				),
			},
			//// Current, there is an api bug that the service will return an extra key: kafka.offsets.retention.minutes in the config
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"config": `{\"enable.vpc_sasl_ssl\":\"false\",\"kafka.log.retention.hours\":\"96\",\"enable.acl\":\"false\",\"kafka.message.max.bytes\":\"1048576\"}`,
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"config": "{\"enable.vpc_sasl_ssl\":\"false\",\"kafka.log.retention.hours\":\"96\",\"enable.acl\":\"false\",\"kafka.message.max.bytes\":\"1048576\"}",
			//		}),
			//	),
			//},

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

			// suspend PrePaid testing
			{
				Config: testAccConfig(map[string]interface{}{
					"name":            "${var.name}",
					"topic_quota":     "50",
					"disk_type":       "1",
					"disk_size":       "500",
					"deploy_type":     "5",
					"io_max":          "20",
					"eip_max":         "0",
					"spec_type":       "professional",
					"service_version": "2.2.0",
					//"config":          `{\"enable.vpc_sasl_ssl\":\"false\",\"kafka.log.retention.hours\":\"96\",\"enable.acl\":\"false\",\"kafka.message.max.bytes\":\"1048576\"}`,
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":            fmt.Sprintf("tf-testacc-alikafkainstancebasic%v", rand),
						"topic_quota":     "50",
						"disk_type":       "1",
						"disk_size":       "500",
						"deploy_type":     "5",
						"io_max":          "20",
						"eip_max":         "0",
						"paid_type":       "PostPaid",
						"spec_type":       "professional",
						"service_version": "2.2.0",
						//"config":          "{\"enable.vpc_sasl_ssl\":\"false\",\"kafka.log.retention.hours\":\"96\",\"enable.acl\":\"false\",\"kafka.message.max.bytes\":\"1048576\"}",
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

func TestAccAlicloudAlikafkaInstance_convert(t *testing.T) {
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
	name := fmt.Sprintf("tf-testacc-alikafkainstanceconvert%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAlikafkaInstanceConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, connectivity.AlikafkaSupportedRegions)
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		//CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":            "${var.name}",
					"topic_quota":     "50",
					"disk_type":       "1",
					"disk_size":       "500",
					"deploy_type":     "4",
					"eip_max":         "3",
					"io_max":          "20",
					"vswitch_id":      "${data.alicloud_vswitches.default.ids.0}",
					"paid_type":       "PostPaid",
					"spec_type":       "normal",
					"service_version": "0.10.2",
					//"config":          `{\"enable.vpc_sasl_ssl\":\"false\",\"kafka.log.retention.hours\":\"72\",\"enable.acl\":\"false\",\"kafka.message.max.bytes\":\"1048576\"}`,
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test",
					},
					"security_group": "${data.alicloud_security_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":            name,
						"topic_quota":     "50",
						"disk_type":       "1",
						"disk_size":       "500",
						"deploy_type":     "4",
						"eip_max":         "3",
						"io_max":          "20",
						"paid_type":       "PostPaid",
						"spec_type":       "normal",
						"service_version": "0.10.2",
						//"config":          "{\"enable.vpc_sasl_ssl\":\"false\",\"kafka.log.retention.hours\":\"72\",\"enable.acl\":\"false\",\"kafka.message.max.bytes\":\"1048576\"}",
						"tags.%":         "2",
						"tags.Created":   "TF",
						"tags.For":       "acceptance test",
						"security_group": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"paid_type": "PrePaid",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"paid_type": "PrePaid",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAlicloudAlikafkaInstance_prepaid(t *testing.T) {
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
	name := fmt.Sprintf("tf-testacc-alikafkainstancepre%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAlikafkaInstanceConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, connectivity.AlikafkaSupportedRegions)
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		//CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":            "${var.name}",
					"topic_quota":     "50",
					"disk_type":       "1",
					"disk_size":       "500",
					"deploy_type":     "4",
					"eip_max":         "3",
					"io_max":          "20",
					"vswitch_id":      "${data.alicloud_vswitches.default.ids.0}",
					"paid_type":       "PrePaid",
					"spec_type":       "normal",
					"service_version": "0.10.2",
					//"config":          `{\"enable.vpc_sasl_ssl\":\"false\",\"kafka.log.retention.hours\":\"72\",\"enable.acl\":\"false\",\"kafka.message.max.bytes\":\"1048576\"}`,
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test",
					},
					"security_group": "${data.alicloud_security_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":            name,
						"topic_quota":     "50",
						"disk_type":       "1",
						"disk_size":       "500",
						"deploy_type":     "4",
						"eip_max":         "3",
						"io_max":          "20",
						"paid_type":       "PrePaid",
						"spec_type":       "normal",
						"service_version": "0.10.2",
						//"config":          "{\"enable.vpc_sasl_ssl\":\"false\",\"kafka.log.retention.hours\":\"72\",\"enable.acl\":\"false\",\"kafka.message.max.bytes\":\"1048576\"}",
						"tags.%":         "2",
						"tags.Created":   "TF",
						"tags.For":       "acceptance test",
						"security_group": CHECKSET,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func resourceAlikafkaInstanceConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%v"
}

data "alicloud_vpcs" "default" {
 name_regex = "^default-NODELETING"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_security_group" "default" {
  name   = var.name
  vpc_id = data.alicloud_vpcs.default.ids.0
}
data "alicloud_security_groups" "default" {
	name_regex = "^default-NODELETING"
}
`, name)
}

var alikafkaInstanceBasicMap = map[string]string{
	"topic_quota":     CHECKSET,
	"disk_type":       CHECKSET,
	"disk_size":       CHECKSET,
	"deploy_type":     CHECKSET,
	"io_max":          CHECKSET,
	"vswitch_id":      CHECKSET,
	"paid_type":       CHECKSET,
	"spec_type":       CHECKSET,
	"vpc_id":          CHECKSET,
	"zone_id":         CHECKSET,
	"end_point":       CHECKSET,
	"service_version": CHECKSET,
	"config":          CHECKSET,
}
