package ens_test

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// The companion alicloud_ens_common_bandwidth_package resource is not yet in
// this provider, so the acceptance test provisions the backing package through
// the ENS OpenAPI directly and injects the returned id into the Terraform
// config. The package id is not known until the dependence runs, so it is
// stored here for the checks and the destroy cleanup.
var testAccEnsCommonBandwidthPackageId string

// testAccEnsCommonBandwidthRegion is the ENS edge region used to provision the
// backing package. ENS edge regions are not the ALICLOUD_REGION the suite
// defaults to, so a stable edge region is used; the SDK v2 ENS suite used the
// same value. Override with ALICLOUD_ENS_REGION when the backend opens the
// CommonBandwidthPackage feature on a different edge region.
func testAccEnsCommonBandwidthRegion() string {
	if r := os.Getenv("ALICLOUD_ENS_REGION"); r != "" {
		return r
	}
	return "cn-chenzhou-telecom_unicom_cmcc"
}

func testAccEnsCreateCommonBandwidthPackage(client *connectivity.AliyunClient, ensRegionId, name string) (string, error) {
	request := map[string]interface{}{
		"EnsRegionId": ensRegionId,
		"Bandwidth":   5,
		"Name":        name,
	}
	response, err := client.RpcPost("Ens", "2017-11-10", "CreateCommonBandwidthPackage", nil, request, true)
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
	_, err := client.RpcPost("Ens", "2017-11-10", "DeleteCommonBandwidthPackage", nil, request, true)
	if err != nil {
		// Treat already-gone packages as cleaned up, matching the rest of the
		// provider's tolerance for the ENS not-found forms.
		if msg := err.Error(); strings.Contains(msg, "InvalidBandwidthPackageId.NotFound") || strings.Contains(msg, "NotFound") {
			return nil
		}
		return fmt.Errorf("deleting ENS common bandwidth package %s: %w", id, err)
	}
	return nil
}

func TestAccAliCloudEnsCommonBandwidthEipAttachment_basic(t *testing.T) {
	r := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-ens-cbwp-eip-att-%d", r)
	resourceId := "alicloud_ens_common_bandwidth_eip_attachment.default"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy: func(*terraform.State) error {
			// The attachment is destroyed by the resource itself; the backing
			// common bandwidth package was provisioned out-of-band and must
			// be cleaned up here.
			client := acctest.Provider.Meta().(*connectivity.AliyunClient)
			if err := testAccEnsDeleteCommonBandwidthPackage(client, testAccEnsCommonBandwidthPackageId); err != nil {
				return err
			}
			testAccEnsCommonBandwidthPackageId = ""
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: testAccEnsCommonBandwidthEipAttachmentConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceId, "id"),
					resource.TestCheckResourceAttr(resourceId, "bandwidth_package_id", testAccEnsCommonBandwidthPackageId),
					resource.TestCheckResourceAttrSet(resourceId, "ip_instance_id"),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// testAccEnsCommonBandwidthEipAttachmentConfig builds the test config. It has
// to provision the backing common bandwidth package out-of-band before the
// config string is assembled, because the attachment resource references the
// package id by value. A provisioning failure aborts the test binary before any
// apply, matching the existing ENS suite's behavior under the SDK v2 harness.
func testAccEnsCommonBandwidthEipAttachmentConfig(name string) string {
	client := acctest.Provider.Meta().(*connectivity.AliyunClient)
	packageId, err := testAccEnsCreateCommonBandwidthPackage(client, testAccEnsCommonBandwidthRegion(), name)
	if err != nil {
		log.Fatalf("[ERROR] building test dependence for ENS common bandwidth package: %v", err)
	}
	testAccEnsCommonBandwidthPackageId = packageId
	return fmt.Sprintf(`
variable "name" {
  default = %[1]q
}

variable "ens_region_id" {
  default = %[2]q
}

resource "alicloud_ens_eip" "default" {
  bandwidth            = "5"
  eip_name             = var.name
  ens_region_id        = var.ens_region_id
  internet_charge_type = "95BandwidthByMonth"
  payment_type         = "PayAsYouGo"
}

resource "alicloud_ens_common_bandwidth_eip_attachment" "default" {
  bandwidth_package_id = %[3]q
  ip_instance_id       = alicloud_ens_eip.default.id
}
`, name, testAccEnsCommonBandwidthRegion(), packageId)
}
