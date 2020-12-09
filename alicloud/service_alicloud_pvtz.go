package alicloud

import (
	"strconv"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/pvtz"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
)

type PvtzService struct {
	client *connectivity.AliyunClient
}

func (s *PvtzService) DescribePvtzZone(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewPvtzClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeZoneInfo"
	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
		"ZoneId":   id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-01-01"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		if IsExpectedErrors(err, []string{"Zone.Invalid.Id", "Zone.Invalid.UserId", "Zone.NotExists", "ZoneVpc.NotExists.VpcId"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("PvtzZone", id)), NotFoundMsg, ProviderERROR)
			return object, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return object, err
	}
	addDebug(action, response, request)
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *PvtzService) DescribePvtzZoneAttachment(id string) (object map[string]interface{}, err error) {
	object, err = s.DescribePvtzZone(id)
	if err != nil {
		err = WrapError(err)
		return
	}
	vpcs := object["BindVpcs"].(map[string]interface{})["Vpc"].([]interface{})
	if len(vpcs) < 1 {
		err = WrapErrorf(Error(GetNotFoundMessage("PvtzZoneAttachment", id)), NotFoundMsg, ProviderERROR)
	}

	return
}

func (s *PvtzService) WaitForZoneAttachment(id string, vpcIdMap map[string]string, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	var vpcId string
	for {
		object, err := s.DescribePvtzZoneAttachment(id)
		if err != nil && !NotFoundError(err) {
			return WrapError(err)
		}

		equal := true
		vpcs := object["BindVpcs"].(map[string]interface{})["Vpc"].([]interface{})
		if len(vpcs) == len(vpcIdMap) {
			for _, vpc := range vpcs {
				vpc := vpc.(map[string]interface{})
				if _, ok := vpcIdMap[vpc["VpcId"].(string)]; !ok {
					equal = false
					vpcId = vpc["VpcId"].(string)
					break
				}
			}
		}
		if equal {
			break
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, "", vpcId, ProviderERROR)
		}
	}
	return nil
}

func (s *PvtzService) DescribePvtzZoneRecord(id string) (record pvtz.Record, err error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return record, WrapError(err)
	}
	request := pvtz.CreateDescribeZoneRecordsRequest()
	request.RegionId = s.client.RegionId
	request.ZoneId = parts[1]
	request.PageNumber = requests.NewInteger(1)
	request.PageSize = requests.NewInteger(PageSizeLarge)

	recordIdStr := parts[0]
	var response *pvtz.DescribeZoneRecordsResponse

	for {
		err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			raw, err := s.client.WithPvtzClient(func(pvtzClient *pvtz.Client) (interface{}, error) {
				return pvtzClient.DescribeZoneRecords(request)
			})
			if err != nil {
				if IsExpectedErrors(err, []string{ServiceUnavailable, ThrottlingUser, "System.Busy"}) {
					time.Sleep(5 * time.Second)
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			response, _ = raw.(*pvtz.DescribeZoneRecordsResponse)
			return nil
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"Zone.NotExists", "ZoneVpc.NotExists.VpcId"}) {
				return record, WrapErrorf(Error(GetNotFoundMessage("ZoneRecord", id)), NotFoundMsg, AlibabaCloudSdkGoERROR)
			}
			return record, WrapErrorf(err, DefaultErrorMsg, recordIdStr, request.GetActionName(), AlibabaCloudSdkGoERROR)
		}

		if len(response.Records.Record) < 1 {
			return record, WrapErrorf(Error(GetNotFoundMessage("ZoneRecord", id)), NotFoundMsg, ProviderERROR)
		}

		for _, rec := range response.Records.Record {
			if strconv.FormatInt(rec.RecordId, 10) == parts[0] {
				record = rec
				return record, nil
			}
		}
		if len(response.Records.Record) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(request.PageNumber); err != nil {
			return record, WrapError(err)
		} else {
			request.PageNumber = page
		}
	}

	return record, WrapErrorf(Error(GetNotFoundMessage("ZoneRecord", recordIdStr)), NotFoundMsg, ProviderERROR)

}

func (s *PvtzService) WaitForPvtzZone(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

	for {
		object, err := s.DescribePvtzZone(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if object["ZoneId"] == id {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object["ZoneId"], id, ProviderERROR)
		}

	}
}

func (s *PvtzService) WaitForPvtzZoneAttachment(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

	for {
		object, err := s.DescribePvtzZoneAttachment(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if object["ZoneId"] == id {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object["ZoneId"], id, ProviderERROR)
		}

	}
}

func (s *PvtzService) WaitForPvtzZoneRecord(id string, status Status, timeout int) error {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return WrapError(err)
	}
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

	for {
		object, err := s.DescribePvtzZoneRecord(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if strconv.FormatInt(object.RecordId, 10) == parts[0] {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, strconv.FormatInt(object.RecordId, 10), id, ProviderERROR)
		}

	}
}
