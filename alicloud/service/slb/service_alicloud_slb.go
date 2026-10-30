package slb

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Service struct {
	client *connectivity.AliyunClient
}

// NewService is the exported constructor: the client field stays unexported,
// so callers outside the package (including the acceptance tests) build values
// with it.
func NewService(client *connectivity.AliyunClient) *Service {
	return &Service{client: client}
}

// Status mirrors the root package status type used by the WaitFor helpers; the
// root package cannot be imported from a service package, so it is redefined
// here together with the status constants the SLB wait helpers rely on.
type Status string

const (
	Running = Status("Running")
	Stopped = Status("Stopped")
	Deleted = Status("Deleted")
	Null    = Status("Null")
)

// DefaultTimeout and DefaultTimeoutMedium mirror the root package timeout
// constants used as the default wait durations by the SLB resources.
const (
	DefaultTimeout       = 120
	DefaultTimeoutMedium = 500
)

type SlbTag struct {
	TagKey   string
	TagValue string
}

const max_num_per_time = 50
const tags_max_num_per_time = 5
const tags_max_page_size = 50

func (s *Service) BuildSlbCommonRequest() (*requests.CommonRequest, error) {
	// Get product code from the built request
	slbReq := slb.CreateCreateLoadBalancerRequest()
	req, err := s.client.NewCommonRequest(slbReq.GetProduct(), slbReq.GetLocationServiceCode(), strings.ToUpper(string(helper.Https)), connectivity.ApiVersion20140515)
	if err != nil {
		return req, errs.WrapError(err)
	}
	req.RegionId = s.client.RegionId
	return req, err
}

// DescribeVSwitchWithTeadsl looks up a VSwitch via the VPC OpenAPI. It is a
// local copy of the root VpcService helper so this package does not import
// the unmigrated root package.
func (s *Service) DescribeVSwitchWithTeadsl(id string) (object map[string]interface{}, err error) {
	client := s.client
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

func (s *Service) DescribeSlb(id string) (*slb.DescribeLoadBalancerAttributeResponse, error) {
	response := &slb.DescribeLoadBalancerAttributeResponse{}
	request := slb.CreateDescribeLoadBalancerAttributeRequest()
	request.RegionId = s.client.RegionId
	request.LoadBalancerId = id
	raw, err := s.client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.DescribeLoadBalancerAttribute(request)
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidLoadBalancerId.NotFound"}) {
			err = errs.WrapErrorf(errs.NotFoundErr("Slb", id), errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		} else {
			err = errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}
		return response, err
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest)
	response, _ = raw.(*slb.DescribeLoadBalancerAttributeResponse)
	if response.LoadBalancerId == "" {
		err = errs.WrapErrorf(errs.NotFoundErr("Slb", id), errs.NotFoundMsg, errs.ProviderERROR)
	}
	return response, err
}

func (s *Service) DescribeSlbRule(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeRuleAttribute"
	request := map[string]interface{}{
		"RuleId":   id,
		"RegionId": s.client.RegionId,
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcGet("Slb", "2014-05-15", action, request, nil)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidRuleId.NotFound"}) {
			return response, errs.WrapErrorf(errs.NotFoundErr("SlbRule", id), errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	if fmt.Sprint(response["RuleId"]) != id {
		return response, errs.WrapErrorf(errs.NotFoundErr("SlbRule", id), errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
	}
	return response, nil
}

func (s *Service) DescribeSlbServerGroup(id string) (*slb.DescribeVServerGroupAttributeResponse, error) {
	response := &slb.DescribeVServerGroupAttributeResponse{}
	request := slb.CreateDescribeVServerGroupAttributeRequest()
	request.RegionId = s.client.RegionId
	request.VServerGroupId = id

	raw, err := s.client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.DescribeVServerGroupAttribute(request)
	})
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"The specified VServerGroupId does not exist", "InvalidParameter"}) {
			return response, errs.WrapErrorf(errs.NotFoundErr("Slb:ServerGroup", id), errs.NotFoundMsg, errs.ProviderERROR, response.RequestId)
		}
		return response, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	response, _ = raw.(*slb.DescribeVServerGroupAttributeResponse)
	if response.VServerGroupId == "" {
		return response, errs.WrapErrorf(errs.NotFoundErr("Slb:ServerGroup", id), errs.NotFoundMsg, errs.ProviderERROR, response.RequestId)
	}

	return response, err
}

