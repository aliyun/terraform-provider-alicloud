package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Ims OidcProvider. >>> Resource test cases, automatically generated.
// Case 4434
func TestAccAliCloudImsOidcProvider_basic4434(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ims_oidc_provider.default"
	ra := resourceAttrInit(resourceId, AlicloudImsOidcProviderMap4434)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ImsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeImsOidcProvider")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%simsoidcprovider%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudImsOidcProviderBasicDependence4434)
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
					"issuer_url":         "https://oauth.aliyun.com",
					"oidc_provider_name": name,
					"fingerprints": []string{
						"902ef2deeb3c5b13ea4c3d5193629309e231ae55", "8A0246A2F6AA51BBC9D32A1353E3E63D0037A9DA"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"issuer_url":         "https://oauth.aliyun.com",
						"oidc_provider_name": name,
						"fingerprints.#":     "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "${var.oidc_provider_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"fingerprints": []string{
						"902ef2deeb3c5b13ea4c3d5193629309e231ae55"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"fingerprints.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"issuance_limit_time": "12",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"issuance_limit_time": "12",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"client_ids": []string{
						"123", "456"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"client_ids.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"issuance_limit_time": "14",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"issuance_limit_time": "14",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"issuance_limit_time": "12",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"issuance_limit_time": "12",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"fingerprints": []string{
						"8A0246A2F6AA51BBC9D32A1353E3E63D0037A9DA"},
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
						"74EF335E5E18788307FB9D89CB704BEC112ABD23487DBFF41C4DED5070F241D9", "902ef2deeb3c5b13ea4c3d5193629309e231ae55"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"fingerprints.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"client_ids": []string{
						"789"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"client_ids.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "${var.oidc_provider_name}",
					"issuer_url":  "https://oauth.aliyun.com",
					"fingerprints": []string{
						"902ef2deeb3c5b13ea4c3d5193629309e231ae55"},
					"issuance_limit_time": "12",
					"oidc_provider_name":  name + "_update",
					"client_ids": []string{
						"123", "456"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":         CHECKSET,
						"issuer_url":          "https://oauth.aliyun.com",
						"fingerprints.#":      "1",
						"issuance_limit_time": "12",
						"oidc_provider_name":  name + "_update",
						"client_ids.#":        "2",
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

var AlicloudImsOidcProviderMap4434 = map[string]string{
	"create_time":         CHECKSET,
	"issuance_limit_time": "12",
	"arn":                 CHECKSET,
}

func AlicloudImsOidcProviderBasicDependence4434(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "oidc_provider_name" {
  default = "amp-resource-test-oidc-provider"
}


`, name)
}

// Case 4434  twin
func TestAccAliCloudImsOidcProvider_basic4434_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ims_oidc_provider.default"
	ra := resourceAttrInit(resourceId, AlicloudImsOidcProviderMap4434)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ImsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeImsOidcProvider")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%simsoidcprovider%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudImsOidcProviderBasicDependence4434)
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
					"description": "test",
					"issuer_url":  "https://oauth.aliyun.com",
					"fingerprints": []string{
						"902ef2deeb3c5b13ea4c3d5193629309e231ae55", "8A0246A2F6AA51BBC9D32A1353E3E63D0037A9DA"},
					"issuance_limit_time": "14",
					"oidc_provider_name":  name,
					"client_ids": []string{
						"123", "456", "789"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":         "test",
						"issuer_url":          "https://oauth.aliyun.com",
						"fingerprints.#":      "2",
						"issuance_limit_time": "14",
						"oidc_provider_name":  name,
						"client_ids.#":        "3",
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

// Test Ims OidcProvider. <<< Resource test cases, automatically generated.
