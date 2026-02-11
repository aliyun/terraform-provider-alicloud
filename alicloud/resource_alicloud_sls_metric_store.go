package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudSlsMetricStore() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudSlsMetricStoreCreate,
		Read:   resourceAliCloudSlsMetricStoreRead,
		Update: resourceAliCloudSlsMetricStoreUpdate,
		Delete: resourceAliCloudSlsMetricStoreDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"auto_split": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"hot_ttl": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntBetween(0, 3650),
			},
			"infrequent_access_ttl": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntBetween(0, 3650),
			},
			"max_split_shard": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntBetween(0, 256),
			},
			"metering_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"ChargeByFunction", "ChargeByDataIngest"}, false),
			},
			"metric_store_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"metric_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"prometheus"}, false),
			},
			"mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"standard"}, false),
			},
			"project_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"shard_count": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: IntAtLeast(1),
			},
			"ttl": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: IntBetween(0, 3650),
			},
		},
	}
}

func resourceAliCloudSlsMetricStoreCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/metricstores")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	hostMap := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	hostMap["project"] = StringPointer(d.Get("project_name").(string))
	if v, ok := d.GetOk("metric_store_name"); ok {
		request["name"] = v
	}

	if v, ok := d.GetOkExists("infrequent_access_ttl"); ok && v.(int) > 0 {
		request["infrequentAccessTTL"] = v
	}
	if v, ok := d.GetOkExists("hot_ttl"); ok && v.(int) > 0 {
		request["hot_ttl"] = v
	}
	request["ttl"] = d.Get("ttl")
	if v, ok := d.GetOk("metric_type"); ok {
		request["metricType"] = v
	}
	if v, ok := d.GetOk("mode"); ok {
		request["mode"] = v
	}
	if v, ok := d.GetOkExists("auto_split"); ok {
		request["autoSplit"] = v
	}
	request["shardCount"] = d.Get("shard_count")
	if v, ok := d.GetOkExists("max_split_shard"); ok {
		request["maxSplitShard"] = v
	}
	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.Do("Sls", roaParam("POST", "2020-12-30", "CreateMetricStore", action), query, body, nil, hostMap, false)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_sls_metric_store", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", *hostMap["project"], request["name"]))

	return resourceAliCloudSlsMetricStoreUpdate(d, meta)
}

func resourceAliCloudSlsMetricStoreRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slsServiceV2 := SlsServiceV2{client}

	objectRaw, err := slsServiceV2.DescribeSlsMetricStore(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_sls_metric_store DescribeSlsMetricStore Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["createTime"])
	d.Set("max_split_shard", objectRaw["maxSplitShard"])
	d.Set("metric_type", objectRaw["metricType"])
	d.Set("mode", objectRaw["mode"])
	d.Set("shard_count", objectRaw["shardCount"])
	d.Set("ttl", objectRaw["ttl"])
	d.Set("metric_store_name", objectRaw["name"])

	parts := strings.Split(d.Id(), ":")
	d.Set("project_name", parts[0])

	return nil
}

func resourceAliCloudSlsMetricStoreUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false
	d.Partial(true)

	var err error
	parts := strings.Split(d.Id(), ":")
	name := parts[1]
	action := fmt.Sprintf("/metricstores/%s", name)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	hostMap := make(map[string]*string)
	hostMap["project"] = StringPointer(parts[0])

	if !d.IsNewResource() && d.HasChange("infrequent_access_ttl") {
		update = true
	}
	if v, ok := d.GetOkExists("infrequent_access_ttl"); (ok || d.HasChange("infrequent_access_ttl")) && v.(int) > 0 {
		request["infrequentAccessTTL"] = v
	}
	if !d.IsNewResource() && d.HasChange("hot_ttl") {
		update = true
	}
	if v, ok := d.GetOkExists("hot_ttl"); (ok || d.HasChange("hot_ttl")) && v.(int) > 0 {
		request["hot_ttl"] = v
	}
	if !d.IsNewResource() && d.HasChange("ttl") {
		update = true
	}
	request["ttl"] = d.Get("ttl")
	if !d.IsNewResource() && d.HasChange("mode") {
		update = true
	}
	if v, ok := d.GetOk("mode"); ok || d.HasChange("mode") {
		request["mode"] = v
	}
	if !d.IsNewResource() && d.HasChange("auto_split") {
		update = true
	}
	if v, ok := d.GetOkExists("auto_split"); ok || d.HasChange("auto_split") {
		request["autoSplit"] = v
	}
	if !d.IsNewResource() && d.HasChange("max_split_shard") {
		update = true
	}
	if v, ok := d.GetOkExists("max_split_shard"); ok || d.HasChange("max_split_shard") {
		request["maxSplitShard"] = v
	}
	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.Do("Sls", roaParam("PUT", "2020-12-30", "UpdateMetricStore", action), query, body, nil, hostMap, false)
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
	}
	update = false
	parts = strings.Split(d.Id(), ":")
	metricStore := parts[1]
	action = fmt.Sprintf("/metricstores/%s/meteringmode", metricStore)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	hostMap = make(map[string]*string)
	hostMap["project"] = StringPointer(parts[0])

	if d.HasChange("metering_mode") {
		update = true
	}
	request["meteringMode"] = d.Get("metering_mode")
	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.Do("Sls", roaParam("PUT", "2020-12-30", "UpdateMetricStoreMeteringMode", action), query, body, nil, hostMap, false)
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
	}

	d.Partial(false)
	return resourceAliCloudSlsMetricStoreRead(d, meta)
}

func resourceAliCloudSlsMetricStoreDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	name := parts[1]
	action := fmt.Sprintf("/metricstores/%s", name)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	hostMap := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	hostMap["project"] = StringPointer(parts[0])

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.Do("Sls", roaParam("DELETE", "2020-12-30", "DeleteMetricStore", action), query, nil, nil, hostMap, false)
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
		if IsExpectedErrors(err, []string{"MetricStoreNotExist"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
