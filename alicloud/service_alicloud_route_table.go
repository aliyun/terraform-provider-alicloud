package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/schema"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
)

func (client *AliyunClient) DescribeRouteTable(routeTableId string) (v vpc.RouterTableListType, err error) {
	request := vpc.CreateDescribeRouteTableListRequest()

	invoker := NewInvoker()
	err = invoker.Run(func() error {
		resp, err := client.vpcconn.DescribeRouteTableList(request)
		if err != nil {
			return err
		}
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

func (client *AliyunClient) DescribeRouteTableAttachment(routeTableId string, vSwitchId string) (err error) {
	invoker := NewInvoker()
	return invoker.Run(func() error {
		routeTable, err := client.DescribeRouteTable(routeTableId)
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

func (client *AliyunClient) WaitForRouteTable(routeTableId string, timeout int) error {

	if timeout <= 0 {
		timeout = DefaultTimeout
	}
	for {
		time.Sleep(3 * time.Second)
		resp, err := client.DescribeRouteTable(routeTableId)

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
	return nil
}

func (client *AliyunClient) WaitForRouteTableAttachment(routeTableId string, vSwitchId string, timeout int) error {

	if timeout <= 0 {
		timeout = DefaultTimeout
	}
	for {
		err := client.DescribeRouteTableAttachment(routeTableId, vSwitchId)
		if err != nil {
			return err
		} else {
			break
		}
		timeout = timeout - DefaultIntervalShort
		if timeout <= 0 {
			return GetTimeErrorFromString(GetTimeoutMessage("Route Table Attachment", string("Unavailable")))
		}
		time.Sleep(DefaultTimeout * time.Second)
	}
	return nil
}

func getRouteTableIdAndVSwitchId(d *schema.ResourceData, meta interface{}) (string, string, error) {
	parts := strings.Split(d.Id(), COLON_SEPARATED)

	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid resource id")
	}
	return parts[0], parts[1], nil
}
