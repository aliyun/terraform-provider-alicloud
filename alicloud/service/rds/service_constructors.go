package rds

import "github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

// NewRdsService and NewRdsServiceV2 are exported constructors for the
// otherwise single-unexported-field service types. The RDS product migration
// moved the ACC tests into the external rds_test package, which cannot
// construct RdsService{client: ...} / RdsServiceV2{client: ...} directly
// because the client field is unexported. This mirrors the hologram.NewService
// precedent used by the other migrated product sub-packages and changes no
// behavior: the constructed value is identical to the positional literal.
func NewRdsService(client *connectivity.AliyunClient) *RdsService {
	return &RdsService{client: client}
}

func NewRdsServiceV2(client *connectivity.AliyunClient) *RdsServiceV2 {
	return &RdsServiceV2{client: client}
}

// NewVpcService mirrors the RdsService constructors for the VpcService type
// (whose client field is likewise unexported). The RDS ACC tests in the
// external rds_test package construct it during sweep and check helpers; this
// preserves the previous behavior identically.
func NewVpcService(client *connectivity.AliyunClient) *VpcService {
	return &VpcService{client: client}
}
