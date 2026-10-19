package ims

import (
	"fmt"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
)

// listVirtualMFADevices lists every virtual MFA device visible to the caller,
// following the IsTruncated/Marker pagination of ListVirtualMFADevices.
//
// The RpcPost body carries the Marker cursor; the API returns devices under
// VirtualMFADevices.VirtualMFADevice.
func listVirtualMFADevices(client *connectivity.AliyunClient) ([]map[string]interface{}, error) {
	const action = "ListVirtualMFADevices"
	var all []map[string]interface{}
	body := map[string]interface{}{}

	for {
		response, err := client.RpcPost("Ims", "2019-08-15", action, nil, body, true)
		if err != nil {
			return nil, err
		}

		raw, _ := response["VirtualMFADevices"].(map[string]interface{})
		devices, _ := raw["VirtualMFADevice"].([]interface{})
		for _, item := range devices {
			if device, ok := item.(map[string]interface{}); ok {
				all = append(all, device)
			}
		}

		if isTruncated, _ := response["IsTruncated"].(bool); !isTruncated {
			break
		}
		marker, _ := response["Marker"].(string)
		if marker == "" || marker == fmt.Sprint(body["Marker"]) {
			break
		}
		body["Marker"] = marker
	}

	return all, nil
}

// mfaDeviceName extracts the device name from a virtual MFA device serial
// number. The serial number has the form acs:ram::<accountId>:mfa/<name>;
// the name is the segment after the last slash.
func mfaDeviceName(serialNumber string) string {
	for i := len(serialNumber) - 1; i >= 0; i-- {
		if serialNumber[i] == '/' {
			return serialNumber[i+1:]
		}
	}
	return serialNumber
}
