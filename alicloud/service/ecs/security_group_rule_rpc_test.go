package ecs_test

import (
	"io"
	"net/http"
	"sync/atomic"
	"testing"

	ecsservice "github.com/aliyun/terraform-provider-alicloud/alicloud/service/ecs"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
)

func TestUnitEcsSecurityGroupRuleFirstPageNotFound(t *testing.T) {
	var calls atomic.Int32
	client := reservedInstancesTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		if r.FormValue("Action") != "DescribeSecurityGroupAttribute" {
			t.Errorf("unexpected action %s", r.FormValue("Action"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		io.WriteString(w, `{"Code":"InvalidSecurityGroupId.NotFound","Message":"group was deleted","RequestId":"gone"}`)
	})
	_, err := ecsservice.NewService(client).DescribeSecurityGroupRule("sg-deleted:ingress:tcp:80/80:internet:10.0.0.0/8:accept:1")
	if !errs.NotFoundError(err) || calls.Load() != 1 {
		t.Fatalf("calls=%d, err=%v", calls.Load(), err)
	}
}
