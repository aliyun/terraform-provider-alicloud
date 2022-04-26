package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudCSKubernetesVersionDataSource(t *testing.T) {
	default_cluster_type := "Kubernetes"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSourceCSKubernetesVersion(default_cluster_type),
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

func dataSourceCSKubernetesVersion(clusterType string) string {
	return fmt.Sprintf(`
variable "cluster_Type" {
 default = "%s"
}
# return all kubernetes version
data "alicloud_cs_kubernetes_version" "default" {
 cluster_type= var.cluster_Type
}
`, clusterType)
}
