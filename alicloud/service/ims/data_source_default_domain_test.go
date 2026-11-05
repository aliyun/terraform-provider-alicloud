package ims_test

import (
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudImsDefaultDomainDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudImsDefaultDomainDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.alicloud_ims_default_domain.default", "id"),
					resource.TestCheckResourceAttrSet("data.alicloud_ims_default_domain.default", "default_domain"),
				),
			},
		},
	})
}

func testAccAliCloudImsDefaultDomainDataSourceConfig() string {
	return `
data "alicloud_ims_default_domain" "default" {
}`
}
