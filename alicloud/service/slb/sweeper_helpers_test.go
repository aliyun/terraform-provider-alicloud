package slb_test

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	cs "github.com/denverdino/aliyungo/cs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// The slb load-balancer sweeper depends on cross-product cleanup helpers that
// in the flat layout live on the root alicloud package's VpcService, CsService
// and SlbService types (alicloud/service_alicloud_common_test.go,
// alicloud/service_alicloud_vpc.go and alicloud/service_alicloud_cs.go). The
// split slb_test package cannot import the root alicloud package, so these
// helpers are reproduced here against *connectivity.AliyunClient with the same
// API calls and control flow. Behavioural differences from the root originals:
//   - error wrapping uses errs.WrapError/NotFoundErr instead of the root
//     unexported message constants; only the NotFound classification and the
//     skip/keep/delete control flow need to match for the sweeper;
//   - describeCsKubernetes calls DescribeKubernetesClusterDetail once instead
//     of through the root NewInvoker retry wrapper; the sweeper only needs the
//     NotFound result and already retries the surrounding delete;
//   - needSweepVpc guards a nil describe result to avoid a nil-map access that
//     the root original would hit on a NotFound vswitch.

// slbTestVpcService mirrors the root VpcService sweeper helpers.
type slbTestVpcService struct {
	client *connectivity.AliyunClient
}

func (s slbTestVpcService) describeVswitch(id string) (map[string]interface{}, error) {
	action := "DescribeVSwitchAttributes"
	request := map[string]interface{}{
		"RegionId":  s.client.RegionId,
		"VSwitchId": id,
	}
	response, err := s.client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidVSwitchId.NotFound", "InvalidVswitchID.NotFound"}) {
			return nil, errs.NotFoundErr("Vswitch", id)
		}
		return nil, errs.WrapError(err)
	}
	logs.AddDebug(action, response, request)
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return nil, errs.WrapError(err)
	}
	object, _ := v.(map[string]interface{})
	if fmt.Sprint(object["VSwitchId"]) != id {
		return nil, errs.NotFoundErr("vswitch", id)
	}
	return object, nil
}

func (s slbTestVpcService) describeVpc(id string) (map[string]interface{}, error) {
	action := "DescribeVpcs"
	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
		"VpcId":    id,
	}
	response, err := s.client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"Forbidden.VpcNotFound", "InvalidVpcID.NotFound"}) {
			return nil, errs.NotFoundErr("Vpc", id)
		}
		return nil, errs.WrapError(err)
	}
	logs.AddDebug(action, response, request)
	v, err := jsonpath.Get("$.Vpcs.Vpc", response)
	if err != nil {
		return nil, errs.WrapError(err)
	}
	arr, _ := v.([]interface{})
	if len(arr) < 1 {
		return nil, errs.NotFoundErr("VPC", id)
	}
	object, _ := arr[0].(map[string]interface{})
	if fmt.Sprint(object["VpcId"]) != id {
		return nil, errs.NotFoundErr("VPC", id)
	}
	return object, nil
}

func (s slbTestVpcService) needSweepVpc(vpcId, vswitchId string) (bool, error) {
	if vpcId == "" && vswitchId != "" {
		object, err := s.describeVswitch(vswitchId)
		if err != nil && !errs.NotFoundError(err) {
			return false, errs.WrapError(err)
		}
		if object != nil {
			name := strings.ToLower(fmt.Sprint(object["VSwitchName"]))
			if strings.HasPrefix(name, "tf-testacc") || strings.HasPrefix(name, "tf_testacc") {
				log.Printf("[DEBUG] Need to sweep the vswitch (%v (%v)).", object["VSwitchId"], object["VSwitchName"])
				return true, nil
			}
			vpcId = fmt.Sprint(object["VpcId"])
		}
	}
	if vpcId != "" {
		object, err := s.describeVpc(vpcId)
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

// slbTestCsService mirrors the root CsService.DescribeCsKubernetes sweeper
// helper; only the NotFound classification matters to the caller.
type slbTestCsService struct {
	client *connectivity.AliyunClient
}

func (s slbTestCsService) describeCsKubernetes(id string) error {
	_, err := s.client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
		return csClient.DescribeKubernetesClusterDetail(id)
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"ErrorClusterNotFound"}) {
			return errs.NotFoundErr("CsKubernetes", id)
		}
		return errs.WrapError(err)
	}
	return nil
}

// slbSweepLoadBalancer mirrors the root SlbService.sweepSlb helper: turn off
// DeleteProtection then delete the load balancer, both with bounded retry.
func slbSweepLoadBalancer(client *connectivity.AliyunClient, id string) error {
	if id == "" {
		return nil
	}
	log.Printf("[DEBUG] Set SLB DeleteProtection to off before deleting %s ...", id)
	action := "SetLoadBalancerDeleteProtection"
	request := map[string]interface{}{
		"RegionId":         client.RegionId,
		"LoadBalancerId":   id,
		"DeleteProtection": "off",
	}
	wait := helper.IncrementalWait(3*time.Second, 10*time.Second)
	var err error
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		_, err = client.RpcPost("Slb", "2014-05-15", action, nil, request, false)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		log.Printf("[ERROR] Set SLB %s DeleteProtection to off failed.", id)
	}
	log.Printf("[DEBUG] Deleting SLB %s ...", id)
	delRequest := map[string]interface{}{
		"RegionId":       client.RegionId,
		"LoadBalancerId": id,
	}
	action = "DeleteLoadBalancer"
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		_, err = client.RpcPost("Slb", "2014-05-15", action, nil, delRequest, false)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	if err == nil {
		time.Sleep(1 * time.Second)
	}
	return errs.WrapError(err)
}
