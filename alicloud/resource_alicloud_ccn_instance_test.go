package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/smartag"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudCcnInstance_basic(t *testing.T) {
	var ccn smartag.CloudConnectNetwork
	resourceId := "alicloud_ccn_instance.default"
	ra := resourceAttrInit(resourceId, nil)
	serviceFunc := func() interface{} {
		return &SagService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &ccn, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	rand := acctest.RandIntRange(1000000, 9999999)
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
				Config: testAccCcnInstanceConfigBasic(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":        fmt.Sprintf("tf-testAcc%sCcnConfig-%d", defaultRegionToTest, rand),
						"description": "tf-testAccCcnConfigDescription",
						"cidr_block":  "192.168.0.0/24",
						"is_default":  "true",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccCcnInstanceConfigBasicUpdateName(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": fmt.Sprintf("tf-testAcc%sCcnConfig-%d-New", defaultRegionToTest, rand),
					}),
				),
			},
			{
				Config: testAccCcnInstanceConfigBasicUpdateDescription(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "tf-testAccCcnConfigDescription-New",
					}),
				),
			},
			{
				Config: testAccCcnInstanceConfigBasicUpdateCidrblock(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cidr_block": "192.168.1.0/24,192.168.2.0/24",
					}),
				),
			},
			{
				Config: testAccCcnInstanceConfigBasicUpdateConfig(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":        fmt.Sprintf("tf-testAcc%sCcnConfig-%d-New", defaultRegionToTest, rand),
						"description": "tf-testAccCcnConfigDescription-New",
						"cidr_block":  "192.168.1.0/24,192.168.2.0/24",
						"is_default":  "true",
					}),
				),
			},
			{
				Config: testAccCcnInstanceConfigBasicUpdateGrant(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":        fmt.Sprintf("tf-testAcc%sCcnConfig-%d-New", defaultRegionToTest, rand),
						"description": "tf-testAccCcnConfigDescription-New",
						"cidr_block":  "192.168.1.0/24,192.168.2.0/24",
						"is_default":  "true",
						"cen_id":      "cen-4vdgx1tyhjisjyjyy2",
						"cen_uid":     "1688401595963306",
						"total_count": "1",
					}),
				),
			},
			{
				Config: testAccCcnInstanceConfigBasicUpdateRevoke(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":        fmt.Sprintf("tf-testAcc%sCcnConfig-%d-New", defaultRegionToTest, rand),
						"description": "tf-testAccCcnConfigDescription-New",
						"cidr_block":  "192.168.1.0/24,192.168.2.0/24",
						"is_default":  "true",
						"total_count": "0",
					}),
				),
			},
		},
	})
}

func testAccCcnInstanceConfigBasic(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_ccn_instance" "default" {
		name = "tf-testAcc%sCcnConfig-%d"
		description = "tf-testAccCcnConfigDescription"
		cidr_block = "192.168.0.0/24"
		is_default = true
}
`, defaultRegionToTest, rand)
}
func testAccCcnInstanceConfigBasicUpdateName(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_ccn_instance" "default" {
		name = "tf-testAcc%sCcnConfig-%d-New"
		description = "tf-testAccCcnConfigDescription"
		cidr_block = "192.168.0.0/24"
		is_default = true
}
`, defaultRegionToTest, rand)
}
func testAccCcnInstanceConfigBasicUpdateDescription(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_ccn_instance" "default" {
		name = "tf-testAcc%sCcnConfig-%d-New"
		description = "tf-testAccCcnConfigDescription-New"
		cidr_block = "192.168.0.0/24"
		is_default = true
}
`, defaultRegionToTest, rand)
}
func testAccCcnInstanceConfigBasicUpdateCidrblock(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_ccn_instance" "default" {
		name = "tf-testAcc%sCcnConfig-%d-New"
		description = "tf-testAccCcnConfigDescription-New"
		cidr_block = "192.168.1.0/24,192.168.2.0/24"
		is_default = true
}
`, defaultRegionToTest, rand)
}
func testAccCcnInstanceConfigBasicUpdateConfig(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_ccn_instance" "default" {
		name = "tf-testAcc%sCcnConfig-%d-New"
		description = "tf-testAccCcnConfigDescription-New"
		cidr_block = "192.168.1.0/24,192.168.2.0/24"
		is_default = true
}
`, defaultRegionToTest, rand)
}
func testAccCcnInstanceConfigBasicUpdateGrant(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_ccn_instance" "default" {
		name = "tf-testAcc%sCcnConfig-%d-New"
		description = "tf-testAccCcnConfigDescription-New"
		cidr_block = "192.168.1.0/24,192.168.2.0/24"
		is_default = true
		cen_id = "cen-4vdgx1tyhjisjyjyy2"
		cen_uid = "1688401595963306"
		total_count = "1"
}
`, defaultRegionToTest, rand)
}
func testAccCcnInstanceConfigBasicUpdateRevoke(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_ccn_instance" "default" {
		name = "tf-testAcc%sCcnConfig-%d-New"
		description = "tf-testAccCcnConfigDescription-New"
		cidr_block = "192.168.1.0/24,192.168.2.0/24"
		is_default = true
		cen_id = "cen-4vdgx1tyhjisjyjyy2"
		cen_uid = "1688401595963306"
		total_count = "0"
}
`, defaultRegionToTest, rand)
}
