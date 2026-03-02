package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/elasticsearch"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const DataNodeSpec = "elasticsearch.sn1ne.large"
const DataNodeSpecForUpdate = "elasticsearch.sn2ne.large"

const DataNodeAmount = "2"
const DataNodeAmountForUpdate = "3"

const DataNodeDisk = "20"
const DataNodeDiskForUpdate = "30"
const DataNodeDiskForEssdUpdate = "461"

const DataNodeSsdDiskType = "cloud_ssd"
const DataNodeEssdDiskType = "cloud_essd"

const DataNodeDiskPerformanceLevel1 = "PL1"
const DataNodeDiskPerformanceLevel2 = "PL2"

const DataNodeAmountForMultiZone = "4"
const DefaultZoneAmount = "2"

const MasterNodeSpec = "elasticsearch.sn1ne.large"
const MasterNodeSpecForUpdate = "elasticsearch.sn2ne.large"

const MasterNodeDiskType = "cloud_ssd"
const MasterNodeEssdDiskType = "cloud_essd"

const ClientNodeSpec = "elasticsearch.sn1ne.large"
const ClientNodeSpecForUpdate = "elasticsearch.sn2ne.large"

const ClientNodeAmount = "2"
const ClientNodeAmountForUpdate = "3"

const KibanaSpec = "elasticsearch.sn1ne.large"
const KibanaSpecForUpdate = "elasticsearch.sn2ne.large"

const WarmNodeSpec = "elasticsearch.sn1ne.large"
const WarmNodeSpecUpdate = "elasticsearch.sn1ne.xlarge"
const WarmNodeAmount = "3"
const WarmNodeAmountUpdate = "4"
const WarmNodeDiskSize = "500"
const WarmNodeDiskSizeUpdate = "500"
const WarmNodeDiskType = "cloud_efficiency"

const AutoRenewal = "AutoRenewal"
const NotRenewal = "NotRenewal"
const ManualRenewal = "ManualRenewal"

const Version55 = "5.5.3_with_X-Pack"
const Version716 = "7.16.2_with_X-Pack"

const Version77 = "7.7.1_with_X-Pack"

func init() {
	resource.AddTestSweepers("alicloud_elasticsearch_instance", &resource.Sweeper{
		Name: "alicloud_elasticsearch_instance",
		F:    testSweepElasticsearch,
	})
}

func testSweepElasticsearch(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("Error getting Alicloud client: %s", err)
	}

	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"",
		"tf-testAcc",
		"tf_testAcc",
	}

	var instances []elasticsearch.Instance
	req := elasticsearch.CreateListInstanceRequest()
	req.RegionId = client.RegionId
	req.Page = requests.NewInteger(1)
	req.Size = requests.NewInteger(PageSizeLarge)

	for {
		raw, err := client.WithElasticsearchClient(func(elasticsearchClient *elasticsearch.Client) (interface{}, error) {
			return elasticsearchClient.ListInstance(req)
		})

		if err != nil {
			log.Printf("[ERROR] %s", WrapError(fmt.Errorf("Error listing Elasticsearch instances: %s", err)))
			break
		}

		resp, _ := raw.(*elasticsearch.ListInstanceResponse)
		if resp == nil || len(resp.Result) < 1 {
			break
		}

		instances = append(instances, resp.Result...)

		if len(resp.Result) < PageSizeLarge {
			break
		}

		page, err := getNextpageNumber(req.Page)
		if err != nil {
			return err
		}
		req.Page = page
	}

	sweeped := false
	service := VpcService{client}
	for _, v := range instances {
		description := v.Description
		id := v.InstanceId
		skip := true
		if !sweepAll() {
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(description), strings.ToLower(prefix)) {
					skip = false
					break
				}
			}
			// If a ES description is not set successfully, it should be fetched by vswitch name and deleted.
			if skip {
				if need, err := service.needSweepVpc(v.NetworkConfig.VpcId, v.NetworkConfig.VswitchId); err == nil {
					skip = !need
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Elasticsearch Instance: %s (%s)", description, id)
				continue
			}
		}
		log.Printf("[INFO] Deleting Elasticsearch Instance: %s (%s)", description, id)
		req := elasticsearch.CreateDeleteInstanceRequest()
		req.InstanceId = id
		_, err := client.WithElasticsearchClient(func(elasticsearchClient *elasticsearch.Client) (interface{}, error) {
			return elasticsearchClient.DeleteInstance(req)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Elasticsearch Instance (%s (%s)): %s", description, id, err)
		} else {
			sweeped = true
		}
	}

	if sweeped {
		// Waiting 30 seconds to eusure these instances have been deleted.
		time.Sleep(30 * time.Second)
	}

	return nil
}

func TestAccAliCloudElasticsearchInstance_basic(t *testing.T) {
	var instance *elasticsearch.DescribeInstanceResponse

	resourceId := "alicloud_elasticsearch_instance.default"
	ra := resourceAttrInit(resourceId, elasticsearchMap)

	serviceFunc := func() interface{} {
		return &ElasticsearchService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &instance, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandInt()
	name := fmt.Sprintf("tf-testAccES%s%d", defaultRegionToTest, rand)
	if len(name) > 30 {
		name = name[:30]
	}
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceElasticsearchInstanceConfigDependence)

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
					"description":                      name,
					"vswitch_id":                       "${local.vswitch_id}",
					"version":                          Version716,
					"password":                         "Yourpassword1234",
					"data_node_spec":                   DataNodeSpec,
					"data_node_amount":                 DataNodeAmount,
					"data_node_disk_size":              DataNodeDisk,
					"data_node_disk_type":              DataNodeEssdDiskType,
					"data_node_disk_performance_level": DataNodeDiskPerformanceLevel1,
					"kibana_node_spec":                 KibanaSpec,
					"instance_charge_type":             string(PostPaid),
					"zone_count":                       "1",
					"resource_group_id":                "${local.resource_group_id}",
					"enable_kibana_public_network":     "true",
					"enable_kibana_private_network":    "true",
					"kibana_private_security_group_id": "${local.security_group}",
					"warm_node_amount":                 WarmNodeAmount,
					"warm_node_disk_size":              WarmNodeDiskSize,
					"warm_node_disk_encrypted":         "false",
					"warm_node_spec":                   WarmNodeSpec,
					"warm_node_disk_type":              WarmNodeDiskType,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":                      name,
						"version":                          Version716,
						"data_node_spec":                   DataNodeSpec,
						"data_node_amount":                 DataNodeAmount,
						"data_node_disk_size":              DataNodeDisk,
						"data_node_disk_type":              DataNodeEssdDiskType,
						"data_node_disk_performance_level": DataNodeDiskPerformanceLevel1,
						"kibana_node_spec":                 KibanaSpec,
						"instance_charge_type":             string(PostPaid),
						"zone_count":                       "1",
						"resource_group_id":                CHECKSET,
						"enable_kibana_public_network":     "true",
						"enable_kibana_private_network":    "true",
						"kibana_private_security_group_id": CHECKSET,
						"warm_node_amount":                 WarmNodeAmount,
						"warm_node_disk_size":              WarmNodeDiskSize,
						"warm_node_disk_encrypted":         "false",
						"warm_node_spec":                   WarmNodeSpec,
						"warm_node_disk_type":              WarmNodeDiskType,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "acceptance Test",
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

			{
				Config: testAccConfig(map[string]interface{}{
					"warm_node_amount":         WarmNodeAmountUpdate,
					"warm_node_disk_size":      WarmNodeDiskSizeUpdate,
					"warm_node_disk_encrypted": "false",
					"warm_node_spec":           WarmNodeSpecUpdate,
					"warm_node_disk_type":      WarmNodeDiskType,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"warm_node_amount":         WarmNodeAmountUpdate,
						"warm_node_disk_size":      WarmNodeDiskSizeUpdate,
						"warm_node_disk_encrypted": "false",
						"warm_node_spec":           WarmNodeSpecUpdate,
						"warm_node_disk_type":      WarmNodeDiskType,
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"enable_kibana_private_network": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_kibana_private_network": "false",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"enable_kibana_private_network":    "true",
					"kibana_private_security_group_id": "${local.security_group}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_kibana_private_network":    "true",
						"kibana_private_security_group_id": CHECKSET,
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"password": "Yourpassword1235",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"password": "Yourpassword1235",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name[:len(name)-1],
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name[:len(name)-1],
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"data_node_spec":                   DataNodeSpecForUpdate,
					"data_node_amount":                 DataNodeAmountForUpdate,
					"data_node_disk_type":              DataNodeEssdDiskType,
					"data_node_disk_size":              DataNodeDiskForEssdUpdate,
					"data_node_disk_performance_level": DataNodeDiskPerformanceLevel2,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_node_spec":                   DataNodeSpecForUpdate,
						"data_node_amount":                 DataNodeAmountForUpdate,
						"data_node_disk_type":              DataNodeEssdDiskType,
						"data_node_disk_size":              DataNodeDiskForEssdUpdate,
						"data_node_disk_performance_level": DataNodeDiskPerformanceLevel2,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"kibana_node_spec": KibanaSpecForUpdate,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"kibana_node_spec": KibanaSpecForUpdate,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"master_node_spec":      MasterNodeSpec,
					"master_node_disk_type": MasterNodeEssdDiskType,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"master_node_spec":      MasterNodeSpec,
						"master_node_disk_type": MasterNodeEssdDiskType,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"client_node_spec":   ClientNodeSpec,
					"client_node_amount": ClientNodeAmount,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"client_node_spec":   ClientNodeSpec,
						"client_node_amount": ClientNodeAmount,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"protocol": "HTTP",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocol": "HTTP",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"protocol": "HTTPS",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocol": "HTTPS",
					}),
				),
			},
		},
	})
}

