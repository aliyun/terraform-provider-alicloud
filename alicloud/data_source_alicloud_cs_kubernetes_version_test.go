package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudCSKubernetesVersionDataSource(t *testing.T) {
	resourceId := "data.alicloud_cs_kubernetes_version.default"
	testAccCheck := resourceAttrInit(resourceId, map[string]string{}).resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSourceCSKubernetesVersion(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"metadata": CHECKSET,
					}),
				),
			},
		},
	})
}

func dataSourceCSKubernetesVersion() string {
	return fmt.Sprintf(`
# return all kubernetes version
data "alicloud_cs_kubernetes_version" "default" {
  cluster_type = "ManagedKubernetes"
  kubernetes_version = "1.22.3-aliyun.1"
  profile = "Default"
}
`)
}
