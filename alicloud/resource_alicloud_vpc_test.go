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
	resource.AddTestSweepers("alicloud_vpc", &resource.Sweeper{
		Name: "alicloud_vpc",
		F:    testSweepVpcs,
		// When implemented, these should be removed firstly
		Dependencies: []string{
			"alicloud_vswitch",
			"alicloud_nat_gateway",
			"alicloud_security_group",
			"alicloud_ots_instance",
			"alicloud_router_interface",
			"alicloud_route_table",
			"alicloud_cen_instance",
			"alicloud_edas_cluster",
			"alicloud_edas_k8s_cluster",
		},
	})
}

func testSweepVpcs(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
		"tf_test_",
		"tf-test-",
		"testAcc",
	}

	var vpcs []vpc.Vpc
	request := vpc.CreateDescribeVpcsRequest()
	request.RegionId = client.RegionId
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	invoker := NewInvoker()
	for {
		var raw interface{}
		if err := invoker.Run(func() error {
			raw, err = client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
				return vpcClient.DescribeVpcs(request)
			})
			return err
		}); err != nil {
			log.Printf("[ERROR] Error retrieving VPCs: %s", WrapError(err))
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*vpc.DescribeVpcsResponse)
		if len(response.Vpcs.Vpc) < 1 {
			break
		}
		vpcs = append(vpcs, response.Vpcs.Vpc...)

		if len(response.Vpcs.Vpc) < PageSizeLarge {
			break
		}

		page, err := getNextpageNumber(request.PageNumber)
		if err != nil {
			log.Printf("[ERROR] %s", WrapError(err))
		}
		request.PageNumber = page
	}

	for _, v := range vpcs {
		name := v.VpcName
		id := v.VpcId
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping VPC: %s (%s)", name, id)
			continue
		}
		log.Printf("[INFO] Deleting VPC: %s (%s)", name, id)
		service := VpcService{client}
		err := service.sweepVpc(id)
		if err != nil {
			log.Printf("[ERROR] Failed to delete VPC (%s (%s)): %s", name, id, err)
		}
	}
	return nil
}

func testAccCheckVpcDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	vpcService := VpcService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_vpc" {
			continue
		}

		// Try to find the VPC
		instance, err := vpcService.DescribeVpc(rs.Primary.ID)

		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}
		return WrapError(fmt.Errorf("VPC %s still exist", instance.VpcId))
	}

	return nil
}

func TestAccAlicloudVpcBasic(t *testing.T) {
	var v vpc.Vpc
	rand := acctest.RandInt()
	resourceId := "alicloud_vpc.default"
	ra := resourceAttrInit(resourceId, testAccCheckVpcCheckMap)
	serviceFunc := func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckVpcDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVpcConfigBasic(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":              fmt.Sprintf("tf_testAccVpcConfigName%d", rand),
						"resource_group_id": CHECKSET,
						"router_id":         CHECKSET,
						"router_table_id":   CHECKSET,
						"route_table_id":    CHECKSET,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccCheckVpcConfig_name(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": fmt.Sprintf("tf_testAccVpcConfigName%d_change", rand),
					}),
				),
			},
			{
				Config: testAccCheckVpcConfig_description(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": fmt.Sprintf("tf_testAccVpcConfigName%d_decription", rand),
					}),
				),
			},
			{
				Config: testAccCheckVpcConfig_tags(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "acceptance test",
					}),
				),
			},
			{
				Config: testAccCheckVpcConfig_all(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":         fmt.Sprintf("tf_testAccVpcConfigName%d_all", rand),
						"description":  fmt.Sprintf("tf_testAccVpcConfigName%d_decription_all", rand),
						"tags.%":       REMOVEKEY,
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})

}

func TestAccAlicloudVpcMulti(t *testing.T) {
	var v vpc.Vpc
	rand := acctest.RandInt()
	resourceId := "alicloud_vpc.default.9"
	ra := resourceAttrInit(resourceId, testAccCheckVpcCheckMap)
	serviceFunc := func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckVpcDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVpcConfigMulti(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": fmt.Sprintf("tf_testAccVpcConfigName%d", rand),
					}),
				),
			},
		},
	})

}

func testAccCheckVpcConfigBasic(rand int) string {
	return fmt.Sprintf(
		`
variable "name" {
	default = "tf_testAccVpcConfigName%d"
}

resource "alicloud_vpc" "default" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}
`, rand)
}

func testAccCheckVpcConfig_name(rand int) string {
	return fmt.Sprintf(
		`
variable "name" {
	default = "tf_testAccVpcConfigName%d"
}

resource "alicloud_vpc" "default" {
	cidr_block = "172.16.0.0/12"
	name = "${var.name}_change"
}
`, rand)
}

func testAccCheckVpcConfig_description(rand int) string {
	return fmt.Sprintf(
		`
variable "name" {
	default = "tf_testAccVpcConfigName%d"
}

resource "alicloud_vpc" "default" {
	cidr_block = "172.16.0.0/12"
	name = "${var.name}_change"
	description = "${var.name}_decription"
}
`, rand)
}

func testAccCheckVpcConfig_tags(rand int) string {
	return fmt.Sprintf(
		`
variable "name" {
	default = "tf_testAccVpcConfigName%d"
}

resource "alicloud_vpc" "default" {
	cidr_block = "172.16.0.0/12"
	name = "${var.name}_change"
	description = "${var.name}_decription"
	tags 		= {
		Created = "TF"
		For 	= "acceptance test"
	}
}
`, rand)
}

func testAccCheckVpcConfig_all(rand int) string {
	return fmt.Sprintf(
		`
variable "name" {
	default = "tf_testAccVpcConfigName%d"
}

resource "alicloud_vpc" "default" {
	cidr_block = "172.16.0.0/12"
	name = "${var.name}_all"
	description = "${var.name}_decription_all"
}
`, rand)
}

func testAccCheckVpcConfigMulti(rand int) string {
	return fmt.Sprintf(
		`
variable "name" {
	default = "tf_testAccVpcConfigName%d"
}

resource "alicloud_vpc" "default" {
	name = "${var.name}"
	count = 10
	cidr_block = "172.16.0.0/12"
}
`, rand)
}

var testAccCheckVpcCheckMap = map[string]string{
	"cidr_block":        "172.16.0.0/12",
	"name":              "",
	"description":       "",
	"resource_group_id": CHECKSET,
	"router_id":         CHECKSET,
	"router_table_id":   CHECKSET,
	"route_table_id":    CHECKSET,
}
