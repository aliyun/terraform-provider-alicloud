// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package threatdetection

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAliCloudThreatDetectionSasTrail() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudThreatDetectionSasTrailCreate,
		Read:   resourceAliCloudThreatDetectionSasTrailRead,
		Delete: resourceAliCloudThreatDetectionSasTrailDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"service_trail": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"config": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"update_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceAliCloudThreatDetectionSasTrailCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateServiceTrail"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("Sas", "2018-12-03", action, query, request, true)

		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_threat_detection_sas_trail", action, errs.AlibabaCloudSdkGoERROR)
	}

	accountId, err := client.AccountId()
	d.SetId(fmt.Sprint(accountId))

	return resourceAliCloudThreatDetectionSasTrailRead(d, meta)
}

func resourceAliCloudThreatDetectionSasTrailRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	threatDetectionServiceV2 := ServiceV2{client}

	objectRaw, err := threatDetectionServiceV2.DescribeThreatDetectionSasTrail(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_threat_detection_sas_trail DescribeThreatDetectionSasTrail Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("create_time", objectRaw["CreateTime"])
	serviceTrailMaps := make([]map[string]interface{}, 0)
	serviceTrailMap := make(map[string]interface{})

	serviceTrailMap["config"] = objectRaw["Config"]
	serviceTrailMap["update_time"] = objectRaw["UpdateTime"]
	serviceTrailMaps = append(serviceTrailMaps, serviceTrailMap)
	err = d.Set("service_trail", serviceTrailMaps)
	if err != nil {
		return err
	}

	return nil
}

func resourceAliCloudThreatDetectionSasTrailDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteServiceTrail"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("Sas", "2018-12-03", action, query, request, true)

		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"ServerError"}) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	return nil
}
