package alicloud

import (
	"fmt"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudThreatDetectionDataSet() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudThreatDetectionDataSetRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Required: true,
				Type:     schema.TypeString,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"data_set_id": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"data_set_name": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"data_set_field_key_name": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"data_set_file_name": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"data_set_description": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"data_set_type": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"data_set_status": {
				Computed: true,
				Type:     schema.TypeInt,
			},
			"role_for": {
				Computed: true,
				Type:     schema.TypeInt,
			},
			"lang": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"region_id": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"create_time": {
				Computed: true,
				Type:     schema.TypeInt,
			},
			"ip_whitelist_recognizers": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auto_recognize_status": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"recognize_scope": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"ip_whitelist_recognizer_type": {
							Computed: true,
							Type:     schema.TypeString,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudThreatDetectionDataSetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	threatDetectionServiceV2 := ThreatDetectionServiceV2{client}

	objectRaw, err := threatDetectionServiceV2.DescribeThreatDetectionDataSet(d.Get("id").(string))
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_threat_detection_data_set", "ListDataSets", AlibabaCloudSdkGoERROR)
	}

	dataSetId := fmt.Sprint(objectRaw["DataSetId"])
	d.SetId(dataSetId)
	d.Set("data_set_id", dataSetId)
	d.Set("data_set_name", objectRaw["DataSetName"])
	d.Set("data_set_field_key_name", objectRaw["DataSetFieldKeyName"])
	d.Set("data_set_file_name", objectRaw["DataSetFileName"])
	d.Set("data_set_description", objectRaw["DataSetDescription"])
	d.Set("data_set_type", objectRaw["DataSetType"])
	d.Set("data_set_status", objectRaw["DataSetStatus"])
	d.Set("role_for", objectRaw["RoleFor"])
	d.Set("lang", objectRaw["Lang"])
	d.Set("region_id", objectRaw["RegionId"])
	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("ip_whitelist_recognizers", flattenThreatDetectionDataSetIpWhitelistRecognizersResponse(objectRaw["IpWhitelistRecognizers"]))

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), []map[string]interface{}{objectRaw})
	}

	return nil
}
