// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!

// Hand-maintained differences from the generator output. The generator has not absorbed them yet;
// keep them when this file is regenerated:
//  1. A transient 400 InternalError is retried with a bound on create, on the status switch, on
//     update and on delete (esaSiteDeliveryTaskWriteRetryable) and on ListSiteDeliveryTasks
//     (esaSiteDeliveryTaskListRetryable); the InternalError code is only retried together with the
//     service's "please try again later" hint, and other 4xx stay non-retryable.
//  2. http_delivery header_param/query_param values are sent as {"<key>":{"StaticValue":"<value>"}}.
//  3. Create finishes with Read and only an explicit status=offline switches the task; Update
//     switches status only on a real status change instead of a create tail-calling Update.
//  4. Existence and destroy are decided by a fully successful, fully paginated
//     ListSiteDeliveryTasks (a malformed page is an error, never an empty page), because
//     GetSiteDeliveryTask reports an absent task as HTTP 400 InternalError instead of 404.
//
// (1) is expressible through the resource-to-API mapping (per-operation retry policies); (2)-(4)
// are generator template behaviours that still need a generator-side change. Regenerating this file
// before the generator carries them would silently drop them.
package alicloud

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAliCloudEsaSiteDeliveryTask() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEsaSiteDeliveryTaskCreate,
		Read:   resourceAliCloudEsaSiteDeliveryTaskRead,
		Update: resourceAliCloudEsaSiteDeliveryTaskUpdate,
		Delete: resourceAliCloudEsaSiteDeliveryTaskDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"business_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"data_center": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"delivery_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"discard_rate": {
				Type:     schema.TypeFloat,
				Optional: true,
			},
			"field_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"http_delivery": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"compress": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"log_body_suffix": {
							Type:     schema.TypeString,
							Optional: true,
						},
						// lintignore: S006
						"header_param": {
							Type:     schema.TypeMap,
							Optional: true,
						},
						"standard_auth_param": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"private_key": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"url_path": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"expired_time": {
										Type:     schema.TypeInt,
										Optional: true,
									},
								},
							},
						},
						"standard_auth_on": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"log_body_prefix": {
							Type:     schema.TypeString,
							Optional: true,
						},
						// lintignore: S006
						"query_param": {
							Type:     schema.TypeMap,
							Optional: true,
						},
						"dest_url": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"max_batch_size": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"max_retry": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"transform_timeout": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"max_batch_mb": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"kafka_delivery": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"compress": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"user_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"machanism_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"brokers": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"balancer": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"topic": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"user_auth": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"password": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"oss_delivery": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bucket_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"region": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"prefix_path": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"aliuid": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"s3_delivery": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"secret_key": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"endpoint": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"vertify_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"region": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"bucket_path": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"server_side_encryption": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"access_key": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"prefix_path": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"s3_cmpt": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"site_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"sls_delivery": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sls_project": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"sls_region": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"sls_log_store": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"task_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudEsaSiteDeliveryTaskCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateSiteDeliveryTask"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("site_id"); ok {
		request["SiteId"] = v
	}
	if v, ok := d.GetOk("task_name"); ok {
		request["TaskName"] = v
	}

	httpDelivery := make(map[string]interface{})

	if v := d.Get("http_delivery"); !IsNil(v) {
		transformTimeout1, _ := jsonpath.Get("$[0].transform_timeout", v)
		if transformTimeout1 != nil && transformTimeout1 != "" {
			httpDelivery["TransformTimeout"] = transformTimeout1
		}
		standardAuthParam := make(map[string]interface{})
		privateKey1, _ := jsonpath.Get("$[0].standard_auth_param[0].private_key", d.Get("http_delivery"))
		if privateKey1 != nil && privateKey1 != "" {
			standardAuthParam["PrivateKey"] = privateKey1
		}
		expiredTime1, _ := jsonpath.Get("$[0].standard_auth_param[0].expired_time", d.Get("http_delivery"))
		if expiredTime1 != nil && expiredTime1 != "" {
			standardAuthParam["ExpiredTime"] = expiredTime1
		}
		urlPath1, _ := jsonpath.Get("$[0].standard_auth_param[0].url_path", d.Get("http_delivery"))
		if urlPath1 != nil && urlPath1 != "" {
			standardAuthParam["UrlPath"] = urlPath1
		}

		httpDelivery["StandardAuthParam"] = standardAuthParam
		maxBatchMB1, _ := jsonpath.Get("$[0].max_batch_mb", v)
		if maxBatchMB1 != nil && maxBatchMB1 != "" {
			httpDelivery["MaxBatchMB"] = maxBatchMB1
		}
		logBodyPrefix1, _ := jsonpath.Get("$[0].log_body_prefix", v)
		if logBodyPrefix1 != nil && logBodyPrefix1 != "" {
			httpDelivery["LogBodyPrefix"] = logBodyPrefix1
		}
		destUrl1, _ := jsonpath.Get("$[0].dest_url", v)
		if destUrl1 != nil && destUrl1 != "" {
			httpDelivery["DestUrl"] = destUrl1
		}
		headerParam1, _ := jsonpath.Get("$[0].header_param", v)
		if headerParam1 != nil && headerParam1 != "" {
			httpDelivery["HeaderParam"] = esaSiteDeliveryTaskHttpParamValues(headerParam1)
		}
		compress1, _ := jsonpath.Get("$[0].compress", v)
		if compress1 != nil && compress1 != "" {
			httpDelivery["Compress"] = compress1
		}
		maxRetry1, _ := jsonpath.Get("$[0].max_retry", v)
		if maxRetry1 != nil && maxRetry1 != "" {
			httpDelivery["MaxRetry"] = maxRetry1
		}
		standardAuthOn1, _ := jsonpath.Get("$[0].standard_auth_on", v)
		if standardAuthOn1 != nil && standardAuthOn1 != "" {
			httpDelivery["StandardAuthOn"] = standardAuthOn1
		}
		queryParam1, _ := jsonpath.Get("$[0].query_param", v)
		if queryParam1 != nil && queryParam1 != "" {
			httpDelivery["QueryParam"] = esaSiteDeliveryTaskHttpParamValues(queryParam1)
		}
		maxBatchSize1, _ := jsonpath.Get("$[0].max_batch_size", v)
		if maxBatchSize1 != nil && maxBatchSize1 != "" {
			httpDelivery["MaxBatchSize"] = maxBatchSize1
		}
		logBodySuffix1, _ := jsonpath.Get("$[0].log_body_suffix", v)
		if logBodySuffix1 != nil && logBodySuffix1 != "" {
			httpDelivery["LogBodySuffix"] = logBodySuffix1
		}

		httpDeliveryJson, err := json.Marshal(httpDelivery)
		if err != nil {
			return WrapError(err)
		}
		request["HttpDelivery"] = string(httpDeliveryJson)
	}

	ossDelivery := make(map[string]interface{})

	if v := d.Get("oss_delivery"); !IsNil(v) {
		bucketName1, _ := jsonpath.Get("$[0].bucket_name", v)
		if bucketName1 != nil && bucketName1 != "" {
			ossDelivery["BucketName"] = bucketName1
		}
		region1, _ := jsonpath.Get("$[0].region", v)
		if region1 != nil && region1 != "" {
			ossDelivery["Region"] = region1
		}
		aliuid1, _ := jsonpath.Get("$[0].aliuid", v)
		if aliuid1 != nil && aliuid1 != "" {
			ossDelivery["Aliuid"] = aliuid1
		}
		prefixPath1, _ := jsonpath.Get("$[0].prefix_path", v)
		if prefixPath1 != nil && prefixPath1 != "" {
			ossDelivery["PrefixPath"] = prefixPath1
		}

		ossDeliveryJson, err := json.Marshal(ossDelivery)
		if err != nil {
			return WrapError(err)
		}
		request["OssDelivery"] = string(ossDeliveryJson)
	}

	kafkaDelivery := make(map[string]interface{})

	if v := d.Get("kafka_delivery"); !IsNil(v) {
		compress3, _ := jsonpath.Get("$[0].compress", v)
		if compress3 != nil && compress3 != "" {
			kafkaDelivery["Compress"] = compress3
		}
		machanismType1, _ := jsonpath.Get("$[0].machanism_type", v)
		if machanismType1 != nil && machanismType1 != "" {
			kafkaDelivery["MachanismType"] = machanismType1
		}
		userAuth1, _ := jsonpath.Get("$[0].user_auth", v)
		if userAuth1 != nil && userAuth1 != "" {
			kafkaDelivery["UserAuth"] = userAuth1
		}
		password1, _ := jsonpath.Get("$[0].password", v)
		if password1 != nil && password1 != "" {
			kafkaDelivery["Password"] = password1
		}
		topic1, _ := jsonpath.Get("$[0].topic", v)
		if topic1 != nil && topic1 != "" {
			kafkaDelivery["Topic"] = topic1
		}
		userName1, _ := jsonpath.Get("$[0].user_name", v)
		if userName1 != nil && userName1 != "" {
			kafkaDelivery["UserName"] = userName1
		}
		brokers1, _ := jsonpath.Get("$[0].brokers", v)
		if brokers1 != nil && brokers1 != "" {
			kafkaDelivery["Brokers"] = brokers1
		}
		balancer1, _ := jsonpath.Get("$[0].balancer", v)
		if balancer1 != nil && balancer1 != "" {
			kafkaDelivery["Balancer"] = balancer1
		}

		kafkaDeliveryJson, err := json.Marshal(kafkaDelivery)
		if err != nil {
			return WrapError(err)
		}
		request["KafkaDelivery"] = string(kafkaDeliveryJson)
	}

	request["FieldName"] = d.Get("field_name")
	slsDelivery := make(map[string]interface{})

	if v := d.Get("sls_delivery"); !IsNil(v) {
		sLSProject1, _ := jsonpath.Get("$[0].sls_project", v)
		if sLSProject1 != nil && sLSProject1 != "" {
			slsDelivery["SLSProject"] = sLSProject1
		}
		sLSRegion1, _ := jsonpath.Get("$[0].sls_region", v)
		if sLSRegion1 != nil && sLSRegion1 != "" {
			slsDelivery["SLSRegion"] = sLSRegion1
		}
		sLSLogStore1, _ := jsonpath.Get("$[0].sls_log_store", v)
		if sLSLogStore1 != nil && sLSLogStore1 != "" {
			slsDelivery["SLSLogStore"] = sLSLogStore1
		}

		slsDeliveryJson, err := json.Marshal(slsDelivery)
		if err != nil {
			return WrapError(err)
		}
		request["SlsDelivery"] = string(slsDeliveryJson)
	}

	request["DataCenter"] = d.Get("data_center")
	s3Delivery := make(map[string]interface{})

	if v := d.Get("s3_delivery"); !IsNil(v) {
		prefixPath3, _ := jsonpath.Get("$[0].prefix_path", v)
		if prefixPath3 != nil && prefixPath3 != "" {
			s3Delivery["PrefixPath"] = prefixPath3
		}
		accessKey1, _ := jsonpath.Get("$[0].access_key", v)
		if accessKey1 != nil && accessKey1 != "" {
			s3Delivery["AccessKey"] = accessKey1
		}
		s3Cmpt1, _ := jsonpath.Get("$[0].s3_cmpt", v)
		if s3Cmpt1 != nil && s3Cmpt1 != "" {
			s3Delivery["S3Cmpt"] = s3Cmpt1
		}
		region3, _ := jsonpath.Get("$[0].region", v)
		if region3 != nil && region3 != "" {
			s3Delivery["Region"] = region3
		}
		serverSideEncryption1, _ := jsonpath.Get("$[0].server_side_encryption", v)
		if serverSideEncryption1 != nil && serverSideEncryption1 != "" {
			s3Delivery["ServerSideEncryption"] = serverSideEncryption1
		}
		vertifyType1, _ := jsonpath.Get("$[0].vertify_type", v)
		if vertifyType1 != nil && vertifyType1 != "" {
			s3Delivery["VertifyType"] = vertifyType1
		}
		bucketPath1, _ := jsonpath.Get("$[0].bucket_path", v)
		if bucketPath1 != nil && bucketPath1 != "" {
			s3Delivery["BucketPath"] = bucketPath1
		}
		endpoint1, _ := jsonpath.Get("$[0].endpoint", v)
		if endpoint1 != nil && endpoint1 != "" {
			s3Delivery["Endpoint"] = endpoint1
		}
		secretKey1, _ := jsonpath.Get("$[0].secret_key", v)
		if secretKey1 != nil && secretKey1 != "" {
			s3Delivery["SecretKey"] = secretKey1
		}

		s3DeliveryJson, err := json.Marshal(s3Delivery)
		if err != nil {
			return WrapError(err)
		}
		request["S3Delivery"] = string(s3DeliveryJson)
	}

	request["BusinessType"] = d.Get("business_type")
	if v, ok := d.GetOk("discard_rate"); ok {
		request["DiscardRate"] = v
	}
	request["DeliveryType"] = d.Get("delivery_type")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)
		if err != nil {
			if esaSiteDeliveryTaskWriteRetryable(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_esa_site_delivery_task", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", response["SiteId"], response["TaskName"]))

	// CreateSiteDeliveryTask brings a new task up online, so only an explicit offline request needs
	// a status switch here. An online create must not call UpdateSiteDeliveryTaskStatus: the extra
	// call races the just-created task and the service answers 400 InternalError.
	if status, _ := d.Get("status").(string); status == "offline" {
		if err := esaSiteDeliveryTaskUpdateStatus(client, d.Id(), "offline", d.Timeout(schema.TimeoutCreate)); err != nil {
			return err
		}
	}

	return resourceAliCloudEsaSiteDeliveryTaskRead(d, meta)
}

