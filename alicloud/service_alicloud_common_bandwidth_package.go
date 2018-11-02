package alicloud

import (
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

type CommonBandwidthPackageService struct {
	client *connectivity.AliyunClient
}

func (s *CommonBandwidthPackageService) DescribeCommonBandwidthPackage(commonBandwidthPackageId string) (v vpc.CommonBandwidthPackage, err error) {
	request := vpc.CreateDescribeCommonBandwidthPackagesRequest()
	request.BandwidthPackageId = commonBandwidthPackageId

	invoker := NewInvoker()
	err = invoker.Run(func() error {
		raw, err := s.client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeCommonBandwidthPackages(request)
		})
		if err != nil {
			return err
		}
		resp, _ := raw.(*vpc.DescribeCommonBandwidthPackagesResponse)
		length := len(resp.CommonBandwidthPackages.CommonBandwidthPackage)
		if resp == nil || length <= 0 {
			return GetNotFoundErrorFromString(GetNotFoundMessage("CommonBandwidthPackage", commonBandwidthPackageId))
		}
		//Finding the commonBandwidthPackageId
		for _, id := range resp.CommonBandwidthPackages.CommonBandwidthPackage {
			if id.BandwidthPackageId == commonBandwidthPackageId {
				v = id
				return nil
			}
		}
		return GetNotFoundErrorFromString(GetNotFoundMessage("CommonBandwidthPackageId", commonBandwidthPackageId))
	})
	return
}

func (s *CommonBandwidthPackageService) WaitForCommonBandwidthPackage(commonBandwidthPackageId string, timeout int) error {

	if timeout <= 0 {
		timeout = DefaultTimeout
	}
	for {
		time.Sleep(DefaultIntervalShort * time.Second)
		resp, err := s.DescribeCommonBandwidthPackage(commonBandwidthPackageId)

		if err != nil {
			return err
		}
		if resp.BandwidthPackageId == commonBandwidthPackageId {
			return nil
		}
		timeout = timeout - DefaultIntervalShort
		if timeout <= 0 {
			return GetTimeErrorFromString(GetTimeoutMessage("CommonBandwidthPackage Attachment", string("Unavailable")))
		}
	}
}
