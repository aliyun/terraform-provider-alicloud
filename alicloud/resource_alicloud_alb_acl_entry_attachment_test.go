package alicloud

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAliCloudALBAclEntryAttachment_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alb_acl_entry_attachment.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAlbAclEntryAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccAlbAclEntryAttachment%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAlbAclEntryAttachmentBasicDependence0)
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
					"acl_id":      "${alicloud_alb_acl.default.id}",
					"entry":       "10.10.10.0/24",
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"acl_id":      CHECKSET,
						"entry":       "10.10.10.0/24",
						"description": name,
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

func TestAccAliCloudALBAclEntryAttachment_entries(t *testing.T) {
	var entries []map[string]interface{}
	resourceId := "alicloud_alb_acl_entry_attachment.default"
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccAlbAclEntryAttachment%d", rand)
	updatedName := fmt.Sprintf("%s-updated", name)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAlbAclEntryAttachmentBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckAlbAclEntryAttachmentBatchDestroy,
		Steps: []resource.TestStep{
			{
				// exactly one of `entry` and `entries` must be specified: a config
				// with zero entry blocks is rejected at plan time (the sdk reports
				// the exactly-one-of keys in sorted order)
				Config:      testAccConfig(map[string]interface{}{"acl_id": "${alicloud_alb_acl.default.id}"}),
				ExpectError: regexp.MustCompile("one of `entries,entry` must be specified"),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"acl_id": "${alicloud_alb_acl.default.id}",
					"entries": []map[string]interface{}{
						{
							"entry":       "10.10.10.0/24",
							"description": name,
						},
						{
							"entry":       "10.10.11.0/24",
							"description": name,
						},
						{
							"entry":       "10.10.12.0/24",
							"description": name,
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlbAclEntryAttachmentBatchExists(resourceId, &entries),
					resource.TestCheckResourceAttr(resourceId, "entries.#", "3"),
					testAccCheckAlbAclEntryAttachmentEntry(resourceId, "10.10.10.0/24", name),
					testAccCheckAlbAclEntryAttachmentEntry(resourceId, "10.10.11.0/24", name),
					testAccCheckAlbAclEntryAttachmentEntry(resourceId, "10.10.12.0/24", name),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"acl_id": "${alicloud_alb_acl.default.id}",
					"entries": []map[string]interface{}{
						{
							"entry":       "10.10.10.0/24",
							"description": name,
						},
						{
							"entry":       "10.10.11.0/24",
							"description": updatedName,
						},
						{
							"entry":       "10.10.13.0/24",
							"description": name,
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlbAclEntryAttachmentBatchExists(resourceId, &entries),
					resource.TestCheckResourceAttr(resourceId, "entries.#", "3"),
					testAccCheckAlbAclEntryAttachmentEntry(resourceId, "10.10.10.0/24", name),
					testAccCheckAlbAclEntryAttachmentEntry(resourceId, "10.10.11.0/24", updatedName),
					testAccCheckAlbAclEntryAttachmentEntry(resourceId, "10.10.13.0/24", name),
				),
			},
			{
				// the order of the entry blocks is not significant: reordering
				// them must not update the resource
				Config: testAccConfig(map[string]interface{}{
					"acl_id": "${alicloud_alb_acl.default.id}",
					"entries": []map[string]interface{}{
						{
							"entry":       "10.10.13.0/24",
							"description": name,
						},
						{
							"entry":       "10.10.10.0/24",
							"description": name,
						},
						{
							"entry":       "10.10.11.0/24",
							"description": updatedName,
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlbAclEntryAttachmentBatchExists(resourceId, &entries),
					resource.TestCheckResourceAttr(resourceId, "entries.#", "3"),
					testAccCheckAlbAclEntryAttachmentEntry(resourceId, "10.10.10.0/24", name),
					testAccCheckAlbAclEntryAttachmentEntry(resourceId, "10.10.11.0/24", updatedName),
					testAccCheckAlbAclEntryAttachmentEntry(resourceId, "10.10.13.0/24", name),
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

func TestAccAliCloudALBAclEntryAttachment_entriesChunking(t *testing.T) {
	var entries []map[string]interface{}
	resourceId := "alicloud_alb_acl_entry_attachment.default"
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccAlbAclEntryAttachment%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAlbAclEntryAttachmentBasicDependence0)
	// the api accepts at most 20 entries per call: volumes above that are
	// chunked, so the scenarios below cross the chunk boundary in every
	// direction (create, add, remove and destroy)
	entryBlocks := func(count, secondOctet int) []map[string]interface{} {
		blocks := make([]map[string]interface{}, 0, count)
		for i := 0; i < count; i++ {
			blocks = append(blocks, map[string]interface{}{
				"entry":       fmt.Sprintf("10.%d.%d.0/24", secondOctet, i),
				"description": name,
			})
		}
		return blocks
	}
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckAlbAclEntryAttachmentBatchDestroy,
		Steps: []resource.TestStep{
			{
				// exactly 20 entries: the single-call boundary of the api
				Config: testAccConfig(map[string]interface{}{
					"acl_id":  "${alicloud_alb_acl.default.id}",
					"entries": entryBlocks(20, 20),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlbAclEntryAttachmentBatchExists(resourceId, &entries),
					resource.TestCheckResourceAttr(resourceId, "entries.#", "20"),
					testAccCheckAlbAclEntryAttachmentEntry(resourceId, "10.20.0.0/24", name),
					testAccCheckAlbAclEntryAttachmentEntry(resourceId, "10.20.19.0/24", name),
				),
			},
			{
				// grow to 43: 23 added entries are sent in 2 chunked calls (20+3)
				Config: testAccConfig(map[string]interface{}{
					"acl_id":  "${alicloud_alb_acl.default.id}",
					"entries": entryBlocks(43, 20),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlbAclEntryAttachmentBatchExists(resourceId, &entries),
					resource.TestCheckResourceAttr(resourceId, "entries.#", "43"),
					testAccCheckAlbAclEntryAttachmentEntry(resourceId, "10.20.20.0/24", name),
					testAccCheckAlbAclEntryAttachmentEntry(resourceId, "10.20.42.0/24", name),
				),
			},
			{
				// swap: 21 removed entries are sent in 2 chunked calls (20+1)
				// and 22 added entries in 2 chunked calls (20+2)
				Config: testAccConfig(map[string]interface{}{
					"acl_id":  "${alicloud_alb_acl.default.id}",
					"entries": append(entryBlocks(22, 20), entryBlocks(22, 21)...),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlbAclEntryAttachmentBatchExists(resourceId, &entries),
					resource.TestCheckResourceAttr(resourceId, "entries.#", "44"),
					testAccCheckAlbAclEntryAttachmentEntry(resourceId, "10.20.0.0/24", name),
					testAccCheckAlbAclEntryAttachmentEntry(resourceId, "10.20.21.0/24", name),
					testAccCheckAlbAclEntryAttachmentEntry(resourceId, "10.21.0.0/24", name),
					testAccCheckAlbAclEntryAttachmentEntry(resourceId, "10.21.21.0/24", name),
				),
			},
		},
	})
}

func testAccCheckAlbAclEntryAttachmentBatchExists(n string, entries *[]map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return WrapError(Error("Not found: %s", n))
		}
		if rs.Primary.ID == "" {
			return WrapError(Error("No alicloud_alb_acl_entry_attachment ID is set"))
		}
		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		albService := AlbService{client}
		list, err := albService.ListAclEntries(rs.Primary.ID)
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

// testAccCheckAlbAclEntryAttachmentEntry asserts that one `entries` set element
// with the given entry and description and the Available status exists in the
// resource state.
func testAccCheckAlbAclEntryAttachmentEntry(n, entry, description string) resource.TestCheckFunc {
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
			actual := rs.Primary.Attributes[fmt.Sprintf("entries.%s.description", hash)]
			if actual != description {
				return WrapError(Error("entry %s has description %q, expected %q", entry, actual, description))
			}
			status := rs.Primary.Attributes[fmt.Sprintf("entries.%s.status", hash)]
			if status != "Available" {
				return WrapError(Error("entry %s has status %q, expected %q", entry, status, "Available"))
			}
			return nil
		}
		return WrapError(Error("entry %s not found in resource %s", entry, n))
	}
}

func testAccCheckAlbAclEntryAttachmentBatchDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_alb_acl_entry_attachment" {
			continue
		}
		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		albService := AlbService{client}
		entries, err := albService.ListAclEntries(rs.Primary.ID)
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

func AlicloudAlbAclEntryAttachmentBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_alb_acl" "default" {
  acl_name          = var.name
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.groups.0.id
}
`, name)
}
