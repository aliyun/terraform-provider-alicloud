package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type AlidnsService struct {
	client *connectivity.AliyunClient
}

func (s *AlidnsService) DescribeAlidnsDomainGroup(id string) (object alidns.DomainGroup, err error) {
	request := alidns.CreateDescribeDomainGroupsRequest()
	request.RegionId = s.client.RegionId

	request.PageNumber = requests.NewInteger(1)
	request.PageSize = requests.NewInteger(20)
	for {

		raw, err := s.client.WithAlidnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
			return alidnsClient.DescribeDomainGroups(request)
		})
		if err != nil {
			err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
			return object, err
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*alidns.DescribeDomainGroupsResponse)

		if len(response.DomainGroups.DomainGroup) < 1 {
			err = WrapErrorf(NotFoundErr("AlidnsDomainGroup", id), NotFoundMsg, ProviderERROR, response.RequestId)
			return object, err
		}
		for _, object := range response.DomainGroups.DomainGroup {
			if object.GroupId == id {
				return object, nil
			}
		}
		if len(response.DomainGroups.DomainGroup) < PageSizeMedium {
			break
		}
		if page, err := getNextpageNumber(request.PageNumber); err != nil {
			return object, WrapError(err)
		} else {
			request.PageNumber = page
		}
	}
	err = WrapErrorf(NotFoundErr("AlidnsDomainGroup", id), NotFoundMsg, ProviderERROR)
	return
}

func (s *AlidnsService) DescribeAlidnsRecord(id string) (object alidns.DescribeDomainRecordInfoResponse, err error) {
	request := alidns.CreateDescribeDomainRecordInfoRequest()
	request.RegionId = s.client.RegionId

	request.RecordId = id

	raw, err := s.client.WithAlidnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
		return alidnsClient.DescribeDomainRecordInfo(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"DomainRecordNotBelongToUser", "InvalidRR.NoExist"}) {
			err = WrapErrorf(NotFoundErr("AlidnsRecord", id), NotFoundMsg, ProviderERROR)
			return
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*alidns.DescribeDomainRecordInfoResponse)
	return *response, nil
}

func (s *AlidnsService) ListTagResources(id string) (object alidns.ListTagResourcesResponse, err error) {
	request := alidns.CreateListTagResourcesRequest()
	request.RegionId = s.client.RegionId

	request.ResourceType = "DOMAIN"
	request.ResourceId = &[]string{id}

	raw, err := s.client.WithAlidnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
		return alidnsClient.ListTagResources(request)
	})
	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*alidns.ListTagResourcesResponse)
	return *response, nil
}

