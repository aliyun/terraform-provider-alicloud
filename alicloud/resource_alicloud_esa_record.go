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
				Type:     schema.TypeInt,
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
	request["RegionId"] = client.RegionId

	if v, ok := d.GetOk("comment"); ok {
		request["Comment"] = v
	}
	request["SiteId"] = d.Get("site_id")
	objectDataLocalMap := make(map[string]interface{})

	if v := d.Get("data"); v != nil {
		value1, _ := jsonpath.Get("$[0].value", v)
		if value1 != nil && value1 != "" {
			objectDataLocalMap["Value"] = value1
		}
		matchingType1, _ := jsonpath.Get("$[0].matching_type", v)
		if matchingType1 != nil && matchingType1 != "" {
			objectDataLocalMap["MatchingType"] = matchingType1
		}
		fingerprint1, _ := jsonpath.Get("$[0].fingerprint", v)
		if fingerprint1 != nil && fingerprint1 != "" {
			objectDataLocalMap["Fingerprint"] = fingerprint1
		}
		usage1, _ := jsonpath.Get("$[0].usage", v)
		if usage1 != nil && usage1 != "" {
			objectDataLocalMap["Usage"] = usage1
		}
		algorithm1, _ := jsonpath.Get("$[0].algorithm", v)
		if algorithm1 != nil && algorithm1 != "" {
			objectDataLocalMap["Algorithm"] = algorithm1
		}
		type1, _ := jsonpath.Get("$[0].type", v)
		if type1 != nil && type1 != "" {
			objectDataLocalMap["Type"] = type1
		}
		keyTag1, _ := jsonpath.Get("$[0].key_tag", v)
		if keyTag1 != nil && keyTag1 != "" {
			objectDataLocalMap["KeyTag"] = keyTag1
		}
		certificate1, _ := jsonpath.Get("$[0].certificate", v)
		if certificate1 != nil && certificate1 != "" {
			objectDataLocalMap["Certificate"] = certificate1
		}
		port1, _ := jsonpath.Get("$[0].port", v)
		if port1 != nil && port1 != "" {
			objectDataLocalMap["Port"] = port1
		}
		weight1, _ := jsonpath.Get("$[0].weight", v)
		if weight1 != nil && weight1 != "" {
			objectDataLocalMap["Weight"] = weight1
		}
		priority1, _ := jsonpath.Get("$[0].priority", v)
		if priority1 != nil && priority1 != "" {
			objectDataLocalMap["Priority"] = priority1
		}
		selector1, _ := jsonpath.Get("$[0].selector", v)
		if selector1 != nil && selector1 != "" {
			objectDataLocalMap["Selector"] = selector1
		}
		tag1, _ := jsonpath.Get("$[0].tag", v)
		if tag1 != nil && tag1 != "" {
			objectDataLocalMap["Tag"] = tag1
		}
		flag1, _ := jsonpath.Get("$[0].flag", v)
		if flag1 != nil && flag1 != "" {
			objectDataLocalMap["Flag"] = flag1
		}

		objectDataLocalMapJson, err := json.Marshal(objectDataLocalMap)
		if err != nil {
			return WrapError(err)
		}
		request["Data"] = string(objectDataLocalMapJson)
	}

	objectDataLocalMap1 := make(map[string]interface{})

	if v := d.Get("auth_conf"); !IsNil(v) {
		accessKey1, _ := jsonpath.Get("$[0].access_key", v)
		if accessKey1 != nil && accessKey1 != "" {
			objectDataLocalMap1["AccessKey"] = accessKey1
		}
		version1, _ := jsonpath.Get("$[0].version", v)
		if version1 != nil && version1 != "" {
			objectDataLocalMap1["Version"] = version1
		}
		authType1, _ := jsonpath.Get("$[0].auth_type", v)
		if authType1 != nil && authType1 != "" {
			objectDataLocalMap1["AuthType"] = authType1
		}
		secretKey1, _ := jsonpath.Get("$[0].secret_key", v)
		if secretKey1 != nil && secretKey1 != "" {
			objectDataLocalMap1["SecretKey"] = secretKey1
		}
		region1, _ := jsonpath.Get("$[0].region", v)
		if region1 != nil && region1 != "" {
			objectDataLocalMap1["Region"] = region1
		}

		objectDataLocalMap1Json, err := json.Marshal(objectDataLocalMap1)
		if err != nil {
			return WrapError(err)
		}
		request["AuthConf"] = string(objectDataLocalMap1Json)
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

	if objectRaw["BizName"] != nil {
		d.Set("biz_name", objectRaw["BizName"])
	}
	if objectRaw["Comment"] != nil {
		d.Set("comment", objectRaw["Comment"])
	}
	if objectRaw["CreateTime"] != nil {
		d.Set("create_time", objectRaw["CreateTime"])
	}
	if objectRaw["HostPolicy"] != nil {
		d.Set("host_policy", objectRaw["HostPolicy"])
	}
	if objectRaw["Proxied"] != nil {
		d.Set("proxied", objectRaw["Proxied"])
	}
	if objectRaw["RecordName"] != nil {
		d.Set("record_name", objectRaw["RecordName"])
	}
	if objectRaw["RecordType"] != nil {
		d.Set("record_type", objectRaw["RecordType"])
	}
	if objectRaw["SiteId"] != nil {
		d.Set("site_id", objectRaw["SiteId"])
	}
	if objectRaw["RecordSourceType"] != nil {
		d.Set("source_type", objectRaw["RecordSourceType"])
	}
	if objectRaw["Ttl"] != nil {
		d.Set("ttl", objectRaw["Ttl"])
	}

	authConfMaps := make([]map[string]interface{}, 0)
	authConfMap := make(map[string]interface{})
	authConf1Raw := make(map[string]interface{})
	if objectRaw["AuthConf"] != nil {
		authConf1Raw = objectRaw["AuthConf"].(map[string]interface{})
	}
	if len(authConf1Raw) > 0 {
		authConfMap["access_key"] = authConf1Raw["AccessKey"]
		authConfMap["auth_type"] = authConf1Raw["AuthType"]
		authConfMap["region"] = authConf1Raw["Region"]
		authConfMap["secret_key"] = authConf1Raw["SecretKey"]
		authConfMap["version"] = authConf1Raw["Version"]
		if v := d.Get("auth_conf"); !IsNil(v) {
			accessKey1, _ := jsonpath.Get("$[0].access_key", v)
			if accessKey1 != nil && accessKey1 != "" {
				authConfMap["access_key"] = accessKey1
			}
			secretKey1, _ := jsonpath.Get("$[0].secret_key", v)
			if secretKey1 != nil && secretKey1 != "" {
				authConfMap["secret_key"] = secretKey1
			}
		}

		authConfMaps = append(authConfMaps, authConfMap)
	}
	if objectRaw["AuthConf"] != nil {
		if err := d.Set("auth_conf", authConfMaps); err != nil {
			return err
		}
	}
	dataMaps := make([]map[string]interface{}, 0)
	dataMap := make(map[string]interface{})
	data1Raw := make(map[string]interface{})
	if objectRaw["Data"] != nil {
		data1Raw = objectRaw["Data"].(map[string]interface{})
	}
	if len(data1Raw) > 0 {
		dataMap["algorithm"] = data1Raw["Algorithm"]
		dataMap["certificate"] = data1Raw["Certificate"]
		dataMap["fingerprint"] = data1Raw["Fingerprint"]
		dataMap["flag"] = data1Raw["Flag"]
		dataMap["key_tag"] = data1Raw["KeyTag"]
		dataMap["matching_type"] = data1Raw["MatchingType"]
		dataMap["port"] = data1Raw["Port"]
		dataMap["priority"] = data1Raw["Priority"]
		dataMap["selector"] = data1Raw["Selector"]
		dataMap["tag"] = data1Raw["Tag"]
		dataMap["type"] = data1Raw["Type"]
		dataMap["usage"] = data1Raw["Usage"]
		dataMap["value"] = data1Raw["Value"]
		dataMap["weight"] = data1Raw["Weight"]

		dataMaps = append(dataMaps, dataMap)
	}
	if objectRaw["Data"] != nil {
		if err := d.Set("data", dataMaps); err != nil {
			return err
		}
	}

	return nil
}

func resourceAliCloudEsaRecordUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	action := "UpdateRecord"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["RecordId"] = d.Id()
	request["RegionId"] = client.RegionId
	if d.HasChange("comment") {
		update = true
		request["Comment"] = d.Get("comment")
	}

	if d.HasChange("data") {
		update = true
	}
	objectDataLocalMap := make(map[string]interface{})

	if v := d.Get("data"); v != nil {
		value1, _ := jsonpath.Get("$[0].value", v)
		if value1 != nil && (d.HasChange("data.0.value") || value1 != "") {
			objectDataLocalMap["Value"] = value1
		}
		matchingType1, _ := jsonpath.Get("$[0].matching_type", v)
		if matchingType1 != nil && (d.HasChange("data.0.matching_type") || matchingType1 != "") {
			objectDataLocalMap["MatchingType"] = matchingType1
		}
		fingerprint1, _ := jsonpath.Get("$[0].fingerprint", v)
		if fingerprint1 != nil && (d.HasChange("data.0.fingerprint") || fingerprint1 != "") {
			objectDataLocalMap["Fingerprint"] = fingerprint1
		}
		usage1, _ := jsonpath.Get("$[0].usage", v)
		if usage1 != nil && (d.HasChange("data.0.usage") || usage1 != "") {
			objectDataLocalMap["Usage"] = usage1
		}
		algorithm1, _ := jsonpath.Get("$[0].algorithm", v)
		if algorithm1 != nil && (d.HasChange("data.0.algorithm") || algorithm1 != "") {
			objectDataLocalMap["Algorithm"] = algorithm1
		}
		type1, _ := jsonpath.Get("$[0].type", v)
		if type1 != nil && (d.HasChange("data.0.type") || type1 != "") {
			objectDataLocalMap["Type"] = type1
		}
		keyTag1, _ := jsonpath.Get("$[0].key_tag", v)
		if keyTag1 != nil && (d.HasChange("data.0.key_tag") || keyTag1 != "") {
			objectDataLocalMap["KeyTag"] = keyTag1
		}
		certificate1, _ := jsonpath.Get("$[0].certificate", v)
		if certificate1 != nil && (d.HasChange("data.0.certificate") || certificate1 != "") {
			objectDataLocalMap["Certificate"] = certificate1
		}
		port1, _ := jsonpath.Get("$[0].port", v)
		if port1 != nil && (d.HasChange("data.0.port") || port1 != "") {
			objectDataLocalMap["Port"] = port1
		}
		weight1, _ := jsonpath.Get("$[0].weight", v)
		if weight1 != nil && (d.HasChange("data.0.weight") || weight1 != "") {
			objectDataLocalMap["Weight"] = weight1
		}
		priority1, _ := jsonpath.Get("$[0].priority", v)
		if priority1 != nil && (d.HasChange("data.0.priority") || priority1 != "") {
			objectDataLocalMap["Priority"] = priority1
		}
		selector1, _ := jsonpath.Get("$[0].selector", v)
		if selector1 != nil && (d.HasChange("data.0.selector") || selector1 != "") {
			objectDataLocalMap["Selector"] = selector1
		}
		tag1, _ := jsonpath.Get("$[0].tag", v)
		if tag1 != nil && (d.HasChange("data.0.tag") || tag1 != "") {
			objectDataLocalMap["Tag"] = tag1
		}
		flag1, _ := jsonpath.Get("$[0].flag", v)
		if flag1 != nil && (d.HasChange("data.0.flag") || flag1 != "") {
			objectDataLocalMap["Flag"] = flag1
		}

		objectDataLocalMapJson, err := json.Marshal(objectDataLocalMap)
		if err != nil {
			return WrapError(err)
		}
		request["Data"] = string(objectDataLocalMapJson)
	}

	if d.HasChange("auth_conf") {
		update = true
		objectDataLocalMap1 := make(map[string]interface{})

		if v := d.Get("auth_conf"); v != nil {
			accessKey1, _ := jsonpath.Get("$[0].access_key", v)
			if accessKey1 != nil && (d.HasChange("auth_conf.0.access_key") || accessKey1 != "") {
				objectDataLocalMap1["AccessKey"] = accessKey1
			}
			version1, _ := jsonpath.Get("$[0].version", v)
			if version1 != nil && (d.HasChange("auth_conf.0.version") || version1 != "") {
				objectDataLocalMap1["Version"] = version1
			}
			authType1, _ := jsonpath.Get("$[0].auth_type", v)
			if authType1 != nil && (d.HasChange("auth_conf.0.auth_type") || authType1 != "") {
				objectDataLocalMap1["AuthType"] = authType1
			}
			secretKey1, _ := jsonpath.Get("$[0].secret_key", v)
			if secretKey1 != nil && (d.HasChange("auth_conf.0.secret_key") || secretKey1 != "") {
				objectDataLocalMap1["SecretKey"] = secretKey1
			}
			region1, _ := jsonpath.Get("$[0].region", v)
			if region1 != nil && (d.HasChange("auth_conf.0.region") || region1 != "") {
				objectDataLocalMap1["Region"] = region1
			}

			objectDataLocalMap1Json, err := json.Marshal(objectDataLocalMap1)
			if err != nil {
				return WrapError(err)
			}
			request["AuthConf"] = string(objectDataLocalMap1Json)
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
				if IsExpectedErrors(err, []string{"Record.ServiceBusy", "TooManyRequests"}) || NeedRetry(err) {
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
	request["RegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)

		if err != nil {
			if IsExpectedErrors(err, []string{"Record.ServiceBusy", "TooManyRequests"}) || NeedRetry(err) {
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
