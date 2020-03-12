package alicloud

import (
	"testing"

	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/hbase"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func init() {
	resource.AddTestSweepers("alicloud_hbase_instance", &resource.Sweeper{
		Name: "alicloud_hbase_instance",
		F:    testSweepHBaseInstances,
	})
}

func testSweepHBaseInstances(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	var insts []hbase.Instance
	req := hbase.CreateDescribeInstancesRequest()
	req.RegionId = client.RegionId
	req.PageSize = requests.NewInteger(PageSizeLarge)
	req.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithHbaseClient(func(hbaseClient *hbase.Client) (interface{}, error) {
			return hbaseClient.DescribeInstances(req)
		})
		if err != nil {
			return fmt.Errorf("Error retrieving HBase Instances: %s", err)
		}
		resp, _ := raw.(*hbase.DescribeInstancesResponse)
		if resp == nil || len(resp.Instances.Instance) < 1 {
			break
		}
		insts = append(insts, resp.Instances.Instance...)

		if len(resp.Instances.Instance) < PageSizeLarge {
			break
		}

		page, err := getNextpageNumber(req.PageNumber)
		if err != nil {
			return err
		}
		req.PageNumber = page
	}

	sweeped := false
	vpcService := VpcService{client}
	for _, v := range insts {
		name := v.InstanceName
		id := v.InstanceId
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		// If a slb name is set by other service, it should be fetched by vswitch name and deleted.
		if skip {
			if need, err := vpcService.needSweepVpc(v.VpcId, ""); err == nil {
				skip = !need
			}

		}

		if skip {
			log.Printf("[INFO] Skipping Hbase Instance: %s (%s)", name, id)
			continue
		}

		log.Printf("[INFO] Deleting HBase Instance: %s (%s)", name, id)
		req := hbase.CreateDeleteInstanceRequest()
		req.ClusterId = id
		_, err := client.WithHbaseClient(func(hbaseClient *hbase.Client) (interface{}, error) {
			return hbaseClient.DeleteInstance(req)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Hbase Instance (%s (%s)): %s", name, id, err)
		} else {
			sweeped = true
		}
	}
	if sweeped {
		// Waiting 30 seconds to eusure these DB instances have been deleted.
		time.Sleep(30 * time.Second)
	}
	return nil
}

const resourceHBaseConfigClassic = `
data "alicloud_zones" "default" {
  available_resource_creation = "HBase"
}

resource "alicloud_hbase_instance" "default" {
  name = "tf_testAccHBase_classic"
  zone_id = "${data.alicloud_zones.default.zones.0.id}"
  engine_version = "2.0"
  master_instance_type = "hbase.sn1.large"
  core_instance_type = "hbase.sn1.large"
  core_instance_quantity = 2
  core_disk_type = "cloud_efficiency"
  core_disk_size = 100
  pay_type = "PostPaid"
  cold_storage_size = 0
  deletion_protection = false
}
`

const resourceHBaseConfigClassicName = `
data "alicloud_zones" "default" {
  available_resource_creation = "HBase"
}

resource "alicloud_hbase_instance" "default" {
  name = "tf_testAccHBase_classic_change_name"
  zone_id = "${data.alicloud_zones.default.zones.0.id}"
  engine_version = "2.0"
  master_instance_type = "hbase.sn1.large"
  core_instance_type = "hbase.sn1.large"
  core_instance_quantity = 2
  core_disk_type = "cloud_efficiency"
  core_disk_size = 100
  pay_type = "PostPaid"
  cold_storage_size = 0
  deletion_protection = false
}
`

const resourceHBaseConfigClassicMainTainTime = `
data "alicloud_zones" "default" {
  available_resource_creation = "HBase"
}

resource "alicloud_hbase_instance" "default" {
  name = "tf_testAccHBase_classic_change_name"
  zone_id = "${data.alicloud_zones.default.zones.0.id}"
  engine_version = "2.0"
  master_instance_type = "hbase.sn1.large"
  core_instance_type = "hbase.sn1.large"
  core_instance_quantity = 2
  core_disk_type = "cloud_efficiency"
  core_disk_size = 100
  pay_type = "PostPaid"
  cold_storage_size = 0
  maintain_start_time = "04:00Z"
  maintain_end_time = "06:00Z"
  deletion_protection = false
}`

