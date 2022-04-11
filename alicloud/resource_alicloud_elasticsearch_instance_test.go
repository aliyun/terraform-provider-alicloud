package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/elasticsearch"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const DataNodeSpec = "elasticsearch.sn1ne.large"
const DataNodeAmount = "2"
const DataNodeDisk = "20"
const DataNodeDiskType = "cloud_ssd"

const DataNodeSpecForUpdate = "elasticsearch.sn2ne.large"
const DataNodeAmountForUpdate = "3"
const DataNodeDiskForUpdate = "30"

const DataNodeAmountForMultiZone = "4"
const DefaultZoneAmount = "2"

const MasterNodeSpec = "elasticsearch.sn2ne.large"
const MasterNodeSpecForUpdate = "elasticsearch.sn2ne.xlarge"

const ClientNodeSpec = "elasticsearch.sn2ne.large"
const ClientNodeAmount = "2"

const ClientNodeSpecForUpdate = "elasticsearch.sn2ne.xlarge"
const ClientNodeAmountForUpdate = "3"

const KibanaSpec = "elasticsearch.sn1ne.large"
const KibanaSpecForUpdate = "elasticsearch.sn2ne.large"

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

func TestAccAlicloudElasticsearchInstance_basic(t *testing.T) {
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
					"description":          name,
					"vswitch_id":           "${local.vswitch_id}",
					"version":              "5.5.3_with_X-Pack",
					"password":             "Yourpassword1234",
					"data_node_spec":       DataNodeSpec,
					"data_node_amount":     DataNodeAmount,
					"data_node_disk_size":  DataNodeDisk,
					"data_node_disk_type":  DataNodeDiskType,
					"instance_charge_type": string(PostPaid),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
						"version":     "5.5.3_with_X-Pack",
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
					"data_node_spec": DataNodeSpecForUpdate,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_node_spec": DataNodeSpecForUpdate,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"data_node_amount": DataNodeAmountForUpdate,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_node_amount": DataNodeAmountForUpdate,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"data_node_disk_size": DataNodeDiskForUpdate,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_node_disk_size": DataNodeDiskForUpdate,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"master_node_spec": MasterNodeSpec,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"master_node_spec": MasterNodeSpec,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"private_whitelist": []string{"192.168.0.0/24", "127.0.0.1"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"private_whitelist.#": "2",
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
					"public_whitelist": []string{"192.168.0.0/24", "127.0.0.1"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"public_whitelist.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"kibana_whitelist": []string{"192.168.0.0/24", "127.0.0.1"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"kibana_whitelist.#": "2",
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
					"kibana_private_whitelist": []string{"192.168.0.0/24", "127.0.0.1"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"kibana_private_whitelist.#": "2",
					}),
				),
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
		},
	})
}

func TestAccAlicloudElasticsearchInstance_multizone(t *testing.T) {
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
					"description":          name,
					"vswitch_id":           "${local.vswitch_id}",
					"version":              "5.5.3_with_X-Pack",
					"password":             "Yourpassword1234",
					"data_node_spec":       DataNodeSpec,
					"data_node_amount":     DataNodeAmountForMultiZone,
					"data_node_disk_size":  DataNodeDisk,
					"data_node_disk_type":  DataNodeDiskType,
					"instance_charge_type": string(PostPaid),
					"master_node_spec":     MasterNodeSpec,
					"zone_count":           DefaultZoneAmount,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":      name,
						"version":          "5.5.3_with_X-Pack",
						"data_node_amount": DataNodeAmountForMultiZone,
						"master_node_spec": MasterNodeSpec,
						"zone_count":       DefaultZoneAmount,
					}),
				),
			},
		},
	})
}

func TestAccAlicloudElasticsearchInstance_version(t *testing.T) {
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
					"description":          name,
					"vswitch_id":           "${local.vswitch_id}",
					"version":              "6.3_with_X-Pack",
					"password":             "Yourpassword1234",
					"data_node_spec":       DataNodeSpec,
					"data_node_amount":     DataNodeAmount,
					"data_node_disk_size":  DataNodeDisk,
					"data_node_disk_type":  DataNodeDiskType,
					"instance_charge_type": string(PostPaid),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
						"version":     REGEXMATCH + "^6.3.*_with_X-Pack",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"version": "6.7_with_X-Pack",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"version": REGEXMATCH + "6.7.*_with_X-Pack",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudElasticsearchInstance_multi(t *testing.T) {
	var instance *elasticsearch.DescribeInstanceResponse

	resourceId := "alicloud_elasticsearch_instance.default.1"
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
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceElasticsearchInstanceConfigDependence_multi)

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
					"description":          name,
					"vswitch_id":           "${local.vswitch_id}",
					"version":              "5.5.3_with_X-Pack",
					"password":             "Yourpassword1234",
					"data_node_spec":       DataNodeSpec,
					"data_node_amount":     DataNodeAmount,
					"data_node_disk_size":  DataNodeDisk,
					"data_node_disk_type":  DataNodeDiskType,
					"instance_charge_type": string(PostPaid),
					"count":                "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func TestAccAlicloudElasticsearchInstance_encrypt_disk(t *testing.T) {
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
					"description":              name,
					"vswitch_id":               "${local.vswitch_id}",
					"version":                  "5.5.3_with_X-Pack",
					"password":                 "Yourpassword1234",
					"data_node_spec":           DataNodeSpec,
					"data_node_amount":         DataNodeAmountForMultiZone,
					"data_node_disk_size":      DataNodeDisk,
					"data_node_disk_type":      DataNodeDiskType,
					"data_node_disk_encrypted": "true",
					"instance_charge_type":     string(PostPaid),
					"master_node_spec":         MasterNodeSpec,
					"zone_count":               DefaultZoneAmount,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":              name,
						"version":                  "5.5.3_with_X-Pack",
						"data_node_amount":         DataNodeAmountForMultiZone,
						"data_node_disk_encrypted": "true",
						"master_node_spec":         MasterNodeSpec,
						"zone_count":               DefaultZoneAmount,
					}),
				),
			},
		},
	})
}

