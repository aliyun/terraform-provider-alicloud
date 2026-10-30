package cs_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudCSKubernetesVersionDataSource(t *testing.T) {
	resourceId := "data.alicloud_cs_kubernetes_version.default"
	testAccCheck := acctest.ResourceAttrInit(resourceId, map[string]string{}).ResourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: dataSourceCSManagedKubernetesVersion(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_type":                 acctest.CHECKSET,
						"kubernetes_version":           "",
						"profile":                      acctest.CHECKSET,
						"metadata.0.version":           acctest.CHECKSET,
						"metadata.0.runtime.0.name":    acctest.CHECKSET,
						"metadata.0.runtime.0.version": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: dataSourceCSKubernetesVersion(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_type":                 acctest.REGEXMATCH + ".+",
						"kubernetes_version":           acctest.REGEXMATCH + ".*",
						"profile":                      acctest.REGEXMATCH + ".+",
						"metadata.0.version":           acctest.CHECKSET,
						"metadata.0.runtime.0.name":    acctest.CHECKSET,
						"metadata.0.runtime.0.version": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: dataSourceCSServerlessKubernetesVersion(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"metadata.0.version":           acctest.CHECKSET,
						"metadata.0.runtime.0.name":    acctest.CHECKSET,
						"metadata.0.runtime.0.version": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: dataSourceCSEdgeKubernetesVersion(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"metadata.0.version":           acctest.CHECKSET,
						"metadata.0.runtime.0.name":    acctest.CHECKSET,
						"metadata.0.runtime.0.version": acctest.CHECKSET,
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
