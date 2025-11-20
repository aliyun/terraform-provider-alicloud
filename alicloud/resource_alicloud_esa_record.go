// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudEsaRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEsaRecordCreate,
		Read:   resourceAliCloudEsaRecordRead,
		Update: resourceAliCloudEsaRecordUpdate,
		Delete: resourceAliCloudEsaRecordDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"auth_conf": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"secret_key": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},
						"version": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"region": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"auth_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: StringInSlice([]string{"public", "private", "private_same_account"}, false),
						},
						"access_key": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},
					},
				},
			},
			"biz_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"comment": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"data": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"usage": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"fingerprint": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"priority": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"flag": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: IntBetween(0, 255),
						},
						"algorithm": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"weight": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"type": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"matching_type": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"key_tag": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"tag": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"certificate": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"selector": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"host_policy": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"proxied": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"record_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"record_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"site_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"source_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ttl": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func resourceAliCloudEsaRecordCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateRecord"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	if v, ok := d.GetOk("comment"); ok {
		request["Comment"] = v
	}
	request["SiteId"] = d.Get("site_id")
	data := make(map[string]interface{})

	if v := d.Get("data"); v != nil {
		value1, _ := jsonpath.Get("$[0].value", v)
		if value1 != nil && value1 != "" {
			data["Value"] = value1
		}
		matchingType1, _ := jsonpath.Get("$[0].matching_type", v)
		if matchingType1 != nil && matchingType1 != "" {
			data["MatchingType"] = matchingType1
		}
		fingerprint1, _ := jsonpath.Get("$[0].fingerprint", v)
		if fingerprint1 != nil && fingerprint1 != "" {
			data["Fingerprint"] = fingerprint1
		}
		usage1, _ := jsonpath.Get("$[0].usage", v)
		if usage1 != nil && usage1 != "" {
			data["Usage"] = usage1
		}
		algorithm1, _ := jsonpath.Get("$[0].algorithm", v)
		if algorithm1 != nil && algorithm1 != "" {
			data["Algorithm"] = algorithm1
		}
		type1, _ := jsonpath.Get("$[0].type", v)
		if type1 != nil && type1 != "" {
			data["Type"] = type1
		}
		keyTag1, _ := jsonpath.Get("$[0].key_tag", v)
		if keyTag1 != nil && keyTag1 != "" {
			data["KeyTag"] = keyTag1
		}
		certificate1, _ := jsonpath.Get("$[0].certificate", v)
		if certificate1 != nil && certificate1 != "" {
			data["Certificate"] = certificate1
		}
		port1, _ := jsonpath.Get("$[0].port", v)
		if port1 != nil && port1 != "" {
			data["Port"] = port1
		}
		weight1, _ := jsonpath.Get("$[0].weight", v)
		if weight1 != nil && weight1 != "" {
			data["Weight"] = weight1
		}
		priority1, _ := jsonpath.Get("$[0].priority", v)
		if priority1 != nil && priority1 != "" {
			data["Priority"] = priority1
		}
		selector1, _ := jsonpath.Get("$[0].selector", v)
		if selector1 != nil && selector1 != "" {
			data["Selector"] = selector1
		}
		tag1, _ := jsonpath.Get("$[0].tag", v)
		if tag1 != nil && tag1 != "" {
			data["Tag"] = tag1
		}
		flag1, _ := jsonpath.Get("$[0].flag", v)
		if flag1 != nil && flag1 != "" {
			data["Flag"] = flag1
		}

		dataJson, err := json.Marshal(data)
		if err != nil {
			return WrapError(err)
		}
		request["Data"] = string(dataJson)
	}

	authConf := make(map[string]interface{})

	if v := d.Get("auth_conf"); !IsNil(v) {
		accessKey1, _ := jsonpath.Get("$[0].access_key", v)
		if accessKey1 != nil && accessKey1 != "" {
			authConf["AccessKey"] = accessKey1
		}
		version1, _ := jsonpath.Get("$[0].version", v)
		if version1 != nil && version1 != "" {
			authConf["Version"] = version1
		}
		authType1, _ := jsonpath.Get("$[0].auth_type", v)
		if authType1 != nil && authType1 != "" {
			authConf["AuthType"] = authType1
		}
		secretKey1, _ := jsonpath.Get("$[0].secret_key", v)
		if secretKey1 != nil && secretKey1 != "" {
			authConf["SecretKey"] = secretKey1
		}
		region1, _ := jsonpath.Get("$[0].region", v)
		if region1 != nil && region1 != "" {
			authConf["Region"] = region1
		}

		authConfJson, err := json.Marshal(authConf)
		if err != nil {
			return WrapError(err)
		}
		request["AuthConf"] = string(authConfJson)
	}

	if v, ok := d.GetOkExists("ttl"); ok {
		request["Ttl"] = v
	}
	request["RecordName"] = d.Get("record_name")
	if v, ok := d.GetOk("host_policy"); ok {
		request["HostPolicy"] = v
	}
	request["Type"] = d.Get("record_type")
	if v, ok := d.GetOk("biz_name"); ok {
		request["BizName"] = v
	}
	if v, ok := d.GetOkExists("proxied"); ok {
		request["Proxied"] = v
	}
	if v, ok := d.GetOk("source_type"); ok {
		request["SourceType"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"Site.ServiceBusy", "TooManyRequests"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_esa_record", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["RecordId"]))

	return resourceAliCloudEsaRecordRead(d, meta)
}

func resourceAliCloudEsaRecordRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	esaServiceV2 := EsaServiceV2{client}

	objectRaw, err := esaServiceV2.DescribeEsaRecord(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_esa_record DescribeEsaRecord Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("biz_name", objectRaw["BizName"])
	d.Set("comment", objectRaw["Comment"])
	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("host_policy", objectRaw["HostPolicy"])
	d.Set("proxied", objectRaw["Proxied"])
	d.Set("record_name", objectRaw["RecordName"])
	d.Set("record_type", objectRaw["RecordType"])
	if v, ok := objectRaw["SiteId"]; ok {
		d.Set("site_id", v)
	}

	d.Set("source_type", objectRaw["RecordSourceType"])
	d.Set("ttl", objectRaw["Ttl"])

	authConfMaps := make([]map[string]interface{}, 0)
	authConfMap := make(map[string]interface{})
	authConfRaw := make(map[string]interface{})
	if objectRaw["AuthConf"] != nil {
		authConfRaw = objectRaw["AuthConf"].(map[string]interface{})
	}
	if len(authConfRaw) > 0 {
		authConfMap["access_key"] = authConfRaw["AccessKey"]
		authConfMap["auth_type"] = authConfRaw["AuthType"]
		authConfMap["region"] = authConfRaw["Region"]
		authConfMap["secret_key"] = authConfRaw["SecretKey"]
		authConfMap["version"] = authConfRaw["Version"]

		authConfMaps = append(authConfMaps, authConfMap)
	}
	if err := d.Set("auth_conf", authConfMaps); err != nil {
		return err
	}
	dataMaps := make([]map[string]interface{}, 0)
	dataMap := make(map[string]interface{})
	dataRaw := make(map[string]interface{})
	if objectRaw["Data"] != nil {
		dataRaw = objectRaw["Data"].(map[string]interface{})
	}
	if len(dataRaw) > 0 {
		dataMap["algorithm"] = dataRaw["Algorithm"]
		dataMap["certificate"] = dataRaw["Certificate"]
		dataMap["fingerprint"] = dataRaw["Fingerprint"]
		dataMap["flag"] = dataRaw["Flag"]
		dataMap["key_tag"] = dataRaw["KeyTag"]
		dataMap["matching_type"] = dataRaw["MatchingType"]
		dataMap["port"] = dataRaw["Port"]
		dataMap["priority"] = dataRaw["Priority"]
		dataMap["selector"] = dataRaw["Selector"]
		dataMap["tag"] = dataRaw["Tag"]
		dataMap["type"] = dataRaw["Type"]
		dataMap["usage"] = dataRaw["Usage"]
		dataMap["value"] = dataRaw["Value"]
		dataMap["weight"] = dataRaw["Weight"]

		dataMaps = append(dataMaps, dataMap)
	}
	if err := d.Set("data", dataMaps); err != nil {
		return err
	}

	return nil
}

func resourceAliCloudEsaRecordUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	action := "UpdateRecord"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["RecordId"] = d.Id()

	if d.HasChange("comment") {
		update = true
		request["Comment"] = d.Get("comment")
	}

	if d.HasChange("data") {
		update = true
	}
	data := make(map[string]interface{})

	if v := d.Get("data"); v != nil {
		value1, _ := jsonpath.Get("$[0].value", v)
		if value1 != nil && (d.HasChange("data.0.value") || value1 != "") {
			data["Value"] = value1
		}
		matchingType1, _ := jsonpath.Get("$[0].matching_type", v)
		if matchingType1 != nil && (d.HasChange("data.0.matching_type") || matchingType1 != "") {
			data["MatchingType"] = matchingType1
		}
		fingerprint1, _ := jsonpath.Get("$[0].fingerprint", v)
		if fingerprint1 != nil && (d.HasChange("data.0.fingerprint") || fingerprint1 != "") {
			data["Fingerprint"] = fingerprint1
		}
		usage1, _ := jsonpath.Get("$[0].usage", v)
		if usage1 != nil && (d.HasChange("data.0.usage") || usage1 != "") {
			data["Usage"] = usage1
		}
		algorithm1, _ := jsonpath.Get("$[0].algorithm", v)
		if algorithm1 != nil && (d.HasChange("data.0.algorithm") || algorithm1 != "") {
			data["Algorithm"] = algorithm1
		}
		type1, _ := jsonpath.Get("$[0].type", v)
		if type1 != nil && (d.HasChange("data.0.type") || type1 != "") {
			data["Type"] = type1
		}
		keyTag1, _ := jsonpath.Get("$[0].key_tag", v)
		if keyTag1 != nil && (d.HasChange("data.0.key_tag") || keyTag1 != "") {
			data["KeyTag"] = keyTag1
		}
		certificate1, _ := jsonpath.Get("$[0].certificate", v)
		if certificate1 != nil && (d.HasChange("data.0.certificate") || certificate1 != "") {
			data["Certificate"] = certificate1
		}
		port1, _ := jsonpath.Get("$[0].port", v)
		if port1 != nil && (d.HasChange("data.0.port") || port1 != "") {
			data["Port"] = port1
		}
		weight1, _ := jsonpath.Get("$[0].weight", v)
		if weight1 != nil && (d.HasChange("data.0.weight") || weight1 != "") {
			data["Weight"] = weight1
		}
		priority1, _ := jsonpath.Get("$[0].priority", v)
		if priority1 != nil && (d.HasChange("data.0.priority") || priority1 != "") {
			data["Priority"] = priority1
		}
		selector1, _ := jsonpath.Get("$[0].selector", v)
		if selector1 != nil && (d.HasChange("data.0.selector") || selector1 != "") {
			data["Selector"] = selector1
		}
		tag1, _ := jsonpath.Get("$[0].tag", v)
		if tag1 != nil && (d.HasChange("data.0.tag") || tag1 != "") {
			data["Tag"] = tag1
		}
		flag1, _ := jsonpath.Get("$[0].flag", v)
		if flag1 != nil && (d.HasChange("data.0.flag") || flag1 != "") {
			data["Flag"] = flag1
		}

		dataJson, err := json.Marshal(data)
		if err != nil {
			return WrapError(err)
		}
		request["Data"] = string(dataJson)
	}

	if d.HasChange("auth_conf") {
		update = true
		authConf := make(map[string]interface{})

		if v := d.Get("auth_conf"); v != nil {
			accessKey1, _ := jsonpath.Get("$[0].access_key", v)
			if accessKey1 != nil && (d.HasChange("auth_conf.0.access_key") || accessKey1 != "") {
				authConf["AccessKey"] = accessKey1
			}
			version1, _ := jsonpath.Get("$[0].version", v)
			if version1 != nil && (d.HasChange("auth_conf.0.version") || version1 != "") {
				authConf["Version"] = version1
			}
			authType1, _ := jsonpath.Get("$[0].auth_type", v)
			if authType1 != nil && (d.HasChange("auth_conf.0.auth_type") || authType1 != "") {
				authConf["AuthType"] = authType1
			}
			secretKey1, _ := jsonpath.Get("$[0].secret_key", v)
			if secretKey1 != nil && (d.HasChange("auth_conf.0.secret_key") || secretKey1 != "") {
				authConf["SecretKey"] = secretKey1
			}
			region1, _ := jsonpath.Get("$[0].region", v)
			if region1 != nil && (d.HasChange("auth_conf.0.region") || region1 != "") {
				authConf["Region"] = region1
			}

			authConfJson, err := json.Marshal(authConf)
			if err != nil {
				return WrapError(err)
			}
			request["AuthConf"] = string(authConfJson)
		}
	}

	if d.HasChange("ttl") {
		update = true
		request["Ttl"] = d.Get("ttl")
	}

	if d.HasChange("host_policy") {
		update = true
		request["HostPolicy"] = d.Get("host_policy")
	}

	if d.HasChange("biz_name") {
		update = true
		request["BizName"] = d.Get("biz_name")
	}

	if d.HasChange("proxied") {
		update = true
		request["Proxied"] = d.Get("proxied")
	}

	if d.HasChange("source_type") {
		update = true
		request["SourceType"] = d.Get("source_type")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"TooManyRequests", "Record.ServiceBusy"}) || NeedRetry(err) {
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

	return resourceAliCloudEsaRecordRead(d, meta)
}

func resourceAliCloudEsaRecordDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteRecord"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RecordId"] = d.Id()

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"TooManyRequests", "Record.ServiceBusy"}) || NeedRetry(err) {
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
