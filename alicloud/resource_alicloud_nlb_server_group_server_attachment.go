// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tidwall/sjson"
)

func resourceAliCloudNlbServerGroupServerAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudNlbServerGroupServerAttachmentCreate,
		Read:   resourceAliCloudNlbServerGroupServerAttachmentRead,
		Update: resourceAliCloudNlbServerGroupServerAttachmentUpdate,
		Delete: resourceAliCloudNlbServerGroupServerAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"server_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"server_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"server_ip": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"server_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"weight": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func resourceAliCloudNlbServerGroupServerAttachmentCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "AddServersToServerGroup"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	conn, err := client.NewNlbClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["ServerGroupId"] = d.Get("server_group_id")
	request, _ = sjson.Set(request, "Servers[0].ServerId", d.Get("server_id"))
	request, _ = sjson.Set(request, "Servers[0].ServerType", d.Get("server_type"))
	request, _ = sjson.Set(request, "Servers[0].Port", d.Get("port"))
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-04-30"), StringPointer("AK"), query, request, &runtime)
		request["ClientToken"] = buildClientToken(action)

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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_nlb_server_group_server_attachment", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v:%v:%v", response["ServerGroupId"], request["Servers[*].ServerId"], request["Servers[*].ServerType"], request["Servers[*].Port"]))

	nlbServiceV2 := NlbServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, nlbServiceV2.NlbServerGroupServerAttachmentStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudNlbServerGroupServerAttachmentRead(d, meta)
}

func resourceAliCloudNlbServerGroupServerAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	nlbServiceV2 := NlbServiceV2{client}

	objectRaw, err := nlbServiceV2.DescribeNlbServerGroupServerAttachment(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_nlb_server_group_server_attachment DescribeNlbServerGroupServerAttachment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("description", objectRaw["Description"])
	d.Set("server_ip", objectRaw["ServerIp"])
	d.Set("status", objectRaw["Status"])
	d.Set("weight", objectRaw["Weight"])
	d.Set("port", objectRaw["Port"])
	d.Set("server_group_id", objectRaw["ServerGroupId"])
	d.Set("server_id", objectRaw["ServerId"])
	d.Set("server_type", objectRaw["ServerType"])

	return nil
}

func resourceAliCloudNlbServerGroupServerAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	parts := strings.Split(d.Id(), ":")
	action := "UpdateServerGroupServersAttribute"
	conn, err := client.NewNlbClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ServerGroupId"] = parts[0]
	request, _ = sjson.Set(request, "Servers[0].ServerId", parts[1])
	request, _ = sjson.Set(request, "Servers[0].ServerType", parts[2])
	request, _ = sjson.Set(request, "Servers[0].Port", parts[3])
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-04-30"), StringPointer("AK"), query, request, &runtime)
			request["ClientToken"] = buildClientToken(action)

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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		nlbServiceV2 := NlbServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, nlbServiceV2.NlbServerGroupServerAttachmentStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return resourceAliCloudNlbServerGroupServerAttachmentRead(d, meta)
}

func resourceAliCloudNlbServerGroupServerAttachmentDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "RemoveServersFromServerGroup"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	conn, err := client.NewNlbClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["ServerGroupId"] = parts[0]
	request, _ = sjson.Set(request, "Servers[0].ServerId", parts[1])
	request, _ = sjson.Set(request, "Servers[0].ServerType", parts[2])
	request, _ = sjson.Set(request, "Servers[0].Port", parts[3])
	request["RegionId"] = client.RegionId

	request["ClientToken"] = buildClientToken(action)

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-04-30"), StringPointer("AK"), query, request, &runtime)
		request["ClientToken"] = buildClientToken(action)

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
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	nlbServiceV2 := NlbServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{" Available"}, d.Timeout(schema.TimeoutDelete), 5*time.Second, nlbServiceV2.NlbServerGroupServerAttachmentStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
