package alicloud

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudAlidnsRecordsWeight() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudAlidnsRecordsWeightRead,
		Schema: map[string]*schema.Schema{
			"record_ids": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"records": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"record_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"weight": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"wrr_status": {
							Type:     schema.TypeBool,
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
	recordIDs := d.Get("record_ids").([]interface{})

	var records []map[string]interface{}
	var idList []string

	for _, rid := range recordIDs {
		recordID := rid.(string)
		subDomain, domainName, err := getSubDomainAndDomain(client, recordID)
		if err != nil {
			return WrapErrorf(err, "failed to get subdomain/domain for record %s", recordID)
		}

		weight, err := getWeight(client, recordID)
		if err != nil {
			return WrapErrorf(err, "failed to get weight for record %s", recordID)
		}

		wrrStatus, err := getWRRStatus(client, subDomain, domainName)
		if err != nil {
			return WrapErrorf(err, "failed to get WRR status for subdomain %s", subDomain)
		}

		records = append(records, map[string]interface{}{
			"record_id":  recordID,
			"weight":     weight,
			"wrr_status": wrrStatus,
		})
		idList = append(idList, recordID)
	}

	d.SetId(dataResourceIdHash(idList))
	if err := d.Set("records", records); err != nil {
		return WrapError(err)
	}
	return nil
}
