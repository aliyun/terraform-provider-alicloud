package nas

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
)

// nasTestClientForRegion delegates to the public sweep.SharedClientForRegion so
// that internal package-nas TestUnit cases (*_unit_test.go) that call unexported
// CRUD share the same per-region endpoint cache (sweep.SharedEndpointCacheForRegion)
// as the baseline root-package sweeper, instead of a package-level global cache
// that could route one region's request to another region's cached endpoint.
//
// The sweep package only imports connectivity and does not reverse-import the
// root or nas package, so delegating here introduces no import cycle. Internal
// (package nas) and external (package nas_test) test files compile into separate
// binaries, so this helper does not collide with the same-named one in
// sweeper_test.go. The returned client is *connectivity.AliyunClient, which
// callers cast as needed; the original Skip-on-error and rawClient assertions
// are preserved because sweep.SharedClientForRegion returns conf.Client().
func nasTestClientForRegion(region string) (interface{}, error) {
	return sweep.SharedClientForRegion(region)
}
