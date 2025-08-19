package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudAlidnsRecordWeight() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudAlidnsRecordWeightCreate,
		Read:   resourceAliCloudAlidnsRecordWeightRead,
		Update: resourceAliCloudAlidnsRecordWeightUpdate,
		Delete: resourceAliCloudAlidnsRecordWeightDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Minute),
			Update: schema.DefaultTimeout(3 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"record_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"weight": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"wrr_status": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudAlidnsRecordWeightCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	recordID := d.Get("record_id").(string)
	weight := d.Get("weight").(int)

	// Retrieve subdomain and domain
	subDomain, domainName, err := getSubDomainAndDomain(client, recordID)
	if err != nil {
		return WrapError(err)
	}

	// Check and set WRR status to true if it's not already enabled
	wrrStatus, err := getWRRStatus(client, subDomain, domainName)
	if err != nil {
		return WrapError(err)
	}

	if !wrrStatus { // If WRR is not enabled, enable it
		if err := setWRRStatus(client, subDomain, domainName, true); err != nil {
			return WrapError(err)
		}
	}

	// Set weight for the record
	if err := setWeight(client, recordID, weight); err != nil {
		return WrapError(err)
	}

	d.SetId(recordID)
	return resourceAliCloudAlidnsRecordWeightRead(d, meta)
}

func resourceAliCloudAlidnsRecordWeightRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	recordID := d.Id()

	if err := d.Set("record_id", recordID); err != nil {
		return WrapError(fmt.Errorf("failed to set record_id during read: %v", err))
	}

	subDomain, domainName, err := getSubDomainAndDomain(client, recordID)
	if err != nil {
		return WrapError(fmt.Errorf("failed to retrieve subdomain and domain: %v", err))
	}

	wrrStatus, err := getWRRStatus(client, subDomain, domainName)
	if err != nil {
		return WrapError(fmt.Errorf("failed to retrieve WRR status: %v", err))
	}
	if err := d.Set("wrr_status", wrrStatus); err != nil {
		return WrapError(fmt.Errorf("failed to set WRR status: %v", err))
	}

	weight, err := getWeight(client, recordID)
	if err != nil {
		return WrapError(fmt.Errorf("failed to retrieve weight: %v", err))
	}
	if err := d.Set("weight", weight); err != nil {
		return WrapError(fmt.Errorf("failed to set weight: %v", err))
	}

	return nil
}

func resourceAliCloudAlidnsRecordWeightUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	recordID := d.Id()

	if d.HasChange("weight") {
		weight := d.Get("weight").(int)

		// Retrieve subdomain and domain
		subDomain, domainName, err := getSubDomainAndDomain(client, recordID)
		if err != nil {
			return WrapError(err)
		}

		// Check and set WRR status to true if it's not already enabled
		wrrStatus, err := getWRRStatus(client, subDomain, domainName)
		if err != nil {
			return WrapError(err)
		}

		if !wrrStatus { // If WRR is not enabled, enable it
			if err := setWRRStatus(client, subDomain, domainName, true); err != nil {
				return WrapError(err)
			}
		}

		// Update weight for the record
		if err := setWeight(client, recordID, weight); err != nil {
			return WrapError(err)
		}
	}

	return resourceAliCloudAlidnsRecordWeightRead(d, meta)
}

func resourceAliCloudAlidnsRecordWeightDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	recordID := d.Id()

	// Retrieve subdomain and domain
	subDomain, domainName, err := getSubDomainAndDomain(client, recordID)
	if err != nil {
		return WrapError(err)
	}

	// Set WRR status to false
	log.Printf("[DEBUG] Disabling WRR Status for subdomain: %s in domain: %s", subDomain, domainName)
	if err := setWRRStatus(client, subDomain, domainName, false); err != nil {
		if IsExpectedErrors(err, []string{"DisableDNSSLB"}) {
			log.Printf("[WARN] WRR is already disabled for subdomain: %s in domain: %s", subDomain, domainName)
		} else {
			return WrapError(err)
		}
	}

	d.SetId("") // Clear the resource ID as the resource is deleted
	return nil
}

func resourceAliCloudAlidnsRecordWeightImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(*connectivity.AliyunClient)
	recordID := d.Id()

	// Set the record_id to the imported ID
	if err := d.Set("record_id", recordID); err != nil {
		return nil, WrapError(fmt.Errorf("failed to set record_id during import: %v", err))
	}

	// Retrieve the subdomain and domain
	subDomain, domainName, err := getSubDomainAndDomain(client, recordID)
	if err != nil {
		return nil, WrapError(fmt.Errorf("failed to retrieve subdomain and domain during import: %v", err))
	}

	// Retrieve WRR status
	wrrStatus, err := getWRRStatus(client, subDomain, domainName)
	if err != nil {
		return nil, WrapError(fmt.Errorf("failed to retrieve WRR status during import: %v", err))
	}
	if err := d.Set("wrr_status", wrrStatus); err != nil {
		return nil, WrapError(fmt.Errorf("failed to set WRR status during import: %v", err))
	}

	// Retrieve weight
	weight, err := getWeight(client, recordID)
	if err != nil {
		return nil, WrapError(fmt.Errorf("failed to retrieve weight during import: %v", err))
	}
	if err := d.Set("weight", weight); err != nil {
		return nil, WrapError(fmt.Errorf("failed to set weight during import: %v", err))
	}

	return []*schema.ResourceData{d}, nil
}