const resourceHBaseConfigClassicTags = `
data "alicloud_zones" "default" {
  available_resource_creation = "HBase"
}

resource "alicloud_hbase_instance" "default" {
  name = "tf_testAccHBase_classic_change_name"
  zone_id = "${data.alicloud_zones.default.zones.0.id}"
  engine_version = "2.0"
  master_instance_type = "hbase.sn1.large"
  core_instance_type = "hbase.sn1.large"
  core_instance_quantity = 2
  core_disk_type = "cloud_efficiency"
  core_disk_size = 100
  pay_type = "PostPaid"
  cold_storage_size = 0
  maintain_start_time = "04:00Z"
  maintain_end_time = "06:00Z"
  deletion_protection = false
  tags = {
    Created = "TF"
    For     = "acceptance test"
  }
}`

var resourceHBaseConfigVpc = fmt.Sprintf(`
data "alicloud_vpcs" "default" {
  is_default = "true"
}

data "alicloud_zones" "default" {
  available_resource_creation = "HBase"
}

data "alicloud_vswitches" "default" {
  vpc_id = "${data.alicloud_vpcs.default.ids.0}"
  zone_id = "${data.alicloud_zones.default.zones.0.id}"
}

resource "alicloud_vswitch" "this" {
  count = "${length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1}"
  name = "tf_testAccHBase_vpc"
  vpc_id = "${data.alicloud_vpcs.default.ids.0}"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  cidr_block = "${cidrsubnet(data.alicloud_vpcs.default.vpcs.0.cidr_block, 8, %d)}"
}

resource "alicloud_hbase_instance" "default" {
  name = "tf_testAccHBase_vpc"
  engine_version = "2.0"
  master_instance_type = "hbase.sn1.large"
  core_instance_type = "hbase.sn1.large"
  core_instance_quantity = 2
  core_disk_type = "cloud_efficiency"
  core_disk_size = 100
  pay_type = "PostPaid"
  vswitch_id = "${length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids.0 : concat(alicloud_vswitch.this.*.id, [""])[0]}"
  cold_storage_size = 0
  deletion_protection = false
}
`, acctest.RandIntRange(10, 100))

var resourceHBaseConfigVpcName = fmt.Sprintf(`
data "alicloud_vpcs" "default" {
  is_default = "true"
}
data "alicloud_zones" "default" {
  available_resource_creation = "HBase"
}

data "alicloud_vswitches" "default" {
  vpc_id = "${data.alicloud_vpcs.default.ids.0}"
  zone_id = "${data.alicloud_zones.default.zones.0.id}"
}

resource "alicloud_vswitch" "this" {
  count = "${length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1}"
  name = "tf_testAccHBase_vpc"
  vpc_id = "${data.alicloud_vpcs.default.ids.0}"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  cidr_block = "${cidrsubnet(data.alicloud_vpcs.default.vpcs.0.cidr_block, 8, %d)}"
}

resource "alicloud_hbase_instance" "default" {
  name = "tf_testAccHBase_vpc_change_name"
  engine_version = "2.0"
  master_instance_type = "hbase.sn1.large"
  core_instance_type = "hbase.sn1.large"
  core_instance_quantity = 2
  core_disk_type = "cloud_efficiency"
  core_disk_size = 100
  pay_type = "PostPaid"
  vswitch_id = "${length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids.0 : concat(alicloud_vswitch.this.*.id, [""])[0]}"
  cold_storage_size = 0
  deletion_protection = false
}
`, acctest.RandIntRange(10, 100))

