package alicloud

import (
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudKMSKeyVersion_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_kms_key_version.default"
	ra := resourceAttrInit(resourceId, KmsKeyVersionMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KmsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKmsKeyVersion")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, "", resourceKMSKeyVersionConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckKMSForKeyIdImport(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"key_id": os.Getenv("ALICLOUD_KMS_KEY_ID"),
				}),
				Check: resource.ComposeTestCheckFunc(
					//testAccCheckKmsKeyVersionExists(resourceId, &l),
					testAccCheck(KmsKeyVersionMap),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var KmsKeyVersionMap = map[string]string{
	"key_version_id": CHECKSET,
	"key_id":         CHECKSET,
}

//func testAccCheckKmsKeyVersionExists(n string, kv *kms.KeyVersion) resource.TestCheckFunc {
//	return func(s *terraform.State) error {
//		rs, ok := s.RootModule().Resources[n]
//		if !ok {
//			return WrapError(fmt.Errorf("Not found: %s", n))
//		}
//
//		if rs.Primary.ID == "" {
//			return WrapError(Error("No Key Version ID is set"))
//		}
//
//		client := testAccProvider.Meta().(*connectivity.AliyunClient)
//
//		request := kms.CreateListKeyVersionsRequest()
//		request.KeyId = rs.Primary.Attributes["key_id"]
//
//		raw, err := client.WithKmsClient(func(kmsClient *kms.Client) (interface{}, error) {
//			return kmsClient.ListKeyVersions(request)
//		})
//
//		if err == nil {
//			response, _ := raw.(*kms.ListKeyVersionsResponse)
//			if len(response.KeyVersions.KeyVersion) > 0 {
//				for _, v := range response.KeyVersions.KeyVersion {
//					if v.KeyVersionId == strings.Split(rs.Primary.ID, ":")[1] {
//						*kv = v
//						return nil
//					}
//				}
//			}
//			return WrapError(fmt.Errorf("Error finding key version %s", rs.Primary.ID))
//		}
//		return WrapError(err)
//	}
//}

func resourceKMSKeyVersionConfigDependence(name string) string {
	return ""
}
