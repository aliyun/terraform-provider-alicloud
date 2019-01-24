package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/elasticsearch"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func init() {
	resource.AddTestSweepers("alicloud_elasticsearch_instance", &resource.Sweeper{
		Name: "alicloud_elasticsearch_instance",
		F:    testSweepElasticsearch,
		Dependencies: []string{
			"alicloud_vswitch",
		},
	})
}

func testSweepElasticsearch(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("Error getting Alicloud client: %s", err)
	}

	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
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
			return fmt.Errorf("Error listing Elasticsearch instances: %s", err)
		}

		resp, _ := raw.(*elasticsearch.ListInstanceResponse)
		if resp == nil || len(resp.Result) < 1 {
			break
		}

		instances = append(instances, resp.Result...)

		if len(resp.Result) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(req.Page); err != nil {
			return err
		} else {
			req.Page = page
		}
	}

	sweeped := false
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
		if skip {
			log.Printf("[INFO] Skipping Elasticsearch Instance: %s (%s)", description, id)
			continue
		}

		sweeped = true
		log.Printf("[INFO] Deleting Elasticsearch Instance: %s (%s)", description, id)
		req := elasticsearch.CreateDeleteInstanceRequest()
		req.InstanceId = id
		_, err := client.WithElasticsearchClient(func(elasticsearchClient *elasticsearch.Client) (interface{}, error) {
			return elasticsearchClient.DeleteInstance(req)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Elasticsearch Instance (%s (%s)): %s", description, id, err)
		}
	}

	if sweeped {
		// Waiting 30 seconds to eusure these instances have been deleted.
		time.Sleep(30 * time.Second)
	}

	return nil
}

func TestAccAlicloudElasticsearchInstance_basic(t *testing.T) {
	var instance elasticsearch.DescribeInstanceResponse

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "alicloud_elasticsearch_instance.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckElasticsearchDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccElasticsearchInstance_basic(ElasticsearchInstanceCommonTestCase, DataNodeSpec, DataNodeAmount, DataNodeDisk, DataNodeDiskType),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckElasticsearchInstanceExists(
						"alicloud_elasticsearch_instance.foo",
						&instance),
					resource.TestCheckResourceAttr("alicloud_elasticsearch_instance.foo", "description", "tf-testAccES_classic"),
					resource.TestCheckResourceAttr("alicloud_elasticsearch_instance.foo", "data_node_spec", DataNodeSpec),
					resource.TestCheckResourceAttr("alicloud_elasticsearch_instance.foo", "data_node_amount", DataNodeAmount),
					resource.TestCheckResourceAttr("alicloud_elasticsearch_instance.foo", "data_node_disk_size", DataNodeDisk),
					resource.TestCheckResourceAttr("alicloud_elasticsearch_instance.foo", "data_node_disk_type", DataNodeDiskType),
					resource.TestCheckResourceAttr("alicloud_elasticsearch_instance.foo", "instance_charge_type", string(PostPaid)),
					resource.TestCheckResourceAttr("alicloud_elasticsearch_instance.foo", "version", string(ESVersion553WithXPack)),
					resource.TestCheckResourceAttrSet("alicloud_elasticsearch_instance.foo", "domain"),
					resource.TestCheckResourceAttrSet("alicloud_elasticsearch_instance.foo", "port"),
					resource.TestCheckResourceAttr("alicloud_elasticsearch_instance.foo", "kibana_whitelist.#", "1"),
				),
			},
			resource.TestStep{
				Config: testAccElasticsearchInstance_basic_with_kibana_whitelist(ElasticsearchInstanceCommonTestCase, DataNodeSpec, DataNodeAmount, DataNodeDisk, DataNodeDiskType),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckElasticsearchInstanceExists(
						"alicloud_elasticsearch_instance.foo",
						&instance),
					resource.TestCheckResourceAttr("alicloud_elasticsearch_instance.foo", "description", "tf-testAccES_classic"),
					resource.TestCheckResourceAttr("alicloud_elasticsearch_instance.foo", "data_node_spec", DataNodeSpec),
					resource.TestCheckResourceAttr("alicloud_elasticsearch_instance.foo", "data_node_amount", DataNodeAmount),
					resource.TestCheckResourceAttr("alicloud_elasticsearch_instance.foo", "data_node_disk_size", DataNodeDisk),
					resource.TestCheckResourceAttr("alicloud_elasticsearch_instance.foo", "data_node_disk_type", DataNodeDiskType),
					resource.TestCheckResourceAttr("alicloud_elasticsearch_instance.foo", "instance_charge_type", string(PostPaid)),
					resource.TestCheckResourceAttr("alicloud_elasticsearch_instance.foo", "version", string(ESVersion553WithXPack)),
					resource.TestCheckResourceAttrSet("alicloud_elasticsearch_instance.foo", "domain"),
					resource.TestCheckResourceAttrSet("alicloud_elasticsearch_instance.foo", "port"),
					resource.TestCheckResourceAttr("alicloud_elasticsearch_instance.foo", "kibana_whitelist.#", "2"),
				),
			},
		},
	})
}

