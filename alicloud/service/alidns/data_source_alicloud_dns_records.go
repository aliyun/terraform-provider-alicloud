package alidns

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	alidnssdk "github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlicloudDnsRecords() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudDnsRecordsRead,

		Schema: map[string]*schema.Schema{
			"domain_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"host_record_regex": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"value_regex": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				// must be one of [A, NS, MX, TXT, CNAME, SRV, AAAA, CAA, REDIRECT_URL, FORWORD_URL]
				ValidateFunc: validation.StringInSlice([]string{"A", "NS", "MX", "TXT", "CNAME", "SRV", "AAAA", "CAA", "REDIRECT_URL", "FORWORD_URL"}, false),
			},
			"line": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"enable", "disable"}, true),
			},
			"is_locked": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"urls": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			// Computed values
			"records": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"record_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domain_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host_record": {
							Type:     schema.TypeString,
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
						"ttl": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"priority": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"line": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"locked": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudDnsRecordsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := alidnssdk.CreateDescribeDomainRecordsRequest()
	request.RegionId = client.RegionId
	request.DomainName = d.Get("domain_name").(string)
	if v, ok := d.GetOk("type"); ok && v.(string) != "" {
		request.TypeKeyWord = v.(string)
	}

	var allRecords []alidnssdk.Record

	request.PageSize = requests.NewInteger(helper.PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithDnsClient(func(dnsClient *alidnssdk.Client) (interface{}, error) {
			return dnsClient.DescribeDomainRecords(request)
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_dns_records", request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*alidnssdk.DescribeDomainRecordsResponse)
		records := response.DomainRecords.Record
		for _, record := range records {
			allRecords = append(allRecords, record)
		}

		if len(records) < helper.PageSizeLarge {
			break
		}
		page, err := helper.GetNextpageNumber(request.PageNumber)
		if err != nil {
			return errs.WrapError(err)
		}
		request.PageNumber = page
	}

	var filteredRecords []alidnssdk.Record
	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}
	for _, record := range allRecords {
		if v, ok := d.GetOk("line"); ok && v.(string) != "" && strings.ToUpper(record.Line) != strings.ToUpper(v.(string)) {
			continue
		}

		if v, ok := d.GetOk("status"); ok && v.(string) != "" && strings.ToUpper(record.Status) != strings.ToUpper(v.(string)) {
			continue
		}

		if v, ok := d.GetOk("is_locked"); ok && record.Locked != v.(bool) {
			continue
		}

		if v, ok := d.GetOk("host_record_regex"); ok && v.(string) != "" {
			r := regexp.MustCompile(v.(string))
			if !r.MatchString(record.RR) {
				continue
			}
		}

		if v, ok := d.GetOk("value_regex"); ok && v.(string) != "" {
			r := regexp.MustCompile(v.(string))
			if !r.MatchString(record.Value) {
				continue
			}
		}
		if len(idsMap) > 0 {
			if _, ok := idsMap[record.RecordId]; !ok {
				continue
			}
		}

		filteredRecords = append(filteredRecords, record)
	}

	return recordsDecriptionAttributes(d, filteredRecords, meta)
}

func recordsDecriptionAttributes(d *schema.ResourceData, recordTypes []alidnssdk.Record, meta interface{}) error {
	var ids []string
	var urls []string
	var s []map[string]interface{}
	for _, record := range recordTypes {
		mapping := map[string]interface{}{
			"record_id":   record.RecordId,
			"domain_name": record.DomainName,
			"line":        record.Line,
			"host_record": record.RR,
			"type":        record.Type,
			"value":       record.Value,
			"status":      strings.ToLower(record.Status),
			"locked":      record.Locked,
			"ttl":         record.TTL,
			"priority":    record.Priority,
		}
		ids = append(ids, record.RecordId)
		urls = append(urls, fmt.Sprintf("%v.%v", record.RR, record.DomainName))
		s = append(s, mapping)
	}

	d.SetId(helper.DataResourceIdHash(ids))
	if err := d.Set("records", s); err != nil {
		return errs.WrapError(err)
	}
	if err := d.Set("urls", urls); err != nil {
		return errs.WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return errs.WrapError(err)
	}
	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		helper.WriteToFile(output.(string), s)
	}
	return nil
}
