package alicloud

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// testAccEnsCommonBandwidthRegion resolves the region used to build the
// standalone client that provisions the backing ENS common bandwidth package.
// The package-level defaultRegionToTest is initialized from
// os.Getenv("ALICLOUD_REGION") at init time and is only re-assigned inside
// testAccPreCheck when ALICLOUD_REGION is already set, so it stays empty in
// remote AccTest runs that do not pre-set ALICLOUD_REGION. Both the Dependence
// step (which runs before PreCheck while constructing the TestCase steps) and
// CheckDestroy would therefore pass an empty region to sharedClientForRegion
// and abort with "Invalid Alibaba Cloud region: .". Mirroring the gpdb /
// resource-manager pattern, fall back to cn-beijing (the same default
// testAccPreCheck sets) so the test reaches create/delete.
func testAccEnsCommonBandwidthRegion() string {
	region := os.Getenv("ALICLOUD_REGION")
	if region == "" {
		region = "cn-beijing"
	}
	return region
}

// The ENS common bandwidth package resource is not yet available in this
// provider, so the acceptance test provisions the backing package through the
// ENS OpenAPI directly and embeds the returned id into the Terraform config.
var testAccEnsCommonBandwidthPackageId string

func testAccEnsCreateCommonBandwidthPackage(client *connectivity.AliyunClient, ensRegionId, name string) (string, error) {
	request := map[string]interface{}{
		"EnsRegionId": ensRegionId,
		"Bandwidth":   5,
		"Name":        name,
	}
	query := make(map[string]interface{})
	response, err := client.RpcPost("Ens", "2017-11-10", "CreateCommonBandwidthPackage", query, request, true)
	if err != nil {
		return "", fmt.Errorf("creating ENS common bandwidth package: %w", err)
	}
	if v, ok := response["BandwidthPackageId"]; ok && fmt.Sprint(v) != "" {
		return fmt.Sprint(v), nil
	}
	return "", fmt.Errorf("BandwidthPackageId missing from CreateCommonBandwidthPackage response: %v", response)
}

func testAccEnsDeleteCommonBandwidthPackage(client *connectivity.AliyunClient, id string) error {
	if id == "" {
		return nil
	}
	request := map[string]interface{}{
		"BandwidthPackageId": id,
	}
	query := make(map[string]interface{})
	_, err := client.RpcPost("Ens", "2017-11-10", "DeleteCommonBandwidthPackage", query, request, true)
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidBandwidthPackageId.NotFound", "NotFound"}) || NotFoundError(err) {
			return nil
		}
		return fmt.Errorf("deleting ENS common bandwidth package %s: %w", id, err)
	}
	return nil
}

func TestAccAliCloudEnsCommonBandwidthEipAttachment_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ens_common_bandwidth_eip_attachment.default"
	ra := resourceAttrInit(resourceId, AlicloudEnsCommonBandwidthEipAttachmentMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EnsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEnsCommonBandwidthEipAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-ens-cbwp-eip-att-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEnsCommonBandwidthEipAttachmentBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy: func(s *terraform.State) error {
			if err := rac.checkResourceDestroy()(s); err != nil {
				return err
			}
			rawClient, err := sharedClientForRegion(testAccEnsCommonBandwidthRegion())
			if err != nil {
				return fmt.Errorf("getting client for cleanup: %w", err)
			}
			client := rawClient.(*connectivity.AliyunClient)
			if derr := testAccEnsDeleteCommonBandwidthPackage(client, testAccEnsCommonBandwidthPackageId); derr != nil {
				return derr
			}
			testAccEnsCommonBandwidthPackageId = ""
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth_package_id": "${local.bandwidth_package_id}",
					"ip_instance_id":       "${alicloud_ens_eip.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth_package_id": testAccEnsCommonBandwidthPackageId,
						"ip_instance_id":       CHECKSET,
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

var AlicloudEnsCommonBandwidthEipAttachmentMap = map[string]string{
	"bandwidth_package_id": CHECKSET,
	"ip_instance_id":       CHECKSET,
}

func AlicloudEnsCommonBandwidthEipAttachmentBasicDependence(name string) string {
	rawClient, err := sharedClientForRegion(testAccEnsCommonBandwidthRegion())
	if err != nil {
		log.Fatalf("[ERROR] failed to build standalone Alicloud client for ENS common bandwidth package setup: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	packageId, err := testAccEnsCreateCommonBandwidthPackage(client, "cn-chenzhou-telecom_unicom_cmcc", name)
	if err != nil {
		log.Fatalf("[ERROR] %s", err)
	}
	testAccEnsCommonBandwidthPackageId = packageId
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

variable "ens_region_id" {
  default = "cn-chenzhou-telecom_unicom_cmcc"
}

locals {
  bandwidth_package_id = "%s"
}

resource "alicloud_ens_eip" "default" {
  bandwidth            = "5"
  eip_name             = var.name
  ens_region_id        = var.ens_region_id
  internet_charge_type = "95BandwidthByMonth"
  payment_type         = "PayAsYouGo"
}
`, name, packageId)
}
