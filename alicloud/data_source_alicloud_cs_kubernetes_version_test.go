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
				Config: dataSourceCSManagedKubernetesVersion(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"metadata.0.version":           CHECKSET,
						"metadata.0.runtime.0.name":    CHECKSET,
						"metadata.0.runtime.0.version": CHECKSET,
					}),
				),
			},
			{
				Config: dataSourceCSKubernetesVersion(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"metadata.0.version":           CHECKSET,
						"metadata.0.runtime.0.name":    CHECKSET,
						"metadata.0.runtime.0.version": CHECKSET,
					}),
				),
			},
			{
				Config: dataSourceCSServerlessKubernetesVersion(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"metadata.0.version":           CHECKSET,
						"metadata.0.runtime.0.name":    CHECKSET,
						"metadata.0.runtime.0.version": CHECKSET,
					}),
				),
			},
			{
				Config: dataSourceCSEdgeKubernetesVersion(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"metadata.0.version":           CHECKSET,
						"metadata.0.runtime.0.name":    CHECKSET,
						"metadata.0.runtime.0.version": CHECKSET,
					}),
				),
			},
		},
	})
}

func dataSourceCSManagedKubernetesVersion() string {
	return fmt.Sprintf(`
# managed kubernetes version
data "alicloud_cs_kubernetes_version" "default" {
  cluster_type = "ManagedKubernetes"
  profile = "Default"
}
`)
}

func dataSourceCSKubernetesVersion() string {
	return fmt.Sprintf(`
# kubernetes version
data "alicloud_cs_kubernetes_version" "default" {
  cluster_type = "Kubernetes"
  profile = "Default"
}
`)
}

func dataSourceCSServerlessKubernetesVersion() string {
	return fmt.Sprintf(`
# serverless kubernetes version
data "alicloud_cs_kubernetes_version" "default" {
  cluster_type = "ManagedKubernetes"
  profile = "Serverless"
}
`)
}

func dataSourceCSEdgeKubernetesVersion() string {
	return fmt.Sprintf(`
# edge kubernetes version
data "alicloud_cs_kubernetes_version" "default" {
  cluster_type = "ManagedKubernetes"
  profile = "Edge"
}
`)
}
