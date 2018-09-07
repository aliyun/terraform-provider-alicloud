package alicloud

import (
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
)

func (client *AliyunClient) DescribeCenInstance(cenId string) (c cbn.Cen, err error) {
	request := cbn.CreateDescribeCensRequest()

	values := []string{cenId}
	filters := []cbn.DescribeCensFilter{cbn.DescribeCensFilter{
		Key:   "CenId",
		Value: &values,
	}}

	request.Filter = &filters

	invoker := NewInvoker()
	err = invoker.Run(func() error {
		resp, err := client.cenconn.DescribeCens(request)
		if err != nil {
			if IsExceptedError(err, ParameterCenInstanceIdNotExist) {
				return GetNotFoundErrorFromString(GetNotFoundMessage("CEN Instance", cenId))
			}
			return err
		}
		if resp == nil || len(resp.Cens.Cen) <= 0 || resp.Cens.Cen[0].CenId != cenId {
			return GetNotFoundErrorFromString(GetNotFoundMessage("CEN Instance", cenId))
		}
		c = resp.Cens.Cen[0]
		return nil
	})

	return
}

func (client *AliyunClient) WaitForCenInstance(cenId string, status Status, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}

	for {
		cen, err := client.DescribeCenInstance(cenId)
		if err != nil {
			return err
		}
		if cen.Status == string(status) {
			break
		}
		timeout = timeout - DefaultIntervalShort
		if timeout <= 0 {
			return GetTimeErrorFromString(GetTimeoutMessage("CEN", string(status)))
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}

	return nil
}
