package polardb

import (
	"reflect"
	"testing"
)

// TestUnitFlattenPolarDBGatewayAttribute exercises the production flatten
// helpers (resource_gateway.go) directly from the package-internal test so the
// unit test guards the real flattening logic rather than a re-homed copy. The
// external acceptance test (resource_gateway_test.go, package polardb_test)
// cannot reach these unexported symbols. Behavior is unchanged.
func TestUnitFlattenPolarDBGatewayAttribute(t *testing.T) {
	endpoints, endpointsErr := flattenPolarDBGatewayEndpoints(map[string]interface{}{
		"Endpoint": []interface{}{
			map[string]interface{}{
				"Address": "gw.example.com", "EndpointId": "gwe-1", "GwClusterId": "gw-1",
				"Port": "3306", "TunnelId": "tunnel-1", "VpcId": "vpc-1", "NetType": "Private",
			},
			map[string]interface{}{
				"Address": "gw2.example.com", "EndpointId": "gwe-2", "GwClusterId": "gw-1",
				"Port": float64(8443), "TunnelId": "tunnel-2", "VpcId": "vpc-1", "NetType": "Public",
			},
		},
	})
	if endpointsErr != nil {
		t.Fatalf("flattenPolarDBGatewayEndpoints got an error: %v", endpointsErr)
	}
	wantEndpoints := []map[string]interface{}{
		{
			"address": "gw.example.com", "endpoint_id": "gwe-1", "gateway_id": "gw-1",
			"port": 3306, "tunnel_id": "tunnel-1", "vpc_id": "vpc-1", "network_type": "Private",
		},
		{
			"address": "gw2.example.com", "endpoint_id": "gwe-2", "gateway_id": "gw-1",
			"port": 8443, "tunnel_id": "tunnel-2", "vpc_id": "vpc-1", "network_type": "Public",
		},
	}
	if !reflect.DeepEqual(endpoints, wantEndpoints) {
		t.Fatalf("unexpected endpoints: %#v", endpoints)
	}

	securityIPArrays := flattenPolarDBGatewaySecurityIPArrays(map[string]interface{}{
		"SecurityIPArray": []interface{}{map[string]interface{}{
			"SecurityIPArrayName": "default", "SecurityIPList": "127.0.0.1",
		}},
	})
	wantSecurityIPArrays := []map[string]interface{}{{"name": "default", "ip_list": "127.0.0.1"}}
	if !reflect.DeepEqual(securityIPArrays, wantSecurityIPArrays) {
		t.Fatalf("unexpected security IP arrays: %#v", securityIPArrays)
	}
}
