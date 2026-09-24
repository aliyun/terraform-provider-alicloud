package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAccAlicloudLogServiceDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudLogServiceDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_log_service.current"),
					resource.TestCheckResourceAttrSet("data.alicloud_log_service.current", "id"),
					resource.TestCheckResourceAttr("data.alicloud_log_service.current", "status", "Opened"),
				),
			},
		},
	})
}

const testAccCheckAlicloudLogServiceDataSource = `
data "alicloud_log_service" "current" {
	enable = "On"
}
`

func TestUnitLogServiceCompatibility(t *testing.T) {
	for _, tc := range []struct {
		name       string
		config     map[string]interface{}
		wantID     string
		wantStatus string
	}{
		{name: "default", config: map[string]interface{}{}, wantID: "LogServiceHasNotBeenOpened"},
		{name: "Off", config: map[string]interface{}{"enable": "Off"}, wantID: "LogServiceHasNotBeenOpened"},
		{name: "On", config: map[string]interface{}{"enable": "On"}, wantID: "LogServiceHasBeenOpened", wantStatus: "Opened"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, dataSourceAlicloudLogService().Schema, tc.config)
			d.SetId("previous-id")
			if err := d.Set("status", "previous-status"); err != nil {
				t.Fatal(err)
			}
			// Even enable=On must work without a provider client or service API calls.
			if err := dataSourceAlicloudLogServiceRead(d, nil); err != nil {
				t.Fatal(err)
			}
			if d.Id() != tc.wantID || d.Get("status") != tc.wantStatus {
				t.Fatalf("id=%q, status=%v; want id=%q, status=%q", d.Id(), d.Get("status"), tc.wantID, tc.wantStatus)
			}
		})
	}
}