func (s *Service) DescribeSlbMasterSlaveServerGroup(id string) (*slb.DescribeMasterSlaveServerGroupAttributeResponse, error) {
	response := &slb.DescribeMasterSlaveServerGroupAttributeResponse{}
	request := slb.CreateDescribeMasterSlaveServerGroupAttributeRequest()
	request.RegionId = s.client.RegionId
	request.MasterSlaveServerGroupId = id
	raw, err := s.client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.DescribeMasterSlaveServerGroupAttribute(request)
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"The specified MasterSlaveGroupId does not exist", "InvalidParameter"}) {
			return response, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return response, errs.WrapErrorf(err, logs.DefaultDebugMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ = raw.(*slb.DescribeMasterSlaveServerGroupAttributeResponse)
	if response.MasterSlaveServerGroupId == "" {
		return response, errs.WrapErrorf(errs.NotFoundErr("SlbMasterSlaveServerGroup", id), errs.NotFoundMsg, errs.ProviderERROR)
	}
	return response, err
}

func (s *Service) DescribeSlbBackendServer(id string) (*slb.DescribeLoadBalancerAttributeResponse, error) {
	response := &slb.DescribeLoadBalancerAttributeResponse{}
	request := slb.CreateDescribeLoadBalancerAttributeRequest()
	request.RegionId = s.client.RegionId
	request.LoadBalancerId = id
	raw, err := s.client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.DescribeLoadBalancerAttribute(request)
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidLoadBalancerId.NotFound"}) {
			err = errs.WrapErrorf(errs.NotFoundErr("SlbBackendServers", id), errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		} else {
			err = errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}
		return response, err
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ = raw.(*slb.DescribeLoadBalancerAttributeResponse)
	if response.LoadBalancerId == "" {
		err = errs.WrapErrorf(errs.NotFoundErr("SlbBackendServers", id), errs.NotFoundMsg, errs.ProviderERROR)
	}
	return response, err
}

func (s *Service) DescribeSlbListener(id string) (listener map[string]interface{}, err error) {
	parts, err := ParseSlbListenerId(id)
	if err != nil {
		return nil, errs.WrapError(err)
	}
	protocol := parts[1]
	request, err := s.BuildSlbCommonRequest()
	request.RegionId = s.client.RegionId
	if err != nil {
		err = errs.WrapError(err)
		return
	}
	request.ApiName = fmt.Sprintf("DescribeLoadBalancer%sListenerAttribute", strings.ToUpper(string(protocol)))
	request.QueryParams["LoadBalancerId"] = parts[0]
	port, _ := strconv.Atoi(parts[2])
	request.QueryParams["ListenerPort"] = string(requests.NewInteger(port))
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		raw, err := s.client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.ProcessCommonRequest(request)
		})

		if err != nil {
			if errs.IsExpectedErrors(err, []string{"The specified resource does not exist"}) {
				return retry.NonRetryableError(errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR))
			} else if errs.IsExpectedErrors(err, errs.SlbIsBusy) {
				return retry.RetryableError(errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR))
			}
			return retry.NonRetryableError(errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR))
		}
		logs.AddDebug(request.GetActionName(), raw, request, request.QueryParams)
		response, _ := raw.(*responses.CommonResponse)
		if err = json.Unmarshal(response.GetHttpContentBytes(), &listener); err != nil {
			return retry.NonRetryableError(errs.WrapError(err))
		}
		if port, ok := listener["ListenerPort"]; ok && port.(float64) > 0 {
			return nil
		} else {
			return retry.RetryableError(errs.WrapErrorf(errs.NotFoundErr("SlbListener", id), errs.NotFoundMsg, errs.ProviderERROR))
		}
	})

	return
}