func setWRRStatus(client *connectivity.AliyunClient, subDomain string, domainName string, enable bool) error {
	log.Printf("[DEBUG] Setting WRR Status for subdomain: %s in domain: %s to %t", subDomain, domainName, enable)

	// Validate if at least two records exist for the subdomain
	records, err := getAllSubdomainRecords(client, subDomain)
	if err != nil {
		return WrapError(err)
	}

	if len(records) < 2 {
		return fmt.Errorf("WRR cannot be enabled. Subdomain %s must have at least two records", subDomain)
	}

	// Create the request to enable/disable WRR
	request := alidns.CreateSetDNSSLBStatusRequest()
	request.SubDomain = subDomain
	request.Open = requests.NewBoolean(enable)

	return resource.Retry(3*time.Minute, func() *resource.RetryError {
		_, err := client.WithDnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
			return alidnsClient.SetDNSSLBStatus(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"LastOperationNotFinished"}) || NeedRetry(err) {
				return resource.RetryableError(fmt.Errorf("retrying SetDNSSLBStatus: %s", err))
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
}

func getWRRStatus(client *connectivity.AliyunClient, subDomain, domainName string) (bool, error) {
	// Log debug information
	log.Printf("[DEBUG] Checking WRR Status for subdomain: %s in domain: %s", subDomain, domainName)

	// Create the request to describe the WRR status
	request := alidns.CreateDescribeDNSSLBSubDomainsRequest()
	request.DomainName = domainName

	// Call the API
	response, err := client.WithDnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
		return alidnsClient.DescribeDNSSLBSubDomains(request)
	})
	if err != nil {
		return false, WrapError(fmt.Errorf("error describing WRR status: %v", err))
	}

	// Parse the response
	descResponse := response.(*alidns.DescribeDNSSLBSubDomainsResponse)
	for _, record := range descResponse.SlbSubDomains.SlbSubDomain {
		if record.SubDomain == subDomain {
			return record.Open, nil // Return the WRR status (true if enabled, false otherwise)
		}
	}

	return false, fmt.Errorf("WRR status not found for subdomain: %s", subDomain)
}

func setWeight(client *connectivity.AliyunClient, recordID string, weight int) error {
	request := alidns.CreateUpdateDNSSLBWeightRequest()
	request.RecordId = recordID
	request.Weight = requests.NewInteger(weight)

	return resource.Retry(3*time.Minute, func() *resource.RetryError {
		_, err := client.WithDnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
			return alidnsClient.UpdateDNSSLBWeight(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"LastOperationNotFinished"}) || NeedRetry(err) {
				return resource.RetryableError(fmt.Errorf("retrying UpdateDNSSLBWeight: %s", err))
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
}

func getWeight(client *connectivity.AliyunClient, recordID string) (int, error) {
	subDomain, _, err := getSubDomainAndDomain(client, recordID)
	if err != nil {
		return 0, err
	}

	request := alidns.CreateDescribeSubDomainRecordsRequest()
	request.SubDomain = subDomain

	response, err := client.WithDnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
		return alidnsClient.DescribeSubDomainRecords(request)
	})
	if err != nil {
		return 0, err
	}

	descResponse := response.(*alidns.DescribeSubDomainRecordsResponse)
	for _, record := range descResponse.DomainRecords.Record {
		if record.RecordId == recordID {
			return int(record.Weight), nil
		}
	}

	return 0, fmt.Errorf("record weight not found for record_id: %s", recordID)
}

func getSubDomainAndDomain(client *connectivity.AliyunClient, recordID string) (string, string, error) {
	request := alidns.CreateDescribeDomainRecordInfoRequest()
	request.RecordId = recordID

	response, err := client.WithDnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
		return alidnsClient.DescribeDomainRecordInfo(request)
	})
	if err != nil {
		return "", "", WrapError(fmt.Errorf("failed to retrieve domain info for record_id %s: %v", recordID, err))
	}

	descResponse := response.(*alidns.DescribeDomainRecordInfoResponse)
	subDomain := fmt.Sprintf("%s.%s", descResponse.RR, descResponse.DomainName)
	return subDomain, descResponse.DomainName, nil
}

func getAllSubdomainRecords(client *connectivity.AliyunClient, subDomain string) ([]alidns.Record, error) {
	request := alidns.CreateDescribeSubDomainRecordsRequest()
	request.SubDomain = subDomain
	request.PageSize = requests.NewInteger(100)
	request.PageNumber = requests.NewInteger(1)

	var records []alidns.Record

	for {
		response, err := client.WithDnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
			return alidnsClient.DescribeSubDomainRecords(request)
		})
		if err != nil {
			return nil, fmt.Errorf("failed to describe subdomain records: %s", err)
		}

		descResponse := response.(*alidns.DescribeSubDomainRecordsResponse)
		records = append(records, descResponse.DomainRecords.Record...)

		if len(records) >= int(descResponse.TotalCount) { // Convert TotalCount to int
			break
		}

		// Retrieve and increment the current page number
		currentPage, err := request.PageNumber.GetValue()
		if err != nil {
			return nil, fmt.Errorf("failed to get current page number: %s", err)
		}
		request.PageNumber = requests.NewInteger(currentPage + 1)
	}

	return records, nil
}