func TestAccAlicloudElasticsearchInstance_master_and_whitelist(t *testing.T) {
	var instance elasticsearch.DescribeInstanceResponse

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "alicloud_elasticsearch_instance.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckElasticsearchDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccElasticsearchInstance_master(ElasticsearchInstanceCommonTestCase, DataNodeSpec, DataNodeAmount, DataNodeDisk, DataNodeDiskType),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckElasticsearchInstanceExists("alicloud_elasticsearch_instance.foo", &instance),
					resource.TestCheckResourceAttr("alicloud_elasticsearch_instance.foo", "description", "tf-testAccES_classic"),
					resource.TestCheckResourceAttr("alicloud_elasticsearch_instance.foo", "data_node_spec", DataNodeSpec),
					resource.TestCheckResourceAttr("alicloud_elasticsearch_instance.foo", "data_node_amount", DataNodeAmount),
					resource.TestCheckResourceAttr("alicloud_elasticsearch_instance.foo", "data_node_disk_size", DataNodeDisk),
					resource.TestCheckResourceAttr("alicloud_elasticsearch_instance.foo", "data_node_disk_type", DataNodeDiskType),
					resource.TestCheckResourceAttr("alicloud_elasticsearch_instance.foo", "instance_charge_type", string(PostPaid)),
					resource.TestCheckResourceAttr("alicloud_elasticsearch_instance.foo", "version", string(ESVersion553WithXPack)),
					resource.TestCheckResourceAttr("alicloud_elasticsearch_instance.foo", "master_node_spec", DataNodeSpecForUpdate),
					resource.TestCheckResourceAttr("alicloud_elasticsearch_instance.foo", "private_whitelist.#", "1"),
					resource.TestCheckResourceAttr("alicloud_elasticsearch_instance.foo", "kibana_whitelist.#", "1"),
				),
			},
			resource.TestStep{
				Config: testAccElasticsearchInstance_master_whitelist(ElasticsearchInstanceCommonTestCase, DataNodeSpec, DataNodeAmount, DataNodeDisk, DataNodeDiskType),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckElasticsearchInstanceExists("alicloud_elasticsearch_instance.foo", &instance),
					resource.TestCheckResourceAttr("alicloud_elasticsearch_instance.foo", "description", "tf-testAccES_classic"),
					resource.TestCheckResourceAttr("alicloud_elasticsearch_instance.foo", "data_node_spec", DataNodeSpec),
					resource.TestCheckResourceAttr("alicloud_elasticsearch_instance.foo", "data_node_amount", DataNodeAmount),
					resource.TestCheckResourceAttr("alicloud_elasticsearch_instance.foo", "data_node_disk_size", DataNodeDisk),
					resource.TestCheckResourceAttr("alicloud_elasticsearch_instance.foo", "data_node_disk_type", DataNodeDiskType),
					resource.TestCheckResourceAttr("alicloud_elasticsearch_instance.foo", "instance_charge_type", string(PostPaid)),
					resource.TestCheckResourceAttr("alicloud_elasticsearch_instance.foo", "version", string(ESVersion553WithXPack)),
					resource.TestCheckResourceAttr("alicloud_elasticsearch_instance.foo", "master_node_spec", DataNodeSpecForUpdate),
					resource.TestCheckResourceAttr("alicloud_elasticsearch_instance.foo", "private_whitelist.#", "2"),
				),
			},
			resource.TestStep{
				Config: testAccElasticsearchInstance_master_xlarge(ElasticsearchInstanceCommonTestCase, DataNodeSpec, DataNodeAmount, DataNodeDisk, DataNodeDiskType),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckElasticsearchInstanceExists("alicloud_elasticsearch_instance.foo", &instance),
					resource.TestCheckResourceAttr("alicloud_elasticsearch_instance.foo", "description", "tf-testAccES_classic"),
					resource.TestCheckResourceAttr("alicloud_elasticsearch_instance.foo", "data_node_spec", DataNodeSpec),
					resource.TestCheckResourceAttr("alicloud_elasticsearch_instance.foo", "data_node_amount", DataNodeAmount),
					resource.TestCheckResourceAttr("alicloud_elasticsearch_instance.foo", "data_node_disk_size", DataNodeDisk),
					resource.TestCheckResourceAttr("alicloud_elasticsearch_instance.foo", "data_node_disk_type", DataNodeDiskType),
					resource.TestCheckResourceAttr("alicloud_elasticsearch_instance.foo", "instance_charge_type", string(PostPaid)),
					resource.TestCheckResourceAttr("alicloud_elasticsearch_instance.foo", "version", string(ESVersion553WithXPack)),
					resource.TestCheckResourceAttr("alicloud_elasticsearch_instance.foo", "master_node_spec", MasterNodeSpecForUpdate),
					resource.TestCheckResourceAttr("alicloud_elasticsearch_instance.foo", "private_whitelist.#", "2"),
				),
			},
		},
	})
}

