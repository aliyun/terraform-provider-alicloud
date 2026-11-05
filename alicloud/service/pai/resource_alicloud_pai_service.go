package pai

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/blues/jsonata-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAliCloudPaiService() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudPaiServiceCreate,
		Read:   resourceAliCloudPaiServiceRead,
		Update: resourceAliCloudPaiServiceUpdate,
		Delete: resourceAliCloudPaiServiceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(16 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"develop": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_config": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsJSON,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := helper.CompareJsonTemplateAreEquivalent(old, new)
					return equal
				},
				StateFunc: func(v interface{}) string {
					jsonString, _ := validators.NormalizeJsonString(v)
					return jsonString
				},
			},
			"service_uid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tags": tags.TagsSchema(),
			"workspace_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAliCloudPaiServiceCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/api/v2/services")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["body"] = convert.ConvertJsonStringToObject(d.Get("service_config"))
	if v, ok := d.GetOk("develop"); ok {
		query["Develop"] = helper.StringPointer(v.(string))
	}

	if v, ok := d.GetOk("tags"); ok {
		query["Labels"] = helper.StringPointer(convert.ConvertMapToJsonStringIgnoreError(v.(map[string]interface{})))
	}

	if v, ok := d.GetOk("workspace_id"); ok {
		query["WorkspaceId"] = helper.StringPointer(v.(string))
	}

	body = request["body"].(map[string]interface{})
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RoaPost("EAS", "2021-07-01", action, query, nil, body, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_pai_service", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["ServiceName"]))

	paiServiceV2 := PaiServiceV2{client}
	stateConf := helper.BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, paiServiceV2.PaiServiceStateRefreshFunc(d.Id(), "Status", []string{"Failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return resourceAliCloudPaiServiceUpdate(d, meta)
}

func resourceAliCloudPaiServiceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	paiServiceV2 := PaiServiceV2{client}

	objectRaw, err := paiServiceV2.DescribePaiService(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_pai_service DescribePaiService Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("region_id", objectRaw["Region"])
	d.Set("service_uid", objectRaw["ServiceUid"])
	d.Set("status", objectRaw["Status"])
	d.Set("workspace_id", objectRaw["WorkspaceId"])

	e := jsonata.MustCompile("$merge($map($.Labels, function($v, $k) {{$lookup($v, \"LabelKey\"): $lookup($v, \"LabelValue\")}}))")
	evaluation, _ := e.Eval(objectRaw)
	d.Set("tags", tags.TagsToMap(evaluation))
	e = jsonata.MustCompile("$.ServiceConfig")
	evaluation, _ = e.Eval(objectRaw)
	d.Set("service_config", evaluation)

	return nil
}

func resourceAliCloudPaiServiceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false

	if d.HasChange("status") {
		paiServiceV2 := PaiServiceV2{client}
		object, err := paiServiceV2.DescribePaiService(d.Id())
		if err != nil {
			return errs.WrapError(err)
		}

		target := d.Get("status").(string)
		if object["Status"].(string) != target {
			if target == "Running" {
				ServiceName := d.Id()
				ClusterId := client.RegionId
				action := fmt.Sprintf("/api/v2/services/%s/%s/start", ClusterId, ServiceName)
				request = make(map[string]interface{})
				query = make(map[string]*string)
				body = make(map[string]interface{})
				request["ServiceName"] = d.Id()
				query["ClusterId"] = helper.StringPointer(client.RegionId)
				body = request
				wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
				err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
					response, err = client.RoaPut("EAS", "2021-07-01", action, query, nil, body, true)
					if err != nil {
						if helper.NeedRetry(err) {
							wait()
							return retry.RetryableError(err)
						}
						return retry.NonRetryableError(err)
					}
					return nil
				})
				logs.AddDebug(action, response, request)
				if err != nil {
					return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
				}
				paiServiceV2 := PaiServiceV2{client}
				stateConf := helper.BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, paiServiceV2.PaiServiceStateRefreshFunc(d.Id(), "Status", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return errs.WrapErrorf(err, errs.IdMsg, d.Id())
				}

			}
			if target == "Stopped" {
				ServiceName := d.Id()
				ClusterId := client.RegionId
				action := fmt.Sprintf("/api/v2/services/%s/%s/stop", ClusterId, ServiceName)
				request = make(map[string]interface{})
				query = make(map[string]*string)
				body = make(map[string]interface{})
				request["ServiceName"] = d.Id()
				query["ClusterId"] = helper.StringPointer(client.RegionId)
				body = request
				wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
				err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
					response, err = client.RoaPut("EAS", "2021-07-01", action, query, nil, body, true)
					if err != nil {
						if helper.NeedRetry(err) {
							wait()
							return retry.RetryableError(err)
						}
						return retry.NonRetryableError(err)
					}
					return nil
				})
				logs.AddDebug(action, response, request)
				if err != nil {
					return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
				}
				paiServiceV2 := PaiServiceV2{client}
				stateConf := helper.BuildStateConf([]string{}, []string{"Stopped"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, paiServiceV2.PaiServiceStateRefreshFunc(d.Id(), "Status", []string{"Failed"}))
				if _, err := stateConf.WaitForState(); err != nil {
					return errs.WrapErrorf(err, errs.IdMsg, d.Id())
				}

			}
		}
	}

	ServiceName := d.Id()
	ClusterId := client.RegionId
	action := fmt.Sprintf("/api/v2/services/%s/%s", ClusterId, ServiceName)
	var err error
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	request["ServiceName"] = d.Id()
	query["ClusterId"] = helper.StringPointer(client.RegionId)
	if !d.IsNewResource() && d.HasChange("service_config") {
		update = true
	}
	request["body"] = convert.ConvertJsonStringToObject(d.Get("service_config"))
	body = request["body"].(map[string]interface{})
	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RoaPut("EAS", "2021-07-01", action, query, nil, body, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
		paiServiceV2 := PaiServiceV2{client}
		stateConf := helper.BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, paiServiceV2.PaiServiceStateRefreshFunc(d.Id(), "Status", []string{"Failed"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}
	}

	if d.HasChange("tags") {
		paiServiceV2 := PaiServiceV2{client}
		if err := paiServiceV2.SetResourceTags(d, "service"); err != nil {
			return errs.WrapError(err)
		}
	}
	return resourceAliCloudPaiServiceRead(d, meta)
}

func resourceAliCloudPaiServiceDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	ServiceName := d.Id()
	ClusterId := client.RegionId
	action := fmt.Sprintf("/api/v2/services/%s/%s", ClusterId, ServiceName)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	request["ServiceName"] = d.Id()
	query["ClusterId"] = helper.StringPointer(client.RegionId)

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RoaDelete("EAS", "2021-07-01", action, query, nil, nil, true)

		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidService.NotFound"}) || errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	paiServiceV2 := PaiServiceV2{client}
	stateConf := helper.BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, paiServiceV2.PaiServiceStateRefreshFunc(d.Id(), "$.ServiceName", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return nil
}
