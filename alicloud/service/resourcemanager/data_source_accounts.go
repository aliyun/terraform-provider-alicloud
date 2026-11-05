package resourcemanager

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	tfTags "github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceAccounts() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceAccountsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validators.StringInSlice([]string{"CreateCancelled", "CreateExpired", "CreateFailed", "CreateSuccess", "CreateVerifying", "InviteSuccess", "PromoteCancelled", "PromoteExpired", "PromoteFailed", "PromoteSuccess", "PromoteVerifying"}, false),
			},
			"tags": tfTags.TagsSchemaForceNew(),
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"accounts": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"account_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"display_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"folder_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_directory_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"account_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"payer_account_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"join_method": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"join_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"modify_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"deletion_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func DataSourceAccountsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListAccounts"
	request := make(map[string]interface{})
	helper.SetPagingRequest(d, request, helper.PageSizeLarge)
	request["IncludeTags"] = true

	if v, ok := d.GetOk("tags"); ok {
		tagsMap := tfTags.ConvertTags(v.(map[string]interface{}))
		request = tfTags.ExpandTagsToMap(request, tagsMap)
	}

	status, statusOk := d.GetOk("status")

	var objects []map[string]interface{}
	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	var response map[string]interface{}
	var err error

	for {
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("ResourceDirectoryMaster", "2022-04-19", action, nil, request, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)

		if err != nil {
			return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_resource_manager_accounts", action, errs.AlibabaCloudSdkGoERROR)
		}

		resp, err := jsonpath.Get("$.Accounts.Account", response)
		if err != nil {
			return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.Accounts.Account", response)
		}

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["AccountId"])]; !ok {
					continue
				}
			}
			if statusOk && status.(string) != "" && status.(string) != item["Status"].(string) {
				continue
			}

			objects = append(objects, item)
		}

		if len(result) < helper.PageSizeLarge {
			break
		}

		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"id":                    fmt.Sprint(object["AccountId"]),
			"account_id":            fmt.Sprint(object["AccountId"]),
			"display_name":          object["DisplayName"],
			"type":                  object["Type"],
			"folder_id":             object["FolderId"],
			"resource_directory_id": object["ResourceDirectoryId"],
			"status":                object["Status"],
			"join_method":           object["JoinMethod"],
			"join_time":             object["JoinTime"],
			"modify_time":           object["ModifyTime"],
			"deletion_status":       object["DeletionStatus"],
		}

		if v, ok := object["Tags"]; ok {
			tags := v.(map[string]interface{})
			if tagMaps, ok := tags["Tag"]; ok {
				mapping["tags"] = tfTags.TagsToMap(tagMaps)
			}
		}

		ids = append(ids, fmt.Sprint(mapping["id"]))

		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			s = append(s, mapping)
			continue
		}

		id := fmt.Sprint(object["AccountId"])
		resourceManagerService := ResourcemanagerService{client}

		resourceManagerAccountDetail, err := resourceManagerService.DescribeResourceManagerAccount(id)
		if err != nil {
			return errs.WrapError(err)
		}

		mapping["account_name"] = resourceManagerAccountDetail["AccountName"]

		payerForAccountDetail, err := resourceManagerService.GetPayerForAccount(id)
		if err != nil {
			return errs.WrapError(err)
		}

		mapping["payer_account_id"] = payerForAccountDetail["PayerAccountId"]

		s = append(s, mapping)
	}

	d.SetId(helper.DataResourceIdHash(ids))

	if err := d.Set("ids", ids); err != nil {
		return errs.WrapError(err)
	}

	if err := d.Set("accounts", s); err != nil {
		return errs.WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		helper.WriteToFile(output.(string), s)
	}

	return nil
}
