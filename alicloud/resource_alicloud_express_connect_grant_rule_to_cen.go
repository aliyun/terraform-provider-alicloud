package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudExpressConnectGrantRuleToCen() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudExpressConnectGrantRuleToCenCreate,
		Read:   resourceAlicloudExpressConnectGrantRuleToCenRead,
		Delete: resourceAlicloudExpressConnectGrantRuleToCenDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cen_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cen_owner_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudExpressConnectGrantRuleToCenCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "GrantInstanceToCen"
	request := make(map[string]interface{})
	var err error

	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken("GrantInstanceToCen")
	request["CenId"] = d.Get("cen_id")
	request["CenOwnerId"] = d.Get("cen_owner_id")
	request["InstanceId"] = d.Get("instance_id")
	request["InstanceType"] = "VBR"

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_express_connect_grant_rule_to_cen", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v:%v", request["CenId"], request["CenOwnerId"], request["InstanceId"]))

	return resourceAlicloudExpressConnectGrantRuleToCenRead(d, meta)
}

func resourceAlicloudExpressConnectGrantRuleToCenRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}

	object, err := vpcService.DescribeExpressConnectGrantRuleToCen(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}

	d.Set("cen_id", object["CenInstanceId"])
	d.Set("cen_owner_id", formatInt(object["CenOwnerId"]))
	d.Set("instance_id", parts[2])

	return nil
}

func resourceAlicloudExpressConnectGrantRuleToCenDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "RevokeInstanceFromCen"
	var response map[string]interface{}

	var err error

	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"RegionId":     client.RegionId,
		"ClientToken":  buildClientToken("RevokeInstanceFromCen"),
		"CenId":        parts[0],
		"CenOwnerId":   parts[1],
		"InstanceId":   parts[2],
		"InstanceType": "VBR",
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
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
