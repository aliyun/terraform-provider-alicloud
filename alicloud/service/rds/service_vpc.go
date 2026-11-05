package rds

import (
	"fmt"
	"log"
	"strings"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
)

// VpcService is the package-rds copy of the root alicloud/service_alicloud_vpc.go
// type. RDS only calls DescribeVSwitch, so only that method is copied here. It
// is a cross-product service type; per stock-method-disposition rule #4 the
// root keeps its full type and RDS keeps this behavior-preserving local copy
// (with only the methods RDS uses) rather than importing the root package.
// The body is verbatim from the root; NewInvoker/IsExpectedErrors/WrapErrorf/
// addDebug/NotFoundErr and the error-code constants resolve via aliases.go.

type VpcService struct {
	client *connectivity.AliyunClient
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

// DescribeVswitch and DescribeVpc are the map-returning VPC attribute readers
// copied verbatim from the root alicloud/service_alicloud_vpc.go (rule #4
// cross-product local copy). NeedSweepVpc is the exported wrapper for the root
// needSweepVpc test helper; the RDS ACC sweep/check helper in the external
// rds_test package calls vpcService.NeedSweepVpc. Bodies are verbatim from the
// root; IsExpectedErrors/WrapErrorf/NotFoundErr/addDebug/NotFoundError/WrapError
// resolve via aliases.go.
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

// NeedSweepVpc is the exported wrapper for the root needSweepVpc helper; the
// body is verbatim from alicloud/service_alicloud_common_test.go. Exposed for
// the external rds_test ACC sweep/check helper.
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
