// Package mns_test.
package mns_test

import (
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAlicloudMnsServiceDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, false, []connectivity.Region{connectivity.EUCentral1})
		},
		ProviderFactories: acctest.SDKProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudMnsServiceDataSource,
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckAlicloudDataSourceID("data.alicloud_mns_service.current"),
					resource.TestCheckResourceAttrSet("data.alicloud_mns_service.current", "id"),
					resource.TestCheckResourceAttr("data.alicloud_mns_service.current", "status", "Opened"),
				),
			},
		},
	})
}

const testAccCheckAlicloudMnsServiceDataSource = `
data "alicloud_mns_service" "current" {
	enable = "On"
}
`

// TestAccAlicloudMnsServiceDataSource_readOnly exercises the read-only branch of
// the alicloud_mns_service data source with enable = "Off": no OpenService write
// API is invoked, and the data source sets id to "MnsServiceHasNotBeenOpened"
// with an empty status. This case only covers the existing Off branch; it does
// not verify the enable = "On"/OpenService path. No new MNS service opening or
// account-permission change is performed.
func TestAccAlicloudMnsServiceDataSource_readOnly(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ProviderFactories: acctest.SDKProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudMnsServiceDataSourceReadOnly,
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckAlicloudDataSourceID("data.alicloud_mns_service.current"),
					resource.TestCheckResourceAttr("data.alicloud_mns_service.current", "id", "MnsServiceHasNotBeenOpened"),
					resource.TestCheckResourceAttr("data.alicloud_mns_service.current", "status", ""),
					resource.TestCheckResourceAttr("data.alicloud_mns_service.current", "enable", "Off"),
				),
			},
		},
	})
}

const testAccCheckAlicloudMnsServiceDataSourceReadOnly = `
data "alicloud_mns_service" "current" {
	enable = "Off"
}
`
