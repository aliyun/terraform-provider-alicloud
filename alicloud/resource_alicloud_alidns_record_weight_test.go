package alicloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudAlidnsRecordWeight_weightOnly(t *testing.T) {

	recordID := os.Getenv("ALICLOUD_TEST_RECORD_ID")
	if recordID == "" {
		t.Skip("Environment variable ALICLOUD_TEST_RECORD_ID must be set with an existing DNS record ID")
	}
	resourceName := "alicloud_alidns_record_weight.default"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAlidnsRecordWeightConfigWeightOnly(recordID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "record_id", recordID),
					resource.TestCheckResourceAttr(resourceName, "weight", "60"),
				),
			},
		},
	})
}

func testAccAlidnsRecordWeightConfigWeightOnly(recordID string) string {
	return fmt.Sprintf(`
resource "alicloud_alidns_record_weight" "default" {
  record_id = "%s"
  weight    = 60
}
`, recordID)
}
