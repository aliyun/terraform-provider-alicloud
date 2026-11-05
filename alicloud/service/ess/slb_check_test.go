package ess_test

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// slbCheckService gives the ess_test ACC package a DescribeSlbLoadBalancer
// method that faithfully replicates root alicloud.SlbService.DescribeSlbLoadBalancer:
// the same DescribeLoadBalancerAttribute RPC, the same incremental-wait retry, and
// the same InvalidLoadBalancerId.NotFound -> NotFound mapping. The root alicloud
// package aggregates ess.ServicePackage(), so importing it from this package would
// form an import cycle; this test-scoped adapter keeps the ess migration
// product-isolated while preserving the original cross-product existence check
// used by TestAccAliCloudEssScalingGroup_slb. acctest.ResourceCheck.callDescribeMethod
// resolves the method name from the alicloud_slb_load_balancer resource id via
// reflection, so the method name and signature must match the root service exactly.
type slbCheckService struct {
	client *connectivity.AliyunClient
}

func (s *slbCheckService) DescribeSlbLoadBalancer(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeLoadBalancerAttribute"
	request := map[string]interface{}{
		"RegionId":       s.client.RegionId,
		"LoadBalancerId": id,
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Slb", "2014-05-15", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidLoadBalancerId.NotFound"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("SLB:LoadBalancer", id), errs.NotFoundMsg, errs.ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}
