package polardb_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/polardb"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestUnitAIDBClusterDatasetSchema(t *testing.T) {
	r := polardb.ResourceAIDBClusterDataset()
	for _, required := range []string{"db_cluster_id", "dataset_name"} {
		s, ok := r.Schema[required]
		if !ok {
			t.Fatalf("schema is missing required attribute %q", required)
		}
		if !s.Required || !s.ForceNew {
			t.Fatalf("attribute %q must be Required and ForceNew", required)
		}
	}
	if r.Create == nil || r.Read == nil || r.Delete == nil {
		t.Fatal("resource must implement Create, Read and Delete")
	}
	if r.Importer == nil || r.Importer.State == nil {
		t.Fatal("resource must implement ImportState")
	}
}

func TestUnitAIDBClusterDatasetImportState(t *testing.T) {
	r := polardb.ResourceAIDBClusterDataset()
	d := r.TestResourceData()
	d.SetId("pc-bp1234567890:ds-bp1234567890")
	results, err := r.Importer.State(d, nil)
	if err != nil {
		t.Fatalf("unexpected import error: %s", err)
	}
	if len(results) != 1 || results[0].Id() != "ds-bp1234567890" {
		t.Fatalf("unexpected import result: %#v", results)
	}
	if got := results[0].Get("db_cluster_id").(string); got != "pc-bp1234567890" {
		t.Fatalf("unexpected db_cluster_id after import: %q", got)
	}

	bad := r.TestResourceData()
	bad.SetId("only-one-part")
	if _, err := r.Importer.State(bad, nil); err == nil {
		t.Fatal("expected an error for an import ID without the cluster prefix")
	}
}

func TestAccAliCloudPolarDBAIDBClusterDataset_basic0(t *testing.T) {
	rand := tfacctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testaccPolardbAIDBDataset%d", rand)
	resourceId := "alicloud_polardb_aidb_cluster_dataset.default"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPolarDBAIDBClusterDatasetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPolarDBAIDBClusterDatasetConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "dataset_name", name+"-ds"),
					resource.TestCheckResourceAttrSet(resourceId, "db_cluster_id"),
					resource.TestCheckResourceAttrSet(resourceId, "creation_time"),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccPolarDBAIDBClusterDatasetImportStateIdFunc(resourceId),
				ImportStateVerifyIgnore: []string{
					"dataset_name",
					"db_cluster_id",
					"import_mode",
					"relative_db_cluster_id",
					"region_id",
				},
			},
		},
	})
}

func testAccPolarDBAIDBClusterDatasetImportStateIdFunc(resourceId string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceId]
		if !ok {
			return "", fmt.Errorf("resource %s is not found in state", resourceId)
		}
		return fmt.Sprintf("%s:%s", rs.Primary.Attributes["db_cluster_id"], rs.Primary.ID), nil
	}
}

func testAccCheckPolarDBAIDBClusterDatasetDestroy(s *terraform.State) error {
	client := acctest.Provider.Meta().(*connectivity.AliyunClient)
	service := polardb.NewService(client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_polardb_aidb_cluster_dataset" {
			continue
		}
		_, err := service.DescribeAIDBClusterDataset(rs.Primary.Attributes["db_cluster_id"], rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("alicloud_polardb_aidb_cluster_dataset %s still exists", rs.Primary.ID)
		}
		if !errs.NotFoundError(err) {
			return fmt.Errorf("unexpected error while checking alicloud_polardb_aidb_cluster_dataset %s: %s", rs.Primary.ID, err)
		}
	}
	return nil
}

func testAccPolarDBAIDBClusterDatasetConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_polardb_node_classes" "this" {
  db_type    = "MySQL"
  db_version = "8.0"
  pay_type   = "PostPaid"
  category   = "Normal"
}

resource "alicloud_vpc" "default" {
  vpc_name = var.name
}

resource "alicloud_vswitch" "default" {
  zone_id    = data.alicloud_polardb_node_classes.this.classes.1.zone_id
  vpc_id     = alicloud_vpc.default.id
  cidr_block = cidrsubnet(alicloud_vpc.default.cidr_block, 8, 4)
}

resource "alicloud_polardb_cluster" "cluster" {
  db_type       = "MySQL"
  db_version    = "8.0"
  pay_type      = "PostPaid"
  db_node_count = "2"
  db_node_class = "polar.mysql.x4.large"
  vswitch_id    = alicloud_vswitch.default.id
  description   = var.name
}

resource "alicloud_security_group" "default" {
  name   = var.name
  vpc_id = alicloud_vpc.default.id
}

resource "alicloud_polardb_aicluster" "default" {
  region_id              = "cn-beijing"
  zone_id                = data.alicloud_polardb_node_classes.this.classes.1.zone_id
  db_node_class          = "polar.mysql.g8.4xlarge.gu50"
  db_cluster_description = var.name
  pay_type               = "Postpaid"
  vswitch_id             = alicloud_vswitch.default.id
  vpc_id                 = alicloud_vpc.default.id
  kube_type              = "ainode"
  model_name             = "Qwen3.5-9B"
  extension              = "maas"
  inference_engine       = "sglang"
  db_cluster_id          = alicloud_polardb_cluster.cluster.id
  security_group_id      = alicloud_security_group.default.id
}

resource "alicloud_polardb_aidb_cluster_dataset" "default" {
  db_cluster_id = alicloud_polardb_aicluster.default.id
  dataset_name  = "${var.name}-ds"
}
`, name)
}
