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
		CheckDestroy: testAccCheckContainerClusterDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccContainerKubernetes_basic,
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"name_prefix", "new_nat_gateway", "master_instance_type", "worker_instance_type",
					"pod_cidr", "service_cidr", "enable_ssh", "password", "master_disk_size", "master_disk_category",
					"worker_disk_size", "worker_disk_category", "install_cloud_monitor"},
			},
		},
	})
}
