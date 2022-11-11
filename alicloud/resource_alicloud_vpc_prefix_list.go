package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudVpcPrefixList() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudVpcPrefixListCreate,
		Read:   resourceAlicloudVpcPrefixListRead,
		Update: resourceAlicloudVpcPrefixListUpdate,
		Delete: resourceAlicloudVpcPrefixListDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"entrys": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cidr": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"description": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.All(validation.StringLenBetween(2, 256), validation.StringDoesNotMatch(regexp.MustCompile(`(^http://.*)|(^https://.*)`), "It cannot begin with \"http://\", \"https://\".")),
						},
					},
				},
			},
			"ip_version": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"IPV4", "IPV6"}, false),
			},
			"max_entries": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"prefix_list_description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.All(validation.StringLenBetween(2, 256), validation.StringDoesNotMatch(regexp.MustCompile(`(^http://.*)|(^https://.*)`), "It cannot begin with \"http://\", \"https://\".")),
			},
			"prefix_list_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_.-]{1,127}$`), "The name must be 2 to 128 characters in length, and must start with a letter. It can contain digits, periods (.), underscores (_), and hyphens (-)."),
			},
		},
	}
}

func resourceAlicloudVpcPrefixListCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateVpcPrefixList"
	request := make(map[string]interface{})
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("entrys"); ok {
		for entrysPtr, entrys := range v.(*schema.Set).List() {
			entrysArg := entrys.(map[string]interface{})
			request["PrefixListEntrys."+fmt.Sprint(entrysPtr+1)+".Cidr"] = entrysArg["cidr"]
			request["PrefixListEntrys."+fmt.Sprint(entrysPtr+1)+".Description"] = entrysArg["description"]
		}
	}
	if v, ok := d.GetOk("ip_version"); ok {
		request["IpVersion"] = v
	}
	if v, ok := d.GetOk("max_entries"); ok {
		request["MaxEntries"] = v
	}
	if v, ok := d.GetOk("prefix_list_description"); ok {
		request["PrefixListDescription"] = v
	}
	if v, ok := d.GetOk("prefix_list_name"); ok {
		request["PrefixListName"] = v
	}
	request["RegionId"] = client.RegionId
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		request["ClientToken"] = buildClientToken("CreateVpcPrefixList")
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"OperationConflict", "SystemBusy", "IncorrectStatus.%s"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_vpc_prefix_list", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["PrefixListId"]))

	return resourceAlicloudVpcPrefixListRead(d, meta)
}

func resourceAlicloudVpcPrefixListRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	object, err := vpcService.DescribeVpcPrefixList(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_vpc_prefix_list vpcService.DescribeVpcPrefixList Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("ip_version", object["IpVersion"])
	if v, ok := object["MaxEntries"]; ok && fmt.Sprint(v) != "0" {
		d.Set("max_entries", formatInt(v))
	}
	d.Set("prefix_list_name", object["PrefixListName"])
	d.Set("prefix_list_description", object["PrefixListDescription"])
	getVpcPrefixListEntriesObject, err := vpcService.GetVpcPrefixListEntries(d.Id())
	if err != nil {
		return WrapError(err)
	}
	entrysMaps := make([]map[string]interface{}, 0)
	for _, prefixListEntryListItem := range getVpcPrefixListEntriesObject {
		entrysMaps = append(entrysMaps, map[string]interface{}{
			"cidr":        prefixListEntryListItem["Cidr"],
			"description": prefixListEntryListItem["Description"],
		})
	}
	d.Set("entrys", entrysMaps)
	return nil
}

func resourceAlicloudVpcPrefixListUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	update := false
	request := map[string]interface{}{}
	request["PrefixListId"] = d.Id()
	request["RegionId"] = client.RegionId
	if d.HasChange("entrys") {
		update = true
		oldEntry, newEntry := d.GetChange("entrys")
		oldEntrySet := oldEntry.(*schema.Set)
		newEntrySet := newEntry.(*schema.Set)
		removed := oldEntrySet.Difference(newEntrySet)
		added := newEntrySet.Difference(oldEntrySet)

		for entrysPtr, entrys := range removed.List() {
			entrysArg := entrys.(map[string]interface{})
			request["RemovePrefixListEntry."+fmt.Sprint(entrysPtr+1)+".Cidr"] = entrysArg["cidr"]
			request["RemovePrefixListEntry."+fmt.Sprint(entrysPtr+1)+".Description"] = entrysArg["description"]
		}

		for entrysPtr, entrys := range added.List() {
			entrysArg := entrys.(map[string]interface{})
			request["AddPrefixListEntry."+fmt.Sprint(entrysPtr+1)+".Cidr"] = entrysArg["cidr"]
			request["AddPrefixListEntry."+fmt.Sprint(entrysPtr+1)+".Description"] = entrysArg["description"]
		}
	}
	if d.HasChange("max_entries") {
		update = true
		if v, ok := d.GetOk("max_entries"); ok {
			request["MaxEntries"] = v
		}
	}
	if d.HasChange("prefix_list_description") {
		update = true
		if v, ok := d.GetOk("prefix_list_description"); ok {
			request["PrefixListDescription"] = v
		}
	}
	if d.HasChange("prefix_list_name") {
		update = true
		if v, ok := d.GetOk("prefix_list_name"); ok {
			request["PrefixListName"] = v
		}
	}
	if update {
		action := "ModifyVpcPrefixList"
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			request["ClientToken"] = buildClientToken("ModifyVpcPrefixList")
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if IsExpectedErrors(err, []string{"IncorrectStatus.PrefixList", "SystemBusy", "LastTokenProcessing", "IncorrectStatus.%s"}) || NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}
	return resourceAlicloudVpcPrefixListRead(d, meta)
}

func resourceAlicloudVpcPrefixListDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteVpcPrefixList"
	var response map[string]interface{}
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"PrefixListId": d.Id(),
	}

	request["RegionId"] = client.RegionId
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		request["ClientToken"] = buildClientToken("DeleteVpcPrefixList")
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"OperationConflict", "SystemBusy"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
