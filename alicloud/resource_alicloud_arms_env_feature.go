// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudArmsEnvFeature() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudArmsEnvFeatureCreate,
		Read:   resourceAliCloudArmsEnvFeatureRead,
		Update: resourceAliCloudArmsEnvFeatureUpdate,
		Delete: resourceAliCloudArmsEnvFeatureDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"env_feature_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"environment_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"feature_version": {
				Type:     schema.TypeString,
				Required: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudArmsEnvFeatureCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "InstallEnvironmentFeature"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["EnvironmentId"] = d.Get("environment_id")
	query["FeatureName"] = d.Get("env_feature_name")
	query["RegionId"] = client.RegionId

	query["FeatureVersion"] = d.Get("feature_version")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("ARMS", "2019-08-08", action, query, request, true)

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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_arms_env_feature", action, AlibabaCloudSdkGoERROR)
	}
	code, _ := jsonpath.Get("$.Code", response)
	if fmt.Sprint(code) != "200" {
		log.Printf("[DEBUG] Resource alicloud_arms_env_feature InstallEnvironmentFeature Failed!!! %s", response)
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_arms_env_feature", action, AlibabaCloudSdkGoERROR, response)
	}

	d.SetId(fmt.Sprintf("%v:%v", query["EnvironmentId"], query["FeatureName"]))

	armsServiceV2 := ArmsServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Success"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, armsServiceV2.ArmsEnvFeatureStateRefreshFunc(d.Id(), "$.Feature.Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudArmsEnvFeatureRead(d, meta)
}

func resourceAliCloudArmsEnvFeatureRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	armsServiceV2 := ArmsServiceV2{client}

	objectRaw, err := armsServiceV2.DescribeArmsEnvFeature(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_arms_env_feature DescribeArmsEnvFeature Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	feature1RawObj, _ := jsonpath.Get("$.Feature", objectRaw)
	feature1Raw := make(map[string]interface{})
	if feature1RawObj != nil {
		feature1Raw = feature1RawObj.(map[string]interface{})
	}
	d.Set("feature_version", feature1Raw["Version"])
	d.Set("status", feature1Raw["Status"])
	d.Set("env_feature_name", feature1Raw["Name"])
	d.Set("environment_id", feature1Raw["EnvironmentId"])
	featureStatus1RawObj, _ := jsonpath.Get("$.FeatureStatus", objectRaw)
	featureStatus1Raw := make(map[string]interface{})
	if featureStatus1RawObj != nil {
		featureStatus1Raw = featureStatus1RawObj.(map[string]interface{})
	}
	d.Set("namespace", featureStatus1Raw["Namespace"])

	return nil
}

func resourceAliCloudArmsEnvFeatureUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	parts := strings.Split(d.Id(), ":")
	action := "UpgradeEnvironmentFeature"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["EnvironmentId"] = parts[0]
	query["FeatureName"] = parts[1]
	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("feature_version") {
		update = true
	}
	request["FeatureVersion"] = d.Get("feature_version")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("ARMS", "2019-08-08", action, query, request, true)

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
		armsServiceV2 := ArmsServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{fmt.Sprint(d.Get("feature_version"))}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, armsServiceV2.ArmsEnvFeatureStateRefreshFunc(d.Id(), "$.Feature.Version", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return resourceAliCloudArmsEnvFeatureRead(d, meta)
}

func resourceAliCloudArmsEnvFeatureDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteEnvironmentFeature"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["EnvironmentId"] = parts[0]
	query["FeatureName"] = parts[1]
	request["RegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("ARMS", "2019-08-08", action, query, request, true)

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

	armsServiceV2 := ArmsServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{""}, d.Timeout(schema.TimeoutDelete), 5*time.Second, armsServiceV2.ArmsEnvFeatureStateRefreshFunc(d.Id(), "$.Feature.Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