func (s *AlidnsService) SetResourceTags(d *schema.ResourceData, resourceType string) error {
	oldItems, newItems := d.GetChange("tags")
	added := make([]alidns.TagResourcesTag, 0)
	for key, value := range newItems.(map[string]interface{}) {
		added = append(added, alidns.TagResourcesTag{
			Key:   key,
			Value: value.(string),
		})
	}
	removed := make([]string, 0)
	for key := range oldItems.(map[string]interface{}) {
		removed = append(removed, key)
	}
	if len(removed) > 0 {
		request := alidns.CreateUntagResourcesRequest()
		request.RegionId = s.client.RegionId
		request.ResourceId = &[]string{d.Id()}
		request.ResourceType = resourceType
		request.TagKey = &removed
		raw, err := s.client.WithAlidnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
			return alidnsClient.UntagResources(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}
	if len(added) > 0 {
		request := alidns.CreateTagResourcesRequest()
		request.RegionId = s.client.RegionId
		request.ResourceId = &[]string{d.Id()}
		request.ResourceType = resourceType
		request.Tag = &added
		raw, err := s.client.WithAlidnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
			return alidnsClient.TagResources(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}
	return nil
}

func (s *AlidnsService) DescribeAlidnsDomain(id string) (object alidns.DescribeDomainInfoResponse, err error) {
	request := alidns.CreateDescribeDomainInfoRequest()
	request.RegionId = s.client.RegionId

	request.DomainName = id

	raw, err := s.client.WithAlidnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
		return alidnsClient.DescribeDomainInfo(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDomainName.NoExist"}) {
			err = WrapErrorf(NotFoundErr("AlidnsDomain", id), NotFoundMsg, ProviderERROR)
			return
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*alidns.DescribeDomainInfoResponse)
	return *response, nil
}

func (s *AlidnsService) DescribeAlidnsInstance(id string) (object map[string]interface{}, err error) {
	client := s.client
	var response map[string]interface{}
	action := "DescribeDnsProductInstance"
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"InstanceId": id,
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(11*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Alidns", "2015-01-09", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDnsProduct"}) {
			return object, WrapErrorf(NotFoundErr("Alidns:Instance", id), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *AlidnsService) DescribeAlidnsCustomLine(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeCustomLine"
	request := map[string]interface{}{
		"LineId": id,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Alidns", "2015-01-09", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"DnsCustomLine.NotExists"}) {
			return object, WrapErrorf(NotFoundErr("Alidns::CustomLine", id), NotFoundMsg, ProviderERROR)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *AlidnsService) DescribeCustomLine(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeCustomLine"
	request := map[string]interface{}{
		"LineId": id,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Alidns", "2015-01-09", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"DnsCustomLine.NotExists"}) {
			return object, WrapErrorf(NotFoundErr("Alidns::CustomLine", id), NotFoundMsg, ProviderERROR)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *AlidnsService) DescribeAlidnsGtmInstance(id string) (object map[string]interface{}, err error) {
	client := s.client
	var response map[string]interface{}
	action := "DescribeDnsGtmInstance"
	request := map[string]interface{}{
		"InstanceId": id,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Alidns", "2015-01-09", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"DnsGtmInstance.NotExists"}) {
			return object, WrapErrorf(NotFoundErr("Alidns:GtmInstance", id), NotFoundWithResponse, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *AlidnsService) DescribeAlidnsAddressPool(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeDnsGtmInstanceAddressPool"
	request := map[string]interface{}{
		"AddrPoolId": id,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Alidns", "2015-01-09", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"DnsGtmAddrPool.NotExists"}) {
			return object, WrapErrorf(NotFoundErr("Alidns::AddressPool", id), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *AlidnsService) DescribeAlidnsAccessStrategy(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeDnsGtmAccessStrategy"
	request := map[string]interface{}{
		"StrategyId": id,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Alidns", "2015-01-09", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"DnsGtmAccessStrategy.NotExists"}) {
			return object, WrapErrorf(NotFoundErr("Alidns:DnsGtmAccessStrategy", id), NotFoundWithResponse, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *AlidnsService) DescribeAlidnsMonitorConfig(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeDnsGtmMonitorConfig"
	request := map[string]interface{}{
		"MonitorConfigId": id,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Alidns", "2015-01-09", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *AlidnsService) DescribeDomainRecordById(recordID, domainName string) (*alidns.Record, error) {
	request := alidns.CreateDescribeDomainRecordInfoRequest()
	request.RecordId = recordID

	response, err := s.client.WithAlidnsClient(func(client *alidns.Client) (interface{}, error) {
		return client.DescribeDomainRecordInfo(request)
	})
	if err != nil {
		return nil, WrapError(err)
	}

	resp, ok := response.(*alidns.DescribeDomainRecordInfoResponse)
	if !ok {
		return nil, fmt.Errorf("unexpected response type for DescribeDomainRecordInfo")
	}

	// Fetch weight from DescribeDomainRecords
	records, err := s.DescribeDomainRecords(domainName)
	if err != nil {
		return nil, WrapError(err)
	}

	var weight int
	for _, r := range records {
		if r.RecordId == recordID {
			weight = r.Weight
			break
		}
	}

	return &alidns.Record{
		RecordId:   resp.RecordId,
		DomainName: domainName,
		RR:         resp.RR,
		Type:       resp.Type,
		Value:      resp.Value,
		TTL:        resp.TTL,
		Line:       resp.Line,
		Priority:   resp.Priority,
		Remark:     resp.Remark,
		Status:     resp.Status,
		Weight:     weight, // Set weight
	}, nil
}

func (s *AlidnsService) UpdateRecordRemark(recordID, remark string) error {
	request := alidns.CreateUpdateDomainRecordRemarkRequest()
	request.RecordId = recordID
	request.Remark = remark

	log.Printf("[DEBUG] Setting record remark: RecordId=%s, Remark=%s", recordID, remark)

	_, err := s.client.WithAlidnsClient(func(client *alidns.Client) (interface{}, error) {
		return client.UpdateDomainRecordRemark(request)
	})
	if err != nil {
		return WrapError(err)
	}

	return nil
}

func (s *AlidnsService) IsOperationFinished(recordID string) (bool, error) {
	request := alidns.CreateDescribeDomainRecordInfoRequest()
	request.RecordId = recordID

	response, err := s.client.WithAlidnsClient(func(client *alidns.Client) (interface{}, error) {
		return client.DescribeDomainRecordInfo(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"LastOperationNotFinished"}) {
			return false, nil
		}
		return false, WrapError(err)
	}

	resp, ok := response.(*alidns.DescribeDomainRecordInfoResponse)
	if !ok {
		return false, fmt.Errorf("unexpected response type for DescribeDomainRecordInfo")
	}

	// Check operation status based on API response (adjust this as per API documentation)
	return resp.RecordId != "", nil
}

func (s *AlidnsService) DescribeDomainRecords(domainName string) ([]alidns.Record, error) {
	request := alidns.CreateDescribeDomainRecordsRequest()
	request.RegionId = s.client.RegionId
	request.DomainName = domainName

	response, err := s.client.WithAlidnsClient(func(client *alidns.Client) (interface{}, error) {
		return client.DescribeDomainRecords(request)
	})
	if err != nil {
		return nil, WrapError(err)
	}

	resp, ok := response.(*alidns.DescribeDomainRecordsResponse)
	if !ok {
		return nil, fmt.Errorf("unexpected response type for DescribeDomainRecords")
	}

	return resp.DomainRecords.Record, nil
}

func (s *AlidnsService) UpdateRecordStatus(recordID, status string) error {
	request := alidns.CreateSetDomainRecordStatusRequest()
	request.RecordId = recordID
	request.Status = status

	_, err := s.client.WithAlidnsClient(func(client *alidns.Client) (interface{}, error) {
		return client.SetDomainRecordStatus(request)
	})
	if err != nil {
		return WrapError(err)
	}

	return nil
}

func (s *AlidnsService) SetRecordWeight(recordID string, weight int, domainName, rr string) error {
	// Ensure WRR is enabled before setting weight
	wrrStatus, err := s.GetWRRStatus(domainName, rr)
	if err != nil {
		return WrapError(err)
	}

	if wrrStatus != "ENABLE" {
		log.Printf("[DEBUG] WRR is disabled for subdomain %s.%s. Enabling it now.", rr, domainName)
		err = s.SetWRRStatus(domainName, rr, "ENABLE")
		if err != nil {
			return WrapError(err)
		}
	}

	request := alidns.CreateUpdateDNSSLBWeightRequest()
	request.RegionId = s.client.RegionId
	request.RecordId = recordID
	request.Weight = requests.NewInteger(weight)

	log.Printf("[DEBUG] Setting record weight: RecordId=%s, Weight=%d", recordID, weight)

	_, err = s.client.WithAlidnsClient(func(client *alidns.Client) (interface{}, error) {
		resp, err := client.UpdateDNSSLBWeight(request)
		if err != nil {
			log.Printf("[ERROR] Failed to update weight: %s", err)
		} else {
			log.Printf("[DEBUG] Successfully updated weight: %+v", resp)
		}
		return resp, err
	})

	return WrapError(err)
}

func (s *AlidnsService) SetRecordStatus(recordID, status string) error {
	request := alidns.CreateSetDomainRecordStatusRequest()
	request.RecordId = recordID
	request.Status = status // "ENABLE" or "DISABLE"

	log.Printf("[DEBUG] Setting record status: RecordId=%s, Status=%s", recordID, status)

	_, err := s.client.WithAlidnsClient(func(client *alidns.Client) (interface{}, error) {
		return client.SetDomainRecordStatus(request)
	})
	if err != nil {
		return WrapError(err)
	}

	return nil
}

func (s *AlidnsService) SetWRRStatus(domainName, rr, status string) error {
	request := alidns.CreateSetDNSSLBStatusRequest()
	request.RegionId = s.client.RegionId
	request.DomainName = domainName
	request.SubDomain = fmt.Sprintf("%s.%s", rr, domainName)
	request.Open = requests.NewBoolean(status == "ENABLE") // Convert status to boolean

	log.Printf("[DEBUG] Setting WRR status: SubDomain=%s, Status=%s", request.SubDomain, status)

	_, err := s.client.WithAlidnsClient(func(client *alidns.Client) (interface{}, error) {
		return client.SetDNSSLBStatus(request)
	})
	if err != nil {
		return WrapError(err)
	}

	return nil
}

func (s *AlidnsService) GetWRRStatus(domainName, rr string) (string, error) {
	request := alidns.CreateDescribeDNSSLBSubDomainsRequest()
	request.RegionId = s.client.RegionId
	request.DomainName = domainName

	log.Printf("[DEBUG] Fetching WRR status for domain: %s, subdomain: %s", domainName, rr)

	response, err := s.client.WithAlidnsClient(func(client *alidns.Client) (interface{}, error) {
		return client.DescribeDNSSLBSubDomains(request)
	})
	if err != nil {
		return "", WrapError(err)
	}

	resp, ok := response.(*alidns.DescribeDNSSLBSubDomainsResponse)
	if !ok {
		return "", fmt.Errorf("unexpected response type for DescribeDNSSLBSubDomains")
	}

	for _, sub := range resp.SlbSubDomains.SlbSubDomain {
		if sub.SubDomain == fmt.Sprintf("%s.%s", rr, domainName) {
			if sub.Open {
				return "ENABLE", nil
			}
			return "DISABLE", nil
		}
	}

	log.Printf("[DEBUG] Subdomain %s.%s not found, defaulting WRR status to DISABLE", rr, domainName)
	return "DISABLE", nil
}

func (s *AlidnsService) WaitForRecordReady(recordID, domainName, rr string, timeout time.Duration) error {
	start := time.Now()
	for {
		if time.Since(start) > timeout {
			return fmt.Errorf("timeout while waiting for record %s to be ready", recordID)
		}

		wrrStatus, err := s.GetWRRStatus(domainName, rr)
		if err != nil {
			log.Printf("[DEBUG] Waiting for record readiness: RecordId=%s, Error=%s", recordID, err)
		} else if wrrStatus != "" {
			log.Printf("[DEBUG] Record %s is ready with WRR status: %s", recordID, wrrStatus)
			return nil
		}

		time.Sleep(2 * time.Second) // Polling interval
	}
}

func (s *AlidnsService) WaitForLastOperation(recordID, domainName string, timeout time.Duration) error {
	if domainName == "" {
		return fmt.Errorf("domainName is required for WaitForLastOperation")
	}

	start := time.Now()
	for {
		if time.Since(start) > timeout {
			return fmt.Errorf("timeout while waiting for last operation to finish for record ID %s in domain %s", recordID, domainName)
		}

		// Use DescribeDomainRecords instead of DescribeDomainRecordById to avoid issues with missing domainName
		records, err := s.DescribeDomainRecords(domainName)
		if err != nil {
			if IsExpectedErrors(err, []string{"LastOperationNotFinished"}) {
				log.Printf("[DEBUG] Last operation not finished for record ID %s, retrying...", recordID)
				time.Sleep(2 * time.Second)
				continue
			}
			return WrapError(err)
		}

		// Check if the record exists and is ready
		for _, record := range records {
			if record.RecordId == recordID {
				log.Printf("[DEBUG] Last operation finished for record ID %s in domain %s", recordID, domainName)
				return nil
			}
		}

		log.Printf("[DEBUG] Record ID %s not found in domain %s, retrying...", recordID, domainName)
		time.Sleep(2 * time.Second)
	}
}

func (s *AlidnsService) ListAllDomains() ([]alidns.DomainInDescribeDomains, error) {
	var allDomains []alidns.DomainInDescribeDomains
	pageNumber := 1
	pageSize := 50

	for {
		request := alidns.CreateDescribeDomainsRequest()
		request.RegionId = s.client.RegionId
		request.PageNumber = requests.NewInteger(pageNumber)
		request.PageSize = requests.NewInteger(pageSize)

		response, err := s.client.WithAlidnsClient(func(client *alidns.Client) (interface{}, error) {
			return client.DescribeDomains(request)
		})
		if err != nil {
			return nil, WrapError(err)
		}

		resp, ok := response.(*alidns.DescribeDomainsResponse)
		if !ok {
			return nil, fmt.Errorf("unexpected response type for DescribeDomains")
		}

		allDomains = append(allDomains, resp.Domains.Domain...)
		if len(resp.Domains.Domain) < pageSize {
			break
		}
		pageNumber++
	}

	return allDomains, nil
}

func (s *AlidnsService) GetSLBSubDomains(domainName string) ([]alidns.SlbSubDomain, error) {
	request := alidns.CreateDescribeDNSSLBSubDomainsRequest()
	request.DomainName = domainName

	response, err := s.client.WithAlidnsClient(func(client *alidns.Client) (interface{}, error) {
		return client.DescribeDNSSLBSubDomains(request)
	})
	if err != nil {
		return nil, WrapError(err)
	}

	resp, ok := response.(*alidns.DescribeDNSSLBSubDomainsResponse)
	if !ok {
		return nil, fmt.Errorf("unexpected response type for DescribeDNSSLBSubDomains")
	}

	return resp.SlbSubDomains.SlbSubDomain, nil
}

func (s *AlidnsService) GetWeight(recordID, domainName string) (int, error) {
	records, err := s.DescribeDomainRecords(domainName)
	if err != nil {
		return 0, WrapError(err)
	}

	for _, record := range records {
		if record.RecordId == recordID {
			return record.Weight, nil
		}
	}

	return 0, fmt.Errorf("record with ID %s not found in domain %s", recordID, domainName)
}

// fetchWRRStatus retrieves the WRR status for a given domain and RR.
func fetchWRRStatus(client *connectivity.AliyunClient, domainName, rr string) (string, error) {
	request := alidns.CreateDescribeDNSSLBSubDomainsRequest()
	request.DomainName = domainName

	raw, err := client.WithAlidnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
		return alidnsClient.DescribeDNSSLBSubDomains(request)
	})
	if err != nil {
		return "", err
	}

	response := raw.(*alidns.DescribeDNSSLBSubDomainsResponse)
	for _, subDomain := range response.SlbSubDomains.SlbSubDomain {
		if subDomain.SubDomain == fmt.Sprintf("%s.%s", rr, domainName) {
			if subDomain.Open {
				return "ENABLE", nil
			}
			return "DISABLE", nil
		}
	}

	// Default to DISABLE if not found
	return "DISABLE", nil
}
