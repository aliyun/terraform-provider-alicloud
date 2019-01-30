package alicloud

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudCSManagedKubernetes_import(t *testing.T) {
	resourceName := "alicloud_cs_managed_kubernetes.k8s"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckWithRegions(t, true, connectivity.ManagedKubernetesSupportedRegions) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckManagedKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccManagedKubernetes_basic,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"name_prefix", "new_nat_gateway", "pod_cidr",
					"service_cidr", "password", "install_cloud_monitor"},
			},
		},
	})
}
