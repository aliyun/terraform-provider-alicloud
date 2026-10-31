package cs

import (
	"encoding/json"
	"strings"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const MASKED_CONFIG_KEY_PREFIX = "x-ui"
const DatasourceCSKubernetesAddonMetadata = "alicloud_cs_kubernetes_addon_metadata"

func DataSourceCSKubernetesAddonMetadata() *schema.Resource {
	return &schema.Resource{
		Read: dataCSKubernetesAddonMetadataRead,

		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"version": {
				Type:     schema.TypeString,
				Required: true,
			},
			"config_schema": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataCSKubernetesAddonMetadataRead(d *schema.ResourceData, meta interface{}) error {
	clusterId := d.Get("cluster_id").(string)
	name := d.Get("name").(string)
	version := d.Get("version").(string)

	component, err := DescribeClusterAddonMetadata(d, meta)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, DatasourceCSKubernetesAddonMetadata, "DescribeClusterAddonMetadata", err)
	}

	config, err := fetchJsonSchema(component.ConfigSchema)
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, DatasourceCSKubernetesAddonMetadata, "DescribeClusterAddonMetadata", err)
	}

	d.Set("cluster_id", clusterId)
	d.Set("name", name)
	d.Set("version", version)
	d.Set("config_schema", config)

	d.SetId(tea.ToString(helper.Hashcode(clusterId)))
	return nil
}

func DescribeClusterAddonMetadata(d *schema.ResourceData, meta interface{}) (*Component, error) {
	clusterId := d.Get("cluster_id").(string)
	name := d.Get("name").(string)
	version := d.Get("version").(string)

	client, err := meta.(*connectivity.AliyunClient).NewRoaCsClient()
	if err != nil {
		return nil, err
	}
	csClient := CsClient{client}

	return csClient.DescribeCsKubernetesAddonMetadata(clusterId, name, version)
}

func fetchJsonSchema(schema string) (string, error) {
	if schema == "" {
		return "", nil
	}
	var i interface{}
	if err := json.Unmarshal([]byte(schema), &i); err != nil {
		return "", errs.WrapErrorf(err, errs.DefaultErrorMsg, DatasourceCSKubernetesAddonMetadata, "fetchJsonSchema", err)
	}
	if v, ok := i.(map[string]interface{}); ok {
		result, err := json.MarshalIndent(parseNode(v), "", "\t")
		if err != nil {
			return "", errs.WrapErrorf(errs.Error("addon config schema marshal error"), errs.DefaultErrorMsg, DatasourceCSKubernetesAddonMetadata, "fetchJsonSchema")
		}
		return string(result), nil
	}
	return "", errs.WrapErrorf(errs.Error("addon config schema parse error"), errs.DefaultErrorMsg, DatasourceCSKubernetesAddonMetadata, "fetchJsonSchema")
}

func parseNode(p map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range p {
		if strings.HasPrefix(k, MASKED_CONFIG_KEY_PREFIX) {
			continue
		}
		if n, ok := v.(map[string]interface{}); ok {
			r := parseNode(n)
			result[k] = r
		} else {
			result[k] = v
		}
	}
	return result
}
