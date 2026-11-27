package expressconnectrouter

import (
	"context"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
)

// SetTRAssociationStateDelayForTest shortens the initial status-poll delay so
// unit tests do not sleep for real. The returned function restores the value.
func SetTRAssociationStateDelayForTest(delay time.Duration) func() {
	previous := trAssociationStateDelay
	trAssociationStateDelay = delay
	return func() { trAssociationStateDelay = previous }
}

// SetTRAssociationNowForTest replaces the clock used by the operation budget so
// tests can move time forward without sleeping. The returned function restores
// the previous clock.
func SetTRAssociationNowForTest(now func() time.Time) func() {
	previous := trAssociationNow
	trAssociationNow = now
	return func() { trAssociationNow = previous }
}

func TRAssociationRPCPostForTest(ctx context.Context, client *connectivity.AliyunClient, product, version, action string, query, body map[string]interface{}) (map[string]interface{}, error) {
	return trAssociationRpcPost(ctx, client, product, version, action, query, body)
}
