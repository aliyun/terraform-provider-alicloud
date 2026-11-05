package drds

import (
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/drds"
)

// TestUnitAliCloudDRDSInstanceFlattenVips guards the DRDS instance Read against a
// panic when the DescribeDrdsInstance response returns an empty VIP list. A valid
// instance can transiently report no VIPs, and indexing Vip[0] directly used to
// crash the provider. The empty case must return zero values; the populated case
// must surface the first VIP's VpcId plus the intranet connection string/port.
func TestUnitAliCloudDRDSInstanceFlattenVips(t *testing.T) {
	// Empty VIP list: must not panic and must yield zero values.
	vpcId, connectionString, port := flattenDrdsInstanceVips([]drds.Vip{})
	if vpcId != "" || connectionString != "" || port != "" {
		t.Fatalf("empty VIP list should yield zero values, got vpcId=%q connectionString=%q port=%q", vpcId, connectionString, port)
	}

	// Nil VIP list: same guarantee.
	vpcId, connectionString, port = flattenDrdsInstanceVips(nil)
	if vpcId != "" || connectionString != "" || port != "" {
		t.Fatalf("nil VIP list should yield zero values, got vpcId=%q connectionString=%q port=%q", vpcId, connectionString, port)
	}

	// Populated VIP list: vpc_id from the first VIP, connection_string/port from the intranet VIP.
	vips := []drds.Vip{
		{Type: "internet", VpcId: "vpc-external", Dns: "public.example.com", Port: "3306"},
		{Type: "intranet", VpcId: "vpc-internal", Dns: "intranet.example.com", Port: "3307"},
	}
	vpcId, connectionString, port = flattenDrdsInstanceVips(vips)
	if vpcId != "vpc-external" {
		t.Fatalf("vpcId should come from the first VIP, got %q", vpcId)
	}
	if connectionString != "intranet.example.com" {
		t.Fatalf("connectionString should come from the intranet VIP, got %q", connectionString)
	}
	if port != "3307" {
		t.Fatalf("port should come from the intranet VIP, got %q", port)
	}
}
