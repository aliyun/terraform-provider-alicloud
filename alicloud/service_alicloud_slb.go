package alicloud

import (
	"fmt"

	"encoding/json"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
)

func (client *AliyunClient) BuildSlbCommonRequest() *requests.CommonRequest {
	request := requests.NewCommonRequest()
	endpoint := LoadEndpoint(client.RegionId, SLBCode)
	if endpoint == "" {
		endpoint, _ = client.DescribeEndpointByCode(client.RegionId, SLBCode)
	}
	if endpoint == "" {
		endpoint = fmt.Sprintf("slb.%s.aliyuncs.com", client.RegionId)
	}
	request.Domain = endpoint
	request.Version = ApiVersion20140515
	request.RegionId = client.RegionId
	return request
}
func (client *AliyunClient) DescribeLoadBalancerAttribute(slbId string) (loadBalancer *slb.DescribeLoadBalancerAttributeResponse, err error) {

	req := slb.CreateDescribeLoadBalancerAttributeRequest()
	req.LoadBalancerId = slbId
	loadBalancer, err = client.slbconn.DescribeLoadBalancerAttribute(req)

	if err != nil {
		if IsExceptedErrors(err, []string{LoadBalancerNotFound}) {
			err = GetNotFoundErrorFromString(GetNotFoundMessage("LoadBalancer", slbId))
		}
		return
	}
	if loadBalancer == nil || loadBalancer.LoadBalancerId == "" {
		err = GetNotFoundErrorFromString(GetNotFoundMessage("LoadBalancer", slbId))
	}
	return
}

func (client *AliyunClient) DescribeLoadBalancerRuleId(slbId string, port int, domain, url string) (string, error) {
	req := slb.CreateDescribeRulesRequest()
	req.LoadBalancerId = slbId
	req.ListenerPort = requests.NewInteger(port)
	if rules, err := client.slbconn.DescribeRules(req); err != nil {
		return "", fmt.Errorf("DescribeRules got an error: %#v", err)
	} else {
		for _, rule := range rules.Rules.Rule {
			if rule.Domain == domain && rule.Url == url {
				return rule.RuleId, nil
			}
		}
	}
	return "", GetNotFoundErrorFromString(fmt.Sprintf("Rule is not found based on domain %s and url %s.", domain, url))
}

func (client *AliyunClient) DescribeLoadBalancerRuleAttribute(ruleId string) (*slb.DescribeRuleAttributeResponse, error) {
	req := slb.CreateDescribeRuleAttributeRequest()
	req.RuleId = ruleId
	rule, err := client.slbconn.DescribeRuleAttribute(req)
	if err != nil {
		if IsExceptedErrors(err, []string{InvalidRuleIdNotFound}) {
			return nil, GetNotFoundErrorFromString(GetNotFoundMessage("SLB Rule", ruleId))
		}
		return nil, fmt.Errorf("DescribeLoadBalancerRuleAttribute got an error: %#v", err)
	}
	if rule == nil || rule.LoadBalancerId == "" {
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("SLB Rule", ruleId))
	}
	return rule, err
}

func (client *AliyunClient) DescribeSlbVServerGroupAttribute(groupId string) (*slb.DescribeVServerGroupAttributeResponse, error) {
	req := slb.CreateDescribeVServerGroupAttributeRequest()
	req.VServerGroupId = groupId
	group, err := client.slbconn.DescribeVServerGroupAttribute(req)
	if err != nil {
		if IsExceptedErrors(err, []string{VServerGroupNotFoundMessage, InvalidParameter}) {
			return nil, GetNotFoundErrorFromString(GetNotFoundMessage("SLB VServer Group", groupId))
		}
		return nil, fmt.Errorf("DescribeSlbVServerGroupAttribute got an error: %#v", err)
	}
	if group == nil || group.VServerGroupId == "" {
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("SLB VServer Group", groupId))
	}
	return group, err
}

func (client *AliyunClient) DescribeLoadBalancerListenerAttribute(loadBalancerId string, port int, protocol Protocol) (listener map[string]interface{}, err error) {
	req := client.BuildSlbCommonRequest()
	req.ApiName = fmt.Sprintf("DescribeLoadBalancer%sListenerAttribute", strings.ToUpper(string(protocol)))
	req.QueryParams["LoadBalancerId"] = loadBalancerId
	req.QueryParams["ListenerPort"] = string(requests.NewInteger(port))
	resp, err := client.slbconn.ProcessCommonRequest(req)
	if err != nil {
		return
	}

	if err = json.Unmarshal(resp.GetHttpContentBytes(), &listener); err != nil {
		err = fmt.Errorf("Unmarshalling body got an error: %#v.", err)
	}

	return

}

func (client *AliyunClient) WaitForLoadBalancer(loadBalancerId string, status Status, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}

	for {
		lb, err := client.DescribeLoadBalancerAttribute(loadBalancerId)

		if err != nil {
			if !NotFoundError(err) {

				return err
			}
		} else if &lb != nil && strings.ToLower(lb.LoadBalancerStatus) == strings.ToLower(string(status)) {
			//TODO
			break
		}
		timeout = timeout - DefaultIntervalShort
		if timeout <= 0 {
			return GetTimeErrorFromString(GetTimeoutMessage("LoadBalancer", string(status)))
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}

func (client *AliyunClient) WaitForListener(loadBalancerId string, port int, protocol Protocol, status Status, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}

	for {
		listener, err := client.DescribeLoadBalancerListenerAttribute(loadBalancerId, port, protocol)
		if err != nil && !IsExceptedErrors(err, []string{LoadBalancerNotFound}) {
			return err
		}

		if value, ok := listener["Status"]; ok && strings.ToLower(value.(string)) == strings.ToLower(string(status)) {
			//TODO
			break
		}
		timeout = timeout - DefaultIntervalShort
		if timeout <= 0 {
			return GetTimeErrorFromString(GetTimeoutMessage("LoadBalancer Listener", string(status)))
		}
		time.Sleep(DefaultIntervalShort * time.Second)

	}
	return nil
}