func TestAccAlicloudElasticsearchInstance_upgrade(t *testing.T) {
	var instance elasticsearch.DescribeInstanceResponse

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "alicloud_elasticsearch_instance.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckElasticsearchDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccElasticsearchInstance_basic(ElasticsearchInstanceCommonTestCase, DataNodeSpec, DataNodeAmount, DataNodeDisk, DataNodeDiskType),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckElasticsearchInstanceExists("alicloud_elasticsearch_instance.foo", &instance),
					resource.TestCheckResourceAttr("alicloud_elasticsearch_instance.foo", "description", "tf-testAccES_classic"),
					resource.TestCheckResourceAttr("alicloud_elasticsearch_instance.foo", "data_node_spec", DataNodeSpec),
					resource.TestCheckResourceAttr("alicloud_elasticsearch_instance.foo", "data_node_amount", DataNodeAmount),
					resource.TestCheckResourceAttr("alicloud_elasticsearch_instance.foo", "data_node_disk_size", DataNodeDisk),
					resource.TestCheckResourceAttr("alicloud_elasticsearch_instance.foo", "data_node_disk_type", DataNodeDiskType),
					resource.TestCheckResourceAttr("alicloud_elasticsearch_instance.foo", "instance_charge_type", string(PostPaid)),
					resource.TestCheckResourceAttr("alicloud_elasticsearch_instance.foo", "version", string(ESVersion553WithXPack)),
				),
			},
			resource.TestStep{
				Config: testAccElasticsearchInstance_basic(ElasticsearchInstanceCommonTestCase, DataNodeSpecForUpdate, DataNodeAmount, DataNodeDisk, DataNodeDiskType),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckElasticsearchInstanceExists("alicloud_elasticsearch_instance.foo", &instance),
					resource.TestCheckResourceAttr("alicloud_elasticsearch_instance.foo", "description", "tf-testAccES_classic"),
					resource.TestCheckResourceAttr("alicloud_elasticsearch_instance.foo", "data_node_spec", DataNodeSpecForUpdate),
					resource.TestCheckResourceAttr("alicloud_elasticsearch_instance.foo", "data_node_amount", DataNodeAmount),
					resource.TestCheckResourceAttr("alicloud_elasticsearch_instance.foo", "data_node_disk_size", DataNodeDisk),
					resource.TestCheckResourceAttr("alicloud_elasticsearch_instance.foo", "data_node_disk_type", DataNodeDiskType),
					resource.TestCheckResourceAttr("alicloud_elasticsearch_instance.foo", "instance_charge_type", string(PostPaid)),
					resource.TestCheckResourceAttr("alicloud_elasticsearch_instance.foo", "version", string(ESVersion553WithXPack)),
				),
			},
			resource.TestStep{
				Config: testAccElasticsearchInstance_basic(ElasticsearchInstanceCommonTestCase, DataNodeSpecForUpdate, DataNodeAmountForUpdate, DataNodeDisk, DataNodeDiskType),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckElasticsearchInstanceExists("alicloud_elasticsearch_instance.foo", &instance),
					resource.TestCheckResourceAttr("alicloud_elasticsearch_instance.foo", "description", "tf-testAccES_classic"),
					resource.TestCheckResourceAttr("alicloud_elasticsearch_instance.foo", "data_node_spec", DataNodeSpecForUpdate),
					resource.TestCheckResourceAttr("alicloud_elasticsearch_instance.foo", "data_node_amount", DataNodeAmountForUpdate),
					resource.TestCheckResourceAttr("alicloud_elasticsearch_instance.foo", "data_node_disk_size", DataNodeDisk),
					resource.TestCheckResourceAttr("alicloud_elasticsearch_instance.foo", "data_node_disk_type", DataNodeDiskType),
					resource.TestCheckResourceAttr("alicloud_elasticsearch_instance.foo", "instance_charge_type", string(PostPaid)),
					resource.TestCheckResourceAttr("alicloud_elasticsearch_instance.foo", "version", string(ESVersion553WithXPack)),
				),
			},
			resource.TestStep{
				Config: testAccElasticsearchInstance_basic(ElasticsearchInstanceCommonTestCase, DataNodeSpecForUpdate, DataNodeAmountForUpdate, DataNodeDiskForUpdate, DataNodeDiskType),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckElasticsearchInstanceExists("alicloud_elasticsearch_instance.foo", &instance),
					resource.TestCheckResourceAttr("alicloud_elasticsearch_instance.foo", "description", "tf-testAccES_classic"),
					resource.TestCheckResourceAttr("alicloud_elasticsearch_instance.foo", "data_node_spec", DataNodeSpecForUpdate),
					resource.TestCheckResourceAttr("alicloud_elasticsearch_instance.foo", "data_node_amount", DataNodeAmountForUpdate),
					resource.TestCheckResourceAttr("alicloud_elasticsearch_instance.foo", "data_node_disk_size", DataNodeDiskForUpdate),
					resource.TestCheckResourceAttr("alicloud_elasticsearch_instance.foo", "data_node_disk_type", DataNodeDiskType),
					resource.TestCheckResourceAttr("alicloud_elasticsearch_instance.foo", "instance_charge_type", string(PostPaid)),
					resource.TestCheckResourceAttr("alicloud_elasticsearch_instance.foo", "version", string(ESVersion553WithXPack)),
				),
			},
		},
	})
}

