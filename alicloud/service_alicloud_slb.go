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
		return nil, err
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
