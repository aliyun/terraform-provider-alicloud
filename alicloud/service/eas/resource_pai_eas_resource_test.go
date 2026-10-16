package eas_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/eas"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	tfacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// testAccPaiEasResourceProvider is the single provider instance shared by the
// provider factory and CheckDestroy, so the post-destroy verification reads the
// same configured client the test steps used.
var testAccPaiEasResourceProvider = alicloud.Provider()

// testAccPaiEasResourceProtoV5ProviderFactory serves the muxed v5 provider (the
// SDK v2 half and the framework half that implements alicloud_pai_eas_resource)
// under the name "alicloud", mirroring acctest.ProtoV5ProviderFactories without
// pulling the testing-only echoprovider dependency into the test binary.
//
// The SDK v2 helper/resource harness is used deliberately: package alicloud still
// has non-test files importing terraform-plugin-sdk/v2/helper/resource (e.g. the
// PolarFS resources, which use resource.Retry / resource.StateRefreshFunc), whose
// init registers the -sweep flag. Linking terraform-plugin-testing/helper/resource
// alongside it would re-define -sweep and panic the test binary at startup before
// any step runs. The SDK v2 harness registers -sweep once, so the split-layout
// eas_test binary starts cleanly.
func testAccPaiEasResourceProtoV5ProviderFactory() map[string]func() (tfprotov5.ProviderServer, error) {
	return map[string]func() (tfprotov5.ProviderServer, error){
		"alicloud": func() (tfprotov5.ProviderServer, error) {
			serverFactory, err := provider.ProtoV5ProviderServerFactory(context.Background(), testAccPaiEasResourceProvider)
			if err != nil {
				return nil, err
			}
			return serverFactory(), nil
		},
	}
}

func testAccPaiEasResourcePreCheck(t *testing.T) {
	if v := os.Getenv("ALICLOUD_ACCESS_KEY"); v == "" {
		t.Fatal("ALICLOUD_ACCESS_KEY must be set for acceptance tests")
	}
	if v := os.Getenv("ALICLOUD_SECRET_KEY"); v == "" {
		t.Fatal("ALICLOUD_SECRET_KEY must be set for acceptance tests")
	}
	if v := os.Getenv("ALICLOUD_REGION"); v == "" {
		os.Setenv("ALICLOUD_REGION", "cn-beijing")
	}
}

// TestAccAliCloudPaiEasResource exercises the alicloud_pai_eas_resource
// resource over create, in-place update (rename) and import. charge_type uses
// PostPaid to avoid any prepaid/subscription financial commitment.
func TestAccAliCloudPaiEasResource(t *testing.T) {
	r := tfacctest.RandIntRange(100000, 999999)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPaiEasResourcePreCheck(t)
		},
		ProtoV5ProviderFactories: testAccPaiEasResourceProtoV5ProviderFactory(),
		CheckDestroy:             testAccCheckPaiEasResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudPaiEasResourceConfig(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_pai_eas_resource.default", "id"),
					resource.TestCheckResourceAttrSet("alicloud_pai_eas_resource.default", "region_id"),
					resource.TestCheckResourceAttrSet("alicloud_pai_eas_resource.default", "status"),
					resource.TestCheckResourceAttrSet("alicloud_pai_eas_resource.default", "create_time"),
					resource.TestCheckResourceAttr("alicloud_pai_eas_resource.default", "resource_name", fmt.Sprintf("tf-testacc%d", r)),
					resource.TestCheckResourceAttr("alicloud_pai_eas_resource.default", "charge_type", "PostPaid"),
				),
			},
			{
				Config: testAccAliCloudPaiEasResourceConfigUpdate(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_pai_eas_resource.default", "resource_name", fmt.Sprintf("tf-testacc%d-upd", r)),
					resource.TestCheckResourceAttr("alicloud_pai_eas_resource.default", "charge_type", "PostPaid"),
				),
			},
			{
				ResourceName:            "alicloud_pai_eas_resource.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"charge_type"},
			},
		},
	})
}

// testAccCheckPaiEasResourceDestroy verifies the EAS resource group is gone
// after the test. It re-reads the resource with DescribeResource and treats only
// a genuine not-found result as destroyed; permission errors, throttling or
// malformed responses surface as failures instead of being misread as deleted,
// so a broken cleanup cannot be masked as a pass.
func testAccCheckPaiEasResourceDestroy(s *terraform.State) error {
	client := testAccPaiEasResourceProvider.Meta().(*connectivity.AliyunClient)
	svc := eas.NewEasService(client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_pai_eas_resource" {
			continue
		}
		if _, err := svc.DescribeResource(rs.Primary.ID); err != nil {
			if eas.IsResourceNotFound(err) {
				continue
			}
			return fmt.Errorf("checking alicloud_pai_eas_resource %s after destroy: %w", rs.Primary.ID, err)
		}
		return fmt.Errorf("EAS resource %s still exists after destroy", rs.Primary.ID)
	}
	return nil
}

// testAccAliCloudPaiEasResourceConfig is the minimal create: only the required
// resource_name plus an explicit PostPaid charge_type. PostPaid avoids any
// subscription purchase or prepaid financial commitment.
func testAccAliCloudPaiEasResourceConfig(r int) string {
	return fmt.Sprintf(`
resource "alicloud_pai_eas_resource" "default" {
  resource_name = "tf-testacc%[1]d"
  charge_type   = "PostPaid"
}
`, r)
}

func testAccAliCloudPaiEasResourceConfigUpdate(r int) string {
	return fmt.Sprintf(`
resource "alicloud_pai_eas_resource" "default" {
  resource_name = "tf-testacc%[1]d-upd"
  charge_type   = "PostPaid"
}
`, r)
}
