// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAliCloudPrivatelinkVpcEndpointServiceUsers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudPrivatelinkVpcEndpointServiceUsersRead,
		Schema: map[string]*schema.Schema{
			"service_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"user_list_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Users", "UserARNs"}, false),
			},
			"ids": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"users": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_arn": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAliCloudPrivatelinkVpcEndpointServiceUsersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListVpcEndpointServiceUsers"
	request := make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ServiceId"] = d.Get("service_id")
	if v, ok := d.GetOk("user_id"); ok {
		request["UserId"] = v
	}
	request["MaxResults"] = PageSizeLarge

	userListType := ""
	if v, ok := d.GetOk("user_list_type"); ok {
		userListType = v.(string)
	}

	// The whitelist consists of account ID entries (Users) and ARN entries (UserARNs).
	// Querying with UserListType=Users or with a UserId filter returns account ID
	// entries; querying with UserListType=UserARNs returns ARN entries. When neither
	// user_id nor user_list_type is specified, both kinds are collected.
	queryUsers := userListType == "" || userListType == "Users"
	queryUserArns := userListType == "UserARNs" || (userListType == "" && request["UserId"] == nil)

	listWhitelist := func(listType string) ([]map[string]interface{}, error) {
		req := make(map[string]interface{})
		for k, v := range request {
			req[k] = v
		}
		if listType != "" {
			req["UserListType"] = listType
		}
		var objects []map[string]interface{}
		var response map[string]interface{}
		var err error
		for {
			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			response, err = client.RpcPost("Privatelink", "2020-04-15", action, nil, req, true)
			if err != nil {
				return nil, WrapErrorf(err, DataDefaultErrorMsg, "alicloud_privatelink_vpc_endpoint_service_users", action, AlibabaCloudSdkGoERROR)
			}
			addDebug(action, response, req)

			jsonPath := "$.Users"
			if listType == "UserARNs" {
				jsonPath = "$.UserARNs"
			}
			resp, err := jsonpath.Get(jsonPath, response)
			if err != nil {
				return nil, WrapErrorf(err, FailedGetAttributeMsg, action, jsonPath, response)
			}
			result, _ := resp.([]interface{})
			for _, v := range result {
				item := v.(map[string]interface{})
				objects = append(objects, item)
			}
			if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
				req["NextToken"] = nextToken
			} else {
				break
			}
		}
		return objects, nil
	}

	serviceId := fmt.Sprint(request["ServiceId"])
	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)

	userIds := make(map[string]bool)
	if queryUsers {
		objects, err := listWhitelist("")
		if err != nil {
			return err
		}
		for _, object := range objects {
			userId := formatPrivatelinkUserId(object["UserId"])
			id := fmt.Sprint(serviceId, ":", userId)
			mapping := map[string]interface{}{
				"id":      id,
				"user_id": userId,
			}
			userIds[userId] = true
			ids = append(ids, id)
			s = append(s, mapping)
		}
	}

	if queryUserArns {
		objects, err := listWhitelist("UserARNs")
		if err != nil {
			return err
		}
		for _, object := range objects {
			userArn := fmt.Sprint(object["UserARN"])
			// An account ID whitelist entry also shows up in the UserARNs view as
			// "acs:ram:*:<uid>:*". When both views are collected, skip such ARN forms
			// of entries already listed by user_id so each entry appears exactly once.
			if queryUsers {
				if uid := parsePrivatelinkAccountArnUserId(userArn); uid != "" && userIds[uid] {
					continue
				}
			}
			id := fmt.Sprint(serviceId, "::", userArn)
			mapping := map[string]interface{}{
				"id":       id,
				"user_arn": userArn,
			}
			ids = append(ids, id)
			s = append(s, mapping)
		}
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("users", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
