package alicloud

import (
	"log"

	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
)

func (client *AliyunClient) DescribeCen(cenId string) (c cbn.Cen, err error) {
	request := cbn.CreateDescribeCensRequest()

	for pageNum := 1; ; pageNum++ {
		request.PageNumber = requests.NewInteger(pageNum)
		resp, err := client.cenconn.DescribeCens(request)

		if err != nil {
			log.Printf("Cen Id %s, Err %s", cenId, err)
			return c, err
		}

		cenList := resp.Cens.Cen

		for cenNum := 0; cenNum <= len(cenList)-1; cenNum++ {
			if cenList[cenNum].CenId == cenId {
				return cenList[cenNum], nil
			}
		}

		if pageNum*resp.PageSize >= resp.TotalCount {
			return c, GetNotFoundErrorFromString(GetNotFoundMessage("CEN", cenId))
		}

	}
}

func (client *AliyunClient) WaitForCen(cenId string, status Status, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}

	for {
		cen, err := client.DescribeCen(cenId)
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
