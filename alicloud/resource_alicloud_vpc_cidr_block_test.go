package alicloud

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"log"
	"strings"
	"testing"
)

func init() {
	resource.AddTestSweepers("alicloud_vpc_cidr", &resource.Sweeper{
		Name: "alicloud_vpc_cidr",
		F:    testSweepVpcsCidr,
	})
}

func testSweepVpcsCidr(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testCidr",
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

func TestCidrAlicloudVpc(t *testing.T) {
	var v vpc.Vpc
	rand := acctest.RandInt()
	resourceId := "alicloud_vpc.default"
	ra := resourceAttrInit(resourceId, testCidrCheckVpcCheckMap)
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
		Steps: []resource.TestStep{
			{
				Config: testCidrCheckVpcConfig(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": fmt.Sprintf("tf_testCidrCheckVpcConfigName%d", rand),
					}),
				),
			},
		},
	})

}

func testCidrCheckVpcConfig(rand int) string {
	return fmt.Sprintf(
		`
variable "name" {
  default = "tf_testCidrCheckVpcConfigName%d"
}

resource "alicloud_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vpc_cidr_block" "default" {
  vpc_id = alicloud_vpc.default.id
  secondary_cidr_block = "10.0.0.0/8"
}
data "alicloud_vpcs" "vpcs" {
  ids = [alicloud_vpc.default.id]
  is_default = false
  name_regex = alicloud_vpc.default.name
}
`, rand)
}

var testCidrCheckVpcCheckMap = map[string]string{
	"cidr_block":        "172.16.0.0/12",
	"name":              "",
	"resource_group_id": CHECKSET,
	"router_id":         CHECKSET,
	"router_table_id":   CHECKSET,
	"route_table_id":    CHECKSET,
}
