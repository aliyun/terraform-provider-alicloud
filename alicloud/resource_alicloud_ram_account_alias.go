package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudRamAccountAlias() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudRamAccountAliasCreate,
		Read:   resourceAlicloudRamAccountAliasRead,
		Delete: resourceAlicloudRamAccountAliasDelete,

		Schema: map[string]*schema.Schema{
			"account_alias": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateRamAlias,
			},
		},
	}
}

func resourceAlicloudRamAccountAliasCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := ram.CreateSetAccountAliasRequest()
	request.AccountAlias = d.Get("account_alias").(string)

	_, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.SetAccountAlias(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "ram_account_alias", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	d.SetId(request.AccountAlias)
	return resourceAlicloudRamAccountAliasRead(d, meta)
}

func resourceAlicloudRamAccountAliasRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := ram.CreateGetAccountAliasRequest()

	raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.GetAccountAlias(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	response, _ := raw.(*ram.GetAccountAliasResponse)

	d.Set("account_alias", response.AccountAlias)
	return nil
}

func resourceAlicloudRamAccountAliasDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := ram.CreateClearAccountAliasRequest()

	_, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.ClearAccountAlias(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return nil
}
