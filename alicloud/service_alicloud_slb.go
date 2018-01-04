package alicloud

import (
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