func TestAccAlicloudElasticsearchInstance_client_node(t *testing.T) {
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
					"description":          name,
					"vswitch_id":           "${local.vswitch_id}",
					"version":              "6.3_with_X-Pack",
					"password":             "Yourpassword1234",
					"data_node_spec":       DataNodeSpec,
					"data_node_amount":     DataNodeAmount,
					"data_node_disk_size":  DataNodeDisk,
					"data_node_disk_type":  DataNodeDiskType,
					"instance_charge_type": string(PostPaid),
					"client_node_spec":     ClientNodeSpec,
					"client_node_amount":   ClientNodeAmount,
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
					"client_node_spec":   ClientNodeSpecForUpdate,
					"client_node_amount": ClientNodeAmountForUpdate,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"client_node_spec":   ClientNodeSpecForUpdate,
						"client_node_amount": ClientNodeAmountForUpdate,
					}),
				),
			},
		},
	})
}

func TestAccAlicloudElasticsearchInstance_kibana_node(t *testing.T) {
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
					"description":          name,
					"vswitch_id":           "${local.vswitch_id}",
					"version":              "6.3_with_X-Pack",
					"password":             "Yourpassword1234",
					"data_node_spec":       DataNodeSpec,
					"data_node_amount":     DataNodeAmount,
					"data_node_disk_size":  DataNodeDisk,
					"data_node_disk_type":  DataNodeDiskType,
					"instance_charge_type": string(PostPaid),
					"kibana_node_spec":     KibanaSpec,
					"zone_count":           "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"kibana_node_spec": KibanaSpec,
						"zone_count":       "1",
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
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func TestAccAlicloudElasticsearchInstance_https(t *testing.T) {
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
	name := fmt.Sprintf("tf-testAccES-keepit%s%d", defaultRegionToTest, rand)
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
					"description":          name,
					"vswitch_id":           "${local.vswitch_id}",
					"version":              "6.3_with_X-Pack",
					"password":             "Yourpassword1234",
					"data_node_spec":       DataNodeSpec,
					"data_node_amount":     DataNodeAmount,
					"data_node_disk_size":  DataNodeDisk,
					"data_node_disk_type":  DataNodeDiskType,
					"instance_charge_type": string(PostPaid),
					"client_node_spec":     ClientNodeSpec,
					"client_node_amount":   ClientNodeAmount,
					"protocol":             "HTTPS",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocol": "HTTPS",
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

func TestAccAlicloudElasticsearchInstance_setting_config(t *testing.T) {
	var instance *elasticsearch.DescribeInstanceResponse

	resourceId := "alicloud_elasticsearch_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudElasticsearchMap)

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
					"description":          name,
					"vswitch_id":           "${local.vswitch_id}",
					"version":              "6.7_with_X-Pack",
					"password":             "Yourpassword1234",
					"data_node_spec":       DataNodeSpec,
					"data_node_amount":     "3",
					"data_node_disk_size":  DataNodeDisk,
					"data_node_disk_type":  DataNodeDiskType,
					"instance_charge_type": string(PostPaid),
					"setting_config": map[string]string{
						"\"action.auto_create_index\"":         "+.*,-*",
						"\"action.destructive_requires_name\"": "false",
						"\"xpack.security.audit.enabled\"":     "true",
						"\"xpack.security.audit.outputs\"":     "index",
						"\"xpack.watcher.enabled\"":            "false",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"setting_config.action.auto_create_index":         "+.*,-*",
						"setting_config.action.destructive_requires_name": "false",
						"setting_config.xpack.security.audit.enabled":     "true",
						"setting_config.xpack.security.audit.outputs":     "index",
						"setting_config.xpack.watcher.enabled":            "false",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

var elasticsearchMap = map[string]string{
	"description":                   CHECKSET,
	"data_node_spec":                DataNodeSpec,
	"data_node_amount":              DataNodeAmount,
	"data_node_disk_size":           DataNodeDisk,
	"data_node_disk_type":           DataNodeDiskType,
	"instance_charge_type":          string(PostPaid),
	"status":                        "active",
	"private_whitelist.#":           "0",
	"public_whitelist.#":            "0",
	"enable_public":                 "false",
	"kibana_whitelist.#":            "0",
	"enable_kibana_public_network":  "true",
	"kibana_private_whitelist.#":    "0",
	"enable_kibana_private_network": "false",
	"master_node_spec":              "",
	"id":                            CHECKSET,
	"domain":                        CHECKSET,
	"port":                          CHECKSET,
	"kibana_domain":                 CHECKSET,
	"kibana_port":                   CHECKSET,
	"vswitch_id":                    CHECKSET,
}

var AlicloudElasticsearchMap = map[string]string{
	"id":                   CHECKSET,
	"domain":               CHECKSET,
	"port":                 CHECKSET,
	"kibana_domain":        CHECKSET,
	"kibana_port":          CHECKSET,
	"vswitch_id":           CHECKSET,
	"description":          CHECKSET,
	"instance_charge_type": string(PostPaid),
}

func resourceElasticsearchInstanceConfigDependence(name string) string {
	return fmt.Sprintf(`
    %s
	variable "name" {
		default = "%s"
	}
	`, ElasticsearchInstanceCommonTestCase, name)
}

func resourceElasticsearchInstanceConfigDependence_multi(name string) string {
	return fmt.Sprintf(`
    %s
	variable "name" {
		default = "%s"
	}
	`, ElasticsearchInstanceCommonTestCase, name)
}
