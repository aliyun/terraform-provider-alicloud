package alicloud

import (
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
		if resp == nil || len(resp.RouterTableList.RouterTableListType) <= 0 {
			return GetNotFoundErrorFromString(GetNotFoundMessage("RouteTableId", routeTableId))
		}
		//Finding the routeTableId
		flag := -1
		length := len(resp.RouterTableList.RouterTableListType)
		for i := 0; i < length; i++ {
			if resp.RouterTableList.RouterTableListType[i].RouteTableId == routeTableId {
				flag = i
			}
		}
		if flag == -1 {
			return GetNotFoundErrorFromString(GetNotFoundMessage("RouteTableId", routeTableId))
		}
		v = resp.RouterTableList.RouterTableListType[flag]
		return nil
	})
	return
}
