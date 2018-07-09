package alicloud

import (
	"fmt"

	"github.com/denverdino/aliyungo/slb"
)

func (client *AliyunClient) DescribeLoadBalancerAttribute(slbId string) (*slb.LoadBalancerType, error) {

	loadBalancer, err := client.slbconn.NewDescribeLoadBalancerAttribute(&slb.NewDescribeLoadBalancerAttributeArgs{
		RegionId:       client.Region,
		LoadBalancerId: slbId,
	})

	if err != nil {
		if IsExceptedErrors(err, []string{LoadBalancerNotFound}) {
			return nil, GetNotFoundErrorFromString(GetNotFoundMessage("Load Balancer", slbId))
		}
		return nil, err
	}
	if loadBalancer == nil || loadBalancer.LoadBalancerId == "" {
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("Load Balancer", slbId))
	}
	return loadBalancer, nil
}

func (client *AliyunClient) DescribeLoadBalancerRuleId(slbId string, port int, domain, url string) (string, error) {

	if rules, err := client.slbconn.DescribeRules(&slb.DescribeRulesArgs{
		RegionId:       client.Region,
		LoadBalancerId: slbId,
		ListenerPort:   port,
	}); err != nil {
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

	rule, err := client.slbconn.DescribeRuleAttribute(&slb.DescribeRuleAttributeArgs{
		RegionId: client.Region,
		RuleId:   ruleId,
	})
	if err != nil {
		if IsExceptedErrors(err, []string{InvalidRuleIdNotFound}) {
			return nil, GetNotFoundErrorFromString(GetNotFoundMessage("SLB Rule", ruleId))
		}
		return nil, fmt.Errorf("DescribeLoadBalancerRuleAttribute got an error: %#v", err)
	}
	if rule == nil || &rule.Rule == nil {
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("SLB Rule", ruleId))
	}
	return rule, err
}

func (client *AliyunClient) DescribeSlbVServerGroupAttribute(groupId string) (*slb.DescribeVServerGroupAttributeResponse, error) {
	group, err := client.slbconn.DescribeVServerGroupAttribute(&slb.DescribeVServerGroupAttributeArgs{
		RegionId:       client.Region,
		VServerGroupId: groupId,
	})
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
