package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func init() {
	resource.AddTestSweepers("alicloud_network_acl", &resource.Sweeper{
		Name: "alicloud_network_acl",
		F:    testSweepNetworkAcl,
		// When implemented, these should be removed firstly
		Dependencies: []string{
			"alicloud_network_acl_attachment",
			"alicloud_network_acl_entries",
		},
	})
}

func testSweepNetworkAcl(region string) error {
	if testSweepPreCheckWithRegions(region, true, connectivity.NetworkAclSupportedRegions) {
		log.Printf("[INFO] Skipping Network Acl unsupported region: %s", region)
		return nil
	}
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
			log.Printf("[ERROR] %s got an error: %#v", request.GetActionName(), err)
			return nil
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
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
		log.Printf("[INFO] Deleting Network Acl: %s (%s)", name, id)
		request := vpc.CreateDeleteNetworkAclRequest()
		request.NetworkAclId = id
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DeleteNetworkAcl(request)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Network Acl (%s (%s)): %s", name, id, err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}
	return nil
}

func TestAccAlicloudNetworkAcl_basic(t *testing.T) {
	var v vpc.NetworkAcl
	resourceId := "alicloud_network_acl.default"
	ra := resourceAttrInit(resourceId, testAccNaclCheckMap)
	serviceFunc := func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandInt()
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, connectivity.NetworkAclSupportedRegions)
		},
		// module name
		IDRefreshName: "alicloud_network_acl.default",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckNetworkAclDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkAcl_create(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_id":      CHECKSET,
						"name":        fmt.Sprintf("tf-testAcc_network_acl%v.abc", rand),
						"description": "tf-testAcc_network_acl",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccNetworkAcl_name(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":        fmt.Sprintf("tf-testAcc_network_acl_modify%v.abc", rand),
						"description": "tf-testAcc_network_acl",
					}),
				),
			},
			{
				Config: testAccNetworkAcl_description(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":        fmt.Sprintf("tf-testAcc_network_acl_modify%v.abc", rand),
						"description": "tf-testAcc_network_acl_modify",
					}),
				),
			},
			{
				Config: testAccNetworkAcl_modify(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":        fmt.Sprintf("tf-testAcc_network_acl%v.abc", rand),
						"description": "tf-testAcc_network_acl",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudNetworkAcl_multi(t *testing.T) {
	var v vpc.NetworkAcl

	ra := resourceAttrInit("alicloud_network_acl.default.2", testAccNaclCheckMap)
	serviceFunc := func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit("alicloud_network_acl.default.2", &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	rand := acctest.RandInt()
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, connectivity.NetworkAclSupportedRegions)
		},
		IDRefreshName: "alicloud_network_acl.default.2",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckNetworkAclDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkAcl_multi(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_id":      CHECKSET,
						"name":        "tf-testAcc_network_acl",
						"description": "tf-testAcc_network_acl2",
					}),
				),
			},
		},
	})
}

func testAccCheckNetworkAclDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_network_acl" {
			continue
		}
		client := testAccProvider.Meta().(*connectivity.AliyunClient)

		vpcService := &VpcService{client: client}
		_, err := vpcService.DescribeNetworkAcl(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}
	}
	return nil
}

func testAccNetworkAcl_create(randInt int) string {
	return fmt.Sprintf(`
resource "alicloud_vpc" "default" {	
  cidr_block = "172.16.0.0/12"	
  name = "tf-testAccVpcConfig"
}	
resource "alicloud_network_acl" "default" {
  vpc_id = "${alicloud_vpc.default.id}"
  name = "tf-testAcc_network_acl%v.abc"
  description = "tf-testAcc_network_acl"
}
`, randInt)
}

func testAccNetworkAcl_name(randInt int) string {
	return fmt.Sprintf(`
resource "alicloud_vpc" "default" {	
  cidr_block = "172.16.0.0/12"	
  name = "tf-testAccVpcConfig"
}	
resource "alicloud_network_acl" "default" {
  vpc_id = "${alicloud_vpc.default.id}"
  name = "tf-testAcc_network_acl_modify%v.abc"
  description = "tf-testAcc_network_acl"
}
`, randInt)
}

func testAccNetworkAcl_description(randInt int) string {
	return fmt.Sprintf(`
resource "alicloud_vpc" "default" {	
  cidr_block = "172.16.0.0/12"	
  name = "tf-testAccVpcConfig"
}	
resource "alicloud_network_acl" "default" {
  vpc_id = "${alicloud_vpc.default.id}"
  name = "tf-testAcc_network_acl_modify%v.abc"
  description = "tf-testAcc_network_acl_modify"
}
`, randInt)
}

func testAccNetworkAcl_modify(randInt int) string {
	return fmt.Sprintf(`
resource "alicloud_vpc" "default" {	
  cidr_block = "172.16.0.0/12"	
  name = "tf-testAccVpcConfig"
}	
resource "alicloud_network_acl" "default" {
  vpc_id = "${alicloud_vpc.default.id}"
  name = "tf-testAcc_network_acl%v.abc"
  description = "tf-testAcc_network_acl"
}
`, randInt)
}

func testAccNetworkAcl_multi(randInt int) string {
	return fmt.Sprintf(`
resource "alicloud_vpc" "default" {	
  cidr_block = "172.16.0.0/12"	
  name = "tf-testAccVpcConfig%v.abc"
}	
resource "alicloud_network_acl" "default" {
  vpc_id = "${alicloud_vpc.default.id}"
  name = "tf-testAcc_network_acl"
  description = "tf-testAcc_network_acl${count.index}"
  count = 3
}
`, randInt)
}

var testAccNaclCheckMap = map[string]string{
	"description": "tf-testAcc_network_acl",
}
