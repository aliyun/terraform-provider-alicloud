package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// NOTE: https://aliyuque.antfin.com/aliyun-subaccount/vkgfrb/ohxat6d7vn7h3gfa

func TestAccAliCloudImsOidcProvider_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ims_oidc_provider.default"
	ra := resourceAttrInit(resourceId, AlicloudImsOidcProviderMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ImsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeImsOidcProvider")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccims%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudImsOidcProviderBasicDependence0)
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
					"description": "${var.oidc_provider_name}",
					"issuer_url":  "https://oauth.aliyun.com",
					"fingerprints": []string{
						"0BBFAB97059595E8D1EC48E89EB8657C0E5AAE71"},
					"issuance_limit_time": "12",
					"oidc_provider_name":  name,
					"client_ids": []string{
						"123", "456"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":         CHECKSET,
						"issuer_url":          "https://oauth.aliyun.com",
						"fingerprints.#":      "1",
						"issuance_limit_time": "12",
						"oidc_provider_name":  name,
						"client_ids.#":        "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":         "test",
					"issuance_limit_time": "14",
					"client_ids": []string{
						"123", "456", "789"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":         "test",
						"issuance_limit_time": "14",
						"client_ids.#":        "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"issuance_limit_time": "12",
					"client_ids": []string{
						"123"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"issuance_limit_time": "12",
						"client_ids.#":        "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"fingerprints": []string{"7328B027C7E2139EE3BA57B528CDB357F8E3B9D0"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"fingerprints.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"fingerprints": []string{
						"0BBFAB97059595E8D1EC48E89EB8657C0E5AAE71", "7328B027C7E2139EE3BA57B528CDB357F8E3B9D0"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"fingerprints.#": "2",
					}),
				),
			},
			{
				// Prepare for scenario 1: increase to 5 fingerprints
				Config: testAccConfig(map[string]interface{}{
					"fingerprints": []string{
						"0BBFAB97059595E8D1EC48E89EB8657C0E5AAE71",
						"7328B027C7E2139EE3BA57B528CDB357F8E3B9D0",
						"902EF2DEEB3C5B13EA4C3D5193629309E231AE55",
						"A031C46782E6E6C662C2C87C76DA9AA62CCABD8E",
						"1FB86B1168EC743154062E8C9CC5B171A4B7CCB4"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"fingerprints.#": "5",
					}),
				),
			},
			{
				// Scenario 1: Replace all 5 fingerprints with 1 new one (would exceed limit if add first)
				Config: testAccConfig(map[string]interface{}{
					"fingerprints": []string{
						"C8A401588C2D1B9BF37833CBE2A2A4E7E8E1C1B2"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"fingerprints.#": "1",
					}),
				),
			},
			{
				// Prepare for scenario 2: increase to 3 fingerprints
				Config: testAccConfig(map[string]interface{}{
					"fingerprints": []string{
						"C8A401588C2D1B9BF37833CBE2A2A4E7E8E1C1B2",
						"D9B512699D3C2C8CF48944DDF3D3B5F8F9F2D2C3",
						"E0C623700E4D3D9DG59055EEG4E4C6G9G0G3E3D4"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"fingerprints.#": "3",
					}),
				),
			},
			{
				// Scenario 2: Replace all 3 fingerprints with 2 (won't exceed limit)
				Config: testAccConfig(map[string]interface{}{
					"fingerprints": []string{
						"0BBFAB97059595E8D1EC48E89EB8657C0E5AAE71",
						"7328B027C7E2139EE3BA57B528CDB357F8E3B9D0"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"fingerprints.#": "2",
					}),
				),
			},
			{
				// Prepare for scenario 3: increase to 5 fingerprints
				Config: testAccConfig(map[string]interface{}{
					"fingerprints": []string{
						"0BBFAB97059595E8D1EC48E89EB8657C0E5AAE71",
						"7328B027C7E2139EE3BA57B528CDB357F8E3B9D0",
						"902EF2DEEB3C5B13EA4C3D5193629309E231AE55",
						"A031C46782E6E6C662C2C87C76DA9AA62CCABD8E",
						"1FB86B1168EC743154062E8C9CC5B171A4B7CCB4"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"fingerprints.#": "5",
					}),
				),
			},
			{
				// Scenario 3: Partial update - remove 2 and add 2, keep 3 (afterRemovalCount >= 1)
				Config: testAccConfig(map[string]interface{}{
					"fingerprints": []string{
						"0BBFAB97059595E8D1EC48E89EB8657C0E5AAE71",
						"7328B027C7E2139EE3BA57B528CDB357F8E3B9D0",
						"902EF2DEEB3C5B13EA4C3D5193629309E231AE55",
						"D9B512699D3C2C8CF48944DDF3D3B5F8F9F2D2C3",
						"E0C623700E4D3D9DG59055EEG4E4C6G9G0G3E3D4"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"fingerprints.#": "5",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudImsOidcProviderMap0 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudImsOidcProviderBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "oidc_provider_name" {
  default = "amp-resource-test-oidc-provider"
}


`, name)
}