func (s *Service) DescribeSlbAcl(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeAccessControlListAttribute"
	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
		"AclId":    id,
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Slb", "2014-05-15", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"AclNotExist"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("SLB:Acl", id), errs.NotFoundMsg, errs.ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *Service) WaitForSlb(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeSlb(id)

		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errs.WrapError(err)
			}
		} else if strings.ToLower(object.LoadBalancerStatus) == strings.ToLower(string(status)) {
			//TODO
			break
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, object.LoadBalancerStatus, status, errs.ProviderERROR)
		}
		time.Sleep(helper.DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *Service) WaitForSlbListener(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeSlbListener(id)
		if err != nil && !errs.IsExpectedErrors(err, []string{"InvalidLoadBalancerId.NotFound"}) {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errs.WrapError(err)
			}
		}
		gotStatus := ""
		if value, ok := object["Status"]; ok {
			gotStatus = strings.ToLower(value.(string))
		}
		if gotStatus == strings.ToLower(string(status)) {
			//TODO
			break
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, gotStatus, status, errs.ProviderERROR)
		}
		time.Sleep(helper.DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *Service) WaitForSlbRule(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeSlbRule(id)

		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errs.WrapError(err)
			}
		}
		if object["RuleId"] == id && status != Deleted {
			break
		}

		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, "", id, errs.ProviderERROR)
		}
		time.Sleep(helper.DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *Service) WaitForSlbServerGroup(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeSlbServerGroup(id)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errs.WrapError(err)
			}
		}
		if object.VServerGroupId == id {
			break
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, object.VServerGroupId, id, errs.ProviderERROR)
		}
		time.Sleep(helper.DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *Service) WaitForSlbMasterSlaveServerGroup(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeSlbMasterSlaveServerGroup(id)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errs.WrapError(err)
			}
		}
		if object.MasterSlaveServerGroupId == id && status != Deleted {
			break
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, object.MasterSlaveServerGroupId, id, errs.ProviderERROR)
		}
		time.Sleep(helper.DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *Service) WaitSlbAttribute(id string, instanceSet *schema.Set, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

RETRY:
	object, err := s.DescribeSlb(id)
	if err != nil {
		if errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapError(err)
	}
	if time.Now().After(deadline) {
		return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, Null, id, errs.ProviderERROR)
	}
	servers := object.BackendServers.BackendServer
	if len(servers) > 0 {
		for _, s := range servers {
			if instanceSet.Contains(s.ServerId) {
				goto RETRY
			}
		}
	}
	return nil
}

func (s *Service) slbRemoveAccessControlListEntryPerTime(list []interface{}, id string) error {
	request := slb.CreateRemoveAccessControlListEntryRequest()
	request.RegionId = s.client.RegionId
	request.AclId = id
	b, err := json.Marshal(list)
	if err != nil {
		return errs.WrapError(err)
	}
	request.AclEntrys = string(b)
	raw, err := s.client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.RemoveAccessControlListEntry(request)
	})
	if err != nil {
		if !errs.IsExpectedErrors(err, []string{"AclEntryEmpty"}) {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return nil
}

func (s *Service) SlbRemoveAccessControlListEntry(list []interface{}, aclId string) error {
	num := len(list)

	if num <= 0 {
		return nil
	}

	t := (num + max_num_per_time - 1) / max_num_per_time
	for i := 0; i < t; i++ {
		start := i * max_num_per_time
		end := (i + 1) * max_num_per_time

		if end > num {
			end = num
		}

		slice := list[start:end]
		if err := s.slbRemoveAccessControlListEntryPerTime(slice, aclId); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) slbAddAccessControlListEntryPerTime(list []interface{}, id string) error {
	request := slb.CreateAddAccessControlListEntryRequest()
	request.RegionId = s.client.RegionId
	request.AclId = id
	b, err := json.Marshal(list)
	if err != nil {
		return errs.WrapError(err)
	}
	request.AclEntrys = string(b)
	raw, err := s.client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.AddAccessControlListEntry(request)
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return nil
}

func (s *Service) SlbAddAccessControlListEntry(list []interface{}, aclId string) error {
	num := len(list)

	if num <= 0 {
		return nil
	}

	t := (num + max_num_per_time - 1) / max_num_per_time
	for i := 0; i < t; i++ {
		start := i * max_num_per_time
		end := (i + 1) * max_num_per_time

		if end > num {
			end = num
		}
		slice := list[start:end]
		if err := s.slbAddAccessControlListEntryPerTime(slice, aclId); err != nil {
			return err
		}
	}

	return nil
}

// Flattens an array of slb.AclEntry into a []map[string]string
func (s *Service) FlattenSlbAclEntryMappings(list []slb.AclEntry) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(list))

	for _, i := range list {
		l := map[string]interface{}{
			"entry":   i.AclEntryIP,
			"comment": i.AclEntryComment,
		}
		result = append(result, l)
	}

	return result
}

