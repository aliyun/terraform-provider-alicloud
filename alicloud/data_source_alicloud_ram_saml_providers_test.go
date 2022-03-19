package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

//  The test parameter encodedsaml_metadata_document should not be exposed
func SkipTestAccAlicloudRAMSamlProvidersDataSource(t *testing.T) {
	resourceId := "data.alicloud_ram_saml_providers.default"
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccSamlProviders%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceSamlProvidersDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex":     "${alicloud_ram_saml_provider.default.saml_provider_name}",
			"enable_details": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex":     "${alicloud_ram_saml_provider.default.saml_provider_name}-fake",
			"enable_details": "true",
		}),
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_ram_saml_provider.default.id}"},
			"enable_details": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_ram_saml_provider.default.id}-fake"},
			"enable_details": "true",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex":     name,
			"ids":            []string{"${alicloud_ram_saml_provider.default.id}"},
			"enable_details": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex":     name + "_fake",
			"ids":            []string{"${alicloud_ram_saml_provider.default.id}-fake"},
			"enable_details": "true",
		}),
	}
	var existRamSamlProvidersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                          "1",
			"ids.0":                          CHECKSET,
			"names.#":                        "1",
			"names.0":                        CHECKSET,
			"providers.#":                    "1",
			"providers.0.arn":                CHECKSET,
			"providers.0.description":        "For Terraform Test",
			"providers.0.saml_provider_name": name,
			"providers.0.update_date":        CHECKSET,
			"providers.0.id":                 CHECKSET,
			"providers.0.encodedsaml_metadata_document": CHECKSET,
		}
	}

	var fakeRamSamlProvidersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"names.#":     "0",
			"providers.#": "0",
		}
	}

	var RamSamlProvidersInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existRamSamlProvidersMapFunc,
		fakeMapFunc:  fakeRamSamlProvidersMapFunc,
	}

	RamSamlProvidersInfo.dataSourceTestCheck(t, 0, nameRegexConf, idsConf, allConf)
}

func dataSourceSamlProvidersDependence(name string) string {
	return fmt.Sprintf(`
    resource "alicloud_ram_saml_provider" "default" {
		saml_provider_name =  "%s"
		encodedsaml_metadata_document = "your encodedsaml metadata document"
		description = "For Terraform Test"
	}`, name)
}
