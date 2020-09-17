package cs

//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.
//
// Code generated by Alibaba Cloud SDK Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

// Node is a nested struct in cs response
type Node struct {
	InstanceType       string   `json:"instance_type" xml:"instance_type"`
	InstanceRole       string   `json:"instance_role" xml:"instance_role"`
	ExpiredTime        string   `json:"expired_time" xml:"expired_time"`
	State              string   `json:"state" xml:"state"`
	InstanceName       string   `json:"instance_name" xml:"instance_name"`
	IsAliyunNode       bool     `json:"is_aliyun_node" xml:"is_aliyun_node"`
	HostName           string   `json:"host_name" xml:"host_name"`
	ImageId            string   `json:"image_id" xml:"image_id"`
	InstanceStatus     string   `json:"instance_status" xml:"instance_status"`
	InstanceChargeType string   `json:"instance_charge_type" xml:"instance_charge_type"`
	Source             string   `json:"source" xml:"source"`
	ErrorMessage       string   `json:"error_message" xml:"error_message"`
	NodeStatus         string   `json:"node_status" xml:"node_status"`
	CreationTime       string   `json:"creation_time" xml:"creation_time"`
	NodeName           string   `json:"node_name" xml:"node_name"`
	InstanceTypeFamily string   `json:"instance_type_family" xml:"instance_type_family"`
	NodepoolId         string   `json:"nodepool_id" xml:"nodepool_id"`
	InstanceId         string   `json:"instance_id" xml:"instance_id"`
	IpAddress          []string `json:"ip_address" xml:"ip_address"`
}
