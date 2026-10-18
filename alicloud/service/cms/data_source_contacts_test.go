package cms_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccAliCloudCmsContactsDataSource_basic exercises the Framework-managed
// alicloud_cms_contacts data source against the real CMS ListContacts API: it
// creates a contact, then reads the data source filtered by contact_name via
// interpolation (which creates the implicit dependency edge that delays the
// data source Read until apply, after the resource exists), and asserts the
// created contact is returned.
func TestAccAliCloudCmsContactsDataSource_basic(t *testing.T) {
	r := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccmscontactds%d", r)
	email := fmt.Sprintf("tf-ds%d@example.com", r)
	phone := fmt.Sprintf("86-137%08d", r)
	resourceId := "data.alicloud_cms_contacts.default"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckCmsContactDestroy("alicloud_cms_contact.default"),
		Steps: []resource.TestStep{
			{
				Config: testAccCmsContactsDataSourceConfig(name, email, phone, ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceId, "contacts.#"),
					resource.TestCheckResourceAttr(resourceId, "contacts.0.contact_name", name),
					resource.TestCheckResourceAttr(resourceId, "contacts.0.email", email),
					resource.TestCheckResourceAttr(resourceId, "contacts.0.phone", phone),
				),
			},
			{
				Config: testAccCmsContactsDataSourceConfig(name, email, phone, "/tmp/contact_output.txt"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceId, "contacts.#"),
				),
			},
		},
	})
}

// testAccCmsContactsDataSourceConfig builds the full Terraform config: a CMS
// Contact resource followed by a CMS Contacts data source. The data source
// references the resource's contact_name via interpolation, creating an
// implicit dependency edge that delays the data source Read until apply time
// (after the resource exists), so ListContacts returns the created contact and
// contacts.0.* is populated.
func testAccCmsContactsDataSourceConfig(name, email, phone, outputFile string) string {
	outputAttr := ""
	if outputFile != "" {
		outputAttr = fmt.Sprintf("  output_file = \"%s\"\n", outputFile)
	}
	return fmt.Sprintf(`
resource "alicloud_cms_contact" "default" {
  contact_name = %[1]q
  email        = %[2]q
  phone        = %[3]q
  lang         = "zh_CN"
  im_user_ids = {
    key1 = "user1"
  }
}

data "alicloud_cms_contacts" "default" {
  contact_name = alicloud_cms_contact.default.contact_name
%[4]s
}
`, name, email, phone, outputAttr)
}