// Flattens an array of slb.AclEntry into a []map[string]string
func (s *Service) flattenSlbRelatedListenerMappings(list []slb.RelatedListener) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(list))

	for _, i := range list {
		l := map[string]interface{}{
			"load_balancer_id": i.LoadBalancerId,
			"protocol":         i.Protocol,
			"frontend_port":    i.ListenerPort,
			"acl_type":         i.AclType,
		}
		result = append(result, l)
	}

	return result
}

func (s *Service) DescribeSlbCACertificate(id string) (*slb.CACertificate, error) {
	certificate := &slb.CACertificate{}
	request := slb.CreateDescribeCACertificatesRequest()
	request.RegionId = s.client.RegionId
	request.CACertificateId = id
	raw, err := s.client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.DescribeCACertificates(request)
	})
	if err != nil {
		return certificate, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*slb.DescribeCACertificatesResponse)
	if len(response.CACertificates.CACertificate) < 1 {
		return certificate, errs.WrapErrorf(errs.NotFoundErr("SlbCACertificate", id), errs.NotFoundMsg, errs.ProviderERROR)
	}
	return &response.CACertificates.CACertificate[0], nil
}

func (s *Service) WaitForSlbCACertificate(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeSlbCACertificate(id)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errs.WrapError(err)
			}
		} else {
			break
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, object.CACertificateId, id, errs.ProviderERROR)
		}
	}
	return nil
}

func (s *Service) DescribeSlbServerCertificate(id string) (*slb.ServerCertificate, error) {
	certificate := &slb.ServerCertificate{}
	request := slb.CreateDescribeServerCertificatesRequest()
	request.RegionId = s.client.RegionId
	request.ServerCertificateId = id

	raw, err := s.client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.DescribeServerCertificates(request)
	})
	if err != nil {
		return certificate, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*slb.DescribeServerCertificatesResponse)

	if len(response.ServerCertificates.ServerCertificate) < 1 || response.ServerCertificates.ServerCertificate[0].ServerCertificateId != id {
		return certificate, errs.WrapErrorf(errs.NotFoundErr("SlbServerCertificate", id), errs.NotFoundMsg, errs.ProviderERROR)
	}

	return &response.ServerCertificates.ServerCertificate[0], nil
}

func (s *Service) WaitForSlbServerCertificate(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeSlbServerCertificate(id)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errs.WrapError(err)
			}
		}
		if object.ServerCertificateId == id {
			break
		}

		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, object.ServerCertificateId, id, errs.ProviderERROR)
		}
	}
	return nil
}

func toSlbTagsString(tags []tags.Tag) string {
	slbTags := make([]SlbTag, 0, len(tags))

	for _, tag := range tags {
		slbTag := SlbTag{
			TagKey:   tag.Key,
			TagValue: tag.Value,
		}
		slbTags = append(slbTags, slbTag)
	}

	b, _ := json.Marshal(slbTags)

	return string(b)
}

func (s *Service) DescribeDomainExtensionAttribute(domainExtensionId string) (*slb.DescribeDomainExtensionAttributeResponse, error) {
	response := &slb.DescribeDomainExtensionAttributeResponse{}
	request := slb.CreateDescribeDomainExtensionAttributeRequest()
	request.DomainExtensionId = domainExtensionId
	var raw interface{}
	var err error
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		raw, err = s.client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.DescribeDomainExtensionAttribute(request)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{errs.AliyunGoClientFailure, "ServiceUnavailable", errs.Throttling}) {
				time.Sleep(10 * time.Second)
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidParameter.DomainExtensionId", "InvalidParameter"}) {
			return response, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return response, errs.WrapErrorf(err, errs.DefaultErrorMsg, domainExtensionId, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	response, _ = raw.(*slb.DescribeDomainExtensionAttributeResponse)
	if response.DomainExtensionId != domainExtensionId {
		return response, errs.WrapErrorf(errs.NotFoundErr("SLBDomainExtension", domainExtensionId), errs.NotFoundMsg, errs.ProviderERROR)
	}
	return response, nil
}

func (s *Service) WaitForSlbDomainExtension(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

	for {
		_, err := s.DescribeDomainExtensionAttribute(id)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errs.WrapError(err)
			}
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, Null, string(status), errs.ProviderERROR)
		}
	}
}

