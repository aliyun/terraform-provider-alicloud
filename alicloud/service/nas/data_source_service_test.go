package nas_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	terraform "github.com/hashicorp/terraform-plugin-testing/terraform"
)

// testAccCheckAlicloudDataSourceID 是 alicloud/provider_test.go 中同名根包测试 helper 的本地副本，
// 供 nas_test 外部测试包在不引用根包的情况下校验数据源 ID 已设置。
func testAccCheckAlicloudDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("data source ID not set")
		}
		return nil
	}
}

func TestAccAliCloudNASServiceDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudNasServiceDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_nas_service.current"),
					resource.TestCheckResourceAttrSet("data.alicloud_nas_service.current", "id"),
					resource.TestCheckResourceAttr("data.alicloud_nas_service.current", "status", "Opened"),
				),
			},
		},
	})
}

const testAccCheckAlicloudNasServiceDataSource = `
data "alicloud_nas_service" "current" {
	enable = "On"
}
`