func TestAccAliCloudElasticsearchInstance_version(t *testing.T) {
	var instance *elasticsearch.DescribeInstanceResponse

	resourceId := "alicloud_elasticsearch_instance.default"
	ra := resourceAttrInit(resourceId, elasticsearchMap)

	serviceFunc := func() interface{} {
		return &ElasticsearchService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &instance, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandInt()
	name := fmt.Sprintf("tf-testAccES%s%d", defaultRegionToTest, rand)
	if len(name) > 30 {
		name = name[:30]
	}
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceElasticsearchInstanceConfigDependence)

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
					"description":                      name,
					"vswitch_id":                       "${local.vswitch_id}",
					"version":                          "7.16_with_X-Pack",
					"password":                         "Yourpassword1234",
					"data_node_spec":                   DataNodeSpec,
					"data_node_amount":                 DataNodeAmount,
					"data_node_disk_size":              DataNodeDisk,
					"data_node_disk_type":              DataNodeEssdDiskType,
					"data_node_disk_performance_level": DataNodeDiskPerformanceLevel1,
					"kibana_node_spec":                 KibanaSpec,

					"instance_charge_type": string(PostPaid),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
						"version":     REGEXMATCH + "^7.16.*_with_X-Pack",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"version": "8.5_with_X-Pack",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"version": REGEXMATCH + "8.5.*_with_X-Pack",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"warm_node_amount":         WarmNodeAmount,
					"warm_node_disk_size":      WarmNodeDiskSize,
					"warm_node_disk_encrypted": "false",
					"warm_node_spec":           WarmNodeSpec,
					"warm_node_disk_type":      WarmNodeDiskType,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"warm_node_amount":         WarmNodeAmount,
						"warm_node_disk_size":      WarmNodeDiskSize,
						"warm_node_disk_encrypted": "false",
						"warm_node_spec":           WarmNodeSpec,
						"warm_node_disk_type":      WarmNodeDiskType,
					}),
				),
			},
		},
	})
}

func TestAccAliCloudElasticsearchInstance_multizone(t *testing.T) {
	var instance *elasticsearch.DescribeInstanceResponse

	resourceId := "alicloud_elasticsearch_instance.default"
	ra := resourceAttrInit(resourceId, elasticsearchMap)

	serviceFunc := func() interface{} {
		return &ElasticsearchService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &instance, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandInt()
	name := fmt.Sprintf("tf-testAccES%s%d", defaultRegionToTest, rand)
	if len(name) > 30 {
		name = name[:30]
	}
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceElasticsearchInstanceConfigDependenceKms)

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
					"description":            name,
					"vswitch_id":             "${local.vswitch_id}",
					"version":                Version716,
					"kms_encrypted_password": "${alicloud_kms_ciphertext.default.ciphertext_blob}",
					"kms_encryption_context": map[string]string{
						"name": name,
					},
					"data_node_spec":                   DataNodeSpec,
					"data_node_amount":                 DataNodeAmountForMultiZone,
					"data_node_disk_size":              DataNodeDisk,
					"data_node_disk_type":              DataNodeEssdDiskType,
					"data_node_disk_performance_level": DataNodeDiskPerformanceLevel1,
					"master_node_spec":                 MasterNodeSpec,
					"master_node_disk_type":            MasterNodeEssdDiskType,
					"kibana_node_spec":                 KibanaSpec,
					"instance_charge_type":             string(PostPaid),
					"zone_count":                       DefaultZoneAmount,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":                      name,
						"version":                          Version716,
						"data_node_spec":                   DataNodeSpec,
						"data_node_amount":                 DataNodeAmountForMultiZone,
						"data_node_disk_size":              DataNodeDisk,
						"data_node_disk_type":              DataNodeEssdDiskType,
						"data_node_disk_performance_level": DataNodeDiskPerformanceLevel1,
						"master_node_spec":                 MasterNodeSpec,
						"master_node_disk_type":            MasterNodeEssdDiskType,
						"kibana_node_spec":                 KibanaSpec,
						"instance_charge_type":             string(PostPaid),
						"zone_count":                       DefaultZoneAmount,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"kms_encrypted_password": "${alicloud_kms_ciphertext.update.ciphertext_blob}",
					"kms_encryption_context": map[string]string{
						"name": name + "update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"kms_encrypted_password":   CHECKSET,
						"kms_encryption_context.%": "1",
					}),
				),
			},
		},
	})
}

func TestAccAliCloudElasticsearchInstance_encrypt_disk(t *testing.T) {
	var instance *elasticsearch.DescribeInstanceResponse

	resourceId := "alicloud_elasticsearch_instance.default"
	ra := resourceAttrInit(resourceId, elasticsearchMap)

	serviceFunc := func() interface{} {
		return &ElasticsearchService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &instance, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandInt()
	name := fmt.Sprintf("tf-testAccES%s%d", defaultRegionToTest, rand)
	if len(name) > 30 {
		name = name[:30]
	}
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceElasticsearchInstanceConfigDependence)

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
					"description":                      name,
					"vswitch_id":                       "${local.vswitch_id}",
					"version":                          Version716,
					"password":                         "Yourpassword1234",
					"data_node_spec":                   DataNodeSpec,
					"data_node_amount":                 DataNodeAmount,
					"data_node_disk_size":              DataNodeDisk,
					"data_node_disk_type":              DataNodeEssdDiskType,
					"data_node_disk_performance_level": DataNodeDiskPerformanceLevel1,
					"data_node_disk_encrypted":         "true",
					"kibana_node_spec":                 KibanaSpec,
					"instance_charge_type":             string(PostPaid),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":              name,
						"data_node_disk_encrypted": "true",
					}),
				),
			},
		},
	})
}

func TestAccAliCloudElasticsearchInstance_prepaid_autorenew(t *testing.T) {
	var instance *elasticsearch.DescribeInstanceResponse

	resourceId := "alicloud_elasticsearch_instance.default"
	ra := resourceAttrInit(resourceId, elasticsearchMap)

	serviceFunc := func() interface{} {
		return &ElasticsearchService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &instance, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandInt()
	name := fmt.Sprintf("tf-testAccES%s%d", defaultRegionToTest, rand)
	if len(name) > 30 {
		name = name[:30]
	}
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceElasticsearchInstanceConfigDependence)

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
					"description":                      name,
					"vswitch_id":                       "${local.vswitch_id}",
					"version":                          "7.16_with_X-Pack",
					"password":                         "Yourpassword1234",
					"data_node_spec":                   DataNodeSpec,
					"data_node_amount":                 DataNodeAmount,
					"data_node_disk_size":              DataNodeDisk,
					"data_node_disk_type":              DataNodeEssdDiskType,
					"data_node_disk_performance_level": DataNodeDiskPerformanceLevel1,
					"kibana_node_spec":                 KibanaSpec,
					"zone_count":                       "1",
					"instance_charge_type":             "PrePaid",
					"period":                           "1",
					"renew_status":                     AutoRenewal,
					"auto_renew_duration":              "1",
					"renewal_duration_unit":            "M",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_charge_type":             "PrePaid",
						"version":                          "7.16.2_with_X-Pack",
						"data_node_spec":                   DataNodeSpec,
						"data_node_amount":                 DataNodeAmount,
						"data_node_disk_size":              DataNodeDisk,
						"data_node_disk_type":              DataNodeEssdDiskType,
						"data_node_disk_performance_level": DataNodeDiskPerformanceLevel1,
						"kibana_node_spec":                 KibanaSpec,
						"private_whitelist.#":              "1",
						"zone_count":                       "1",
						"renew_status":                     AutoRenewal,
						"auto_renew_duration":              "1",
						"renewal_duration_unit":            "M",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "period"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_renew_duration":   "1",
					"renewal_duration_unit": "Y",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_renew_duration":   "1",
						"renewal_duration_unit": "Y",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_renew_duration":   "3",
					"renewal_duration_unit": "M",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_renew_duration":   "3",
						"renewal_duration_unit": "M",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"renew_status": ManualRenewal,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"renew_status": ManualRenewal,
					}),
				),
			},
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"renew_status": NotRenewal,
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"renew_status": NotRenewal,
			//		}),
			//	),
			//},
			// pre paid to post paid
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_charge_type": string(PostPaid),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_charge_type": string(PostPaid),
					}),
				),
			},
		},
	})
}

func TestAccAliCloudElasticsearchInstance_network(t *testing.T) {
	var instance *elasticsearch.DescribeInstanceResponse

	resourceId := "alicloud_elasticsearch_instance.default"
	ra := resourceAttrInit(resourceId, elasticsearchNetworkMap)

	serviceFunc := func() interface{} {
		return &ElasticsearchService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &instance, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandInt()
	name := fmt.Sprintf("tf-testAccES%s%d", defaultRegionToTest, rand)
	if len(name) > 30 {
		name = name[:30]
	}
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceElasticsearchInstanceConfigDependence)

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
					"description":                      name,
					"vswitch_id":                       "${local.vswitch_id}",
					"version":                          Version716,
					"password":                         "Yourpassword1234",
					"data_node_spec":                   DataNodeSpec,
					"data_node_amount":                 DataNodeAmount,
					"data_node_disk_size":              DataNodeDisk,
					"data_node_disk_type":              DataNodeEssdDiskType,
					"data_node_disk_performance_level": DataNodeDiskPerformanceLevel1,
					"kibana_node_spec":                 KibanaSpec,
					"instance_charge_type":             string(PostPaid),
					"zone_count":                       "1",
					"enable_public":                    "true",
					"enable_kibana_private_network":    "false",
					"enable_kibana_public_network":     "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":                      name,
						"version":                          Version716,
						"password":                         "Yourpassword1234",
						"data_node_spec":                   DataNodeSpec,
						"data_node_amount":                 DataNodeAmount,
						"data_node_disk_size":              DataNodeDisk,
						"data_node_disk_type":              DataNodeEssdDiskType,
						"data_node_disk_performance_level": DataNodeDiskPerformanceLevel1,
						"kibana_node_spec":                 KibanaSpec,
						"instance_charge_type":             string(PostPaid),
						"zone_count":                       "1",
						"enable_public":                    "true",
						"enable_kibana_private_network":    "false",
						"enable_kibana_public_network":     "false",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"enable_kibana_public_network": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_kibana_public_network": "true",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"enable_kibana_private_network":    "true",
					"kibana_private_security_group_id": "${local.security_group}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_kibana_private_network":    "true",
						"kibana_private_security_group_id": CHECKSET,
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"enable_kibana_public_network": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_kibana_public_network": "false",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"enable_kibana_private_network": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_kibana_private_network": "false",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"enable_kibana_public_network": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_kibana_public_network": "true",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"enable_kibana_private_network": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_kibana_private_network": "true",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"kibana_whitelist": []string{"192.0.0.1/32", "127.0.0.1/32"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"kibana_whitelist.#": "2",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"public_whitelist": []string{"192.168.0.0/24", "127.0.0.1/32"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"public_whitelist.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"private_whitelist": []string{"192.168.0.0/24", "127.0.0.1/32"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"private_whitelist.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"setting_config": map[string]string{
						"\"action.auto_create_index\"":         "+.*,-*",
						"\"action.destructive_requires_name\"": "false",
						"\"xpack.security.audit.enabled\"":     "true",
						"\"xpack.watcher.enabled\"":            "false",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"setting_config.action.auto_create_index":         "+.*,-*",
						"setting_config.action.destructive_requires_name": "false",
						"setting_config.xpack.security.audit.enabled":     "true",
						"setting_config.xpack.watcher.enabled":            "false",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"enable_kibana_private_network":    "true",
					"kibana_private_security_group_id": "${local.security_group}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_kibana_private_network":    "true",
						"kibana_private_security_group_id": CHECKSET,
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"enable_public": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_public": "false",
					}),
				),
			},
		},
	})
}

