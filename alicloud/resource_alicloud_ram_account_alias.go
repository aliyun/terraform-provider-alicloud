package alicloud

import (
	"fmt"

	"github.com/denverdino/aliyungo/ram"
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

	args := ram.AccountAliasRequest{
		AccountAlias: d.Get("account_alias").(string),
	}

	_, err := client.WithRamClient(func(ramClient ram.RamClientInterface) (interface{}, error) {
		return ramClient.SetAccountAlias(args)
	})
	if err != nil {
		return fmt.Errorf("SetAccountAlias got an error: %#v", err)
	}

	d.SetId(args.AccountAlias)
	return resourceAlicloudRamAccountAliasRead(d, meta)
}

func resourceAlicloudRamAccountAliasRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	raw, err := client.WithRamClient(func(ramClient ram.RamClientInterface) (interface{}, error) {
		return ramClient.GetAccountAlias()
	})
	if err != nil {
		return fmt.Errorf("GetAccountAlias got an error: %#v", err)
	}
	response, _ := raw.(ram.AccountAliasResponse)

	d.Set("account_alias", response.AccountAlias)
	return nil
}

func resourceAlicloudRamAccountAliasDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	_, err := client.WithRamClient(func(ramClient ram.RamClientInterface) (interface{}, error) {
		return ramClient.ClearAccountAlias()
	})
	if err != nil {
		return fmt.Errorf("ClearAccountAlias got an error: %#v", err)
	}
	return nil
}