func (s *Service) setInstanceTags(d *schema.ResourceData, resourceType tags.TagResourceType) error {
	oraw, nraw := d.GetChange("tags")
	o := oraw.(map[string]interface{})
	n := nraw.(map[string]interface{})
	create, remove := s.diffTags(s.tagsFromMap(o), s.tagsFromMap(n))

	if len(remove) > 0 {
		var tagKey []string
		for _, v := range remove {
			tagKey = append(tagKey, v.Key)
		}
		request := slb.CreateUntagResourcesRequest()
		request.ResourceId = &[]string{d.Id()}
		request.ResourceType = string(resourceType)
		request.TagKey = &tagKey
		request.RegionId = s.client.RegionId

		wait := helper.IncrementalWait(1*time.Second, 1*time.Second)
		err := retry.Retry(10*time.Minute, func() *retry.RetryError {
			raw, err := s.client.WithSlbClient(func(client *slb.Client) (interface{}, error) {
				return client.UntagResources(request)
			})
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)

				}
				return retry.NonRetryableError(err)
			}
			logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
			return nil
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}
	}

	if len(create) > 0 {
		request := slb.CreateTagResourcesRequest()
		request.ResourceId = &[]string{d.Id()}
		request.Tag = &create
		request.ResourceType = string(resourceType)
		request.RegionId = s.client.RegionId

		wait := helper.IncrementalWait(1*time.Second, 1*time.Second)
		err := retry.Retry(10*time.Minute, func() *retry.RetryError {
			raw, err := s.client.WithSlbClient(func(client *slb.Client) (interface{}, error) {
				return client.TagResources(request)
			})
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)

				}
				return retry.NonRetryableError(err)
			}
			logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
			return nil
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}
	}

	return nil
}

func (s *Service) tagsToMap(tags []slb.TagResource) map[string]string {
	result := make(map[string]string)
	for _, t := range tags {
		if !s.ignoreTag(t) {
			result[t.TagKey] = t.TagValue
		}
	}
	return result
}

func (s *Service) ignoreTag(t slb.TagResource) bool {
	filter := []string{"^aliyun", "^acs:", "^http://", "^https://"}
	for _, v := range filter {
		log.Printf("[DEBUG] Matching prefix %v with %v\n", v, t.TagKey)
		ok, _ := regexp.MatchString(v, t.TagKey)
		if ok {
			log.Printf("[DEBUG] Found Alibaba Cloud specific t %s (val: %s), ignoring.\n", t.TagKey, t.TagValue)
			return true
		}
	}
	return false
}

func (s *Service) diffTags(oldTags, newTags []slb.TagResourcesTag) ([]slb.TagResourcesTag, []slb.TagResourcesTag) {
	// First, we're creating everything we have
	create := make(map[string]interface{})
	for _, t := range newTags {
		create[t.Key] = t.Value
	}

	// Build the list of what to remove
	var remove []slb.TagResourcesTag
	for _, t := range oldTags {
		old, ok := create[t.Key]
		if !ok || old != t.Value {
			// Delete it!
			remove = append(remove, t)
		}
	}

	return s.tagsFromMap(create), remove
}

func (s *Service) tagsFromMap(m map[string]interface{}) []slb.TagResourcesTag {
	result := make([]slb.TagResourcesTag, 0, len(m))
	for k, v := range m {
		result = append(result, slb.TagResourcesTag{
			Key:   k,
			Value: v.(string),
		})
	}

	return result
}

func (s *Service) ListTagResources(id string, resourceType string) (object interface{}, err error) {
	client := s.client
	action := "ListTagResources"
	request := map[string]interface{}{
		"RegionId":     s.client.RegionId,
		"ResourceType": resourceType,
		"ResourceId.1": id,
	}
	tags := make([]interface{}, 0)
	var response map[string]interface{}

	for {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err := client.RpcPost("Slb", "2014-05-15", action, nil, request, false)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			logs.AddDebug(action, response, request)
			v, err := jsonpath.Get("$.TagResources.TagResource", response)
			if err != nil {
				return retry.NonRetryableError(errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.TagResources.TagResource", response))
			}
			if v != nil {
				tags = append(tags, v.([]interface{})...)
			}
			return nil
		})
		if err != nil {
			err = errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
			return
		}
		if response["NextToken"] == nil {
			break
		}
		request["NextToken"] = response["NextToken"]
	}

	return tags, nil
}