func testAccCheckElasticsearchInstanceExists(n string, d *elasticsearch.DescribeInstanceResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No elasticsearch instance ID is set")
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		elasticsearchService := ElasticsearchService{client}
		raw, err := elasticsearchService.DescribeInstance(rs.Primary.ID)
		log.Printf("[DEBUG] check instance %s in %#v", rs.Primary.ID, raw)

		if err != nil {
			return err
		}

		*d = raw
		return nil
	}
}

func testAccCheckElasticsearchDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	elasticsearchService := ElasticsearchService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_elasticsearch_instance" {
			continue
		}

		_, err := elasticsearchService.DescribeInstance(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}
	}

	return nil
}

func testAccElasticsearchInstance_basic(common, spec string, amount string, disk string, diskType string) string {

	return fmt.Sprintf(`
    %s
    variable "creation" {
		default = "Elasticsearch"
	}

	variable "name" {
		default = "tf-testAccES_classic"
	}

	resource "alicloud_elasticsearch_instance" "foo" {
    vswitch_id           = "${alicloud_vswitch.default.id}"
	password             = "Test@Elastic"
    instance_charge_type = "PostPaid"
    description          = "${var.name}"
    version              = "5.5.3_with_X-Pack"
    data_node_spec       = "%s"
    data_node_amount     = "%s"
	data_node_disk_size  = "%s"
    data_node_disk_type  = "%s"
	}
	`, common, spec, amount, disk, diskType)
}