var resourceHBaseConfigVpcMaintainTime = fmt.Sprintf(`
data "alicloud_vpcs" "default" {
  is_default = "true"
}
data "alicloud_zones" "default" {
  available_resource_creation = "HBase"
}

data "alicloud_vswitches" "default" {
  vpc_id = "${data.alicloud_vpcs.default.ids.0}"
  zone_id = "${data.alicloud_zones.default.zones.0.id}"
}

resource "alicloud_vswitch" "this" {
  count = "${length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1}"
  name = "tf_testAccHBase_vpc"
  vpc_id = "${data.alicloud_vpcs.default.ids.0}"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  cidr_block = "${cidrsubnet(data.alicloud_vpcs.default.vpcs.0.cidr_block, 8, %d)}"
}

resource "alicloud_hbase_instance" "default" {
  name = "tf_testAccHBase_vpc_change_name"
  engine_version = "2.0"
  master_instance_type = "hbase.sn1.large"
  core_instance_type = "hbase.sn1.large"
  core_instance_quantity = 2
  core_disk_type = "cloud_efficiency"
  core_disk_size = 100
  pay_type = "PostPaid"
  vswitch_id = "${length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids.0 : concat(alicloud_vswitch.this.*.id, [""])[0]}"
  cold_storage_size = 0
  maintain_start_time = "04:00Z"
  maintain_end_time = "06:00Z"
  deletion_protection = false
}
`, acctest.RandIntRange(10, 100))

var resourceHBaseConfigVpcTags = fmt.Sprintf(`
data "alicloud_vpcs" "default" {
  is_default = "true"
}
data "alicloud_zones" "default" {
  available_resource_creation = "HBase"
}

data "alicloud_vswitches" "default" {
  vpc_id = "${data.alicloud_vpcs.default.ids.0}"
  zone_id = "${data.alicloud_zones.default.zones.0.id}"
}

resource "alicloud_vswitch" "this" {
  count = "${length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1}"
  name = "tf_testAccHBase_vpc"
  vpc_id = "${data.alicloud_vpcs.default.ids.0}"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  cidr_block = "${cidrsubnet(data.alicloud_vpcs.default.vpcs.0.cidr_block, 8, %d)}"
}

resource "alicloud_hbase_instance" "default" {
  name = "tf_testAccHBase_vpc_change_name"
  engine_version = "2.0"
  master_instance_type = "hbase.sn1.large"
  core_instance_type = "hbase.sn1.large"
  core_instance_quantity = 2
  core_disk_type = "cloud_efficiency"
  core_disk_size = 100
  pay_type = "PostPaid"
  vswitch_id = "${length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids.0 : concat(alicloud_vswitch.this.*.id, [""])[0]}"
  cold_storage_size = 0
  maintain_start_time = "04:00Z"
  maintain_end_time = "06:00Z"
  deletion_protection = false
  tags = {
    Created = "TF"
    For     = "acceptance test"
  }
}
`, acctest.RandIntRange(10, 100))

var resourceHBaseConfigMultiInstance = fmt.Sprintf(`
data "alicloud_vpcs" "default" {
  is_default = "true"
}
data "alicloud_zones" "default" {
  available_resource_creation = "HBase"
}

data "alicloud_vswitches" "default" {
  vpc_id = "${data.alicloud_vpcs.default.ids.0}"
  zone_id = "${data.alicloud_zones.default.zones.0.id}"
}

resource "alicloud_vswitch" "this" {
  count = "${length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1}"
  name = "tf_testAccHBase_vpc"
  vpc_id = "${data.alicloud_vpcs.default.ids.0}"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  cidr_block = "${cidrsubnet(data.alicloud_vpcs.default.vpcs.0.cidr_block, 8, %d)}"
}

resource "alicloud_hbase_instance" "default" {
  count = 2
  name = "tf_testAccHBase_multi"
  zone_id = "${data.alicloud_zones.default.zones.0.id}"
  engine_version = "2.0"
  master_instance_type = "hbase.sn1.large"
  core_instance_type = "hbase.sn1.large"
  core_instance_quantity = 2
  core_disk_type = "cloud_efficiency"
  core_disk_size = 100
  pay_type = "PostPaid"
  vswitch_id = "${length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids.0 : concat(alicloud_vswitch.this.*.id, [""])[0]}"
  cold_storage_size = 0
  deletion_protection = false
}
`, acctest.RandIntRange(10, 100))

