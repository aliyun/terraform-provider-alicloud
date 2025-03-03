// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/alibabacloud-go/tea/tea"
	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudSlsLogStore() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudSlsLogStoreCreate,
		Read:   resourceAliCloudSlsLogStoreRead,
		Update: resourceAliCloudSlsLogStoreUpdate,
		Delete: resourceAliCloudSlsLogStoreDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"append_meta": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"auto_split": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"enable_web_tracking": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"encrypt_conf": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"encrypt_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"enable": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"user_cmk_info": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							ForceNew: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cmk_key_id": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										ForceNew: true,
									},
									"region_id": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										ForceNew: true,
									},
									"arn": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										ForceNew: true,
									},
								},
							},
						},
					},
				},
			},
			"hot_ttl": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"infrequent_access_ttl": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"logstore_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"logstore_name", "name"},
				ForceNew:     true,
			},
			"max_split_shard_count": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntBetween(0, 256),
			},
			"metering_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"mode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if new == "" {
						return true
					}
					return old != "" && new != "" && old == new
				},
			},
			"project_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"project_name", "project"},
				ForceNew:     true,
			},
			"retention_period": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  30,
			},
			"shard_count": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if old == "" {
						return false
					}
					return true
				},
				Default: 2,
			},
			"telemetry_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"project": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: "Field 'project' has been deprecated since provider version 1.215.0. New field 'project_name' instead.",
				ForceNew:   true,
			},
			"name": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: "Field 'name' has been deprecated since provider version 1.215.0. New field 'logstore_name' instead.",
				ForceNew:   true,
			},
			"shards": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"begin_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"end_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceAliCloudSlsLogStoreCreate(d *schema.ResourceData, meta interface{}) error {

	if v, ok := d.GetOk("telemetry_type"); ok && v == "Metrics" {
		client := meta.(*connectivity.AliyunClient)
		projectName := d.Get("project_name").(string)
		if v, ok := d.GetOkExists("project"); ok {
			projectName = v.(string)
		}
		logstoreName := d.Get("logstore_name").(string)
		if v, ok := d.GetOkExists("name"); ok {
			logstoreName = v.(string)
		}

		logstore := buildLogStore(d)
		var requestinfo *sls.Client
		err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
			raw, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
				return nil, slsClient.CreateMetricStore(projectName, logstore)
			})
			if err != nil {
				if IsExpectedErrors(err, []string{"InternalServerError", LogClientTimeout}) {
					time.Sleep(10 * time.Second)
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug("CreateMetricStore", raw, requestinfo, map[string]interface{}{
				"project":  projectName,
				"logstore": logstore,
			})
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "alicloud_log_store", "CreateLogStoreV2", AliyunLogGoSdkERROR)
		}
		d.SetId(fmt.Sprintf("%s%s%s", projectName, COLON_SEPARATED, logstoreName))
		return resourceAliCloudSlsLogStoreUpdate(d, meta)
	}

	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/logstores")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	hostMap := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	request["logstoreName"] = d.Get("logstore_name")
	hostMap["project"] = StringPointer(d.Get("project_name").(string))
	if v, ok := d.GetOkExists("project"); ok {
		hostMap["project"] = tea.String(v.(string))
	}
	if v, ok := d.GetOkExists("name"); ok {
		request["logstoreName"] = v
	}

	request["shardCount"] = 2
	if v, ok := d.GetOkExists("shard_count"); ok {
		request["shardCount"] = v
	}
	if v, ok := d.GetOk("auto_split"); ok {
		request["autoSplit"] = v
	}
	if v, ok := d.GetOk("append_meta"); ok {
		request["appendMeta"] = v
	}
	if v, ok := d.GetOk("telemetry_type"); ok {
		request["telemetryType"] = v
	}
	if v, ok := d.GetOk("hot_ttl"); ok {
		request["hot_ttl"] = v
	}
	if v, ok := d.GetOk("mode"); ok {
		request["mode"] = v
	}
	objectDataLocalMap := make(map[string]interface{})

	if v := d.Get("encrypt_conf"); !IsNil(v) {
		enable1, _ := jsonpath.Get("$[0].enable", d.Get("encrypt_conf"))
		if enable1 != nil && enable1 != "" {
			objectDataLocalMap["enable"] = enable1
		}
		encryptType, _ := jsonpath.Get("$[0].encrypt_type", d.Get("encrypt_conf"))
		if encryptType != nil && encryptType != "" {
			objectDataLocalMap["encrypt_type"] = encryptType
		}
		user_cmk_info := make(map[string]interface{})
		cmkKeyId, _ := jsonpath.Get("$[0].user_cmk_info[0].cmk_key_id", d.Get("encrypt_conf"))
		if cmkKeyId != nil && cmkKeyId != "" {
			user_cmk_info["cmk_key_id"] = cmkKeyId
		}
		arn1, _ := jsonpath.Get("$[0].user_cmk_info[0].arn", d.Get("encrypt_conf"))
		if arn1 != nil && arn1 != "" {
			user_cmk_info["arn"] = arn1
		}
		regionId, _ := jsonpath.Get("$[0].user_cmk_info[0].region_id", d.Get("encrypt_conf"))
		if regionId != nil && regionId != "" {
			user_cmk_info["region_id"] = regionId
		}

		user_cmk_info_map, _ := jsonpath.Get("$[0].user_cmk_info[0]", v)
		if !IsNil(user_cmk_info_map) {
			objectDataLocalMap["user_cmk_info"] = user_cmk_info
		}

		request["encrypt_conf"] = objectDataLocalMap
	}

	request["ttl"] = 30
	if v, ok := d.GetOk("retention_period"); ok {
		request["ttl"] = v
	}
	if v, ok := d.GetOk("max_split_shard_count"); ok {
		request["maxSplitShard"] = v
	}
	if v, ok := d.GetOk("enable_web_tracking"); ok {
		request["enable_tracking"] = v
	}
	if v, ok := d.GetOk("infrequent_access_ttl"); ok {
		request["infrequentAccessTTL"] = v
	}

	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.Do("Sls", roaParam("POST", "2020-12-30", "CreateLogStore", action), query, body, nil, hostMap, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_log_store", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", *hostMap["project"], request["logstoreName"]))

	return resourceAliCloudSlsLogStoreUpdate(d, meta)
}

func resourceAliCloudSlsLogStoreRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slsServiceV2 := SlsServiceV2{client}

	objectRaw, err := slsServiceV2.DescribeSlsLogStore(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_log_store DescribeSlsLogStore Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["appendMeta"] != nil {
		d.Set("append_meta", objectRaw["appendMeta"])
	}
	if objectRaw["autoSplit"] != nil {
		d.Set("auto_split", objectRaw["autoSplit"])
	}
	if objectRaw["createTime"] != nil {
		d.Set("create_time", objectRaw["createTime"])
	}
	if objectRaw["enable_tracking"] != nil {
		d.Set("enable_web_tracking", objectRaw["enable_tracking"])
	}
	if objectRaw["hot_ttl"] != nil {
		d.Set("hot_ttl", objectRaw["hot_ttl"])
	}
	if objectRaw["infrequentAccessTTL"] != nil {
		d.Set("infrequent_access_ttl", objectRaw["infrequentAccessTTL"])
	}
	if objectRaw["maxSplitShard"] != nil {
		d.Set("max_split_shard_count", objectRaw["maxSplitShard"])
	}
	if objectRaw["mode"] != nil {
		d.Set("mode", objectRaw["mode"])
	}
	if objectRaw["ttl"] != nil {
		d.Set("retention_period", objectRaw["ttl"])
	}
	if objectRaw["shardCount"] != nil {
		d.Set("shard_count", objectRaw["shardCount"])
	}
	if objectRaw["telemetryType"] != nil {
		d.Set("telemetry_type", objectRaw["telemetryType"])
	}
	if objectRaw["logstoreName"] != nil {
		d.Set("logstore_name", objectRaw["logstoreName"])
	}

	encryptConfMaps := make([]map[string]interface{}, 0)
	encryptConfMap := make(map[string]interface{})
	encrypt_conf1Raw := make(map[string]interface{})
	if objectRaw["encrypt_conf"] != nil {
		encrypt_conf1Raw = objectRaw["encrypt_conf"].(map[string]interface{})
	}
	if len(encrypt_conf1Raw) > 0 {
		encryptConfMap["enable"] = encrypt_conf1Raw["enable"]
		encryptConfMap["encrypt_type"] = encrypt_conf1Raw["encrypt_type"]

		userCmkInfoMaps := make([]map[string]interface{}, 0)
		userCmkInfoMap := make(map[string]interface{})
		user_cmk_info1Raw := make(map[string]interface{})
		if encrypt_conf1Raw["user_cmk_info"] != nil {
			user_cmk_info1Raw = encrypt_conf1Raw["user_cmk_info"].(map[string]interface{})
		}
		if len(user_cmk_info1Raw) > 0 {
			userCmkInfoMap["arn"] = user_cmk_info1Raw["arn"]
			userCmkInfoMap["cmk_key_id"] = user_cmk_info1Raw["cmk_key_id"]
			userCmkInfoMap["region_id"] = user_cmk_info1Raw["region_id"]

			userCmkInfoMaps = append(userCmkInfoMaps, userCmkInfoMap)
		}
		encryptConfMap["user_cmk_info"] = userCmkInfoMaps
		encryptConfMaps = append(encryptConfMaps, encryptConfMap)
	}
	if objectRaw["encrypt_conf"] != nil {
		if err := d.Set("encrypt_conf", encryptConfMaps); err != nil {
			return err
		}
	}

	objectRaw, err = slsServiceV2.DescribeGetLogStoreMeteringMode(d.Id())
	if err != nil {
		return WrapError(err)
	}

	if objectRaw["meteringMode"] != nil {
		d.Set("metering_mode", objectRaw["meteringMode"])
	}

	parts := strings.Split(d.Id(), ":")
	d.Set("project_name", parts[0])
	d.Set("logstore_name", parts[1])

	d.Set("project", d.Get("project_name"))
	d.Set("name", d.Get("logstore_name"))
	logService := LogService{client}
	object, err := logService.DescribeLogStore(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}
	var shards []*sls.Shard
	err = resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
		shards, err = object.ListShards()
		if err != nil {
			if IsExpectedErrors(err, []string{"InternalServerError"}) {
				time.Sleep(10 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug("ListShards", shards)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_log_store", "ListShards", AliyunLogGoSdkERROR)
	}
	var shardList []map[string]interface{}
	for _, s := range shards {
		mapping := map[string]interface{}{
			"id":        s.ShardID,
			"status":    s.Status,
			"begin_key": s.InclusiveBeginKey,
			"end_key":   s.ExclusiveBeginKey,
		}
		shardList = append(shardList, mapping)
	}
	d.Set("shards", shardList)
	return nil
}

func resourceAliCloudSlsLogStoreUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false
	d.Partial(true)
	parts := strings.Split(d.Id(), ":")
	logstore := parts[1]
	action := fmt.Sprintf("/logstores/%s", logstore)
	var err error
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	hostMap := make(map[string]*string)
	hostMap["project"] = StringPointer(parts[0])

	if !d.IsNewResource() && d.HasChange("auto_split") {
		update = true
	}
	if v, ok := d.GetOk("auto_split"); ok || d.HasChange("auto_split") {
		request["autoSplit"] = v
	}
	if !d.IsNewResource() && d.HasChange("append_meta") {
		update = true
	}
	if v, ok := d.GetOk("append_meta"); ok || d.HasChange("append_meta") {
		request["appendMeta"] = v
	}
	if !d.IsNewResource() && d.HasChange("hot_ttl") {
		update = true
	}
	if v, ok := d.GetOk("hot_ttl"); ok || d.HasChange("hot_ttl") {
		request["hot_ttl"] = v
	}
	if !d.IsNewResource() && d.HasChange("mode") {
		update = true
	}
	if v, ok := d.GetOk("mode"); ok || d.HasChange("mode") {
		request["mode"] = v
	}
	if !d.IsNewResource() && d.HasChange("retention_period") {
		update = true
	}
	request["ttl"] = 30
	if v, ok := d.GetOk("retention_period"); ok {
		request["ttl"] = v
	}
	if !d.IsNewResource() && d.HasChange("max_split_shard_count") {
		update = true
	}
	if v, ok := d.GetOk("max_split_shard_count"); ok || d.HasChange("max_split_shard_count") {
		request["maxSplitShard"] = v
	}
	if !d.IsNewResource() && d.HasChange("enable_web_tracking") {
		update = true
	}
	if v, ok := d.GetOk("enable_web_tracking"); ok || d.HasChange("enable_web_tracking") {
		request["enable_tracking"] = v
	}
	if !d.IsNewResource() && d.HasChange("encrypt_conf") {
		update = true
	}
	objectDataLocalMap := make(map[string]interface{})

	if v := d.Get("encrypt_conf"); !IsNil(v) || d.HasChange("encrypt_conf") {
		enable1, _ := jsonpath.Get("$[0].enable", v)
		if enable1 != nil && (d.HasChange("encrypt_conf.0.enable") || enable1 != "") {
			objectDataLocalMap["enable"] = enable1
		}
		encryptType, _ := jsonpath.Get("$[0].encrypt_type", v)
		if encryptType != nil && (d.HasChange("encrypt_conf.0.encrypt_type") || encryptType != "") {
			objectDataLocalMap["encrypt_type"] = encryptType
		}
		user_cmk_info := make(map[string]interface{})
		cmkKeyId, _ := jsonpath.Get("$[0].user_cmk_info[0].cmk_key_id", v)
		if cmkKeyId != nil && (d.HasChange("encrypt_conf.0.user_cmk_info.0.cmk_key_id") || cmkKeyId != "") {
			user_cmk_info["cmk_key_id"] = cmkKeyId
		}
		arn1, _ := jsonpath.Get("$[0].user_cmk_info[0].arn", v)
		if arn1 != nil && (d.HasChange("encrypt_conf.0.user_cmk_info.0.arn") || arn1 != "") {
			user_cmk_info["arn"] = arn1
		}
		regionId, _ := jsonpath.Get("$[0].user_cmk_info[0].region_id", v)
		if regionId != nil && (d.HasChange("encrypt_conf.0.user_cmk_info.0.region_id") || regionId != "") {
			user_cmk_info["region_id"] = regionId
		}

		user_cmk_info_map, _ := jsonpath.Get("$[0].user_cmk_info[0]", v)
		if !IsNil(user_cmk_info_map) {
			objectDataLocalMap["user_cmk_info"] = user_cmk_info
		}

		request["encrypt_conf"] = objectDataLocalMap
	}

	if !d.IsNewResource() && d.HasChange("infrequent_access_ttl") {
		update = true
	}
	if v, ok := d.GetOk("infrequent_access_ttl"); ok || d.HasChange("infrequent_access_ttl") {
		request["infrequentAccessTTL"] = v
	}

	if v, ok := d.GetOk("telemetry_type"); ok && v == "Metrics" {

		projectName := d.Get("project_name").(string)
		if v, ok := d.GetOkExists("project"); ok {
			projectName = v.(string)
		}

		logstore := buildLogStore(d)
		var requestinfo *sls.Client
		err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
			raw, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
				return nil, slsClient.UpdateMetricStore(projectName, logstore)
			})
			if err != nil {
				if IsExpectedErrors(err, []string{"InternalServerError", LogClientTimeout}) {
					time.Sleep(10 * time.Second)
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug("UpdateMetricStore", raw, requestinfo, map[string]interface{}{
				"project":  projectName,
				"logstore": logstore,
			})
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "alicloud_log_store", "CreateLogStoreV2", AliyunLogGoSdkERROR)
		}

		update = false
	}
	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.Do("Sls", roaParam("PUT", "2020-12-30", "UpdateLogStore", action), query, body, nil, hostMap, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}
	update = false
	parts = strings.Split(d.Id(), ":")
	logstore = parts[1]
	action = fmt.Sprintf("/logstores/%s/meteringmode", logstore)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	hostMap = make(map[string]*string)
	hostMap["project"] = StringPointer(parts[0])

	slsServiceV2 := SlsServiceV2{client}
	objectRaw, _ := slsServiceV2.DescribeGetLogStoreMeteringMode(d.Id())
	if d.HasChange("metering_mode") && objectRaw["meteringMode"] != d.Get("metering_mode") {
		update = true
	}
	request["meteringMode"] = d.Get("metering_mode")

	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.Do("Sls", roaParam("PUT", "2020-12-30", "UpdateLogStoreMeteringMode", action), query, body, nil, hostMap, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	d.Partial(false)
	return resourceAliCloudSlsLogStoreRead(d, meta)
}

func resourceAliCloudSlsLogStoreDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	logstore := parts[1]
	action := fmt.Sprintf("/logstores/%s", logstore)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	hostMap := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	hostMap["project"] = StringPointer(parts[0])

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.Do("Sls", roaParam("DELETE", "2020-12-30", "DeleteLogStore", action), query, nil, nil, hostMap, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}

func buildLogStore(d *schema.ResourceData) *sls.LogStore {
	logstore := &sls.LogStore{
		Name:          d.Get("logstore_name").(string),
		TTL:           d.Get("retention_period").(int),
		ShardCount:    d.Get("shard_count").(int),
		WebTracking:   d.Get("enable_web_tracking").(bool),
		AutoSplit:     d.Get("auto_split").(bool),
		MaxSplitShard: d.Get("max_split_shard_count").(int),
		AppendMeta:    d.Get("append_meta").(bool),
		TelemetryType: d.Get("telemetry_type").(string),
		Mode:          d.Get("mode").(string),
	}
	if v, ok := d.GetOkExists("name"); ok {
		logstore.Name = v.(string)
	}
	if hotTTL, ok := d.GetOk("hot_ttl"); ok {
		logstore.HotTTL = int32(hotTTL.(int))
	}
	if encrypt := buildEncrypt(d); encrypt != nil {
		logstore.EncryptConf = encrypt
	}

	return logstore
}

func buildEncrypt(d *schema.ResourceData) *sls.EncryptConf {
	var encryptConf *sls.EncryptConf
	if field, ok := d.GetOk("encrypt_conf"); ok {
		encryptConf = new(sls.EncryptConf)
		value := field.([]interface{})[0].(map[string]interface{})
		encryptConf.Enable = value["enable"].(bool)
		encryptConf.EncryptType = value["encrypt_type"].(string)
		cmkInfo := value["user_cmk_info"].([]interface{})
		if len(cmkInfo) > 0 {
			cmk := cmkInfo[0].(map[string]interface{})
			encryptConf.UserCmkInfo = &sls.EncryptUserCmkConf{
				CmkKeyId: cmk["cmk_key_id"].(string),
				Arn:      cmk["arn"].(string),
				RegionId: cmk["region_id"].(string),
			}
		}
	}
	return encryptConf
}
