package alicloud

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAliCloudSLBAclEntryAttachment_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_slb_acl_entry_attachment.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSlbAclEntryAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccSlbAclEntryAttachment%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSlbAclEntryAttachmentBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"acl_id":  "${alicloud_slb_acl.default.id}",
					"entry":   "10.10.10.0/24",
					"comment": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"acl_id":  CHECKSET,
						"entry":   "10.10.10.0/24",
						"comment": name,
					}),
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

func TestAccAliCloudSLBAclEntryAttachment_entries(t *testing.T) {
	var entries []map[string]interface{}
	resourceId := "alicloud_slb_acl_entry_attachment.default"
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccSlbAclEntryAttachment%d", rand)
	updatedName := fmt.Sprintf("%s-updated", name)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSlbAclEntryAttachmentBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSlbAclEntryAttachmentBatchDestroy,
		Steps: []resource.TestStep{
			{
				// exactly one of `entry` and `entries` must be specified: a config
				// with zero entry blocks is rejected at plan time (the sdk reports
				// the exactly-one-of keys in sorted order)
				Config:      testAccConfig(map[string]interface{}{"acl_id": "${alicloud_slb_acl.default.id}"}),
				ExpectError: regexp.MustCompile("one of `entries,entry` must be specified"),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"acl_id": "${alicloud_slb_acl.default.id}",
					"entries": []map[string]interface{}{
						{
							"entry":   "10.10.10.0/24",
							"comment": name,
						},
						{
							"entry":   "10.10.11.0/24",
							"comment": name,
						},
						{
							"entry":   "10.10.12.0/24",
							"comment": name,
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSlbAclEntryAttachmentBatchExists(resourceId, &entries),
					resource.TestCheckResourceAttr(resourceId, "entries.#", "3"),
					testAccCheckSlbAclEntryAttachmentEntry(resourceId, "10.10.10.0/24", name),
					testAccCheckSlbAclEntryAttachmentEntry(resourceId, "10.10.11.0/24", name),
					testAccCheckSlbAclEntryAttachmentEntry(resourceId, "10.10.12.0/24", name),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"acl_id": "${alicloud_slb_acl.default.id}",
					"entries": []map[string]interface{}{
						{
							"entry":   "10.10.10.0/24",
							"comment": name,
						},
						{
							"entry":   "10.10.11.0/24",
							"comment": updatedName,
						},
						{
							"entry":   "10.10.13.0/24",
							"comment": name,
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSlbAclEntryAttachmentBatchExists(resourceId, &entries),
					resource.TestCheckResourceAttr(resourceId, "entries.#", "3"),
					testAccCheckSlbAclEntryAttachmentEntry(resourceId, "10.10.10.0/24", name),
					testAccCheckSlbAclEntryAttachmentEntry(resourceId, "10.10.11.0/24", updatedName),
					testAccCheckSlbAclEntryAttachmentEntry(resourceId, "10.10.13.0/24", name),
				),
			},
			{
				// the order of the entry blocks is not significant: reordering
				// them must not update the resource
				Config: testAccConfig(map[string]interface{}{
					"acl_id": "${alicloud_slb_acl.default.id}",
					"entries": []map[string]interface{}{
						{
							"entry":   "10.10.13.0/24",
							"comment": name,
						},
						{
							"entry":   "10.10.10.0/24",
							"comment": name,
						},
						{
							"entry":   "10.10.11.0/24",
							"comment": updatedName,
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSlbAclEntryAttachmentBatchExists(resourceId, &entries),
					resource.TestCheckResourceAttr(resourceId, "entries.#", "3"),
					testAccCheckSlbAclEntryAttachmentEntry(resourceId, "10.10.10.0/24", name),
					testAccCheckSlbAclEntryAttachmentEntry(resourceId, "10.10.11.0/24", updatedName),
					testAccCheckSlbAclEntryAttachmentEntry(resourceId, "10.10.13.0/24", name),
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

func TestAccAliCloudSLBAclEntryAttachment_entriesChunking(t *testing.T) {
	var entries []map[string]interface{}
	resourceId := "alicloud_slb_acl_entry_attachment.default"
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccSlbAclEntryAttachment%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSlbAclEntryAttachmentBasicDependence0)
	// the api accepts at most 50 entries per call: volumes above that are
	// chunked, so the scenarios below cross the chunk boundary in every
	// direction (create, add, remove and destroy)
	entryBlocks := func(count, secondOctet int) []map[string]interface{} {
		blocks := make([]map[string]interface{}, 0, count)
		for i := 0; i < count; i++ {
			blocks = append(blocks, map[string]interface{}{
				"entry":   fmt.Sprintf("10.%d.%d.0/24", secondOctet, i),
				"comment": name,
			})
		}
		return blocks
	}
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSlbAclEntryAttachmentBatchDestroy,
		Steps: []resource.TestStep{
			{
				// exactly 50 entries: the single-call boundary of the api
				Config: testAccConfig(map[string]interface{}{
					"acl_id":  "${alicloud_slb_acl.default.id}",
					"entries": entryBlocks(50, 20),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSlbAclEntryAttachmentBatchExists(resourceId, &entries),
					resource.TestCheckResourceAttr(resourceId, "entries.#", "50"),
					testAccCheckSlbAclEntryAttachmentEntry(resourceId, "10.20.0.0/24", name),
					testAccCheckSlbAclEntryAttachmentEntry(resourceId, "10.20.49.0/24", name),
				),
			},
			{
				// grow to 103: 53 added entries are sent in 2 chunked calls (50+3)
				Config: testAccConfig(map[string]interface{}{
					"acl_id":  "${alicloud_slb_acl.default.id}",
					"entries": entryBlocks(103, 20),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSlbAclEntryAttachmentBatchExists(resourceId, &entries),
					resource.TestCheckResourceAttr(resourceId, "entries.#", "103"),
					testAccCheckSlbAclEntryAttachmentEntry(resourceId, "10.20.50.0/24", name),
					testAccCheckSlbAclEntryAttachmentEntry(resourceId, "10.20.102.0/24", name),
				),
			},
			{
				// swap: 51 removed entries are sent in 2 chunked calls (50+1)
				// and 52 added entries in 2 chunked calls (50+2)
				Config: testAccConfig(map[string]interface{}{
					"acl_id":  "${alicloud_slb_acl.default.id}",
					"entries": append(entryBlocks(52, 20), entryBlocks(52, 21)...),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSlbAclEntryAttachmentBatchExists(resourceId, &entries),
					resource.TestCheckResourceAttr(resourceId, "entries.#", "104"),
					testAccCheckSlbAclEntryAttachmentEntry(resourceId, "10.20.0.0/24", name),
					testAccCheckSlbAclEntryAttachmentEntry(resourceId, "10.20.51.0/24", name),
					testAccCheckSlbAclEntryAttachmentEntry(resourceId, "10.21.0.0/24", name),
					testAccCheckSlbAclEntryAttachmentEntry(resourceId, "10.21.51.0/24", name),
				),
			},
		},
	})
}

func testAccCheckSlbAclEntryAttachmentBatchExists(n string, entries *[]map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return WrapError(Error("Not found: %s", n))
		}
		if rs.Primary.ID == "" {
			return WrapError(Error("No alicloud_slb_acl_entry_attachment ID is set"))
		}
		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		slbService := SlbService{client}
		list, err := slbService.ListSlbAclEntries(rs.Primary.ID)
		if err != nil {
			return WrapError(err)
		}
		if len(list) == 0 {
			return WrapError(Error("the acl %s has no entries", rs.Primary.ID))
		}
		*entries = list
		return nil
	}
}

// testAccCheckSlbAclEntryAttachmentEntry asserts that one `entries` set element
// with the given entry and comment exists in the resource state.
func testAccCheckSlbAclEntryAttachmentEntry(n, entry, comment string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return WrapError(Error("Not found: %s", n))
		}
		for attributeName, attributeValue := range rs.Primary.Attributes {
			if !strings.HasPrefix(attributeName, "entries.") || !strings.HasSuffix(attributeName, ".entry") || attributeValue != entry {
				continue
			}
			hash := strings.TrimSuffix(strings.TrimPrefix(attributeName, "entries."), ".entry")
			actual := rs.Primary.Attributes[fmt.Sprintf("entries.%s.comment", hash)]
			if actual != comment {
				return WrapError(Error("entry %s has comment %q, expected %q", entry, actual, comment))
			}
			return nil
		}
		return WrapError(Error("entry %s not found in resource %s", entry, n))
	}
}

func testAccCheckSlbAclEntryAttachmentBatchDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_slb_acl_entry_attachment" {
			continue
		}
		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		slbService := SlbService{client}
		entries, err := slbService.ListSlbAclEntries(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}
		if len(entries) > 0 {
			return WrapError(Error("the acl %s still has %d entries", rs.Primary.ID, len(entries)))
		}
	}
	return nil
}

func AlicloudSlbAclEntryAttachmentBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_slb_acl" "default" {
  name       = var.name
  ip_version = "ipv4"
}
`, name)
}
