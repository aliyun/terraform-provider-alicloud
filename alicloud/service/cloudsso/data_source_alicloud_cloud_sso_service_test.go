package cloudsso_test

import (
	"testing"

	acctest "github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudCloudSSOServiceDataSource(t *testing.T) {
	resourceId := "data.alicloud_cloud_sso_service.current"
	testAccCheck := acctest.ResourceAttrInit(resourceId, map[string]string{}).ResourceAttrMapUpdateSet()
	acctest.CheckoutSupportedRegions(t, true, connectivity.CloudSsoSupportRegions)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckEnterpriseAccountEnabled(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudCloudSsoServiceDataOnSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"id":     acctest.CHECKSET,
						"enable": "On",
						"status": "Opened",
					}),
				),
			},
			{
				Config: testAccCheckAlicloudCloudSsoServiceOffDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"id":     acctest.CHECKSET,
						"enable": "Off",
						"status": "",
					}),
				),
			},
		},
	})
}

const testAccCheckAlicloudCloudSsoServiceDataOnSource = `
data "alicloud_cloud_sso_service" "current" {
	enable = "On"
}
`

const testAccCheckAlicloudCloudSsoServiceOffDataSource = `
data "alicloud_cloud_sso_service" "current" {
	enable = "Off"
}
`
