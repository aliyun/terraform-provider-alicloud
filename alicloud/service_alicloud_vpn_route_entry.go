package alicloud

import (
	"crypto/md5"
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

type VpnRouteEntryService struct {
	client *connectivity.AliyunClient
}

type VpnState struct {
	State        string
	CreateTime   int64
	Status       Status
	VpnGatewayId string
	RouteDest    string
	NextHop      string
	Weight       int
	PublishVpc   string
}

func (s *VpnRouteEntryService) DescribeVpnRouteEntry(ids string) (v VpnState, err error) {
	request := vpc.CreateDescribeVpnRouteEntriesRequest()

	a := strings.Split(ids, ":")
	gatewayId := a[0]
	id := a[1]
	request.VpnGatewayId = gatewayId

	raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.DescribeVpnRouteEntries(request)
	})
	if err != nil {
		if IsExceptedErrors(err, []string{VpnForbidden, VpnNotFound}) {
			return v, WrapErrorf(Error(GetNotFoundMessage("VpnRouterEntry", id)), NotFoundMsg, ProviderERROR)
		}
		return v, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	response, _ := raw.(*vpc.DescribeVpnRouteEntriesResponse)

	for _, resp := range response.VpnRouteEntries.VpnRouteEntry {
		i := gatewayId + ":" + resp.NextHop + resp.RouteDest

		if ids == i {
			data := VpnState{
				resp.State,
				resp.CreateTime,
				Active,
				resp.VpnInstanceId,
				resp.RouteDest,
				resp.NextHop,
				resp.Weight,
				resp.State,
			}
			return data, nil
		}
	}
	return v, WrapErrorf(Error(GetNotFoundMessage("VpnRouterEntry", gatewayId)), NotFoundMsg, ProviderERROR)
}

func (s *VpnRouteEntryService) WaitForVpnRouteEntry(ids string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeVpnRouteEntry(ids)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}

		if strings.EqualFold(string(object.Status), string(status)) {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, ids, GetFunc(1), timeout, object.Status, string(status), ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
}

// func getMd5FromStr(str string) string {
// 	Md5Inst := md5.New()
// 	Md5Inst.Write([]byte(str))
// 	Result := Md5Inst.Sum([]byte(""))
// 	return fmt.Sprintf("%x", Result)
// }
