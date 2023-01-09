package alicloud

import (
	"fmt"
	"github.com/PaesslerAG/jsonpath"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudGaCustomRoutingEndpointGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudGaCustomRoutingEndpointGroupCreate,
		Read:   resourceAlicloudGaCustomRoutingEndpointGroupRead,
		Update: resourceAlicloudGaCustomRoutingEndpointGroupUpdate,
		Delete: resourceAlicloudGaCustomRoutingEndpointGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"accelerator_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"listener_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"endpoint_group_region": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"custom_routing_endpoint_group_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudGaCustomRoutingEndpointGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}
	var response map[string]interface{}
	action := "CreateCustomRoutingEndpointGroups"
	request := make(map[string]interface{})
	conn, err := client.NewGaplusClient()
	if err != nil {
		return WrapError(err)
	}

	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken("CreateCustomRoutingEndpointGroups")
	request["AcceleratorId"] = d.Get("accelerator_id")
	request["ListenerId"] = d.Get("listener_id")

	endpointGroupConfigurationsMaps := make([]map[string]interface{}, 0)
	endpointGroupConfigurationsMap := map[string]interface{}{}
	endpointGroupConfigurationsMap["EndpointGroupRegion"] = d.Get("endpoint_group_region")

	if v, ok := d.GetOk("custom_routing_endpoint_group_name"); ok {
		endpointGroupConfigurationsMap["Name"] = v
	}

	if v, ok := d.GetOk("description"); ok {
		endpointGroupConfigurationsMap["Description"] = v
	}

	endpointGroupConfigurationsMaps = append(endpointGroupConfigurationsMaps, endpointGroupConfigurationsMap)
	request["EndpointGroupConfigurations"] = endpointGroupConfigurationsMaps

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-20"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"NotActive.Listener", "StateError.Accelerator"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ga_custom_routing_endpoint_group", action, AlibabaCloudSdkGoERROR)
	}

	if resp, err := jsonpath.Get("$.EndpointGroupIds", response); err != nil || resp == nil {
		return WrapErrorf(err, IdMsg, "alicloud_ga_custom_routing_endpoint_group")
	} else {
		endpointGroupId := resp.([]interface{})[0]
		d.SetId(fmt.Sprint(endpointGroupId))
	}

	stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, gaService.GaCustomRoutingEndpointGroupStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudGaCustomRoutingEndpointGroupRead(d, meta)
}

func resourceAlicloudGaCustomRoutingEndpointGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}

	object, err := gaService.DescribeGaCustomRoutingEndpointGroup(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("accelerator_id", object["AcceleratorId"])
	d.Set("listener_id", object["ListenerId"])
	d.Set("endpoint_group_region", object["EndpointGroupRegion"])
	d.Set("custom_routing_endpoint_group_name", object["Name"])
	d.Set("description", object["Description"])
	d.Set("status", object["State"])

	return nil
}

func resourceAlicloudGaCustomRoutingEndpointGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}
	var response map[string]interface{}
	update := false

	request := map[string]interface{}{
		"RegionId":        client.RegionId,
		"ClientToken":     buildClientToken("UpdateCustomRoutingEndpointGroupAttribute"),
		"EndpointGroupId": d.Id(),
	}

	if d.HasChange("custom_routing_endpoint_group_name") {
		update = true
	}
	if v, ok := d.GetOk("custom_routing_endpoint_group_name"); ok {
		request["Name"] = v
	}

	if d.HasChange("description") {
		update = true
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	if update {
		action := "UpdateCustomRoutingEndpointGroupAttribute"
		conn, err := client.NewGaplusClient()
		if err != nil {
			return WrapError(err)
		}

		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-20"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if IsExpectedErrors(err, []string{"NotActive.Listener", "StateError.Accelerator"}) || NeedRetry(err) {
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

		stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, gaService.GaCustomRoutingEndpointGroupStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return resourceAlicloudGaCustomRoutingEndpointGroupRead(d, meta)
}

func resourceAlicloudGaCustomRoutingEndpointGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}
	action := "DeleteCustomRoutingEndpointGroups"
	var response map[string]interface{}

	conn, err := client.NewGaplusClient()
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"RegionId":         client.RegionId,
		"ClientToken":      buildClientToken("DeleteCustomRoutingEndpointGroups"),
		"EndpointGroupIds": []string{d.Id()},
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-20"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"NotActive.Listener", "StateError.Accelerator"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"NotExist.EndPointGroup"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, gaService.GaCustomRoutingEndpointGroupStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
