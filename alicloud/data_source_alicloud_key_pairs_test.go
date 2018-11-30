package alicloud

import (
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudKeyPairsDataSource_basic(t *testing.T) {
	var keypair ecs.KeyPair

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKeyPairDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudKeyPairsDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeyPairExists("alicloud_key_pair.basic", &keypair),
					testAccCheckAlicloudDataSourceID("data.alicloud_key_pairs.name_regex"),
					resource.TestCheckResourceAttrSet("data.alicloud_key_pairs.name_regex", "key_pairs.0.id"),
					resource.TestCheckResourceAttr("data.alicloud_key_pairs.name_regex", "key_pairs.0.key_name", "tf-testAcc-key-pair-datasource"),
					resource.TestCheckResourceAttr("data.alicloud_key_pairs.name_regex", "key_pairs.0.instances.#", "0"),
				),
			},
		},
	})
}

func TestAccAlicloudKeyPairsDataSource_empty(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKeyPairDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudKeyPairsDataSourceEmpty,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_key_pairs.name_regex"),
					resource.TestCheckNoResourceAttr("data.alicloud_key_pairs.name_regex", "key_pairs.0.id"),
					resource.TestCheckNoResourceAttr("data.alicloud_key_pairs.name_regex", "key_pairs.0.key_name"),
					resource.TestCheckNoResourceAttr("data.alicloud_key_pairs.name_regex", "key_pairs.0.instances.#"),
				),
			},
		},
	})
}

const testAccCheckAlicloudKeyPairsDataSourceBasic = `
resource "alicloud_key_pair" "basic" {
	key_name = "tf-testAcc-key-pair-datasource"
}
data "alicloud_key_pairs" "name_regex" {
	name_regex = "${alicloud_key_pair.basic.id}"
}
`

const testAccCheckAlicloudKeyPairsDataSourceEmpty = `
data "alicloud_key_pairs" "name_regex" {
	name_regex = "^tf-testacc-fake-name"
}
`
