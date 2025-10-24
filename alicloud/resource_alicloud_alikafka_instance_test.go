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
		return WrapErrorf(err, "error getting AliCloud client.")
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

func TestAccAliCloudAlikafkaInstance_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alikafka_instance.default"
	ra := resourceAttrInit(resourceId, alikafkaInstanceBasicMap)
	serviceFunc := func() interface{} {
		return &AlikafkaService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeAliKafkaInstance")
	rac := resourceAttrCheckInit(rc, ra)
	rand := acctest.RandInt()
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testacc-alikafkainstancebasic%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAlikafkaInstanceConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rc.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":           name,
					"topic_quota":    "50",
					"disk_type":      "1",
					"disk_size":      "500",
					"deploy_type":    "5",
					"io_max":         "20",
					"vswitch_id":     "${data.alicloud_vswitches.default.ids.0}",
					"security_group": "${alicloud_security_group.default.id}",
					"kms_key_id":     "${alicloud_kms_key.key.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":           name,
						"security_group": CHECKSET,
						"kms_key_id":     CHECKSET,
						"partition_num":  "0",
						"topic_quota":    "1000",
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
					"partition_num": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"partition_num": "1",
						"topic_quota":   "1001",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"disk_size": "800",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disk_size": "800",
					}),
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
			{
				Config: testAccConfig(map[string]interface{}{
					"config": `{\"enable.vpc_sasl_ssl\":\"false\",\"kafka.log.retention.hours\":\"96\",\"enable.acl\":\"false\",\"kafka.message.max.bytes\":\"1048576\"}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"config": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_auto_group": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_auto_group": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_auto_topic": "enable",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_auto_topic": "enable",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"default_topic_partition_num": "6",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"default_topic_partition_num": "6",
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
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       REMOVEKEY,
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}

func TestAccAliCloudAlikafkaInstance_convert(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alikafka_instance.default"
	ra := resourceAttrInit(resourceId, alikafkaInstanceBasicMap)
	serviceFunc := func() interface{} {
		return &AlikafkaService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeAliKafkaInstance")
	rac := resourceAttrCheckInit(rc, ra)
	rand := acctest.RandInt()
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testacc-alikafkainstanceconvert%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAlikafkaInstancePrePaidConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":              name,
					"partition_num":     "50",
					"disk_type":         "1",
					"disk_size":         "500",
					"deploy_type":       "4",
					"instance_type":     "alikafka",
					"eip_max":           "3",
					"io_max":            "20",
					"vswitch_id":        "${data.alicloud_vswitches.default.ids.0}",
					"paid_type":         "PostPaid",
					"spec_type":         "normal",
					"service_version":   "2.2.0",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
					//"config":          `{\"enable.vpc_sasl_ssl\":\"false\",\"kafka.log.retention.hours\":\"72\",\"enable.acl\":\"false\",\"kafka.message.max.bytes\":\"1048576\"}`,
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":              name,
						"partition_num":     "50",
						"topic_quota":       "1050",
						"disk_type":         "1",
						"disk_size":         "500",
						"deploy_type":       "4",
						"instance_type":     "alikafka",
						"eip_max":           "3",
						"io_max":            "20",
						"paid_type":         "PostPaid",
						"spec_type":         "normal",
						"service_version":   "2.2.0",
						"resource_group_id": CHECKSET,
						//"config":          "{\"enable.vpc_sasl_ssl\":\"false\",\"kafka.log.retention.hours\":\"72\",\"enable.acl\":\"false\",\"kafka.message.max.bytes\":\"1048576\"}",
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "acceptance test",
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

func TestAccAliCloudAlikafkaInstance_prepaid(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alikafka_instance.default"
	ra := resourceAttrInit(resourceId, alikafkaInstanceBasicMap)
	serviceFunc := func() interface{} {
		return &AlikafkaService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeAliKafkaInstance")
	rac := resourceAttrCheckInit(rc, ra)
	rand := acctest.RandInt()
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testacc-alikafkainstancepre%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAlikafkaInstancePrePaidConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":            name,
					"partition_num":   "50",
					"disk_type":       "1",
					"disk_size":       "500",
					"deploy_type":     "4",
					"instance_type":   "alikafka",
					"eip_max":         "3",
					"io_max":          "20",
					"vpc_id":          "${data.alicloud_vpcs.default.ids.0}",
					"vswitch_id":      "${data.alicloud_vswitches.default.ids.0}",
					"paid_type":       "PrePaid",
					"spec_type":       "normal",
					"service_version": "2.2.0",
					//"config":          `{\"enable.vpc_sasl_ssl\":\"false\",\"kafka.log.retention.hours\":\"72\",\"enable.acl\":\"false\",\"kafka.message.max.bytes\":\"1048576\"}`,
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":            name,
						"partition_num":   "50",
						"topic_quota":     "1050",
						"disk_type":       "1",
						"disk_size":       "500",
						"deploy_type":     "4",
						"instance_type":   "alikafka",
						"eip_max":         "3",
						"io_max":          "20",
						"paid_type":       "PrePaid",
						"spec_type":       "normal",
						"service_version": "2.2.0",
						//"config":          "{\"enable.vpc_sasl_ssl\":\"false\",\"kafka.log.retention.hours\":\"72\",\"enable.acl\":\"false\",\"kafka.message.max.bytes\":\"1048576\"}",
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "acceptance test",
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

func TestAccAliCloudAlikafkaInstance_VpcId(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alikafka_instance.default"
	ra := resourceAttrInit(resourceId, alikafkaInstanceBasicMap)
	serviceFunc := func() interface{} {
		return &AlikafkaService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeAliKafkaInstance")
	rac := resourceAttrCheckInit(rc, ra)
	rand := acctest.RandInt()
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testacc-alikafkainstancepre%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAlikafkaInstanceConfigDependence)
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
					"name":                        name,
					"topic_quota":                 "50",
					"disk_type":                   "1",
					"disk_size":                   "800",
					"deploy_type":                 "4",
					"eip_max":                     "3",
					"io_max_spec":                 "alikafka.hw.2xlarge",
					"vswitch_id":                  "${data.alicloud_vswitches.default.ids.0}",
					"paid_type":                   "PostPaid",
					"spec_type":                   "professional",
					"service_version":             "2.2.0",
					"enable_auto_group":           "true",
					"enable_auto_topic":           "enable",
					"default_topic_partition_num": "6",
					"config":                      `{\"enable.vpc_sasl_ssl\":\"true\",\"kafka.log.retention.hours\":\"72\",\"enable.acl\":\"true\",\"kafka.message.max.bytes\":\"1048576\"}`,
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test",
					},
					"security_group": "${alicloud_security_group.default.id}",
					"vpc_id":         "${data.alicloud_vpcs.default.ids.0}",
					"vswitch_ids":    []string{"${data.alicloud_vswitches.default.ids.0}", "${data.alicloud_vswitches.default.ids.1}"},
					"selected_zones": []string{"zoneb", "zonec"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                        name,
						"partition_num":               "0",
						"topic_quota":                 "1000",
						"disk_type":                   "1",
						"disk_size":                   "800",
						"deploy_type":                 "4",
						"eip_max":                     "3",
						"io_max_spec":                 "alikafka.hw.2xlarge",
						"paid_type":                   "PostPaid",
						"spec_type":                   "professional",
						"service_version":             "2.2.0",
						"enable_auto_group":           "true",
						"enable_auto_topic":           "enable",
						"default_topic_partition_num": "6",
						"config":                      CHECKSET,
						"vswitch_ids.#":               "2",
						"tags.%":                      "2",
						"tags.Created":                "TF",
						"tags.For":                    "acceptance test",
						"ssl_endpoint":                CHECKSET,
						"ssl_domain_endpoint":         CHECKSET,
						"sasl_domain_endpoint":        CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"io_max_spec": "alikafka.hw.3xlarge",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"io_max_spec": "alikafka.hw.3xlarge",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"partition_num": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"partition_num": "1",
						"topic_quota":   "1001",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"selected_zones"},
			},
		},
	})
}

func TestAccAliCloudAlikafkaInstance_Serverless(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	resourceId := "alicloud_alikafka_instance.default"
	ra := resourceAttrInit(resourceId, alikafkaInstanceServerlessMap)
	serviceFunc := func() interface{} {
		return &AlikafkaService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeAliKafkaInstance")
	rac := resourceAttrCheckInit(rc, ra)
	rand := acctest.RandInt()
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testacc-alikafkainstancepre%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAlikafkaInstanceConfigDependence)
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
					"deploy_type":     "4",
					"instance_type":   "alikafka_serverless",
					"vswitch_id":      "${data.alicloud_vswitches.default.ids.0}",
					"paid_type":       "PostPaid",
					"spec_type":       "normal",
					"service_version": "3.3.1",
					"config":          `{\"enable.vpc_sasl_ssl\":\"true\",\"kafka.log.retention.hours\":\"72\",\"enable.acl\":\"true\",\"kafka.message.max.bytes\":\"1048576\"}`,
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test",
					},
					"security_group": "${alicloud_security_group.default.id}",
					"vswitch_ids":    []string{"${data.alicloud_vswitches.default.ids.0}", "${data.alicloud_vswitches.default.ids.1}"},
					"selected_zones": []string{"zoneb", "zonec"},
					"serverless_config": []map[string]interface{}{
						{
							"reserved_publish_capacity":   "60",
							"reserved_subscribe_capacity": "60",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deploy_type":         "4",
						"instance_type":       "alikafka_serverless",
						"paid_type":           "PostPaid",
						"spec_type":           "normal",
						"service_version":     "3.3.1",
						"config":              CHECKSET,
						"vswitch_ids.#":       "2",
						"serverless_config.#": "1",
						"tags.%":              "2",
						"tags.Created":        "TF",
						"tags.For":            "acceptance test",
						"ssl_endpoint":        CHECKSET,
						"ssl_domain_endpoint": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"serverless_config": []map[string]interface{}{
						{
							"reserved_publish_capacity":   "100",
							"reserved_subscribe_capacity": "100",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"serverless_config.#": "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"selected_zones"},
			},
		},
	})
}

