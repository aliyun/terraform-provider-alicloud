package sls

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestMain(m *testing.M) { resource.TestMain(m) }

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
