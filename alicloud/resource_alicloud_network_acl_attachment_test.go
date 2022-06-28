package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func init() {
	resource.AddTestSweepers("alicloud_network_acl_attachment", &resource.Sweeper{
		Name: "alicloud_network_acl_attachment",
		F:    testSweepNetworkAclAttachment,
	})
}

func testSweepNetworkAclAttachment(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	var networkAcls []vpc.NetworkAcl
	request := vpc.CreateDescribeNetworkAclsRequest()
	request.RegionId = client.RegionId
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeNetworkAcls(request)
		})
		if err != nil {
			log.Printf("[ERROR] %s get an error: %#v", request.GetActionName(), err)
		}
		response, _ := raw.(*vpc.DescribeNetworkAclsResponse)
		if len(response.NetworkAcls.NetworkAcl) < 1 {
			break
		}
		networkAcls = append(networkAcls, response.NetworkAcls.NetworkAcl...)

		if len(response.NetworkAcls.NetworkAcl) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(request.PageNumber); err != nil {
			return WrapError(err)
		} else {
			request.PageNumber = page
		}
	}

	for _, nacl := range networkAcls {
		name := nacl.NetworkAclName
		id := nacl.NetworkAclId
		resources := nacl.Resources.Resource
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping Network Acl: %s (%s)", name, id)
			continue
		}
		log.Printf("[INFO] Unassociating Network Acl: %s (%s)", name, id)
		request := vpc.CreateUnassociateNetworkAclRequest()
		request.NetworkAclId = id

		unassociateNetworkAclResource := []vpc.UnassociateNetworkAclResource{}
		for i := 0; i < len(resources); i++ {
			vpcSource := vpc.UnassociateNetworkAclResource{
				ResourceId:   resources[i].ResourceId,
				ResourceType: resources[i].ResourceType,
			}
			unassociateNetworkAclResource = append(unassociateNetworkAclResource, vpcSource)
		}
		request.Resource = &unassociateNetworkAclResource

		_, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.UnassociateNetworkAcl(request)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to unassociate Network Acl (%s (%s)): %s", name, id, err)
		}
	}
	return nil
}

// Skip the test because 'resources' is conflict with 'alicloud_network_acl'.
func SkipTestAccAlicloudVPCNetworkAclAttachment_basic(t *testing.T) {
	resourceId := "alicloud_network_acl_attachment.default"
	ra := resourceAttrInit(resourceId, testAccNaclAttachmentCheckMap)
	rand := acctest.RandInt()
	testAccCheck := ra.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckNetworkAclAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkAclAttachment_create(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkAclAttachmentExists(resourceId),
					testAccCheck(map[string]string{
						"network_acl_id": CHECKSET,
						"resources.#":    "1",
					}),
				),
			},
			{
				Config: testAccNetworkAclAttachment_associate(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkAclAttachmentExists(resourceId),
					testAccCheck(map[string]string{
						"network_acl_id": CHECKSET,
						"resources.#":    "2",
					}),
				),
			},
			{
				Config: testAccNetworkAclAttachment_unassociate(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkAclAttachmentExists(resourceId),
					testAccCheck(map[string]string{
						"network_acl_id": CHECKSET,
						"resources.#":    "1",
					}),
				),
			},
		},
	})
}

func testAccCheckNetworkAclAttachmentExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return WrapError(Error("Not found: %s", n))
		}
		if rs.Primary.ID == "" {
			return WrapError(Error("No Network Acl Attachment ID is set"))
		}
		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		vpcService := VpcService{client}

		parts, err := ParseResourceId(rs.Primary.ID, 2)
		if err != nil {
			return WrapError(err)
		}
		networkAclId := parts[0]

		object, err := vpcService.DescribeNetworkAcl(networkAclId)
		res, _ := object["Resources"].(map[string]interface{})["Resource"].([]interface{})
		resources := make([]vpc.Resource, 0)
		for _, v := range res {
			item := v.(map[string]interface{})
			resources = append(resources, vpc.Resource{
				Status:       fmt.Sprint(item["Status"]),
				ResourceId:   fmt.Sprint(item["ResourceId"]),
				ResourceType: fmt.Sprint(item["ResourceType"]),
			})
		}
		err = vpcService.DescribeNetworkAclAttachment(networkAclId, resources)
		if err != nil {
			return WrapError(err)
		}
		return nil
	}
}

