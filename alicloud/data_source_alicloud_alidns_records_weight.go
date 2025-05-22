package alicloud

import (
	"log"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudAlidnsRecordsWeight() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudAlidnsRecordsWeightRead,
		Schema: map[string]*schema.Schema{
			"domain_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"line": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "default",
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"ENABLE", "DISABLE"}, false),
				Default:      "ENABLE",
			},
			"weight_greater_than": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Default:  0,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"records": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"line": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rr": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"record_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ttl": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"weight": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"wrr_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"remark": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"priority": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudAlidnsRecordsWeightRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := alidns.CreateDescribeDomainRecordsRequest()
	request.DomainName = d.Get("domain_name").(string)
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)

	var records []map[string]interface{}
	var recordIDs []string // To store unique IDs for hash generation

	for {
		raw, err := client.WithAlidnsClient(func(client *alidns.Client) (interface{}, error) {
			return client.DescribeDomainRecords(request)
		})
		if err != nil {
			return WrapError(err)
		}

		response := raw.(*alidns.DescribeDomainRecordsResponse)
		for _, record := range response.DomainRecords.Record {
			wrrStatus, err := fetchWRRStatus(client, record.DomainName, record.RR)
			if err != nil {
				log.Printf("[DEBUG] Failed to fetch WRR status for record %s: %v", record.RecordId, err)
				wrrStatus = "DISABLE"
			}

			detailedRecord, err := client.WithAlidnsClient(func(client *alidns.Client) (interface{}, error) {
				request := alidns.CreateDescribeDomainRecordInfoRequest()
				request.RecordId = record.RecordId
				return client.DescribeDomainRecordInfo(request)
			})

			remark := ""
			if err == nil {
				remark = detailedRecord.(*alidns.DescribeDomainRecordInfoResponse).Remark
			} else {
				log.Printf("[DEBUG] Failed to fetch remark for record %s: %v", record.RecordId, err)
			}

			records = append(records, map[string]interface{}{
				"domain_name": record.DomainName,
				"line":        record.Line,
				"rr":          record.RR,
				"record_id":   record.RecordId,
				"status":      record.Status,
				"ttl":         record.TTL,
				"type":        record.Type,
				"value":       record.Value,
				"weight":      record.Weight,
				"remark":      remark,
				"wrr_status":  wrrStatus,
			})

			recordIDs = append(recordIDs, record.RecordId) // Add the record ID to the list
		}

		if len(response.DomainRecords.Record) < PageSizeLarge {
			break
		}

		page, err := getNextpageNumber(request.PageNumber)
		if err != nil {
			return WrapError(err)
		}
		request.PageNumber = page
	}

	if err := d.Set("records", records); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), records)
	}

	// Use the list of record IDs to generate a unique hash
	d.SetId(dataResourceIdHash(recordIDs))
	return nil
}