func TestAccAliCloudElasticsearchInstance_onecs(t *testing.T) {
	var instance *elasticsearch.DescribeInstanceResponse

	resourceId := "alicloud_elasticsearch_instance.default"
	ra := resourceAttrInit(resourceId, elasticsearchNetworkMap)

	serviceFunc := func() interface{} {
		return &ElasticsearchService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &instance, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandInt()
	name := fmt.Sprintf("tf-testAccES%s%d", defaultRegionToTest, rand)
	if len(name) > 30 {
		name = name[:30]
	}
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceElasticsearchInstanceConfigDependence)

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
					"description":                      name,
					"vswitch_id":                       "${local.vswitch_id}",
					"version":                          "5.5.3_with_X-Pack",
					"password":                         "Yourpassword1234",
					"data_node_spec":                   DataNodeSpec,
					"data_node_amount":                 DataNodeAmount,
					"data_node_disk_size":              DataNodeDisk,
					"data_node_disk_type":              DataNodeEssdDiskType,
					"data_node_disk_performance_level": DataNodeDiskPerformanceLevel1,
					"kibana_node_spec":                 KibanaSpec,
					"instance_charge_type":             string(PostPaid),
					"zone_count":                       "1",
					"enable_public":                    "false",
					"enable_kibana_private_network":    "true",
					"enable_kibana_public_network":     "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":                      name,
						"version":                          "5.5.3_with_X-Pack",
						"password":                         "Yourpassword1234",
						"data_node_spec":                   DataNodeSpec,
						"data_node_amount":                 DataNodeAmount,
						"data_node_disk_size":              DataNodeDisk,
						"data_node_disk_type":              DataNodeEssdDiskType,
						"data_node_disk_performance_level": DataNodeDiskPerformanceLevel1,
						"kibana_node_spec":                 KibanaSpec,
						"instance_charge_type":             string(PostPaid),
						"zone_count":                       "1",
						"enable_public":                    "false",
						"enable_kibana_private_network":    "true",
						"enable_kibana_public_network":     "true",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"kibana_private_whitelist": []string{"192.168.0.0/24", "127.0.0.1/32"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"kibana_private_whitelist.#": "2",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"enable_public": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_public": "true",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"enable_kibana_public_network": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_kibana_public_network": "false",
					}),
				),
			},
		},
	})
}

var elasticsearchMap = map[string]string{
	"description":                      CHECKSET,
	"data_node_spec":                   CHECKSET,
	"data_node_amount":                 CHECKSET,
	"data_node_disk_size":              CHECKSET,
	"data_node_disk_type":              CHECKSET,
	"data_node_disk_performance_level": CHECKSET,
	"instance_charge_type":             CHECKSET,
	"status":                           "active",
	"enable_public":                    "false",
	"enable_kibana_public_network":     "true",
	"enable_kibana_private_network":    "false",
	"id":                               CHECKSET,
	"vswitch_id":                       CHECKSET,
}

var elasticsearchNetworkMap = map[string]string{}

func resourceElasticsearchInstanceConfigDependence(name string) string {
	return fmt.Sprintf(`
    %s
	variable "name" {
		default = "%s"
	}
	`, ElasticsearchInstanceCommonTestCase, name)
}

func resourceElasticsearchInstanceConfigDependenceKms(name string) string {
	return fmt.Sprintf(`
    %s
	variable "name" {
		default = "%s"
	}

	resource "alicloud_kms_key" "default" {
  		description            = var.name
  		status                 = "Enabled"
  		pending_window_in_days = 7
	}

	resource "alicloud_kms_ciphertext" "default" {
  		key_id    = alicloud_kms_key.default.id
  		plaintext = "YourPassword1234!"
  		encryption_context = {
    		"name" = var.name
  		}
	}

	resource "alicloud_kms_ciphertext" "update" {
  		key_id    = alicloud_kms_key.default.id
  		plaintext = "YourPassword1234!update"
  		encryption_context = {
    		"name" = "${var.name}update"
  		}
	}
	
	`, ElasticsearchInstanceCommonTestCase, name)
}

// Test Elasticsearch Instance. >>> Resource test cases, automatically generated.
// Case 创建包含-描述-资源组-Protocol 11468
func TestAccAliCloudElasticsearchInstance_basic11468(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_elasticsearch_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudElasticsearchInstanceMap11468)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ElasticsearchServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeElasticsearchInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccelasticsearch%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudElasticsearchInstanceBasicDependence11468)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_count":  "1",
					"description": "创建包含-描述-资源组-Protocol",
					"kibana_configuration": []map[string]interface{}{
						{
							"amount": "1",
							"spec":   "elasticsearch.n4.small",
							"disk":   "0",
						},
					},
					"payment_type":             "PayAsYouGo",
					"vswitch_id":               "${alicloud_vswitch.defaultAislbL.id}",
					"instance_category":        "x-pack",
					"kibana_private_whitelist": []string{},
					"protocol":                 "HTTP",
					"version":                  "7.10.0_with_X-Pack",
					"password":                 "Admain@123",
					"resource_group_id":        "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"data_node_configuration": []map[string]interface{}{
						{
							"disk_type":         "cloud_essd",
							"spec":              "elasticsearch.sn1ne.large.new",
							"disk":              "20",
							"performance_level": "PL1",
							"disk_encryption":   "true",
							"amount":            "2",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_count":                 "1",
						"description":                "创建包含-描述-资源组-Protocol",
						"payment_type":               "PayAsYouGo",
						"instance_category":          "x-pack",
						"kibana_private_whitelist.#": "0",
						"protocol":                   "HTTP",
						"version":                    "7.10.0_with_X-Pack",
						"password":                   "Admain@123",
						"resource_group_id":          CHECKSET,
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
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force", "order_action_type", "password", "update_strategy"},
			},
		},
	})
}

var AlicloudElasticsearchInstanceMap11468 = map[string]string{
	"status":        CHECKSET,
	"kibana_domain": CHECKSET,
	"domain":        CHECKSET,
	"kibana_port":   CHECKSET,
	"create_time":   CHECKSET,
	"public_port":   CHECKSET,
}

func AlicloudElasticsearchInstanceBasicDependence11468(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "defaultLPXHTQ" {
  description = "es-instance-case"
  vpc_name    = "es-instance-case"
  dry_run     = false
  cidr_block  = "10.0.10.0/24"
  is_default  = false
}

resource "alicloud_vswitch" "defaultAislbL" {
  vpc_id       = alicloud_vpc.defaultLPXHTQ.id
  zone_id      = "cn-hangzhou-i"
  description  = "tf-case"
  cidr_block   = "10.0.10.0/24"
  vswitch_name = "es-instanc-case"
}


`, name)
}

// Case IS 11477
func TestAccAliCloudElasticsearchInstance_basic11477(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_elasticsearch_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudElasticsearchInstanceMap11477)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ElasticsearchServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeElasticsearchInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccelasticsearch%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudElasticsearchInstanceBasicDependence11477)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_count":  "1",
					"description": "IS",
					"kibana_configuration": []map[string]interface{}{
						{
							"amount": "1",
							"spec":   "elasticsearch.n4.small",
							"disk":   "0",
						},
					},
					"payment_type":             "PayAsYouGo",
					"vswitch_id":               "${alicloud_vswitch.defaultAislbL.id}",
					"instance_category":        "IS",
					"kibana_private_whitelist": []string{},
					"public_whitelist":         []string{},
					"version":                  "7.10.0_with_X-Pack",
					"password":                 "Admain@123",
					"data_node_configuration": []map[string]interface{}{
						{
							"spec":   "openstore.hybrid.i2.2xlarge",
							"amount": "2",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_count":                 "1",
						"description":                "IS",
						"payment_type":               "PayAsYouGo",
						"instance_category":          "IS",
						"kibana_private_whitelist.#": "0",
						"public_whitelist.#":         "0",
						"version":                    "7.10.0_with_X-Pack",
						"password":                   "Admain@123",
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
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force", "order_action_type", "password", "update_strategy"},
			},
		},
	})
}

var AlicloudElasticsearchInstanceMap11477 = map[string]string{
	"status":        CHECKSET,
	"kibana_domain": CHECKSET,
	"domain":        CHECKSET,
	"kibana_port":   CHECKSET,
	"create_time":   CHECKSET,
	"public_port":   CHECKSET,
}

func AlicloudElasticsearchInstanceBasicDependence11477(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc" "defaultLPXHTQ" {
  description = "es-instance-case"
  vpc_name    = "es-instance-case"
  dry_run     = false
  cidr_block  = "10.0.10.0/24"
  is_default  = false
}

resource "alicloud_vswitch" "defaultAislbL" {
  vpc_id       = alicloud_vpc.defaultLPXHTQ.id
  zone_id      = "cn-hangzhou-i"
  description  = "tf-case"
  cidr_block   = "10.0.10.0/24"
  vswitch_name = "es-instanc-case"
}


`, name)
}

// Case 开Kibana私网并修改白名单 11518
func TestAccAliCloudElasticsearchInstance_basic11518(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_elasticsearch_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudElasticsearchInstanceMap11518)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ElasticsearchServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeElasticsearchInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccelasticsearch%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudElasticsearchInstanceBasicDependence11518)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_count":  "1",
					"description": "开Kibana私网并修改白名单",
					"kibana_configuration": []map[string]interface{}{
						{
							"amount": "1",
							"spec":   "elasticsearch.sn1ne.large",
							"disk":   "0",
						},
					},
					"payment_type": "PayAsYouGo",
					"vswitch_id":   "${alicloud_vswitch.defaultAislbL.id}",
					"version":      "5.5.3_with_X-Pack",
					"password":     "Admain@123",
					"data_node_configuration": []map[string]interface{}{
						{
							"disk_type": "cloud_ssd",
							"spec":      "elasticsearch.sn1ne.xlarge.new",
							"disk":      "20",
							"amount":    "3",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_count":   "1",
						"description":  "开Kibana私网并修改白名单",
						"payment_type": "PayAsYouGo",
						"version":      "5.5.3_with_X-Pack",
						"password":     "Admain@123",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"order_action_type":             "UPGRADE",
					"public_whitelist":              []string{},
					"enable_kibana_private_network": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"order_action_type":             "UPGRADE",
						"public_whitelist.#":            "0",
						"enable_kibana_private_network": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_public": "false",
					"kibana_private_whitelist": []string{
						"0.0.0.0/24"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_public":              "false",
						"kibana_private_whitelist.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"kibana_private_whitelist": []string{
						"1.0.0.0/24", "3.0.0.0/24", "2.0.0.0/24"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"kibana_private_whitelist.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_kibana_private_network": "false",
					"kibana_private_whitelist":      []string{},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_kibana_private_network": "false",
						"kibana_private_whitelist.#":    CHECKSET,
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
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force", "order_action_type", "password", "update_strategy"},
			},
		},
	})
}

