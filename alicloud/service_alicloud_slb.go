package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type SlbService struct {
	client *connectivity.AliyunClient
}

type SlbTag struct {
	TagKey   string
	TagValue string
}

const max_num_per_time = 50
const tags_max_num_per_time = 5
const tags_max_page_size = 50

func (s *SlbService) BuildSlbCommonRequest() (*requests.CommonRequest, error) {
	// Get product code from the built request
	slbReq := slb.CreateCreateLoadBalancerRequest()
	req, err := s.client.NewCommonRequest(slbReq.GetProduct(), slbReq.GetLocationServiceCode(), strings.ToUpper(string(Https)), connectivity.ApiVersion20140515)
	if err != nil {
		return req, WrapError(err)
	}
	req.RegionId = s.client.RegionId
	return req, err
}

func (s *SlbService) DescribeSlb(id string) (*slb.DescribeLoadBalancerAttributeResponse, error) {
	response := &slb.DescribeLoadBalancerAttributeResponse{}
	request := slb.CreateDescribeLoadBalancerAttributeRequest()
	request.RegionId = s.client.RegionId
	request.LoadBalancerId = id
	raw, err := s.client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.DescribeLoadBalancerAttribute(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidLoadBalancerId.NotFound"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("Slb", id)), NotFoundMsg, AlibabaCloudSdkGoERROR)
		} else {
			err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		return response, err
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest)
	response, _ = raw.(*slb.DescribeLoadBalancerAttributeResponse)
	if response.LoadBalancerId == "" {
		err = WrapErrorf(Error(GetNotFoundMessage("Slb", id)), NotFoundMsg, ProviderERROR)
	}
	return response, err
}

