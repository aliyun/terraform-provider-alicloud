package alicloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func dataSourceAlicloudStsCallerIdentity() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudStsCallerIdentityRead,
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"role_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"identity_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"principal_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAlicloudStsCallerIdentityRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	resp, err := client.GetCallerIdentity()
	if err != nil {
		return err
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		s := map[string]interface{}{
			"account_id":     resp.AccountId,
			"user_id":        resp.UserId,
			"role_id":        resp.RoleId,
			"arn":            resp.Arn,
			"identity_type ": resp.IdentityType,
			"principal_id":   resp.PrincipalId,
		}
		writeToFile(output.(string), s)
	}
	return nil
}
