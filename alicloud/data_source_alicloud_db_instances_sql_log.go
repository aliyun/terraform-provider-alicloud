package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func dataSourceAlicloudDBInstancesSQLLog() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudDBInstancesSQLLogRead,

		Schema: map[string]*schema.Schema{
			"db_instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"start_time": {
				Type:     schema.TypeString,
				Required: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Required: true,
			},
			"sql_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"query_key_words": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"data_base": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"user": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"page_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(30, 100),
				Default:      30,
			},
			"page_number": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// Computed values
			"total_record_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"page_record_number": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"page_record_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"sql_record": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"db_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"account_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sql_text": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"total_execution_times": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"return_row_counts": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"thread_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"execute_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudDBInstancesSQLLogRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := rds.CreateDescribeSQLLogRecordsRequest()

	request.RegionId = client.RegionId
	request.DBInstanceId = d.Get("db_instance_id").(string)
	request.StartTime = d.Get("start_time").(string)
	request.EndTime = d.Get("end_time").(string)
	request.Form = "Stream"
	if v, ok := d.GetOk("sql_id"); ok && v.(string) != "" {
		request.SQLId = requests.NewInteger(d.Get("sql_id").(int))
	}
	if v, ok := d.GetOk("query_key_words"); ok && v.(string) != "" {
		request.QueryKeywords = d.Get("query_key_words").(string)
	}
	if v, ok := d.GetOk("data_base"); ok && v.(string) != "" {
		request.Database = d.Get("data_base").(string)
	}
	if v, ok := d.GetOk("user"); ok && v.(string) != "" {
		request.User = d.Get("user").(string)
	}
	if v, ok := d.GetOk("page_size"); ok && v.(int) != 0 {
		request.PageSize = requests.NewInteger(d.Get("page_size").(int))
	}
	if v, ok := d.GetOk("page_number"); ok && v.(int) != 0 {
		request.PageNumber = requests.NewInteger(d.Get("page_number").(int))
	}
	raw, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
		return rdsClient.DescribeSQLLogRecords(request)
	})
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_db_instances_sql_log", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*rds.DescribeSQLLogRecordsResponse)
	if len(response.Items.SQLRecord) > 0 {
		var ids []string
		var sqlRecord []map[string]interface{}
		for _, item := range response.Items.SQLRecord {
			mapping := map[string]interface{}{
				"db_name":               item.DBName,
				"account_name":          item.AccountName,
				"host_address":          item.HostAddress,
				"sql_text":              item.SQLText,
				"total_execution_times": item.TotalExecutionTimes,
				"return_row_counts":     item.ReturnRowCounts,
				"thread_id":             item.ThreadID,
				"execute_time":          item.ExecuteTime,
			}
			ids = append(ids, item.SQLText)
			sqlRecord = append(sqlRecord, mapping)
		}
		d.SetId(dataResourceIdHash(ids))
		d.Set("total_record_count", response.TotalRecordCount)
		d.Set("page_record_number", response.PageNumber)
		d.Set("page_record_count", response.PageRecordCount)
		if err := d.Set("sql_record", sqlRecord); err != nil {
			return WrapError(err)
		}
		// create a json file in current directory and write data source to it
		if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
			writeToFile(output.(string), sqlRecord)
		}
		return nil
	}
	return nil
}
