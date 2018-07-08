package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudCSKubernetes_import(t *testing.T) {
	resourceName := "alicloud_cs_kubernetes.k8s"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccContainerKubernetes_basic,
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"name_prefix", "new_nat_gateway", "pod_cidr",
					"service_cidr", "enable_ssh", "password", "install_cloud_monitor"},
			},
		},
	})
}