var AlicloudElasticsearchInstanceMap11518 = map[string]string{
	"status":        CHECKSET,
	"kibana_domain": CHECKSET,
	"domain":        CHECKSET,
	"kibana_port":   CHECKSET,
	"create_time":   CHECKSET,
	"public_port":   CHECKSET,
}

func AlicloudElasticsearchInstanceBasicDependence11518(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc" "defaultLPXHTQ" {
  description = "es-instance-case"
  vpc_name    = "es-instance-case"
  dry_run     = false
  cidr_block  = "10.0.10.0/24"
  is_default  = false
}

resource "alicloud_vswitch" "defaultAislbL" {
  vpc_id       = alicloud_vpc.defaultLPXHTQ.id
  zone_id      = "cn-hangzhou-i"
  description  = "tf-case"
  cidr_block   = "10.0.10.0/24"
  vswitch_name = "es-instanc-case"
}


`, name)
}

// Case 开Kibana公网并修改白名单 11519
func TestAccAliCloudElasticsearchInstance_basic11519(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_elasticsearch_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudElasticsearchInstanceMap11519)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ElasticsearchServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeElasticsearchInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccelasticsearch%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudElasticsearchInstanceBasicDependence11519)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_count":  "1",
					"description": "开Kibana公网并修改白名单",
					"kibana_configuration": []map[string]interface{}{
						{
							"amount": "1",
							"spec":   "elasticsearch.sn1ne.large",
							"disk":   "0",
						},
					},
					"payment_type": "PayAsYouGo",
					"vswitch_id":   "${alicloud_vswitch.defaultAislbL.id}",
					"version":      "5.5.3_with_X-Pack",
					"password":     "Admain@123",
					"data_node_configuration": []map[string]interface{}{
						{
							"disk_type": "cloud_ssd",
							"spec":      "elasticsearch.sn1ne.xlarge.new",
							"disk":      "20",
							"amount":    "2",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_count":   "1",
						"description":  "开Kibana公网并修改白名单",
						"payment_type": "PayAsYouGo",
						"version":      "5.5.3_with_X-Pack",
						"password":     "Admain@123",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"order_action_type":            "UPGRADE",
					"public_whitelist":             []string{},
					"enable_kibana_public_network": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"order_action_type":            "UPGRADE",
						"public_whitelist.#":           "0",
						"enable_kibana_public_network": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"kibana_private_whitelist": []string{},
					"kibana_whitelist": []string{
						"::1", "0.0.0.0/24"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"kibana_private_whitelist.#": "0",
						"kibana_whitelist.#":         "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"kibana_whitelist": []string{
						"::1", "1.0.0.0/24", "3.0.0.0/24", "2.0.0.0/24"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"kibana_whitelist.#": "4",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_kibana_public_network": "false",
					"kibana_whitelist":             []string{},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_kibana_public_network": "false",
						"kibana_whitelist.#":           CHECKSET,
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
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force", "order_action_type", "password", "update_strategy"},
			},
		},
	})
}

var AlicloudElasticsearchInstanceMap11519 = map[string]string{
	"status":      CHECKSET,
	"domain":      CHECKSET,
	"kibana_port": CHECKSET,
	"create_time": CHECKSET,
	"public_port": CHECKSET,
}

func AlicloudElasticsearchInstanceBasicDependence11519(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc" "defaultLPXHTQ" {
  description = "es-instance-case"
  vpc_name    = "es-instance-case"
  dry_run     = false
  cidr_block  = "10.0.10.0/24"
  is_default  = false
}

resource "alicloud_vswitch" "defaultAislbL" {
  vpc_id       = alicloud_vpc.defaultLPXHTQ.id
  zone_id      = "cn-hangzhou-i"
  description  = "tf-case"
  cidr_block   = "10.0.10.0/24"
  vswitch_name = "es-instanc-case"
}


`, name)
}

// Case tag 11456
func TestAccAliCloudElasticsearchInstance_basic11456(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_elasticsearch_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudElasticsearchInstanceMap11456)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ElasticsearchServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeElasticsearchInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccelasticsearch%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudElasticsearchInstanceBasicDependence11456)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_count": "1",
					"vswitch_id": "${alicloud_vswitch.defaultAislbL.id}",
					"kibana_configuration": []map[string]interface{}{
						{
							"amount": "1",
							"spec":   "elasticsearch.sn1ne.large",
							"disk":   "0",
						},
					},
					"payment_type": "PayAsYouGo",
					"description":  "Tag",
					"version":      "7.10.0_with_X-Pack",
					"password":     "Admain@123",
					"data_node_configuration": []map[string]interface{}{
						{
							"disk_type": "cloud_ssd",
							"spec":      "elasticsearch.sn1ne.xlarge.new",
							"disk":      "20",
							"amount":    "2",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_count":   "1",
						"payment_type": "PayAsYouGo",
						"description":  "Tag",
						"version":      "7.10.0_with_X-Pack",
						"password":     "Admain@123",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"order_action_type": "UPGRADE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"order_action_type": "UPGRADE",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
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
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force", "order_action_type", "password", "update_strategy"},
			},
		},
	})
}

var AlicloudElasticsearchInstanceMap11456 = map[string]string{
	"status":        CHECKSET,
	"kibana_domain": CHECKSET,
	"domain":        CHECKSET,
	"kibana_port":   CHECKSET,
	"create_time":   CHECKSET,
	"public_port":   CHECKSET,
}

func AlicloudElasticsearchInstanceBasicDependence11456(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc" "defaultLPXHTQ" {
  description = "es-instance-case"
  vpc_name    = "es-instance-case"
  dry_run     = false
  cidr_block  = "10.0.10.0/24"
  is_default  = false
}

resource "alicloud_vswitch" "defaultAislbL" {
  vpc_id       = alicloud_vpc.defaultLPXHTQ.id
  zone_id      = "cn-hangzhou-i"
  description  = "tf-case"
  cidr_block   = "10.0.10.0/24"
  vswitch_name = "es-instanc-case"
}


`, name)
}

// Case 开Kibana私网并修改白名单 11520
func TestAccAliCloudElasticsearchInstance_basic11520(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_elasticsearch_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudElasticsearchInstanceMap11520)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ElasticsearchServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeElasticsearchInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccelasticsearch%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudElasticsearchInstanceBasicDependence11520)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_count":  "1",
					"description": "开Kibana私网并修改白名单",
					"kibana_configuration": []map[string]interface{}{
						{
							"amount": "1",
							"spec":   "elasticsearch.sn1ne.large",
							"disk":   "0",
						},
					},
					"payment_type": "PayAsYouGo",
					"vswitch_id":   "${alicloud_vswitch.defaultAislbL.id}",
					"version":      "5.5.3_with_X-Pack",
					"password":     "Admain@123",
					"data_node_configuration": []map[string]interface{}{
						{
							"disk_type": "cloud_ssd",
							"spec":      "elasticsearch.sn1ne.xlarge.new",
							"disk":      "20",
							"amount":    "2",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_count":   "1",
						"description":  "开Kibana私网并修改白名单",
						"payment_type": "PayAsYouGo",
						"version":      "5.5.3_with_X-Pack",
						"password":     "Admain@123",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"order_action_type": "UPGRADE",
					"public_whitelist":  []string{},
					"private_whitelist": []string{
						"1.0.0.0/24", "3.0.0.0/24", "2.0.0.0/24"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"order_action_type":   "UPGRADE",
						"public_whitelist.#":  "0",
						"private_whitelist.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"private_whitelist": []string{
						"1.0.0.0/24"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"private_whitelist.#": "1",
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
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force", "order_action_type", "password", "update_strategy"},
			},
		},
	})
}

var AlicloudElasticsearchInstanceMap11520 = map[string]string{
	"status":        CHECKSET,
	"kibana_domain": CHECKSET,
	"domain":        CHECKSET,
	"kibana_port":   CHECKSET,
	"create_time":   CHECKSET,
	"public_port":   CHECKSET,
}

func AlicloudElasticsearchInstanceBasicDependence11520(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc" "defaultLPXHTQ" {
  description = "es-instance-case"
  vpc_name    = "es-instance-case"
  dry_run     = false
  cidr_block  = "10.0.10.0/24"
  is_default  = false
}

resource "alicloud_vswitch" "defaultAislbL" {
  vpc_id       = alicloud_vpc.defaultLPXHTQ.id
  zone_id      = "cn-hangzhou-i"
  description  = "tf-case"
  cidr_block   = "10.0.10.0/24"
  vswitch_name = "es-instanc-case"
}


`, name)
}

// Case 创建包含公私网-白名单 11469
func TestAccAliCloudElasticsearchInstance_basic11469(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_elasticsearch_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudElasticsearchInstanceMap11469)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ElasticsearchServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeElasticsearchInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccelasticsearch%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudElasticsearchInstanceBasicDependence11469)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_count": "1",
					"master_configuration": []map[string]interface{}{
						{
							"disk_type": "cloud_essd",
							"amount":    "3",
							"spec":      "elasticsearch.sn1ne.large.new",
							"disk":      "20",
						},
					},
					"description": "创建包含公私网-白名单",
					"kibana_configuration": []map[string]interface{}{
						{
							"amount": "1",
							"spec":   "elasticsearch.n4.small",
							"disk":   "0",
						},
					},
					"payment_type":             "PayAsYouGo",
					"vswitch_id":               "${alicloud_vswitch.defaultAislbL.id}",
					"instance_category":        "x-pack",
					"enable_public":            "false",
					"kibana_private_whitelist": []string{},
					"public_whitelist":         []string{},
					"kibana_whitelist": []string{
						"127.0.0.1/32"},
					"enable_kibana_public_network":  "true",
					"enable_kibana_private_network": "false",
					"private_whitelist": []string{
						"10.0.10.0/24"},
					"version":  "7.10.0_with_X-Pack",
					"password": "Admain@123",
					"data_node_configuration": []map[string]interface{}{
						{
							"amount":            "2",
							"disk_type":         "cloud_essd",
							"spec":              "elasticsearch.sn1ne.large.new",
							"disk":              "20",
							"performance_level": "PL1",
							"disk_encryption":   "true",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_count":                    "1",
						"description":                   "创建包含公私网-白名单",
						"payment_type":                  "PayAsYouGo",
						"instance_category":             "x-pack",
						"enable_public":                 "false",
						"kibana_private_whitelist.#":    "0",
						"public_whitelist.#":            "0",
						"kibana_whitelist.#":            "1",
						"enable_kibana_public_network":  "true",
						"enable_kibana_private_network": "false",
						"private_whitelist.#":           "1",
						"version":                       "7.10.0_with_X-Pack",
						"password":                      "Admain@123",
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
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force", "order_action_type", "password", "update_strategy"},
			},
		},
	})
}