func (s *SlbService) DescribeSlbRule(id string) (*slb.DescribeRuleAttributeResponse, error) {
	response := &slb.DescribeRuleAttributeResponse{}
	request := slb.CreateDescribeRuleAttributeRequest()
	request.RegionId = s.client.RegionId
	request.RuleId = id
	raw, err := s.client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.DescribeRuleAttribute(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidRuleId.NotFound"}) {
			return response, WrapErrorf(Error(GetNotFoundMessage("SlbRule", id)), NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return response, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ = raw.(*slb.DescribeRuleAttributeResponse)
	if response.RuleId != id {
		return response, WrapErrorf(Error(GetNotFoundMessage("SlbRule", id)), NotFoundMsg, AlibabaCloudSdkGoERROR)
	}
	return response, nil
}

func (s *SlbService) DescribeSlbServerGroup(id string) (*slb.DescribeVServerGroupAttributeResponse, error) {
	response := &slb.DescribeVServerGroupAttributeResponse{}
	request := slb.CreateDescribeVServerGroupAttributeRequest()
	request.RegionId = s.client.RegionId
	request.VServerGroupId = id
	raw, err := s.client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.DescribeVServerGroupAttribute(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"The specified VServerGroupId does not exist", "InvalidParameter"}) {
			return response, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return response, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ = raw.(*slb.DescribeVServerGroupAttributeResponse)
	if response.VServerGroupId == "" {
		return response, WrapErrorf(Error(GetNotFoundMessage("SlbServerGroup", id)), NotFoundMsg, ProviderERROR)
	}
	return response, err
}

func (s *SlbService) DescribeSlbMasterSlaveServerGroup(id string) (*slb.DescribeMasterSlaveServerGroupAttributeResponse, error) {
	response := &slb.DescribeMasterSlaveServerGroupAttributeResponse{}
	request := slb.CreateDescribeMasterSlaveServerGroupAttributeRequest()
	request.RegionId = s.client.RegionId
	request.MasterSlaveServerGroupId = id
	raw, err := s.client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.DescribeMasterSlaveServerGroupAttribute(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"The specified MasterSlaveGroupId does not exist", "InvalidParameter"}) {
			return response, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return response, WrapErrorf(err, DefaultDebugMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ = raw.(*slb.DescribeMasterSlaveServerGroupAttributeResponse)
	if response.MasterSlaveServerGroupId == "" {
		return response, WrapErrorf(Error(GetNotFoundMessage("SlbMasterSlaveServerGroup", id)), NotFoundMsg, ProviderERROR)
	}
	return response, err
}

func (s *SlbService) DescribeSlbBackendServer(id string) (*slb.DescribeLoadBalancerAttributeResponse, error) {
	response := &slb.DescribeLoadBalancerAttributeResponse{}
	request := slb.CreateDescribeLoadBalancerAttributeRequest()
	request.RegionId = s.client.RegionId
	request.LoadBalancerId = id
	raw, err := s.client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.DescribeLoadBalancerAttribute(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidLoadBalancerId.NotFound"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("SlbBackendServers", id)), NotFoundMsg, AlibabaCloudSdkGoERROR)
		} else {
			err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		return response, err
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ = raw.(*slb.DescribeLoadBalancerAttributeResponse)
	if response.LoadBalancerId == "" {
		err = WrapErrorf(Error(GetNotFoundMessage("SlbBackendServers", id)), NotFoundMsg, ProviderERROR)
	}
	return response, err
}

func (s *SlbService) DescribeSlbListener(id string) (listener map[string]interface{}, err error) {
	parts, err := ParseSlbListenerId(id)
	if err != nil {
		return nil, WrapError(err)
	}
	protocol := parts[1]
	request, err := s.BuildSlbCommonRequest()
	request.RegionId = s.client.RegionId
	if err != nil {
		err = WrapError(err)
		return
	}
	request.ApiName = fmt.Sprintf("DescribeLoadBalancer%sListenerAttribute", strings.ToUpper(string(protocol)))
	request.QueryParams["LoadBalancerId"] = parts[0]
	port, _ := strconv.Atoi(parts[2])
	request.QueryParams["ListenerPort"] = string(requests.NewInteger(port))
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := s.client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.ProcessCommonRequest(request)
		})

		if err != nil {
			if IsExpectedErrors(err, []string{"The specified resource does not exist"}) {
				return resource.NonRetryableError(WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR))
			} else if IsExpectedErrors(err, SlbIsBusy) {
				return resource.RetryableError(WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR))
			}
			return resource.NonRetryableError(WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR))
		}
		addDebug(request.GetActionName(), raw, request, request.QueryParams)
		response, _ := raw.(*responses.CommonResponse)
		if err = json.Unmarshal(response.GetHttpContentBytes(), &listener); err != nil {
			return resource.NonRetryableError(WrapError(err))
		}
		if port, ok := listener["ListenerPort"]; ok && port.(float64) > 0 {
			return nil
		} else {
			return resource.RetryableError(WrapErrorf(Error(GetNotFoundMessage("SlbListener", id)), NotFoundMsg, ProviderERROR))
		}
	})

	return
}

func (s *SlbService) DescribeSlbAcl(id string) (*slb.DescribeAccessControlListAttributeResponse, error) {
	response := &slb.DescribeAccessControlListAttributeResponse{}
	request := slb.CreateDescribeAccessControlListAttributeRequest()
	request.RegionId = s.client.RegionId
	request.AclId = id

	raw, err := s.client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.DescribeAccessControlListAttribute(request)
	})
	if err != nil {
		if err != nil {
			if IsExpectedErrors(err, []string{"AclNotExist"}) {
				return response, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
			}
			return response, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ = raw.(*slb.DescribeAccessControlListAttributeResponse)
	return response, nil
}

func (s *SlbService) WaitForSlbAcl(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeSlbAcl(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		} else {
			return nil
		}

		time.Sleep(DefaultIntervalShort * time.Second)
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.AclId, id, ProviderERROR)
		}
	}
}

