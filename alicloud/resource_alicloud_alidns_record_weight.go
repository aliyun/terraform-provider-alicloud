package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudAlidnsRecordWeight() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudAlidnsRecordWeightCreate,
		Read:   resourceAlicloudAlidnsRecordWeightRead,
		Update: resourceAlicloudAlidnsRecordWeightUpdate,
		Delete: resourceAlicloudAlidnsRecordWeightDelete,
		Importer: &schema.ResourceImporter{
			State: resourceAlicloudAlidnsRecordWeightImport,
		},
		Schema: map[string]*schema.Schema{
			"domain_name": {
				Type:     schema.TypeString,
				Required: true, // This makes domain_name mandatory
			},
			"line": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "default",
			},
			"priority": {
				Type:             schema.TypeInt,
				Optional:         true,
				DiffSuppressFunc: dnsPriorityDiffSuppressFunc,
			},
			"rr": {
				Type:     schema.TypeString,
				Required: true,
			},
			"remark": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"ENABLE", "DISABLE"}, false),
				Default:      "ENABLE",
			},
			"ttl": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  600,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"A", "NS", "MX", "TXT", "CNAME", "SRV", "AAAA", "CAA", "REDIRECT_URL", "FORWORD_URL"}, false),
			},
			"value": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: dnsValueDiffSuppressFunc,
			},
			"weight": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      0,
				ValidateFunc: validation.IntBetween(0, 100),
			},
			"wrr_status": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"ENABLE", "DISABLE"}, false),
				Default:      "ENABLE",
			},
		},
	}
}

func resourceAlicloudAlidnsRecordWeightCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	alidnsService := AlidnsService{client}

	domainName := d.Get("domain_name").(string)
	rr := d.Get("rr").(string)
	recordType := d.Get("type").(string)
	recordValue := d.Get("value").(string)
	ttl := d.Get("ttl").(int)

	// Step 1: Create the record
	request := alidns.CreateAddDomainRecordRequest()
	request.RegionId = client.RegionId
	request.DomainName = domainName
	request.RR = rr
	request.Type = recordType
	request.Value = recordValue
	request.TTL = requests.NewInteger(ttl)

	response, err := client.WithAlidnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
		return alidnsClient.AddDomainRecord(request)
	})
	if err != nil {
		return WrapError(err)
	}

	resp, ok := response.(*alidns.AddDomainRecordResponse)
	if !ok {
		return fmt.Errorf("unexpected response type for AddDomainRecord")
	}

	d.SetId(resp.RecordId)

	err = alidnsService.WaitForRecordReady(resp.RecordId, domainName, rr, 2*time.Minute)
	if err != nil {
		return WrapError(err)
	}

	recordID := d.Id()

	err = alidnsService.WaitForLastOperation(recordID, domainName, 2*time.Minute)
	if err != nil {
		return WrapError(err)
	}

	// Step 2: Update remark
	if d.HasChange("remark") {
		newRemark := d.Get("remark").(string)
		log.Printf("[DEBUG] Updating remark: RecordId=%s, NewRemark=%s", recordID, newRemark)

		err := alidnsService.UpdateRecordRemark(recordID, newRemark)
		if err != nil {
			return WrapError(err)
		}
	}

	// Step 3: Set WRR status
	if wrrStatus, ok := d.GetOk("wrr_status"); ok && wrrStatus.(string) == "ENABLE" {
		err = alidnsService.SetWRRStatus(domainName, rr, "ENABLE")
		if err != nil {
			return WrapError(err)
		}

		// Wait for WRR status to be set
		err = alidnsService.WaitForLastOperation(resp.RecordId, domainName, 2*time.Minute)
		if err != nil {
			return WrapError(err)
		}
	}

	// Step 4: Set weight
	if weight, ok := d.GetOk("weight"); ok && weight.(int) > 0 {
		err = alidnsService.SetRecordWeight(resp.RecordId, weight.(int), domainName, rr)
		if err != nil {
			return WrapError(err)
		}
	}

	return resourceAlicloudAlidnsRecordWeightRead(d, meta)
}

func resourceAlicloudAlidnsRecordWeightUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	alidnsService := AlidnsService{client}
	recordID := d.Id()
	domainName := d.Get("domain_name").(string)
	rr := d.Get("rr").(string)

	// Step 1: Update remark
	if d.HasChange("remark") {
		newRemark := d.Get("remark").(string)
		log.Printf("[DEBUG] Updating remark: RecordId=%s, NewRemark=%s", recordID, newRemark)

		err := alidnsService.UpdateRecordRemark(recordID, newRemark)
		if err != nil {
			return WrapError(err)
		}
	}

	// Step 2: Update status
	if d.HasChange("status") {
		newStatus := d.Get("status").(string)
		log.Printf("[DEBUG] Updating status: RecordId=%s, NewStatus=%s", recordID, newStatus)

		err := alidnsService.SetRecordStatus(recordID, newStatus)
		if err != nil {
			return WrapError(err)
		}
	}

	// Step 3: Update other record attributes (ttl, value, type)
	if d.HasChange("ttl") || d.HasChange("value") || d.HasChange("type") {
		updateRequest := alidns.CreateUpdateDomainRecordRequest()
		updateRequest.RecordId = recordID
		updateRequest.RR = rr
		updateRequest.Type = d.Get("type").(string)
		updateRequest.Value = d.Get("value").(string)
		updateRequest.TTL = requests.NewInteger(d.Get("ttl").(int))

		_, err := client.WithAlidnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
			return alidnsClient.UpdateDomainRecord(updateRequest)
		})
		if err != nil {
			return WrapError(err)
		}
	}

	// Step 4: Update WRR status
	if d.HasChange("wrr_status") {
		newWRRStatus := d.Get("wrr_status").(string)
		log.Printf("[DEBUG] Updating WRR status: DomainName=%s, SubDomain=%s, NewStatus=%s", domainName, rr, newWRRStatus)

		err := alidnsService.SetWRRStatus(domainName, rr, newWRRStatus)
		if err != nil {
			return WrapError(err)
		}

		err = alidnsService.WaitForLastOperation(recordID, domainName, 2*time.Minute)
		if err != nil {
			return WrapError(err)
		}
	}

	// Step 5: Update weight
	if d.HasChange("weight") {
		newWeight := d.Get("weight").(int)
		log.Printf("[DEBUG] Updating weight: RecordId=%s, NewWeight=%d", recordID, newWeight)

		err := alidnsService.SetRecordWeight(recordID, newWeight, domainName, rr)
		if err != nil {
			return WrapError(err)
		}
	}

	return resourceAlicloudAlidnsRecordWeightRead(d, meta)
}

func resourceAlicloudAlidnsRecordWeightRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	alidnsService := AlidnsService{client}

	recordID := d.Id()
	domainName := d.Get("domain_name").(string)
	rr := d.Get("rr").(string)

	// Retry logic for API response delay
	var record *alidns.Record
	var err error
	for i := 0; i < 5; i++ { // Retry up to 5 times
		record, err = alidnsService.DescribeDomainRecordById(recordID, domainName)
		if err == nil && record != nil {
			break
		}
		log.Printf("[DEBUG] Retrying to fetch record: RecordId=%s, Attempt=%d", recordID, i+1)
		time.Sleep(2 * time.Second) // Wait 2 seconds before retrying
	}

	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidRecordId.NotFound"}) {
			log.Printf("[DEBUG] Record not found: %s", recordID)
			d.SetId("") // Mark resource as deleted
			return nil
		}
		return WrapError(err)
	}

	// Fetch WRR status
	wrrStatus, err := alidnsService.GetWRRStatus(domainName, rr)
	if err != nil {
		log.Printf("[WARN] Failed to fetch WRR status: %s", err)
		wrrStatus = "DISABLE" // Default to DISABLE if fetch fails
	}
	d.Set("wrr_status", wrrStatus)

	// Update Terraform state
	d.Set("status", record.Status)
	d.Set("remark", record.Remark)
	d.Set("weight", record.Weight)
	d.Set("rr", record.RR)
	d.Set("type", record.Type)
	d.Set("value", record.Value)
	d.Set("ttl", record.TTL)
	d.Set("line", record.Line)
	d.Set("priority", record.Priority)

	return nil
}

func resourceAlicloudAlidnsRecordWeightDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	alidnsService := AlidnsService{client}
	recordID := d.Id()
	domainName := d.Get("domain_name").(string)

	// Wait for the last operation to finish
	err := alidnsService.WaitForLastOperation(recordID, domainName, 2*time.Minute)
	if err != nil {
		return WrapError(err)
	}

	// Proceed with the delete operation
	request := alidns.CreateDeleteDomainRecordRequest()
	request.RegionId = client.RegionId
	request.RecordId = recordID

	_, err = client.WithAlidnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
		return alidnsClient.DeleteDomainRecord(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"DomainRecordNotBelongToUser", "InvalidRecordId.NotFound"}) {
			return nil
		}
		return WrapError(err)
	}

	d.SetId("")
	return nil
}

func resourceAlicloudAlidnsRecordWeightImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(*connectivity.AliyunClient)
	alidnsService := AlidnsService{client}

	importID := d.Id()
	var domainName, recordID string

	// Parse the import ID to support both formats
	parts := strings.Split(importID, "/")
	if len(parts) == 2 {
		domainName = parts[0]
		recordID = parts[1]
	} else if len(parts) == 1 {
		recordID = parts[0]

		// If domainName is not provided, attempt to discover it
		domains, err := alidnsService.ListAllDomains()
		if err != nil {
			return nil, fmt.Errorf("failed to list all domains: %w", err)
		}

		for _, domain := range domains {
			records, err := alidnsService.DescribeDomainRecords(domain.DomainName)
			if err != nil {
				return nil, fmt.Errorf("failed to describe records for domain %s: %w", domain.DomainName, err)
			}

			for _, record := range records {
				if record.RecordId == recordID {
					domainName = domain.DomainName
					break
				}
			}

			if domainName != "" {
				break
			}
		}

		if domainName == "" {
			return nil, fmt.Errorf("record with ID %s not found in any domain", recordID)
		}
	} else {
		return nil, fmt.Errorf("invalid import ID format: %s", importID)
	}

	// Fetch the record details
	record, err := alidnsService.DescribeDomainRecordById(recordID, domainName)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch record %s in domain %s: %w", recordID, domainName, err)
	}

	// Fetch WRR configuration
	weight := 0
	wrrStatus := "DISABLE"

	slbRecords, err := alidnsService.GetSLBSubDomains(domainName)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch WRR configuration for domain %s: %w", domainName, err)
	}

	for _, subRecord := range slbRecords {
		if subRecord.SubDomain == fmt.Sprintf("%s.%s", record.RR, domainName) {
			wrrStatus = "ENABLE"
			weight, err = alidnsService.GetWeight(recordID, domainName)
			if err != nil {
				return nil, fmt.Errorf("failed to fetch weight for record ID %s: %w", recordID, err)
			}
			break
		}
	}

	// Set Terraform resource data
	d.SetId(recordID)
	d.Set("domain_name", domainName)
	d.Set("rr", record.RR)
	d.Set("type", record.Type)
	d.Set("value", record.Value)
	d.Set("ttl", record.TTL)
	d.Set("line", record.Line)
	d.Set("remark", record.Remark)
	d.Set("weight", weight)
	d.Set("wrr_status", wrrStatus)

	return []*schema.ResourceData{d}, nil
}