var AlicloudElasticsearchInstanceMap11469 = map[string]string{
	"status":        CHECKSET,
	"kibana_domain": CHECKSET,
	"domain":        CHECKSET,
	"kibana_port":   CHECKSET,
	"create_time":   CHECKSET,
	"public_port":   CHECKSET,
}

func AlicloudElasticsearchInstanceBasicDependence11469(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc" "defaultLPXHTQ" {
  description = "es-instance-case"
  vpc_name    = "es-instance-case"
  dry_run     = false
  cidr_block  = "10.0.10.0/24"
  is_default  = false
}

resource "alicloud_vswitch" "defaultAislbL" {
  vpc_id       = alicloud_vpc.defaultLPXHTQ.id
  zone_id      = "cn-hangzhou-i"
  description  = "tf-case"
  cidr_block   = "10.0.10.0/24"
  vswitch_name = "es-instanc-case"
}


`, name)
}

// Case 修改SettingConfig 11508
func TestAccAliCloudElasticsearchInstance_basic11508(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_elasticsearch_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudElasticsearchInstanceMap11508)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ElasticsearchServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeElasticsearchInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccelasticsearch%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudElasticsearchInstanceBasicDependence11508)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description":       "修改SettingConfig",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"version":           "7.10.0_with_X-Pack",
					"kibana_configuration": []map[string]interface{}{
						{
							"amount": "1",
							"spec":   "elasticsearch.sn1ne.large",
							"disk":   "0",
						},
					},
					"payment_type": "PayAsYouGo",
					"data_node_configuration": []map[string]interface{}{
						{
							"disk_type": "cloud_ssd",
							"spec":      "elasticsearch.sn1ne.xlarge.new",
							"disk":      "20",
							"amount":    "2",
						},
					},
					"password":          "TerraformAdmin@123!",
					"vswitch_id":        "${alicloud_vswitch.defaultAislbL.id}",
					"instance_category": "x-pack",
					"zone_count":        "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":       "修改SettingConfig",
						"resource_group_id": CHECKSET,
						"version":           "7.10.0_with_X-Pack",
						"payment_type":      "PayAsYouGo",
						"password":          "TerraformAdmin@123!",
						"instance_category": "x-pack",
						"zone_count":        "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"order_action_type": "UPGRADE",
					"update_strategy":   "intelligent",
					"setting_config": map[string]string{
						"\"action.auto_create_index\"":         "+.*,-*",
						"\"action.destructive_requires_name\"": "false",
						"\"xpack.security.audit.enabled\"":     "true",
						"\"xpack.watcher.enabled\"":            "false",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"order_action_type":                               "UPGRADE",
						"setting_config.action.auto_create_index":         "+.*,-*",
						"setting_config.action.destructive_requires_name": "false",
						"setting_config.xpack.security.audit.enabled":     "true",
						"setting_config.xpack.watcher.enabled":            "false",
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
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force", "order_action_type", "password", "update_strategy"},
			},
		},
	})
}

var AlicloudElasticsearchInstanceMap11508 = map[string]string{
	"status":        CHECKSET,
	"kibana_domain": CHECKSET,
	"domain":        CHECKSET,
	"kibana_port":   CHECKSET,
	"create_time":   CHECKSET,
	"public_port":   CHECKSET,
}

func AlicloudElasticsearchInstanceBasicDependence11508(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "defaultLPXHTQ" {
  description = "es-instance-case"
  vpc_name    = "es-instance-case"
  dry_run     = false
  cidr_block  = "10.0.10.0/24"
  is_default  = false
}

resource "alicloud_vswitch" "defaultAislbL" {
  vpc_id       = alicloud_vpc.defaultLPXHTQ.id
  zone_id      = "cn-hangzhou-i"
  description  = "tf-case"
  cidr_block   = "10.0.10.0/24"
  vswitch_name = "es-instanc-case"
}


`, name)
}

// Case 修改KibanaConfiguration 11514
func TestAccAliCloudElasticsearchInstance_basic11514(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_elasticsearch_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudElasticsearchInstanceMap11514)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ElasticsearchServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeElasticsearchInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccelasticsearch%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudElasticsearchInstanceBasicDependence11514)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description":       "修改KibanaConfiguration",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"version":           "7.10.0_with_X-Pack",
					"kibana_configuration": []map[string]interface{}{
						{
							"amount": "1",
							"spec":   "elasticsearch.sn1ne.large",
							"disk":   "0",
						},
					},
					"payment_type": "PayAsYouGo",
					"data_node_configuration": []map[string]interface{}{
						{
							"disk_type": "cloud_ssd",
							"spec":      "elasticsearch.sn1ne.xlarge.new",
							"disk":      "20",
							"amount":    "2",
						},
					},
					"password":          "TerraformAdmin@123!",
					"vswitch_id":        "${alicloud_vswitch.defaultAislbL.id}",
					"instance_category": "x-pack",
					"zone_count":        "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":       "修改KibanaConfiguration",
						"resource_group_id": CHECKSET,
						"version":           "7.10.0_with_X-Pack",
						"payment_type":      "PayAsYouGo",
						"password":          "TerraformAdmin@123!",
						"instance_category": "x-pack",
						"zone_count":        "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"kibana_configuration": []map[string]interface{}{
						{
							"spec": "elasticsearch.sn2ne.xlarge",
						},
					},
					"order_action_type": "UPGRADE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"order_action_type": "UPGRADE",
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
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force", "order_action_type", "password", "update_strategy"},
			},
		},
	})
}

var AlicloudElasticsearchInstanceMap11514 = map[string]string{
	"status":        CHECKSET,
	"kibana_domain": CHECKSET,
	"domain":        CHECKSET,
	"kibana_port":   CHECKSET,
	"create_time":   CHECKSET,
	"public_port":   CHECKSET,
}

func AlicloudElasticsearchInstanceBasicDependence11514(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "defaultLPXHTQ" {
  description = "es-instance-case"
  vpc_name    = "es-instance-case"
  dry_run     = false
  cidr_block  = "10.0.10.0/24"
  is_default  = false
}

resource "alicloud_vswitch" "defaultAislbL" {
  vpc_id       = alicloud_vpc.defaultLPXHTQ.id
  zone_id      = "cn-hangzhou-i"
  description  = "tf-case"
  cidr_block   = "10.0.10.0/24"
  vswitch_name = "es-instanc-case"
}


`, name)
}

// Case 开es公网并修改白名单 11516
func TestAccAliCloudElasticsearchInstance_basic11516(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_elasticsearch_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudElasticsearchInstanceMap11516)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ElasticsearchServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeElasticsearchInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccelasticsearch%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudElasticsearchInstanceBasicDependence11516)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_count":  "1",
					"description": "开es公网并修改白名单",
					"kibana_configuration": []map[string]interface{}{
						{
							"amount": "1",
							"spec":   "elasticsearch.sn1ne.large",
							"disk":   "0",
						},
					},
					"payment_type": "PayAsYouGo",
					"vswitch_id":   "${alicloud_vswitch.defaultAislbL.id}",
					"version":      "5.5.3_with_X-Pack",
					"password":     "Admain@123",
					"data_node_configuration": []map[string]interface{}{
						{
							"disk_type": "cloud_ssd",
							"spec":      "elasticsearch.sn1ne.xlarge.new",
							"disk":      "20",
							"amount":    "2",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_count":   "1",
						"description":  "开es公网并修改白名单",
						"payment_type": "PayAsYouGo",
						"version":      "5.5.3_with_X-Pack",
						"password":     "Admain@123",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_public":     "true",
					"order_action_type": "UPGRADE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_public":     "true",
						"order_action_type": "UPGRADE",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"public_whitelist": []string{
						"::1", "1.0.0.0/24", "3.0.0.0/24", "2.0.0.0/24"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"public_whitelist.#": "4",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"public_whitelist": []string{
						"::1", "1.0.0.0/24"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"public_whitelist.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_public":    "false",
					"public_whitelist": []string{},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_public":      "false",
						"public_whitelist.#": CHECKSET,
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
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force", "order_action_type", "password", "update_strategy"},
			},
		},
	})
}

var AlicloudElasticsearchInstanceMap11516 = map[string]string{
	"status":        CHECKSET,
	"kibana_domain": CHECKSET,
	"domain":        CHECKSET,
	"kibana_port":   CHECKSET,
	"create_time":   CHECKSET,
	"public_port":   CHECKSET,
}

func AlicloudElasticsearchInstanceBasicDependence11516(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc" "defaultLPXHTQ" {
  description = "es-instance-case"
  vpc_name    = "es-instance-case"
  dry_run     = false
  cidr_block  = "10.0.10.0/24"
  is_default  = false
}

resource "alicloud_vswitch" "defaultAislbL" {
  vpc_id       = alicloud_vpc.defaultLPXHTQ.id
  zone_id      = "cn-hangzhou-i"
  description  = "tf-case"
  cidr_block   = "10.0.10.0/24"
  vswitch_name = "es-instanc-case"
}


`, name)
}

