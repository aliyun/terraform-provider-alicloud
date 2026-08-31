package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAliCloudGaBasicEndpoint() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudGaBasicEndpointCreate,
		Read:   resourceAliCloudGaBasicEndpointRead,
		Update: resourceAliCloudGaBasicEndpointUpdate,
		Delete: resourceAliCloudGaBasicEndpointDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Minute),
			Update: schema.DefaultTimeout(3 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"accelerator_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"endpoint_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"endpoint_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"ENI", "SLB", "ECS", "NLB"}, false),
			},
			"endpoint_address": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"endpoint_sub_address_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"primary", "secondary"}, false),
			},
			"endpoint_sub_address": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"endpoint_zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"basic_endpoint_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"endpoint_id": {
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

func resourceAliCloudGaBasicEndpointCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}
	var response map[string]interface{}
	action := "CreateBasicEndpoint"
	request := make(map[string]interface{})
	var err error

	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken("CreateBasicEndpoint")
	request["AcceleratorId"] = d.Get("accelerator_id")
	request["EndpointGroupId"] = d.Get("endpoint_group_id")
	request["EndpointType"] = d.Get("endpoint_type")
	request["EndpointAddress"] = d.Get("endpoint_address")

	if v, ok := d.GetOk("endpoint_sub_address_type"); ok {
		request["EndpointSubAddressType"] = v
	}

	if v, ok := d.GetOk("endpoint_sub_address"); ok {
		request["EndpointSubAddress"] = v
	}

	if v, ok := d.GetOk("endpoint_zone_id"); ok {
		request["EndpointZoneId"] = v
	}

	if v, ok := d.GetOk("basic_endpoint_name"); ok {
		request["Name"] = v
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"NotExist.Accelerator", "NotExist.EndPointGroup", "NotExist.Nlb", "NotFound.ECSInstance", "StateError.Accelerator"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ga_basic_endpoint", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["EndpointGroupId"], response["EndpointId"]))

	stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, gaService.GaBasicEndpointStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudGaBasicEndpointRead(d, meta)
}

func resourceAliCloudGaBasicEndpointRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}

	object, err := gaService.DescribeGaBasicEndpoint(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("accelerator_id", object["AcceleratorId"])
	d.Set("endpoint_group_id", object["EndpointGroupId"])
	d.Set("endpoint_type", object["EndpointType"])
	d.Set("endpoint_address", object["EndpointAddress"])
	d.Set("endpoint_sub_address_type", object["EndpointSubAddressType"])
	d.Set("endpoint_sub_address", object["EndpointSubAddress"])
	d.Set("endpoint_zone_id", object["EndpointZoneId"])
	d.Set("basic_endpoint_name", object["Name"])
	d.Set("endpoint_id", object["EndPointId"])
	d.Set("status", object["State"])

	return nil
}

func resourceAliCloudGaBasicEndpointUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	update := false

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"RegionId":        client.RegionId,
		"ClientToken":     buildClientToken("UpdateBasicEndpoint"),
		"EndpointGroupId": parts[0],
		"EndpointId":      parts[1],
	}

	if d.HasChange("basic_endpoint_name") {
		update = true
	}
	if v, ok := d.GetOk("basic_endpoint_name"); ok {
		request["Name"] = v
	}

	if update {
		action := "UpdateBasicEndpoint"
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
	}

	return resourceAliCloudGaBasicEndpointRead(d, meta)
}

func resourceAliCloudGaBasicEndpointDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}
	action := "DeleteBasicEndpoint"
	var response map[string]interface{}

	var err error

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"RegionId":        client.RegionId,
		"ClientToken":     buildClientToken("DeleteBasicEndpoint"),
		"EndpointGroupId": parts[0],
		"EndpointId":      parts[1],
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"StateError.Accelerator", "StateError.EndPointGroup"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"NotExist.EndPoints"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, gaService.GaBasicEndpointStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
