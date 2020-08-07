package alicloud

import (
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/resourcemanager"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudResourceManagerPolicyAttachments() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudResourceManagerPolicyAttachmentsRead,
		Schema: map[string]*schema.Schema{
			"language": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"policy_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"policy_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Custom", "System"}, false),
			},
			"principal_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"principal_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"IMSUser", "IMSGroup", "ServiceRole"}, false),
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
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
			"attachments": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"attach_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policy_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policy_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"principal_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"principal_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudResourceManagerPolicyAttachmentsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := resourcemanager.CreateListPolicyAttachmentsRequest()
	if v, ok := d.GetOk("language"); ok {
		request.Language = v.(string)
	}

	if v, ok := d.GetOk("policy_name"); ok {
		request.PolicyName = v.(string)
	}

	if v, ok := d.GetOk("principal_name"); ok {
		request.PrincipalName = v.(string)
	}

	if v, ok := d.GetOk("principal_type"); ok {
		request.PrincipalType = v.(string)
	}

	if v, ok := d.GetOk("resource_group_id"); ok {
		request.ResourceGroupId = v.(string)
	}
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	var objects []resourcemanager.PolicyAttachment
	var response *resourcemanager.ListPolicyAttachmentsResponse
	for {
		raw, err := client.WithResourcemanagerClient(func(resourcemanagerClient *resourcemanager.Client) (interface{}, error) {
			return resourcemanagerClient.ListPolicyAttachments(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_resource_manager_policy_attachments", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw)
		response, _ = raw.(*resourcemanager.ListPolicyAttachmentsResponse)

		for _, item := range response.PolicyAttachments.PolicyAttachment {
			objects = append(objects, item)
		}
		if len(response.PolicyAttachments.PolicyAttachment) < PageSizeLarge {
			break
		}

		page, err := getNextpageNumber(request.PageNumber)
		if err != nil {
			return WrapError(err)
		}
		request.PageNumber = page
	}
	ids := make([]string, len(objects))
	s := make([]map[string]interface{}, len(objects))
	for i, object := range objects {
		mapping := map[string]interface{}{
			"id":                fmt.Sprintf("%v:%v:%v:%v", object.PolicyName, object.PolicyType, object.PrincipalName, object.PrincipalType),
			"attach_date":       object.AttachDate,
			"description":       object.Description,
			"policy_name":       object.PolicyName,
			"policy_type":       object.PolicyType,
			"principal_name":    object.PrincipalName,
			"principal_type":    object.PrincipalType,
			"resource_group_id": object.ResourceGroupId,
		}
		ids[i] = fmt.Sprintf("%v:%v:%v:%v", object.PolicyName, object.PolicyType, object.PrincipalName, object.PrincipalType)
		s[i] = mapping
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("attachments", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
