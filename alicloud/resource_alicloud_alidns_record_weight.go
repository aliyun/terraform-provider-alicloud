package alicloud

import (
	"fmt"
	"log"
	"strings"

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
				Required: true,
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

	// Step 1: Check for existing records
	domainName := d.Get("domain_name").(string)
	rr := d.Get("rr").(string)
	recordType := d.Get("type").(string)
	recordValue := d.Get("value").(string)

	records, err := alidnsService.DescribeDomainRecords(domainName)
	if err != nil {
		return WrapError(err)
	}

	// Look for a matching record
	var existingRecord *alidns.Record
	for _, record := range records {
		if record.RR == rr && record.Type == recordType && record.Value == recordValue {
			existingRecord = &record
			break
		}
	}

	// Step 2: If record exists, update weight and WRR status
	if existingRecord != nil {
		log.Printf("[DEBUG] Found existing record: %s", existingRecord.RecordId)

		// Enable WRR for the subdomain if needed
		wrrStatus := d.Get("wrr_status").(string)
		if wrrStatus == "ENABLE" {
			err = alidnsService.SetWRRStatus(domainName, rr, wrrStatus)
			if err != nil {
				return WrapError(err)
			}

			// Update the weight for the existing record
			if weight, ok := d.GetOk("weight"); ok && weight.(int) > 0 {
				err = alidnsService.SetRecordWeight(existingRecord.RecordId, weight.(int))
				if err != nil {
					return WrapError(err)
				}
			}
		}

		// Set the existing record's ID
		d.SetId(existingRecord.RecordId)
		return resourceAlicloudAlidnsRecordWeightRead(d, meta)
	}

	// Step 3: If no existing record, create a new one
	log.Printf("[DEBUG] No existing record found. Creating a new record for %s", rr)

	request := alidns.CreateAddDomainRecordRequest()
	request.RegionId = client.RegionId
	request.DomainName = domainName
	request.RR = rr
	request.Type = recordType
	request.Value = recordValue
	request.TTL = requests.NewInteger(d.Get("ttl").(int))
	request.Line = d.Get("line").(string)

	response, err := client.WithAlidnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
		return alidnsClient.AddDomainRecord(request)
	})
	if err != nil {
		return WrapError(err)
	}

	resp, _ := response.(*alidns.AddDomainRecordResponse)
	d.SetId(resp.RecordId)

	// Enable WRR and set weight for the new record
	wrrStatus := d.Get("wrr_status").(string)
	if wrrStatus == "ENABLE" {
		err = alidnsService.SetWRRStatus(domainName, rr, wrrStatus)
		if err != nil {
			return WrapError(err)
		}

		if weight, ok := d.GetOk("weight"); ok && weight.(int) > 0 {
			err = alidnsService.SetRecordWeight(resp.RecordId, weight.(int))
			if err != nil {
				return WrapError(err)
			}
		}
	}

	return resourceAlicloudAlidnsRecordWeightRead(d, meta)
}

func resourceAlicloudAlidnsRecordWeightUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	alidnsService := AlidnsService{client}

	// Step 1: Update WRR status if changed
	if d.HasChange("wrr_status") {
		wrrStatus := d.Get("wrr_status").(string)
		err := alidnsService.SetWRRStatus(d.Get("domain_name").(string), d.Get("rr").(string), wrrStatus)
		if err != nil {
			return WrapError(err)
		}
	}

	// Step 2: Update weight only if WRR is enabled
	if d.HasChange("weight") && d.Get("wrr_status").(string) == "ENABLE" {
		err := alidnsService.SetRecordWeight(d.Id(), d.Get("weight").(int))
		if err != nil {
			return WrapError(err)
		}
	} else if d.Get("wrr_status").(string) == "DISABLE" {
		log.Printf("[DEBUG] WRR is disabled, skipping weight update for record ID: %s", d.Id())
	}

	return resourceAlicloudAlidnsRecordWeightRead(d, meta)
}

func resourceAlicloudAlidnsRecordWeightRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	alidnsService := AlidnsService{client}

	recordID := d.Id()

	// Use DescribeDomainRecordInfo to fetch detailed record info
	record, err := alidnsService.DescribeDomainRecordById(recordID)
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidRecordId.NotFound"}) {
			log.Printf("[DEBUG] Record not found: %s", recordID)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	// Populate resource data
	d.Set("domain_name", record.DomainName)
	d.Set("rr", record.RR)
	d.Set("type", record.Type)
	d.Set("value", record.Value)
	d.Set("ttl", record.TTL)
	d.Set("line", record.Line)
	d.Set("weight", record.Weight)
	d.Set("priority", record.Priority) // Map Priority
	d.Set("remark", record.Remark)     // Map Remark

	// Determine WRR status based on weight
	if record.Weight > 0 {
		d.Set("wrr_status", "ENABLE")
	} else {
		d.Set("wrr_status", "DISABLE")
	}

	return nil
}

func resourceAlicloudAlidnsRecordWeightDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := alidns.CreateDeleteDomainRecordRequest()
	request.RegionId = client.RegionId
	request.RecordId = d.Id()

	_, err := client.WithAlidnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
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
	parts := strings.Split(importID, "/")
	var domainName, recordID string
	if len(parts) == 2 {
		domainName = parts[0]
		recordID = parts[1]
	} else {
		recordID = importID
		record, err := alidnsService.DescribeDomainRecordById(recordID)
		if err != nil {
			return nil, WrapError(err)
		}
		domainName = record.DomainName
	}

	records, err := alidnsService.DescribeDomainRecords(domainName)
	if err != nil {
		return nil, WrapError(err)
	}

	var importedRecord *alidns.Record
	for _, record := range records {
		if record.RecordId == recordID {
			importedRecord = &record
			break
		}
	}
	if importedRecord == nil {
		return nil, fmt.Errorf("record with ID %s not found in domain %s", recordID, domainName)
	}

	d.SetId(importedRecord.RecordId)
	d.Set("domain_name", domainName)
	d.Set("rr", importedRecord.RR)
	d.Set("type", importedRecord.Type)
	d.Set("value", importedRecord.Value)
	d.Set("ttl", importedRecord.TTL)
	d.Set("line", importedRecord.Line)
	d.Set("weight", importedRecord.Weight)
	d.Set("priority", 0) // Default or API-provided value
	d.Set("remark", "")  // Default or API-provided value

	// Fetch WRR status
	if importedRecord.Weight > 0 {
		d.Set("wrr_status", "ENABLE")
	} else {
		d.Set("wrr_status", "DISABLE")
	}

	return []*schema.ResourceData{d}, nil
}

func (s *AlidnsService) SetWRRStatus(domainName, rr, status string) error {
	request := alidns.CreateSetDNSSLBStatusRequest()
	request.RegionId = s.client.RegionId
	request.DomainName = domainName
	request.SubDomain = fmt.Sprintf("%s.%s", rr, domainName)
	request.Open = requests.NewBoolean(status == "ENABLE")

	_, err := s.client.WithAlidnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
		return alidnsClient.SetDNSSLBStatus(request)
	})
	return WrapError(err)
}

func (s *AlidnsService) SetRecordWeight(recordID string, weight int) error {
	request := alidns.CreateUpdateDNSSLBWeightRequest()
	request.RegionId = s.client.RegionId
	request.RecordId = recordID
	request.Weight = requests.NewInteger(weight)

	_, err := s.client.WithAlidnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
		return alidnsClient.UpdateDNSSLBWeight(request)
	})
	return WrapError(err)
}

func (s *AlidnsService) DescribeDomainRecordById(recordID string) (*alidns.Record, error) {
	request := alidns.CreateDescribeDomainRecordInfoRequest()
	request.RecordId = recordID

	response, err := s.client.WithAlidnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
		return alidnsClient.DescribeDomainRecordInfo(request)
	})
	if err != nil {
		return nil, WrapError(err)
	}

	resp, ok := response.(*alidns.DescribeDomainRecordInfoResponse)
	if !ok {
		return nil, fmt.Errorf("failed to cast response to DescribeDomainRecordInfoResponse")
	}

	return &alidns.Record{
		RecordId:   resp.RecordId,
		DomainName: resp.DomainName,
		RR:         resp.RR,
		Type:       resp.Type,
		Value:      resp.Value,
		TTL:        resp.TTL,
		Line:       resp.Line,
		Priority:   resp.Priority, // Ensure this is set
		Remark:     resp.Remark,   // Ensure this is set
	}, nil
}
