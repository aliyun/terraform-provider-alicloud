package alicloud

import (
	"log"

	"fmt"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
	"github.com/hashicorp/terraform/helper/schema"
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

func (client *AliyunClient) DescribeCenBandwidthLimit(cenId, localRegionId, oppositeRegionId string) (c cbn.CenInterRegionBandwidthLimit, err error) {
	request := cbn.CreateDescribeCenInterRegionBandwidthLimitsRequest()
	request.CenId = cenId

	for pageNum := 1; ; pageNum++ {
		request.PageNumber = requests.NewInteger(pageNum)
		resp, err := client.cenconn.DescribeCenInterRegionBandwidthLimits(request)

		if err != nil {
			log.Printf("Describe Cen Bandwidth Limit, CEN id %s, Err %s", cenId, err)
			return c, err
		}

		cenBandwidthLimitList := resp.CenInterRegionBandwidthLimits.CenInterRegionBandwidthLimit

		for limitNum := 0; limitNum <= len(cenBandwidthLimitList)-1; limitNum++ {
			ifMatch := cenBandwidthLimitList[limitNum].LocalRegionId == localRegionId && cenBandwidthLimitList[limitNum].OppositeRegionId == oppositeRegionId
			if !ifMatch {
				ifMatch = cenBandwidthLimitList[limitNum].LocalRegionId == oppositeRegionId && cenBandwidthLimitList[limitNum].OppositeRegionId == localRegionId
			}
			if ifMatch {
				return cenBandwidthLimitList[limitNum], nil
			}
		}

		if pageNum*resp.PageSize >= resp.TotalCount {
			return c, GetNotFoundErrorFromString(fmt.Sprintf("The specified CEN bandwith limit cenId %s localRegionId %s oppositeRegionId %s is not found", cenId, localRegionId, oppositeRegionId))
		}
	}
}

func (client *AliyunClient) SetCenInterRegionBandwidthLimit(d *schema.ResourceData, bandwidthLimit int) (err error) {
	request := cbn.CreateSetCenInterRegionBandwidthLimitRequest()
	request.CenId = d.Get("cen_id").(string)
	regionsId := d.Get("regions_id").(*schema.Set).List()
	request.LocalRegionId = regionsId[0].(string)
	request.OppositeRegionId = regionsId[1].(string)
	request.BandwidthLimit = requests.NewInteger(bandwidthLimit)

	_, err = client.cenconn.SetCenInterRegionBandwidthLimit(request)
	return err
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

func (client *AliyunClient) DescribeCenAttachedChildInstanceById(instanceId, cenId string) (c cbn.ChildInstance, err error) {
	request := cbn.CreateDescribeCenAttachedChildInstancesRequest()
	request.CenId = cenId

	for pageNum := 1; ; pageNum++ {
		request.PageNumber = requests.NewInteger(pageNum)
		response, err := client.cenconn.DescribeCenAttachedChildInstances(request)

		if err != nil {
			log.Printf("Cen Id %s, Err %s", cenId, err)
			return c, err
		}

		instanceList := response.ChildInstances.ChildInstance

		for instanceNum := 0; instanceNum <= len(instanceList)-1; instanceNum++ {
			if instanceList[instanceNum].ChildInstanceId == instanceId {
				return instanceList[instanceNum], nil
			}
		}

		if pageNum*response.PageSize >= response.TotalCount {
			return c, GetNotFoundErrorFromString(GetNotFoundMessage("CEN Attached Child Instance", instanceId))
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

func (client *AliyunClient) WaitForCenChildInstanceAttached(instanceId string, cenId string, status Status, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}

	for {
		instance, err := client.DescribeCenAttachedChildInstanceById(instanceId, cenId)
		if err != nil {
			return err
		}
		if instance.Status == string(status) {
			break
		}
		timeout = timeout - DefaultIntervalShort
		if timeout <= 0 {
			return GetTimeErrorFromString(GetTimeoutMessage("CEN Child Instance Attach", string(status)))
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}

	return nil
}

func (client *AliyunClient) WaitForCenChildInstanceDetached(instanceId string, cenId string, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}

	for {
		_, err := client.DescribeCenAttachedChildInstanceById(instanceId, cenId)
		if err != nil {
			if NotFoundError(err) {
				break
			} else {
				return err
			}
		}

		timeout = timeout - DefaultIntervalShort
		if timeout <= 0 {
			return GetTimeErrorFromString(fmt.Sprintf("Waitting for %s detach timeout.", instanceId))
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}

	return nil
}

func (client *AliyunClient) WaitForCenInterRegionBandwidthLimitActive(cenId string, localRegionId string, oppositeRegionId string, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}

	for {
		cenBandwidthLimit, err := client.DescribeCenBandwidthLimit(cenId, localRegionId, oppositeRegionId)
		if err != nil {
			return err
		}

		if cenBandwidthLimit.Status == string(Active) {
			break
		}

		timeout = timeout - DefaultIntervalShort
		if timeout <= 0 {
			return GetTimeErrorFromString(fmt.Sprintf("Waitting for bandwidth limit CenId %s localRegionId %s oppositeRegionId %s timeout.", cenId, localRegionId, oppositeRegionId))
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}

	return nil
}

func (client *AliyunClient) WaitForCenInterRegionBandwidthLimitDestroy(cenId string, localRegionId string, oppositeRegionId string, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}

	for {
		_, err := client.DescribeCenBandwidthLimit(cenId, localRegionId, oppositeRegionId)
		if err != nil {
			if NotFoundError(err) {
				return nil
			}
			return err
		}

		timeout = timeout - DefaultIntervalShort
		if timeout <= 0 {
			return GetTimeErrorFromString(fmt.Sprintf("Waitting for bandwidth limit CenId %s localRegionId %s oppositeRegionId %s timeout.", cenId, localRegionId, oppositeRegionId))
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