// Case PVL 11480
func TestAccAliCloudElasticsearchInstance_basic11480(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_elasticsearch_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudElasticsearchInstanceMap11480)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ElasticsearchServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeElasticsearchInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccelasticsearch%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudElasticsearchInstanceBasicDependence11480)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_count":  "1",
					"description": "PVL",
					"kibana_configuration": []map[string]interface{}{
						{
							"amount": "1",
							"spec":   "elasticsearch.sn1ne.large",
							"disk":   "0",
						},
					},
					"payment_type": "PayAsYouGo",
					"vswitch_id":   "${alicloud_vswitch.defaultAislbL.id}",
					"data_node_configuration": []map[string]interface{}{
						{
							"disk_type": "cloud_ssd",
							"spec":      "elasticsearch.sn1ne.large.new",
							"disk":      "20",
							"amount":    "3",
						},
					},
					"version":                          "7.10.0_with_X-Pack",
					"password":                         "Admain@123",
					"kibana_private_security_group_id": "${alicloud_security_group.defaultfyX8IE.id}",
					"enable_kibana_private_network":    "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_count":                       "1",
						"description":                      "PVL",
						"payment_type":                     "PayAsYouGo",
						"version":                          "7.10.0_with_X-Pack",
						"password":                         "Admain@123",
						"kibana_private_security_group_id": CHECKSET,
						"enable_kibana_private_network":    "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vswitch_id":                       "${alicloud_vswitch.defaultAislbL.id}",
					"kibana_private_security_group_id": "${alicloud_security_group.default5tuoQf.id}",
					"order_action_type":                "UPGRADE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"kibana_private_security_group_id": CHECKSET,
						"order_action_type":                "UPGRADE",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"kibana_private_security_group_id": "${alicloud_security_group.defaultfyX8IE.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"kibana_private_security_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_kibana_private_network": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_kibana_private_network": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_kibana_private_network": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_kibana_private_network": "true",
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
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force", "order_action_type", "password", "update_strategy"},
			},
		},
	})
}

var AlicloudElasticsearchInstanceMap11480 = map[string]string{
	"status":        CHECKSET,
	"kibana_domain": CHECKSET,
	"domain":        CHECKSET,
	"kibana_port":   CHECKSET,
	"create_time":   CHECKSET,
	"public_port":   CHECKSET,
}

func AlicloudElasticsearchInstanceBasicDependence11480(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc" "defaultLPXHTQ" {
  description = "es-instance-case"
  vpc_name    = "es-instance-case"
  dry_run     = false
  cidr_block  = "10.0.10.0/24"
  is_default  = false
}

resource "alicloud_vswitch" "defaultAislbL" {
  vpc_id       = alicloud_vpc.defaultLPXHTQ.id
  zone_id      = "cn-hangzhou-i"
  description  = "tf-case"
  cidr_block   = "10.0.10.0/24"
  vswitch_name = "es-instanc-case"
}

resource "alicloud_security_group" "default5tuoQf" {
  description         = "terraform"
  security_group_name = "terraform"
  security_group_type = "normal"
  vpc_id              = alicloud_vpc.defaultLPXHTQ.id
}

resource "alicloud_security_group" "defaultfyX8IE" {
  description         = "terraform2"
  security_group_name = "terraform2"
  security_group_type = "normal"
  vpc_id              = alicloud_vpc.defaultLPXHTQ.id
}


`, name)
}

// Case Update_描述-密码-https-http 11506
func TestAccAliCloudElasticsearchInstance_basic11506(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_elasticsearch_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudElasticsearchInstanceMap11506)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ElasticsearchServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeElasticsearchInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccelasticsearch%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudElasticsearchInstanceBasicDependence11506)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description":       "Update_描述-密码-https-http",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"version":           "7.10.0_with_X-Pack",
					"kibana_configuration": []map[string]interface{}{
						{
							"amount": "1",
							"spec":   "elasticsearch.sn1ne.large",
							"disk":   "0",
						},
					},
					"payment_type": "PayAsYouGo",
					"data_node_configuration": []map[string]interface{}{
						{
							"disk_type": "cloud_ssd",
							"spec":      "elasticsearch.sn1ne.xlarge.new",
							"disk":      "20",
							"amount":    "3",
						},
					},
					"password":          "TerraformAdmin@123!",
					"vswitch_id":        "${alicloud_vswitch.defaultAislbL.id}",
					"instance_category": "x-pack",
					"zone_count":        "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":       "Update_描述-密码-https-http",
						"resource_group_id": CHECKSET,
						"version":           "7.10.0_with_X-Pack",
						"payment_type":      "PayAsYouGo",
						"password":          "TerraformAdmin@123!",
						"instance_category": "x-pack",
						"zone_count":        "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":       "描述-密码-https-http",
					"force":             "false",
					"order_action_type": "UPGRADE",
					"protocol":          "HTTPS",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":       "描述-密码-https-http",
						"force":             "false",
						"order_action_type": "UPGRADE",
						"protocol":          "HTTPS",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"password": "Admain@321",
					"protocol": "HTTP",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"password": "Admain@321",
						"protocol": "HTTP",
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
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force", "order_action_type", "password", "update_strategy"},
			},
		},
	})
}

var AlicloudElasticsearchInstanceMap11506 = map[string]string{
	"status":        CHECKSET,
	"kibana_domain": CHECKSET,
	"domain":        CHECKSET,
	"kibana_port":   CHECKSET,
	"create_time":   CHECKSET,
	"public_port":   CHECKSET,
}

func AlicloudElasticsearchInstanceBasicDependence11506(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "defaultLPXHTQ" {
  description = "es-instance-case"
  vpc_name    = "es-instance-case"
  dry_run     = false
  cidr_block  = "10.0.10.0/24"
  is_default  = false
}

resource "alicloud_vswitch" "defaultAislbL" {
  vpc_id       = alicloud_vpc.defaultLPXHTQ.id
  zone_id      = "cn-hangzhou-i"
  description  = "tf-case"
  cidr_block   = "10.0.10.0/24"
  vswitch_name = "es-instanc-case"
}


`, name)
}

// Case Update_NodeAmount 11515
func TestAccAliCloudElasticsearchInstance_basic11515(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_elasticsearch_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudElasticsearchInstanceMap11515)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ElasticsearchServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeElasticsearchInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccelasticsearch%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudElasticsearchInstanceBasicDependence11515)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description":       "Update_NodeAmount",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"version":           "7.10.0_with_X-Pack",
					"kibana_configuration": []map[string]interface{}{
						{
							"amount": "1",
							"spec":   "elasticsearch.sn1ne.large",
							"disk":   "0",
						},
					},
					"payment_type": "PayAsYouGo",
					"data_node_configuration": []map[string]interface{}{
						{
							"disk_type": "cloud_ssd",
							"spec":      "elasticsearch.sn1ne.xlarge.new",
							"disk":      "20",
							"amount":    "3",
						},
					},
					"password":          "TerraformAdmin@123!",
					"vswitch_id":        "${alicloud_vswitch.defaultAislbL.id}",
					"instance_category": "x-pack",
					"zone_count":        "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":       "Update_NodeAmount",
						"resource_group_id": CHECKSET,
						"version":           "7.10.0_with_X-Pack",
						"payment_type":      "PayAsYouGo",
						"password":          "TerraformAdmin@123!",
						"instance_category": "x-pack",
						"zone_count":        "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"data_node_configuration": []map[string]interface{}{
						{
							"disk_type": "cloud_ssd",
							"spec":      "elasticsearch.sn1ne.xlarge.new",
							"disk":      "20",
							"amount":    "4",
						},
					},
					"order_action_type": "UPGRADE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"order_action_type": "UPGRADE",
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
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force", "order_action_type", "password", "update_strategy"},
			},
		},
	})
}

var AlicloudElasticsearchInstanceMap11515 = map[string]string{
	"status":        CHECKSET,
	"kibana_domain": CHECKSET,
	"domain":        CHECKSET,
	"kibana_port":   CHECKSET,
	"create_time":   CHECKSET,
	"public_port":   CHECKSET,
}

func AlicloudElasticsearchInstanceBasicDependence11515(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "defaultLPXHTQ" {
  description = "es-instance-case"
  vpc_name    = "es-instance-case"
  dry_run     = false
  cidr_block  = "10.0.10.0/24"
  is_default  = false
}

resource "alicloud_vswitch" "defaultAislbL" {
  vpc_id       = alicloud_vpc.defaultLPXHTQ.id
  zone_id      = "cn-hangzhou-i"
  description  = "tf-case"
  cidr_block   = "10.0.10.0/24"
  vswitch_name = "es-instanc-case"
}


`, name)
}

// Case Update_DataNodeConfiguration 11512
func TestAccAliCloudElasticsearchInstance_basic11512(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_elasticsearch_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudElasticsearchInstanceMap11512)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ElasticsearchServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeElasticsearchInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccelasticsearch%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudElasticsearchInstanceBasicDependence11512)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description":       "Update_DataNodeConfiguration",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"version":           "7.10.0_with_X-Pack",
					"kibana_configuration": []map[string]interface{}{
						{
							"amount": "1",
							"spec":   "elasticsearch.sn1ne.large",
							"disk":   "0",
						},
					},
					"payment_type": "PayAsYouGo",
					"data_node_configuration": []map[string]interface{}{
						{
							"disk_type": "cloud_ssd",
							"spec":      "elasticsearch.sn1ne.xlarge.new",
							"disk":      "20",
							"amount":    "3",
						},
					},
					"password":          "TerraformAdmin@123!",
					"vswitch_id":        "${alicloud_vswitch.defaultAislbL.id}",
					"instance_category": "x-pack",
					"zone_count":        "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":       "Update_DataNodeConfiguration",
						"resource_group_id": CHECKSET,
						"version":           "7.10.0_with_X-Pack",
						"payment_type":      "PayAsYouGo",
						"password":          "TerraformAdmin@123!",
						"instance_category": "x-pack",
						"zone_count":        "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"data_node_configuration": []map[string]interface{}{
						{
							"disk_type":         "cloud_essd",
							"disk_encryption":   "false",
							"spec":              "elasticsearch.r7a.xlarge",
							"disk":              "461",
							"performance_level": "PL2",
							"amount":            "3",
						},
					},
					"order_action_type": "UPGRADE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"order_action_type": "UPGRADE",
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
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force", "order_action_type", "password", "update_strategy"},
			},
		},
	})
}

var AlicloudElasticsearchInstanceMap11512 = map[string]string{
	"status":        CHECKSET,
	"kibana_domain": CHECKSET,
	"domain":        CHECKSET,
	"kibana_port":   CHECKSET,
	"create_time":   CHECKSET,
	"public_port":   CHECKSET,
}

func AlicloudElasticsearchInstanceBasicDependence11512(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "defaultLPXHTQ" {
  description = "es-instance-case"
  vpc_name    = "es-instance-case"
  dry_run     = false
  cidr_block  = "10.0.10.0/24"
  is_default  = false
}

resource "alicloud_vswitch" "defaultAislbL" {
  vpc_id       = alicloud_vpc.defaultLPXHTQ.id
  zone_id      = "cn-hangzhou-i"
  description  = "tf-case"
  cidr_block   = "10.0.10.0/24"
  vswitch_name = "es-instanc-case"
}


`, name)
}

