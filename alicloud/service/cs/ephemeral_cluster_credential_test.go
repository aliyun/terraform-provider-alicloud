package cs_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// matches the RFC3339 UTC timestamps the CS API returns, e.g. 2029-10-07T11:14:54Z
var csClusterCredentialExpirationRegexp = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`)

func TestAccAliCloudCsClusterCredentialEphemeral_basic(t *testing.T) {
	rand := sdkacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccCsClusterCredentialEphemeral%d", rand)
	dataPath := tfjsonpath.New("data")
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_10_0),
		},
		CheckDestroy: acctest.CheckDestroyNoop,
		Steps: []resource.TestStep{
			{
				Config: hclCsClusterCredentialEphemeralClusterConfig(name),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("alicloud_cs_managed_kubernetes.default", tfjsonpath.New("id"), knownvalue.NotNull()),
				},
			},
			{
				Config: hclCsClusterCredentialEphemeralOpenConfig(name),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("echo.open", dataPath.AtMapKey("cluster_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("echo.open", dataPath.AtMapKey("temporary_duration_minutes"), knownvalue.Null()),
					statecheck.ExpectKnownValue("echo.open", dataPath.AtMapKey("kube_config"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("echo.open", dataPath.AtMapKey("expiration"), knownvalue.StringRegexp(csClusterCredentialExpirationRegexp)),
				},
			},
			{
				Config: hclCsClusterCredentialEphemeralTemporaryConfig(name),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("echo.open_temp", dataPath.AtMapKey("cluster_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("echo.open_temp", dataPath.AtMapKey("temporary_duration_minutes"), knownvalue.Int64Exact(60)),
					statecheck.ExpectKnownValue("echo.open_temp", dataPath.AtMapKey("kube_config"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("echo.open_temp", dataPath.AtMapKey("expiration"), knownvalue.StringRegexp(csClusterCredentialExpirationRegexp)),
				},
			},
		},
	})
}

func hclCsClusterCredentialEphemeralClusterConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "vpc" {
  vpc_name   = var.name
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "vswitch" {
  vpc_id       = alicloud_vpc.vpc.id
  cidr_block   = "192.168.1.0/24"
  zone_id      = data.alicloud_zones.default.zones.0.id
  vswitch_name = var.name
}

resource "alicloud_cs_managed_kubernetes" "default" {
  name                 = var.name
  cluster_spec         = "ack.pro.small"
  vswitch_ids          = [alicloud_vswitch.vswitch.id]
  new_nat_gateway      = true
  pod_cidr             = "10.93.0.0/16"
  service_cidr         = "172.21.0.0/16"
  slb_internet_enabled = true
}
`, name)
}

func hclCsClusterCredentialEphemeralOpenConfig(name string) string {
	return fmt.Sprintf(`
%s

ephemeral "alicloud_cs_cluster_credential" "default" {
  cluster_id = alicloud_cs_managed_kubernetes.default.id
}

provider "echo" {
  data = ephemeral.alicloud_cs_cluster_credential.default
}

resource "echo" "open" {}
`, hclCsClusterCredentialEphemeralClusterConfig(name))
}

func hclCsClusterCredentialEphemeralTemporaryConfig(name string) string {
	return fmt.Sprintf(`
%s

ephemeral "alicloud_cs_cluster_credential" "default" {
  cluster_id                 = alicloud_cs_managed_kubernetes.default.id
  temporary_duration_minutes = 60
}

provider "echo" {
  data = ephemeral.alicloud_cs_cluster_credential.default
}

resource "echo" "open_temp" {}
`, hclCsClusterCredentialEphemeralClusterConfig(name))
}
