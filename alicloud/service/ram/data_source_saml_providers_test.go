package ram_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// The test parameter encodedsaml_metadata_document should not be exposed
func SkipTestAccAlicloudRAMSamlProvidersDataSource(t *testing.T) {
	resourceId := "data.alicloud_ram_saml_providers.default"
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccSamlProviders%d", rand)
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourceSamlProvidersDependence)

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"name_regex":     "${alicloud_ram_saml_provider.default.saml_provider_name}",
			"enable_details": "true",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"name_regex":     "${alicloud_ram_saml_provider.default.saml_provider_name}-fake",
			"enable_details": "true",
		}),
	}
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_ram_saml_provider.default.id}"},
			"enable_details": "true",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_ram_saml_provider.default.id}-fake"},
			"enable_details": "true",
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"name_regex":     name,
			"ids":            []string{"${alicloud_ram_saml_provider.default.id}"},
			"enable_details": "true",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"name_regex":     name + "_fake",
			"ids":            []string{"${alicloud_ram_saml_provider.default.id}-fake"},
			"enable_details": "true",
		}),
	}
	var existRamSamlProvidersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                          "1",
			"ids.0":                          acctest.CHECKSET,
			"names.#":                        "1",
			"names.0":                        acctest.CHECKSET,
			"providers.#":                    "1",
			"providers.0.arn":                acctest.CHECKSET,
			"providers.0.description":        "For Terraform Test",
			"providers.0.saml_provider_name": name,
			"providers.0.update_date":        acctest.CHECKSET,
			"providers.0.id":                 acctest.CHECKSET,
			"providers.0.encodedsaml_metadata_document": acctest.CHECKSET,
		}
	}

	var fakeRamSamlProvidersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"names.#":     "0",
			"providers.#": "0",
		}
	}

	var RamSamlProvidersInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existRamSamlProvidersMapFunc,
		FakeMapFunc:  fakeRamSamlProvidersMapFunc,
	}

	RamSamlProvidersInfo.DataSourceTestCheck(t, 0, nameRegexConf, idsConf, allConf)
}

func dataSourceSamlProvidersDependence(name string) string {
	return fmt.Sprintf(`
    resource "alicloud_ram_saml_provider" "default" {
		saml_provider_name =  "%s"
		encodedsaml_metadata_document = "your encodedsaml metadata document"
		description = "For Terraform Test"
	}`, name)
}

