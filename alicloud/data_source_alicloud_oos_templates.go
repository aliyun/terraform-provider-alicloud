package alicloud

import (
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/oos"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudOosTemplates() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudOosTemplatesRead,
		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"category": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"created_by": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"created_date": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"created_date_after": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"has_trigger": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"share_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Private", "Public"}, false),
			},
			"sort_field": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"CreatedDate", "Popularity", "TemplateName", "TotalExecutionCount"}, false),
				Default:      "TotalExecutionCount",
			},
			"sort_order": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Ascending", "Descending"}, false),
				Default:      "Descending",
			},
			"tags": tagsSchema(),
			"template_format": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"JSON", "YAML"}, false),
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"template_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"templates": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"category": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_by": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"has_trigger": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"share_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
						},
						"template_format": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"template_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"template_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"template_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"template_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updated_by": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updated_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudOosTemplatesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := oos.CreateListTemplatesRequest()
	if v, ok := d.GetOk("category"); ok {
		request.Category = v.(string)
	}
	if v, ok := d.GetOk("created_by"); ok {
		request.CreatedBy = v.(string)
	}
	if v, ok := d.GetOk("created_date"); ok {
		request.CreatedDateBefore = v.(string)
	}
	if v, ok := d.GetOk("created_date_after"); ok {
		request.CreatedDateAfter = v.(string)
	}
	if v, ok := d.GetOkExists("has_trigger"); ok {
		request.HasTrigger = requests.NewBoolean(v.(bool))
	}
	request.RegionId = client.RegionId
	if v, ok := d.GetOk("share_type"); ok {
		request.ShareType = v.(string)
	}
	if v, ok := d.GetOk("sort_field"); ok {
		request.SortField = v.(string)
	}
	if v, ok := d.GetOk("sort_order"); ok {
		request.SortOrder = v.(string)
	}
	if v, ok := d.GetOk("tags"); ok {
		request.Tags = v.(map[string]interface{})
	}
	if v, ok := d.GetOk("template_format"); ok {
		request.TemplateFormat = v.(string)
	}
	if v, ok := d.GetOk("template_type"); ok {
		request.TemplateType = v.(string)
	}
	request.MaxResults = requests.NewInteger(PageSizeLarge)
	var objects []oos.Template
	var templateNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		templateNameRegex = r
	}

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}
	var response *oos.ListTemplatesResponse
	for {
		raw, err := client.WithOosClient(func(oosClient *oos.Client) (interface{}, error) {
			return oosClient.ListTemplates(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_oos_templates", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw)
		response, _ = raw.(*oos.ListTemplatesResponse)

		for _, item := range response.Templates {
			if templateNameRegex != nil {
				if !templateNameRegex.MatchString(item.TemplateName) {
					continue
				}
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[item.TemplateName]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}
		NextToken := response.NextToken
		if NextToken == "" {
			break
		}
		request.NextToken = NextToken
	}
	ids := make([]string, len(objects))
	s := make([]map[string]interface{}, len(objects))
	for i, object := range objects {
		mapping := map[string]interface{}{
			"category":         object.Category,
			"created_by":       object.CreatedBy,
			"created_date":     object.CreatedDate,
			"description":      object.Description,
			"has_trigger":      object.HasTrigger,
			"share_type":       object.ShareType,
			"template_format":  object.TemplateFormat,
			"template_id":      object.TemplateId,
			"id":               object.TemplateName,
			"template_name":    object.TemplateName,
			"template_type":    object.TemplateType,
			"template_version": object.TemplateVersion,
			"updated_by":       object.UpdatedBy,
			"updated_date":     object.UpdatedDate,
		}
		ids[i] = object.TemplateName
		s[i] = mapping
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("templates", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
