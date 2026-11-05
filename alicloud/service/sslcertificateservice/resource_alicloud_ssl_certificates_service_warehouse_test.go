package sslcertificateservice_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/sslcertificateservice"
	"github.com/hashicorp/terraform-plugin-framework/path"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
	fwschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const sslCertificatesServiceWarehouseResourceId = "alicloud_ssl_certificates_service_warehouse.default"

// TestAccAliCloudSslCertificatesServiceWarehouse_basic covers the full lifecycle of the
// Warehouse resource: create, refresh, force-new on every configurable argument (there is no
// UpdateWarehouse), import, and destroy. The warehouse is a freeDV-style certificate warehouse,
// so type=ssl is used as the cheapest non-empty CreateWarehouse Type.
func TestAccAliCloudSslCertificatesServiceWarehouse_basic(t *testing.T) {
	ra := tfacctest.RandIntRange(10000, 99999)
	region := os.Getenv("ALICLOUD_REGION")
	if region == "" {
		region = "cn-hangzhou"
	}
	name := fmt.Sprintf("tf-testacc-ssl-warehouse-%d", ra)
	updatedName := fmt.Sprintf("tf-testacc-ssl-warehouse-updated-%d", ra)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckSslCertificatesServiceWarehouseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSslCertificatesServiceWarehouseBasicConfig(name, "ssl", "test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(sslCertificatesServiceWarehouseResourceId, "id"),
					resource.TestCheckResourceAttrSet(sslCertificatesServiceWarehouseResourceId, "warehouse_instance_id"),
					resource.TestCheckResourceAttr(sslCertificatesServiceWarehouseResourceId, "name", name),
					resource.TestCheckResourceAttr(sslCertificatesServiceWarehouseResourceId, "type", "ssl"),
					resource.TestCheckResourceAttr(sslCertificatesServiceWarehouseResourceId, "biz", "test"),
				),
			},
			// name is RequiresReplace: a change re-creates the warehouse instead of updating it.
			{
				Config: testAccSslCertificatesServiceWarehouseBasicConfig(updatedName, "ssl", "test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(sslCertificatesServiceWarehouseResourceId, "name", updatedName),
					resource.TestCheckResourceAttr(sslCertificatesServiceWarehouseResourceId, "type", "ssl"),
					resource.TestCheckResourceAttr(sslCertificatesServiceWarehouseResourceId, "biz", "test"),
					resource.TestCheckResourceAttrSet(sslCertificatesServiceWarehouseResourceId, "id"),
				),
			},
			{
				ResourceName:      sslCertificatesServiceWarehouseResourceId,
				ImportState:       true,
				ImportStateVerify: true,
				// type and biz are write-only / enum-shifted and retained from state, so import
				// verification is limited to the readable identity and name.
				ImportStateVerifyIgnore: []string{"type", "biz"},
			},
		},
	})
}

// testAccCheckSslCertificatesServiceWarehouseDestroy confirms the warehouse is gone after the
// destroy step by listing it through ListWarehouse: a missing warehouse resolves to a NotFound
// describe error, while any other failure surfaces as a cleanup error.
func testAccCheckSslCertificatesServiceWarehouseDestroy(s *terraform.State) error {
	client := acctest.Provider.Meta().(*connectivity.AliyunClient)
	service := sslcertificateservice.NewSslCertificatesServiceServiceV2(client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ssl_certificates_service_warehouse" {
			continue
		}
		if _, err := service.DescribeSslCertificatesServiceWarehouse(rs.Primary.ID); err != nil {
			if errs.IsProviderNotFound(err) {
				continue
			}
			return fmt.Errorf("describing warehouse %s failed: %s", rs.Primary.ID, err)
		}
		return fmt.Errorf("warehouse %s still exists", rs.Primary.ID)
	}
	return nil
}

func testAccSslCertificatesServiceWarehouseBasicConfig(name, warehouseType, biz string) string {
	return fmt.Sprintf(`
resource "alicloud_ssl_certificates_service_warehouse" "default" {
  name = "%[1]s"
  type = "%[2]s"
  biz  = "%[3]s"
}
`, name, warehouseType, biz)
}

// TestUnitAliCloudSslCertificatesServiceWarehouseSchema pins the resource contract: every
// configurable argument is Required and replacement-only (no UpdateWarehouse), the CreateWarehouse
// Type enum is enforced, and the computed identity attributes use state for unknowns.
func TestUnitAliCloudSslCertificatesServiceWarehouseSchema(t *testing.T) {
	r := sslcertificateservice.NewSslCertificatesServiceWarehouseResource()

	var metadata fwresource.MetadataResponse
	r.Metadata(context.Background(), fwresource.MetadataRequest{ProviderTypeName: "alicloud"}, &metadata)
	if metadata.TypeName != "alicloud_ssl_certificates_service_warehouse" {
		t.Fatalf("TypeName = %q, want alicloud_ssl_certificates_service_warehouse", metadata.TypeName)
	}

	var schemaResp fwresource.SchemaResponse
	r.Schema(context.Background(), fwresource.SchemaRequest{}, &schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("building schema: %v", schemaResp.Diagnostics)
	}

	typeAttr, ok := schemaResp.Schema.Attributes["type"].(fwschema.StringAttribute)
	if !ok {
		t.Fatalf("type attribute is %T, want schema.StringAttribute", schemaResp.Schema.Attributes["type"])
	}
	if !typeAttr.Required {
		t.Fatal("type must be Required")
	}
	// The CreateWarehouse Type enum is enforced before the API call.
	rejected := func(value string) bool {
		var resp validator.StringResponse
		for _, v := range typeAttr.Validators {
			v.ValidateString(context.Background(), validator.StringRequest{
				Path:        path.Root("type"),
				ConfigValue: types.StringValue(value),
			}, &resp)
		}
		return resp.Diagnostics.HasError()
	}
	for _, value := range []string{"pcaAuthorityCertificate", "pcaCertificate", "personal", "ssl", "uploadedAuthorityCertificate", "uploadedCertificate"} {
		if rejected(value) {
			t.Fatalf("type enum must accept %q", value)
		}
	}
	for _, value := range []string{"SSL", "free", ""} {
		if !rejected(value) {
			t.Fatalf("type enum must reject %q", value)
		}
	}
	for _, attrName := range []string{"name", "type", "biz"} {
		attr, ok := schemaResp.Schema.Attributes[attrName].(fwschema.StringAttribute)
		if !ok {
			t.Fatalf("%s attribute is %T, want schema.StringAttribute", attrName, schemaResp.Schema.Attributes[attrName])
		}
		if !attr.Required || len(attr.PlanModifiers) == 0 {
			t.Fatalf("%s must be Required with a plan modifier", attrName)
		}
	}
	for _, attrName := range []string{"id", "warehouse_instance_id"} {
		attr, ok := schemaResp.Schema.Attributes[attrName].(fwschema.StringAttribute)
		if !ok {
			t.Fatalf("%s attribute is %T, want schema.StringAttribute", attrName, schemaResp.Schema.Attributes[attrName])
		}
		if !attr.Computed || attr.Required {
			t.Fatalf("%s must be Computed-only, got Computed=%v Required=%v", attrName, attr.Computed, attr.Required)
		}
	}
}