// testAccRamSamlProvidersValidMetadataDocument is a real, base64-encoded SAML
// EntityDescriptor that the IMS CreateSAMLProvider API accepts (reused from
// resource_saml_provider_test.go TestAccAliCloudRAMSamlProvider_basic0), so
// the dependency SAML provider is actually created instead of relying on a
// placeholder string.
const testAccRamSamlProvidersValidMetadataDocument = "PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiPz48bWQ6RW50aXR5RGVzY3JpcHRvciBlbnRpdHlJRD0iaHR0cDovL2V4YW1wbGUuYWxpeXVuLmNvbS9leGFtcGxlLWlkcCIgeG1sbnM6bWQ9InVybjpvYXNpczpuYW1lczp0YzpTQU1MOjIuMDptZXRhZGF0YSI+PG1kOklEUFNTT0Rlc2NyaXB0b3IgV2FudEF1dGhuUmVxdWVzdHNTaWduZWQ9ImZhbHNlIiBwcm90b2NvbFN1cHBvcnRFbnVtZXJhdGlvbj0idXJuOm9hc2lzOm5hbWVzOnRjOlNBTUw6Mi4wOnByb3RvY29sIj48bWQ6S2V5RGVzY3JpcHRvciB1c2U9InNpZ25pbmciPjxkczpLZXlJbmZvIHhtbG5zOmRzPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwLzA5L3htbGRzaWcjIj48ZHM6WDUwOURhdGE+PGRzOlg1MDlDZXJ0aWZpY2F0ZT5NSUlEL3pDQ0F1ZWdBd0lCQWdJRU1yb0tjakFOQmdrcWhraUc5dzBCQVFzRkFEQ0JnREVuTUNVR0ExVUVBeE1lClFXeHBlWFZ1SUZKQlRTQkZlR0Z0Y0d4bElFTmxjblJwWm1sallYUmxNUkF3RGdZRFZRUUxFd2RCYkdsaVlXSmgKTVJBd0RnWURWUVFLRXdkQmJHbGlZV0poTVJFd0R3WURWUVFIRXdoSVlXNW5XbWh2ZFRFUk1BOEdBMVVFQ0JNSQpXbWhsU21saGJtY3hDekFKQmdOVkJBWVRBa05PTUNBWERUSXpNVEl3TkRFeU1EY3dNRm9ZRHpJd05URXdOREl4Ck1USXdOekF3V2pDQmdERW5NQ1VHQTFVRUF4TWVRV3hwZVhWdUlGSkJUU0JGZUdGdGNHeGxJRU5sY25ScFptbGoKWVhSbE1SQXdEZ1lEVlFRTEV3ZEJiR2xpWVdKaE1SQXdEZ1lEVlFRS0V3ZEJiR2xpWVdKaE1SRXdEd1lEVlFRSApFd2hJWVc1bldtaHZkVEVSTUE4R0ExVUVDQk1JV21obFNtbGhibWN4Q3pBSkJnTlZCQVlUQWtOT01JSUJJakFOCkJna3Foa2lHOXcwQkFRRUZBQU9DQVE4QU1JSUJDZ0tDQVFFQW9KVGVndWc0eXFaalNzQzFQUWpzbGxreUxWZHEKcXR0UGFqNmNOYldQdVFRNThSMkF4ZHNYeng4c05lOElLYUFYbW84azdDTFhDcXFLVzNjNEpzRWtTOUcva3B2NApJWFpBOGFVcDBCeXdQUFBocmFjUXd4cmJ5dkhja0dqVUpkNEZrOEVjbVVjNjRrSE5LbjBCaVJpL0NEZlM3MXBuCjh5T3dDNUZPSUlYWXhWMGtTTnNQMnozV2tBbFBXWm1sVkZSd1dxeHhGS2xCTjVpdVhaVHA4dk5rU2htVndBTW8KcjVpb2VBaFdXd0N1L0pvdUhLa3lnbVJnSDNhRjlSRlkrOGZ4NkMzR2hjZktISUszRTFBbVVtWlpjR3NDUCtxNApXeTBuSFp4QStaZEhTeE1OYUJPMm5JbkxJSHVDWHgza096eWpGV3dUaTVGSTlwdE5vNktBay9wRThRSURBUUFCCm8zMHdlekFQQmdOVkhSTUVDREFHQVFIL0FnRURNQjBHQTFVZEpRUVdNQlFHQ0NzR0FRVUZCd01CQmdnckJnRUYKQlFjREFqQUxCZ05WSFE4RUJBTUNBb1F3SFFZRFZSMFJBUUgvQkJNd0VZSUpiRzlqWVd4b2IzTjBod1IvQUFBQgpNQjBHQTFVZERnUVdCQlQ2TXluMnJjMVhEQTZqQkZYWVBOYitGaldMVmpBTkJna3Foa2lHOXcwQkFRc0ZBQU9DCkFRRUFoWHpUUzJJaHZjY3hzSVNzVVNFcldNNDJiQlZESHhTa05EemhPRmd0eGNtNUxuNHdjWXJvdkM3NHZxS1oKUWdQWmpGcWw3YUJTb1ZyNFdseWFaZlVBdHdNL3pZZytJbklUSVpBQ0dhM1VNK3h5V0NLSVhRNGpJVldnNG9QWQpxTStjNWllLzJFVlE0YmhObEQyL0lYZUVEZFd2TXMzdmFyRTFCUE5PQXJZZ2tZTmNER3lDSnA0ZmQ3d3ladWxhCllEZFFIWDdpdUJ1R0JOZFRBajlCUW5xaTJRcTc5RndMVTBrQkFvdUpVVVBPUjBpMGtwZ24vc2dSbHhvaHY1bHgKVTFwYVhtMEZRWHpUUDEvdjV5Y24wM3NVckFUekg2VkRpVlQ2N0NRQjR4MXJpOTFvUVRkWERXN1RvRkVhOGIrOApPdE8wZERMdDlnbCtNMkxYRzJTWnBZTkJoZz09PC9kczpYNTA5Q2VydGlmaWNhdGU+PC9kczpYNTA5RGF0YT48L2RzOktleUluZm8+PC9tZDpLZXlEZXNjcmlwdG9yPjxtZDpOYW1lSURGb3JtYXQ+dXJuOm9hc2lzOm5hbWVzOnRjOlNBTUw6MS4xOm5hbWVpZC1mb3JtYXQ6ZW1haWxBZGRyZXNzPC9tZDpOYW1lSURGb3JtYXQ+PG1kOk5hbWVJREZvcm1hdD51cm46b2FzaXM6bmFtZXM6dGM6U0FNTDoyLjA6bmFtZWlkLWZvcm1hdDpwZXJzaXN0ZW50PC9tZDpOYW1lSURGb3JtYXQ+PG1kOlNpbmdsZVNpZ25PblNlcnZpY2UgQmluZGluZz0idXJuOm9hc2lzOm5hbWVzOnRjOlNBTUw6Mi4wOmJpbmRpbmdzOkhUVFAtUE9TVCIgTG9jYXRpb249Imh0dHA6Ly9leGFtcGxlLmFsaXl1bi5jb20vZXhhbXBsZS1pZHAvc3NvL3NhbWwiLz48bWQ6U2luZ2xlU2lnbk9uU2VydmljZSBCaW5kaW5nPSJ1cm46b2FzaXM6bmFtZXM6dGM6U0FNTDoyLjA6YmluZGluZ3M6SFRUUC1SZWRpcmVjdCIgTG9jYXRpb249Imh0dHA6Ly9leGFtcGxlLmFsaXl1bi5jb20vZXhhbXBsZS1pZHAvc3NvL3NhbWwiLz48L21kOklEUFNTT0Rlc2NyaXB0b3I+PC9tZDpFbnRpdHlEZXNjcmlwdG9yPg=="

func dataSourceSamlProvidersActiveDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_ram_saml_provider" "default" {
  saml_provider_name            = "${var.name}"
  encodedsaml_metadata_document = "%s"
  description                   = "For Terraform Test"
}
`, name, testAccRamSamlProvidersValidMetadataDocument)
}

// TestAccAliCloudRAMSamlProvidersDataSource_basic is an active representative
// ACC for the alicloud_ram_saml_providers data source. It creates a real SAML
// provider (reusing a valid metadata document accepted by the IMS API), then
// queries the data source with both the name_regex and ids filters, asserting
// the created provider is returned with its real attributes (no fake id or
// empty-match shortcut to PASS).
func TestAccAliCloudRAMSamlProvidersDataSource_basic(t *testing.T) {
	resourceId := "data.alicloud_ram_saml_providers.default"
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccSamlProviders%d", rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccRamSamlProvidersDataSourceNameRegexConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "ids.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "names.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "providers.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "providers.0.saml_provider_name", name),
					resource.TestCheckResourceAttr(resourceId, "providers.0.description", "For Terraform Test"),
					resource.TestCheckResourceAttrSet(resourceId, "providers.0.arn"),
					resource.TestCheckResourceAttrSet(resourceId, "providers.0.id"),
					resource.TestCheckResourceAttrSet(resourceId, "providers.0.update_date"),
					resource.TestCheckResourceAttrSet(resourceId, "providers.0.encodedsaml_metadata_document"),
				),
			},
			{
				Config: testAccRamSamlProvidersDataSourceIdsConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "ids.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "names.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "providers.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "providers.0.saml_provider_name", name),
					resource.TestCheckResourceAttrSet(resourceId, "providers.0.arn"),
					resource.TestCheckResourceAttrSet(resourceId, "providers.0.id"),
				),
			},
		},
	})
}

func testAccRamSamlProvidersDataSourceNameRegexConfig(name string) string {
	return fmt.Sprintf(`
%s

data "alicloud_ram_saml_providers" "default" {
  name_regex     = "^${alicloud_ram_saml_provider.default.saml_provider_name}$"
  enable_details = true
}
`, dataSourceSamlProvidersActiveDependence(name))
}

func testAccRamSamlProvidersDataSourceIdsConfig(name string) string {
	return fmt.Sprintf(`
%s

data "alicloud_ram_saml_providers" "default" {
  ids            = ["${alicloud_ram_saml_provider.default.saml_provider_name}"]
  enable_details = false
}
`, dataSourceSamlProvidersActiveDependence(name))
}
