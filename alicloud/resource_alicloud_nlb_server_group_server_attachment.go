package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudNlbServerGroupServerAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudNlbServerGroupServerAttachmentCreate,
		Read:   resourceAlicloudNlbServerGroupServerAttachmentRead,
		Update: resourceAlicloudNlbServerGroupServerAttachmentUpdate,
		Delete: resourceAlicloudNlbServerGroupServerAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"description": {
				Optional:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringLenBetween(2, 256),
			},
			"port": {
				Required:     true,
				ForceNew:     true,
				Type:         schema.TypeInt,
				ValidateFunc: validation.IntBetween(0, 65535),
			},
			"server_group_id": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"server_id": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"server_ip": {
				Optional: true,
				ForceNew: true,
				Computed: true,
				Type:     schema.TypeString,
			},
			"server_type": {
				Required:     true,
				ForceNew:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"Ecs", "Eni", "Eci", "Ip"}, false),
			},
			"status": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"weight": {
				Optional:     true,
				Computed:     true,
				Type:         schema.TypeInt,
				ValidateFunc: validation.IntBetween(0, 100),
			},
			"zone_id": {
				Computed: true,
				Type:     schema.TypeString,
			},
		},
	}
}

func resourceAlicloudNlbServerGroupServerAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	nlbService := NlbService{client}
	request := make(map[string]interface{})
	conn, err := client.NewNlbClient()
	if err != nil {
		return WrapError(err)
	}

	request["Servers.1.Port"] = d.Get("port")
	request["ServerGroupId"] = d.Get("server_group_id")
	request["Servers.1.ServerId"] = d.Get("server_id")
	request["Servers.1.ServerType"] = d.Get("server_type")
	if v, ok := d.GetOk("server_ip"); ok {
		request["Servers.1.ServerIp"] = v
	}
	if v, ok := d.GetOk("weight"); ok {
		request["Servers.1.Weight"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		request["Servers.1.Description"] = v
	}

	request["ClientToken"] = buildClientToken("AddServersToServerGroup")
	var response map[string]interface{}
	action := "AddServersToServerGroup"
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-04-30"), StringPointer("AK"), nil, request, &runtime)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_nlb_server_group_server_attachment", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["ServerGroupId"], ":", request["Servers.1.ServerId"], ":", request["Servers.1.ServerType"], ":", request["Servers.1.Port"]))

	stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, nlbService.NlbServerGroupServerAttachmentStateRefreshFunc(d, []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAlicloudNlbServerGroupServerAttachmentRead(d, meta)
}

func resourceAlicloudNlbServerGroupServerAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	nlbService := NlbService{client}

	object, err := nlbService.DescribeNlbServerGroupServerAttachment(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_nlb_server_group_server_attachment nlbService.DescribeNlbServerGroupServerAttachment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 4)
	if err != nil {
		return WrapError(err)
	}
	d.Set("server_group_id", parts[0])
	d.Set("server_id", parts[1])
	d.Set("server_type", parts[2])
	d.Set("port", formatInt(parts[3]))
	d.Set("description", object["Description"])
	d.Set("server_ip", object["ServerIp"])
	d.Set("status", object["Status"])
	d.Set("weight", formatInt(object["Weight"]))
	d.Set("zone_id", object["ZoneId"])

	return nil
}

func resourceAlicloudNlbServerGroupServerAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	nlbService := NlbService{client}
	conn, err := client.NewNlbClient()
	if err != nil {
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 4)
	if err != nil {
		return WrapError(err)
	}
	update := false
	request := map[string]interface{}{
		"ServerGroupId":        parts[0],
		"Servers.1.ServerId":   parts[1],
		"Servers.1.ServerType": parts[2],
		"Servers.1.Port":       parts[3],
	}

	if v, ok := d.GetOk("server_ip"); ok {
		request["Servers.1.ServerIp"] = v
	}

	if d.HasChange("description") {
		update = true
		if v, ok := d.GetOk("description"); ok {
			request["Servers.1.Description"] = v
		}
	}
	if d.HasChange("weight") {
		update = true
		if v, ok := d.GetOk("weight"); ok {
			request["Servers.1.Weight"] = v
		}
	}

	if update {
		action := "UpdateServerGroupServersAttribute"
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-04-30"), StringPointer("AK"), nil, request, &runtime)
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
		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, nlbService.NlbServerGroupServerAttachmentStateRefreshFunc(d, []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return resourceAlicloudNlbServerGroupServerAttachmentRead(d, meta)
}

func resourceAlicloudNlbServerGroupServerAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	nlbService := NlbService{client}
	conn, err := client.NewNlbClient()
	if err != nil {
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 4)
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"ServerGroupId":        parts[0],
		"Servers.1.ServerId":   parts[1],
		"Servers.1.ServerType": parts[2],
		"Servers.1.Port":       parts[3],
	}

	if v, ok := d.GetOk("server_ip"); ok {
		request["Servers.1.ServerIp"] = v
	}

	request["ClientToken"] = buildClientToken("RemoveServersFromServerGroup")
	action := "RemoveServersFromServerGroup"
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-04-30"), StringPointer("AK"), nil, request, &runtime)
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
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, nlbService.NlbServerGroupServerAttachmentStateRefreshFunc(d, []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
