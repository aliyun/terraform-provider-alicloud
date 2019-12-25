package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAliCloudCopyImageBasic(t *testing.T) {
	var v ecs.Image

	resourceId := "alicloud_copy_image.default"
	// multi provideris
	var providers []*schema.Provider
	providerFactories := map[string]terraform.ResourceProviderFactory{
		"alicloud": func() (terraform.ResourceProvider, error) {
			p := Provider()
			providers = append(providers, p.(*schema.Provider))
			return p, nil
		},
	}
	ra := resourceAttrInit(resourceId, testAccCopyImageCheckMap)
	rand := acctest.RandIntRange(1000, 9999)
	testAccCheck := ra.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEcsCopyImageConfigBasic%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCopyImageBasicConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckImageDestroyWithProviders(&providers),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"provider":         "alicloud.sh",
					"source_image_id":  alicloud_image.default.id,
					"source_region_id": "cn-hangzhou",
					"description":      fmt.Sprintf("tf-testAccEcsImageConfigBasic%ddescription", rand),
					"name":             name,
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test123",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImageExistsWithProviders(resourceId, &v, &providers),
					testAccCheck(map[string]string{
						"name":         name,
						"description":  fmt.Sprintf("tf-testAccEcsImageConfigBasic%ddescription", rand),
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "acceptance test123",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": fmt.Sprintf("tf-testAccEcsImageConfigBasic%ddescriptionChange", rand),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImageExistsWithProviders(resourceId, &v, &providers),
					testAccCheck(map[string]string{
						"description": fmt.Sprintf("tf-testAccEcsImageConfigBasic%ddescriptionChange", rand),
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": fmt.Sprintf("tf-testAccEcsImageConfigBasic%dchange", rand),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImageExistsWithProviders(resourceId, &v, &providers),
					testAccCheck(map[string]string{
						"name": fmt.Sprintf("tf-testAccEcsImageConfigBasic%dchange", rand),
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF1",
						"For":     "acceptance test1232",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImageExistsWithProviders(resourceId, &v, &providers),
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF1",
						"tags.For":     "acceptance test1232",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": fmt.Sprintf("tf-testAccEcsImageConfigBasic%ddescription", rand),
					"name":        name,
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test123",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImageExistsWithProviders(resourceId, &v, &providers),
					testAccCheck(map[string]string{
						"description":  fmt.Sprintf("tf-testAccEcsImageConfigBasic%ddescription", rand),
						"name":         name,
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "acceptance test123",
					}),
				),
			},
		},
	})
}

func testAccCheckImageExistsWithProviders(n string, image *ecs.Image, providers *[]*schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No image  ID is set")
		}
		for _, provider := range *providers {
			if provider.Meta() == nil {
				continue
			}

			client := provider.Meta().(*connectivity.AliyunClient)
			ecsService := EcsService{client}

			resp, err := ecsService.DescribeImageById(rs.Primary.ID)
			if err != nil {
				if NotFoundError(err) {
					continue
				}
				return err
			}

			*image = resp
			return nil
		}
		return fmt.Errorf("image not found")
	}
}

func testAccCheckImageDestroyWithProviders(providers *[]*schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, provider := range *providers {
			if provider.Meta() == nil {
				continue
			}
			if err := testAccCheckImageDestroyWithProvider(s, provider); err != nil {
				return err
			}
		}
		return nil
	}
}

func testAccCheckImageDestroyWithProvider(s *terraform.State, provider *schema.Provider) error {

	client := provider.Meta().(*connectivity.AliyunClient)
	ecsService := EcsService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_copy_image" {
			continue
		}

		resp, err := ecsService.DescribeImageById(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		} else {
			return fmt.Errorf("image still exist,  ID %s ", resp.ImageId)
		}
	}

	return nil
}

var testAccCopyImageCheckMap = map[string]string{}

func resourceCopyImageBasicConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
provider "alicloud" {
  alias = "sh"
  region = "cn-shanghai"
}
provider "alicloud" {
  alias = "hz"
  region = "cn-hangzhou"
}
data "alicloud_instance_types" "default" {
    provider = "alicloud.hz"
 	cpu_core_count    = 1
	memory_size       = 2
}
data "alicloud_images" "default" {
  provider = "alicloud.hz"
  name_regex  = "^ubuntu_18.*_64"
  owners      = "system"
}
resource "alicloud_vpc" "default" {
  provider = "alicloud.hz"
  name       = var.name
  cidr_block = "172.16.0.0/16"
}
resource "alicloud_vswitch" "default" {
  provider = "alicloud.hz"
  vpc_id            = alicloud_vpc.default.id
  cidr_block        = "172.16.0.0/24"
  availability_zone = data.alicloud_instance_types.default.instance_types.0.availability_zones.0
  name              = var.name
}
resource "alicloud_security_group" "default" {
  provider = "alicloud.hz"
  name   = var.name
  vpc_id = alicloud_vpc.default.id
}
resource "alicloud_instance" "default" {
  provider = "alicloud.hz"
  image_id = "${data.alicloud_images.default.ids[0]}"
  instance_type = "${data.alicloud_instance_types.default.ids[0]}"
  security_groups = "${[alicloud_security_group.default.id]}"
  vswitch_id = alicloud_vswitch.default.id
  instance_name = var.name
}
resource "alicloud_image" "default" {
  provider = "alicloud.hz"
  instance_id = alicloud_instance.default.id
  name        = var.name
}
`, name)
}
