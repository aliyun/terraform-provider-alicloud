package gpdb_test

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// needSweepVpc replicates the root alicloud package's VpcService.needSweepVpc
// heuristic used by the elastic_instance sweeper: a GPDB instance is swept when
// its own description misses the test prefix but its VPC/vswitch name carries
// it. The flat-layout original lives in the root package and cannot be imported
// from this external test package, so the subset used here is kept local.
//
// The describe helpers below mirror VpcService.DescribeVpc/DescribeVswitch:
// they convert the SDK not-found responses into errs typed not-found errors and
// re-verify that the returned VpcId/VSwitchId equals the requested id, so
// needSweepVpc classifies a vanished VPC/vswitch via errs.NotFoundError exactly
// as the original common_test.go path did, instead of surfacing it as a hard
// error that aborts the sweep.
func needSweepVpc(client *connectivity.AliyunClient, vpcId, vswitchId string) (bool, error) {
	if vpcId == "" && vswitchId != "" {
		object, err := sweepDescribeVswitch(client, vswitchId)
		if err != nil {
			if !errs.NotFoundError(err) {
				return false, errs.WrapError(err)
			}
		} else {
			name := strings.ToLower(fmt.Sprint(object["VSwitchName"]))
			if strings.HasPrefix(name, "tf-testacc") || strings.HasPrefix(name, "tf_testacc") {
				log.Printf("[DEBUG] Need to sweep the vswitch (%v (%v)).", object["VSwitchId"], object["VSwitchName"])
				return true, nil
			}
			vpcId = fmt.Sprint(object["VpcId"])
		}
	}
	if vpcId != "" {
		object, err := sweepDescribeVpc(client, vpcId)
		if err != nil {
			if errs.NotFoundError(err) {
				return false, nil
			}
			return false, errs.WrapError(err)
		}
		name := strings.ToLower(fmt.Sprint(object["VpcName"]))
		if strings.HasPrefix(name, "tf-testacc") || strings.HasPrefix(name, "tf_testacc") {
			log.Printf("[DEBUG] Need to sweep the VPC (%v (%v)).", object["VpcId"], object["VpcName"])
			return true, nil
		}
	}
	return false, nil
}

func sweepDescribeVswitch(client *connectivity.AliyunClient, id string) (map[string]interface{}, error) {
	action := "DescribeVSwitchAttributes"
	request := map[string]interface{}{
		"RegionId":  client.RegionId,
		"VSwitchId": id,
	}
	response, err := client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidVSwitchId.NotFound", "InvalidVswitchID.NotFound"}) {
			return nil, errs.WrapErrorf(errs.NotFoundErr("Vswitch", id), errs.NotFoundWithError, err)
		}
		return nil, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(action, response, request)
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return nil, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$", response)
	}
	object := v.(map[string]interface{})
	if fmt.Sprint(object["VSwitchId"]) != id {
		return nil, errs.WrapErrorf(errs.NotFoundErr("vswitch", id), errs.NotFoundWithResponse, response)
	}
	return object, nil
}

func sweepDescribeVpc(client *connectivity.AliyunClient, id string) (map[string]interface{}, error) {
	action := "DescribeVpcs"
	request := map[string]interface{}{
		"RegionId": client.RegionId,
		"VpcId":    id,
	}
	response, err := client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"Forbidden.VpcNotFound", "InvalidVpcID.NotFound"}) {
			return nil, errs.WrapErrorf(errs.NotFoundErr("Vpc", id), errs.NotFoundMsg, errs.ProviderERROR)
		}
		return nil, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(action, response, request)
	v, err := jsonpath.Get("$.Vpcs.Vpc", response)
	if err != nil {
		return nil, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Vpcs.Vpc", response)
	}
	items := v.([]interface{})
	if len(items) < 1 {
		return nil, errs.WrapErrorf(errs.NotFoundErr("VPC", id), errs.NotFoundWithResponse, response)
	}
	object := items[0].(map[string]interface{})
	returnedID, ok := object["VpcId"].(string)
	if !ok || returnedID != id {
		return nil, errs.WrapErrorf(errs.NotFoundErr("VPC", id), errs.NotFoundWithResponse, response)
	}
	return object, nil
}

// TestMain wires the sweepers registered in this external test package through
// the same helper/resource registry used by the flat-layout original; without it
// the AddTestSweepers registrations in gpdb_test are inert. It only delegates to
// resource.TestMain and does not introduce a second registry or run sweeps itself.
func TestMain(m *testing.M) {
	resource.TestMain(m)
}
