package gpdb

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
)

// describeVSwitchWithTeadsl looks up a VSwitch via the Vpc DescribeVSwitchAttributes
// RPC. It is a 1:1 local copy of the root alicloud package's
// VpcService.DescribeVSwitchWithTeadsl, kept here because a service package may
// not import the root package or another product's service package. Behavior is
// preserved, including the InvalidVswitchID.NotFound mapping and the response
// VSwitchId guard.
func describeVSwitchWithTeadsl(client *connectivity.AliyunClient, id string) (map[string]interface{}, error) {
	action := "DescribeVSwitchAttributes"
	request := map[string]interface{}{
		"RegionId":  client.RegionId,
		"VSwitchId": id,
	}

	response, err := client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidVswitchID.NotFound"}) {
			return nil, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return nil, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(action, response, request)
	if v, ok := response["VSwitchId"].(string); ok && v != id {
		return nil, errs.WrapErrorf(errs.NotFoundErr("vswitch", id), errs.NotFoundMsg, errs.ProviderERROR)
	}
	return response, nil
}
