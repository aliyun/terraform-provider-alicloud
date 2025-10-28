package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudDtsInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudDtsInstanceCreate,
		Read:   resourceAliCloudDtsInstanceRead,
		Update: resourceAliCloudDtsInstanceUpdate,
		Delete: resourceAliCloudDtsInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"auto_pay": {
				Optional: true,
				Type:     schema.TypeBool,
			},
			"auto_start": {
				Optional: true,
				Type:     schema.TypeBool,
			},
			"compute_unit": {
				Optional: true,
				Type:     schema.TypeInt,
			},
			"create_time": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"database_count": {
				Optional: true,
				Computed: true,
				Type:     schema.TypeInt,
			},
			"destination_endpoint_engine_name": {
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				Type:         schema.TypeString,
				AtLeastOneOf: []string{"destination_endpoint_engine_name", "job_id"},
				ValidateFunc: StringInSlice([]string{"MySQL", "PolarDB", "polardb_o", "polardb_pg", "Redis", "DRDS", "PostgreSQL", "odps", "oracle", "mongodb", "tidb", "ADS", "ADB30", "Greenplum", "MSSQL", "kafka", "DataHub", "clickhouse", "DB2", "as400", "Tablestore"}, false),
			},
			"dts_instance_id": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"du": {
				Optional:     true,
				Type:         schema.TypeInt,
				ValidateFunc: IntBetween(1, 100),
			},
			"fee_type": {
				Optional:     true,
				Type:         schema.TypeString,
				ValidateFunc: StringInSlice([]string{"ONLY_CONFIGURATION_FEE", "CONFIGURATION_FEE_AND_DATA_FEE"}, false),
			},
			"instance_class": {
				Optional:     true,
				ForceNew:     true,
				Type:         schema.TypeString,
				ValidateFunc: StringInSlice([]string{"4xlarge", "2xlarge", "xlarge", "large", "medium", "small", "micro"}, false),
			},
			"instance_name": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"job_id": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"payment_type": {
				Optional:     true,
				ForceNew:     true,
				Type:         schema.TypeString,
				ValidateFunc: StringInSlice([]string{"PayAsYouGo", "Subscription"}, false),
			},
			"period": {
				Optional:     true,
				Type:         schema.TypeString,
				ValidateFunc: StringInSlice([]string{"Month", "Year"}, false),
			},
			"resource_group_id": {
				Optional: true,
				Computed: true,
				Type:     schema.TypeString,
			},
			"source_endpoint_engine_name": {
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				Type:         schema.TypeString,
				AtLeastOneOf: []string{"source_endpoint_engine_name", "job_id"},
				ValidateFunc: StringInSlice([]string{"MySQL", "PolarDB", "polardb_o", "polardb_pg", "Redis", "DRDS", "PostgreSQL", "odps", "oracle", "mongodb", "tidb", "ADS", "ADB30", "Greenplum", "MSSQL", "kafka", "DataHub", "clickhouse", "DB2", "as400", "Tablestore"}, false),
			},
			"source_region": {
				Optional:     true,
				ForceNew:     true,
				Type:         schema.TypeString,
				AtLeastOneOf: []string{"source_region", "job_id"},
			},
			"destination_region": {
				Optional:     true,
				ForceNew:     true,
				Type:         schema.TypeString,
				AtLeastOneOf: []string{"destination_region", "job_id"},
			},
			"status": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"sync_architecture": {
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"bidirectional", "oneway"}, false),
				Type:         schema.TypeString,
			},
			"synchronization_direction": {
				Optional:     true,
				Type:         schema.TypeString,
				ValidateFunc: StringInSlice([]string{"Forward", "Reverse"}, false),
			},
			"tags": tagsSchema(),
			"type": {
				Optional:     true,
				ForceNew:     true,
				Type:         schema.TypeString,
				AtLeastOneOf: []string{"type", "job_id"},
				ValidateFunc: StringInSlice([]string{"migration", "sync", "subscribe"}, false),
			},
			"used_time": {
				Optional:     true,
				Type:         schema.TypeInt,
				ValidateFunc: IntInSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("payment_type").(string) != "Subscription"
				},
			},
		},
	}
}

func resourceAliCloudDtsInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := map[string]interface{}{
		"RegionId": client.RegionId,
	}
	var err error

	if v, ok := d.GetOk("auto_start"); ok {
		request["AutoStart"] = v
	}
	if v, ok := d.GetOk("du"); ok {
		request["Du"] = v
	}
	if v, ok := d.GetOk("auto_pay"); ok {
		request["AutoPay"] = v
	}
	if v, ok := d.GetOk("compute_unit"); ok {
		request["ComputeUnit"] = v
	}
	if v, ok := d.GetOk("database_count"); ok {
		request["DatabaseCount"] = v
	}
	if v, ok := d.GetOk("destination_endpoint_engine_name"); ok {
		request["DestinationEndpointEngineName"] = v
	}
	if v, ok := d.GetOk("fee_type"); ok {
		request["FeeType"] = v
	}
	if v, ok := d.GetOk("instance_class"); ok {
		request["InstanceClass"] = v
	}
	if v, ok := d.GetOk("job_id"); ok {
		request["JobId"] = v
	}
	if v, ok := d.GetOk("payment_type"); ok {
		request["PayType"] = convertDTSInstancePaymentTypeRequest(v.(string))
	}
	if v, ok := d.GetOk("period"); ok {
		request["Period"] = v
	}
	if v, ok := d.GetOk("source_endpoint_engine_name"); ok {
		request["SourceEndpointEngineName"] = v
	}
	if v, ok := d.GetOk("source_region"); ok {
		request["SourceRegion"] = v
	}
	if v, ok := d.GetOk("destination_region"); ok {
		request["DestinationRegion"] = v
	}
	if v, ok := d.GetOk("sync_architecture"); ok {
		request["SyncArchitecture"] = v
	}
	if v, ok := d.GetOk("type"); ok {
		request["Type"] = v
	}
	if v, ok := d.GetOk("used_time"); ok {
		request["UsedTime"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}

	var response map[string]interface{}
	action := "CreateDtsInstance"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("Dts", "2020-01-01", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_dts_instance", action, AlibabaCloudSdkGoERROR)
	}

	if v, err := jsonpath.Get("$.InstanceId", response); err != nil || v == nil {
		return WrapErrorf(err, IdMsg, "alicloud_dts_instance")
	} else {
		d.SetId(fmt.Sprint(v))
	}

	return resourceAliCloudDtsInstanceUpdate(d, meta)
}

func resourceAliCloudDtsInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dtsService := DtsService{client}

	object, err := dtsService.DescribeDtsInstance(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_dts_instance dtsService.DescribeDtsInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("create_time", object["CreateTime"])
	d.Set("dts_instance_id", object["DtsInstanceId"])
	d.Set("destination_endpoint_engine_name", object["DestEndpointEngineType"])
	d.Set("instance_class", object["InstanceClass"])
	d.Set("instance_name", object["InstanceName"])
	d.Set("payment_type", convertDTSInstancePaymentTypeResponse(object["PayType"]))
	d.Set("resource_group_id", object["ResourceGroupId"])
	d.Set("source_endpoint_engine_name", object["SourceEndpointEngineType"])
	d.Set("source_region", object["SourceEndpointRegion"])
	d.Set("status", object["Status"])
	d.Set("destination_region", object["DestEndpointRegion"])
	tagsMap := make(map[string]interface{})
	tagsRaw, _ := jsonpath.Get("$.Tags", object)
	if tagsRaw != nil {
		for _, value0 := range tagsRaw.([]interface{}) {
			tags := value0.(map[string]interface{})
			key := tags["TagKey"].(string)
			value := tags["TagValue"]
			if !ignoredTags(key, value) {
				tagsMap[key] = value
			}
		}
	}
	if len(tagsMap) > 0 {
		d.Set("tags", tagsMap)
	}
	d.Set("type", object["Type"])

	return nil
}

func resourceAliCloudDtsInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var err error
	dtsService := DtsService{client}
	var response map[string]interface{}
	d.Partial(true)
	update := false
	request := map[string]interface{}{
		"ResourceId": d.Id(),
	}

	if !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
		if v, ok := d.GetOk("resource_group_id"); ok {
			request["NewResourceGroupId"] = v
		}
	}

	if update {
		action := "ConvertInstanceResourceGroup"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("Dts", "2020-01-01", action, nil, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("resource_group_id")
	}

	if d.HasChange("tags") {
		if err := dtsService.SetResourceTags(d, "ALIYUN::DTS::INSTANCE"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}
	d.Partial(false)
	return resourceAliCloudDtsInstanceRead(d, meta)
}

func resourceAliCloudDtsInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	if d.Get("payment_type").(string) == "Subscription" {
		log.Printf("[WARN] Cannot destroy resource AliCloud Resource DTS Instance. Terraform will remove this resource from the state file, however resources may remain.")
		return nil
	}
	client := meta.(*connectivity.AliyunClient)
	var err error
	var response map[string]interface{}

	request := map[string]interface{}{
		"DtsInstanceId": d.Id(),
		"RegionId":      client.RegionId,
	}

	if v, ok := d.GetOk("synchronization_direction"); ok {
		request["SynchronizationDirection"] = v
	}

	action := "DeleteDtsJob"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Dts", "2020-01-01", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}

func convertDTSInstancePaymentTypeResponse(source interface{}) interface{} {
	switch source {
	case "POSTPAY":
		return "PayAsYouGo"
	case "PREPAY":
		return "Subscription"
	}
	return source
}

func convertDTSInstancePaymentTypeRequest(source interface{}) interface{} {
	switch source {
	case "PayAsYouGo":
		return "PostPaid"
	case "Subscription":
		return "PrePaid"
	}
	return source
}
