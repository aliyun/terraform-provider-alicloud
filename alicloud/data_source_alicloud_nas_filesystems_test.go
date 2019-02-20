package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloud_FileSystem_DataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudFileSystemsDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_nas_filesystems.fs"),
					resource.TestCheckResourceAttrSet("data.alicloud_nas_filesystems.fs", "filesystems.0.id"),
					resource.TestCheckResourceAttrSet("data.alicloud_nas_filesystems.fs", "filesystems.0.region_id"),
					resource.TestCheckResourceAttr("data.alicloud_nas_filesystems.fs", "filesystems.0.protocol_type", "NFS"),
					resource.TestCheckResourceAttr("data.alicloud_nas_filesystems.fs", "filesystems.0.storage_type", "Performance"),
					resource.TestCheckResourceAttrSet("data.alicloud_nas_filesystems.fs", "filesystems.0.metered_size"),
					resource.TestCheckResourceAttrSet("data.alicloud_nas_filesystems.fs", "filesystems.0.mounttarget_domain"),
					resource.TestCheckResourceAttrSet("data.alicloud_nas_filesystems.fs", "filesystems.0.create_time"),
				),
			},
		},
	})
}

const testAccCheckAlicloudFileSystemsDataSource = `
variable "description" {
  default = "tf-testAccCheckAlicloudFileSystemsDataSource"
}
resource "alicloud_nas_filesystem" "foo" {
  description = "${var.description}"
  storage_type = "Performance"
  protocol_type = "NFS"
}
data "alicloud_nas_filesystems" "fs" {
  storage_type = "Performance"
  protocol_type = "NFS"
}
`
