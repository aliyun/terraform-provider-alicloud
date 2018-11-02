package alicloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
	"time"
)

func queryPrivateIps(d *schema.ResourceData, meta interface{}) ([]string, error) {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}

	if eni, err := ecsService.DescribeNetworkInterfaceById("", d.Id()); err != nil {
		return nil, fmt.Errorf("Describe NetworkInterface(%s) failed, %s", d.Id(), err)
	} else {
		filterIps := make([]string, 0, len(eni.PrivateIpSets.PrivateIpSet))
		for i := range eni.PrivateIpSets.PrivateIpSet {
			if eni.PrivateIpSets.PrivateIpSet[i].Primary {
				continue
			}
			filterIps = append(filterIps, eni.PrivateIpSets.PrivateIpSet[i].PrivateIpAddress)
		}
		return filterIps, nil
	}
}

func waitForPrivateIpsCountChanged(d *schema.ResourceData, meta interface{}) error {
	deadline := time.Now().Add(DefaultTimeout * time.Second)
	for {
		if time.Now().After(deadline) {
			return fmt.Errorf("Wait for private IP addrsses count changed timeout")
		}
		time.Sleep(DefaultIntervalShort * time.Second)

		ips, err := queryPrivateIps(d, meta)
		if err != nil {
			return fmt.Errorf("Query private IP failed, %s", err)
		}
		if len(ips) == d.Get("private_ips_count").(int) {
			return nil
		}
	}
}

func waitForPrivateIpsListChanged(d *schema.ResourceData, meta interface{}) error {
	deadline := time.Now().Add(DefaultTimeout * time.Second)
	for {
		if time.Now().After(deadline) {
			return fmt.Errorf("Wait for private IP addrsses list changed timeout")
		}
		time.Sleep(DefaultIntervalShort * time.Second)

		ips, err := queryPrivateIps(d, meta)
		if err != nil {
			return fmt.Errorf("Query private IP failed, %s", err)
		}

		ns := d.Get("private_ips").(*schema.Set)
		if len(ips) != ns.Len() {
			continue
		}

		diff := false
		for i := range ips {
			if !ns.Contains(ips[i]) {
				diff = true
				break
			}
		}

		if !diff {
			return nil
		}
	}
}
