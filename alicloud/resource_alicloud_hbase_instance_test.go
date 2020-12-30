package alicloud

import (
	"testing"

	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/hbase"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
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

const resourceHBaseConfigVpc = `
data "alicloud_zones" "default" {
  available_resource_creation = "HBase"
}

data "alicloud_vswitches" "default" {
  zone_id = "${data.alicloud_zones.default.zones.0.id}"
}

resource "alicloud_hbase_instance" "default" {
  name = "tf_testAccHBase_vpc"
  engine_version = "2.0"
  master_instance_type = "hbase.sn1.large"
  core_instance_type = "hbase.sn1.large"
  core_instance_quantity = 2
  core_disk_type = "cloud_efficiency"
  core_disk_size = 400
  pay_type = "PostPaid"
  vswitch_id = "${data.alicloud_vswitches.default.ids.0}"
  cold_storage_size = 0
  deletion_protection = false
  immediate_delete_flag = true
}
`

const resourceHBaseConfigVpcName = `
data "alicloud_zones" "default" {
  available_resource_creation = "HBase"
}

data "alicloud_vswitches" "default" {
  zone_id = "${data.alicloud_zones.default.zones.0.id}"
}

resource "alicloud_hbase_instance" "default" {
  name = "tf_testAccHBase_vpc_change_name"
  engine_version = "2.0"
  master_instance_type = "hbase.sn1.large"
  core_instance_type = "hbase.sn1.large"
  core_instance_quantity = 2
  core_disk_type = "cloud_efficiency"
  core_disk_size = 400
  pay_type = "PostPaid"
  vswitch_id = "${data.alicloud_vswitches.default.ids.0}"
  cold_storage_size = 0
  deletion_protection = false
  immediate_delete_flag = true
}
`

const resourceHBaseConfigVpcMaintainTime = `
data "alicloud_zones" "default" {
  available_resource_creation = "HBase"
}

data "alicloud_vswitches" "default" {
  zone_id = "${data.alicloud_zones.default.zones.0.id}"
}

resource "alicloud_hbase_instance" "default" {
  name = "tf_testAccHBase_vpc_change_name"
  engine_version = "2.0"
  master_instance_type = "hbase.sn1.large"
  core_instance_type = "hbase.sn1.large"
  core_instance_quantity = 2
  core_disk_type = "cloud_efficiency"
  core_disk_size = 400
  pay_type = "PostPaid"
  vswitch_id = "${data.alicloud_vswitches.default.ids.0}"
  cold_storage_size = 0
  maintain_start_time = "04:00Z"
  maintain_end_time = "06:00Z"
  deletion_protection = false
  immediate_delete_flag = true
}
`

const resourceHBaseConfigVpcTags = `
data "alicloud_zones" "default" {
  available_resource_creation = "HBase"
}

data "alicloud_vswitches" "default" {
  zone_id = "${data.alicloud_zones.default.zones.0.id}"
}

resource "alicloud_hbase_instance" "default" {
  name = "tf_testAccHBase_vpc_change_name"
  engine_version = "2.0"
  master_instance_type = "hbase.sn1.large"
  core_instance_type = "hbase.sn1.large"
  core_instance_quantity = 2
  core_disk_type = "cloud_efficiency"
  core_disk_size = 400
  pay_type = "PostPaid"
  vswitch_id = "${data.alicloud_vswitches.default.ids.0}"
  cold_storage_size = 0
  maintain_start_time = "04:00Z"
  maintain_end_time = "06:00Z"
  deletion_protection = false
  immediate_delete_flag = true
  tags = {
    Created = "TF"
    For     = "acceptance test"
  }
}
`
const resourceHBaseConfigVpcDiskSize = `
data "alicloud_zones" "default" {
  available_resource_creation = "HBase"
}

data "alicloud_vswitches" "default" {
  zone_id = "${data.alicloud_zones.default.zones.0.id}"
}

resource "alicloud_hbase_instance" "default" {
  name = "tf_testAccHBase_vpc_change_name"
  engine_version = "2.0"
  master_instance_type = "hbase.sn1.large"
  core_instance_type = "hbase.sn1.large"
  core_instance_quantity = 2
  core_disk_type = "cloud_efficiency"
  core_disk_size = 440
  pay_type = "PostPaid"
  vswitch_id = "${data.alicloud_vswitches.default.ids.0}"
  cold_storage_size = 0
  maintain_start_time = "04:00Z"
  maintain_end_time = "06:00Z"
  deletion_protection = false
  immediate_delete_flag = true
  tags = {
    Created = "TF"
    For     = "acceptance test"
  }
}
`

const resourceHBaseConfigVpcIpWhite = `
data "alicloud_zones" "default" {
  available_resource_creation = "HBase"
}

data "alicloud_vswitches" "default" {
  zone_id = "${data.alicloud_zones.default.zones.0.id}"
}

resource "alicloud_hbase_instance" "default" {
  name = "tf_testAccHBase_vpc_change_name"
  engine_version = "2.0"
  master_instance_type = "hbase.sn1.large"
  core_instance_type = "hbase.sn1.large"
  core_instance_quantity = 2
  core_disk_type = "cloud_efficiency"
  core_disk_size = 440
  pay_type = "PostPaid"
  vswitch_id = "${data.alicloud_vswitches.default.ids.0}"
  cold_storage_size = 0
  maintain_start_time = "04:00Z"
  maintain_end_time = "06:00Z"
  deletion_protection = false
  immediate_delete_flag = true
  ip_white  = "192.168.1.1"
  tags = {
    Created = "TF"
    For     = "acceptance test"
  }
}
`

const resourceHBaseConfigVpcSecurityAccount = `
data "alicloud_zones" "default" {
  available_resource_creation = "HBase"
}

data "alicloud_vswitches" "default" {
  zone_id = "${data.alicloud_zones.default.zones.0.id}"
}

resource "alicloud_hbase_instance" "default" {
  name = "tf_testAccHBase_vpc_change_name"
  engine_version = "2.0"
  master_instance_type = "hbase.sn1.large"
  core_instance_type = "hbase.sn1.large"
  core_instance_quantity = 2
  core_disk_type = "cloud_efficiency"
  core_disk_size = 440
  pay_type = "PostPaid"
  vswitch_id = "${data.alicloud_vswitches.default.ids.0}"
  cold_storage_size = 0
  maintain_start_time = "04:00Z"
  maintain_end_time = "06:00Z"
  deletion_protection = false
  immediate_delete_flag = true
  account = "admin"
  password = "admin!@#"
  ip_white  = "192.168.1.1"
  tags = {
    Created = "TF"
    For     = "acceptance test"
  }
}
`

const resourceHBaseConfigColdStorage = `
data "alicloud_zones" "default" {
  available_resource_creation = "HBase"
}

data "alicloud_vswitches" "default" {
  zone_id = "${data.alicloud_zones.default.zones.0.id}"
}

resource "alicloud_hbase_instance" "default" {
  name = "tf_testAccHBase_vpc_change_name"
  engine_version = "2.0"
  master_instance_type = "hbase.sn1.large"
  core_instance_type = "hbase.sn1.large"
  core_instance_quantity = 2
  core_disk_type = "cloud_efficiency"
  core_disk_size = 440
  pay_type = "PostPaid"
  vswitch_id = "${data.alicloud_vswitches.default.ids.0}"
  cold_storage_size = 800
  maintain_start_time = "04:00Z"
  maintain_end_time = "06:00Z"
  deletion_protection = false
  immediate_delete_flag = true
  account = "admin"
  password = "admin!@#"
  ip_white  = "192.168.1.1"
  tags = {
    Created = "TF"
    For     = "acceptance test"
  }
}
`

const resourceHBaseConfigPrePaid = `
data "alicloud_zones" "default" {
  available_resource_creation = "HBase"
}

data "alicloud_vswitches" "default" {
  zone_id = "${data.alicloud_zones.default.zones.0.id}"
}

resource "alicloud_hbase_instance" "default" {
  name = "tf_testAccHBase_vpc_change_name"
  engine_version = "2.0"
  master_instance_type = "hbase.sn1.large"
  core_instance_type = "hbase.sn1.large"
  core_instance_quantity = 2
  core_disk_type = "cloud_efficiency"
  core_disk_size = 440
  pay_type = "PrePaid"
  duration = 1
  vswitch_id = "${data.alicloud_vswitches.default.ids.0}"
  cold_storage_size = 800
  maintain_start_time = "04:00Z"
  maintain_end_time = "06:00Z"
  deletion_protection = false
  immediate_delete_flag = true
  account = "admin"
  password = "admin!@#"
  ip_white  = "192.168.1.1"
  tags = {
    Created = "TF"
    For     = "acceptance test"
  }
}
`

const resourceHBaseConfigPostPaid = `
data "alicloud_zones" "default" {
  available_resource_creation = "HBase"
}

data "alicloud_vswitches" "default" {
  zone_id = "${data.alicloud_zones.default.zones.0.id}"
}

resource "alicloud_hbase_instance" "default" {
  name = "tf_testAccHBase_vpc_change_name"
  engine_version = "2.0"
  master_instance_type = "hbase.sn1.large"
  core_instance_type = "hbase.sn1.large"
  core_instance_quantity = 2
  core_disk_type = "cloud_efficiency"
  core_disk_size = 440
  pay_type = "PostPaid"
  duration = 1
  vswitch_id = "${data.alicloud_vswitches.default.ids.0}"
  cold_storage_size = 800
  maintain_start_time = "04:00Z"
  maintain_end_time = "06:00Z"
  deletion_protection = false
  immediate_delete_flag = true
  account = "admin"
  password = "admin!@#"
  ip_white  = "192.168.1.1"
  tags = {
    Created = "TF"
    For     = "acceptance test"
  }
}
`

const resourceHBaseConfigVpcSecurityGroups = `
data "alicloud_zones" "default" {
  available_resource_creation = "HBase"
}

data "alicloud_vpcs" "default" {
  is_default = "true"
}

data "alicloud_vswitches" "default" {
  zone_id = "${data.alicloud_zones.default.zones.0.id}"
}

resource "alicloud_security_group" "default" {
  vpc_id = "${data.alicloud_vpcs.default.vpcs.0.id}"
  inner_access_policy = "Accept"
  name = "tf_testAccHBase_vpc-s-g"
}

resource "alicloud_hbase_instance" "default" {
  name = "tf_testAccHBase_vpc_change_name"
  engine_version = "2.0"
  master_instance_type = "hbase.sn1.large"
  core_instance_type = "hbase.sn1.large"
  core_instance_quantity = 2
  core_disk_type = "cloud_efficiency"
  core_disk_size = 440
  pay_type = "PostPaid"
  duration = 1
  vswitch_id = "${data.alicloud_vswitches.default.ids.0}"
  cold_storage_size = 800
  maintain_start_time = "04:00Z"
  maintain_end_time = "06:00Z"
  deletion_protection = false
  immediate_delete_flag = true
  account = "admin"
  password = "admin!@#"
  ip_white  = "192.168.1.1"
  security_groups = ["${alicloud_security_group.default.id}"]
  tags = {
    Created = "TF"
    For     = "acceptance test"
  }
}
`

const resourceHBaseConfigMultiInstance = `
data "alicloud_zones" "default" {
  available_resource_creation = "HBase"
}

data "alicloud_vswitches" "default" {
  zone_id = "${data.alicloud_zones.default.zones.0.id}"
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
  core_disk_size = 400
  pay_type = "PostPaid"
  vswitch_id = "${data.alicloud_vswitches.default.ids.0}"
  cold_storage_size = 0
  deletion_protection = false
  immediate_delete_flag = true
}
`

const resourceHBaseConfigMultiInstanceChangeName = `
data "alicloud_zones" "default" {
  available_resource_creation = "HBase"
}

data "alicloud_vswitches" "default" {
  zone_id = "${data.alicloud_zones.default.zones.0.id}"
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
  core_disk_size = 400
  pay_type = "PostPaid"
  vswitch_id = "${data.alicloud_vswitches.default.ids.0}"
  cold_storage_size = 0
  deletion_protection = false
  immediate_delete_flag = true
}
`

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
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"immediate_delete_flag"},
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
			{
				Config: resourceHBaseConfigVpcDiskSize,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                "tf_testAccHBase_vpc_change_name",
						"maintain_start_time": "04:00Z",
						"maintain_end_time":   "06:00Z",
						"tags.%":              "2",
						"tags.Created":        "TF",
						"tags.For":            "acceptance test",
						"core_disk_size":      "440",
					}),
				),
			},
			{
				Config: resourceHBaseConfigVpcIpWhite,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                "tf_testAccHBase_vpc_change_name",
						"maintain_start_time": "04:00Z",
						"maintain_end_time":   "06:00Z",
						"tags.%":              "2",
						"tags.Created":        "TF",
						"tags.For":            "acceptance test",
						"ip_white":            "192.168.1.1",
					}),
				),
			},
			{
				Config: resourceHBaseConfigVpcSecurityAccount,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                "tf_testAccHBase_vpc_change_name",
						"maintain_start_time": "04:00Z",
						"maintain_end_time":   "06:00Z",
						"tags.%":              "2",
						"tags.Created":        "TF",
						"tags.For":            "acceptance test",
						"ip_white":            "192.168.1.1",
						"account":             "admin",
					}),
				),
			},
			{
				Config: resourceHBaseConfigColdStorage,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                "tf_testAccHBase_vpc_change_name",
						"maintain_start_time": "04:00Z",
						"maintain_end_time":   "06:00Z",
						"tags.%":              "2",
						"tags.Created":        "TF",
						"tags.For":            "acceptance test",
						"ip_white":            "192.168.1.1",
						"cold_storage_size":   "800",
					}),
				),
			},
			{
				Config: resourceHBaseConfigPrePaid,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                "tf_testAccHBase_vpc_change_name",
						"maintain_start_time": "04:00Z",
						"maintain_end_time":   "06:00Z",
						"tags.%":              "2",
						"tags.Created":        "TF",
						"tags.For":            "acceptance test",
						"ip_white":            "192.168.1.1",
						"cold_storage_size":   "800",
						"pay_type":            "PrePaid",
					}),
				),
			},
			{
				Config: resourceHBaseConfigPostPaid,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                "tf_testAccHBase_vpc_change_name",
						"maintain_start_time": "04:00Z",
						"maintain_end_time":   "06:00Z",
						"tags.%":              "2",
						"tags.Created":        "TF",
						"tags.For":            "acceptance test",
						"ip_white":            "192.168.1.1",
						"cold_storage_size":   "800",
						"pay_type":            "PostPaid",
					}),
				),
			},
			{
				Config: resourceHBaseConfigVpcSecurityGroups,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                "tf_testAccHBase_vpc_change_name",
						"maintain_start_time": "04:00Z",
						"maintain_end_time":   "06:00Z",
						"tags.%":              "2",
						"tags.Created":        "TF",
						"tags.For":            "acceptance test",
						"ip_white":            "192.168.1.1",
						"cold_storage_size":   "800",
						"pay_type":            "PostPaid",
						"security_groups.#":   "1",
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
