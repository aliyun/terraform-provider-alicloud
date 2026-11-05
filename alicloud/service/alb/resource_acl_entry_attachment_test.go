package alb_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/alb"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccAliCloudALBAclEntryAttachment_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alb_acl_entry_attachment.default"
	ra := acctest.ResourceAttrInit(resourceId, nil)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return alb.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeAclEntryAttachment")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccAlbAclEntryAttachment%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AclEntryAttachmentBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"acl_id":      "${alicloud_alb_acl.default.id}",
					"entry":       "10.10.10.0/24",
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"acl_id":      acctest.CHECKSET,
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
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccAlbAclEntryAttachment%d", rand)
	updatedName := fmt.Sprintf("%s-updated", name)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AclEntryAttachmentBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckAclEntryAttachmentBatchDestroy,
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
					testAccCheckAclEntryAttachmentBatchExists(resourceId, &entries),
					resource.TestCheckResourceAttr(resourceId, "entries.#", "3"),
					testAccCheckAclEntryAttachmentEntry(resourceId, "10.10.10.0/24", name),
					testAccCheckAclEntryAttachmentEntry(resourceId, "10.10.11.0/24", name),
					testAccCheckAclEntryAttachmentEntry(resourceId, "10.10.12.0/24", name),
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
					testAccCheckAclEntryAttachmentBatchExists(resourceId, &entries),
					resource.TestCheckResourceAttr(resourceId, "entries.#", "3"),
					testAccCheckAclEntryAttachmentEntry(resourceId, "10.10.10.0/24", name),
					testAccCheckAclEntryAttachmentEntry(resourceId, "10.10.11.0/24", updatedName),
					testAccCheckAclEntryAttachmentEntry(resourceId, "10.10.13.0/24", name),
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
					testAccCheckAclEntryAttachmentBatchExists(resourceId, &entries),
					resource.TestCheckResourceAttr(resourceId, "entries.#", "3"),
					testAccCheckAclEntryAttachmentEntry(resourceId, "10.10.10.0/24", name),
					testAccCheckAclEntryAttachmentEntry(resourceId, "10.10.11.0/24", updatedName),
					testAccCheckAclEntryAttachmentEntry(resourceId, "10.10.13.0/24", name),
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
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccAlbAclEntryAttachment%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AclEntryAttachmentBasicDependence0)
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
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckAclEntryAttachmentBatchDestroy,
		Steps: []resource.TestStep{
			{
				// exactly 20 entries: the single-call boundary of the api
				Config: testAccConfig(map[string]interface{}{
					"acl_id":  "${alicloud_alb_acl.default.id}",
					"entries": entryBlocks(20, 20),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAclEntryAttachmentBatchExists(resourceId, &entries),
					resource.TestCheckResourceAttr(resourceId, "entries.#", "20"),
					testAccCheckAclEntryAttachmentEntry(resourceId, "10.20.0.0/24", name),
					testAccCheckAclEntryAttachmentEntry(resourceId, "10.20.19.0/24", name),
				),
			},
			{
				// grow to 43: 23 added entries are sent in 2 chunked calls (20+3)
				Config: testAccConfig(map[string]interface{}{
					"acl_id":  "${alicloud_alb_acl.default.id}",
					"entries": entryBlocks(43, 20),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAclEntryAttachmentBatchExists(resourceId, &entries),
					resource.TestCheckResourceAttr(resourceId, "entries.#", "43"),
					testAccCheckAclEntryAttachmentEntry(resourceId, "10.20.20.0/24", name),
					testAccCheckAclEntryAttachmentEntry(resourceId, "10.20.42.0/24", name),
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
					testAccCheckAclEntryAttachmentBatchExists(resourceId, &entries),
					resource.TestCheckResourceAttr(resourceId, "entries.#", "44"),
					testAccCheckAclEntryAttachmentEntry(resourceId, "10.20.0.0/24", name),
					testAccCheckAclEntryAttachmentEntry(resourceId, "10.20.21.0/24", name),
					testAccCheckAclEntryAttachmentEntry(resourceId, "10.21.0.0/24", name),
					testAccCheckAclEntryAttachmentEntry(resourceId, "10.21.21.0/24", name),
				),
			},
		},
	})
}

func testAccCheckAclEntryAttachmentBatchExists(n string, entries *[]map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return errs.WrapError(errs.Error("Not found: %s", n))
		}
		if rs.Primary.ID == "" {
			return errs.WrapError(errs.Error("No alicloud_alb_acl_entry_attachment ID is set"))
		}
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)
		albService := alb.NewService(client)
		list, err := albService.ListAclEntries(rs.Primary.ID)
		if err != nil {
			return errs.WrapError(err)
		}
		if len(list) == 0 {
			return errs.WrapError(errs.Error("the acl %s has no entries", rs.Primary.ID))
		}
		*entries = list
		return nil
	}
}

// testAccCheckAclEntryAttachmentEntry asserts that one `entries` set element
// with the given entry and description and the Available status exists in the
// resource state.
func testAccCheckAclEntryAttachmentEntry(n, entry, description string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return errs.WrapError(errs.Error("Not found: %s", n))
		}
		for attributeName, attributeValue := range rs.Primary.Attributes {
			if !strings.HasPrefix(attributeName, "entries.") || !strings.HasSuffix(attributeName, ".entry") || attributeValue != entry {
				continue
			}
			hash := strings.TrimSuffix(strings.TrimPrefix(attributeName, "entries."), ".entry")
			actual := rs.Primary.Attributes[fmt.Sprintf("entries.%s.description", hash)]
			if actual != description {
				return errs.WrapError(errs.Error("entry %s has description %q, expected %q", entry, actual, description))
			}
			status := rs.Primary.Attributes[fmt.Sprintf("entries.%s.status", hash)]
			if status != "Available" {
				return errs.WrapError(errs.Error("entry %s has status %q, expected %q", entry, status, "Available"))
			}
			return nil
		}
		return errs.WrapError(errs.Error("entry %s not found in resource %s", entry, n))
	}
}

func testAccCheckAclEntryAttachmentBatchDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_alb_acl_entry_attachment" {
			continue
		}
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)
		albService := alb.NewService(client)
		entries, err := albService.ListAclEntries(rs.Primary.ID)
		if err != nil {
			if errs.NotFoundError(err) {
				continue
			}
			return errs.WrapError(err)
		}
		if len(entries) > 0 {
			return errs.WrapError(errs.Error("the acl %s still has %d entries", rs.Primary.ID, len(entries)))
		}
	}
	return nil
}

func AclEntryAttachmentBasicDependence0(name string) string {
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
