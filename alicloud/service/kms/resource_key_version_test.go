package kms_test

import (
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/kms"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudKmsKeyVersion_basic(t *testing.T) {
	if os.Getenv("ALICLOUD_KMS_KEY_ID") == "" {
		t.Skip("Skipping as ALICLOUD_KMS_KEY_ID not set")
	}
	var v map[string]interface{}
	resourceId := "alicloud_kms_key_version.default"
	ra := acctest.ResourceAttrInit(resourceId, KmsKeyVersionMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return kms.NewKmsService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeKmsKeyVersion")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, "", resourceKMSKeyVersionConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			testAccPreCheckKMSForKeyIdImport(t)
		},

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             nil,
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
	"key_version_id": acctest.CHECKSET,
	"key_id":         acctest.CHECKSET,
}

//func testAccCheckKmsKeyVersionExists(n string, kv *kms.KeyVersion) resource.TestCheckFunc {
//	return func(s *terraform.State) error {
//		rs, ok := s.RootModule().Resources[n]
//		if !ok {
//			return errs.WrapError(fmt.Errorf("Not found: %s", n))
//		}
//
//		if rs.Primary.ID == "" {
//			return errs.WrapError(errs.Error("No Key Version ID is set"))
//		}
//
//		client := acctest.Provider.Meta().(*connectivity.AliyunClient)
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
//			return errs.WrapError(fmt.Errorf("Error finding key version %s", rs.Primary.ID))
//		}
//		return errs.WrapError(err)
//	}
//}

func resourceKMSKeyVersionConfigDependence(name string) string {
	return ""
}

// lintignore: R001
