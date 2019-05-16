package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudRamAccessKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudRamAccessKeyCreate,
		Read:   resourceAlicloudRamAccessKeyRead,
		Update: resourceAlicloudRamAccessKeyUpdate,
		Delete: resourceAlicloudRamAccessKeyDelete,

		Schema: map[string]*schema.Schema{
			"user_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"secret_file": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      Active,
				ValidateFunc: validateRamAKStatus,
			},
		},
	}
}

func resourceAlicloudRamAccessKeyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ramService := RamService{client}
	request := ram.CreateCreateAccessKeyRequest()
	if v, ok := d.GetOk("user_name"); ok && v.(string) != "" {
		request.UserName = v.(string)
	}

	raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.CreateAccessKey(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ram_access_key", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	response, _ := raw.(*ram.CreateAccessKeyResponse)

	// create a secret_file and write access key to it.
	if output, ok := d.GetOk("secret_file"); ok && output != nil {
		writeToFile(output.(string), response.AccessKey)
	}
	d.SetId(response.AccessKey.AccessKeyId)
	err = ramService.WaitForRamAccessKey(d.Id(), request.UserName, Active, DefaultTimeout)
	if err != nil {
		return WrapError(err)
	}
	return resourceAlicloudRamAccessKeyUpdate(d, meta)
}

func resourceAlicloudRamAccessKeyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := ram.CreateUpdateAccessKeyRequest()
	request.UserAccessKeyId = d.Id()
	request.Status = d.Get("status").(string)

	if v, ok := d.GetOk("user_name"); ok && v.(string) != "" {
		request.UserName = v.(string)
	}

	if d.HasChange("status") {
		raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.UpdateAccessKey(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw)
	}
	return resourceAlicloudRamAccessKeyRead(d, meta)
}

func resourceAlicloudRamAccessKeyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ramservice := RamService{client}
	userName := ""
	if v, ok := d.GetOk("user_name"); ok && v.(string) != "" {
		userName = v.(string)
	}
	object, err := ramservice.DescribeRamAccessKey(d.Id(), userName)
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("status", object.Status)
	return nil
}

func resourceAlicloudRamAccessKeyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ramService := RamService{client}
	request := ram.CreateDeleteAccessKeyRequest()
	request.UserAccessKeyId = d.Id()

	if v, ok := d.GetOk("user_name"); ok && v.(string) != "" {
		request.UserName = v.(string)
	}

	raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.DeleteAccessKey(request)
	})
	if err != nil {
		if RamEntityNotExist(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, request.UserName, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	return WrapError(ramService.WaitForRamAccessKey(d.Id(), request.UserName, Deleted, DefaultTimeoutMedium))

}
