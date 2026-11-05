package alikafka

import (
	"fmt"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
)

// VpcService is a local product-package copy of the root alicloud.VpcService
// (service_alicloud_vpc.go) limited to DescribeVswitch — the only VPC lookup
// the alikafka instance resource needs (resource_instance.go). The root type
// cannot be imported without a reverse dependency on the root package; this
// copy preserves the original behavior 1:1 (stock-method-disposition:
// necessary local copy, no behavior change). The root bare aliases are
// qualified to the errs/logs leaf packages they already resolve to.
type VpcService struct {
	client *connectivity.AliyunClient
}

// DescribeVswitch mirrors alicloud.VpcService.DescribeVswitch
// (service_alicloud_vpc.go:787). Identical control flow, error wrapping and
// not-found semantics.
func (s *VpcService) DescribeVswitch(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeVSwitchAttributes"
	request := map[string]interface{}{
		"RegionId":  s.client.RegionId,
		"VSwitchId": id,
	}
	response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidVSwitchId.NotFound", "InvalidVswitchID.NotFound"}) {
			err = errs.WrapErrorf(errs.NotFoundErr("Vswitch", id), errs.NotFoundWithError, err)
			return object, err
		}
		err = errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
		return object, err
	}
	logs.AddDebug(action, response, request)
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	if fmt.Sprint(object["VSwitchId"]) != id {
		return object, errs.WrapErrorf(errs.NotFoundErr("vswitch", id), errs.NotFoundWithResponse, response)
	}
	return object, nil
}
