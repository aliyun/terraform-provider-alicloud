package cloudmonitorservice

import (
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlicloudCmsAlarmContacts() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCmsAlarmContactsRead,
		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"chanel_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"chanel_value": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"contacts": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"alarm_contact_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"channels_aliim": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"channels_ding_web_hook": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"channels_mail": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"channels_sms": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"channels_state_aliim": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"channels_state_ding_web_hook": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"channels_state_mail": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"channels_status_sms": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"contact_groups": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"describe": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"lang": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudCmsAlarmContactsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := cms.CreateDescribeContactListRequest()
	if v, ok := d.GetOk("chanel_type"); ok {
		request.ChanelType = v.(string)
	}
	if v, ok := d.GetOk("chanel_value"); ok {
		request.ChanelValue = v.(string)
	}
	request.PageSize = requests.NewInteger(helper.PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	var objects []cms.Contact
	var alarmContactNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return errs.WrapError(err)
		}
		alarmContactNameRegex = r
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
	var response *cms.DescribeContactListResponse
	for {
		raw, err := client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
			return cmsClient.DescribeContactList(request)
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_cms_alarm_contacts", request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}
		logs.AddDebug(request.GetActionName(), raw)
		response, _ = raw.(*cms.DescribeContactListResponse)

		for _, item := range response.Contacts.Contact {
			if alarmContactNameRegex != nil {
				if !alarmContactNameRegex.MatchString(item.Name) {
					continue
				}
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[item.Name]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}
		if len(response.Contacts.Contact) < helper.PageSizeLarge {
			break
		}

		page, err := helper.GetNextpageNumber(request.PageNumber)
		if err != nil {
			return errs.WrapError(err)
		}
		request.PageNumber = page
	}
	ids := make([]string, 0)
	names := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"id":                           object.Name,
			"alarm_contact_name":           object.Name,
			"channels_aliim":               object.Channels.AliIM,
			"channels_ding_web_hook":       object.Channels.DingWebHook,
			"channels_mail":                object.Channels.Mail,
			"channels_sms":                 object.Channels.SMS,
			"channels_state_aliim":         object.ChannelsState.AliIM,
			"channels_state_ding_web_hook": object.ChannelsState.DingWebHook,
			"channels_state_mail":          object.ChannelsState.Mail,
			"channels_status_sms":          object.ChannelsState.SMS,
			"contact_groups":               object.ContactGroups.ContactGroup,
			"describe":                     object.Desc,
			"lang":                         object.Lang,
		}
		ids = append(ids, object.Name)
		names = append(names, object.Name)
		s = append(s, mapping)
	}

	d.SetId(helper.DataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return errs.WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return errs.WrapError(err)
	}

	if err := d.Set("contacts", s); err != nil {
		return errs.WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		helper.WriteToFile(output.(string), s)
	}

	return nil
}
