package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAlicloudAlidnsRecordsWeightDataSource(t *testing.T) {
	rand := acctest.RandInt()
	resourceId := "data.alicloud_alidns_records_weight.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testacc%salidnsweight%v.abc", defaultRegionToTest, rand),
		dataSourceAlidnsRecordsWeightConfigDependence)

	testAccCheck := func(expected map[string]string) resource.TestCheckFunc {
		return resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceId, "records.#", "1"),
			resource.TestCheckResourceAttr(resourceId, "records.0.domain_name", expected["domain_name"]),
			resource.TestCheckResourceAttr(resourceId, "records.0.rr", expected["rr"]),
			resource.TestCheckResourceAttr(resourceId, "records.0.type", expected["type"]),
			resource.TestCheckResourceAttr(resourceId, "records.0.value", expected["value"]),
			resource.TestCheckResourceAttr(resourceId, "records.0.ttl", expected["ttl"]),
			resource.TestCheckResourceAttr(resourceId, "records.0.weight", expected["weight"]),
			resource.TestCheckResourceAttr(resourceId, "records.0.wrr_status", expected["wrr_status"]),
			resource.TestCheckResourceAttr(resourceId, "records.0.remark", expected["remark"]),
			resource.TestCheckResourceAttr(resourceId, "records.0.priority", expected["priority"]),
		)
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAliCloudAlidnsRecordWeightDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"domain_name": "${alicloud_dns_domain.default.name}",
					"rr":          "test1",
					"type":        "A",
					"value":       "1.1.1.1",
					"ttl":         "900",
					"weight":      "60",
					"wrr_status":  "ENABLE",
					"remark":      "test new record",
					"priority":    "5",
				}),
				Check: testAccCheck(map[string]string{
					"domain_name": fmt.Sprintf("tf-testacc%salidnsweight%v.abc", defaultRegionToTest, rand),
					"rr":          "test1",
					"type":        "A",
					"value":       "1.1.1.1",
					"ttl":         "900",
					"weight":      "60",
					"wrr_status":  "ENABLE",
					"remark":      "test new record",
					"priority":    "5",
				}),
			},
		},
	})
}

func dataSourceAlidnsRecordsWeightConfigDependence(name string) string {
	return fmt.Sprintf(`
resource "alicloud_dns_domain" "default" {
  domain_name = "%s"
}

resource "alicloud_alidns_record_weight" "default" {
  domain_name = "${alicloud_dns_domain.default.domain_name}"
  rr          = "test1"
  type        = "A"
  value       = "1.1.1.1"
  ttl         = 900
  weight      = 60
  wrr_status  = "ENABLE"
  remark      = "test new record"
  priority    = 5
}
`, name)
}

func testAccCheckAliCloudAlidnsRecordWeightDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_alidns_record_weight" {
			continue
		}

		// Check if the record still exists
		request := alidns.CreateDescribeDomainRecordInfoRequest()
		request.RecordId = rs.Primary.ID

		_, err := client.WithAlidnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
			return alidnsClient.DescribeDomainRecordInfo(request)
		})

		if err == nil {
			return fmt.Errorf("record weight %s still exists", rs.Primary.ID)
		}

		// If error is not related to a missing record, return it
		if !IsExpectedErrors(err, []string{"InvalidRecordId.NotFound"}) {
			return err
		}
	}

	return nil
}