func testAccCheckNetworkAclAttachmentDestroy(s *terraform.State) error {

	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_network_acl_attachment" {
			continue
		}

		parts, err := ParseResourceId(rs.Primary.ID, 2)
		if err != nil {
			return WrapError(err)
		}
		networkAclId := parts[0]

		object, err := vpcService.DescribeNetworkAcl(networkAclId)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}
		vpcResource := []vpc.Resource{}
		resources, _ := object["Resources"].(map[string]interface{})["Resource"].([]interface{})
		for _, e := range resources {
			item := e.(map[string]interface{})

			vpcResource = append(vpcResource, vpc.Resource{
				ResourceId:   item["ResourceId"].(string),
				ResourceType: item["ResourceType"].(string),
			})
		}
		err = vpcService.WaitForNetworkAclAttachment(networkAclId, vpcResource, Deleted, DefaultTimeout)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}
	}
	return nil
}

func testAccNetworkAclAttachment_create(randInt int) string {
	return fmt.Sprintf(`
variable "name" {
	default = "tf-testAcc_network_acl"
}

data "alicloud_zones" "default" {
	available_resource_creation= "VSwitch"
}

resource "alicloud_vpc" "default" {
	vpc_name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}

resource "alicloud_network_acl" "default" {
	vpc_id = "${alicloud_vpc.default.id}"
	network_acl_name = "${var.name}%d"
}


resource "alicloud_vswitch" "default" {
	vpc_id = "${alicloud_vpc.default.id}"
	cidr_block = "172.16.0.0/24"
	zone_id = "${data.alicloud_zones.default.zones.0.id}"
	vswitch_name = "${var.name}"
}

resource "alicloud_vswitch" "default2" {
	vpc_id = "${alicloud_vpc.default.id}"
	cidr_block = "172.16.1.0/24"
	zone_id = "${data.alicloud_zones.default.zones.0.id}"
	vswitch_name = "${var.name}"
}

resource "alicloud_network_acl_attachment" "default" {
	network_acl_id = "${alicloud_network_acl.default.id}"
    resources {
          resource_id = "${alicloud_vswitch.default.id}"
          resource_type = "VSwitch"
        }
}
`, randInt)
}

func testAccNetworkAclAttachment_associate(randInt int) string {
	return fmt.Sprintf(`
variable "name" {
	default = "tf-testAcc_network_acl"
}

data "alicloud_zones" "default" {
	available_resource_creation= "VSwitch"
}

resource "alicloud_vpc" "default" {
	vpc_name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}

resource "alicloud_network_acl" "default" {
	vpc_id = "${alicloud_vpc.default.id}"
	network_acl_name = "${var.name}%d"
}


resource "alicloud_vswitch" "default" {
	vpc_id = "${alicloud_vpc.default.id}"
	cidr_block = "172.16.0.0/24"
	zone_id = "${data.alicloud_zones.default.zones.0.id}"
	vswitch_name = "${var.name}"
}

resource "alicloud_vswitch" "default2" {
	vpc_id = "${alicloud_vpc.default.id}"
	cidr_block = "172.16.1.0/24"
	zone_id = "${data.alicloud_zones.default.zones.0.id}"
	vswitch_name = "${var.name}"
}

resource "alicloud_network_acl_attachment" "default" {
	network_acl_id = "${alicloud_network_acl.default.id}"
    resources {
          resource_id = "${alicloud_vswitch.default.id}"
          resource_type = "VSwitch"
        }
	resources {
          resource_id = "${alicloud_vswitch.default2.id}"
          resource_type = "VSwitch"
        }
}
`, randInt)
}

func testAccNetworkAclAttachment_unassociate(randInt int) string {
	return fmt.Sprintf(`
variable "name" {
	default = "tf-testAcc_network_acl"
}

data "alicloud_zones" "default" {
	available_resource_creation= "VSwitch"
}

resource "alicloud_vpc" "default" {
	vpc_name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}

resource "alicloud_network_acl" "default" {
	vpc_id = "${alicloud_vpc.default.id}"
	network_acl_name = "${var.name}%d"
}


resource "alicloud_vswitch" "default" {
	vpc_id = "${alicloud_vpc.default.id}"
	cidr_block = "172.16.0.0/24"
	zone_id = "${data.alicloud_zones.default.zones.0.id}"
	vswitch_name = "${var.name}"
}

resource "alicloud_vswitch" "default2" {
	vpc_id = "${alicloud_vpc.default.id}"
	cidr_block = "172.16.1.0/24"
	zone_id = "${data.alicloud_zones.default.zones.0.id}"
	vswitch_name = "${var.name}"
}

resource "alicloud_network_acl_attachment" "default" {
	network_acl_id = "${alicloud_network_acl.default.id}"
    resources {
          resource_id = "${alicloud_vswitch.default.id}"
          resource_type = "VSwitch"
        }
}
`, randInt)
}

var testAccNaclAttachmentCheckMap = map[string]string{
	"network_acl_id": CHECKSET,
}
