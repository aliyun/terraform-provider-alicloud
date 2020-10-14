package alicloud

import (
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/kms"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudKmsKeys() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudKmsKeysRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				MinItems: 1,
			},

			"description_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.ValidateRegexp,
			},

			"status": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				// must contain a valid status, expected Enabled, Disabled, PendingDeletion
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

			//Computed value
			"keys": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"arn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"delete_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creator": {
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
	request.RegionId = client.RegionId

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok && len(v.([]interface{})) > 0 {
		for _, i := range v.([]interface{}) {
			if i == nil {
				continue
			}
			idsMap[i.(string)] = i.(string)
		}
	}

	var s []map[string]interface{}
	var ids []string
	var keyIds []string

	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	for true {
		raw, err := client.WithKmsClient(func(kmsClient *kms.Client) (interface{}, error) {
			return kmsClient.ListKeys(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_kms_keys", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*kms.ListKeysResponse)
		for _, key := range response.Keys.Key {
			if len(idsMap) > 0 {
				if _, ok := idsMap[key.KeyId]; ok {
					keyIds = append(keyIds, key.KeyId)
					continue
				}
			} else {
				keyIds = append(keyIds, key.KeyId)
				continue
			}
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

	descriptionRegex, ok := d.GetOk("description_regex")
	var r *regexp.Regexp
	if ok && descriptionRegex.(string) != "" {
		r = regexp.MustCompile(descriptionRegex.(string))
	}
	status, statusOk := d.GetOk("status")
	for _, k := range keyIds {

		request := kms.CreateDescribeKeyRequest()
		request.KeyId = k
		raw, err := client.WithKmsClient(func(kmsClient *kms.Client) (interface{}, error) {
			return kmsClient.DescribeKey(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, k, request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		key, _ := raw.(*kms.DescribeKeyResponse)
		if r != nil && !r.MatchString(key.KeyMetadata.Description) {
			continue
		}
		if statusOk && status != "" && status != key.KeyMetadata.KeyState {
			continue
		}
		mapping := map[string]interface{}{
			"id":            key.KeyMetadata.KeyId,
			"arn":           key.KeyMetadata.Arn,
			"description":   key.KeyMetadata.Description,
			"status":        key.KeyMetadata.KeyState,
			"creation_date": key.KeyMetadata.CreationDate,
			"delete_date":   key.KeyMetadata.DeleteDate,
			"creator":       key.KeyMetadata.Creator,
		}

		s = append(s, mapping)
		ids = append(ids, key.KeyMetadata.KeyId)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("keys", s); err != nil {
		return WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
