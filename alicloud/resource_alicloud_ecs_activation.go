package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudEcsActivation() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudEcsActivationCreate,
		Read:   resourceAlicloudEcsActivationRead,
		Delete: resourceAlicloudEcsActivationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.All(validation.StringLenBetween(1, 100), validation.StringDoesNotMatch(regexp.MustCompile(`(^http://.*)|(^https://.*)`), "It must be `1` to `100` characters in length and cannot start with `https://` or `https://`.")),
			},
			"instance_count": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 1000),
			},
			"instance_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.All(validation.StringMatch(regexp.MustCompile(`^[a-zA-Z][A-Za-z0-9.:_-]{1,49}$`), "The name must be `1` to `50` characters in length and can contain digits, periods (.), colons (:), underscores (_), and hyphens (-)."), validation.StringDoesNotMatch(regexp.MustCompile(`(^http://.*)|(^https://.*)`), "It must cannot start with `https://`, `https://`")),
			},
			"ip_address_range": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"time_to_live_in_hours": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 24),
			},
		},
	}
}

func resourceAlicloudEcsActivationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateActivation"
	request := make(map[string]interface{})
	conn, err := client.NewEcsClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("instance_count"); ok {
		request["InstanceCount"] = v
	}
	if v, ok := d.GetOk("instance_name"); ok {
		request["InstanceName"] = v
	}
	if v, ok := d.GetOk("ip_address_range"); ok {
		request["IpAddressRange"] = v
	}
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("time_to_live_in_hours"); ok {
		request["TimeToLiveInHours"] = v
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ecs_activation", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["ActivationId"]))

	return resourceAlicloudEcsActivationRead(d, meta)
}
func resourceAlicloudEcsActivationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	object, err := ecsService.DescribeEcsActivation(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ecs_activation ecsService.DescribeEcsActivation Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("description", object["Description"])
	if v, ok := object["InstanceCount"]; ok && fmt.Sprint(v) != "0" {
		d.Set("instance_count", formatInt(v))
	}
	d.Set("instance_name", object["InstanceName"])
	d.Set("ip_address_range", object["IpAddressRange"])
	d.Set("time_to_live_in_hours", formatInt(object["TimeToLiveInHours"]))

	return nil
}
func resourceAlicloudEcsActivationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteActivation"
	var response map[string]interface{}
	conn, err := client.NewEcsClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"ActivationId": d.Id(),
	}

	request["RegionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	return nil
}
