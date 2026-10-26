package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
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
		if !sweepAll() {
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

func TestAccAliCloudVPCNetworkAclAttachment_basic(t *testing.T) {
	resourceId := "alicloud_network_acl_attachment.default"
	ra := resourceAttrInit(resourceId, map[string]string{
		"network_acl_id": CHECKSET,
	})
	rand := acctest.RandInt()
	testAccCheck := ra.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      testAccCheckNetworkAclAttachmentDestroy,
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
						"network_acl_id":            CHECKSET,
						"resources.#":               "2",
						"resources.0.resource_id":   CHECKSET,
						"resources.0.resource_type": "VSwitch",
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

// TestAccAliCloudVPCNetworkAclAttachment_multiInstance reproduces the
// "inconsistent result after apply" reported in #3350: two attachment instances
// share the same network_acl_id, each managing a different VSwitch. The ACL
// returns every bound resource, so each instance only owns a subset. The
// previous DescribeNetworkAclAttachment logic required all bound resources to
// belong to a single instance and reported the coexisting attachment as absent.
func TestAccAliCloudVPCNetworkAclAttachment_multiInstance(t *testing.T) {
	resourceId1 := "alicloud_network_acl_attachment.default"
	resourceId2 := "alicloud_network_acl_attachment.default2"
	ra1 := resourceAttrInit(resourceId1, map[string]string{
		"network_acl_id": CHECKSET,
	})
	ra2 := resourceAttrInit(resourceId2, map[string]string{
		"network_acl_id": CHECKSET,
	})
	rand := acctest.RandInt()
	testAccCheck1 := ra1.resourceAttrMapUpdateSet()
	testAccCheck2 := ra2.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      testAccCheckNetworkAclAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkAclAttachment_multiInstance(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkAclAttachmentExists(resourceId1),
					testAccCheckNetworkAclAttachmentExists(resourceId2),
					// Each attachment manages a single VSwitch while the ACL has both
					// VSwitches bound; this is the scenario that produced "inconsistent
					// result after apply" when matching required every bound resource to
					// belong to a single instance.
					testAccCheckNetworkAclAttachmentExistsForVSwitch(resourceId1, "alicloud_vswitch.default"),
					testAccCheckNetworkAclAttachmentExistsForVSwitch(resourceId2, "alicloud_vswitch.default2"),
					testAccCheck1(map[string]string{
						"network_acl_id": CHECKSET,
						"resources.#":    "1",
					}),
					testAccCheck2(map[string]string{
						"network_acl_id": CHECKSET,
						"resources.#":    "1",
					}),
				),
			},
			{
				Config:   testAccNetworkAclAttachment_multiInstance(rand),
				PlanOnly: true,
			},
		},
	})
}

// testAccCheckNetworkAclAttachmentExistsForVSwitch verifies that the single
// VSwitch managed by an attachment instance is still bound to the shared ACL.
// It passes only that subset to DescribeNetworkAclAttachment, exercising the
// multi-instance matching path instead of the full bound-resource list used by
// testAccCheckNetworkAclAttachmentExists.
func testAccCheckNetworkAclAttachmentExistsForVSwitch(attResource, vswitchResource string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		attRs, ok := s.RootModule().Resources[attResource]
		if !ok {
			return WrapError(Error("Not found: %s", attResource))
		}
		if attRs.Primary.ID == "" {
			return WrapError(Error("No Network Acl Attachment ID is set"))
		}
		vswRs, ok := s.RootModule().Resources[vswitchResource]
		if !ok {
			return WrapError(Error("Not found: %s", vswitchResource))
		}
		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		vpcService := VpcService{client}

		parts, err := ParseResourceId(attRs.Primary.ID, 2)
		if err != nil {
			return WrapError(err)
		}
		vpcResources := []vpc.Resource{{
			ResourceId:   vswRs.Primary.ID,
			ResourceType: "VSwitch",
		}}
		return vpcService.DescribeNetworkAclAttachment(parts[0], vpcResources)
	}
}

func testAccNetworkAclAttachment_create(randInt int) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAcc_network_acl"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  vpc_name   = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_network_acl" "default" {
  vpc_id           = "${alicloud_vpc.default.id}"
  network_acl_name = "${var.name}%d"
}


