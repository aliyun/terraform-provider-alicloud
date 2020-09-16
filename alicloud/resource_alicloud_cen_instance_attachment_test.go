package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAlicloudCenInstanceAttachment_basic(t *testing.T) {
	var v *cbn.DescribeCenAttachedChildInstanceAttributeResponse
	resourceId := "alicloud_cen_instance_attachment.default"
	var providers []*schema.Provider
	providerFactories := map[string]terraform.ResourceProviderFactory{
		"alicloud": func() (terraform.ResourceProvider, error) {
			p := Provider()
			providers = append(providers, p.(*schema.Provider))
			return p, nil
		},
	}
	ra := resourceAttrInit(resourceId, cenInstanceAttachmentMap)
	rand := acctest.RandIntRange(1000000, 9999999)
	testAccCheck := ra.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.CenNoSkipRegions)
		},

		// module name
		IDRefreshName:     resourceId,
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckCenInstanceAttachmentDestroyWithProviders(&providers),

		Steps: []resource.TestStep{
			{
				Config: testAccCenInstanceAttachmentBasic(rand, defaultRegionToTest),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCenInstanceAttachmentExistsWithProviders(resourceId, v, &providers),
					testAccCheck(nil),
				),
			},
		},
	})
}
func TestAccAlicloudCenInstanceAttachment_multi_same_region(t *testing.T) {
	var v *cbn.DescribeCenAttachedChildInstanceAttributeResponse
	resourceId := "alicloud_cen_instance_attachment.default"
	var providers []*schema.Provider
	providerFactories := map[string]terraform.ResourceProviderFactory{
		"alicloud": func() (terraform.ResourceProvider, error) {
			p := Provider()
			providers = append(providers, p.(*schema.Provider))
			return p, nil
		},
	}
	ra := resourceAttrInit(resourceId, cenInstanceAttachmentMap)
	rand := acctest.RandIntRange(1000000, 9999999)
	testAccCheck := ra.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.CenNoSkipRegions)
		},

		// module name
		IDRefreshName:     resourceId,
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckCenInstanceAttachmentDestroyWithProviders(&providers),

		Steps: []resource.TestStep{
			{
				Config: testAccCenInstanceAttachmentMultiSameRegion(rand, defaultRegionToTest),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCenInstanceAttachmentExistsWithProviders(resourceId, v, &providers),
					testAccCheck(nil),
				),
			},
		},
	})
}

func TestAccAlicloudCenInstanceAttachment_multi_different_region(t *testing.T) {
	var v *cbn.DescribeCenAttachedChildInstanceAttributeResponse
	resourceId := "alicloud_cen_instance_attachment.default"
	var providers []*schema.Provider
	providerFactories := map[string]terraform.ResourceProviderFactory{
		"alicloud": func() (terraform.ResourceProvider, error) {
			p := Provider()
			providers = append(providers, p.(*schema.Provider))
			return p, nil
		},
	}
	ra := resourceAttrInit(resourceId, cenInstanceAttachmentMap)
	rand := acctest.RandIntRange(1000000, 9999999)
	testAccCheck := ra.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName:     resourceId,
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckCenInstanceAttachmentDestroyWithProviders(&providers),

		Steps: []resource.TestStep{
			{
				Config: testAccCenInstanceAttachmentMultiDifferentRegion(rand, defaultRegionToTest),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCenInstanceAttachmentExistsWithProviders(resourceId, v, &providers),
					testAccCheck(nil),
				),
			},
		},
	})
}

var cenInstanceAttachmentMap = map[string]string{
	"instance_id":              CHECKSET,
	"child_instance_id":        CHECKSET,
	"child_instance_region_id": CHECKSET,
}

