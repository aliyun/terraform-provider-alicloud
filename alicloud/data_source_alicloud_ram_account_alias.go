package alicloud

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceAlicloudRamAccountAlias() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudRamAccountAliasRead,

		Schema: map[string]*schema.Schema{
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"account_alias": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAlicloudRamAccountAliasRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ramconn

	resp, err := conn.GetAccountAlias()
	if err != nil {
		return err
	}
	d.SetId(resp.AccountAlias)
	d.Set("account_alias", resp.AccountAlias)

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		s := map[string]interface{}{"account_alias": resp.AccountAlias}
		writeToFile(output.(string), s)
	}
	return nil
}
