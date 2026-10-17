package coverage

import (
	"github.com/aliyun/terraform-provider-alicloud/scripts/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// FrameworkResource synthesizes an SDKv2 *schema.Resource from a
// framework-native resource's Schema() method, so the collection-order
// candidate and coverage checks apply to framework resources with the same
// fail-closed semantics as SDKv2 ones.
//
// The synthesis itself lives in scripts/internal/fwschema and is shared with
// scripts/testing, so the attribute coverage checker measures framework
// resources against their real schema too instead of skipping them as
// "registered framework-native". A parse failure, a missing Schema method or an
// unrecognized attribute/block kind is reported as an error: the caller
// fail-closes rather than silently treating a framework resource as having
// nothing to cover.
func FrameworkResource(implPath string) (*schema.Resource, error) {
	return fwschema.FrameworkResource(implPath)
}

// FrameworkResourceSource is the file-read-free core of FrameworkResource,
// exposed so tests can feed in-memory source.
func FrameworkResourceSource(implPath string, src []byte) (*schema.Resource, error) {
	return fwschema.FrameworkResourceSource(implPath, src)
}
