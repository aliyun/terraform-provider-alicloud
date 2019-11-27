package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/kms"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"

	"testing"
)

func TestAccAlicloudKmsCiphertext_basic(t *testing.T) {
	var v *kms.EncryptResponse

	resourceId := "alicloud_kms_ciphertext.default"
	ra := resourceAttrInit(resourceId, kmsCiphertextBasicMap)

	serviceFunc := func() interface{} {
		return &KmsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccKmsCiphertextConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				ResourceName: resourceId,
			},
		},
	})
}

const testAccKmsCiphertextConfigBasic = `
resource "alicloud_kms_key" "default" {
	is_enabled  = true
}

resource "alicloud_kms_ciphertext" "default" {
	key_id = "${alicloud_kms_key.default.id}"
	plaintext = "plaintext"
}
`

var kmsCiphertextBasicMap = map[string]string{
	"plaintext":       CHECKSET,
	"key_id":          CHECKSET,
	"context":         NOSET,
	"ciphertext_blob": CHECKSET,
}
