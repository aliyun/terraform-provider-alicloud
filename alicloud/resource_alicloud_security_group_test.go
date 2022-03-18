package alicloud

import (
	"fmt"
	"log"
	"os"
	"testing"

	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func init() {
	resource.AddTestSweepers("alicloud_security_group", &resource.Sweeper{
		Name: "alicloud_security_group",
		F:    testSweepSecurityGroups,
		//When implemented, these should be removed firstly
		Dependencies: []string{
			"alicloud_instance",
			"alicloud_ecs_network_interface",
			"alicloud_yundun_bastionhost_instance",
			"alicloud_cs_kubernetes",
		},
	})
}

func testSweepSecurityGroups(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	var groups []ecs.SecurityGroup
	req := ecs.CreateDescribeSecurityGroupsRequest()
	req.RegionId = client.RegionId
	req.PageSize = requests.NewInteger(PageSizeLarge)
	req.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DescribeSecurityGroups(req)
		})
		if err != nil {
			return fmt.Errorf("Error retrieving Security Groups: %s", err)
		}
		resp, _ := raw.(*ecs.DescribeSecurityGroupsResponse)
		if resp == nil || len(resp.SecurityGroups.SecurityGroup) < 1 {
			break
		}
		groups = append(groups, resp.SecurityGroups.SecurityGroup...)

		if len(resp.SecurityGroups.SecurityGroup) < PageSizeLarge {
			break
		}

		page, err := getNextpageNumber(req.PageNumber)
		if err != nil {
			return err
		}
		req.PageNumber = page
	}

	vpcService := VpcService{client}
	ecsService := EcsService{client}
	for _, v := range groups {
		name := v.SecurityGroupName
		id := v.SecurityGroupId
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		// If a Security Group created by other service, it should be fetched by vpc name and deleted.
		if skip {
			if need, err := vpcService.needSweepVpc(v.VpcId, ""); err == nil {
				skip = !need
			}
		}
		if skip {
			log.Printf("[INFO] Skipping Security Group: %s (%s)", name, id)
			continue
		}
		log.Printf("[INFO] Deleting Security Group: %s (%s)", name, id)
		if err := ecsService.sweepSecurityGroup(id); err != nil {
			log.Printf("[ERROR] Failed to delete Security Group (%s (%s)): %s", name, id, err)
		}
	}
	return nil
}

func testAccCheckSecurityGroupDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	ecsService := EcsService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_security_group" {
			continue
		}

		_, err := ecsService.DescribeSecurityGroup(rs.Primary.ID)

		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}
		return WrapError(Error("Error SecurityGroup still exist"))
	}
	return nil
}

func TestAccAlicloudECSSecurityGroupBasic(t *testing.T) {
	var v ecs.DescribeSecurityGroupAttributeResponse
	resourceId := "alicloud_security_group.default"
	ra := resourceAttrInit(resourceId, testAccCheckSecurityBasicMap)
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSecurityGroupConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccCheckSecurityGroupConfigInnerAccess(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"inner_access":        "true",
						"inner_access_policy": "Accept",
					}),
				),
			},
			{
				Config: testAccCheckSecurityGroupConfigName(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": "tf-testAccCheckSecurityGroupName_change",
					}),
				),
			},

			{
				Config: testAccCheckSecurityGroupConfigDescribe(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "tf-testAccCheckSecurityGroupName_describe_change",
					}),
				),
			},
			{
				Config: testAccCheckSecurityGroupConfigTags(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":    "1",
						"tags.Test": REMOVEKEY,
					}),
				),
			},

			{
				Config: testAccCheckSecurityGroupConfigAll(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(testAccCheckSecurityBasicMap),
				),
			},
		},
	})
}

func TestAccAlicloudECSSecurityGroupMulti(t *testing.T) {
	var v ecs.DescribeSecurityGroupAttributeResponse
	resourceId := "alicloud_security_group.default.9"
	ra := resourceAttrInit(resourceId, testAccCheckSecurityBasicMap)
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSecurityGroupConfigMulti(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func testAccCheckSecurityGroupConfigBasic() string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAccCheckSecurityGroupName"
}


resource "alicloud_vpc" "default" {
  vpc_name = "${var.name}_vpc"
  cidr_block = "10.1.0.0/21"
}

resource "alicloud_security_group" "default" {
  vpc_id = "${alicloud_vpc.default.id}"
  resource_group_id = "%s"
  inner_access = false
  name = "${var.name}"
  description = "${var.name}_describe"
  tags = {
		foo  = "foo"
        Test = "Test"
  }
}
`, os.Getenv("ALICLOUD_RESOURCE_GROUP_ID"))
}

func testAccCheckSecurityGroupConfigInnerAccess() string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAccCheckSecurityGroupName"
}


resource "alicloud_vpc" "default" {
  vpc_name = "${var.name}_vpc"
  cidr_block = "10.1.0.0/21"
}

resource "alicloud_security_group" "default" {
  vpc_id = "${alicloud_vpc.default.id}"
  resource_group_id = "%s"
  inner_access_policy = "Accept"
  name = "${var.name}"
  description = "${var.name}_describe"
  tags = {
		foo  = "foo"
        Test = "Test"
  }
}`, os.Getenv("ALICLOUD_RESOURCE_GROUP_ID"))
}

