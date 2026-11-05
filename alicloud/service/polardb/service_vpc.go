package polardb

import (
	"fmt"
	"log"
	"strings"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
)

// VpcService is a minimal local replica of the root-package VpcService
// (alicloud/service_alicloud_vpc.go) carrying only the Describe methods and
// the needSweepVpc sweeper helper that PolarDB resources (and their acceptance
// sweepers) call. The product package does not import other product packages
// or the root package; the original root definition is retained for unmigrated
// consumers. Behavior is unchanged.
type VpcService struct {
	client *connectivity.AliyunClient
}

// NewVpcService constructs a VpcService for the acceptance sweepers, which live
// in the external test package and therefore cannot populate the unexported
// client field via a struct literal. Behavior matches the root package's
// implicit construction.
func NewVpcService(client *connectivity.AliyunClient) *VpcService {
	return &VpcService{client: client}
}

func (s *VpcService) DescribeVSwitch(id string) (v vpc.DescribeVSwitchAttributesResponse, err error) {
	request := vpc.CreateDescribeVSwitchAttributesRequest()
	request.RegionId = s.client.RegionId
	request.VSwitchId = id

	invoker := NewInvoker()
	err = invoker.Run(func() error {
		raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeVSwitchAttributes(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"InvalidVswitchID.NotFound"}) {
				return WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
			}
			return WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*vpc.DescribeVSwitchAttributesResponse)
		if response.VSwitchId != id {
			return WrapErrorf(NotFoundErr("vswitch", id), NotFoundMsg, ProviderERROR)
		}
		v = *response
		return nil
	})
	return
}

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
		if IsExpectedErrors(err, []string{"InvalidVSwitchId.NotFound", "InvalidVswitchID.NotFound"}) {
			err = WrapErrorf(NotFoundErr("Vswitch", id), NotFoundWithError, err)
			return object, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return object, err
	}
	addDebug(action, response, request)
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	if fmt.Sprint(object["VSwitchId"]) != id {
		return object, WrapErrorf(NotFoundErr("vswitch", id), NotFoundWithResponse, response)
	}
	return object, nil
}

func (s *VpcService) DescribeVpc(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeVpcs"
	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
		"VpcId":    id,
	}
	response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
	if err != nil {
		if IsExpectedErrors(err, []string{"Forbidden.VpcNotFound", "InvalidVpcID.NotFound"}) {
			err = WrapErrorf(NotFoundErr("Vpc", id), NotFoundMsg, ProviderERROR)
			return object, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return object, err
	}
	addDebug(action, response, request)
	v, err := jsonpath.Get("$.Vpcs.Vpc", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Vpcs.Vpc", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(NotFoundErr("VPC", id), NotFoundWithResponse, response)
	} else {
		if v.([]interface{})[0].(map[string]interface{})["VpcId"].(string) != id {
			return object, WrapErrorf(NotFoundErr("VPC", id), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

// NeedSweepVpc decides whether a VPC/VSwitch created by a tf-testacc run should
// be swept. Migrated verbatim from the root test helper
// (alicloud/service_alicloud_common_test.go); product-own sweeper logic stays
// with the resource. Exported so the acceptance sweepers in the external test
// package can reach it; the root original remained unexported because its
// sweepers lived in the root internal test package.
func (s *VpcService) NeedSweepVpc(vpcId, vswitchId string) (bool, error) {
	if vpcId == "" && vswitchId != "" {
		object, err := s.DescribeVswitch(vswitchId)
		if err != nil && !NotFoundError(err) {
			return false, WrapError(err)
		}
		name := strings.ToLower(object["VSwitchName"].(string))
		if strings.HasPrefix(name, "tf-testacc") || strings.HasPrefix(name, "tf_testacc") {
			log.Printf("[DEBUG] Need to sweep the vswitch (%v (%v)).", object["VSwitchId"], object["VSwitchName"])
			return true, nil
		}
		vpcId = fmt.Sprint(object["VpcId"])
	}
	if vpcId != "" {
		object, err := s.DescribeVpc(vpcId)
		if err != nil {
			if NotFoundError(err) {
				return false, nil
			}
			return false, WrapError(err)
		}
		name := strings.ToLower(object["VpcName"].(string))
		if strings.HasPrefix(name, "tf-testacc") || strings.HasPrefix(name, "tf_testacc") {
			log.Printf("[DEBUG] Need to sweep the VPC (%v (%v)).", object["VpcId"], object["VpcName"])
			return true, nil
		}
	}
	return false, nil
}
