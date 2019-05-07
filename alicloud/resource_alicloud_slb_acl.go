package alicloud

import (
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudSlbAcl() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudSlbAclCreate,
		Read:   resourceAlicloudSlbAclRead,
		Update: resourceAlicloudSlbAclUpdate,
		Delete: resourceAlicloudSlbAclDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"ip_version": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      IPVersion4,
				ValidateFunc: validateAllowedStringValue([]string{string(IPVersion4), string(IPVersion6)}),
			},
			"entry_list": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"entry": {
							Type:     schema.TypeString,
							Required: true,
						},
						"comment": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				MaxItems: 300,
				MinItems: 0,
			},
		},
	}
}

func resourceAlicloudSlbAclCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := slb.CreateCreateAccessControlListRequest()
	request.AclName = strings.TrimSpace(d.Get("name").(string))
	request.AddressIPVersion = d.Get("ip_version").(string)

	raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.CreateAccessControlList(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_slb_acl", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName, raw)
	response, _ := raw.(*slb.CreateAccessControlListResponse)

	d.SetId(response.AclId)
	return resourceAlicloudSlbAclUpdate(d, meta)
}

func resourceAlicloudSlbAclRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slbService := SlbService{client}

	object, err := slbService.DescribeSlbAcl(d.Id())
	if err != nil {
		if IsExceptedError(err, SlbAclNotExists) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("name", object.AclName)
	d.Set("ip_version", object.AddressIPVersion)
	if err := d.Set("entry_list", slbService.FlattenSlbAclEntryMappings(object.AclEntrys.AclEntry)); err != nil {
		return WrapError(err)
	}
	return nil
}

func resourceAlicloudSlbAclUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slbService := SlbService{client}

	d.Partial(true)

	if !d.IsNewResource() && d.HasChange("name") {
		request := slb.CreateSetAccessControlListAttributeRequest()
		request.AclId = d.Id()
		request.AclName = d.Get("name").(string)
		raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.SetAccessControlListAttribute(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)

		}
		addDebug(request.GetActionName(), raw)
		d.SetPartial("name")
	}

	if d.HasChange("entry_list") {
		o, n := d.GetChange("entry_list")
		oe := o.(*schema.Set)
		ne := n.(*schema.Set)
		remove := oe.Difference(ne).List()
		add := ne.Difference(oe).List()

		if len(remove) > 0 {
			if err := slbService.SlbRemoveAccessControlListEntry(remove, d.Id()); err != nil {
				return WrapError(err)
			}
		}

		if len(add) > 0 {
			if err := slbService.SlbAddAccessControlListEntry(add, d.Id()); err != nil {
				return WrapError(err)
			}
		}
		d.SetPartial("entry_list")
	}

	d.Partial(false)

	return resourceAlicloudSlbAclRead(d, meta)
}

func resourceAlicloudSlbAclDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slbService := SlbService{client}
	request := slb.CreateDeleteAccessControlListRequest()
	request.AclId = d.Id()
	raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.DeleteAccessControlList(request)
	})
	if err != nil {
		if !IsExceptedError(err, SlbAclNotExists) {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}
	addDebug(request.GetActionName(), raw)
	return WrapError(slbService.WaitForSlbAcl(d.Id(), Deleted, DefaultTimeoutMedium))
}
