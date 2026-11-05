package threatdetection

import (
	"strings"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAliCloudThreatDetectionRdDefaultSyncList() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudThreatDetectionRdDefaultSyncListRead,
		Schema: map[string]*schema.Schema{
			"folder_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceAliCloudThreatDetectionRdDefaultSyncListRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	threatDetectionServiceV2 := ServiceV2{client}

	// The default synchronization list is an account-level singleton, so the
	// data source id is the Alibaba Cloud account id, the same as the resource.
	accountId, err := client.AccountId()
	if err != nil {
		return errs.WrapError(err)
	}

	object, err := threatDetectionServiceV2.DescribeThreatDetectionRdDefaultSyncList(accountId)
	if err != nil {
		return errs.WrapError(err)
	}

	d.SetId(accountId)

	folderIds := make([]interface{}, 0)
	if v, ok := object["FolderIds"]; ok && v != nil {
		if str, ok := v.(string); ok {
			str = strings.TrimSpace(str)
			if str != "" {
				for _, id := range strings.Split(str, helper.CommaSeparated) {
					id = strings.TrimSpace(id)
					if id != "" {
						folderIds = append(folderIds, id)
					}
				}
			}
		}
	}
	if err := d.Set("folder_ids", folderIds); err != nil {
		return errs.WrapError(err)
	}

	return nil
}
