package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudGaServiceExtension() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudGaServiceExtensionCreate,
		Read:   resourceAlicloudGaServiceExtensionRead,
		Update: resourceAlicloudGaServiceExtensionUpdate,
		Delete: resourceAlicloudGaServiceExtensionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"service_extension_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tags": tagsSchema(),
			"components": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service_component_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"priority": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"timeout": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"config": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"fail_policy": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"component_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"resources": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"accelerator_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"resource_type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"resource_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"associate_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudGaServiceExtensionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateServiceExtension"
	request := make(map[string]interface{})
	var err error

	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken("CreateServiceExtension")
	request["Name"] = d.Get("name")

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ga_service_extension", action, AlibabaCloudSdkGoERROR)
	}

	id, err := jsonpath.Get("$.ServiceExtensionId", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, "alicloud_ga_service_extension", "$.ServiceExtensionId", response)
	}
	d.SetId(fmt.Sprint(id))

	gaService := GaService{client}
	stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, gaService.GaServiceExtensionStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudGaServiceExtensionUpdate(d, meta)
}

func resourceAlicloudGaServiceExtensionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}

	object, err := gaService.DescribeGaServiceExtension(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ga_service_extension gaService.DescribeGaServiceExtension Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("service_extension_id", object["ServiceExtensionId"])
	d.Set("name", object["Name"])
	d.Set("description", object["Description"])
	d.Set("type", object["Type"])
	d.Set("state", object["State"])
	d.Set("resource_group_id", object["ResourceGroupId"])
	d.Set("create_time", object["CreateTime"])
	d.Set("update_time", object["UpdateTime"])

	if v, ok := object["Components"].([]interface{}); ok {
		components := make([]map[string]interface{}, 0)
		for _, val := range v {
			item := val.(map[string]interface{})
			temp := map[string]interface{}{
				"service_component_id": item["ServiceComponentId"],
				"priority":             item["Priority"],
				"timeout":              item["Timeout"],
				"config":               item["Config"],
				"fail_policy":          item["FailPolicy"],
				"component_name":       item["ComponentName"],
			}
			components = append(components, temp)
		}
		if err := d.Set("components", components); err != nil {
			return WrapError(err)
		}
	}

	if v, ok := object["Resources"].([]interface{}); ok {
		resources := make([]map[string]interface{}, 0)
		for _, val := range v {
			item := val.(map[string]interface{})
			temp := map[string]interface{}{
				"accelerator_id": item["AcceleratorId"],
				"resource_type":  item["ResourceType"],
				"resource_id":    item["ResourceId"],
				"associate_id":   item["AssociateId"],
			}
			resources = append(resources, temp)
		}
		if err := d.Set("resources", resources); err != nil {
			return WrapError(err)
		}
	}

	listTagResourcesObject, err := gaService.ListTagResources(d.Id(), "serviceextension")
	if err != nil {
		// GA TagResources/ListTagResources has not registered a ResourceType for
		// service extension (Invalid.ResourceType). Tolerate it so Read does not
		// block refresh/destroy and DeleteServiceExtension stays reachable.
		if IsExpectedErrors(err, []string{"Invalid.ResourceType"}) {
			log.Printf("[WARN] Resource alicloud_ga_service_extension tagging is not supported by the backend (Invalid.ResourceType), skipping tags for %s: %s", d.Id(), err)
		} else {
			return WrapError(err)
		}
	} else {
		d.Set("tags", tagsToMap(listTagResourcesObject))
	}

	return nil
}

func resourceAlicloudGaServiceExtensionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}
	var response map[string]interface{}
	d.Partial(true)

	if d.HasChange("tags") {
		if err := gaService.SetResourceTags(d, "serviceextension"); err != nil {
			// GA tagging backend has no valid ResourceType for service extension
			// (Invalid.ResourceType). Best-effort: do not fatal-block updates.
			if IsExpectedErrors(err, []string{"Invalid.ResourceType"}) {
				log.Printf("[WARN] Resource alicloud_ga_service_extension tagging is not supported by the backend (Invalid.ResourceType), skipping tags update for %s: %s", d.Id(), err)
			} else {
				return WrapError(err)
			}
		}
		d.SetPartial("tags")
	}

	// UpdateServiceExtension for name/description
	update := false
	request := map[string]interface{}{
		"RegionId":           client.RegionId,
		"ClientToken":        buildClientToken("UpdateServiceExtension"),
		"ServiceExtensionId": d.Id(),
	}
	if !d.IsNewResource() && d.HasChange("name") {
		update = true
	}
	request["Name"] = d.Get("name")
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if !d.IsNewResource() && d.HasChange("description") {
		update = true
	}
	if update {
		action := "UpdateServiceExtension"
		var err error
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
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
		// UpdateServiceExtension is asynchronous: the resource enters state=updating
		// after a successful call and must converge to active before any subsequent
		// update/delete is issued, otherwise the backend rejects it with
		// UpdateFailed.ServiceExtension. Wait for active here.
		stateConf := BuildStateConf([]string{"updating"}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, gaService.GaServiceExtensionStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("name")
		d.SetPartial("description")
	}

	// Components add/remove
	if !d.IsNewResource() && d.HasChange("components") {
		oraw, nraw := d.GetChange("components")
		remove := oraw.(*schema.Set).Difference(nraw.(*schema.Set)).List()
		create := nraw.(*schema.Set).Difference(oraw.(*schema.Set)).List()

		if len(remove) > 0 {
			removeReq := map[string]interface{}{
				"RegionId":           client.RegionId,
				"ClientToken":        buildClientToken("RemoveServiceComponentsFromServiceExtension"),
				"ServiceExtensionId": d.Id(),
			}
			for k, v := range remove {
				item := v.(map[string]interface{})
				removeReq[fmt.Sprintf("ServiceComponentIds.%d", k+1)] = item["service_component_id"].(string)
			}
			action := "RemoveServiceComponentsFromServiceExtension"
			var err error
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
				response, err = client.RpcPost("Ga", "2019-11-20", action, nil, removeReq, true)
				if err != nil {
					if NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			addDebug(action, response, removeReq)
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
			// RemoveServiceComponentsFromServiceExtension is asynchronous: wait for
			// the resource to converge from updating back to active.
			stateConf := BuildStateConf([]string{"updating"}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, gaService.GaServiceExtensionStateRefreshFunc(d.Id(), []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
		}

		if len(create) > 0 {
			addReq := map[string]interface{}{
				"RegionId":           client.RegionId,
				"ClientToken":        buildClientToken("AddServiceComponentsToServiceExtension"),
				"ServiceExtensionId": d.Id(),
			}
			for k, v := range create {
				item := v.(map[string]interface{})
				addReq[fmt.Sprintf("ServiceComponents.%d.ServiceComponentId", k+1)] = item["service_component_id"].(string)
				addReq[fmt.Sprintf("ServiceComponents.%d.Priority", k+1)] = item["priority"]
				addReq[fmt.Sprintf("ServiceComponents.%d.Timeout", k+1)] = item["timeout"]
				if v, ok := item["config"].(string); ok && v != "" {
					addReq[fmt.Sprintf("ServiceComponents.%d.Config", k+1)] = v
				}
				if v, ok := item["fail_policy"].(string); ok && v != "" {
					addReq[fmt.Sprintf("ServiceComponents.%d.FailPolicy", k+1)] = v
				}
			}
			action := "AddServiceComponentsToServiceExtension"
			var err error
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
				response, err = client.RpcPost("Ga", "2019-11-20", action, nil, addReq, true)
				if err != nil {
					if NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			addDebug(action, response, addReq)
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
			// AddServiceComponentsToServiceExtension is asynchronous: wait for
			// the resource to converge from updating back to active.
			stateConf := BuildStateConf([]string{"updating"}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, gaService.GaServiceExtensionStateRefreshFunc(d.Id(), []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
		}
		d.SetPartial("components")
	}

	// Resources associate/dissociate
	if !d.IsNewResource() && d.HasChange("resources") {
		oraw, nraw := d.GetChange("resources")
		remove := oraw.(*schema.Set).Difference(nraw.(*schema.Set)).List()
		create := nraw.(*schema.Set).Difference(oraw.(*schema.Set)).List()

		if len(remove) > 0 {
			removeReq := map[string]interface{}{
				"RegionId":           client.RegionId,
				"ClientToken":        buildClientToken("DissociateResourcesFromServiceExtension"),
				"ServiceExtensionId": d.Id(),
			}
			for k, v := range remove {
				item := v.(map[string]interface{})
				if associateId, ok := item["associate_id"].(string); ok && associateId != "" {
					removeReq[fmt.Sprintf("AssociateIds.%d", k+1)] = associateId
				}
			}
			action := "DissociateResourcesFromServiceExtension"
			var err error
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
				response, err = client.RpcPost("Ga", "2019-11-20", action, nil, removeReq, true)
				if err != nil {
					if NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			addDebug(action, response, removeReq)
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
			// DissociateResourcesFromServiceExtension is asynchronous: wait for
			// the resource to converge from updating back to active.
			stateConf := BuildStateConf([]string{"updating"}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, gaService.GaServiceExtensionStateRefreshFunc(d.Id(), []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
		}

		if len(create) > 0 {
			addReq := map[string]interface{}{
				"RegionId":           client.RegionId,
				"ClientToken":        buildClientToken("AssociateResourcesWithServiceExtension"),
				"ServiceExtensionId": d.Id(),
			}
			for k, v := range create {
				item := v.(map[string]interface{})
				addReq[fmt.Sprintf("Resources.%d.AcceleratorId", k+1)] = item["accelerator_id"].(string)
				addReq[fmt.Sprintf("Resources.%d.ResourceType", k+1)] = item["resource_type"].(string)
				addReq[fmt.Sprintf("Resources.%d.ResourceId", k+1)] = item["resource_id"].(string)
			}
			action := "AssociateResourcesWithServiceExtension"
			var err error
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
				response, err = client.RpcPost("Ga", "2019-11-20", action, nil, addReq, true)
				if err != nil {
					if NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			addDebug(action, response, addReq)
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
			// AssociateResourcesWithServiceExtension is asynchronous: wait for
			// the resource to converge from updating back to active.
			stateConf := BuildStateConf([]string{"updating"}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, gaService.GaServiceExtensionStateRefreshFunc(d.Id(), []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
		}
		d.SetPartial("resources")
	}

	// resource_group_id change
	update = false
	changeResourceGroupReq := map[string]interface{}{
		"RegionId":     client.RegionId,
		"ClientToken":  buildClientToken("ChangeResourceGroup"),
		"ResourceId":   d.Id(),
		"ResourceType": "serviceextension",
	}
	if !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		changeResourceGroupReq["NewResourceGroupId"] = v
	}
	if update {
		action := "ChangeResourceGroup"
		var err error
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("Ga", "2019-11-20", action, nil, changeResourceGroupReq, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, changeResourceGroupReq)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("resource_group_id")
	}

	d.Partial(false)
	return resourceAlicloudGaServiceExtensionRead(d, meta)
}

func resourceAlicloudGaServiceExtensionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}
	action := "DeleteServiceExtension"
	var response map[string]interface{}
	var err error

	request := map[string]interface{}{
		"RegionId":           client.RegionId,
		"ClientToken":        buildClientToken("DeleteServiceExtension"),
		"ServiceExtensionId": d.Id(),
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
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
		if NotFoundError(err) || IsExpectedErrors(err, []string{"NotExist.ServiceExtension"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, gaService.GaServiceExtensionStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
