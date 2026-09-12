package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAliCloudCmsContact_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_contact.default"
	ra := resourceAttrInit(resourceId, AliCloudCmsContactMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsContact")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccmscontact%d", rand)
	email := fmt.Sprintf("tf-test%d@example.com", rand)
	phone := fmt.Sprintf("86-138%08d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCmsContactBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"contact_name": name,
					"email":        email,
					"phone":        phone,
					"lang":         "zh_CN",
					"im_user_ids": map[string]interface{}{
						"key1": "user1",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"contact_name":     name,
						"email":            email,
						"phone":            phone,
						"lang":             "zh_CN",
						"im_user_ids.%":    "1",
						"im_user_ids.key1": "user1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"contact_name": name + "-updated",
					"email":        fmt.Sprintf("tf-test2%d@example.com", rand),
					"phone":        fmt.Sprintf("86-139%08d", rand),
					"lang":         "en_US",
					"im_user_ids": map[string]interface{}{
						"key2": "user2",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"contact_name":     name + "-updated",
						"email":            fmt.Sprintf("tf-test2%d@example.com", rand),
						"phone":            fmt.Sprintf("86-139%08d", rand),
						"lang":             "en_US",
						"im_user_ids.%":    "1",
						"im_user_ids.key2": "user2",
						// CMS PutContact update fully replaces im_user_ids (key1->key2),
						// so key1 legitimately disappears from state. resourceAttrMapUpdateSet
						// merges cumulatively, so the key1 asserted in step 0 would otherwise
						// carry into this step and fail TestCheckResourceAttr. NOSET remaps it
						// to TestCheckNoResourceAttr, asserting the stale key is gone.
						"im_user_ids.key1": NOSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"contact_name": name + "-clear",
					"email":        "",
					"phone":        "",
					"lang":         "zh_CN",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"contact_name": name + "-clear",
						"lang":         "zh_CN",
						// Step3 clear config sets email/phone to "". The CMS API
						// (PutContact PATCH) persists "" for both fields, so state is
						// present-empty. resourceAttrMapUpdateSet merges cumulatively,
						// so the Step2 email/phone assertions would otherwise carry
						// forward and fail TestCheckResourceAttr. Asserting "" here
						// overwrites the stale carried value and matches state.
						"email": "",
						"phone": "",
					}),
				),
			},
		},
	})
}

func testAccCheckAlicloudCmsContactDestroy(resourceId string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		cmsServiceV2 := CmsServiceV2{client}
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "alicloud_cms_contact" {
				continue
			}
			_, err := cmsServiceV2.DescribeCmsContact(rs.Primary.ID)
			if err != nil {
				if NotFoundError(err) {
					continue
				}
				return WrapError(err)
			}
			return WrapError(Error("CMS Contact %s still exists", rs.Primary.ID))
		}
		return nil
	}
}

// AliCloudCmsContactMap is the attribute map for the CMS Contact resource.
// im_user_ids is a TypeMap: the bare map name is not set by the SDK (it
// stores .%/.<key>), so CHECKSET on the bare name is unreliable. The per-key
// assertions (im_user_ids.% / im_user_ids.key1) in each step already cover it.
// workspace is a real @readonly field in CloudSpec (Contact_get/ListContacts
// mapping responsePath "$.contacts[*].workspace"), but ListContacts only
// returns the key when the contact was created with a workspace value; a basic
// contact leaves it absent, so it is not asserted here.
var AliCloudCmsContactMap = map[string]string{
	"contact_name": CHECKSET,
	"email":        CHECKSET,
	"phone":        CHECKSET,
	"lang":         CHECKSET,
	"contact_id":   CHECKSET,
	"region_id":    CHECKSET,
}

// AliCloudCmsContactBasicDependence returns the dependence config (none for Contact).
func AliCloudCmsContactBasicDependence(name string) string {
	return ""
}
