package alicloud

import (
	"fmt"
	"regexp"

	"github.com/denverdino/aliyungo/kms"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func dataSourceAlicloudKmsKeys() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudKmsKeysRead,

		Schema: map[string]*schema.Schema{
			"ids": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				MinItems: 1,
			},

			"description_regex": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateNameRegex,
			},

			"status": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateKmsKeyStatus,
			},

			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			//Computed value
			"keys": &schema.Schema{
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

	args := &kms.ListKeysArgs{}

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok && len(v.([]interface{})) > 0 {
		for _, i := range v.([]interface{}) {
			idsMap[i.(string)] = i.(string)
		}
	}

	var keyIds []string
	pagination := getPagination(1, 50)
	for true {
		args.Pagination = pagination
		raw, err := client.WithKmsClient(func(kmsClient *kms.Client) (interface{}, error) {
			return kmsClient.ListKeys(args)
		})
		if err != nil {
			return fmt.Errorf("Error ListKeys: %#v", err)
		}
		results, _ := raw.(*kms.ListKeysResponse)
		for _, key := range results.Keys.Key {
			if idsMap != nil {
				if _, ok := idsMap[key.KeyId]; ok {
					keyIds = append(keyIds, key.KeyId)
					continue
				}
			}
			keyIds = append(keyIds, key.KeyId)
		}
		if len(results.Keys.Key) < pagination.PageSize {
			break
		}
		pagination.PageNumber += 1
	}

	if len(keyIds) < 1 {
		return fmt.Errorf("Your query kms keys returned no results. Please change your search criteria and try again.")
	}

	var s []map[string]interface{}
	var ids []string
	descriptionRegex, ok := d.GetOk("description_regex")
	var r *regexp.Regexp
	if ok && descriptionRegex.(string) != "" {
		r = regexp.MustCompile(descriptionRegex.(string))
	}
	status, statusOk := d.GetOk("status")

	for _, k := range keyIds {
		raw, err := client.WithKmsClient(func(kmsClient *kms.Client) (interface{}, error) {
			return kmsClient.DescribeKey(k)
		})
		if err != nil {
			return fmt.Errorf("DescribeKey got an error: %#v", err)
		}
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
		return err
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
