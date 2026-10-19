package cms_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/cms"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TestAccAliCloudCmsContact_basic exercises the Framework-managed
// alicloud_cms_contact resource end to end against the real CMS v2024-03-30
// RESTful API via the muxed provider: create, read, update (full replace of
// im_user_ids and clear of email/phone), and import.
func TestAccAliCloudCmsContact_basic(t *testing.T) {
	r := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccmscontact%d", r)
	resourceId := "alicloud_cms_contact.default"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckCmsContactDestroy(resourceId),
		Steps: []resource.TestStep{
			{
				Config: testAccCmsContactConfig(name, "tf-test"+fmt.Sprintf("%d@example.com", r), fmt.Sprintf("86-138%08d", r), "zh_CN", `"key1" = "user1"`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "contact_name", name),
					resource.TestCheckResourceAttr(resourceId, "email", fmt.Sprintf("tf-test%d@example.com", r)),
					resource.TestCheckResourceAttr(resourceId, "phone", fmt.Sprintf("86-138%08d", r)),
					resource.TestCheckResourceAttr(resourceId, "lang", "zh_CN"),
					resource.TestCheckResourceAttr(resourceId, "im_user_ids.%", "1"),
					resource.TestCheckResourceAttr(resourceId, "im_user_ids.key1", "user1"),
					resource.TestCheckResourceAttrSet(resourceId, "contact_id"),
					resource.TestCheckResourceAttrSet(resourceId, "region_id"),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				// im_user_ids is write-only with respect to the CMS ListContacts
				// API: CreateContact/UpdateContact accept imUserIds in the body,
				// but ListContacts (the only read path) does not echo a populated
				// imUserIds for the contact (confirmed against the live API: the
				// Framework Read returns null, preserved across refresh only via
				// the contactState plan/state fallback). A fresh `terraform import`
				// starts with no prior im_user_ids, so the field cannot be
				// reconstructed and is ignored during ImportStateVerify. This
				// mirrors the effective legacy SDKv2 behavior, which never set
				// im_user_ids on Read when the API omitted it.
				ImportStateVerifyIgnore: []string{"im_user_ids"},
			},
			{
				Config: testAccCmsContactConfig(name+"-updated", fmt.Sprintf("tf-test2%d@example.com", r), fmt.Sprintf("86-139%08d", r), "en_US", `"key2" = "user2"`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "contact_name", name+"-updated"),
					resource.TestCheckResourceAttr(resourceId, "email", fmt.Sprintf("tf-test2%d@example.com", r)),
					resource.TestCheckResourceAttr(resourceId, "phone", fmt.Sprintf("86-139%08d", r)),
					resource.TestCheckResourceAttr(resourceId, "lang", "en_US"),
					resource.TestCheckResourceAttr(resourceId, "im_user_ids.%", "1"),
					resource.TestCheckResourceAttr(resourceId, "im_user_ids.key2", "user2"),
					// CMS PutContact fully replaces im_user_ids (key1 -> key2),
					// so key1 legitimately disappears from state on update.
					resource.TestCheckNoResourceAttr(resourceId, "im_user_ids.key1"),
				),
			},
			{
				Config: testAccCmsContactConfig(name+"-clear", "", "", "zh_CN", ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "contact_name", name+"-clear"),
					resource.TestCheckResourceAttr(resourceId, "lang", "zh_CN"),
					// CMS PutContact PATCH persists "" for email/phone when the
					// config clears them, so state is present-empty.
					resource.TestCheckResourceAttr(resourceId, "email", ""),
					resource.TestCheckResourceAttr(resourceId, "phone", ""),
				),
			},
		},
	})
}

func testAccCmsContactConfig(name, email, phone, lang, imUserIds string) string {
	imBlock := ""
	if imUserIds != "" {
		imBlock = fmt.Sprintf(`  im_user_ids = {
    %s
  }`, imUserIds)
	}
	return fmt.Sprintf(`
resource "alicloud_cms_contact" "default" {
  contact_name = "%[1]s"
  email        = "%[2]s"
  phone        = "%[3]s"
  lang         = "%[4]s"
%[5]s
}
`, name, email, phone, lang, imBlock)
}

func testAccCheckCmsContactDestroy(resourceId string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)
		svc := cms.NewService(client)
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "alicloud_cms_contact" {
				continue
			}
			contact, err := svc.GetContact(context.Background(), rs.Primary.ID)
			if err != nil {
				return err
			}
			if contact != nil {
				return fmt.Errorf("CMS Contact %s still exists", rs.Primary.ID)
			}
		}
		return nil
	}
}
