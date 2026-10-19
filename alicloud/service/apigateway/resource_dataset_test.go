package apigateway_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccAliCloudApiGatewayDataset_basic0 covers a JWT_BLOCKING dataset: create
// with two tags, then update by renaming, changing description, swapping the tag
// set (remove For, change Created, add Env), then import. Tags round-trip through
// TagResources/UntagResources/ListTagResources.
func TestAccAliCloudApiGatewayDataset_basic0(t *testing.T) {
	r := tfacctest.RandIntRange(10000, 99999)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccApiGatewayDatasetConfig_basic0(r, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_api_gateway_dataset.default", "dataset_name", fmt.Sprintf("tf-testacc-ds-%d", r)),
					resource.TestCheckResourceAttr("alicloud_api_gateway_dataset.default", "dataset_type", "JWT_BLOCKING"),
					resource.TestCheckResourceAttr("alicloud_api_gateway_dataset.default", "description", "tf test acc dataset 0"),
					resource.TestCheckResourceAttr("alicloud_api_gateway_dataset.default", "tags.%", "2"),
					resource.TestCheckResourceAttr("alicloud_api_gateway_dataset.default", "tags.Created", "tfTestAcc0"),
					resource.TestCheckResourceAttr("alicloud_api_gateway_dataset.default", "tags.For", "tf testacc 0"),
					resource.TestCheckResourceAttrSet("alicloud_api_gateway_dataset.default", "create_time"),
					resource.TestCheckResourceAttrSet("alicloud_api_gateway_dataset.default", "modified_time"),
				),
			},
			{
				Config: testAccApiGatewayDatasetConfig_basic0(r, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_api_gateway_dataset.default", "dataset_name", fmt.Sprintf("tf-testacc-ds-%d-u", r)),
					resource.TestCheckResourceAttr("alicloud_api_gateway_dataset.default", "description", "tf test acc dataset 0 updated"),
					resource.TestCheckResourceAttr("alicloud_api_gateway_dataset.default", "tags.%", "2"),
					resource.TestCheckResourceAttr("alicloud_api_gateway_dataset.default", "tags.Created", "tfTestAcc0u"),
					resource.TestCheckResourceAttr("alicloud_api_gateway_dataset.default", "tags.Env", "tf testacc 0 updated"),
					resource.TestCheckNoResourceAttr("alicloud_api_gateway_dataset.default", "tags.For"),
				),
			},
			{
				ResourceName:      "alicloud_api_gateway_dataset.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccAliCloudApiGatewayDataset_basic1 covers an IP_WHITELIST_CIDR dataset
// created with no tags, then updated to add one tag, then import.
func TestAccAliCloudApiGatewayDataset_basic1(t *testing.T) {
	r := tfacctest.RandIntRange(10000, 99999)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccApiGatewayDatasetConfig_basic1(r, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_api_gateway_dataset.default", "dataset_name", fmt.Sprintf("tf-testacc-ds-%d", r)),
					resource.TestCheckResourceAttr("alicloud_api_gateway_dataset.default", "dataset_type", "IP_WHITELIST_CIDR"),
					resource.TestCheckResourceAttr("alicloud_api_gateway_dataset.default", "description", "tf test acc dataset cidr"),
				),
			},
			{
				Config: testAccApiGatewayDatasetConfig_basic1(r, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_api_gateway_dataset.default", "dataset_name", fmt.Sprintf("tf-testacc-ds-%d-u", r)),
					resource.TestCheckResourceAttr("alicloud_api_gateway_dataset.default", "description", "tf test acc dataset cidr updated"),
					resource.TestCheckResourceAttr("alicloud_api_gateway_dataset.default", "tags.%", "1"),
					resource.TestCheckResourceAttr("alicloud_api_gateway_dataset.default", "tags.Owner", "tfTestAcc1"),
				),
			},
			{
				ResourceName:      "alicloud_api_gateway_dataset.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccApiGatewayDatasetConfig_basic0(r int, updated bool) string {
	if !updated {
		return fmt.Sprintf(`
resource "alicloud_api_gateway_dataset" "default" {
  dataset_name = "tf-testacc-ds-%[1]d"
  dataset_type = "JWT_BLOCKING"
  description  = "tf test acc dataset 0"
  tags = {
    Created = "tfTestAcc0"
    For     = "tf testacc 0"
  }
}
`, r)
	}
	return fmt.Sprintf(`
resource "alicloud_api_gateway_dataset" "default" {
  dataset_name = "tf-testacc-ds-%[1]d-u"
  dataset_type = "JWT_BLOCKING"
  description  = "tf test acc dataset 0 updated"
  tags = {
    Created = "tfTestAcc0u"
    Env     = "tf testacc 0 updated"
  }
}
`, r)
}

func testAccApiGatewayDatasetConfig_basic1(r int, updated bool) string {
	if !updated {
		return fmt.Sprintf(`
resource "alicloud_api_gateway_dataset" "default" {
  dataset_name = "tf-testacc-ds-%[1]d"
  dataset_type = "IP_WHITELIST_CIDR"
  description  = "tf test acc dataset cidr"
}
`, r)
	}
	return fmt.Sprintf(`
resource "alicloud_api_gateway_dataset" "default" {
  dataset_name = "tf-testacc-ds-%[1]d-u"
  dataset_type = "IP_WHITELIST_CIDR"
  description  = "tf test acc dataset cidr updated"
  tags = {
    Owner = "tfTestAcc1"
  }
}
`, r)
}
