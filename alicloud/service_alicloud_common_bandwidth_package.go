package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/schema"

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

func (s *CommonBandwidthPackageService) WaitForCommonBandwidthPackageAttachment(bandwidthPackageId string, ipInstanceId string, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}
	for {
		err := s.DescribeCommonBandwidthPackageAttachment(bandwidthPackageId, ipInstanceId)

		if err != nil {
			if !NotFoundError(err) {
				return err
			}
		} else {
			break
		}
		timeout = timeout - DefaultIntervalShort
		if timeout <= 0 {
			return GetTimeErrorFromString(GetTimeoutMessage("Common Bandwidth Package Attachment", string("Unavailable")))
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *CommonBandwidthPackageService) DescribeCommonBandwidthPackageAttachment(bandwidthPackageId string, ipInstanceId string) (err error) {
	invoker := NewInvoker()
	return invoker.Run(func() error {
		commonBandwidthPackage, err := s.DescribeCommonBandwidthPackage(bandwidthPackageId)
		if err != nil {
			return err
		}
		for _, id := range commonBandwidthPackage.PublicIpAddresses.PublicIpAddresse {
			if ipInstanceId == id.AllocationId {
				return nil
			}
		}
		return GetNotFoundErrorFromString(GetNotFoundMessage("CommonBandwidthPackageAttachment", bandwidthPackageId+COLON_SEPARATED+ipInstanceId))
	})
}

func GetBandwidthPackageIdAndIpInstanceId(d *schema.ResourceData, meta interface{}) (string, string, error) {
	parts := strings.Split(d.Id(), COLON_SEPARATED)

	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid resource id")
	}
	return parts[0], parts[1], nil
}
