package alicloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"os"
	"testing"
)

func TestAccAlicloudFileCRC64DataSourceBasic(t *testing.T) {
	rand := acctest.RandIntRange(10000, 999999)
	name := fmt.Sprintf("tf-testaccalicloudFileCRC64Checksum-%d", rand)
	resourceId := "data.alicloud_file_crc64_checksum.default"
	ra := resourceAttrInit(resourceId, nil)
	testAccCheck := ra.resourceAttrMapUpdateSet()

	path, file, err := createTempFile(name)
	file.WriteString(`    # -*- coding: utf-8 -*-`)
	file.WriteString("\n")
	file.WriteString(`	def handler(event, context):`)
	file.WriteString("\n")
	file.WriteString(`	    print "hello world"`)
	file.WriteString("\n")

	if err != nil {
		t.Fatal(WrapError(err))
	}
	defer func() {
		file.Close()
		os.Remove(path)
	}()

	config := fmt.Sprintf(`
		data "alicloud_file_crc64_checksum" "default" {
			filename = "%s"
		}`, path)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"filename": path,
						"checksum": CHECKSET,
					}),
				),
			},
			{
				PreConfig: func() {
					file.WriteString(`	    return "hello world"`)
				},
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
		},
	})
}
