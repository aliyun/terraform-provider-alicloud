package ess

import (
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
)

// describeSlb mirrors alicloud.SlbService.DescribeSlb. The slb product
// subpackage is not migrated, and ess must not import the flat root alicloud
// package (it would form an import cycle via provider/registry), so the needed
// cross-product method logic is copied here as a package-local helper.
// See provider-resource-standards/references/stock-method-disposition.md.
func describeSlb(client *connectivity.AliyunClient, id string) (*slb.DescribeLoadBalancerAttributeResponse, error) {
	response := &slb.DescribeLoadBalancerAttributeResponse{}
	request := slb.CreateDescribeLoadBalancerAttributeRequest()
	request.RegionId = client.RegionId
	request.LoadBalancerId = id
	raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.DescribeLoadBalancerAttribute(request)
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidLoadBalancerId.NotFound"}) {
			err = errs.WrapErrorf(errs.NotFoundErr("Slb", id), errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		} else {
			err = errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}
		return response, err
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest)
	response, _ = raw.(*slb.DescribeLoadBalancerAttributeResponse)
	if response.LoadBalancerId == "" {
		err = errs.WrapErrorf(errs.NotFoundErr("Slb", id), errs.NotFoundMsg, errs.ProviderERROR)
	}
	return response, err
}

// waitForSlb mirrors alicloud.SlbService.WaitForSlb. Copied locally for the
// same reason as describeSlb. Uses the ess-local Status and Active/Deleted
// constants from extension.go and helper.DefaultIntervalShort.
func waitForSlb(client *connectivity.AliyunClient, id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := describeSlb(client, id)

		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errs.WrapError(err)
			}
		} else if strings.ToLower(object.LoadBalancerStatus) == strings.ToLower(string(status)) {
			//TODO
			break
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, object.LoadBalancerStatus, status, errs.ProviderERROR)
		}
		time.Sleep(helper.DefaultIntervalShort * time.Second)
	}
	return nil
}