func (s *Service) SetResourceTags(d *schema.ResourceData, resourceType string) error {

	if d.HasChange("tags") {
		added, removed := tags.ParsingTags(d)
		client := s.client

		removedTagKeys := make([]string, 0)
		for _, v := range removed {
			if !tags.IgnoredTags(v, "") {
				removedTagKeys = append(removedTagKeys, v)
			}
		}
		if len(removedTagKeys) > 0 {
			action := "UntagResources"
			request := map[string]interface{}{
				"RegionId":     s.client.RegionId,
				"ResourceType": resourceType,
				"ResourceId.1": d.Id(),
			}
			for i, key := range removedTagKeys {
				request[fmt.Sprintf("TagKey.%d", i+1)] = key
			}
			wait := helper.IncrementalWait(2*time.Second, 1*time.Second)
			err := retry.Retry(10*time.Minute, func() *retry.RetryError {
				response, err := client.RpcPost("Slb", "2014-05-15", action, nil, request, false)
				if err != nil {
					if helper.NeedRetry(err) {
						wait()
						return retry.RetryableError(err)

					}
					return retry.NonRetryableError(err)
				}
				logs.AddDebug(action, response, request)
				return nil
			})
			if err != nil {
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
			}
		}
		if len(added) > 0 {
			action := "TagResources"
			request := map[string]interface{}{
				"RegionId":     s.client.RegionId,
				"ResourceType": resourceType,
				"ResourceId.1": d.Id(),
			}
			count := 1
			for key, value := range added {
				request[fmt.Sprintf("Tag.%d.Key", count)] = key
				request[fmt.Sprintf("Tag.%d.Value", count)] = value
				count++
			}

			wait := helper.IncrementalWait(2*time.Second, 1*time.Second)
			err := retry.Retry(10*time.Minute, func() *retry.RetryError {
				response, err := client.RpcPost("Slb", "2014-05-15", action, nil, request, false)
				if err != nil {
					if helper.NeedRetry(err) {
						wait()
						return retry.RetryableError(err)

					}
					return retry.NonRetryableError(err)
				}
				logs.AddDebug(action, response, request)
				return nil
			})
			if err != nil {
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
			}
		}
	}
	return nil
}

func (s *Service) DescribeSlbLoadBalancer(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeLoadBalancerAttribute"
	request := map[string]interface{}{
		"RegionId":       s.client.RegionId,
		"LoadBalancerId": id,
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Slb", "2014-05-15", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidLoadBalancerId.NotFound"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("SLB:LoadBalancer", id), errs.NotFoundMsg, errs.ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *Service) SlbLoadBalancerStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeSlbLoadBalancer(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["LoadBalancerStatus"]) == failState {
				return object, fmt.Sprint(object["LoadBalancerStatus"]), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprint(object["LoadBalancerStatus"])))
			}
		}
		return object, fmt.Sprint(object["LoadBalancerStatus"]), nil
	}
}

