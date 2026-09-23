package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudDmsEnterpriseProxy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudDmsEnterpriseProxyCreate,
		Read:   resourceAlicloudDmsEnterpriseProxyRead,
		Update: resourceAlicloudDmsEnterpriseProxyUpdate,
		Delete: resourceAlicloudDmsEnterpriseProxyDelete,
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

func resourceAlicloudDmsEnterpriseProxyCreate(d *schema.ResourceData, meta interface{}) error {
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
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("dms-enterprise", "2018-11-01", action, nil, request, false)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_dms_enterprise_proxy", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(formatInt(response["ProxyId"])))

	return resourceAlicloudDmsEnterpriseProxyRead(d, meta)
}
func resourceAlicloudDmsEnterpriseProxyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dmsEnterpriseService := DmsEnterpriseService{client}
	object, err := dmsEnterpriseService.DescribeDmsEnterpriseProxy(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_dms_enterprise_proxy dmsEnterpriseService.DescribeDmsEnterpriseProxy Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("instance_id", fmt.Sprint(formatInt(object["InstanceId"])))
	return nil
}
func resourceAlicloudDmsEnterpriseProxyUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Println(fmt.Sprintf("[WARNING] The resouce has not update operation."))
	return resourceAlicloudDmsEnterpriseProxyRead(d, meta)
}
func resourceAlicloudDmsEnterpriseProxyDelete(d *schema.ResourceData, meta interface{}) error {
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
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.RpcPost("dms-enterprise", "2018-11-01", action, nil, request, false)
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
