package dmsenterprise

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tferrs "github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	tfhelper "github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	tflogs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlicloudDmsEnterpriseProxyAccess() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudDmsEnterpriseProxyAccessCreate,
		Read:   resourceAlicloudDmsEnterpriseProxyAccessRead,
		Delete: resourceAlicloudDmsEnterpriseProxyAccessDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"access_id": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"access_secret": {
				Computed:  true,
				Sensitive: true,
				Type:      schema.TypeString,
			},
			"create_time": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"indep_account": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"indep_password": {
				Optional:  true,
				ForceNew:  true,
				Sensitive: true,
				Type:      schema.TypeString,
			},
			"instance_id": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"origin_info": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"proxy_access_id": {
				Optional: true,
				Computed: true,
				Type:     schema.TypeString,
			},
			"proxy_id": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"user_id": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"user_name": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"user_uid": {
				Computed: true,
				Type:     schema.TypeString,
			},
		},
	}
}

func resourceAlicloudDmsEnterpriseProxyAccessCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := map[string]interface{}{
		"RegionId": client.RegionId,
	}
	var err error

	if v, ok := d.GetOk("indep_account"); ok {
		request["IndepAccount"] = v
	}
	if v, ok := d.GetOk("indep_password"); ok {
		request["IndepPassword"] = v
	}
	if v, ok := d.GetOk("proxy_id"); ok {
		request["ProxyId"] = v
	}
	if v, ok := d.GetOk("user_id"); ok {
		request["UserId"] = v
	}

	var response map[string]interface{}
	action := "CreateProxyAccess"
	wait := tfhelper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *retry.RetryError {
		resp, err := client.RpcPost("dms-enterprise", "2018-11-01", action, nil, request, false)
		if err != nil {
			if tfhelper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		response = resp
		tflogs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		return tferrs.WrapErrorf(err, tferrs.DefaultErrorMsg, "alicloud_dms_enterprise_proxy_access", action, tferrs.AlibabaCloudSdkGoERROR)
	}

	if v, err := jsonpath.Get("$.ProxyAccessId", response); err != nil || v == nil {
		return tferrs.WrapErrorf(err, tferrs.IdMsg, "alicloud_dms_enterprise_proxy_access")
	} else {
		d.SetId(fmt.Sprint(v))
	}

	return resourceAlicloudDmsEnterpriseProxyAccessRead(d, meta)
}

func resourceAlicloudDmsEnterpriseProxyAccessRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dmsEnterpriseService := DmsEnterpriseService{client}

	object, err := dmsEnterpriseService.DescribeDmsEnterpriseProxyAccess(d.Id())
	if err != nil {
		if tferrs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_dms_enterprise_proxy_access dmsEnterpriseService.DescribeDmsEnterpriseProxyAccess Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return tferrs.WrapError(err)
	}
	d.Set("access_id", object["AccessId"])
	d.Set("create_time", object["GmtCreate"])
	d.Set("indep_account", object["IndepAccount"])
	d.Set("instance_id", object["InstanceId"])
	d.Set("origin_info", object["OriginInfo"])
	d.Set("proxy_id", object["ProxyId"])
	d.Set("user_id", object["UserId"])
	d.Set("user_name", object["UserName"])
	d.Set("user_uid", object["UserUid"])
	d.Set("proxy_access_id", object["ProxyAccessId"])

	object, err = dmsEnterpriseService.InspectProxyAccessSecret(d.Id())
	if err != nil {
		return tferrs.WrapError(err)
	}
	d.Set("access_secret", object["AccessSecret"])
	return nil
}

func resourceAlicloudDmsEnterpriseProxyAccessDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var err error

	request := map[string]interface{}{

		"ProxyAccessId": d.Id(),
		"RegionId":      client.RegionId,
	}

	action := "DeleteProxyAccess"
	wait := tfhelper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *retry.RetryError {
		resp, err := client.RpcPost("dms-enterprise", "2018-11-01", action, nil, request, false)
		if err != nil {
			if tfhelper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		tflogs.AddDebug(action, resp, request)
		return nil
	})
	if err != nil {
		if tferrs.NotFoundError(err) {
			return nil
		}
		return tferrs.WrapErrorf(err, tferrs.DefaultErrorMsg, d.Id(), action, tferrs.AlibabaCloudSdkGoERROR)
	}
	return nil
}