func testAccCheckSecurityGroupConfigName() string {
	return fmt.Sprintf(`

variable "name" {
  default = "tf-testAccCheckSecurityGroupName"
}


resource "alicloud_vpc" "default" {
  vpc_name = "${var.name}_vpc"
  cidr_block = "10.1.0.0/21"
}

resource "alicloud_security_group" "default" {
  vpc_id = "${alicloud_vpc.default.id}"
  resource_group_id = "%s"
  inner_access = true
  name = "${var.name}_change"
  description = "${var.name}_describe"
  tags = {
		foo  = "foo"
        Test = "Test"
  }
}`, os.Getenv("ALICLOUD_RESOURCE_GROUP_ID"))
}

func testAccCheckSecurityGroupConfigDescribe() string {
	return fmt.Sprintf(`

variable "name" {
  default = "tf-testAccCheckSecurityGroupName"
}


resource "alicloud_vpc" "default" {
  vpc_name = "${var.name}_vpc"
  cidr_block = "10.1.0.0/21"
}

resource "alicloud_security_group" "default" {
  vpc_id = "${alicloud_vpc.default.id}"
  resource_group_id = "%s"
  inner_access = true
  name = "${var.name}_change"
  description = "${var.name}_describe_change"
  tags = {
		foo  = "foo"
        Test = "Test"
  }
}`, os.Getenv("ALICLOUD_RESOURCE_GROUP_ID"))
}
func testAccCheckSecurityGroupConfigTags() string {
	return fmt.Sprintf(`

variable "name" {
  default = "tf-testAccCheckSecurityGroupName"
}


resource "alicloud_vpc" "default" {
  vpc_name = "${var.name}_vpc"
  cidr_block = "10.1.0.0/21"
}

resource "alicloud_security_group" "default" {
  vpc_id = "${alicloud_vpc.default.id}"
  resource_group_id = "%s"
  inner_access = true
  name = "${var.name}_change"
  description = "${var.name}_describe_change"
  tags = {
		foo  = "foo"
  }
}`, os.Getenv("ALICLOUD_RESOURCE_GROUP_ID"))
}

func testAccCheckSecurityGroupConfigAll() string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAccCheckSecurityGroupName"
}


resource "alicloud_vpc" "default" {
  vpc_name = "${var.name}_vpc"
  cidr_block = "10.1.0.0/21"
}

resource "alicloud_security_group" "default" {
  vpc_id = "${alicloud_vpc.default.id}"
  resource_group_id = "%s"
  inner_access_policy = "Drop"
  name = "${var.name}"
  description = "${var.name}_describe"
  tags = {
		foo  = "foo"
        Test = "Test"
  }
}`, os.Getenv("ALICLOUD_RESOURCE_GROUP_ID"))
}

func testAccCheckSecurityGroupConfigMulti() string {
	return fmt.Sprintf(`

variable "name" {
  default = "tf-testAccCheckSecurityGroupName"
}


resource "alicloud_vpc" "default" {
  vpc_name = "${var.name}_vpc"
  cidr_block = "10.1.0.0/21"
}

resource "alicloud_security_group" "default" {
  count = 10
  vpc_id = "${alicloud_vpc.default.id}"
  resource_group_id = "%s"
  inner_access = false
  name = "${var.name}"
  description = "${var.name}_describe"
  tags = {
		foo  = "foo"
        Test = "Test"
  }
}`, os.Getenv("ALICLOUD_RESOURCE_GROUP_ID"))
}

var testAccCheckSecurityBasicMap = map[string]string{
	"vpc_id":              CHECKSET,
	"resource_group_id":   os.Getenv("ALICLOUD_RESOURCE_GROUP_ID"),
	"inner_access":        "false",
	"inner_access_policy": "Drop",
	"name":                "tf-testAccCheckSecurityGroupName",
	"description":         "tf-testAccCheckSecurityGroupName_describe",
	"security_group_type": "normal",
	"tags.%":              "2",
	"tags.foo":            "foo",
	"tags.Test":           "Test",
}
