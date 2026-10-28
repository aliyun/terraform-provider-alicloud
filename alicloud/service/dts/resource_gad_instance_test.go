// Package dts. This file is generated automatically. Please do not modify it manually, thank you!
package dts_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

const dtsGadInstanceResourceId = "alicloud_dts_gad_instance.default"

// TestAccAliCloudDtsGadInstance_basic7553 covers the full lifecycle of the GAD instance: create
// (CreateGadInstance), in-place rename (ModifyGadInstanceName), resource-group replacement
// (resource_group_id is ForceNew) and import.
//
// The GAD instance and its master RDS are both created in the provider region (cn-hangzhou). The
// generator emitted a second byte-identical case (basic7560); it is folded into this one to avoid
// provisioning a second RDS+GAD pair per run without losing any coverage.
//
// CreateGadInstance's RegionId (the GAD instance region) and MasterDbInstanceRegion are independent
// request parameters and the API declares no cross-region constraint (its errorCodes cover only
// throttling, expired security tokens and service unavailability), so a same-region master+GAD pair
// is accepted and exercises the same Create/Read/Update/Delete/Import paths a cross-region pair would.
// A cross-region fixture was attempted first (mirroring dtsCheckJobDependence with a cn-beijing
// provider alias), but the framework half of the muxed provider reads its region back out of the
// single shared alicloud.Provider() instance (see acctest.ProtoV5ProviderFactories and
// provider/factory.go), so ConfigureProvider for the default and the alias instances overwrite one
// shared meta and strand the alias's VPC/vswitch in the wrong region (Forbidden.VpcNotFound). No
// framework service-subpackage test uses a provider alias today; cross-region alias support under
// ProtoV5 would need a per-instance factory and is tracked separately rather than blocking this
// resource.
//
// The generator's hardcoded master_db_instance_id (rm-2ze3772m0g292g33g, an API-doc example value
// absent from the test account, causing CreateGadInstance -> RdsGetInstanceAttr.InvokeFail /
// InvalidDBInstanceId.NotFound) is replaced by a dynamically created RDS, mirroring the
// dtsCheckJobDependence fixture pattern. resource_group_id is covered end-to-end with two dedicated
// resource_manager_resource_group resources: the instance is created inside one resource group and,
// because resource_group_id is ForceNew, recreated inside the other on change. GAD has no in-place
// resource-group relocation API (Dts/2020-01-01 GAD surface: Create/Delete/ModifyName/Describe/Detach
// only); ConvertInstanceResourceGroup — the relocation API alicloud_dts_instance uses for regular DTS
// instances (dts-*) — returns Forbidden.InstanceNotFound for GAD instance IDs (gad-*), so a change
// forces replacement rather than an in-place call. ModifyGadInstanceName is reserved for the rename
// step only; its ResourceGroupId query parameter is an authorization scope, not a relocation target
// (documented as "特殊业务专用字段,无需传参"), so rename never sends ResourceGroupId and cannot drift
// resource_group_id.
func TestAccAliCloudDtsGadInstance_basic7553(t *testing.T) {
	name := fmt.Sprintf("tf-test-dts-gad-%d", tfacctest.RandIntRange(10000, 99999))

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			t.Setenv("ALICLOUD_REGION", "cn-hangzhou")
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				// Create the GAD instance inside a dedicated resource group so the create path carries
				// ResourceGroupId to CreateGadInstance, then assert the refresh keeps it.
				Config: dtsGadInstanceDependence(name) + dtsGadInstanceBasicBody(name, "alicloud_resource_manager_resource_group.default.id"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dtsGadInstanceResourceId, "gad_instance_name", name),
					resource.TestCheckResourceAttr(dtsGadInstanceResourceId, "db_engine_type", "MySQL"),
					resource.TestCheckResourceAttr(dtsGadInstanceResourceId, "gad_instance_type", "DR"),
					resource.TestCheckResourceAttr(dtsGadInstanceResourceId, "master_db_instance_region", "cn-hangzhou"),
					resource.TestCheckResourceAttrPair(dtsGadInstanceResourceId, "master_db_instance_id", "alicloud_db_instance.master", "id"),
					resource.TestCheckResourceAttrPair(dtsGadInstanceResourceId, "master_db_zone_id", "data.alicloud_db_zones.master", "zones.0.id"),
					resource.TestCheckResourceAttrSet(dtsGadInstanceResourceId, "id"),
					resource.TestCheckResourceAttrPair(dtsGadInstanceResourceId, "resource_group_id", "alicloud_resource_manager_resource_group.default", "id"),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(dtsGadInstanceResourceId, tfjsonpath.New("create_time"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(dtsGadInstanceResourceId, tfjsonpath.New("region_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(dtsGadInstanceResourceId, tfjsonpath.New("status"), knownvalue.NotNull()),
				},
			},
			{
				// A rename goes through ModifyGadInstanceName and must be visible after the refresh;
				// the resource group must stay put (rename must not relocate the instance).
				Config: dtsGadInstanceDependence(name) + dtsGadInstanceBasicBody(name+"-update", "alicloud_resource_manager_resource_group.default.id"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dtsGadInstanceResourceId, "gad_instance_name", name+"-update"),
					resource.TestCheckResourceAttr(dtsGadInstanceResourceId, "db_engine_type", "MySQL"),
					resource.TestCheckResourceAttr(dtsGadInstanceResourceId, "gad_instance_type", "DR"),
					resource.TestCheckResourceAttr(dtsGadInstanceResourceId, "master_db_instance_region", "cn-hangzhou"),
					resource.TestCheckResourceAttrPair(dtsGadInstanceResourceId, "master_db_instance_id", "alicloud_db_instance.master", "id"),
					resource.TestCheckResourceAttrPair(dtsGadInstanceResourceId, "master_db_zone_id", "data.alicloud_db_zones.master", "zones.0.id"),
					resource.TestCheckResourceAttrPair(dtsGadInstanceResourceId, "resource_group_id", "alicloud_resource_manager_resource_group.default", "id"),
				),
			},
			{
				// resource_group_id is ForceNew (V1): changing it forces replacement, so the instance is
				// recreated inside the new resource group and the refresh must confirm the new id. GAD has
				// no in-place relocation API; ConvertInstanceResourceGroup (alicloud_dts_instance's path
				// for dts-* ids) returns Forbidden.InstanceNotFound for gad-* ids.
				Config: dtsGadInstanceDependence(name) + dtsGadInstanceBasicBody(name+"-update", "alicloud_resource_manager_resource_group.change.id"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dtsGadInstanceResourceId, "gad_instance_name", name+"-update"),
					resource.TestCheckResourceAttrPair(dtsGadInstanceResourceId, "resource_group_id", "alicloud_resource_manager_resource_group.change", "id"),
				),
			},
			{
				ResourceName:      dtsGadInstanceResourceId,
				ImportState:       true,
				ImportStateVerify: true,
				// master_db_zone_id is write-only (it is sent to CreateGadInstance but never returned
				// by DescribeGadInstanceDetail), so an imported state cannot carry it.
				ImportStateVerifyIgnore: []string{"master_db_zone_id"},
			},
		},
	})
}