resource "alicloud_vswitch" "default" {
  vpc_id       = "${alicloud_vpc.default.id}"
  cidr_block   = "172.16.0.0/24"
  zone_id      = "${data.alicloud_zones.default.zones.0.id}"
  vswitch_name = "${var.name}"
}

resource "alicloud_vswitch" "default2" {
  vpc_id       = "${alicloud_vpc.default.id}"
  cidr_block   = "172.16.1.0/24"
  zone_id      = "${data.alicloud_zones.default.zones.0.id}"
  vswitch_name = "${var.name}"
}

resource "alicloud_network_acl_attachment" "default" {
  network_acl_id = "${alicloud_network_acl.default.id}"
  resources {
    resource_id   = "${alicloud_vswitch.default.id}"
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
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  vpc_name   = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_network_acl" "default" {
  vpc_id           = "${alicloud_vpc.default.id}"
  network_acl_name = "${var.name}%d"
}


resource "alicloud_vswitch" "default" {
  vpc_id       = "${alicloud_vpc.default.id}"
  cidr_block   = "172.16.0.0/24"
  zone_id      = "${data.alicloud_zones.default.zones.0.id}"
  vswitch_name = "${var.name}"
}

resource "alicloud_vswitch" "default2" {
  vpc_id       = "${alicloud_vpc.default.id}"
  cidr_block   = "172.16.1.0/24"
  zone_id      = "${data.alicloud_zones.default.zones.0.id}"
  vswitch_name = "${var.name}"
}

resource "alicloud_network_acl_attachment" "default" {
  network_acl_id = "${alicloud_network_acl.default.id}"
  resources {
    resource_id   = "${alicloud_vswitch.default.id}"
    resource_type = "VSwitch"
  }
  resources {
    resource_id   = "${alicloud_vswitch.default2.id}"
    resource_type = "VSwitch"
  }
}
`, randInt)
}

func testAccNetworkAclAttachment_multiInstance(randInt int) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAcc_network_acl"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  vpc_name   = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_network_acl" "default" {
  vpc_id           = "${alicloud_vpc.default.id}"
  network_acl_name = "${var.name}%d"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = "${alicloud_vpc.default.id}"
  cidr_block   = "172.16.0.0/24"
  zone_id      = "${data.alicloud_zones.default.zones.0.id}"
  vswitch_name = "${var.name}"
}

resource "alicloud_vswitch" "default2" {
  vpc_id       = "${alicloud_vpc.default.id}"
  cidr_block   = "172.16.1.0/24"
  zone_id      = "${data.alicloud_zones.default.zones.0.id}"
  vswitch_name = "${var.name}"
}

resource "alicloud_network_acl_attachment" "default" {
  network_acl_id = "${alicloud_network_acl.default.id}"
  resources {
    resource_id   = "${alicloud_vswitch.default.id}"
    resource_type = "VSwitch"
  }
}

resource "alicloud_network_acl_attachment" "default2" {
  network_acl_id = "${alicloud_network_acl.default.id}"
  // Serialize creation against the first attachment so both VSwitches are bound
  // to the shared ACL by the time this instance is read back, and to avoid
  // TaskConflict from concurrent AssociateNetworkAcl calls on the same ACL.
  depends_on = [alicloud_network_acl_attachment.default]
  resources {
    resource_id   = "${alicloud_vswitch.default2.id}"
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
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  vpc_name   = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_network_acl" "default" {
  vpc_id           = "${alicloud_vpc.default.id}"
  network_acl_name = "${var.name}%d"
}


resource "alicloud_vswitch" "default" {
  vpc_id       = "${alicloud_vpc.default.id}"
  cidr_block   = "172.16.0.0/24"
  zone_id      = "${data.alicloud_zones.default.zones.0.id}"
  vswitch_name = "${var.name}"
}

resource "alicloud_vswitch" "default2" {
  vpc_id       = "${alicloud_vpc.default.id}"
  cidr_block   = "172.16.1.0/24"
  zone_id      = "${data.alicloud_zones.default.zones.0.id}"
  vswitch_name = "${var.name}"
}

resource "alicloud_network_acl_attachment" "default" {
  network_acl_id = "${alicloud_network_acl.default.id}"
  resources {
    resource_id   = "${alicloud_vswitch.default.id}"
    resource_type = "VSwitch"
  }
}
`, randInt)
}