func testAccElasticsearchInstance_basic_with_kibana_whitelist(common, spec string, amount string, disk string, diskType string) string {

	return fmt.Sprintf(`
    %s
    variable "creation" {
		default = "Elasticsearch"
	}

	variable "name" {
		default = "tf-testAccES_classic"
	}

	resource "alicloud_elasticsearch_instance" "foo" {
    vswitch_id           = "${alicloud_vswitch.default.id}"
	password             = "Test@Elastic"
    instance_charge_type = "PostPaid"
    description          = "${var.name}"
    version              = "5.5.3_with_X-Pack"
    data_node_spec       = "%s"
    data_node_amount     = "%s"
	data_node_disk_size  = "%s"
    data_node_disk_type  = "%s"
    kibana_whitelist    = ["192.168.0.0/24", "127.0.0.1"]
	}
	`, common, spec, amount, disk, diskType)
}

func testAccElasticsearchInstance_master(common, spec string, amount string, disk string, diskType string) string {

	return fmt.Sprintf(`
    %s
    variable "creation" {
		default = "Elasticsearch"
	}

	variable "name" {
		default = "tf-testAccES_classic"
	}

	resource "alicloud_elasticsearch_instance" "foo" {
    vswitch_id           = "${alicloud_vswitch.default.id}"
	password             = "Test@Elastic"
    instance_charge_type = "PostPaid"
    description          = "${var.name}"
    version              = "5.5.3_with_X-Pack"
    data_node_spec       = "%s"
    data_node_amount     = "%s"
	data_node_disk_size  = "%s"
    data_node_disk_type  = "%s"
    master_node_spec   = "elasticsearch.sn2ne.large"
	}
	`, common, spec, amount, disk, diskType)
}

func testAccElasticsearchInstance_master_whitelist(common, spec string, amount string, disk string, diskType string) string {

	return fmt.Sprintf(`
    %s
    variable "creation" {
		default = "Elasticsearch"
	}

	variable "name" {
		default = "tf-testAccES_classic"
	}

	resource "alicloud_elasticsearch_instance" "foo" {
    vswitch_id           = "${alicloud_vswitch.default.id}"
	password             = "Test@Elastic"
    instance_charge_type = "PostPaid"
    description          = "${var.name}"
    version              = "5.5.3_with_X-Pack"
    data_node_spec       = "%s"
    data_node_amount     = "%s"
	data_node_disk_size  = "%s"
    data_node_disk_type  = "%s"
    master_node_spec     = "elasticsearch.sn2ne.large"
    private_whitelist    = ["192.168.0.0/24", "127.0.0.1"]
	}
	`, common, spec, amount, disk, diskType)
}

func testAccElasticsearchInstance_master_xlarge(common, spec string, amount string, disk string, diskType string) string {

	return fmt.Sprintf(`
    %s
    variable "creation" {
		default = "Elasticsearch"
	}

	variable "name" {
		default = "tf-testAccES_classic"
	}

	resource "alicloud_elasticsearch_instance" "foo" {
    vswitch_id           = "${alicloud_vswitch.default.id}"
	password             = "Test@Elastic"
    instance_charge_type = "PostPaid"
    description          = "${var.name}"
    version              = "5.5.3_with_X-Pack"
    data_node_spec       = "%s"
    data_node_amount     = "%s"
	data_node_disk_size  = "%s"
    data_node_disk_type  = "%s"
    master_node_spec     = "elasticsearch.sn2ne.xlarge"
    private_whitelist    = ["192.168.0.0/24", "127.0.0.1"]
	}
	`, common, spec, amount, disk, diskType)
}
