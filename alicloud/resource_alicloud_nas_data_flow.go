package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudNasDataFlow() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudNasDataFlowCreate,
		Read:   resourceAlicloudNasDataFlowRead,
		Update: resourceAlicloudNasDataFlowUpdate,
		Delete: resourceAlicloudNasDataFlowDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"data_flow_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.All(validation.StringLenBetween(2, 256), validation.StringDoesNotMatch(regexp.MustCompile(`(^http://.*)|(^https://.*)`), "It cannot begin with \"http://\", \"https://\".")),
			},
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"file_system_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"fset_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"source_security_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"SSL", "NONE"}, false),
			},
			"source_storage": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"source_storage_path": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"file_system_path": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"auto_refresh_policy": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"None", "ImportChanged"}, false),
			},
			"auto_refresh_interval": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"auto_refresh": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"refresh_path": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"Running", "Stopped"}, false),
			},
			"throughput": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntInSlice([]int{600, 1200, 1500}),
			},
		},
	}
}

func resourceAlicloudNasDataFlowCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateDataFlow"
	request := make(map[string]interface{})
	var err error
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	request["FileSystemId"] = d.Get("file_system_id")
	request["FsetId"] = d.Get("fset_id")
	if v, ok := d.GetOk("source_security_type"); ok {
		request["SourceSecurityType"] = v
	}
	request["SourceStorage"] = d.Get("source_storage")
	if v, ok := d.GetOk("source_storage_path"); ok {
		request["SourceStoragePath"] = v
	}
	if v, ok := d.GetOk("file_system_path"); ok {
		request["FileSystemPath"] = v
	}
	if v, ok := d.GetOk("auto_refresh_policy"); ok {
		request["AutoRefreshPolicy"] = v
	}
	if v, ok := d.GetOkExists("auto_refresh_interval"); ok {
		request["AutoRefreshInterval"] = v
	}
	if v, ok := d.GetOk("auto_refresh"); ok {
		autoRefreshsMaps := make([]map[string]interface{}, 0)
		for _, raw := range v.([]interface{}) {
			if raw == nil {
				continue
			}
			m := raw.(map[string]interface{})
			autoRefreshsMaps = append(autoRefreshsMaps, map[string]interface{}{
				"RefreshPath": m["refresh_path"],
			})
		}
		request["AutoRefreshs"] = autoRefreshsMaps
	}
	request["Throughput"] = d.Get("throughput")
	request["ClientToken"] = buildClientToken("CreateDataFlow")
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("NAS", "2017-06-26", action, nil, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_nas_data_flow", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["FileSystemId"], ":", response["DataFlowId"]))
	nasService := NasService{client}
	stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, nasService.NasDataFlowStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudNasDataFlowRead(d, meta)
}
func resourceAlicloudNasDataFlowRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	nasService := NasService{client}
	object, err := nasService.DescribeNasDataFlow(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_nas_data_flow nasService.DescribeNasDataFlow Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("data_flow_id", object["DataFlowId"])
	d.Set("file_system_id", object["FileSystemId"])
	d.Set("description", object["Description"])
	d.Set("fset_id", object["FsetId"])
	d.Set("source_security_type", object["SourceSecurityType"])
	d.Set("source_storage", object["SourceStorage"])
	d.Set("source_storage_path", object["SourceStoragePath"])
	d.Set("file_system_path", object["FileSystemPath"])
	d.Set("auto_refresh_policy", object["AutoRefreshPolicy"])
	d.Set("auto_refresh_interval", formatInt(object["AutoRefreshInterval"]))
	autoRefreshList := make([]map[string]interface{}, 0)
	if ar, ok := object["AutoRefresh"]; ok && ar != nil {
		if arMap, ok := ar.(map[string]interface{}); ok {
			if innerList, ok := arMap["AutoRefresh"].([]interface{}); ok {
				for _, item := range innerList {
					if m, ok := item.(map[string]interface{}); ok {
						autoRefreshList = append(autoRefreshList, map[string]interface{}{
							"refresh_path": m["RefreshPath"],
						})
					}
				}
			}
		}
	}
	if err := d.Set("auto_refresh", autoRefreshList); err != nil {
		return WrapError(err)
	}
	d.Set("status", object["Status"])
	d.Set("throughput", formatInt(object["Throughput"]))
	return nil
}
func resourceAlicloudNasDataFlowUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	nasService := NasService{client}
	var err error
	var response map[string]interface{}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Partial(true)

	update := false
	modifyDataFlowReq := map[string]interface{}{
		"DataFlowId":   parts[1],
		"FileSystemId": parts[0],
	}
	if d.HasChange("description") {
		update = true
		if v, ok := d.GetOk("description"); ok {
			modifyDataFlowReq["Description"] = v
		}
	}
	if d.HasChange("throughput") {
		update = true
		modifyDataFlowReq["Throughput"] = d.Get("throughput")
	}
	if update {
		if v, ok := d.GetOkExists("dry_run"); ok {
			modifyDataFlowReq["DryRun"] = v
		}
		action := "ModifyDataFlow"
		modifyDataFlowReq["ClientToken"] = buildClientToken("ModifyDataFlow")
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("NAS", "2017-06-26", action, nil, modifyDataFlowReq, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, modifyDataFlowReq)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, nasService.NasDataFlowStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("description")
		d.SetPartial("throughput")
	}

	if d.HasChange("status") {
		object, err := nasService.DescribeNasDataFlow(d.Id())
		if err != nil {
			return WrapError(err)
		}
		target := d.Get("status").(string)
		if object["Status"].(string) != target {
			if target == "Running" {
				request := map[string]interface{}{
					"DataFlowId":   parts[1],
					"FileSystemId": parts[0],
				}
				if v, ok := d.GetOkExists("dry_run"); ok {
					request["DryRun"] = v
				}
				action := "StartDataFlow"
				request["ClientToken"] = buildClientToken("StartDataFlow")
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcPost("NAS", "2017-06-26", action, nil, request, true)
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
				stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, nasService.NasDataFlowStateRefreshFunc(d.Id(), []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}
			}
			if target == "Stopped" {
				request := map[string]interface{}{
					"DataFlowId":   parts[1],
					"FileSystemId": parts[0],
				}
				if v, ok := d.GetOkExists("dry_run"); ok {
					request["DryRun"] = v
				}
				action := "StopDataFlow"
				request["ClientToken"] = buildClientToken("StopDataFlow")
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcPost("NAS", "2017-06-26", action, nil, request, true)
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
				stateConf := BuildStateConf([]string{}, []string{"Stopped"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, nasService.NasDataFlowStateRefreshFunc(d.Id(), []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}
			}
			d.SetPartial("status")
		}
	}

	if d.HasChange("auto_refresh") || d.HasChange("auto_refresh_policy") || d.HasChange("auto_refresh_interval") {
		if d.HasChange("auto_refresh") {
			// AutoRefresh list changed: apply the full new auto refresh configuration via ApplyDataFlowAutoRefresh.
			request := map[string]interface{}{
				"DataFlowId":   parts[1],
				"FileSystemId": parts[0],
			}
			if v, ok := d.GetOk("auto_refresh_policy"); ok {
				request["AutoRefreshPolicy"] = v
			}
			if v, ok := d.GetOkExists("auto_refresh_interval"); ok {
				request["AutoRefreshInterval"] = v
			}
			autoRefreshsMaps := make([]map[string]interface{}, 0)
			if v, ok := d.GetOk("auto_refresh"); ok {
				for _, raw := range v.([]interface{}) {
					if raw == nil {
						continue
					}
					m := raw.(map[string]interface{})
					autoRefreshsMaps = append(autoRefreshsMaps, map[string]interface{}{
						"RefreshPath": m["refresh_path"],
					})
				}
			}
			request["AutoRefreshs"] = autoRefreshsMaps
			if v, ok := d.GetOkExists("dry_run"); ok {
				request["DryRun"] = v
			}
			action := "ApplyDataFlowAutoRefresh"
			request["ClientToken"] = buildClientToken("ApplyDataFlowAutoRefresh")
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("NAS", "2017-06-26", action, nil, request, true)
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
			stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, nasService.NasDataFlowStateRefreshFunc(d.Id(), []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
			d.SetPartial("auto_refresh")
			d.SetPartial("auto_refresh_policy")
			d.SetPartial("auto_refresh_interval")
		} else {
			// Only AutoRefreshPolicy or AutoRefreshInterval changed: modify via ModifyDataFlowAutoRefresh.
			request := map[string]interface{}{
				"DataFlowId":   parts[1],
				"FileSystemId": parts[0],
			}
			if v, ok := d.GetOk("auto_refresh_policy"); ok {
				request["AutoRefreshPolicy"] = v
			}
			if v, ok := d.GetOkExists("auto_refresh_interval"); ok {
				request["AutoRefreshInterval"] = v
			}
			if v, ok := d.GetOkExists("dry_run"); ok {
				request["DryRun"] = v
			}
			action := "ModifyDataFlowAutoRefresh"
			request["ClientToken"] = buildClientToken("ModifyDataFlowAutoRefresh")
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("NAS", "2017-06-26", action, nil, request, true)
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
			stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, nasService.NasDataFlowStateRefreshFunc(d.Id(), []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
			d.SetPartial("auto_refresh_policy")
			d.SetPartial("auto_refresh_interval")
		}
	}
	d.Partial(false)
	return resourceAlicloudNasDataFlowRead(d, meta)
}
func resourceAlicloudNasDataFlowDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	action := "DeleteDataFlow"
	var response map[string]interface{}
	request := map[string]interface{}{
		"DataFlowId":   parts[1],
		"FileSystemId": parts[0],
	}

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	request["ClientToken"] = buildClientToken("DeleteDataFlow")
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("NAS", "2017-06-26", action, nil, request, true)
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
	nasService := NasService{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, nasService.NasDataFlowStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
