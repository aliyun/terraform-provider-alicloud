package expressconnect_test

import (
	"fmt"
	"log"
	"strings"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
)

// sweep_vpc_helpers_test.go is a verbatim local replica of the root-package
// sweep-skip helpers used by testSweepRouterInterfaces:
//   - (*VpcService).needSweepVpc      (alicloud/service_alicloud_common_test.go)
//   - (*VpcService).DescribeVpc       (alicloud/service_alicloud_vpc.go)
//   - (*VpcService).DescribeVswitch   (alicloud/service_alicloud_vpc.go)
//
// The migrated external expressconnect_test package cannot import the root
// alicloud package (that would form the acctest -> alicloud ->
// provider/registry -> service/expressconnect cycle), so the exact bodies are
// re-declared here re-pointed at the leaf packages (errs/logs) that the root
// aliases resolve to (addDebug = logs.AddDebug; WrapError = errs.WrapError;
// NotFoundError = errs.NotFoundError; etc.). Behaviour is identical: the same
// Vpc DescribeVpcs / DescribeVSwitchAttributes calls, the same not-found
// handling and the same tf-test prefix matching used to decide whether a
// Router Interface should be swept when its name did not match the prefix.

func needSweepVpc(client *connectivity.AliyunClient, vpcId, vswitchId string) (bool, error) {
	if vpcId == "" && vswitchId != "" {
		object, err := describeVswitchForSweep(client, vswitchId)
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
		object, err := describeVpcForSweep(client, vpcId)
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

func describeVpcForSweep(client *connectivity.AliyunClient, id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "DescribeVpcs"
	request := map[string]interface{}{
		"RegionId": client.RegionId,
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

func describeVswitchForSweep(client *connectivity.AliyunClient, id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "DescribeVSwitchAttributes"
	request := map[string]interface{}{
		"RegionId":  client.RegionId,
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
