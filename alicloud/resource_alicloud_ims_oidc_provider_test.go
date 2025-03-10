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