// dtsGadInstanceDependence provisions the master RDS the GAD instance references. The master RDS is
// created in the provider region (cn-hangzhou) by the default provider so the GAD instance (also
// created by the default provider in cn-hangzhou) can attach to it as its DR primary. The fixture
// mirrors dtsCheckJobDependence: PostPaid HighAvailability MySQL 8.0 on cloud_essd, with db_zones
// and db_instance_classes data sources resolving the zone/class and a 60m create timeout for a busy
// account. Unlike check_job, GAD rejects general-purpose/standard RDS specs
// (DTS.Msg.GADUnsupportedRdsStandard), so a dedicated (mysql.x*) class is selected from the zone's
// available classes instead of taking the first (general-purpose) one; storage is set to 100 GB.
func dtsGadInstanceDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%[1]s"
}

data "alicloud_db_zones" "master" {
  engine                   = "MySQL"
  engine_version           = "8.0"
  instance_charge_type     = "PostPaid"
  category                 = "HighAvailability"
  db_instance_storage_type = "cloud_essd"
}

data "alicloud_db_instance_classes" "master" {
  zone_id                  = data.alicloud_db_zones.master.zones.0.id
  engine                   = "MySQL"
  engine_version           = "8.0"
  instance_charge_type     = "PostPaid"
  category                 = "HighAvailability"
  db_instance_storage_type = "cloud_essd"
}

locals {
  # GAD rejects general-purpose/standard RDS specs (DTS.Msg.GADUnsupportedRdsStandard on a 4c/8g
  # mysql.n2.* class); the dedicated GAD role must use a dedicated (mysql.x*) class. The data source
  # returns general-purpose classes first, so the first returned class is replaced by the first
  # dedicated class available in this zone.
  gad_master_class = [
    for c in data.alicloud_db_instance_classes.master.instance_classes : c.instance_class
    if startswith(c.instance_class, "mysql.x")
  ][0]
}

resource "alicloud_vpc" "master" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "master" {
  vpc_id       = alicloud_vpc.master.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = data.alicloud_db_zones.master.zones.0.id
  vswitch_name = var.name
}

resource "alicloud_db_instance" "master" {
  engine                   = "MySQL"
  engine_version           = "8.0"
  instance_type            = local.gad_master_class
  instance_storage         = "100"
  db_instance_storage_type = "cloud_essd"
  vswitch_id               = alicloud_vswitch.master.id
  instance_name            = var.name

  timeouts {
    create = "60m"
    update = "60m"
  }
}

resource "alicloud_resource_manager_resource_group" "default" {
  display_name        = "tf-test-dts-gad-default"
  resource_group_name = "%[1]s-rg"
}

resource "alicloud_resource_manager_resource_group" "change" {
  display_name        = "tf-test-dts-gad-change"
  resource_group_name = "%[1]s-rgc"
}
`, name)
}

func dtsGadInstanceBasicBody(gadName, resourceGroupIdExpr string) string {
	return fmt.Sprintf(`
resource "alicloud_dts_gad_instance" "default" {
  gad_instance_name          = "%[1]s"
  db_engine_type            = "MySQL"
  gad_instance_type         = "DR"
  master_db_zone_id          = data.alicloud_db_zones.master.zones.0.id
  master_db_instance_region = "cn-hangzhou"
  master_db_instance_id     = alicloud_db_instance.master.id
  resource_group_id         = %[2]s
}
`, gadName, resourceGroupIdExpr)
}
