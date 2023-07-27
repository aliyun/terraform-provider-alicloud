package alicloud

import (
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudRamAccountAlias() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudRamAccountAliasCreate,
		Read:   resourceAlicloudRamAccountAliasRead,
		Delete: resourceAlicloudRamAccountAliasDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"account_alias": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudRamAccountAliasCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := ram.CreateSetAccountAliasRequest()
	request.RegionId = client.RegionId
	request.AccountAlias = d.Get("account_alias").(string)

	var err error
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.SetAccountAlias(request)
		})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ram_account_alias", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	d.SetId(request.AccountAlias)
	return resourceAlicloudRamAccountAliasRead(d, meta)
}

func resourceAlicloudRamAccountAliasRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ramService := RamService{client}

	object, err := ramService.DescribeRamAccountAlias(d.Id())
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapError(err)
	}
	d.Set("account_alias", object.AccountAlias)
	return nil
}

func resourceAlicloudRamAccountAliasDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := ram.CreateClearAccountAliasRequest()
	request.RegionId = client.RegionId

	var err error
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.ClearAccountAlias(request)
		})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	return nil
}
