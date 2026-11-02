package cloudsso_test

import (
	"fmt"
	"testing"

	acctest "github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/cloudsso"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudCloudSSOSCIMServerCredential_basic0(t *testing.T) {
	var v map[string]interface{}
	acctest.CheckoutSupportedRegions(t, true, connectivity.CloudSsoSupportRegions)
	resourceId := "alicloud_cloud_sso_scim_server_credential.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudCloudSSOSCIMServerCredentialMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudsso.NewCloudssoService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCloudSsoScimServerCredential")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scloudssoscimservercredential%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudCloudSSOSCIMServerCredentialBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"directory_id": "${alicloud_cloud_sso_directory.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"directory_id": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Disabled",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Disabled",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Enabled",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Enabled",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"credential_secret_file"},
			},
		},
	})
}

func TestAccAliCloudCloudSSOSCIMServerCredential_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	acctest.CheckoutSupportedRegions(t, true, connectivity.CloudSsoSupportRegions)
	resourceId := "alicloud_cloud_sso_scim_server_credential.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudCloudSSOSCIMServerCredentialMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudsso.NewCloudssoService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCloudSsoScimServerCredential")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scloudssoscimservercredential%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudCloudSSOSCIMServerCredentialBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"directory_id":           "${alicloud_cloud_sso_directory.default.id}",
					"status":                 "Enabled",
					"credential_secret_file": "./credential_secret_file.txt",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"directory_id": acctest.CHECKSET,
						"status":       "Enabled",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"credential_secret_file"},
			},
		},
	})
}

func AliCloudCloudSSOSCIMServerCredentialBasicDependence0(name string) string {
	return fmt.Sprintf(` 
	variable "name" {
  		default = "%s"
	}

	resource "alicloud_cloud_sso_directory" "default" {
	  directory_name    = var.name
	}
`, name)
}

var AliCloudCloudSSOSCIMServerCredentialMap0 = map[string]string{
	"status":          acctest.CHECKSET,
	"credential_id":   acctest.CHECKSET,
	"credential_type": acctest.CHECKSET,
	"create_time":     acctest.CHECKSET,
	"expire_time":     acctest.CHECKSET,
}
