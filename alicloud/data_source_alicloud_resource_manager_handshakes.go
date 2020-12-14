package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/resourcemanager"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudResourceManagerHandshakes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudResourceManagerHandshakesRead,
		Schema: map[string]*schema.Schema{
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
			"handshakes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"expire_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"handshake_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"master_account_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"master_account_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"modify_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"note": {
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
						"target_entity": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"target_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudResourceManagerHandshakesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := resourcemanager.CreateListHandshakesForResourceDirectoryRequest()
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	var objects []resourcemanager.Handshake
	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}
	for {
		raw, err := client.WithResourcemanagerClient(func(resourcemanagerClient *resourcemanager.Client) (interface{}, error) {
			return resourcemanagerClient.ListHandshakesForResourceDirectory(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_resource_manager_handshakes", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw)
		response, _ := raw.(*resourcemanager.ListHandshakesForResourceDirectoryResponse)

		for _, item := range response.Handshakes.Handshake {
			if len(idsMap) > 0 {
				if _, ok := idsMap[item.HandshakeId]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}
		if len(response.Handshakes.Handshake) < PageSizeLarge {
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
			"expire_time":           object.ExpireTime,
			"id":                    object.HandshakeId,
			"handshake_id":          object.HandshakeId,
			"master_account_id":     object.MasterAccountId,
			"master_account_name":   object.MasterAccountName,
			"modify_time":           object.ModifyTime,
			"note":                  object.Note,
			"resource_directory_id": object.ResourceDirectoryId,
			"status":                object.Status,
			"target_entity":         object.TargetEntity,
			"target_type":           object.TargetType,
		}
		ids[i] = object.HandshakeId
		s[i] = mapping
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("handshakes", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
