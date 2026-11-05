package drds

import (
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/drds"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

type DrdsService struct {
	client *connectivity.AliyunClient
}

// NewService is the exported constructor used by external test packages
// (package drds_test) to build a DrdsService, since the client field is
// unexported. It mirrors the established migrated-package convention.
func NewService(client *connectivity.AliyunClient) *DrdsService {
	return &DrdsService{client: client}
}

func (s *DrdsService) DescribeDrdsInstance(id string) (*drds.DescribeDrdsInstanceResponse, error) {
	response := &drds.DescribeDrdsInstanceResponse{}
	request := drds.CreateDescribeDrdsInstanceRequest()
	request.RegionId = s.client.RegionId
	request.DrdsInstanceId = id
	raw, err := s.client.WithDrdsClient(func(drdsClient *drds.Client) (interface{}, error) {
		return drdsClient.DescribeDrdsInstance(request)
	})

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidDrdsInstanceId.NotFound"}) {
			return response, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return response, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ = raw.(*drds.DescribeDrdsInstanceResponse)
	if response.Data.Status == "5" {
		return response, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
	}
	return response, nil
}

func (s *DrdsService) DrdsInstanceStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeDrdsInstance(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if object.Data.Status == failState {
				return object, object.Data.Status, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, object.Data.Status))
			}
		}

		return object, object.Data.Status, nil
	}
}

func (s *DrdsService) WaitDrdsInstanceConfigEffect(id string, item map[string]string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for {
		effected := false
		object, err := s.DescribeDrdsInstance(id)

		if err != nil {
			if errs.NotFoundError(err) {
				return errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
			}
			return errs.WrapError(err)
		}

		if value, ok := item["description"]; ok {
			if object.Data.Description == value {
				effected = true
			}
		}

		if effected {
			break
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, object.Data, item, errs.ProviderERROR)
		}
		time.Sleep(helper.DefaultIntervalShort * time.Second)
	}

	return nil
}

// describeVSwitch is a local 1:1 reproduction of root VpcService.DescribeVSwitch.
// VpcService has not been migrated to a leaf package; the product package cannot
// import root alicloud or another product package, so the minimal VSwitch lookup is
// reproduced here using only leaf packages (connectivity, helper, errs, logs) to
// preserve the Create behavior that resolves VpcId/ZoneId from a VSwitchId.
func (s *DrdsService) describeVSwitch(id string) (v vpc.DescribeVSwitchAttributesResponse, err error) {
	request := vpc.CreateDescribeVSwitchAttributesRequest()
	request.RegionId = s.client.RegionId
	request.VSwitchId = id

	invoker := helper.NewInvoker()
	err = invoker.Run(func() error {
		raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeVSwitchAttributes(request)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"InvalidVswitchID.NotFound"}) {
				return errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
			}
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*vpc.DescribeVSwitchAttributesResponse)
		if response.VSwitchId != id {
			return errs.WrapErrorf(errs.NotFoundErr("vswitch", id), errs.NotFoundMsg, errs.ProviderERROR)
		}
		v = *response
		return nil
	})
	return
}
