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
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		resp, err := client.RpcPost("dms-enterprise", "2018-11-01", action, nil, request, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		response = resp
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_dms_enterprise_proxy_access", action, AlibabaCloudSdkGoERROR)
	}

	if v, err := jsonpath.Get("$.ProxyAccessId", response); err != nil || v == nil {
		return WrapErrorf(err, IdMsg, "alicloud_dms_enterprise_proxy_access")
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
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_dms_enterprise_proxy_access dmsEnterpriseService.DescribeDmsEnterpriseProxyAccess Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
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
		return WrapError(err)
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
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		resp, err := client.RpcPost("dms-enterprise", "2018-11-01", action, nil, request, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, resp, request)
		return nil
	})
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
