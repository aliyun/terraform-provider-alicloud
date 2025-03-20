package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudSlbServerGroupServerAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudSlbServerGroupServerAttachmentCreate,
		Read:   resourceAliCloudSlbServerGroupServerAttachmentRead,
		Delete: resourceAliCloudSlbServerGroupServerAttachmentDelete,
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
				ValidateFunc: IntBetween(1, 65535),
			},
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"ecs", "eni", "eci"}, false),
			},
			"weight": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: IntBetween(0, 100),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudSlbServerGroupServerAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "AddVServerGroupBackendServers"
	request := make(map[string]interface{})
	var err error

	request["RegionId"] = client.RegionId
	request["VServerGroupId"] = d.Get("server_group_id")

	serverMaps := make([]map[string]interface{}, 1)
	server := make(map[string]interface{}, 2)
	server["ServerId"] = d.Get("server_id")
	server["Port"] = d.Get("port")

	if v, ok := d.GetOk("type"); ok {
		server["Type"] = v
	}

	if v, ok := d.GetOkExists("weight"); ok {
		server["Weight"] = v
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
		response, err = client.RpcPost("Slb", "2014-05-15", action, nil, request, false)
		if err != nil {
			if IsExpectedErrors(err, []string{"VServerGroupProcessing", "ServiceIsStopping", "ServiceIsConfiguring"}) || NeedRetry(err) {
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

	d.SetId(fmt.Sprintf("%v:%v:%v", request["VServerGroupId"], server["ServerId"], server["Port"]))

	return resourceAliCloudSlbServerGroupServerAttachmentRead(d, meta)
}

func resourceAliCloudSlbServerGroupServerAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slbService := SlbService{client}

	object, err := slbService.DescribeSlbServerGroupServerAttachment(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
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
	d.Set("server_id", object["ServerId"])
	d.Set("port", formatInt(object["Port"]))
	d.Set("type", object["Type"])
	d.Set("weight", formatInt(object["Weight"]))
	d.Set("description", object["Description"])

	return nil
}

func resourceAliCloudSlbServerGroupServerAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "RemoveVServerGroupBackendServers"
	var response map[string]interface{}
	var err error

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

	if v, ok := d.GetOk("type"); ok {
		server["Type"] = v
	}

	if v, ok := d.GetOkExists("weight"); ok {
		server["Weight"] = v
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
		response, err = client.RpcPost("Slb", "2014-05-15", action, nil, request, false)
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
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