func resourceAliCloudEsaSiteDeliveryTaskRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	esaServiceV2 := EsaServiceV2{client}

	objectRaw, err := esaServiceV2.DescribeEsaSiteDeliveryTask(d.Id())
	if err != nil {
		// GetSiteDeliveryTask reports an absent task as HTTP 400 InternalError instead of 404, so a
		// failed Get does not prove the task is gone. Only a fully successful ListSiteDeliveryTasks
		// result can prove absence: a failed List, or a task that is still listed, surfaces the error.
		parts := strings.SplitN(d.Id(), ":", 2)
		if d.IsNewResource() || len(parts) != 2 {
			return WrapError(err)
		}
		exists, listErr := esaServiceV2.EsaSiteDeliveryTaskExists(parts[0], parts[1])
		if listErr != nil {
			return WrapErrorf(listErr, "failed to confirm whether %s still exists after GetSiteDeliveryTask failed with %v", d.Id(), err)
		}
		if !exists {
			log.Printf("[DEBUG] Resource alicloud_esa_site_delivery_task %s is absent from ListSiteDeliveryTasks, removing from state. GetSiteDeliveryTask reported: %v", d.Id(), err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("business_type", objectRaw["BusinessType"])
	d.Set("data_center", objectRaw["DataCenter"])
	d.Set("delivery_type", objectRaw["DeliveryType"])
	d.Set("discard_rate", objectRaw["DiscardRate"])
	d.Set("field_name", objectRaw["FieldList"])
	d.Set("status", objectRaw["Status"])
	if v, ok := objectRaw["SiteId"]; ok {
		d.Set("site_id", v)
	}

	d.Set("task_name", objectRaw["TaskName"])

	return nil
}

func resourceAliCloudEsaSiteDeliveryTaskUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)

	var err error
	parts := strings.Split(d.Id(), ":")
	action := "UpdateSiteDeliveryTask"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["SiteId"] = parts[0]
	request["TaskName"] = parts[1]

	if d.HasChange("discard_rate") {
		update = true
		request["DiscardRate"] = d.Get("discard_rate")
	}

	if d.HasChange("business_type") {
		update = true
	}
	request["BusinessType"] = d.Get("business_type")
	if d.HasChange("field_name") {
		update = true
	}
	request["FieldName"] = d.Get("field_name")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)
			if err != nil {
				if esaSiteDeliveryTaskWriteRetryable(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}
	update = false
	if d.HasChange("status") {
		update = true
	}

	if update {
		if err := esaSiteDeliveryTaskUpdateStatus(client, d.Id(), d.Get("status").(string), d.Timeout(schema.TimeoutUpdate)); err != nil {
			return err
		}
	}

	d.Partial(false)
	return resourceAliCloudEsaSiteDeliveryTaskRead(d, meta)
}

