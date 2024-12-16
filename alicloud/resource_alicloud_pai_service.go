package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
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
			Create: schema.DefaultTimeout(15 * time.Minute),
			Update: schema.DefaultTimeout(15 * time.Minute),
			Delete: schema.DefaultTimeout(15 * time.Minute),
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
			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_config": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.ValidateJsonString,
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
	conn, err := client.NewPaiClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})

	request["body"] = convertJsonStringToObject(d.Get("service_config"))
	if v, ok := d.GetOk("develop"); ok {
		query["Develop"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("labels"); ok {
		query["Labels"] = StringPointer(convertMapToJsonStringIgnoreError(v.(map[string]interface{})))
	}

	if v, ok := d.GetOk("workspace_id"); ok {
		query["WorkspaceId"] = StringPointer(v.(string))
	}

	body = request["body"].(map[string]interface{})
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer("2021-07-01"), nil, StringPointer("POST"), StringPointer("AK"), StringPointer(action), query, nil, body, &runtime)
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

	id, _ := jsonpath.Get("$.body.ServiceName", response)
	d.SetId(fmt.Sprint(id))

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

	if objectRaw["CreateTime"] != nil {
		d.Set("create_time", objectRaw["CreateTime"])
	}
	if objectRaw["Region"] != nil {
		d.Set("region_id", objectRaw["Region"])
	}
	if objectRaw["Status"] != nil {
		d.Set("status", objectRaw["Status"])
	}
	if objectRaw["WorkspaceId"] != nil {
		d.Set("workspace_id", objectRaw["WorkspaceId"])
	}

	e := jsonata.MustCompile("$merge($map($.Labels, function($v, $k) {{$lookup($v, \"LabelKey\"): $lookup($v, \"LabelValue\")}}))")
	evaluation, _ := e.Eval(objectRaw)
	d.Set("labels", evaluation)
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
				conn, err := client.NewPaiClient()
				if err != nil {
					return WrapError(err)
				}
				request = make(map[string]interface{})
				query = make(map[string]*string)
				body = make(map[string]interface{})
				request["ServiceName"] = d.Id()
				query["ClusterId"] = StringPointer(client.RegionId)
				body = request
				runtime := util.RuntimeOptions{}
				runtime.SetAutoretry(true)
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer("2021-07-01"), nil, StringPointer("PUT"), StringPointer("AK"), StringPointer(action), query, nil, body, &runtime)
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
				conn, err := client.NewPaiClient()
				if err != nil {
					return WrapError(err)
				}
				request = make(map[string]interface{})
				query = make(map[string]*string)
				body = make(map[string]interface{})
				request["ServiceName"] = d.Id()
				query["ClusterId"] = StringPointer(client.RegionId)
				body = request
				runtime := util.RuntimeOptions{}
				runtime.SetAutoretry(true)
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer("2021-07-01"), nil, StringPointer("PUT"), StringPointer("AK"), StringPointer(action), query, nil, body, &runtime)
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
				stateConf := BuildStateConf([]string{}, []string{"Stopped"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, paiServiceV2.PaiServiceStateRefreshFunc(d.Id(), "Status", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}

			}
		}
	}

	ServiceName := d.Id()
	ClusterId := client.RegionId
	action := fmt.Sprintf("/api/v2/services/%s/%s", ClusterId, ServiceName)
	conn, err := client.NewPaiClient()
	if err != nil {
		return WrapError(err)
	}
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
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer("2021-07-01"), nil, StringPointer("PUT"), StringPointer("AK"), StringPointer(action), query, nil, body, &runtime)
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

	if !d.IsNewResource() && d.HasChange("labels") {
		oldEntry, newEntry := d.GetChange("labels")
		removed := oldEntry.(map[string]interface{})
		added := newEntry.(map[string]interface{})

		if len(removed) > 0 {
			ServiceName := d.Id()
			ClusterId := client.RegionId
			action := fmt.Sprintf("/api/v2/services/%s/%s/label", ClusterId, ServiceName)
			conn, err := client.NewPaiClient()
			if err != nil {
				return WrapError(err)
			}
			request = make(map[string]interface{})
			query = make(map[string]*string)
			body = make(map[string]interface{})
			request["ServiceName"] = d.Id()
			query["ClusterId"] = StringPointer(client.RegionId)

			removedLabels := make([]interface{}, 0)
			for key, value := range removed {
				old, ok := added[key]
				if !ok || old != value {
					removedLabels = append(removedLabels, key)
				}
			}
			query["Keys"] = StringPointer(convertListToJsonString(removedLabels))

			body = request
			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = conn.DoRequest(StringPointer("2021-07-01"), nil, StringPointer("DELETE"), StringPointer("AK"), StringPointer(action), query, nil, body, &runtime)
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

		if len(added) > 0 {
			ServiceName := d.Id()
			ClusterId := client.RegionId
			action := fmt.Sprintf("/api/v2/services/%s/%s/label", ClusterId, ServiceName)
			conn, err := client.NewPaiClient()
			if err != nil {
				return WrapError(err)
			}
			request = make(map[string]interface{})
			query = make(map[string]*string)
			body = make(map[string]interface{})
			request["ServiceName"] = d.Id()
			query["ClusterId"] = StringPointer(client.RegionId)
			request["Labels"] = added

			body = request
			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = conn.DoRequest(StringPointer("2021-07-01"), nil, StringPointer("PUT"), StringPointer("AK"), StringPointer(action), query, nil, body, &runtime)
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
	conn, err := client.NewPaiClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["ServiceName"] = d.Id()
	query["ClusterId"] = StringPointer(client.RegionId)

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer("2021-07-01"), nil, StringPointer("DELETE"), StringPointer("AK"), StringPointer(action), query, nil, nil, &runtime)

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
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, paiServiceV2.PaiServiceStateRefreshFunc(d.Id(), "", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
