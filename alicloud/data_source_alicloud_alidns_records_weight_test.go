package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccDataSourceAlicloudAlidnsRecordsWeight_basic(t *testing.T) {
	var record alidns.DescribeDomainRecordInfoResponse
	resourceName := "alicloud_alidns_record_weight.default"
	dataSourceName := "data.alicloud_alidns_records_weight.default"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAliCloudAlidnsRecordWeightDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAlidnsRecordsWeightConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAliCloudAlidnsRecordWeightExists(resourceName, &record),
					resource.TestCheckResourceAttr(dataSourceName, "records.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "records.0.record_id", record.RecordId),
					resource.TestCheckResourceAttr(dataSourceName, "records.0.weight", "60"),
					resource.TestCheckResourceAttr(dataSourceName, "records.0.wrr_status", "true"),
				),
			},
		},
	})
}

func testAccDataSourceAlidnsRecordsWeightConfig() string {
	return `

resource "alicloud_alidns_record" "default" {
  domain_name = "snapp.services"
  rr          = "test1"
  type        = "A"
  value       = "1.1.1.1"
  ttl         = 600
}

resource "alicloud_alidns_record_weight" "default" {
  record_id = alicloud_alidns_record.default.id
  weight    = 60
}

data "alicloud_alidns_records_weight" "default" {
  record_ids = [alicloud_alidns_record_weight.default.record_id]
}
`
}

func testAccCheckAliCloudAlidnsRecordWeightExists(n string, v *alidns.DescribeDomainRecordInfoResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource %s not found", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no record ID is set")
		}

		request := alidns.CreateDescribeDomainRecordInfoRequest()
		request.RecordId = rs.Primary.ID

		raw, err := client.WithAlidnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
			return alidnsClient.DescribeDomainRecordInfo(request)
		})
		if err != nil {
			return err
		}

		response, ok := raw.(*alidns.DescribeDomainRecordInfoResponse)
		if !ok {
			return fmt.Errorf("error asserting DescribeDomainRecordInfoResponse")
		}

		*v = *response
		return nil
	}
}

func testAccCheckAliCloudAlidnsRecordWeightDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_alidns_record_weight" {
			continue
		}

		request := alidns.CreateDescribeDomainRecordInfoRequest()
		request.RecordId = rs.Primary.ID

		_, err := client.WithAlidnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
			return alidnsClient.DescribeDomainRecordInfo(request)
		})

		if err == nil {
			return fmt.Errorf("record weight %s still exists", rs.Primary.ID)
		}

		if !IsExpectedErrors(err, []string{"InvalidRecordId.NotFound"}) {
			return err
		}
	}
	return nil
}
