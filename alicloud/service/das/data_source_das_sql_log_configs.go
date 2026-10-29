package das

import (
	"fmt"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceDasSqlLogConfigs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDasSqlLogConfigsRead,
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

func dataSourceDasSqlLogConfigsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dasService := NewService(client)

	instanceId := d.Get("instance_id").(string)

	objectRaw, err := dasService.DescribeDasSqlLogConfig(instanceId)
	if err != nil {
		if errs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	dataRawObj, _ := jsonpath.Get("$.Data", objectRaw)
	dataRaw, _ := dataRawObj.(map[string]interface{})
	if dataRaw == nil {
		dataRaw = map[string]interface{}{}
	}

	enableVal, err := convert.FormatBool(dataRaw["SqlLogEnable"])
	if err != nil {
		return err
	}
	request_enableVal, err := convert.FormatBool(dataRaw["RequestEnable"])
	if err != nil {
		return err
	}
	retentionVal, err := convert.FormatInt(dataRaw["Retention"])
	if err != nil {
		return err
	}
	hot_retentionVal, err := convert.FormatInt(dataRaw["HotRetention"])
	if err != nil {
		return err
	}
	cold_retentionVal, err := convert.FormatInt(dataRaw["ColdRetention"])
	if err != nil {
		return err
	}

	sql_log_visible_timeVal, err := convert.FormatInt(dataRaw["SqlLogVisibleTime"])
	if err != nil {
		return err
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

	d.SetId(helper.DataResourceIdHash([]string{instanceId}))
	if err := d.Set("instance_id", instanceId); err != nil {
		return errs.WrapError(err)
	}
	if err := d.Set("ids", []string{instanceId}); err != nil {
		return errs.WrapError(err)
	}
	if err := d.Set("configs", []map[string]interface{}{mapping}); err != nil {
		return errs.WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		if err := helper.WriteToFile(output.(string), []map[string]interface{}{mapping}); err != nil {
			return errs.WrapError(fmt.Errorf("writing output file: %w", err))
		}
	}

	return nil
}
