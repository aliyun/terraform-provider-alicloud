package dmsenterprise

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfconvert "github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	tferrs "github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	tfhelper "github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	tflogs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlicloudDmsEnterpriseProxy() *schema.Resource {
	return &schema.Resource{
		Create: ResourceAlicloudDmsEnterpriseProxyCreate,
		Read:   ResourceAlicloudDmsEnterpriseProxyRead,
		Update: resourceAlicloudDmsEnterpriseProxyUpdate,
		Delete: ResourceAlicloudDmsEnterpriseProxyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"tid": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"username": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
		},
	}
}

func ResourceAlicloudDmsEnterpriseProxyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateProxy"
	request := make(map[string]interface{})
	var err error
	request["InstanceId"] = d.Get("instance_id")
	request["Password"] = d.Get("password")
	if v, ok := d.GetOk("tid"); ok {
		request["Tid"] = v
	}
	request["Username"] = d.Get("username")
	wait := tfhelper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *retry.RetryError {
		response, err = client.RpcPost("dms-enterprise", "2018-11-01", action, nil, request, false)
		if err != nil {
			if tfhelper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	tflogs.AddDebug(action, response, request)
	if err != nil {
		return tferrs.WrapErrorf(err, tferrs.DefaultErrorMsg, "alicloud_dms_enterprise_proxy", action, tferrs.AlibabaCloudSdkGoERROR)
	}

	proxyId, err := tfconvert.FormatInt(response["ProxyId"])
	if err != nil {
		return err
	}
	d.SetId(fmt.Sprint(proxyId))

	return ResourceAlicloudDmsEnterpriseProxyRead(d, meta)
}
func ResourceAlicloudDmsEnterpriseProxyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dmsEnterpriseService := DmsEnterpriseService{client}
	object, err := dmsEnterpriseService.DescribeDmsEnterpriseProxy(d.Id())
	if err != nil {
		if tferrs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_dms_enterprise_proxy dmsEnterpriseService.DescribeDmsEnterpriseProxy Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return tferrs.WrapError(err)
	}
	instanceId, err := tfconvert.FormatInt(object["InstanceId"])
	if err != nil {
		return err
	}
	d.Set("instance_id", fmt.Sprint(instanceId))
	return nil
}
func resourceAlicloudDmsEnterpriseProxyUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Println(fmt.Sprintf("[WARNING] The resouce has not update operation."))
	return ResourceAlicloudDmsEnterpriseProxyRead(d, meta)
}
func ResourceAlicloudDmsEnterpriseProxyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteProxy"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"ProxyId": d.Id(),
	}

	if v, ok := d.GetOk("tid"); ok {
		request["Tid"] = v
	}
	wait := tfhelper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *retry.RetryError {
		response, err = client.RpcPost("dms-enterprise", "2018-11-01", action, nil, request, false)
		if err != nil {
			if tfhelper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	tflogs.AddDebug(action, response, request)
	if err != nil {
		return tferrs.WrapErrorf(err, tferrs.DefaultErrorMsg, d.Id(), action, tferrs.AlibabaCloudSdkGoERROR)
	}
	return nil
}
