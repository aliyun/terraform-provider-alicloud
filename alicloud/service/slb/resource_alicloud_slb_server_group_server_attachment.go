package slb

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceServerGroupServerAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceServerGroupServerAttachmentCreate,
		Read:   resourceServerGroupServerAttachmentRead,
		Delete: resourceServerGroupServerAttachmentDelete,
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
				ValidateFunc: validators.IntBetween(1, 65535),
			},
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validators.StringInSlice([]string{"ecs", "eni", "eci"}, false),
			},
			"weight": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validators.IntBetween(0, 100),
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

func resourceServerGroupServerAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
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

	serverEntries, err := convert.ConvertListMapToJsonString(serverMaps)
	if err != nil {
		return errs.WrapError(err)
	}

	request["BackendServers"] = serverEntries

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("Slb", "2014-05-15", action, nil, request, false)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"VServerGroupProcessing", "ServiceIsStopping", "ServiceIsConfiguring"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_slb_server_group_server_attachment", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v:%v", request["VServerGroupId"], server["ServerId"], server["Port"]))

	return resourceServerGroupServerAttachmentRead(d, meta)
}

func resourceServerGroupServerAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slbService := NewService(client)

	object, err := slbService.DescribeSlbServerGroupServerAttachment(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_slb_server_group_server_attachment slbService.DescribeSlbServerGroupServerAttachment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	parts, err := helper.ParseResourceId(d.Id(), 3)
	if err != nil {
		return errs.WrapError(err)
	}

	d.Set("server_group_id", parts[0])
	d.Set("server_id", object["ServerId"])
	portVal, err := convert.FormatInt(object["Port"])
	if err != nil {
		return err
	}
	d.Set("port", portVal)
	d.Set("type", object["Type"])
	weightVal, err := convert.FormatInt(object["Weight"])
	if err != nil {
		return err
	}
	d.Set("weight", weightVal)
	d.Set("description", object["Description"])

	return nil
}

func resourceServerGroupServerAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "RemoveVServerGroupBackendServers"
	var response map[string]interface{}
	var err error

	parts, err := helper.ParseResourceId(d.Id(), 3)
	if err != nil {
		return errs.WrapError(err)
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

	serverEntries, err := convert.ConvertListMapToJsonString(serverMaps)
	if err != nil {
		return errs.WrapError(err)
	}

	request["BackendServers"] = serverEntries

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("Slb", "2014-05-15", action, nil, request, false)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"VServerGroupProcessing"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		if errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	return nil
}