// Case Update_ClientNodeConfiguration 11511
func TestAccAliCloudElasticsearchInstance_basic11511(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_elasticsearch_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudElasticsearchInstanceMap11511)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ElasticsearchServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeElasticsearchInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccelasticsearch%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudElasticsearchInstanceBasicDependence11511)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description":       "Update_ClientNodeConfiguration",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"version":           "7.10.0_with_X-Pack",
					"kibana_configuration": []map[string]interface{}{
						{
							"amount": "1",
							"spec":   "elasticsearch.sn1ne.large",
							"disk":   "0",
						},
					},
					"payment_type": "PayAsYouGo",
					"data_node_configuration": []map[string]interface{}{
						{
							"disk_type": "cloud_ssd",
							"spec":      "elasticsearch.sn1ne.xlarge.new",
							"disk":      "20",
							"amount":    "3",
						},
					},
					"password":          "TerraformAdmin@123!",
					"vswitch_id":        "${alicloud_vswitch.defaultAislbL.id}",
					"instance_category": "x-pack",
					"zone_count":        "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":       "Update_ClientNodeConfiguration",
						"resource_group_id": CHECKSET,
						"version":           "7.10.0_with_X-Pack",
						"payment_type":      "PayAsYouGo",
						"password":          "TerraformAdmin@123!",
						"instance_category": "x-pack",
						"zone_count":        "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"order_action_type": "UPGRADE",
					"client_node_configuration": []map[string]interface{}{
						{
							"disk_type": "cloud_efficiency",
							"amount":    "3",
							"spec":      "elasticsearch.sn2ne.large",
							"disk":      "20",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"order_action_type": "UPGRADE",
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
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force", "order_action_type", "password", "update_strategy"},
			},
		},
	})
}

var AlicloudElasticsearchInstanceMap11511 = map[string]string{
	"status":        CHECKSET,
	"kibana_domain": CHECKSET,
	"domain":        CHECKSET,
	"kibana_port":   CHECKSET,
	"create_time":   CHECKSET,
	"public_port":   CHECKSET,
}

func AlicloudElasticsearchInstanceBasicDependence11511(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "defaultLPXHTQ" {
  description = "es-instance-case"
  vpc_name    = "es-instance-case"
  dry_run     = false
  cidr_block  = "10.0.10.0/24"
  is_default  = false
}

resource "alicloud_vswitch" "defaultAislbL" {
  vpc_id       = alicloud_vpc.defaultLPXHTQ.id
  zone_id      = "cn-hangzhou-i"
  description  = "tf-case"
  cidr_block   = "10.0.10.0/24"
  vswitch_name = "es-instanc-case"
}


`, name)
}

// Case Update_WarmNodeConfiguration 11510
func TestAccAliCloudElasticsearchInstance_basic11510(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_elasticsearch_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudElasticsearchInstanceMap11510)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ElasticsearchServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeElasticsearchInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccelasticsearch%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudElasticsearchInstanceBasicDependence11510)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description":       "Update_WarmNodeConfiguration",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"version":           "7.10.0_with_X-Pack",
					"kibana_configuration": []map[string]interface{}{
						{
							"amount": "1",
							"spec":   "elasticsearch.sn1ne.large",
							"disk":   "0",
						},
					},
					"payment_type": "PayAsYouGo",
					"data_node_configuration": []map[string]interface{}{
						{
							"disk_type": "cloud_ssd",
							"spec":      "elasticsearch.sn1ne.xlarge.new",
							"disk":      "20",
							"amount":    "3",
						},
					},
					"password":          "TerraformAdmin@123!",
					"vswitch_id":        "${alicloud_vswitch.defaultAislbL.id}",
					"instance_category": "x-pack",
					"zone_count":        "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":       "Update_WarmNodeConfiguration",
						"resource_group_id": CHECKSET,
						"version":           "7.10.0_with_X-Pack",
						"payment_type":      "PayAsYouGo",
						"password":          "TerraformAdmin@123!",
						"instance_category": "x-pack",
						"zone_count":        "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"order_action_type": "UPGRADE",
					"warm_node_configuration": []map[string]interface{}{
						{
							"disk_type":       "cloud_efficiency",
							"disk_encryption": "true",
							"amount":          "3",
							"spec":            "elasticsearch.sn1ne.large",
							"disk":            "500",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"order_action_type": "UPGRADE",
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
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force", "order_action_type", "password", "update_strategy"},
			},
		},
	})
}

var AlicloudElasticsearchInstanceMap11510 = map[string]string{
	"status":        CHECKSET,
	"kibana_domain": CHECKSET,
	"domain":        CHECKSET,
	"kibana_port":   CHECKSET,
	"create_time":   CHECKSET,
	"public_port":   CHECKSET,
}

func AlicloudElasticsearchInstanceBasicDependence11510(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "defaultLPXHTQ" {
  description = "es-instance-case"
  vpc_name    = "es-instance-case"
  dry_run     = false
  cidr_block  = "10.0.10.0/24"
  is_default  = false
}

resource "alicloud_vswitch" "defaultAislbL" {
  vpc_id       = alicloud_vpc.defaultLPXHTQ.id
  zone_id      = "cn-hangzhou-i"
  description  = "tf-case"
  cidr_block   = "10.0.10.0/24"
  vswitch_name = "es-instanc-case"
}


`, name)
}

// Case CreateInstance_ Master-EsConfig-Warm-ClientNode 11470
func TestAccAliCloudElasticsearchInstance_basic11470(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_elasticsearch_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudElasticsearchInstanceMap11470)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ElasticsearchServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeElasticsearchInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccelasticsearch%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudElasticsearchInstanceBasicDependence11470)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_count": "1",
					"master_configuration": []map[string]interface{}{
						{
							"disk_type": "cloud_essd",
							"amount":    "3",
							"spec":      "elasticsearch.sn1ne.large.new",
							"disk":      "20",
						},
					},
					"description": "CreateInstance_ Master-EsConfig-Warm-ClientNode",
					"kibana_configuration": []map[string]interface{}{
						{
							"amount": "1",
							"spec":   "elasticsearch.n4.small",
							"disk":   "0",
						},
					},
					"payment_type": "PayAsYouGo",
					"vswitch_id":   "${alicloud_vswitch.defaultAislbL.id}",
					"warm_node_configuration": []map[string]interface{}{
						{
							"disk_type":       "cloud_efficiency",
							"disk_encryption": "true",
							"amount":          "2",
							"spec":            "elasticsearch.sn1ne.xlarge",
							"disk":            "501",
						},
					},
					"client_node_configuration": []map[string]interface{}{
						{
							"disk_type": "cloud_efficiency",
							"amount":    "2",
							"spec":      "elasticsearch.sn1ne.large",
							"disk":      "20",
						},
					},
					"instance_category": "x-pack",
					"setting_config": map[string]interface{}{
						"\"action.destructive_requires_name\"": "true",
						"\"xpack.watcher.enabled\"":            "false",
						"\"xpack.security.audit.enabled\"":     "false",
						"\"action.auto_create_index\"":         "+.*,-*",
					},
					"version":  "7.10.0_with_X-Pack",
					"password": "Admain@123",
					"data_node_configuration": []map[string]interface{}{
						{
							"disk_type":         "cloud_essd",
							"disk_encryption":   "true",
							"spec":              "elasticsearch.sn1ne.large.new",
							"disk":              "20",
							"performance_level": "PL1",
							"amount":            "2",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_count":        "1",
						"description":       "CreateInstance_ Master-EsConfig-Warm-ClientNode",
						"payment_type":      "PayAsYouGo",
						"instance_category": "x-pack",
						"version":           "7.10.0_with_X-Pack",
						"password":          "Admain@123",
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
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force", "order_action_type", "password", "update_strategy"},
			},
		},
	})
}

var AlicloudElasticsearchInstanceMap11470 = map[string]string{
	"status":        CHECKSET,
	"kibana_domain": CHECKSET,
	"domain":        CHECKSET,
	"kibana_port":   CHECKSET,
	"create_time":   CHECKSET,
	"public_port":   CHECKSET,
}

func AlicloudElasticsearchInstanceBasicDependence11470(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc" "defaultLPXHTQ" {
  description = "es-instance-case"
  vpc_name    = "es-instance-case"
  dry_run     = false
  cidr_block  = "10.0.10.0/24"
  is_default  = false
}

resource "alicloud_vswitch" "defaultAislbL" {
  vpc_id       = alicloud_vpc.defaultLPXHTQ.id
  zone_id      = "cn-hangzhou-i"
  description  = "tf-case"
  cidr_block   = "10.0.10.0/24"
  vswitch_name = "es-instanc-case"
}


`, name)
}

// Case Update_MasterConfiguration 11509
func TestAccAliCloudElasticsearchInstance_basic11509(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_elasticsearch_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudElasticsearchInstanceMap11509)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ElasticsearchServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeElasticsearchInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccelasticsearch%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudElasticsearchInstanceBasicDependence11509)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description":       "Update_MasterConfiguration",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"version":           "7.10.0_with_X-Pack",
					"kibana_configuration": []map[string]interface{}{
						{
							"amount": "1",
							"spec":   "elasticsearch.sn1ne.large",
							"disk":   "0",
						},
					},
					"payment_type": "PayAsYouGo",
					"data_node_configuration": []map[string]interface{}{
						{
							"disk_type": "cloud_ssd",
							"spec":      "elasticsearch.sn1ne.xlarge.new",
							"disk":      "20",
							"amount":    "3",
						},
					},
					"password":          "TerraformAdmin@123!",
					"vswitch_id":        "${alicloud_vswitch.defaultAislbL.id}",
					"instance_category": "x-pack",
					"zone_count":        "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":       "Update_MasterConfiguration",
						"resource_group_id": CHECKSET,
						"version":           "7.10.0_with_X-Pack",
						"payment_type":      "PayAsYouGo",
						"password":          "TerraformAdmin@123!",
						"instance_category": "x-pack",
						"zone_count":        "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"order_action_type": "UPGRADE",
					"master_configuration": []map[string]interface{}{
						{
							"disk_type": "cloud_essd",
							"amount":    "3",
							"spec":      "elasticsearch.sn1ne.xlarge.new",
							"disk":      "20",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"order_action_type": "UPGRADE",
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
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force", "order_action_type", "password", "update_strategy"},
			},
		},
	})
}

var AlicloudElasticsearchInstanceMap11509 = map[string]string{
	"status":        CHECKSET,
	"kibana_domain": CHECKSET,
	"domain":        CHECKSET,
	"kibana_port":   CHECKSET,
	"create_time":   CHECKSET,
	"public_port":   CHECKSET,
}

func AlicloudElasticsearchInstanceBasicDependence11509(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "defaultLPXHTQ" {
  description = "es-instance-case"
  vpc_name    = "es-instance-case"
  dry_run     = false
  cidr_block  = "10.0.10.0/24"
  is_default  = false
}

resource "alicloud_vswitch" "defaultAislbL" {
  vpc_id       = alicloud_vpc.defaultLPXHTQ.id
  zone_id      = "cn-hangzhou-i"
  description  = "tf-case"
  cidr_block   = "10.0.10.0/24"
  vswitch_name = "es-instanc-case"
}


`, name)
}

