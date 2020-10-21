package alicloud

import (
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/resourcemanager"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudResourceManagerPolicyVersions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudResourceManagerPolicyVersionsRead,
		Schema: map[string]*schema.Schema{
			"policy_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"policy_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Custom", "System"}, false),
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"versions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"create_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_default_version": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudResourceManagerPolicyVersionsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := resourcemanager.CreateListPolicyVersionsRequest()
	request.PolicyName = d.Get("policy_name").(string)
	request.PolicyType = d.Get("policy_type").(string)
	var objects []resourcemanager.PolicyVersion
	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}
	raw, err := client.WithResourcemanagerClient(func(resourcemanagerClient *resourcemanager.Client) (interface{}, error) {
		return resourcemanagerClient.ListPolicyVersions(request)
	})
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_resource_manager_policy_versions", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	response, _ := raw.(*resourcemanager.ListPolicyVersionsResponse)

	for _, item := range response.PolicyVersions.PolicyVersion {
		if len(idsMap) > 0 {
			if _, ok := idsMap[request.PolicyName+":"+item.VersionId]; !ok {
				continue
			}
		}
		objects = append(objects, item)
	}
	ids := make([]string, len(objects))
	s := make([]map[string]interface{}, len(objects))

	for i, object := range objects {
		mapping := map[string]interface{}{
			"create_date":        object.CreateDate,
			"is_default_version": object.IsDefaultVersion,
			"id":                 fmt.Sprintf("%v:%v", request.PolicyName, object.VersionId),
			"version_id":         object.VersionId,
		}
		ids[i] = fmt.Sprintf("%v:%v", request.PolicyName, object.VersionId)
		s[i] = mapping
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("versions", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
