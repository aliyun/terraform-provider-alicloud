package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudSlbServerGroupServerAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudSlbServerGroupServerAttachmentCreate,
		Read:   resourceAlicloudSlbServerGroupServerAttachmentRead,
		Delete: resourceAlicloudSlbServerGroupServerAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
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
			"port": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(1, 65535),
			},
			"weight": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(0, 100),
			},
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"eni", "ecs"}, false),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudSlbServerGroupServerAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "AddVServerGroupBackendServers"
	request := make(map[string]interface{})
	conn, err := client.NewSlbClient()
	if err != nil {
		return WrapError(err)
	}
	request = map[string]interface{}{
		"RegionId": client.RegionId,
	}
	request["VServerGroupId"] = d.Get("server_group_id")

	serverMaps := make([]map[string]interface{}, 1)
	server := make(map[string]interface{}, 2)
	server["ServerId"] = d.Get("server_id")
	server["Port"] = d.Get("port")

	if v, ok := d.GetOkExists("weight"); ok {
		server["Weight"] = v
	}
	if v, ok := d.GetOk("type"); ok {
		server["Type"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		server["Description"] = v
	}

	serverMaps = append(serverMaps, server)
	serverEntries, err := convertListMapToJsonString(serverMaps)
	if err != nil {
		return WrapError(err)
	}
	request["BackendServers"] = serverEntries
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"VServerGroupProcessing"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_slb_server_group_server_attachment", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["VServerGroupId"], ":", server["ServerId"], ":", server["Port"]))
	return resourceAlicloudSlbServerGroupServerAttachmentRead(d, meta)
}

func resourceAlicloudSlbServerGroupServerAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slbService := SlbService{client}
	object, err := slbService.DescribeSlbServerGroupServerAttachment(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_slb_server_group_server_attachment slbService.DescribeSlbServerGroupServerAttachment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}

	d.Set("server_group_id", parts[0])
	d.Set("server_id", parts[1])
	d.Set("port", formatInt(parts[2]))
	d.Set("weight", formatInt(object["Weight"]))
	d.Set("type", object["Type"])
	d.Set("description", object["Description"])
	return nil
}

func resourceAlicloudSlbServerGroupServerAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "RemoveVServerGroupBackendServers"
	var response map[string]interface{}
	conn, err := client.NewSlbClient()
	if err != nil {
		return WrapError(err)
	}

	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"RegionId": client.RegionId,
	}

	request["VServerGroupId"] = parts[0]
	serverMaps := make([]map[string]interface{}, 1)
	server := make(map[string]interface{}, 1)
	server["ServerId"] = parts[1]
	server["Port"] = parts[2]

	if v, ok := d.GetOkExists("weight"); ok {
		server["Weight"] = v
	}
	if v, ok := d.GetOk("type"); ok {
		server["Type"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		server["Description"] = v
	}

	serverMaps = append(serverMaps, server)
	serverEntries, err := convertListMapToJsonString(serverMaps)
	if err != nil {
		return WrapError(err)
	}
	request["BackendServers"] = serverEntries

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"VServerGroupProcessing"}) || NeedRetry(err) {
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