var resourceHBaseConfigMultiInstanceChangeName = fmt.Sprintf(`
data "alicloud_vpcs" "default" {
  is_default = "true"
}

data "alicloud_zones" "default" {
  available_resource_creation = "HBase"
}

data "alicloud_vswitches" "default" {
  vpc_id = "${data.alicloud_vpcs.default.ids.0}"
  zone_id = "${data.alicloud_zones.default.zones.0.id}"
}

resource "alicloud_vswitch" "this" {
  count = "${length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1}"
  name = "tf_testAccHBase_vpc"
  vpc_id = "${data.alicloud_vpcs.default.ids.0}"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  cidr_block = "${cidrsubnet(data.alicloud_vpcs.default.vpcs.0.cidr_block, 8, %d)}"
}

resource "alicloud_hbase_instance" "default" {
  count = 2
  name = "tf_testAccHBase_multi_change_name"
  zone_id = "${data.alicloud_zones.default.zones.0.id}"
  engine_version = "2.0"
  master_instance_type = "hbase.sn1.large"
  core_instance_type = "hbase.sn1.large"
  core_instance_quantity = 2
  core_disk_type = "cloud_efficiency"
  core_disk_size = 100
  pay_type = "PostPaid"
  vswitch_id = "${length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids.0 : concat(alicloud_vswitch.this.*.id, [""])[0]}"
  cold_storage_size = 0
  deletion_protection = false
}
`, acctest.RandIntRange(10, 100))

func TestAccAlicloudHBaseInstanceClassic(t *testing.T) {
	var instance hbase.DescribeInstanceResponse

	resourceId := "alicloud_hbase_instance.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &HBaseService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeHBaseInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.HBaseClassicSupportedRegions)
		},

		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: resourceHBaseConfigClassic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":               "tf_testAccHBase_classic",
						"engine_version":     "2.0",
						"core_instance_type": "hbase.sn1.large",
						"core_disk_type":     "cloud_efficiency",
						"pay_type":           "PostPaid",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: resourceHBaseConfigClassicName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": "tf_testAccHBase_classic_change_name",
					}),
				),
			},
			{
				Config: resourceHBaseConfigClassicMainTainTime,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                "tf_testAccHBase_classic_change_name",
						"maintain_start_time": "04:00Z",
						"maintain_end_time":   "06:00Z",
					}),
				),
			},
			{
				Config: resourceHBaseConfigClassicTags,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                "tf_testAccHBase_classic_change_name",
						"maintain_start_time": "04:00Z",
						"maintain_end_time":   "06:00Z",
						"tags.%":              "2",
						"tags.Created":        "TF",
						"tags.For":            "acceptance test",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudHBaseInstanceVpc(t *testing.T) {
	var instance hbase.DescribeInstanceResponse

	resourceId := "alicloud_hbase_instance.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &HBaseService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeHBaseInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithNoDefaultVpc(t)
		},

		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: resourceHBaseConfigVpc,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":               "tf_testAccHBase_vpc",
						"engine_version":     "2.0",
						"core_instance_type": "hbase.sn1.large",
						"core_disk_type":     "cloud_efficiency",
						"pay_type":           "PostPaid",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: resourceHBaseConfigVpcName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": "tf_testAccHBase_vpc_change_name",
					}),
				),
			},
			{
				Config: resourceHBaseConfigVpcMaintainTime,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                "tf_testAccHBase_vpc_change_name",
						"maintain_start_time": "04:00Z",
						"maintain_end_time":   "06:00Z",
					}),
				),
			},
			{
				Config: resourceHBaseConfigVpcTags,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                "tf_testAccHBase_vpc_change_name",
						"maintain_start_time": "04:00Z",
						"maintain_end_time":   "06:00Z",
						"tags.%":              "2",
						"tags.Created":        "TF",
						"tags.For":            "acceptance test",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudHBaseInstanceMultiInstance(t *testing.T) {
	var instance hbase.DescribeInstanceResponse

	resourceId := "alicloud_hbase_instance.default.1"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &HBaseService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeHBaseInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithNoDefaultVpc(t)
		},

		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: resourceHBaseConfigMultiInstance,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":               "tf_testAccHBase_multi",
						"engine_version":     "2.0",
						"core_instance_type": "hbase.sn1.large",
						"core_disk_type":     "cloud_efficiency",
						"pay_type":           "PostPaid",
					}),
				),
			},
			{
				Config: resourceHBaseConfigMultiInstanceChangeName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": "tf_testAccHBase_multi_change_name",
					}),
				),
			},
		},
	})
}
