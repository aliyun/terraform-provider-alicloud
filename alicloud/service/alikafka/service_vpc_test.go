package alikafka

import (
	"fmt"
	"log"
	"strings"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
)

// DescribeVpc mirrors alicloud.VpcService.DescribeVpc
// (service_alicloud_vpc.go:52). 1:1 copy of the root method; the root bare
// aliases are qualified to the errs/logs leaf packages they resolve to
// (NotFoundMsg, ProviderERROR, NotFoundWithError, IsExpectedErrors, WrapErrorf,
// NotFoundErr, DefaultErrorMsg, AlibabaCloudSdkGoERROR, addDebug,
// FailedGetAttributeMsg, NotFoundWithResponse all resolve to the errs/logs
// leaf packages, exactly as the migrated DescribeVswitch in service_vpc.go).
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
		if errs.IsExpectedErrors(err, []string{"Forbidden.VpcNotFound", "InvalidVpcID.NotFound"}) {
			err = errs.WrapErrorf(errs.NotFoundErr("Vpc", id), errs.NotFoundMsg, errs.ProviderERROR)
			return object, err
		}
		err = errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
		return object, err
	}
	logs.AddDebug(action, response, request)
	v, err := jsonpath.Get("$.Vpcs.Vpc", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Vpcs.Vpc", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("VPC", id), errs.NotFoundWithResponse, response)
	} else {
		if v.([]interface{})[0].(map[string]interface{})["VpcId"].(string) != id {
			return object, errs.WrapErrorf(errs.NotFoundErr("VPC", id), errs.NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

// needSweepVpc mirrors alicloud.VpcService.needSweepVpc
// (service_alicloud_common_test.go:722). The root bare aliases NotFoundError
// and WrapError resolve to errs.NotFoundError and errs.WrapError.
func (s *VpcService) needSweepVpc(vpcId, vswitchId string) (bool, error) {
	if vpcId == "" && vswitchId != "" {
		object, err := s.DescribeVswitch(vswitchId)
		if err != nil && !errs.NotFoundError(err) {
			return false, errs.WrapError(err)
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
			if errs.NotFoundError(err) {
				return false, nil
			}
			return false, errs.WrapError(err)
		}
		name := strings.ToLower(object["VpcName"].(string))
		if strings.HasPrefix(name, "tf-testacc") || strings.HasPrefix(name, "tf_testacc") {
			log.Printf("[DEBUG] Need to sweep the VPC (%v (%v)).", object["VpcId"], object["VpcName"])
			return true, nil
		}
	}
	return false, nil
}
