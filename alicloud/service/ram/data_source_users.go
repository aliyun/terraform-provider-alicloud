package ram

import (
	"log"
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

func DataSourceUsers() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceUsersRead,
		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"group_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"policy_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validators.StringLenBetween(0, 128),
			},
			"policy_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validators.StringInSlice([]string{"System", "Custom"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			// Computed values
			"users": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func DataSourceUsersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ramService := RamService{client}
	var allUsers []interface{}

	allUsersMap := make(map[string]interface{})
	groupFilterUsersMap := make(map[string]interface{})
	policyFilterUsersMap := make(map[string]interface{})

	var dataMap []map[string]interface{}

	groupName, groupNameOk := d.GetOk("group_name")
	policyName, policyNameOk := d.GetOk("policy_name")
	policyType, policyTypeOk := d.GetOk("policy_type")
	nameRegex, nameRegexOk := d.GetOk("name_regex")
	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}
	if policyTypeOk && !policyNameOk {
		return errs.WrapError(errs.Error("You must set 'policy_name' at one time when you set 'policy_type'."))
	}

	// all users
	request := ramsdk.CreateListUsersRequest()
	request.RegionId = client.RegionId
	request.MaxItems = requests.NewInteger(1000)
	for {
		var raw interface{}
		var err error
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutRead)), func() *retry.RetryError {
			raw, err = client.WithRamClient(func(ramClient *ramsdk.Client) (interface{}, error) {
				return ramClient.ListUsers(request)
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
			return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_ram_users", request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}

		response, _ := raw.(*ramsdk.ListUsersResponse)
		for _, v := range response.Users.User {
			if nameRegexOk {
				r, err := regexp.Compile(nameRegex.(string))
				if err != nil {
					return errs.WrapError(err)
				}
				if !r.MatchString(v.UserName) {
					continue
				}
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[v.UserId]; !ok {
					continue
				}
			}
			allUsersMap[v.UserName] = v
		}
		if !response.IsTruncated {
			break
		}
		request.Marker = response.Marker
	}

	// users for group
	if groupNameOk {
		request := ramsdk.CreateListUsersForGroupRequest()
		request.GroupName = groupName.(string)
		users := []ramsdk.User{}

		var raw interface{}
		var err error
		for {
			wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
			err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutRead)), func() *retry.RetryError {
				raw, err = client.WithRamClient(func(ramClient *ramsdk.Client) (interface{}, error) {
					return ramClient.ListUsersForGroup(request)
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
				return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_ram_users", request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
			}

			response, _ := raw.(*ramsdk.ListUsersForGroupResponse)
			users = append(users, response.Users.User...)
			if !response.IsTruncated {
				break
			}
			request.Marker = response.Marker
		}

		for _, v := range users {
			groupFilterUsersMap[v.UserName] = v
		}

		if len(groupFilterUsersMap) > 0 {
			dataMap = append(dataMap, groupFilterUsersMap)
		}
	}

	// users which attach with this policy
	if policyNameOk {
		pType := "System"
		if policyTypeOk {
			pType = policyType.(string)
		}
		request := ramsdk.CreateListEntitiesForPolicyRequest()
		request.PolicyName = policyName.(string)
		request.PolicyType = pType

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
			return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_ram_users", request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}

		response, _ := raw.(*ramsdk.ListEntitiesForPolicyResponse)

		for _, v := range response.Users.User {
			policyFilterUsersMap[v.UserName] = v
		}

		if len(policyFilterUsersMap) > 0 {
			dataMap = append(dataMap, policyFilterUsersMap)
		}
	}

	// GetIntersection of each map
	allUsers = ramService.GetIntersection(dataMap, allUsersMap, groupNameOk, policyNameOk)

	return ramUsersDescriptionAttributes(d, allUsers)
}

func ramUsersDescriptionAttributes(d *schema.ResourceData, users []interface{}) error {
	var ids []string
	var names []string
	var s []map[string]interface{}
	for _, v := range users {
		user, ok := v.(ramsdk.User)
		if !ok {
			return errs.WrapError(errs.Error("wrong interface convince"))
		}
		mapping := map[string]interface{}{
			"id":          user.UserId,
			"name":        user.UserName,
			"create_date": user.CreateDate,
		}
		log.Printf("[DEBUG] alicloud_ram_users - adding user: %v", mapping)
		ids = append(ids, user.UserId)
		names = append(names, user.UserName)
		s = append(s, mapping)
	}

	d.SetId(helper.DataResourceIdHash(ids))

	if err := d.Set("users", s); err != nil {
		return errs.WrapError(err)
	}

	if err := d.Set("ids", ids); err != nil {
		return errs.WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return errs.WrapError(err)
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		helper.WriteToFile(output.(string), s)
	}

	return nil
}