func resourceAliCloudEsaSiteDeliveryTaskDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	esaServiceV2 := EsaServiceV2{client}
	parts := strings.SplitN(d.Id(), ":", 2)
	if len(parts) != 2 {
		return WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", d.Id(), 2, len(parts)))
	}
	siteId, taskName := parts[0], parts[1]
	action := "DeleteSiteDeliveryTask"

	// DeleteSiteDeliveryTask answers 400 InvalidParameter (or 406) for a task that is already gone,
	// so the error alone is not proof of deletion and an already-deleted task is only accepted after
	// a fully successful ListSiteDeliveryTasks shows it absent.
	err := esaSiteDeliveryTaskDeleteWithApis(d.Id(), siteId, taskName, esaServiceV2.EsaSiteDeliveryTaskExists, func() error {
		request := map[string]interface{}{
			"SiteId":   siteId,
			"TaskName": taskName,
		}
		query := make(map[string]interface{})
		var response map[string]interface{}
		var err error
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
			response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)
			if err != nil {
				if esaSiteDeliveryTaskWriteRetryable(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		return err
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}

// esaSiteDeliveryTaskDeleteWithApis runs the List-guarded delete sequence with the existence probe
// and the delete call injected, so the decision order stays verifiable without a live client. A
// failed pre-delete List fails closed: the task is not deleted and the List error is returned.
func esaSiteDeliveryTaskDeleteWithApis(id, siteId, taskName string, exists func(siteId, taskName string) (bool, error), remove func() error) error {
	existedBefore, listErr := exists(siteId, taskName)
	if listErr != nil {
		// Fail closed: a failed List proves nothing about the task, so deleting anyway would hide a
		// real List/pagination/permission error behind a successful DeleteSiteDeliveryTask. The
		// error reaches the resource Delete and the out-of-band cleanup that share this helper.
		return WrapErrorf(listErr, DefaultErrorMsg, id, "ListSiteDeliveryTasks", AlibabaCloudSdkGoERROR)
	}
	if !existedBefore {
		log.Printf("[DEBUG] Resource alicloud_esa_site_delivery_task %s is absent from ListSiteDeliveryTasks, skip delete", id)
		return nil
	}

	deleteErr := remove()
	if deleteErr == nil {
		return nil
	}

	// Converge only when the post-error List succeeds and no longer sees the task. A failed or
	// still-listed result keeps both failures: the delete error alone does not explain why the
	// disappearance could not be proven.
	existsAfter, afterErr := exists(siteId, taskName)
	if esaSiteDeliveryTaskDeleteConverged(existsAfter, afterErr) {
		log.Printf("[DEBUG] Resource alicloud_esa_site_delivery_task %s is gone from ListSiteDeliveryTasks although DeleteSiteDeliveryTask reported: %v", id, deleteErr)
		return nil
	}
	if afterErr != nil {
		return errors.Join(deleteErr, WrapErrorf(afterErr, DefaultErrorMsg, id, "ListSiteDeliveryTasks", AlibabaCloudSdkGoERROR))
	}
	return deleteErr
}

// esaSiteDeliveryTaskWriteRetryable reports whether a write call on a delivery task may be retried.
// Only the service's transient backend failure (InternalError plus the "please try again later"
// hint) is retried; a persistent InternalError with another message and every other 4xx reach the
// caller. A retried create cannot duplicate the task: TaskName is unique per site, so the worst case
// is the same terminal error as without a retry.
func esaSiteDeliveryTaskWriteRetryable(err error) bool {
	return esaSiteDeliveryTaskTransientInternalError(err) || NeedRetry(err)
}

// esaSiteDeliveryTaskUpdateStatus switches an existing task between online and offline.
func esaSiteDeliveryTaskUpdateStatus(client *connectivity.AliyunClient, id, method string, timeout time.Duration) error {
	parts := strings.SplitN(id, ":", 2)
	if len(parts) != 2 {
		return WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
	}
	action := "UpdateSiteDeliveryTaskStatus"
	request := make(map[string]interface{})
	query := map[string]interface{}{
		"SiteId":   parts[0],
		"TaskName": parts[1],
		"Method":   method,
	}
	var response map[string]interface{}
	var err error
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(timeout, func() *retry.RetryError {
		response, err = client.RpcGet("ESA", "2024-09-10", action, query, request)
		if err != nil {
			if esaSiteDeliveryTaskWriteRetryable(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	return nil
}

// esaSiteDeliveryTaskHttpParamValues converts the flat header_param/query_param map into the
// {"<key>":{"StaticValue":"<value>"}} shape CreateSiteDeliveryTask expects.
func esaSiteDeliveryTaskHttpParamValues(params interface{}) map[string]interface{} {
	values := make(map[string]interface{})
	paramMap, ok := params.(map[string]interface{})
	if !ok {
		return values
	}
	for key, value := range paramMap {
		if strValue, ok := value.(string); ok {
			value = map[string]interface{}{"StaticValue": strValue}
		}
		values[key] = value
	}
	return values
}

// esaSiteDeliveryTaskDeleteConverged reports whether a failed DeleteSiteDeliveryTask may be treated
// as success. The caller reaches it only after a successful pre-delete List saw the task, so absence
// is proven only when the post-error List succeeds and no longer sees the task.
func esaSiteDeliveryTaskDeleteConverged(existsAfter bool, afterErr error) bool {
	return afterErr == nil && !existsAfter
}
