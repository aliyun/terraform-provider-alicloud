package nas_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/nas"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudNASSmbAclAttachment_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nas_smb_acl_attachment.default"
	acctest.CheckoutSupportedRegions(t, true, connectivity.NASSupportRegions)
	ra := acctest.ResourceAttrInit(resourceId, AlicloudNASSmbAclMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return nas.NewNasService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeNasSmbAcl")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%ssmbacl%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudNASSmbAclBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"file_system_id": "${alicloud_nas_file_system.default.id}",
					"keytab":         "BQIAAABHAAIADUFMSUFEVEVTVC5DT00ABGNpZnMAGXNtYnNlcnZlcjI0LmFsaWFkdGVzdC5jb20AAAABAAAAAAEAAQAIqIx6v7p11oUAAABHAAIADUFMSUFEVEVTVC5DT00ABGNpZnMAGXNtYnNlcnZlcjI0LmFsaWFkdGVzdC5jb20AAAABAAAAAAEAAwAIqIx6v7p11oUAAABPAAIADUFMSUFEVEVTVC5DT00ABGNpZnMAGXNtYnNlcnZlcjI0LmFsaWFkdGVzdC5jb20AAAABAAAAAAEAFwAQnQZWB3RAPHU7PMIJyBWePAAAAF8AAgANQUxJQURURVNULkNPTQAEY2lmcwAZc21ic2VydmVyMjQuYWxpYWR0ZXN0LmNvbQAAAAEAAAAAAQASACAGJ7F0s+bcBjf6jD5HlvlRLmPSOW+qDZe0Qk0lQcf8WwAAAE8AAgANQUxJQURURVNULkNPTQAEY2lmcwAZc21ic2VydmVyMjQuYWxpYWR0ZXN0LmNvbQAAAAEAAAAAAQARABDdFmanrSIatnDDhoOXYadj",
					"keytab_md5":     "E3CCF7E2416DF04FA958AA4513EA29E8",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"file_system_id":  acctest.CHECKSET,
						"auth_method":     "Kerberos",
						"enabled":         "true",
						"super_admin_sid": "",
						"home_dir_path":   "",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"keytab":     "BQIAAABZAAIAC1RFU1ROQVMuQ09NAARjaWZzAC0wMjUyYzRiZmFlLXR4YTc5LmNuLWhhbmd6aG91Lm5hcy5hbGl5dW5jcy5jb20AAAABAAAAAAMAAQAIvJLQj+BefzEAAABZAAIAC1RFU1ROQVMuQ09NAARjaWZzAC0wMjUyYzRiZmFlLXR4YTc5LmNuLWhhbmd6aG91Lm5hcy5hbGl5dW5jcy5jb20AAAABAAAAAAMAAwAIvJLQj+BefzEAAABhAAIAC1RFU1ROQVMuQ09NAARjaWZzAC0wMjUyYzRiZmFlLXR4YTc5LmNuLWhhbmd6aG91Lm5hcy5hbGl5dW5jcy5jb20AAAABAAAAAAMAFwAQmnGFGg9wrGIckBGRDz3ZEAAAAHEAAgALVEVTVE5BUy5DT00ABGNpZnMALTAyNTJjNGJmYWUtdHhhNzkuY24taGFuZ3pob3UubmFzLmFsaXl1bmNzLmNvbQAAAAEAAAAAAwASACB8NytHr3V2sFcEmV0p/k+7gxk4nbRvRotJ19NTmP1GigAAAGEAAgALVEVTVE5BUy5DT00ABGNpZnMALTAyNTJjNGJmYWUtdHhhNzkuY24taGFuZ3pob3UubmFzLmFsaXl1bmNzLmNvbQAAAAEAAAAAAwARABBz8lgg+jF54cqR2WZ7nltq",
					"keytab_md5": "91224BCF548214F4DF3F7CA1898354CE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"keytab_md5": acctest.CHECKSET,
						"keytab":     acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_anonymous_access": true,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_anonymous_access": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"encrypt_data": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"encrypt_data": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"reject_unencrypted_access": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"reject_unencrypted_access": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"super_admin_sid": "S-1-1-1-2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"super_admin_sid": "S-1-1-1-2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"home_dir_path": "/wang",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"home_dir_path": "/wang",
					}),
				),
			},
		},
	})
}

var AlicloudNASSmbAclMap0 = map[string]string{
	"file_system_id":  acctest.CHECKSET,
	"auth_method":     acctest.CHECKSET,
	"enabled":         acctest.CHECKSET,
	"super_admin_sid": acctest.CHECKSET,
	"home_dir_path":   acctest.CHECKSET,
}

func AlicloudNASSmbAclBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_nas_zones" "default" {
  file_system_type = "standard"
}

locals {
  count_size = length(data.alicloud_nas_zones.default.zones)
  zone_id    = data.alicloud_nas_zones.default.zones[local.count_size - 1].zone_id
}

resource "alicloud_nas_file_system" "default" {
  protocol_type = "SMB"
  storage_type  = "Capacity"
  description   =  var.name
  zone_id          = local.zone_id
  encrypt_type     = "0"
  file_system_type = "standard"
}
`, name)
}

// lintignore: R001
