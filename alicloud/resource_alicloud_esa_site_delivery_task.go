// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
			httpDelivery["HeaderParam"] = headerParam1
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
			httpDelivery["QueryParam"] = queryParam1
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
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_esa_site_delivery_task", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", response["SiteId"], response["TaskName"]))

	return resourceAliCloudEsaSiteDeliveryTaskUpdate(d, meta)
}

func resourceAliCloudEsaSiteDeliveryTaskRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	esaServiceV2 := EsaServiceV2{client}

	objectRaw, err := esaServiceV2.DescribeEsaSiteDeliveryTask(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_esa_site_delivery_task DescribeEsaSiteDeliveryTask Failed!!! %s", err)
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

	if !d.IsNewResource() && d.HasChange("discard_rate") {
		update = true
		request["DiscardRate"] = d.Get("discard_rate")
	}

	if !d.IsNewResource() && d.HasChange("business_type") {
		update = true
	}
	request["BusinessType"] = d.Get("business_type")
	if !d.IsNewResource() && d.HasChange("field_name") {
		update = true
	}
	request["FieldName"] = d.Get("field_name")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)
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
	action = "UpdateSiteDeliveryTaskStatus"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["SiteId"] = parts[0]
	query["TaskName"] = parts[1]

	if d.HasChange("status") {
		update = true
	}
	if v, ok := d.GetOk("status"); ok {
		query["Method"] = v.(string)
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcGet("ESA", "2024-09-10", action, query, request)
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
	return resourceAliCloudEsaSiteDeliveryTaskRead(d, meta)
}

func resourceAliCloudEsaSiteDeliveryTaskDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteSiteDeliveryTask"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["SiteId"] = parts[0]
	request["TaskName"] = parts[1]

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)
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
