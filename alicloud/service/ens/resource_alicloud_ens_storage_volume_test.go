package ens_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ens"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TestAccAliCloudEnsStorageVolume_basic0 creates an ENS storage volume against a
// pre-existing storage gateway and storage medium. The storage gateway
// (sgw-***) and storage medium (d-***) dependencies are not yet modeled as
// Terraform resources, so the test reads their IDs and the ENS node from
// environment variables and skips when they are absent.
func TestAccAliCloudEnsStorageVolume_basic0(t *testing.T) {
	gatewayId := os.Getenv("ALICLOUD_ENS_STORAGE_GATEWAY_ID")
	storageId := os.Getenv("ALICLOUD_ENS_STORAGE_ID")
	ensRegionId := os.Getenv("ALICLOUD_ENS_REGION_ID")
	if gatewayId == "" || storageId == "" || ensRegionId == "" {
		t.Skip("Environment variables ALICLOUD_ENS_STORAGE_GATEWAY_ID, ALICLOUD_ENS_STORAGE_ID and ALICLOUD_ENS_REGION_ID must be set for ENS storage volume acceptance tests")
	}

	resourceId := "alicloud_ens_storage_volume.default"
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testaccensstoragevolume%d", rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckEnsStorageVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudEnsStorageVolumeBasic0(name, gatewayId, storageId, ensRegionId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEnsStorageVolumeExists(resourceId),
					resource.TestCheckResourceAttrSet(resourceId, "id"),
					resource.TestCheckResourceAttr(resourceId, "gateway_id", gatewayId),
					resource.TestCheckResourceAttr(resourceId, "storage_id", storageId),
					resource.TestCheckResourceAttr(resourceId, "ens_region_id", ensRegionId),
					resource.TestCheckResourceAttr(resourceId, "volume_name", name),
					resource.TestCheckResourceAttr(resourceId, "description", name),
					resource.TestCheckResourceAttr(resourceId, "is_enable", "true"),
					resource.TestCheckResourceAttr(resourceId, "is_auth", "true"),
					resource.TestCheckResourceAttr(resourceId, "auth_protocol", "CHAP"),
					resource.TestCheckResourceAttr(resourceId, "auth_user", name),
					resource.TestCheckResourceAttr(resourceId, "status", "available"),
					resource.TestCheckResourceAttrSet(resourceId, "creation_time"),
					resource.TestCheckResourceAttrSet(resourceId, "target_name"),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				// auth_user and auth_password are write-only: DescribeStorageVolume
				// never returns them, so import cannot verify them.
				ImportStateVerifyIgnore: []string{"auth_user", "auth_password"},
			},
		},
	})
}

func testAccCheckEnsStorageVolumeExists(resourceId string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceId]
		if !ok {
			return fmt.Errorf("resource %s not found in state", resourceId)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource %s has no ID set", resourceId)
		}
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)
		_, err := ens.NewService(client).DescribeStorageVolume(rs.Primary.ID)
		if err != nil {
			if errs.NotFoundError(err) {
				return fmt.Errorf("storage volume %s does not exist", rs.Primary.ID)
			}
			return fmt.Errorf("describing storage volume %s failed: %s", rs.Primary.ID, err)
		}
		return nil
	}
}

func testAccCheckEnsStorageVolumeDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ens_storage_volume" {
			continue
		}
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)
		_, err := ens.NewService(client).DescribeStorageVolume(rs.Primary.ID)
		if err != nil {
			if errs.NotFoundError(err) {
				continue
			}
			return fmt.Errorf("describing storage volume %s failed: %s", rs.Primary.ID, err)
		}
		return fmt.Errorf("storage volume %s still exists", rs.Primary.ID)
	}
	return nil
}

func testAccAliCloudEnsStorageVolumeBasic0(name, gatewayId, storageId, ensRegionId string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_ens_storage_volume" "default" {
  gateway_id    = "%s"
  storage_id    = "%s"
  ens_region_id = "%s"
  volume_name   = var.name
  description   = var.name
  is_enable     = true
  is_auth       = true
  auth_protocol = "CHAP"
  auth_user     = var.name
  auth_password = var.name
}
`, name, gatewayId, storageId, ensRegionId)
}
