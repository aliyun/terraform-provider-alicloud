package alicloud

import "fmt"

// These OSS bucket test fixtures are retained at the root package for the
// provider-level test TestAccAlicloudProviderOss in provider_test.go, which
// depends on the package-internal providerCommon and root test helpers and
// therefore cannot move to alicloud/service/oss (package oss_test). The bucket
// resource test migrated with its own copies; these local copies keep the
// provider-level OSS scenario covered without importing a product test package.
func resourceOssBucketConfigDependence(name string) string {
	return fmt.Sprintf(`
resource "alicloud_oss_bucket" "target"{
  bucket = "%s-t"
}
`, name)
}

var ossBucketBasicMap = map[string]string{
	"creation_date":    CHECKSET,
	"lifecycle_rule.#": "0",
}