func (s *SlbService) WaitForSlb(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeSlb(id)

		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		} else if strings.ToLower(object.LoadBalancerStatus) == strings.ToLower(string(status)) {
			//TODO
			break
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.LoadBalancerStatus, status, ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *SlbService) WaitForSlbListener(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeSlbListener(id)
		if err != nil && !IsExpectedErrors(err, []string{"InvalidLoadBalancerId.NotFound"}) {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
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
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, gotStatus, status, ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *SlbService) WaitForSlbRule(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeSlbRule(id)

		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if object.RuleId == id && status != Deleted {
			break
		}

		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, "", id, ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *SlbService) WaitForSlbServerGroup(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeSlbServerGroup(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if object.VServerGroupId == id {
			break
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.VServerGroupId, id, ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *SlbService) WaitForSlbMasterSlaveServerGroup(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeSlbMasterSlaveServerGroup(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if object.MasterSlaveServerGroupId == id && status != Deleted {
			break
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.MasterSlaveServerGroupId, id, ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *SlbService) WaitSlbAttribute(id string, instanceSet *schema.Set, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

RETRY:
	object, err := s.DescribeSlb(id)
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapError(err)
	}
	if time.Now().After(deadline) {
		return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, Null, id, ProviderERROR)
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

func (s *SlbService) slbRemoveAccessControlListEntryPerTime(list []interface{}, id string) error {
	request := slb.CreateRemoveAccessControlListEntryRequest()
	request.RegionId = s.client.RegionId
	request.AclId = id
	b, err := json.Marshal(list)
	if err != nil {
		return WrapError(err)
	}
	request.AclEntrys = string(b)
	raw, err := s.client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.RemoveAccessControlListEntry(request)
	})
	if err != nil {
		if !IsExpectedErrors(err, []string{"AclEntryEmpty"}) {
			return WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return nil
}

func (s *SlbService) SlbRemoveAccessControlListEntry(list []interface{}, aclId string) error {
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

func (s *SlbService) slbAddAccessControlListEntryPerTime(list []interface{}, id string) error {
	request := slb.CreateAddAccessControlListEntryRequest()
	request.RegionId = s.client.RegionId
	request.AclId = id
	b, err := json.Marshal(list)
	if err != nil {
		return WrapError(err)
	}
	request.AclEntrys = string(b)
	raw, err := s.client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.AddAccessControlListEntry(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return nil
}

func (s *SlbService) SlbAddAccessControlListEntry(list []interface{}, aclId string) error {
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
func (s *SlbService) FlattenSlbAclEntryMappings(list []slb.AclEntry) []map[string]interface{} {
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
func (s *SlbService) flattenSlbRelatedListenerMappings(list []slb.RelatedListener) []map[string]interface{} {
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

func (s *SlbService) DescribeSlbCACertificate(id string) (*slb.CACertificate, error) {
	certificate := &slb.CACertificate{}
	request := slb.CreateDescribeCACertificatesRequest()
	request.RegionId = s.client.RegionId
	request.CACertificateId = id
	raw, err := s.client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.DescribeCACertificates(request)
	})
	if err != nil {
		return certificate, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*slb.DescribeCACertificatesResponse)
	if len(response.CACertificates.CACertificate) < 1 {
		return certificate, WrapErrorf(Error(GetNotFoundMessage("SlbCACertificate", id)), NotFoundMsg, ProviderERROR)
	}
	return &response.CACertificates.CACertificate[0], nil
}

func (s *SlbService) WaitForSlbCACertificate(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeSlbCACertificate(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		} else {
			break
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.CACertificateId, id, ProviderERROR)
		}
	}
	return nil
}

func (s *SlbService) DescribeSlbServerCertificate(id string) (*slb.ServerCertificate, error) {
	certificate := &slb.ServerCertificate{}
	request := slb.CreateDescribeServerCertificatesRequest()
	request.RegionId = s.client.RegionId
	request.ServerCertificateId = id

	raw, err := s.client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.DescribeServerCertificates(request)
	})
	if err != nil {
		return certificate, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*slb.DescribeServerCertificatesResponse)

	if len(response.ServerCertificates.ServerCertificate) < 1 || response.ServerCertificates.ServerCertificate[0].ServerCertificateId != id {
		return certificate, WrapErrorf(Error(GetNotFoundMessage("SlbServerCertificate", id)), NotFoundMsg, ProviderERROR)
	}

	return &response.ServerCertificates.ServerCertificate[0], nil
}

func (s *SlbService) WaitForSlbServerCertificate(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeSlbServerCertificate(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if object.ServerCertificateId == id {
			break
		}

		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.ServerCertificateId, id, ProviderERROR)
		}
	}
	return nil
}

func toSlbTagsString(tags []Tag) string {
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

func (s *SlbService) DescribeDomainExtensionAttribute(domainExtensionId string) (*slb.DescribeDomainExtensionAttributeResponse, error) {
	response := &slb.DescribeDomainExtensionAttributeResponse{}
	request := slb.CreateDescribeDomainExtensionAttributeRequest()
	request.DomainExtensionId = domainExtensionId
	var raw interface{}
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err = s.client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.DescribeDomainExtensionAttribute(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{AliyunGoClientFailure, "ServiceUnavailable", Throttling}) {
				time.Sleep(10 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidParameter.DomainExtensionId", "InvalidParameter"}) {
			return response, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return response, WrapErrorf(err, DefaultErrorMsg, domainExtensionId, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	response, _ = raw.(*slb.DescribeDomainExtensionAttributeResponse)
	if response.DomainExtensionId != domainExtensionId {
		return response, WrapErrorf(Error(GetNotFoundMessage("SLBDomainExtension", domainExtensionId)), NotFoundMsg, ProviderERROR)
	}
	return response, nil
}

func (s *SlbService) WaitForSlbDomainExtension(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

	for {
		_, err := s.DescribeDomainExtensionAttribute(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, Null, string(status), ProviderERROR)
		}
	}
	return nil
}

func (s *SlbService) setInstanceTags(d *schema.ResourceData, resourceType TagResourceType) error {
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

		wait := incrementalWait(1*time.Second, 1*time.Second)
		err := resource.Retry(10*time.Minute, func() *resource.RetryError {
			raw, err := s.client.WithSlbClient(func(client *slb.Client) (interface{}, error) {
				return client.UntagResources(request)
			})
			if err != nil {
				if IsThrottling(err) {
					wait()
					return resource.RetryableError(err)

				}
				return resource.NonRetryableError(err)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}

	if len(create) > 0 {
		request := slb.CreateTagResourcesRequest()
		request.ResourceId = &[]string{d.Id()}
		request.Tag = &create
		request.ResourceType = string(resourceType)
		request.RegionId = s.client.RegionId

		wait := incrementalWait(1*time.Second, 1*time.Second)
		err := resource.Retry(10*time.Minute, func() *resource.RetryError {
			raw, err := s.client.WithSlbClient(func(client *slb.Client) (interface{}, error) {
				return client.TagResources(request)
			})
			if err != nil {
				if IsThrottling(err) {
					wait()
					return resource.RetryableError(err)

				}
				return resource.NonRetryableError(err)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}

	d.SetPartial("tags")

	return nil
}

func (s *SlbService) tagsToMap(tags []slb.TagResource) map[string]string {
	result := make(map[string]string)
	for _, t := range tags {
		if !s.ignoreTag(t) {
			result[t.TagKey] = t.TagValue
		}
	}
	return result
}

func (s *SlbService) ignoreTag(t slb.TagResource) bool {
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

func (s *SlbService) diffTags(oldTags, newTags []slb.TagResourcesTag) ([]slb.TagResourcesTag, []slb.TagResourcesTag) {
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

func (s *SlbService) tagsFromMap(m map[string]interface{}) []slb.TagResourcesTag {
	result := make([]slb.TagResourcesTag, 0, len(m))
	for k, v := range m {
		result = append(result, slb.TagResourcesTag{
			Key:   k,
			Value: v.(string),
		})
	}

	return result
}

func (s *SlbService) DescribeTags(resourceId string, resourceTags map[string]interface{}, resourceType TagResourceType) (tags []slb.TagResource, err error) {
	request := slb.CreateListTagResourcesRequest()
	request.RegionId = s.client.RegionId
	request.ResourceType = string(resourceType)
	request.ResourceId = &[]string{resourceId}
	if resourceTags != nil && len(resourceTags) > 0 {
		var reqTags []slb.ListTagResourcesTag
		for key, value := range resourceTags {
			reqTags = append(reqTags, slb.ListTagResourcesTag{
				Key:   key,
				Value: value.(string),
			})
		}
		request.Tag = &reqTags
	}
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := s.client.WithSlbClient(func(Client *slb.Client) (interface{}, error) {
			return Client.ListTagResources(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{Throttling}) {
				time.Sleep(2 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*slb.ListTagResourcesResponse)
		tags = response.TagResources.TagResource
		return nil
	})
	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, resourceId, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}

	return
}