// Case 预付费转自动续费 11529
func TestAccAliCloudElasticsearchInstance_basic11529(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_elasticsearch_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudElasticsearchInstanceMap11529)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ElasticsearchServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeElasticsearchInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccelasticsearch%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudElasticsearchInstanceBasicDependence11529)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description":  "预付费转自动续费",
					"renew_status": "ManualRenewal",
					"version":      "7.10.0_with_X-Pack",
					"kibana_configuration": []map[string]interface{}{
						{
							"amount": "1",
							"spec":   "elasticsearch.sn1ne.large",
							"disk":   "0",
						},
					},
					"payment_type": "Subscription",
					"data_node_configuration": []map[string]interface{}{
						{
							"disk_type":       "cloud_essd",
							"disk_encryption": "false",
							"spec":            "elasticsearch.sn1ne.2xlarge.new",
							"disk":            "150",
							"amount":          "3",
						},
					},
					"password":          "Admain@123",
					"period":            "1",
					"vswitch_id":        "${alicloud_vswitch.defaultmJ42pm.id}",
					"instance_category": "x-pack",
					"zone_count":        "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":       "预付费转自动续费",
						"renew_status":      "ManualRenewal",
						"version":           "7.10.0_with_X-Pack",
						"payment_type":      "Subscription",
						"password":          "Admain@123",
						"instance_category": "x-pack",
						"zone_count":        "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"renew_status":          "AutoRenewal",
					"auto_renew_duration":   "1",
					"order_action_type":     "UPGRADE",
					"renewal_duration_unit": "M",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"renew_status":          "AutoRenewal",
						"auto_renew_duration":   "1",
						"order_action_type":     "UPGRADE",
						"renewal_duration_unit": "M",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type": "PayAsYouGo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type":          "PayAsYouGo",
						"renew_status":          REMOVEKEY,
						"auto_renew_duration":   REMOVEKEY,
						"renewal_duration_unit": REMOVEKEY,
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
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force", "order_action_type", "password", "update_strategy", "period", "renew_status", "auto_renew_duration", "renewal_duration_unit"},
			},
		},
	})
}

var AlicloudElasticsearchInstanceMap11529 = map[string]string{
	"status":        CHECKSET,
	"kibana_domain": CHECKSET,
	"domain":        CHECKSET,
	"kibana_port":   CHECKSET,
	"create_time":   CHECKSET,
	"public_port":   CHECKSET,
}

func AlicloudElasticsearchInstanceBasicDependence11529(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc" "defaultdZWzwG" {
  is_default  = false
  description = "es-instance-case"
  cidr_block  = "10.0.10.0/24"
  vpc_name    = "es-instance-case"
}

resource "alicloud_vswitch" "defaultmJ42pm" {
  description  = "tf-case"
  vpc_id       = alicloud_vpc.defaultdZWzwG.id
  zone_id      = "cn-hangzhou-i"
  cidr_block   = "10.0.10.0/24"
  vswitch_name = "es-instanc-case"
}


`, name)
}

// Case 后转预再转后 11482
func TestAccAliCloudElasticsearchInstance_basic11482(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_elasticsearch_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudElasticsearchInstanceMap11482)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ElasticsearchServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeElasticsearchInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccelasticsearch%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudElasticsearchInstanceBasicDependence11482)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_count":  "1",
					"description": "后转预再转后",
					"kibana_configuration": []map[string]interface{}{
						{
							"amount": "1",
							"spec":   "elasticsearch.sn1ne.large",
							"disk":   "0",
						},
					},
					"payment_type":      "PayAsYouGo",
					"vswitch_id":        "${alicloud_vswitch.defaultAislbL.id}",
					"version":           "7.10.0_with_X-Pack",
					"instance_category": "x-pack",
					"password":          "Admain@123",
					"data_node_configuration": []map[string]interface{}{
						{
							"disk_type": "cloud_ssd",
							"spec":      "elasticsearch.sn1ne.xlarge.new",
							"disk":      "20",
							"amount":    "3",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_count":        "1",
						"description":       "后转预再转后",
						"payment_type":      "PayAsYouGo",
						"version":           "7.10.0_with_X-Pack",
						"instance_category": "x-pack",
						"password":          "Admain@123",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type":      "Subscription",
					"order_action_type": "UPGRADE",
					"period":            "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type":      "Subscription",
						"order_action_type": "UPGRADE",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type": "PayAsYouGo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type": "PayAsYouGo",
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
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}

var AlicloudElasticsearchInstanceMap11482 = map[string]string{
	"status":        CHECKSET,
	"kibana_domain": CHECKSET,
	"domain":        CHECKSET,
	"kibana_port":   CHECKSET,
	"create_time":   CHECKSET,
	"public_port":   CHECKSET,
}

func AlicloudElasticsearchInstanceBasicDependence11482(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc" "defaultLPXHTQ" {
  description = "es-instance-case"
  vpc_name    = "es-instance-case"
  dry_run     = false
  cidr_block  = "10.0.10.0/24"
  is_default  = false
}

resource "alicloud_vswitch" "defaultAislbL" {
  vpc_id       = alicloud_vpc.defaultLPXHTQ.id
  zone_id      = "cn-hangzhou-i"
  description  = "tf-case"
  cidr_block   = "10.0.10.0/24"
  vswitch_name = "es-instanc-case"
}


`, name)
}

// Test Elasticsearch Instance. <<< Resource test cases, automatically generated.

func TestAccAliCloudElasticsearchInstance_deprecatedToNewFields(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_elasticsearch_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudElasticsearchInstanceMapDeprecatedFields)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ElasticsearchServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeElasticsearchInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccelasticsearch%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudElasticsearchInstanceBasicDependenceDeprecatedFields)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_count":          "1",
					"description":         name,
					"vswitch_id":          "${alicloud_vswitch.default.id}",
					"version":             "7.10.0_with_X-Pack",
					"instance_category":   "x-pack",
					"password":            "Admain@123",
					"payment_type":        "PayAsYouGo",
					"data_node_spec":      "elasticsearch.sn1ne.large.new",
					"data_node_amount":    "2",
					"data_node_disk_size": "20",
					"data_node_disk_type": "cloud_ssd",
					"kibana_node_spec":    "elasticsearch.sn1ne.large",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_count":          "1",
						"description":         name,
						"version":             "7.10.0_with_X-Pack",
						"instance_category":   "x-pack",
						"payment_type":        "PayAsYouGo",
						"data_node_spec":      "elasticsearch.sn1ne.large.new",
						"data_node_amount":    "2",
						"data_node_disk_size": "20",
						"data_node_disk_type": "cloud_ssd",
						"kibana_node_spec":    "elasticsearch.sn1ne.large",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"data_node_spec":      REMOVEKEY,
					"data_node_amount":    REMOVEKEY,
					"data_node_disk_size": REMOVEKEY,
					"data_node_disk_type": REMOVEKEY,
					"kibana_node_spec":    REMOVEKEY,
					"data_node_configuration": []map[string]interface{}{
						{
							"disk_type": "cloud_ssd",
							"spec":      "elasticsearch.sn1ne.large.new",
							"disk":      "20",
							"amount":    "2",
						},
					},
					"kibana_configuration": []map[string]interface{}{
						{
							"amount": "1",
							"spec":   "elasticsearch.sn1ne.large",
							"disk":   "0",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_node_configuration.#":           "1",
						"data_node_configuration.0.disk_type": "cloud_ssd",
						"data_node_configuration.0.spec":      "elasticsearch.sn1ne.large.new",
						"data_node_configuration.0.disk":      "20",
						"data_node_configuration.0.amount":    "2",
						"data_node_spec":                      "elasticsearch.sn1ne.large.new",
						"data_node_amount":                    "2",
						"data_node_disk_size":                 "20",
						"kibana_configuration.#":              "1",
						"kibana_configuration.0.amount":       "1",
						"kibana_configuration.0.spec":         "elasticsearch.sn1ne.large",
						"kibana_configuration.0.disk":         "0",
						"kibana_node_spec":                    "elasticsearch.sn1ne.large",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"data_node_configuration": []map[string]interface{}{
						{
							"disk_type": "cloud_ssd",
							"spec":      "elasticsearch.sn1ne.xlarge.new",
							"disk":      "30",
							"amount":    "3",
						},
					},
					"order_action_type": "UPGRADE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_node_configuration.#":           "1",
						"data_node_configuration.0.disk_type": "cloud_ssd",
						"data_node_configuration.0.spec":      "elasticsearch.sn1ne.xlarge.new",
						"data_node_configuration.0.disk":      "30",
						"data_node_configuration.0.amount":    "3",
						"data_node_spec":                      "elasticsearch.sn1ne.xlarge.new",
						"data_node_amount":                    "3",
						"data_node_disk_size":                 "30",
						"order_action_type":                   "UPGRADE",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"kibana_configuration": []map[string]interface{}{
						{
							"amount": "1",
							"spec":   "elasticsearch.sn2ne.large",
							"disk":   "0",
						},
					},
					"order_action_type": "UPGRADE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"kibana_configuration.#":        "1",
						"kibana_configuration.0.amount": "1",
						"kibana_configuration.0.spec":   "elasticsearch.sn2ne.large",
						"kibana_configuration.0.disk":   "0",
						"kibana_node_spec":              "elasticsearch.sn2ne.large",
						"order_action_type":             "UPGRADE",
					}),
				),
			},
		},
	})
}

var AlicloudElasticsearchInstanceMapDeprecatedFields = map[string]string{
	"status":        CHECKSET,
	"kibana_domain": CHECKSET,
	"domain":        CHECKSET,
	"kibana_port":   CHECKSET,
	"create_time":   CHECKSET,
}

func AlicloudElasticsearchInstanceBasicDependenceDeprecatedFields(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc" "default" {
  description = "es-instance-deprecated-test"
  vpc_name    = var.name
  cidr_block  = "10.0.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  zone_id      = "cn-hangzhou-i"
  description  = "es-instance-deprecated-test"
  cidr_block   = "10.0.1.0/24"
  vswitch_name = var.name
}
`, name)
}
