package alicloud

import (
	"fmt"
	"regexp"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudCenTransitRouterCidr() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCenTransitRouterCidrCreate,
		Read:   resourceAlicloudCenTransitRouterCidrRead,
		Update: resourceAlicloudCenTransitRouterCidrUpdate,
		Delete: resourceAlicloudCenTransitRouterCidrDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"transit_router_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cidr": {
				Type:     schema.TypeString,
				Required: true,
			},
			"transit_router_cidr_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[a-zA-Z][-_a-zA-Z0-9]{1,127}$`), "The name can be up to 128 characters in length and can contain digits, letters, hyphens (-), and underscores (_). It must start with a digit or letter."),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.All(validation.StringLenBetween(2, 256), validation.StringDoesNotMatch(regexp.MustCompile(`(^http://.*)|(^https://.*)`), "It must be `2` to `256` characters in length and cannot start with `https://` or `https://`.")),
			},
			"publish_cidr_route": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"transit_router_cidr_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudCenTransitRouterCidrCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateTransitRouterCidr"
	request := make(map[string]interface{})
	var err error

	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken("CreateTransitRouterCidr")
	request["TransitRouterId"] = d.Get("transit_router_id")
	request["Cidr"] = d.Get("cidr")

	if v, ok := d.GetOk("transit_router_cidr_name"); ok {
		request["Name"] = v
	}

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	if v, ok := d.GetOkExists("publish_cidr_route"); ok {
		request["PublishCidrRoute"] = v
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("Cbn", "2017-09-12", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, []string{"Operation.Blocking", "IncorrectStatus.Status"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cen_transit_router_cidr", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["TransitRouterId"], response["TransitRouterCidrId"]))

	return resourceAlicloudCenTransitRouterCidrRead(d, meta)
}

func resourceAlicloudCenTransitRouterCidrRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := CbnService{client}
	object, err := cbnService.DescribeCenTransitRouterCidr(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("transit_router_id", object["TransitRouterId"])
	d.Set("transit_router_cidr_id", object["TransitRouterCidrId"])
	d.Set("cidr", object["Cidr"])
	d.Set("transit_router_cidr_name", object["Name"])
	d.Set("description", object["Description"])
	d.Set("publish_cidr_route", object["PublishCidrRoute"])

	return nil
}

func resourceAlicloudCenTransitRouterCidrUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	update := false

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"RegionId":            client.RegionId,
		"ClientToken":         buildClientToken("ModifyTransitRouterCidr"),
		"TransitRouterId":     parts[0],
		"TransitRouterCidrId": parts[1],
	}

	if d.HasChange("cidr") {
		update = true
	}
	request["Cidr"] = d.Get("cidr")

	if d.HasChange("transit_router_cidr_name") {
		update = true
	}
	if v, ok := d.GetOk("transit_router_cidr_name"); ok {
		request["Name"] = v
	}

	if d.HasChange("description") {
		update = true
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	if d.HasChange("publish_cidr_route") {
		update = true
	}
	if v, ok := d.GetOkExists("publish_cidr_route"); ok {
		request["PublishCidrRoute"] = v
	}

	if update {
		action := "ModifyTransitRouterCidr"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("Cbn", "2017-09-12", action, nil, request, false)
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

	return resourceAlicloudCenTransitRouterCidrRead(d, meta)
}

func resourceAlicloudCenTransitRouterCidrDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteTransitRouterCidr"
	var response map[string]interface{}

	var err error

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"RegionId":            client.RegionId,
		"ClientToken":         buildClientToken("DeleteTransitRouterCidr"),
		"TransitRouterId":     parts[0],
		"TransitRouterCidrId": parts[1],
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.RpcPost("Cbn", "2017-09-12", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, []string{"Operation.Blocking", "IncorrectStatus.Status"}) {
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
