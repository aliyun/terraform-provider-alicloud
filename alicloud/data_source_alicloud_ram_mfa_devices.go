package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudRamMfaDevices() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudRamMfaDevicesRead,
		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"mfa_devices": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"serial_number": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"virtual_mfa_device_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"activate_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudRamMfaDevicesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var nameRegexOk bool
	var nameRegexStr string
	if v, ok := d.GetOk("name_regex"); ok {
		nameRegexOk = true
		nameRegexStr = v.(string)
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

	action := "ListVirtualMFADevices"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	request = make(map[string]interface{})

	var allDevices []interface{}

	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		var err error
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutRead)), func() *resource.RetryError {
			response, err = client.RpcPost("Ims", "2019-08-15", action, query, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ram_mfa_devices", action, AlibabaCloudSdkGoERROR)
		}

		v, err := jsonpath.Get("$.VirtualMFADevices.VirtualMFADevice[*]", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, "alicloud_ram_mfa_devices", "$.VirtualMFADevices.VirtualMFADevice[*]", response)
		}

		result, _ := v.([]interface{})
		allDevices = append(allDevices, result...)

		if isTruncated, _ := response["IsTruncated"].(bool); !isTruncated {
			break
		}

		marker, _ := response["Marker"].(string)
		if marker == "" || marker == fmt.Sprint(request["Marker"]) {
			break
		}
		request["Marker"] = marker
	}

	var ids []string
	var mfaDevices []map[string]interface{}

	for _, item := range allDevices {
		device, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		serialNumber := fmt.Sprint(device["SerialNumber"])

		// Extract device name from the SerialNumber:
		// acs:ram::<accountId>:mfa/<device-name>
		var deviceName string
		parts := strings.Split(serialNumber, "/")
		if len(parts) > 0 {
			deviceName = parts[len(parts)-1]
		}

		if nameRegexOk {
			r, err := regexp.Compile(nameRegexStr)
			if err != nil {
				return WrapError(err)
			}
			if !r.MatchString(deviceName) {
				continue
			}
		}

		if len(idsMap) > 0 {
			if _, ok := idsMap[serialNumber]; !ok {
				continue
			}
		}

		mapping := map[string]interface{}{
			"serial_number":           serialNumber,
			"virtual_mfa_device_name": deviceName,
			"activate_date":           fmt.Sprint(device["ActivateDate"]),
		}
		log.Printf("[DEBUG] alicloud_ram_mfa_devices - adding device: %v", mapping)
		ids = append(ids, serialNumber)
		mfaDevices = append(mfaDevices, mapping)
	}

	d.SetId(dataResourceIdHash(ids))

	if err := d.Set("mfa_devices", mfaDevices); err != nil {
		return WrapError(err)
	}

	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), mfaDevices)
	}

	return nil
}
