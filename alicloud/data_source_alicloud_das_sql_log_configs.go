package alicloud

import (
	"fmt"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAliCloudDasSqlLogConfigs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudDasSqlLogConfigsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"configs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"request_enable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"retention": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"hot_retention": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"cold_retention": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"log_filter": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sql_log_visible_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAliCloudDasSqlLogConfigsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dasServiceV2 := DasServiceV2{client}

	instanceId := d.Get("instance_id").(string)

	objectRaw, err := dasServiceV2.DescribeDasSqlLogConfig(instanceId)
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	dataRawObj, _ := jsonpath.Get("$.Data", objectRaw)
	dataRaw, _ := dataRawObj.(map[string]interface{})
	if dataRaw == nil {
		dataRaw = map[string]interface{}{}
	}

	enableVal, err100 := formatBool(dataRaw["SqlLogEnable"])
	if err100 != nil {
		return err100
	}
	request_enableVal, err101 := formatBool(dataRaw["RequestEnable"])
	if err101 != nil {
		return err101
	}
	retentionVal, err102 := formatInt(dataRaw["Retention"])
	if err102 != nil {
		return err102
	}
	hot_retentionVal, err103 := formatInt(dataRaw["HotRetention"])
	if err103 != nil {
		return err103
	}
	cold_retentionVal, err104 := formatInt(dataRaw["ColdRetention"])
	if err104 != nil {
		return err104
	}

	sql_log_visible_timeVal, err107 := formatInt(dataRaw["SqlLogVisibleTime"])
	if err107 != nil {
		return err107
	}
	mapping := map[string]interface{}{
		"instance_id":          instanceId,
		"enable":               enableVal,
		"request_enable":       request_enableVal,
		"retention":            retentionVal,
		"hot_retention":        hot_retentionVal,
		"cold_retention":       cold_retentionVal,
		"version":              dataRaw["Version"],
		"log_filter":           dataRaw["LogFilter"],
		"sql_log_visible_time": sql_log_visible_timeVal,
	}

	d.SetId(dataResourceIdHash([]string{instanceId}))
	if err := d.Set("instance_id", instanceId); err != nil {
		return WrapError(err)
	}
	if err := d.Set("ids", []string{instanceId}); err != nil {
		return WrapError(err)
	}
	if err := d.Set("configs", []map[string]interface{}{mapping}); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		if err := writeToFile(output.(string), []map[string]interface{}{mapping}); err != nil {
			return WrapError(fmt.Errorf("writing output file: %w", err))
		}
	}

	return nil
}
