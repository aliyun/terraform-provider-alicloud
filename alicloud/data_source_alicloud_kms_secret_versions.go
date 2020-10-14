package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/kms"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudKmsSecretVersions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudKmsSecretVersionsRead,
		Schema: map[string]*schema.Schema{
			"include_deprecated": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"secret_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"version_stage": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
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
						"secret_data": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"secret_data_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"secret_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version_stages": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudKmsSecretVersionsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := kms.CreateListSecretVersionIdsRequest()

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok && len(v.([]interface{})) > 0 {
		for _, i := range v.([]interface{}) {
			if i == nil {
				continue
			}
			idsMap[i.(string)] = i.(string)
		}
	}

	if v, ok := d.GetOk("include_deprecated"); ok {
		request.IncludeDeprecated = v.(string)
	}

	VersionStage, okStage := d.GetOk("version_stage")

	request.SecretName = d.Get("secret_name").(string)
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)

	var ids []string
	var objects []kms.VersionId
	var response *kms.ListSecretVersionIdsResponse
	for {
		raw, err := client.WithKmsClient(func(kmsClient *kms.Client) (interface{}, error) {
			return kmsClient.ListSecretVersionIds(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_kms_secret_versions", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw)
		response, _ = raw.(*kms.ListSecretVersionIdsResponse)

		for _, item := range response.VersionIds.VersionId {
			if len(idsMap) > 0 {
				if _, ok := idsMap[item.VersionId]; !ok {
					continue
				}
			}
			if okStage && VersionStage.(string) != "" {
				hasVersionStage := false
				for _, VStage := range item.VersionStages.VersionStage {
					if VStage == VersionStage {
						hasVersionStage = true
						break
					}
				}
				if !hasVersionStage {
					continue
				}
			}

			objects = append(objects, item)
		}
		if len(response.VersionIds.VersionId) < PageSizeLarge {
			break
		}

		page, err := getNextpageNumber(request.PageNumber)
		if err != nil {
			return WrapError(err)
		}
		request.PageNumber = page
	}
	s := make([]map[string]interface{}, len(objects))
	for i, object := range objects {
		mapping := map[string]interface{}{
			"secret_name":    response.SecretName,
			"version_id":     object.VersionId,
			"version_stages": object.VersionStages.VersionStage,
		}

		ids = append(ids, object.VersionId)
		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			s[i] = mapping
			continue
		}
		request := kms.CreateGetSecretValueRequest()
		request.RegionId = client.RegionId
		request.VersionId = object.VersionId
		if okStage && VersionStage.(string) != "" {
			request.VersionStage = VersionStage.(string)
		}
		request.SecretName = d.Get("secret_name").(string)
		raw, err := client.WithKmsClient(func(kmsClient *kms.Client) (interface{}, error) {
			return kmsClient.GetSecretValue(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_kms_secret_versions", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		responseGet, _ := raw.(*kms.GetSecretValueResponse)
		mapping["secret_data"] = responseGet.SecretData
		mapping["secret_data_type"] = responseGet.SecretDataType
		s[i] = mapping
	}

	d.SetId(dataResourceIdHash(ids))

	if err := d.Set("versions", s); err != nil {
		return WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
