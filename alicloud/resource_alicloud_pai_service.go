package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/blues/jsonata-go"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
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
					equal, _ := compareJsonTemplateAreEquivalent(old, new)
					return equal
				},
				StateFunc: func(v interface{}) string {
					jsonString, _ := normalizeJsonString(v)
					return jsonString
				},
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tags": tagsSchema(),
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

	request["body"] = convertJsonStringToObject(d.Get("service_config"))
	if v, ok := d.GetOk("develop"); ok {
		query["Develop"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("tags"); ok {
		query["Labels"] = StringPointer(convertMapToJsonStringIgnoreError(v.(map[string]interface{})))
	}

	if v, ok := d.GetOk("workspace_id"); ok {
		query["WorkspaceId"] = StringPointer(v.(string))
	}

	body = request["body"].(map[string]interface{})
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RoaPost("EAS", "2021-07-01", action, query, nil, body, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_pai_service", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["ServiceName"]))

	paiServiceV2 := PaiServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, paiServiceV2.PaiServiceStateRefreshFunc(d.Id(), "Status", []string{"Failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudPaiServiceUpdate(d, meta)
}

func resourceAliCloudPaiServiceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	paiServiceV2 := PaiServiceV2{client}

	objectRaw, err := paiServiceV2.DescribePaiService(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_pai_service DescribePaiService Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("region_id", objectRaw["Region"])
	d.Set("service_uid", objectRaw["ServiceUid"])
	d.Set("status", objectRaw["Status"])
	d.Set("workspace_id", objectRaw["WorkspaceId"])

	e := jsonata.MustCompile("$merge($map($.Labels, function($v, $k) {{$lookup($v, \"LabelKey\"): $lookup($v, \"LabelValue\")}}))")
	evaluation, _ := e.Eval(objectRaw)
	d.Set("tags", evaluation)
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
			return WrapError(err)
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
				query["ClusterId"] = StringPointer(client.RegionId)
				body = request
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RoaPut("EAS", "2021-07-01", action, query, nil, body, true)
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
				paiServiceV2 := PaiServiceV2{client}
				stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, paiServiceV2.PaiServiceStateRefreshFunc(d.Id(), "Status", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
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
				query["ClusterId"] = StringPointer(client.RegionId)
				body = request
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RoaPut("EAS", "2021-07-01", action, query, nil, body, true)
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
				paiServiceV2 := PaiServiceV2{client}
				stateConf := BuildStateConf([]string{}, []string{"Stopped"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, paiServiceV2.PaiServiceStateRefreshFunc(d.Id(), "Status", []string{"Failed"}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
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
	query["ClusterId"] = StringPointer(client.RegionId)
	if !d.IsNewResource() && d.HasChange("service_config") {
		update = true
	}
	request["body"] = convertJsonStringToObject(d.Get("service_config"))
	body = request["body"].(map[string]interface{})
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPut("EAS", "2021-07-01", action, query, nil, body, true)
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
		paiServiceV2 := PaiServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, paiServiceV2.PaiServiceStateRefreshFunc(d.Id(), "Status", []string{"Failed"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	if d.HasChange("tags") {
		paiServiceV2 := PaiServiceV2{client}
		if err := paiServiceV2.SetResourceTags(d, "service"); err != nil {
			return WrapError(err)
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
	query["ClusterId"] = StringPointer(client.RegionId)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RoaDelete("EAS", "2021-07-01", action, query, nil, nil, true)

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
		if IsExpectedErrors(err, []string{"InvalidService.NotFound"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	paiServiceV2 := PaiServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, paiServiceV2.PaiServiceStateRefreshFunc(d.Id(), "$.ServiceName", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
