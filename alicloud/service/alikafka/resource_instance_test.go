// collection-order: bypass selected_zones -- pre-existing attribute outside this change
// collection-order: bypass vswitch_ids -- pre-existing attribute outside this change

package alikafka_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/alikafka"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudAlikafkaInstance_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alikafka_instance.default"
	ra := acctest.ResourceAttrInit(resourceId, alikafkaInstanceBasicMap)
	serviceFunc := func() interface{} {
		return alikafka.NewAlikafkaService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeAliKafkaInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	rand := tfacctest.RandInt()
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testacc-alikafkainstancebasic%v", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceAlikafkaInstanceConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"password": "YourPassword123!",
					"zone_id":  "${data.alicloud_vswitches.default.vswitches.0.zone_id}",
					"confluent_config": []map[string]interface{}{
						{
							"connect_cu":               "1",
							"connect_replica":          "1",
							"control_center_cu":        "1",
							"control_center_replica":   "1",
							"control_center_storage":   "1",
							"kafka_cu":                 "1",
							"kafka_replica":            "1",
							"kafka_rest_proxy_cu":      "1",
							"kafka_rest_proxy_replica": "1",
							"kafka_storage":            "1",
							"ksql_cu":                  "1",
							"ksql_replica":             "1",
							"ksql_storage":             "1",
							"schema_registry_cu":       "1",
							"schema_registry_replica":  "1",
							"zookeeper_cu":             "1",
							"zookeeper_replica":        "1",
							"zookeeper_storage":        "1",
						},
					},
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
						"security_group": acctest.CHECKSET,
						"kms_key_id":     acctest.CHECKSET,
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
					"password":        "YourPassword123!-upd",
					"service_version": "2.2.0",
					"confluent_config": []map[string]interface{}{
						{
							"connect_cu":               "2",
							"connect_replica":          "2",
							"control_center_cu":        "2",
							"control_center_storage":   "2",
							"kafka_cu":                 "2",
							"kafka_replica":            "2",
							"kafka_rest_proxy_cu":      "2",
							"kafka_rest_proxy_replica": "2",
							"kafka_storage":            "2",
							"ksql_cu":                  "2",
							"ksql_replica":             "2",
							"ksql_storage":             "2",
							"schema_registry_cu":       "2",
							"schema_registry_replica":  "2",
							"zookeeper_cu":             "2",
							"zookeeper_storage":        "2",
						},
					},
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
					"spec_type": "professional",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"spec_type": "professional",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"config": `{\"enable.vpc_sasl_ssl\":\"false\",\"kafka.log.retention.hours\":\"96\",\"enable.acl\":\"false\",\"kafka.message.max.bytes\":\"1048576\"}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"config": acctest.CHECKSET,
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
					"tags": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       acctest.REMOVEKEY,
						"tags.Created": acctest.REMOVEKEY,
						"tags.For":     acctest.REMOVEKEY,
					}),
				),
			},
		},
	})
}

func TestAccAliCloudAlikafkaInstance_convert(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alikafka_instance.default"
	ra := acctest.ResourceAttrInit(resourceId, alikafkaInstanceBasicMap)
	serviceFunc := func() interface{} {
		return alikafka.NewAlikafkaService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeAliKafkaInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	rand := tfacctest.RandInt()
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testacc-alikafkainstanceconvert%v", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceAlikafkaInstancePrePaidConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPrePaidResources(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      nil,
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
						"resource_group_id": acctest.CHECKSET,
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
	ra := acctest.ResourceAttrInit(resourceId, alikafkaInstanceBasicMap)
	serviceFunc := func() interface{} {
		return alikafka.NewAlikafkaService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeAliKafkaInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	rand := tfacctest.RandInt()
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testacc-alikafkainstancepre%v", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceAlikafkaInstancePrePaidConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPrePaidResources(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      nil,
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
	ra := acctest.ResourceAttrInit(resourceId, alikafkaInstanceBasicMap)
	serviceFunc := func() interface{} {
		return alikafka.NewAlikafkaService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeAliKafkaInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	rand := tfacctest.RandInt()
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testacc-alikafkainstancepre%v", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceAlikafkaInstanceConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
						"config":                      acctest.CHECKSET,
						"vswitch_ids.#":               "2",
						"tags.%":                      "2",
						"tags.Created":                "TF",
						"tags.For":                    "acceptance test",
						"ssl_endpoint":                acctest.CHECKSET,
						"ssl_domain_endpoint":         acctest.CHECKSET,
						"sasl_domain_endpoint":        acctest.CHECKSET,
						"vpc_sasl_domain_endpoint":    acctest.CHECKSET,
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
						"resource_group_id": acctest.CHECKSET,
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
	acctest.CheckoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	resourceId := "alicloud_alikafka_instance.default"
	ra := acctest.ResourceAttrInit(resourceId, alikafkaInstanceServerlessMap)
	serviceFunc := func() interface{} {
		return alikafka.NewAlikafkaService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeAliKafkaInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	rand := tfacctest.RandInt()
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testacc-alikafkainstancepre%v", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceAlikafkaInstanceConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
						"deploy_type":              "4",
						"instance_type":            "alikafka_serverless",
						"paid_type":                "PostPaid",
						"spec_type":                "normal",
						"service_version":          "3.3.1",
						"config":                   acctest.CHECKSET,
						"vswitch_ids.#":            "2",
						"serverless_config.#":      "1",
						"tags.%":                   "2",
						"tags.Created":             "TF",
						"tags.For":                 "acceptance test",
						"ssl_endpoint":             acctest.CHECKSET,
						"ssl_domain_endpoint":      acctest.CHECKSET,
						"sasl_domain_endpoint":     acctest.CHECKSET,
						"vpc_sasl_domain_endpoint": acctest.CHECKSET,
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
				Config: testAccConfig(map[string]interface{}{
					"deploy_type": "5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deploy_type":         "5",
						"ssl_endpoint":        acctest.REMOVEKEY,
						"ssl_domain_endpoint": acctest.REMOVEKEY,
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
	acctest.CheckoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	resourceId := "alicloud_alikafka_instance.default"
	ra := acctest.ResourceAttrInit(resourceId, alikafkaInstanceServerlessMap)
	serviceFunc := func() interface{} {
		return alikafka.NewAlikafkaService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeAliKafkaInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	rand := tfacctest.RandInt()
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testacc-alikafkainstancepre%v", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceAlikafkaInstancePrePaidConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      nil,
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
	"topic_quota":      acctest.CHECKSET,
	"partition_num":    acctest.CHECKSET,
	"disk_type":        acctest.CHECKSET,
	"disk_size":        acctest.CHECKSET,
	"deploy_type":      acctest.CHECKSET,
	"instance_type":    acctest.CHECKSET,
	"io_max":           acctest.CHECKSET,
	"vswitch_id":       acctest.CHECKSET,
	"paid_type":        acctest.CHECKSET,
	"spec_type":        acctest.CHECKSET,
	"vpc_id":           acctest.CHECKSET,
	"zone_id":          acctest.CHECKSET,
	"service_version":  acctest.CHECKSET,
	"config":           acctest.CHECKSET,
	"end_point":        acctest.CHECKSET,
	"domain_endpoint":  acctest.CHECKSET,
	"topic_num_of_buy": acctest.CHECKSET,
	"topic_used":       acctest.CHECKSET,
	"topic_left":       acctest.CHECKSET,
	"partition_used":   acctest.CHECKSET,
	"partition_left":   acctest.CHECKSET,
	"group_used":       acctest.CHECKSET,
	"group_left":       acctest.CHECKSET,
	"is_partition_buy": acctest.CHECKSET,
	"status":           acctest.CHECKSET,
}

var alikafkaInstanceServerlessMap = map[string]string{
	"deploy_type":   acctest.CHECKSET,
	"instance_type": acctest.CHECKSET,
	"io_max":        acctest.CHECKSET,
	"vswitch_id":    acctest.CHECKSET,
	"paid_type":     acctest.CHECKSET,
	"spec_type":     acctest.CHECKSET,
	"vpc_id":        acctest.CHECKSET,
	"zone_id":       acctest.CHECKSET,
	"status":        acctest.CHECKSET,
}