func (s *Service) DescribeSlbCaCertificate(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeCACertificates"
	request := map[string]interface{}{
		"RegionId":        s.client.RegionId,
		"CACertificateId": id,
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Slb", "2014-05-15", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.CACertificates.CACertificate", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.CACertificates.CACertificate", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("SLB", id), errs.NotFoundWithResponse, response)
	} else {
		if v.([]interface{})[0].(map[string]interface{})["CACertificateId"].(string) != id {
			return object, errs.WrapErrorf(errs.NotFoundErr("SLB", id), errs.NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *Service) DescribeSlbTlsCipherPolicy(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "ListTLSCipherPolicies"
	request := map[string]interface{}{
		"RegionId":          s.client.RegionId,
		"TLSCipherPolicyId": id,
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Slb", "2014-05-15", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.TLSCipherPolicies", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.TLSCipherPolicies", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("SLB", id), errs.NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["InstanceId"]) != id {
			return object, errs.WrapErrorf(errs.NotFoundErr("SLB", id), errs.NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *Service) convertAclEntriesToString(v []interface{}) (string, error) {
	arrayMaps := make([]interface{}, len(v))
	for i, vv := range v {
		item := vv.(map[string]interface{})
		temp := map[string]interface{}{
			"comment": item["comment"],
			"entry":   item["entry"],
		}

		arrayMaps[i] = temp
	}
	maps, err := json.Marshal(arrayMaps)
	if err != nil {
		return "", errs.WrapError(err)
	}
	return string(maps), nil
}

func (s *Service) DescribeSlbAclEntryAttachment(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return object, errs.WrapError(err)
	}

	action := "DescribeAccessControlListAttribute"
	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
		"AclId":    parts[0],
	}
	idExist := false
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Slb", "2014-05-15", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"AclNotExist"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("SLB", id), errs.NotFoundWithResponse, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.AclEntrys.AclEntry", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.AclEntrys.AclEntry", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("SLB", id), errs.NotFoundWithResponse, response)
	}
	for _, v := range v.([]interface{}) {
		if fmt.Sprint(v.(map[string]interface{})["AclEntryIP"]) == parts[1] {
			idExist = true
			return v.(map[string]interface{}), nil
		}
	}
	if !idExist {
		return object, errs.WrapErrorf(errs.NotFoundErr("SLB", id), errs.NotFoundWithResponse, response)
	}
	return object, nil
}

// ListSlbAclEntries lists all entries of an acl through DescribeAccessControlListAttribute.
// It is the batch-mode read path for alicloud_slb_acl_entry_attachment.
func (s *Service) ListSlbAclEntries(aclId string) (entries []map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeAccessControlListAttribute"
	request := map[string]interface{}{
		"RegionId": client.RegionId,
		"AclId":    aclId,
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Slb", "2014-05-15", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"AclNotExist"}) {
			return entries, errs.WrapErrorf(errs.NotFoundErr("SLB:Acl", aclId), errs.NotFoundWithResponse, response)
		}
		return entries, errs.WrapErrorf(err, errs.DefaultErrorMsg, aclId, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.AclEntrys.AclEntry", response)
	if err != nil {
		return entries, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, aclId, "$.AclEntrys.AclEntry", response)
	}
	list, ok := v.([]interface{})
	if !ok || len(list) == 0 {
		return entries, nil
	}
	entries = make([]map[string]interface{}, 0, len(list))
	for _, item := range list {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		entries = append(entries, map[string]interface{}{
			"entry":   fmt.Sprint(m["AclEntryIP"]),
			"comment": slbAclEntryComment(m["AclEntryComment"]),
		})
	}
	return entries, nil
}

// slbAclEntryComment coerces the AclEntryComment field of a DescribeAccessControlListAttribute
// entry into a string, returning "" when the field is absent or not a string.
func slbAclEntryComment(v interface{}) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func (s *Service) DescribeSlbServerGroupServerAttachment(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	parts, err := helper.ParseResourceId(id, 3)
	if err != nil {
		return object, errs.WrapError(err)
	}

	action := "DescribeVServerGroupAttribute"
	request := map[string]interface{}{
		"RegionId":       s.client.RegionId,
		"VServerGroupId": parts[0],
	}
	idExist := false
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Slb", "2014-05-15", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"The specified VServerGroupId does not exist", "InvalidParameter"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("SLB", id), errs.NotFoundWithResponse, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	resp, err := jsonpath.Get("$.BackendServers.BackendServer", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.BackendServers.BackendServer", response)
	}

	if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("Slb:ServerGroupServerAttachment", id), errs.NotFoundWithResponse, response)
	}

	for _, v := range resp.([]interface{}) {
		if fmt.Sprint(v.(map[string]interface{})["ServerId"]) == parts[1] && fmt.Sprint(v.(map[string]interface{})["Port"]) == parts[2] {
			idExist = true
			return v.(map[string]interface{}), nil
		}
	}

	if !idExist {
		return object, errs.WrapErrorf(errs.NotFoundErr("Slb:ServerGroupServerAttachment", id), errs.NotFoundWithResponse, response)
	}

	return object, nil
}
