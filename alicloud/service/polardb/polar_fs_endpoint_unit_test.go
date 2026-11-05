package polardb

import (
	"reflect"
	"testing"
)

// TestUnitFlattenPolarDBPolarFsEndpointAddresses exercises the production
// flatten helper (resource_polar_fs_endpoint.go) directly from the
// package-internal test so the unit case guards the real flattening logic
// rather than a re-homed copy. The external acceptance test
// (resource_polar_fs_endpoint_test.go, package polardb_test) cannot reach this
// unexported symbol. Behavior is unchanged.
func TestUnitFlattenPolarDBPolarFsEndpointAddresses(t *testing.T) {
	raw := []interface{}{map[string]interface{}{
		"ConnectionString": "pfs.example.com", "PrivateZoneConnectionString": "pfs.internal",
		"IPAddress": "10.0.0.1", "Port": "443", "VPCId": "vpc-1", "VSwitchId": "vsw-1", "NetType": "Private",
	}}
	want := []map[string]interface{}{{
		"connection_string": "pfs.example.com", "private_zone_connection_string": "pfs.internal",
		"ip_address": "10.0.0.1", "port": "443", "vpc_id": "vpc-1", "vswitch_id": "vsw-1", "net_type": "Private",
	}}
	if got := flattenPolarDBPolarFsEndpointAddresses(raw); !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected endpoint addresses: %#v", got)
	}
}
