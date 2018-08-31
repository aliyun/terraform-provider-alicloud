package alicloud

import (
	"fmt"
	"log"
	"strings"

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

func (client *AliyunClient) DescribeCenBandwidthPackage(cenBwpId string) (c cbn.CenBandwidthPackage, err error) {
	request := cbn.CreateDescribeCenBandwidthPackagesRequest()

	for pageNum := 1; ; pageNum++ {
		request.PageNumber = requests.NewInteger(pageNum)
		resp, err := client.cenconn.DescribeCenBandwidthPackages(request)

		if err != nil {
			log.Printf("CEN Bandwidth Package Id %s, Err %s", cenBwpId, err)
			return c, err
		}

		cenBwpList := resp.CenBandwidthPackages.CenBandwidthPackage

		for cenNum := 0; cenNum <= len(cenBwpList)-1; cenNum++ {
			if cenBwpList[cenNum].CenBandwidthPackageId == cenBwpId {
				return cenBwpList[cenNum], nil
			}
		}

		if pageNum*resp.PageSize >= resp.TotalCount {
			return c, GetNotFoundErrorFromString(GetNotFoundMessage("CEN Bandwidth Package", cenBwpId))
		}
	}
}
func (client *AliyunClient) DescribeCenBandwidthPackageById(cenBwpId, cenId string) (c cbn.CenBandwidthPackage, err error) {
	resp, err := client.DescribeCenBandwidthPackage(cenBwpId)

	if err != nil {
		return c, err
	}

	if len(resp.CenIds.CenId) == 0 || cenId != resp.CenIds.CenId[0] {
		return c, GetNotFoundErrorFromString(GetNotFoundMessage("CEN Bandwidth Package", cenBwpId))
	}

	return resp, nil
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
func (client *AliyunClient) WaitForCenBandwidthPackage(cenBwpId string, status Status, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}

	for {
		cenBwp, err := client.DescribeCenBandwidthPackage(cenBwpId)
		if err != nil && !NotFoundError(err) {
			return err
		}
		if cenBwp.Status == string(status) {
			break
		}
		timeout = timeout - DefaultIntervalShort
		if timeout <= 0 {
			return GetTimeErrorFromString(GetTimeoutMessage("CEN Bandwidth Package", string(status)))
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}

	return nil
}
func (client *AliyunClient) WaitForCenBandwidthPackageUpdate(cenBwpId string, bandwidth int, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}

	for {
		cenBwp, err := client.DescribeCenBandwidthPackage(cenBwpId)
		if err != nil {
			return err
		}
		if cenBwp.Bandwidth == bandwidth {
			break
		}

		timeout = timeout - DefaultIntervalShort
		if timeout <= 0 {
			return GetTimeErrorFromString(fmt.Sprintf("BandwidthPackage Update timeout"))
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}

	return nil
}
func (client *AliyunClient) WaitForCenBandwidthPackageAssociate(cenBwpId string, status Status, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}

	for {
		cenBwp, err := client.DescribeCenBandwidthPackage(cenBwpId)
		if err != nil {
			return err
		}
		if cenBwp.Status == string(status) {
			break
		}
		timeout = timeout - DefaultIntervalShort
		if timeout <= 0 {
			return GetTimeErrorFromString(GetTimeoutMessage("CEN Bandwidth Package Associate", string(status)))
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}

	return nil
}
func getCenIdAndAnotherId(id string) (string, string, error) {
	parts := strings.Split(id, ":")

	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid resource id")
	}
	return parts[0], parts[1], nil
}
