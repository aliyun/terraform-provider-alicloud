package alicloud

import (
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/kms"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func dataSourceAlicloudKmsKeys() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudKmsKeysRead,
		Schema: map[string]*schema.Schema{
			"description_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.ValidateRegexp,
			},

			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(Enabled),
					string(Disabled),
					string(PendingDeletion),
				}, false),
			},

			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"keys": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"arn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"automatic_rotation": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creator": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"delete_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"key_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"key_spec": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"key_state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"key_usage": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_rotation_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"material_expire_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"next_rotation_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"origin": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"primary_key_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"protection_level": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rotation_interval": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudKmsKeysRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := kms.CreateListKeysRequest()
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	var objects []kms.KeyMetadata
	descriptionRegex, ok := d.GetOk("description_regex")
	var r *regexp.Regexp
	if ok && descriptionRegex.(string) != "" {
		r = regexp.MustCompile(descriptionRegex.(string))
	}
	status, statusOk := d.GetOk("status")

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			idsMap[vv.(string)] = vv.(string)
		}
	}
	for {
		raw, err := client.WithKmsClient(func(kmsClient *kms.Client) (interface{}, error) {
			return kmsClient.ListKeys(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_kms_keys", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw)
		response, _ := raw.(*kms.ListKeysResponse)

		for _, item := range response.Keys.Key {
			if len(idsMap) > 0 {
				if _, ok := idsMap[item.KeyId]; !ok {
					continue
				}
			}
			request := kms.CreateDescribeKeyRequest()
			request.KeyId = item.KeyId
			raw, err := client.WithKmsClient(func(kmsClient *kms.Client) (interface{}, error) {
				return kmsClient.DescribeKey(request)
			})
			if err != nil {
				return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_kms_keys", request.GetActionName(), AlibabaCloudSdkGoERROR)
			}
			project, _ := raw.(*kms.DescribeKeyResponse)
			if r != nil {
				if !r.MatchString(project.KeyMetadata.Description) {
					continue
				}
			}
			if status != nil {
				if statusOk && status != project.KeyMetadata.KeyState {
					continue
				}
			}
			objects = append(objects, project.KeyMetadata)
		}
		if len(response.Keys.Key) < PageSizeLarge {
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
			"arn":                  object.Arn,
			"automatic_rotation":   object.AutomaticRotation,
			"creation_date":        object.CreationDate,
			"creator":              object.Creator,
			"delete_date":          object.DeleteDate,
			"description":          object.Description,
			"id":                   object.KeyId,
			"key_id":               object.KeyId,
			"key_spec":             object.KeySpec,
			"key_state":            object.KeyState,
			"key_usage":            object.KeyUsage,
			"last_rotation_date":   object.LastRotationDate,
			"material_expire_time": object.MaterialExpireTime,
			"next_rotation_date":   object.NextRotationDate,
			"origin":               object.Origin,
			"primary_key_version":  object.PrimaryKeyVersion,
			"protection_level":     object.ProtectionLevel,
			"rotation_interval":    object.RotationInterval,
		}
		ids[i] = object.KeyId
		s[i] = mapping
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("keys", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
