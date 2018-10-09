package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/schema"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
)

func (client *AliyunClient) DescribeHaVip(haVipId string) (v vpc.HaVip, err error) {
	request := vpc.CreateDescribeHaVipsRequest()
	values := []string{haVipId}
	filter := []vpc.DescribeHaVipsFilter{vpc.DescribeHaVipsFilter{
		Key:   "HaVipId",
		Value: &values,
	},
	}
	request.Filter = &filter

	invoker := NewInvoker()
	err = invoker.Run(func() error {
		resp, err := client.vpcconn.DescribeHaVips(request)
		if err != nil {
			return err
		}
		if resp == nil || len(resp.HaVips.HaVip) <= 0 ||
			resp.HaVips.HaVip[0].HaVipId != haVipId {
			return GetNotFoundErrorFromString(GetNotFoundMessage("HaVip", haVipId))
		}
		v = resp.HaVips.HaVip[0]
		return nil
	})
	return
}

func (client *AliyunClient) WaitForHaVip(haVipId string, status Status, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}
	for {
		//wait the order effective
		havip, err := client.DescribeHaVip(haVipId)
		if err != nil {
			return err
		}
		if strings.ToLower(havip.Status) == strings.ToLower(string(status)) {
			break
		}
		timeout = timeout - DefaultIntervalShort
		if timeout <= 0 {
			return GetTimeErrorFromString(GetTimeoutMessage("HaVip", string(status)))
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}

func (client *AliyunClient) DescribeHaVipAttachment(haVipId string, instanceId string) (err error) {
	invoker := NewInvoker()
	return invoker.Run(func() error {
		haVip, err := client.DescribeHaVip(haVipId)
		if err != nil {
			return err
		}
		for _, id := range haVip.AssociatedInstances.AssociatedInstance {
			if id == instanceId {
				return nil
			}
		}
		return GetNotFoundErrorFromString(GetNotFoundMessage("HaVipAttachment", haVipId+COLON_SEPARATED+instanceId))
	})
}

func (client *AliyunClient) WaitForHaVipAttachment(haVipId string, instanceId string, timeout int) error {

	if timeout <= 0 {
		timeout = DefaultTimeout
	}
	for {
		err := client.DescribeHaVipAttachment(haVipId, instanceId)

		if err != nil {
			if !NotFoundError(err) {
				return err
			}
		} else {
			break
		}
		timeout = timeout - DefaultIntervalShort
		if timeout <= 0 {
			return GetTimeErrorFromString(GetTimeoutMessage("HaVip Attachment", string("Unavailable")))
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}

func getHaVipIdAndInstanceId(d *schema.ResourceData, meta interface{}) (string, string, error) {
	parts := strings.Split(d.Id(), COLON_SEPARATED)

	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid resource id")
	}
	return parts[0], parts[1], nil
}
