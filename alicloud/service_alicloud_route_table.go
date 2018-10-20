package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/schema"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

type RouteTableService struct {
	client *connectivity.AliyunClient
}

func (s *RouteTableService) DescribeRouteTable(routeTableId string) (v vpc.RouterTableListType, err error) {
	request := vpc.CreateDescribeRouteTableListRequest()
	request.RouteTableId = routeTableId

	invoker := NewInvoker()
	err = invoker.Run(func() error {
		raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeRouteTableList(request)
		})
		if err != nil {
			return err
		}
		resp, _ := raw.(*vpc.DescribeRouteTableListResponse)
		length := len(resp.RouterTableList.RouterTableListType)
		if resp == nil || length <= 0 {
			return GetNotFoundErrorFromString(GetNotFoundMessage("RouteTable", routeTableId))
		}
		//Finding the routeTableId
		for _, id := range resp.RouterTableList.RouterTableListType {
			if id.RouteTableId == routeTableId {
				v = id
				return nil
			}
		}
		return GetNotFoundErrorFromString(GetNotFoundMessage("RouteTableId", routeTableId))
	})
	return
}

func (s *RouteTableService) DescribeRouteTableAttachment(routeTableId string, vSwitchId string) (err error) {
	invoker := NewInvoker()
	return invoker.Run(func() error {
		routeTable, err := s.DescribeRouteTable(routeTableId)
		if err != nil {
			return err
		}
		for _, id := range routeTable.VSwitchIds.VSwitchId {
			if id == vSwitchId {
				return nil
			}
		}
		return GetNotFoundErrorFromString(GetNotFoundMessage("RouteTableAttachment", routeTableId+COLON_SEPARATED+vSwitchId))
	})
}

func (s *RouteTableService) WaitForRouteTable(routeTableId string, timeout int) error {

	if timeout <= 0 {
		timeout = DefaultTimeout
	}
	for {
		time.Sleep(DefaultIntervalShort * time.Second)
		resp, err := s.DescribeRouteTable(routeTableId)

		if err != nil {
			return err
		}
		if resp.RouteTableId == routeTableId {
			return nil
		}
		timeout = timeout - DefaultIntervalShort
		if timeout <= 0 {
			return GetTimeErrorFromString(GetTimeoutMessage("Route Table Attachment", string("Unavailable")))
		}
	}
}

func (s *RouteTableService) WaitForRouteTableAttachment(routeTableId string, vSwitchId string, timeout int) error {

	if timeout <= 0 {
		timeout = DefaultTimeout
	}
	for {
		err := s.DescribeRouteTableAttachment(routeTableId, vSwitchId)
		if err != nil {
			if !NotFoundError(err) {
				return err
			}
		} else {
			break
		}
		timeout = timeout - DefaultIntervalShort
		if timeout <= 0 {
			return GetTimeErrorFromString(GetTimeoutMessage("Route Table Attachment", string("Unavailable")))
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *RouteTableService) GetRouteTableIdAndVSwitchId(d *schema.ResourceData, meta interface{}) (string, string, error) {
	parts := strings.Split(d.Id(), COLON_SEPARATED)

	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid resource id")
	}
	return parts[0], parts[1], nil
}
