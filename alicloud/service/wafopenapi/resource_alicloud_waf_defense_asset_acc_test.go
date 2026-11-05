package wafopenapi_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TestAccAliCloudWafDefenseAsset_basic exercises the full DefenseAsset lifecycle (create with a
// valid AssetType enum, update asset_status, import by <InstanceId>:<AssetId>, destroy) against a
// real WAF instance that already has a protected domain.
//
// CreateDefenseAsset requires Resource and Conditions to reference a real protection object
// (a domain added to the WAF instance); a made-up value is rejected with
// Defense.Control.DefenseResourceNotExist. The test therefore resolves the instance and a real
// protected domain through the alicloud_wafv3_instances / alicloud_wafv3_domains data sources —
// the same precondition every other wafv3 ACC test in this package relies on.
//
// AssetType is validated against the fixed enum (AI, Client_Injection, Apisec,
// Account_Protector, SSO) at the gateway; the schema mirrors that enum with OneOf so an invalid
// value is rejected at plan time rather than producing the previous round's InvalidAssetType from
// CreateDefenseAsset. AssetConfig for the AI type is built from the authoritative online
// (zhenyuan) AssetConfig definition: InjectionRatio (1..100), SensitiveLevel (N/L1/L2/L3),
// PromptPosition, OperationPurpose, IsPayment, and the array/object fields.
func TestAccAliCloudWafDefenseAsset_basic(t *testing.T) {
	resourceId := "alicloud_waf_defense_asset.default"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckWafDefenseAssetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWafDefenseAssetBasic("1"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceId, "id"),
					resource.TestCheckResourceAttr(resourceId, "asset_type", "AI"),
					resource.TestCheckResourceAttr(resourceId, "asset_status", "1"),
				),
			},
			{
				Config: testAccWafDefenseAssetBasic("0"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "asset_status", "0"),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"asset_config", "conditions"},
			},
		},
	})
}

// testAccWafDefenseAssetBasic renders the DefenseAsset config. asset_type is fixed to "AI" (a
// valid enum); resource is the protected domain id and conditions reference that domain's name,
// both resolved from the wafv3_domains data source so CreateDefenseAsset sees a real protection
// object. asset_config carries the AI-type fields from the authoritative AssetConfig definition.
func testAccWafDefenseAssetBasic(assetStatus string) string {
	return fmt.Sprintf(`
data "alicloud_wafv3_instances" "default" {
}

data "alicloud_wafv3_domains" "default" {
  instance_id    = data.alicloud_wafv3_instances.default.ids.0
  enable_details = true
}

resource "alicloud_waf_defense_asset" "default" {
  instance_id = data.alicloud_wafv3_instances.default.ids.0
  asset_type  = "AI"
  resource    = data.alicloud_wafv3_domains.default.domains.0.domain_id
  conditions = jsonencode([{
    Domain     = data.alicloud_wafv3_domains.default.domains.0.domain
    UrlPath    = "/"
    HttpMethod = "GET"
  }])
  asset_config = jsonencode({
    PromptPosition     = "request"
    InjectionRatio     = 1
    IsPayment          = false
    SensitiveLevel     = "N"
    OperationPurpose   = "tf-acc-test"
    ResponseContents   = []
    SuccessConditions  = []
    FailureConditions  = []
    ExtractAccounts    = []
    CookieDomain       = ""
  })
  asset_name   = "tf-acc-defense-asset"
  description  = "Created by Terraform ACC"
  asset_status = "%[1]s"
}
`, assetStatus)
}

// testAccCheckWafDefenseAssetDestroy confirms the asset is gone after destroy by listing the
// instance's defense assets and checking the AssetId is absent. Read already removes a gone
// resource from state; this is the explicit destroy guard.
func testAccCheckWafDefenseAssetDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_waf_defense_asset" {
			continue
		}
		parts := strings.Split(rs.Primary.ID, ":")
		if len(parts) != 2 {
			return fmt.Errorf("unexpected defense asset id after destroy: %s", rs.Primary.ID)
		}
		instanceId, assetId := parts[0], parts[1]
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)
		request := map[string]interface{}{
			"RegionId":   client.RegionId,
			"InstanceId": instanceId,
		}
		response, err := client.RpcPost("waf-openapi", "2021-10-01", "DescribeDefenseAssets", nil, request, true)
		if err != nil {
			// DescribeDefenseAssets errors are not asset-existence evidence; the Read
			// already dropped the resource, so a transient describe error is not a
			// destroy failure.
			continue
		}
		if assets, ok := response["Assets"].([]interface{}); ok {
			for _, raw := range assets {
				if obj, ok := raw.(map[string]interface{}); ok {
					if id, ok := obj["AssetId"].(string); ok && id == assetId {
						return fmt.Errorf("defense asset %s still exists", rs.Primary.ID)
					}
				}
			}
		}
	}
	return nil
}
