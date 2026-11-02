package tsdb

import (
	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// TsdbService holds the connectivity client and encapsulates the tsdb (HiTSDB)
// service calls shared by the migrated tsdb resources and data sources. The
// tsdb-specific methods move here from the root HitsdbService; the lindorm
// methods remain on the root HitsdbService for the still-root lindorm types.
type TsdbService struct {
	client *connectivity.AliyunClient
}

// NewTsdbService constructs a TsdbService from the provider client.
func NewTsdbService(client *connectivity.AliyunClient) *TsdbService {
	return &TsdbService{client: client}
}

func (s *TsdbService) DescribeTsdbInstance(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeHiTSDBInstance"
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"InstanceId": id,
	}
	response, err = client.RpcPost("hitsdb", "2017-06-01", action, nil, request, true)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"Instance.IsNotValidError", "TSDB.Errorcode.InstanceDeleted", "TSDB.Errorcode.InstanceNotFound", "TSDB.Errorcode.ParameterInvaild"}) {
			err = errs.WrapErrorf(errs.NotFoundErr("TsdbInstance", id), errs.NotFoundMsg, errs.ProviderERROR)
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
	return object, nil
}

func (s *TsdbService) TsdbInstanceStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeTsdbInstance(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if object["Status"].(string) == failState {
				return object, object["Status"].(string), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, object["Status"].(string)))
			}
		}
		return object, object["Status"].(string), nil
	}
}

// describeVSwitch resolves a VPC vswitch to its attributes (VpcId, ZoneId, ...)
// via the public DescribeVSwitchAttributes API. It is a tsdb-local copy of the
// root VpcService.DescribeVSwitchWithTeadsl behavior so the tsdb package does
// not reverse-import the root alicloud package. Behavior matches the baseline.
func (s *TsdbService) describeVSwitch(id string) (object map[string]interface{}, err error) {
	client := s.client
	action := "DescribeVSwitchAttributes"
	request := map[string]interface{}{
		"RegionId":  s.client.RegionId,
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