func SkipTestAccAliCloudAlikafkaInstance_Confluent(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	resourceId := "alicloud_alikafka_instance.default"
	ra := resourceAttrInit(resourceId, alikafkaInstanceServerlessMap)
	serviceFunc := func() interface{} {
		return &AlikafkaService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeAliKafkaInstance")
	rac := resourceAttrCheckInit(rc, ra)
	rand := acctest.RandInt()
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testacc-alikafkainstancepre%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAlikafkaInstancePrePaidConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"deploy_type":     "5",
					"instance_type":   "alikafka_confluent",
					"paid_type":       "PrePaid",
					"spec_type":       "professional",
					"service_version": "7.4.0",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test",
					},
					"password":       "YourPassword123!",
					"vswitch_ids":    []string{"${data.alicloud_vswitches.default.ids.0}", "${data.alicloud_vswitches.default.ids.1}"},
					"selected_zones": []string{"zoneb", "zonec"},
					"confluent_config": []map[string]interface{}{
						{
							"kafka_cu":                 "4",
							"kafka_storage":            "800",
							"kafka_replica":            "3",
							"kafka_rest_proxy_cu":      "4",
							"kafka_rest_proxy_replica": "2",
							"zookeeper_cu":             "2",
							"zookeeper_storage":        "100",
							"zookeeper_replica":        "3",
							"control_center_cu":        "4",
							"control_center_storage":   "300",
							"control_center_replica":   "1",
							"schema_registry_cu":       "2",
							"schema_registry_replica":  "2",
							"connect_cu":               "4",
							"connect_replica":          "2",
							"ksql_cu":                  "4",
							"ksql_storage":             "100",
							"ksql_replica":             "2",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deploy_type":        "5",
						"instance_type":      "alikafka_confluent",
						"paid_type":          "PrePaid",
						"spec_type":          "professional",
						"service_version":    "7.4.0",
						"vswitch_ids.#":      "2",
						"confluent_config.#": "1",
						"tags.%":             "2",
						"tags.Created":       "TF",
						"tags.For":           "acceptance test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"confluent_config": []map[string]interface{}{
						{
							"kafka_cu":                 "8",
							"kafka_storage":            "2000",
							"kafka_replica":            "5",
							"kafka_rest_proxy_cu":      "5",
							"kafka_rest_proxy_replica": "4",
							"zookeeper_cu":             "5",
							"zookeeper_storage":        "300",
							"zookeeper_replica":        "3",
							"control_center_cu":        "5",
							"control_center_storage":   "400",
							"control_center_replica":   "1",
							"schema_registry_cu":       "5",
							"schema_registry_replica":  "3",
							"connect_cu":               "5",
							"connect_replica":          "4",
							"ksql_cu":                  "5",
							"ksql_storage":             "200",
							"ksql_replica":             "4",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"confluent_config.#": "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "selected_zones"},
			},
		},
	})
}

func resourceAlikafkaInstanceConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}

	data "alicloud_resource_manager_resource_groups" "default" {
	}

	data "alicloud_vpcs" "default" {
  		name_regex = "^default-NODELETING$"
	}

	data "alicloud_vswitches" "default" {
  		vpc_id = data.alicloud_vpcs.default.ids.0
	}

	resource "alicloud_security_group" "default" {
		security_group_name = var.name
		vpc_id              = data.alicloud_vpcs.default.ids.0
	}

	resource "alicloud_kms_key" "key" {
  		description            = var.name
  		pending_window_in_days = "7"
  		status                 = "Enabled"
	}
`, name)
}

func resourceAlikafkaInstancePrePaidConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}

	data "alicloud_resource_manager_resource_groups" "default" {
	}

	data "alicloud_vpcs" "default" {
  		name_regex = "^default-NODELETING$"
	}

	data "alicloud_vswitches" "default" {
  		vpc_id = data.alicloud_vpcs.default.ids.0
	}
`, name)
}

var alikafkaInstanceBasicMap = map[string]string{
	"topic_quota":      CHECKSET,
	"partition_num":    CHECKSET,
	"disk_type":        CHECKSET,
	"disk_size":        CHECKSET,
	"deploy_type":      CHECKSET,
	"instance_type":    CHECKSET,
	"io_max":           CHECKSET,
	"vswitch_id":       CHECKSET,
	"paid_type":        CHECKSET,
	"spec_type":        CHECKSET,
	"vpc_id":           CHECKSET,
	"zone_id":          CHECKSET,
	"service_version":  CHECKSET,
	"config":           CHECKSET,
	"end_point":        CHECKSET,
	"domain_endpoint":  CHECKSET,
	"topic_num_of_buy": CHECKSET,
	"topic_used":       CHECKSET,
	"topic_left":       CHECKSET,
	"partition_used":   CHECKSET,
	"partition_left":   CHECKSET,
	"group_used":       CHECKSET,
	"group_left":       CHECKSET,
	"is_partition_buy": CHECKSET,
	"status":           CHECKSET,
}

var alikafkaInstanceServerlessMap = map[string]string{
	"deploy_type":   CHECKSET,
	"instance_type": CHECKSET,
	"io_max":        CHECKSET,
	"vswitch_id":    CHECKSET,
	"paid_type":     CHECKSET,
	"spec_type":     CHECKSET,
	"vpc_id":        CHECKSET,
	"zone_id":       CHECKSET,
	"status":        CHECKSET,
}
