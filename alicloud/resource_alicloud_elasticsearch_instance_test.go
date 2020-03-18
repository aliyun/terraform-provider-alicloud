package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/acctest"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/elasticsearch"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

const DataNodeSpec = "elasticsearch.n4.small"
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
		fmt.Sprintf("tf-testAcc%s", region),
		fmt.Sprintf("tf_testAcc%s", region),
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
			testAccPreCheckWithNoDefaultVpc(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description":          name,
					"vswitch_id":           "${data.alicloud_vswitches.default.ids[0]}",
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
					"master_node_spec": MasterNodeSpec,
					"kibana_whitelist": []string{"192.168.0.0/24", "127.0.0.1"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"master_node_spec":   MasterNodeSpec,
						"kibana_whitelist.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"master_node_spec":  MasterNodeSpecForUpdate,
					"private_whitelist": []string{"192.168.0.0/24", "127.0.0.1"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"master_node_spec":    MasterNodeSpecForUpdate,
						"private_whitelist.#": "2",
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
			testAccPreCheckWithNoDefaultVpc(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description":          name,
					"vswitch_id":           "${data.alicloud_vswitches.default.ids[0]}",
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
			testAccPreCheckWithNoDefaultVpc(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description":          name,
					"vswitch_id":           "${data.alicloud_vswitches.default.ids[0]}",
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
			testAccPreCheckWithNoDefaultVpc(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description":          name,
					"vswitch_id":           "${alicloud_vswitch.default.id}",
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

var elasticsearchMap = map[string]string{
	"description":          CHECKSET,
	"data_node_spec":       DataNodeSpec,
	"data_node_amount":     DataNodeAmount,
	"data_node_disk_size":  DataNodeDisk,
	"data_node_disk_type":  DataNodeDiskType,
	"instance_charge_type": string(PostPaid),
	"status":               "active",
	"kibana_whitelist.#":   "0",
	"private_whitelist.#":  "0",
	"public_whitelist.#":   "0",
	"master_node_spec":     "",
	"id":                   CHECKSET,
	"domain":               CHECKSET,
	"port":                 CHECKSET,
	"kibana_domain":        CHECKSET,
	"kibana_port":          CHECKSET,
	"vswitch_id":           CHECKSET,
}

func resourceElasticsearchInstanceConfigDependence(name string) string {
	return fmt.Sprintf(`
    %s
    variable "creation" {
		default = "Elasticsearch"
	}

	variable "name" {
		default = "%s"
	}
	`, ElasticsearchInstanceCommonTestCase, name)
}

func resourceElasticsearchInstanceConfigDependence_multi(name string) string {
	return fmt.Sprintf(`
    %s
    variable "creation" {
		default = "Elasticsearch"
	}

	variable "name" {
		default = "%s"
	}
	`, ElasticsearchInstanceCommonTestCase, name)
}