func testAccCenInstanceAttachmentBasic(rand int, region string) string {
	return fmt.Sprintf(`
	variable "name"{
	    default = "tf-testAcc%sCenInstanceAttachmentBasic-%d"
	}

	resource "alicloud_cen_instance" "default" {
	    name = "${var.name}"
	    description = "tf-testAccCenInstanceAttachmentBasicDescription"
	}

	resource "alicloud_vpc" "default" {
	    name = "${var.name}"
	    cidr_block = "192.168.0.0/16"
	}

	resource "alicloud_cen_instance_attachment" "default" {
	    instance_id = "${alicloud_cen_instance.default.id}"
	    child_instance_id = "${alicloud_vpc.default.id}"
	    child_instance_type = "VPC"
	    child_instance_region_id = "%s"
	}
	`, region, rand, region)
}
func testAccCenInstanceAttachmentMultiSameRegion(rand int, region string) string {
	return fmt.Sprintf(`
	variable "name"{
	    default = "tf-testAcc%sCenInstanceAttachmentBasic-%d"
	}

	resource "alicloud_cen_instance" "default" {
	    name = "${var.name}"
	    description = "tf-testAccCenInstanceAttachmentBasicDescription"
	}

	resource "alicloud_vpc" "default" {
	    name = "${var.name}"
	    cidr_block = "192.168.0.0/16"
	}

	resource "alicloud_vpc" "default1" {
	    name = "${var.name}"
	    cidr_block = "172.16.0.0/12"
	}

	resource "alicloud_cen_instance_attachment" "default" {
	    instance_id = "${alicloud_cen_instance.default.id}"
	    child_instance_id = "${alicloud_vpc.default.id}"
	    child_instance_type = "VPC"
	    child_instance_region_id = "%s"
	}

	resource "alicloud_cen_instance_attachment" "default1" {
	    instance_id = "${alicloud_cen_instance.default.id}"
	    child_instance_id = "${alicloud_vpc.default1.id}"
	    child_instance_type = "VPC"
	    child_instance_region_id = "%s"
	}
	`, region, rand, region, region)
}

func testAccCenInstanceAttachmentMultiDifferentRegion(rand int, region string) string {
	return fmt.Sprintf(
		`
variable "name"{
    default = "tf-testAccCen%sInstanceAttachmentMultiDifferentRegions-%d"
}

provider "alicloud" {
    alias = "fra"
    region = "eu-central-1"
}

provider "alicloud" {
    alias = "sh"
    region = "cn-shanghai"
}

resource "alicloud_vpc" "default" {
    provider = "alicloud.fra"
    name = "${var.name}"
    cidr_block = "192.168.0.0/16"
}

resource "alicloud_vpc" "default1" {
    provider = "alicloud.sh"
    name = "${var.name}"
    cidr_block = "172.16.0.0/12"
}

resource "alicloud_cen_instance" "default" {
    name = "${var.name}"
    description = "tf-testAccCenInstanceAttachmentMultiDifferentRegionsDescription"
}

resource "alicloud_cen_instance_attachment" "default" {
    instance_id = "${alicloud_cen_instance.default.id}"
    child_instance_id = "${alicloud_vpc.default.id}"
	child_instance_type = "VPC"
    child_instance_region_id = "eu-central-1"
}

resource "alicloud_cen_instance_attachment" "default1" {
    instance_id = "${alicloud_cen_instance.default.id}"
    child_instance_id = "${alicloud_vpc.default1.id}"
	child_instance_type = "VPC"
    child_instance_region_id = "cn-shanghai"
}
`, region, rand)

}

func testAccCheckCenInstanceAttachmentExistsWithProviders(n string, instance *cbn.DescribeCenAttachedChildInstanceAttributeResponse, providers *[]*schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Cen Child Instance ID is set")
		}
		for _, provider := range *providers {
			// Ignore if Meta is empty, this can happen for validation providers
			if provider.Meta() == nil {
				continue
			}

			client := provider.Meta().(*connectivity.AliyunClient)
			cbnService := CbnService{client}
			childInstance, err := cbnService.DescribeCenInstanceAttachment(rs.Primary.ID)
			if err != nil {
				return err
			}

			if childInstance.Status != "Attached" {
				return fmt.Errorf("CEN id %s instance id %s status error", childInstance.CenId, childInstance.ChildInstanceId)
			}

			instance = &childInstance
			return nil
		}
		return fmt.Errorf("Cen Child Instance not found")
	}
}

func testAccCheckCenInstanceAttachmentDestroyWithProviders(providers *[]*schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, provider := range *providers {
			if provider.Meta() == nil {
				continue
			}
			if err := testAccCheckCenInstanceAttachmentDestroyWithProvider(s, provider); err != nil {
				return err
			}
		}
		return nil
	}
}

func testAccCheckCenInstanceAttachmentDestroyWithProvider(s *terraform.State, provider *schema.Provider) error {
	client := provider.Meta().(*connectivity.AliyunClient)
	cbnService := CbnService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_instance_attachment" {
			continue
		}

		instance, err := cbnService.DescribeCenInstanceAttachment(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		} else {
			return fmt.Errorf("CEN %s child instance %s still attach", instance.CenId, instance.ChildInstanceId)
		}
	}

	return nil
}
