package alicloud

import (
	"reflect"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

type AlidnsService struct {
	client *connectivity.AliyunClient
}

func (s *AlidnsService) DescribeInstanceDomains(id string) (object alidns.DescribeInstanceDomainsResponse, err error) {
	request := alidns.CreateDescribeInstanceDomainsRequest()
	request.RegionId = s.client.RegionId

	request.InstanceId = id

	raw, err := s.client.WithAlidnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
		return alidnsClient.DescribeInstanceDomains(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDnsProduct"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("alidnsinstance", id)), NotFoundMsg, ProviderERROR)
			return
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*alidns.DescribeInstanceDomainsResponse)

	if len(response.InstanceDomains) < 1 {
		err = WrapErrorf(Error(GetNotFoundMessage("alidnsinstance", id)), NotFoundMsg, ProviderERROR)
		return
	}
	return *response, nil
}

func (s *AlidnsService) setAlidnsInstanceDomains(d *schema.ResourceData) error {
	o, n := d.GetChange("domain_name")
	oldNames := strings.Split(o.(string), ",")
	newNames := strings.Split(n.(string), ",")
	oldmap := make(map[string]string)
	newmap := make(map[string]string)
	add := make([]string, 0)
	remove := make([]string, 0)
	for _, v := range oldNames {
		oldmap[v] = v
	}
	for _, v := range newNames {
		if _, ok := oldmap[v]; !ok {
			add = append(add, v)
		}
	}

	for _, v := range newNames {
		newmap[v] = v
	}
	for _, v := range oldNames {
		if _, ok := newmap[v]; !ok {
			remove = append(remove, v)
		}
	}
	if len(remove) > 0 {
		removeNames := strings.Join(remove, ",")
		request := alidns.CreateUnbindInstanceDomainsRequest()
		request.InstanceId = d.Id()
		if lang, ok := d.GetOk("lang"); ok {
			request.Lang = lang.(string)
		}
		request.DomainNames = removeNames
		raw, err := s.client.WithAlidnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
			return alidnsClient.UnbindInstanceDomains(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}
	if len(add) > 0 {
		addNames := strings.Join(add, ",")
		request := alidns.CreateBindInstanceDomainsRequest()
		request.InstanceId = d.Id()
		request.DomainNames = addNames
		raw, err := s.client.WithAlidnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
			return alidnsClient.BindInstanceDomains(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}
	return nil
}

func (s *AlidnsService) DescribeAlidnsInstance(id string) (object alidns.DescribeDnsProductInstanceResponse, err error) {
	request := alidns.CreateDescribeDnsProductInstanceRequest()
	request.RegionId = s.client.RegionId

	request.InstanceId = id

	raw, err := s.client.WithAlidnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
		return alidnsClient.DescribeDnsProductInstance(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDnsProduct"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("AlidnsInstance", id)), NotFoundMsg, ProviderERROR)
			return
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*alidns.DescribeDnsProductInstanceResponse)
	return *response, nil
}

func (s *AlidnsService) WaitForAlidnsInstance(id string, expected map[string]interface{}, isDelete bool, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeInstanceDomains(id)
		if err != nil {
			if NotFoundError(err) {
				if isDelete {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		domainNames := make(map[string]interface{}, 0)
		for _, v := range object.InstanceDomains {
			domainNames[v.DomainName] = v.DomainName
		}

		exceptDomainNames := make(map[string]interface{}, 0)
		for _, v := range expected {
			names := strings.Split(v.(string), ",")
			for _, vv := range names {
				exceptDomainNames[vv] = vv
			}
		}

		if reflect.DeepEqual(domainNames, exceptDomainNames) {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, "", expected, ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}
