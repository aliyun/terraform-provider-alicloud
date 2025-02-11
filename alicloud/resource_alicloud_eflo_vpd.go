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

func resourceAlicloudEfloVpd() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudEfloVpdCreate,
		Read:   resourceAlicloudEfloVpdRead,
		Update: resourceAlicloudEfloVpdUpdate,
		Delete: resourceAlicloudEfloVpdDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cidr": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"create_time": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"gmt_modified": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"resource_group_id": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"status": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"vpd_name": {
				Required: true,
				Type:     schema.TypeString,
			},
		},
	}
}

func resourceAlicloudEfloVpdCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	efloService := EfloService{client}
	request := make(map[string]interface{})
	var err error

	request["ClientToken"] = buildClientToken("CreateVpd")
	request["VpdName"] = d.Get("vpd_name")
	request["Cidr"] = d.Get("cidr")
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	var response map[string]interface{}
	action := "CreateVpd"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		resp, err := client.RpcPost("eflo", "2022-05-30", action, nil, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_eflo_vpd", action, AlibabaCloudSdkGoERROR)
	}

	if v, err := jsonpath.Get("$.Content.VpdId", response); err != nil || v == nil {
		return WrapErrorf(err, IdMsg, "alicloud_eflo_vpd")
	} else {
		d.SetId(fmt.Sprint(v))
	}
	stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, efloService.EfloVpdStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAlicloudEfloVpdRead(d, meta)
}

func resourceAlicloudEfloVpdRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	efloService := EfloService{client}

	object, err := efloService.DescribeEfloVpd(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_eflo_vpd efloService.DescribeEfloVpd Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("cidr", object["Cidr"])
	d.Set("create_time", object["CreateTime"])
	d.Set("gmt_modified", object["GmtModified"])
	d.Set("resource_group_id", object["ResourceGroupId"])
	d.Set("status", object["Status"])
	d.Set("vpd_name", object["VpdName"])

	return nil
}

func resourceAlicloudEfloVpdUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var err error
	update := false
	request := map[string]interface{}{
		"VpdId":    d.Id(),
		"RegionId": client.RegionId,
	}

	if d.HasChange("vpd_name") {
		update = true
		if v, ok := d.GetOk("vpd_name"); ok {
			request["VpdName"] = v
		}
	}

	if update {
		action := "UpdateVpd"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			resp, err := client.RpcPost("eflo", "2022-05-30", action, nil, request, true)
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAlicloudEfloVpdRead(d, meta)
}

func resourceAlicloudEfloVpdDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	efloService := EfloService{client}
	var response map[string]interface{}
	var err error

	request := map[string]interface{}{

		"VpdId": d.Id(),
	}

	action := "DeleteVpd"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.RpcPost("eflo", "2022-05-30", action, nil, request, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		if NotFoundError(err) || IsExpectedErrors(err, []string{"1003"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, efloService.EfloVpdStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
