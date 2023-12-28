// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
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

	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/logstores")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	hostMap := make(map[string]*string)
	conn, err := client.NewSlsClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["logstoreName"] = d.Get("logstore_name")
	hostMap["project"] = StringPointer(d.Get("project_name").(string))
	if v, ok := d.GetOkExists("project"); ok {
		hostMap["project"] = tea.String(v.(string))
	}
	if v, ok := d.GetOkExists("name"); ok {
		request["logstoreName"] = v
	}

	if v, ok := d.GetOk("shard_count"); ok {
		request["shardCount"] = v
	}
	if v, ok := d.GetOkExists("auto_split"); ok {
		request["autoSplit"] = v
	}
	if v, ok := d.GetOkExists("append_meta"); ok {
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
		nodeNative, _ := jsonpath.Get("$[0].enable", v)
		if nodeNative != "" {
			objectDataLocalMap["enable"] = nodeNative
		}
		nodeNative1, _ := jsonpath.Get("$[0].encrypt_type", v)
		if nodeNative1 != "" {
			objectDataLocalMap["encrypt_type"] = nodeNative1
		}
		user_cmk_info := make(map[string]interface{})
		nodeNative2, _ := jsonpath.Get("$[0].user_cmk_info[0].cmk_key_id", v)
		if nodeNative2 != "" {
			user_cmk_info["cmk_key_id"] = nodeNative2
		}
		nodeNative3, _ := jsonpath.Get("$[0].user_cmk_info[0].arn", v)
		if nodeNative3 != "" {
			user_cmk_info["arn"] = nodeNative3
		}
		nodeNative4, _ := jsonpath.Get("$[0].user_cmk_info[0].region_id", v)
		if nodeNative4 != "" {
			user_cmk_info["region_id"] = nodeNative4
		}
		objectDataLocalMap["user_cmk_info"] = user_cmk_info
		request["encrypt_conf"] = objectDataLocalMap
	}

	if v, ok := d.GetOk("retention_period"); ok {
		request["ttl"] = v
	}
	if v, ok := d.GetOk("max_split_shard_count"); ok {
		request["maxSplitShard"] = v
	}
	if v, ok := d.GetOkExists("enable_web_tracking"); ok {
		request["enable_tracking"] = v
	}
	body = request
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.Execute(genRoaParam("CreateLogStore", "POST", "2020-12-30", action), &openapi.OpenApiRequest{Query: query, Body: body, HostMap: hostMap}, &util.RuntimeOptions{})

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

	return resourceAliCloudSlsLogStoreRead(d, meta)
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

	d.Set("append_meta", objectRaw["appendMeta"])
	d.Set("auto_split", objectRaw["autoSplit"])
	d.Set("create_time", objectRaw["createTime"])
	d.Set("enable_web_tracking", objectRaw["enable_tracking"])
	d.Set("hot_ttl", objectRaw["hot_ttl"])
	d.Set("max_split_shard_count", objectRaw["maxSplitShard"])
	d.Set("mode", objectRaw["mode"])
	d.Set("retention_period", objectRaw["ttl"])
	d.Set("shard_count", objectRaw["shardCount"])
	d.Set("telemetry_type", objectRaw["telemetryType"])
	d.Set("logstore_name", objectRaw["logstoreName"])

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
	d.Set("encrypt_conf", encryptConfMaps)
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
	var hostMap map[string]*string
	update := false
	parts := strings.Split(d.Id(), ":")
	logstore := parts[1]
	action := fmt.Sprintf("/logstores/%s", logstore)
	conn, err := client.NewSlsClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	hostMap = make(map[string]*string)
	hostMap["project"] = StringPointer(parts[0])
	hostMap["logstore"] = StringPointer(parts[1])
	if d.HasChange("auto_split") {
		update = true
	}
	request["autoSplit"] = d.Get("auto_split")
	if d.HasChange("append_meta") {
		update = true
	}
	request["appendMeta"] = d.Get("append_meta")
	if d.HasChange("hot_ttl") {
		update = true
	}
	request["hot_ttl"] = d.Get("hot_ttl")
	if d.HasChange("mode") {
		update = true
	}
	request["mode"] = d.Get("mode")
	if d.HasChange("retention_period") {
		update = true
	}
	request["ttl"] = 30
	if v, ok := d.GetOk("retention_period"); ok {
		request["ttl"] = v
	}
	if d.HasChange("max_split_shard_count") {
		update = true
	}
	if v, ok := d.GetOk("max_split_shard_count"); ok {
		request["maxSplitShard"] = v
	}
	if d.HasChange("enable_web_tracking") {
		update = true
	}
	request["enable_tracking"] = d.Get("enable_web_tracking")
	if d.HasChange("encrypt_conf") {
		update = true
	}
	objectDataLocalMap := make(map[string]interface{})
	if v := d.Get("encrypt_conf"); !IsNil(v) {
		nodeNative, _ := jsonpath.Get("$[0].enable", v)
		if nodeNative != "" {
			objectDataLocalMap["enable"] = nodeNative
		}
		nodeNative1, _ := jsonpath.Get("$[0].encrypt_type", v)
		if nodeNative1 != "" {
			objectDataLocalMap["encrypt_type"] = nodeNative1
		}
		user_cmk_info := make(map[string]interface{})
		nodeNative2, _ := jsonpath.Get("$[0].user_cmk_info[0].cmk_key_id", v)
		if nodeNative2 != "" {
			user_cmk_info["cmk_key_id"] = nodeNative2
		}
		nodeNative3, _ := jsonpath.Get("$[0].user_cmk_info[0].arn", v)
		if nodeNative3 != "" {
			user_cmk_info["arn"] = nodeNative3
		}
		nodeNative4, _ := jsonpath.Get("$[0].user_cmk_info[0].region_id", v)
		if nodeNative4 != "" {
			user_cmk_info["region_id"] = nodeNative4
		}
		objectDataLocalMap["user_cmk_info"] = user_cmk_info
		request["encrypt_conf"] = objectDataLocalMap
	}

	body = request
	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.Execute(genRoaParam("UpdateLogStore", "PUT", "2020-12-30", action), &openapi.OpenApiRequest{Query: query, Body: body, HostMap: hostMap}, &util.RuntimeOptions{})

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
	conn, err := client.NewSlsClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	hostMap["logstore"] = StringPointer(parts[1])
	hostMap["project"] = StringPointer(parts[0])

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.Execute(genRoaParam("DeleteLogStore", "DELETE", "2020-12-30", action), &openapi.OpenApiRequest{Query: query, Body: nil, HostMap: hostMap}, &util.RuntimeOptions{})

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
