package ram

import (
	"regexp"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	ramsdk "github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceGroups() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceGroupsRead,
		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"user_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validators.StringLenBetween(0, 64),
			},
			"policy_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validators.StringLenBetween(0, 128),
			},
			"policy_type": {
				Type:     schema.TypeString,
				Optional: true,
				// must be ramsdk.System, ramsdk.Custom
				ValidateFunc: validators.StringInSlice([]string{"System", "Custom"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			// Computed values
			"groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"comments": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func DataSourceGroupsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ramService := RamService{client}
	var allGroups []interface{}

	allGroupsMap := make(map[string]interface{})
	userFilterGroupsMap := make(map[string]interface{})
	policyFilterGroupsMap := make(map[string]interface{})

	var dataMap []map[string]interface{}

	userName, userNameOk := d.GetOk("user_name")
	policyName, policyNameOk := d.GetOk("policy_name")
	policyType, policyTypeOk := d.GetOk("policy_type")
	nameRegex, nameRegexOk := d.GetOk("name_regex")

	if policyTypeOk && !policyNameOk {
		return errs.WrapError(errs.Error("You must set 'policy_name' at one time when you set 'policy_type'."))
	}

	// groups filtered by name_regex
	request := ramsdk.CreateListGroupsRequest()
	request.RegionId = client.RegionId
	request.MaxItems = requests.NewInteger(1000)
	for {
		var raw interface{}
		var err error
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutRead)), func() *retry.RetryError {
			raw, err = client.WithRamClient(func(ramClient *ramsdk.Client) (interface{}, error) {
				return ramClient.ListGroups(request)
			})
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
			return nil
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_ram_groups", request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}

		response, _ := raw.(*ramsdk.ListGroupsResponse)
		for _, v := range response.Groups.Group {
			if nameRegexOk {
				r, err := regexp.Compile(nameRegex.(string))
				if err != nil {
					return errs.WrapError(err)
				}
				if !r.MatchString(v.GroupName) {
					continue
				}
			}
			allGroupsMap[v.GroupName] = v
		}
		if !response.IsTruncated {
			break
		}
		request.Marker = response.Marker
	}

	// groups for user
	if userNameOk {
		request := ramsdk.CreateListGroupsForUserRequest()
		request.UserName = userName.(string)

		var raw interface{}
		var err error
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutRead)), func() *retry.RetryError {
			raw, err = client.WithRamClient(func(ramClient *ramsdk.Client) (interface{}, error) {
				return ramClient.ListGroupsForUser(request)
			})
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
			return nil
		})
		if err != nil && !errs.NotFoundError(err) {
			return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_ram_groups", request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}

		response, _ := raw.(*ramsdk.ListGroupsForUserResponse)

		for _, v := range response.Groups.Group {
			userFilterGroupsMap[v.GroupName] = v
		}

		if len(userFilterGroupsMap) > 0 {
			dataMap = append(dataMap, userFilterGroupsMap)
		}
	}

	// groups which attach with this policy
	if policyNameOk {
		pType := "System"
		if policyTypeOk {
			pType = policyType.(string)
		}
		request := ramsdk.CreateListEntitiesForPolicyRequest()
		request.PolicyType = pType
		request.PolicyName = policyName.(string)

		var raw interface{}
		var err error
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutRead)), func() *retry.RetryError {
			raw, err = client.WithRamClient(func(ramClient *ramsdk.Client) (interface{}, error) {
				return ramClient.ListEntitiesForPolicy(request)
			})
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
			return nil
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_ram_groups", request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}

		response, _ := raw.(*ramsdk.ListEntitiesForPolicyResponse)

		for _, v := range response.Groups.Group {
			policyFilterGroupsMap[v.GroupName] = v
		}

		if len(policyFilterGroupsMap) > 0 {
			dataMap = append(dataMap, policyFilterGroupsMap)
		}
	}

	// GetIntersection of each map
	allGroups = ramService.GetIntersection(dataMap, allGroupsMap, userNameOk, policyNameOk)

	return ramGroupsDescriptionAttributes(d, allGroups)
}

func ramGroupsDescriptionAttributes(d *schema.ResourceData, groups []interface{}) error {
	var ids []string
	var s []map[string]interface{}
	for _, v := range groups {
		group := v.(ramsdk.Group)
		mapping := map[string]interface{}{
			"name":     group.GroupName,
			"comments": group.Comments,
		}
		ids = append(ids, v.(ramsdk.Group).GroupName)
		s = append(s, mapping)
	}

	d.SetId(helper.DataResourceIdHash(ids))

	if err := d.Set("names", ids); err != nil {
		return errs.WrapError(err)
	}

	if err := d.Set("groups", s); err != nil {
		return errs.WrapError(err)
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		helper.WriteToFile(output.(string), s)
	}

	return nil
}
