package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudCmsContactsDataSource_basic(t *testing.T) {
	resourceId := "data.alicloud_cms_contacts.default"
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccmscontactds%d", rand)
	email := fmt.Sprintf("tf-ds%d@example.com", rand)
	phone := fmt.Sprintf("86-137%08d", rand)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSourceCmsContactsConfig(name, email, phone, ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceId, "contacts.#"),
					resource.TestCheckResourceAttr(resourceId, "contacts.0.contact_name", name),
					resource.TestCheckResourceAttr(resourceId, "contacts.0.email", email),
					resource.TestCheckResourceAttr(resourceId, "contacts.0.phone", phone),
				),
			},
			{
				Config: dataSourceCmsContactsConfig(name, email, phone, "/tmp/contact_output.txt"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceId, "contacts.#"),
				),
			},
		},
	})
}

// dataSourceCmsContactsConfig builds the full Terraform config: a CMS Contact
// resource (via dataSourceCmsContactsDependence) followed by a CMS Contacts
// data source. The data source references the resource's contact_name via
// interpolation, creating an implicit dependency edge that delays the data
// source Read until apply time (after the resource exists). This also ensures
// the data source state is persisted in post-apply state, so ListCmsContacts
// returns the created contact and contacts.0.* is populated.
func dataSourceCmsContactsConfig(name, email, phone, outputFile string) string {
	resourceBlock := dataSourceCmsContactsDependence(name, email, phone)("")
	var outputAttr string
	if outputFile != "" {
		outputAttr = fmt.Sprintf("  output_file  = \"%s\"\n", outputFile)
	}
	return fmt.Sprintf(`%s
data "alicloud_cms_contacts" "default" {
  contact_name = alicloud_cms_contact.default.contact_name
%s}
`, resourceBlock, outputAttr)
}

func dataSourceCmsContactsDependence(name string, email string, phone string) func(string) string {
	return func(_ string) string {
		return fmt.Sprintf(`
resource "alicloud_cms_contact" "default" {
  contact_name = "%s"
  email        = "%s"
  phone        = "%s"
  lang         = "zh_CN"
  im_user_ids  = {
    key1 = "user1"
  }
}
`, name, email, phone)
	}
}
