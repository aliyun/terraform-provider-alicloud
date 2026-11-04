package cs

import (
	"context"
	"reflect"
	"testing"

	roacs "github.com/alibabacloud-go/cs-20151215/v8/client"
	"github.com/alibabacloud-go/tea/tea"
)

func TestUnitCSKubernetesStateUpgradeV0(t *testing.T) {
	fields := []string{"runtime", "certificate_authority", "connections"}
	cases := []struct {
		name     string
		input    map[string]interface{}
		expected map[string][]interface{}
	}{
		{
			name: "map with data",
			input: map[string]interface{}{
				"id": "c-k8s-123",
				"runtime": map[string]interface{}{
					"name":    "containerd",
					"version": "1.6.20",
				},
				"certificate_authority": map[string]interface{}{
					"cluster_cert": "cert",
					"client_cert":  "client",
					"client_key":   "key",
				},
				"connections": map[string]interface{}{
					"api_server_internet": "https://2.3.4.5:6443",
					"master_public_ip":    "2.3.4.5",
				},
			},
			expected: map[string][]interface{}{
				"runtime": {
					map[string]interface{}{
						"name":    "containerd",
						"version": "1.6.20",
					},
				},
				"certificate_authority": {
					map[string]interface{}{
						"cluster_cert": "cert",
						"client_cert":  "client",
						"client_key":   "key",
					},
				},
				"connections": {
					map[string]interface{}{
						"api_server_internet": "https://2.3.4.5:6443",
						"master_public_ip":    "2.3.4.5",
					},
				},
			},
		},
		{
			name: "empty map",
			input: map[string]interface{}{
				"id":                    "c-k8s-456",
				"runtime":               map[string]interface{}{},
				"certificate_authority": map[string]interface{}{},
				"connections":           map[string]interface{}{},
			},
			expected: map[string][]interface{}{
				"runtime":               {},
				"certificate_authority": {},
				"connections":           {},
			},
		},
		{
			name: "nil value",
			input: map[string]interface{}{
				"id":                    "c-k8s-789",
				"runtime":               nil,
				"certificate_authority": nil,
				"connections":           nil,
			},
			expected: nil,
		},
		{
			name: "field not present",
			input: map[string]interface{}{
				"id": "c-k8s-000",
			},
			expected: nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := resourceCSKubernetesStateUpgradeV0(context.Background(), tc.input, nil)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			for _, field := range fields {
				got := result[field]
				if tc.expected == nil {
					if got != nil {
						t.Errorf("%s: expected nil, got %v", field, got)
					}
					continue
				}
				exp := tc.expected[field]
				gotList, ok := got.([]interface{})
				if !ok {
					t.Fatalf("%s: expected []interface{}, got %T", field, got)
				}
				if len(gotList) != len(exp) {
					t.Fatalf("%s: expected length %d, got %d", field, len(exp), len(gotList))
				}
				for i, item := range gotList {
					gotMap := item.(map[string]interface{})
					expMap := exp[i].(map[string]interface{})
					for k, v := range expMap {
						if gotMap[k] != v {
							t.Errorf("%s[%d].%s: expected %v, got %v", field, i, k, v, gotMap[k])
						}
					}
				}
			}
		})
	}
}

func TestUnitCSKubernetesSchemaVersionV0ToV1(t *testing.T) {
	r := ResourceCSKubernetes()
	if r.SchemaVersion != 1 {
		t.Errorf("expected SchemaVersion 1, got %d", r.SchemaVersion)
	}
	if len(r.StateUpgraders) != 1 {
		t.Fatalf("expected 1 StateUpgrader, got %d", len(r.StateUpgraders))
	}
	if r.StateUpgraders[0].Version != 0 {
		t.Errorf("expected StateUpgrader version 0, got %d", r.StateUpgraders[0].Version)
	}
}

func TestUnit_parseRRSAMetadata(t *testing.T) {
	type args struct {
		meta string
	}
	tests := []struct {
		name    string
		args    args
		want    []map[string]interface{}
		wantErr bool
	}{
		{
			name: "empty error",
			args: args{
				meta: "",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid error",
			args: args{
				meta: `
{
	"RRSAConfig": 1234
}`,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "valid enabled",
			args: args{
				meta: `{
	"RRSAConfig": {
		"enabled": true,
		"issuer": "https://example.com,https://kubernetes.default.svc,kubernetes.default.svc",
		"oidc_name": "ack-rrsa-c12345",
		"oidc_arn": "acs:ram::12345:oidc-provider/ack-rrsa-c12345"
	}
}`,
			},
			want: []map[string]interface{}{{
				"enabled":                true,
				"rrsa_oidc_issuer_url":   "https://example.com",
				"ram_oidc_provider_name": "ack-rrsa-c12345",
				"ram_oidc_provider_arn":  "acs:ram::12345:oidc-provider/ack-rrsa-c12345",
			}},
			wantErr: false,
		},
		{
			name: "valid not enabled",
			args: args{
				meta: `{
	"RRSAConfig": {
		"enabled": false,
		"issuer": "",
		"oidc_name": "",
		"oidc_arn": ""
	}
}`,
			},
			want: []map[string]interface{}{{
				"enabled":                false,
				"rrsa_oidc_issuer_url":   "",
				"ram_oidc_provider_name": "",
				"ram_oidc_provider_arn":  "",
			}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FlattenRRSAMetadata(tt.args.meta)
			if tt.wantErr && err == nil {
				t.Errorf("flattenRRSAMetadata(%v) want error got nil", tt.args.meta)
			}
			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("flattenRRSAMetadata(%v) want %v got %v", tt.args.meta, tt.want, got)
			}
		})
	}
}

func TestUnitFlattenCsMasterNodes(t *testing.T) {
	masterWithIp := &roacs.DescribeClusterNodesResponseBodyNodes{
		InstanceId:   tea.String("i-master-1"),
		NodeName:     tea.String("master-1"),
		InstanceRole: tea.String("Master"),
		IpAddress:    []*string{tea.String("192.168.0.1"), tea.String("192.168.0.2")},
	}
	masterEmptyIp := &roacs.DescribeClusterNodesResponseBodyNodes{
		InstanceId:   tea.String("i-master-2"),
		NodeName:     tea.String("master-2"),
		InstanceRole: tea.String("Master"),
		IpAddress:    []*string{},
	}
	masterNilIp := &roacs.DescribeClusterNodesResponseBodyNodes{
		InstanceId:   tea.String("i-master-3"),
		NodeName:     tea.String("master-3"),
		InstanceRole: tea.String("Master"),
	}
	masterNilRole := &roacs.DescribeClusterNodesResponseBodyNodes{
		InstanceId: tea.String("i-master-4"),
		NodeName:   tea.String("master-4"),
		IpAddress:  []*string{tea.String("192.168.0.4")},
	}
	worker := &roacs.DescribeClusterNodesResponseBodyNodes{
		InstanceId:   tea.String("i-worker-1"),
		NodeName:     tea.String("worker-1"),
		InstanceRole: tea.String("Worker"),
		IpAddress:    []*string{tea.String("192.168.1.1")},
	}
	workerEmptyIp := &roacs.DescribeClusterNodesResponseBodyNodes{
		InstanceId:   tea.String("i-worker-2"),
		NodeName:     tea.String("worker-2"),
		InstanceRole: tea.String("Worker"),
		IpAddress:    []*string{},
	}

	tests := []struct {
		name  string
		nodes []*roacs.DescribeClusterNodesResponseBodyNodes
		want  []map[string]interface{}
	}{
		{
			name:  "nil node list",
			nodes: nil,
			want:  nil,
		},
		{
			name:  "empty node list",
			nodes: []*roacs.DescribeClusterNodesResponseBodyNodes{},
			want:  nil,
		},
		{
			name:  "nil node element",
			nodes: []*roacs.DescribeClusterNodesResponseBodyNodes{nil},
			want:  nil,
		},
		{
			name:  "master with empty ip address list",
			nodes: []*roacs.DescribeClusterNodesResponseBodyNodes{masterEmptyIp},
			want: []map[string]interface{}{
				{"id": "i-master-2", "name": "master-2", "private_ip": ""},
			},
		},
		{
			name:  "master with nil ip address list",
			nodes: []*roacs.DescribeClusterNodesResponseBodyNodes{masterNilIp},
			want: []map[string]interface{}{
				{"id": "i-master-3", "name": "master-3", "private_ip": ""},
			},
		},
		{
			name:  "master with nil instance role is skipped",
			nodes: []*roacs.DescribeClusterNodesResponseBodyNodes{masterNilRole},
			want:  nil,
		},
		{
			name:  "master with addresses keeps the first one",
			nodes: []*roacs.DescribeClusterNodesResponseBodyNodes{masterWithIp},
			want: []map[string]interface{}{
				{"id": "i-master-1", "name": "master-1", "private_ip": "192.168.0.1"},
			},
		},
		{
			name:  "worker nodes are not master nodes",
			nodes: []*roacs.DescribeClusterNodesResponseBodyNodes{worker, workerEmptyIp},
			want:  nil,
		},
		{
			name:  "mixed list keeps masters only",
			nodes: []*roacs.DescribeClusterNodesResponseBodyNodes{worker, masterWithIp, nil, masterEmptyIp, masterNilRole},
			want: []map[string]interface{}{
				{"id": "i-master-1", "name": "master-1", "private_ip": "192.168.0.1"},
				{"id": "i-master-2", "name": "master-2", "private_ip": ""},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FlattenCsMasterNodes(tt.nodes)
			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("flattenCsMasterNodes() want %v got %v", tt.want, got)
			}
		})
	}
}

const testKubeConfigYaml = `apiVersion: v1
kind: Config
clusters:
- cluster:
    certificate-authority-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCg==
    server: https://1.2.3.4:6443
  name: kubernetes
contexts:
- context:
    cluster: kubernetes
    user: kubernetes-admin
  name: kubernetes-admin@kubernetes
current-context: kubernetes-admin@kubernetes
users:
- name: kubernetes-admin
  user:
    client-certificate-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCg==
    client-key-data: LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQo=
`

func TestUnit_flattenAlicloudCSCertificate(t *testing.T) {
	normal := map[string]string{
		"cluster_cert": "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCg==",
		"client_cert":  "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCg==",
		"client_key":   "LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQo=",
	}
	body := func(config string) *roacs.DescribeClusterUserKubeconfigResponseBody {
		return &roacs.DescribeClusterUserKubeconfigResponseBody{Config: tea.String(config)}
	}
	tests := []struct {
		name        string
		certificate *roacs.DescribeClusterUserKubeconfigResponseBody
		want        map[string]string
		wantErr     bool
	}{
		{
			name:        "nil response body",
			certificate: nil,
			wantErr:     true,
		},
		{
			name:        "empty config",
			certificate: body(""),
			wantErr:     true,
		},
		{
			name:        "invalid yaml",
			certificate: body(": ["),
			wantErr:     true,
		},
		{
			name: "missing clusters",
			certificate: body(`users:
- name: kubernetes-admin
  user:
    client-certificate-data: YQ==
    client-key-data: Yg==
`),
			wantErr: true,
		},
		{
			name: "missing users",
			certificate: body(`clusters:
- cluster:
    certificate-authority-data: Yw==
  name: kubernetes
`),
			wantErr: true,
		},
		{
			name: "empty clusters array",
			certificate: body(`clusters: []
users:
- name: kubernetes-admin
  user:
    client-certificate-data: YQ==
    client-key-data: Yg==
`),
			wantErr: true,
		},
		{
			name: "clusters not a list",
			certificate: body(`clusters: kubernetes
users:
- name: kubernetes-admin
  user:
    client-certificate-data: YQ==
    client-key-data: Yg==
`),
			wantErr: true,
		},
		{
			name: "leaf field has wrong type",
			certificate: body(`clusters:
- cluster:
    certificate-authority-data: 123456
  name: kubernetes
users:
- name: kubernetes-admin
  user:
    client-certificate-data: YQ==
    client-key-data: Yg==
`),
			wantErr: true,
		},
		{
			name: "missing leaf field",
			certificate: body(`clusters:
- cluster:
    server: https://1.2.3.4:6443
  name: kubernetes
users:
- name: kubernetes-admin
  user:
    client-certificate-data: YQ==
    client-key-data: Yg==
`),
			wantErr: true,
		},
		{
			name: "empty leaf field",
			certificate: body(`clusters:
- cluster:
    certificate-authority-data: ""
  name: kubernetes
users:
- name: kubernetes-admin
  user:
    client-certificate-data: YQ==
    client-key-data: Yg==
`),
			wantErr: true,
		},
		{
			name:        "normal kubeconfig",
			certificate: body(testKubeConfigYaml),
			want:        normal,
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FlattenAlicloudCSCertificate(tt.certificate)
			if tt.wantErr {
				if err == nil {
					t.Errorf("flattenAlicloudCSCertificate(%v) want error got nil", tt.certificate)
				}
				return
			}
			if err != nil {
				t.Errorf("flattenAlicloudCSCertificate(%v) unexpected error: %v", tt.certificate, err)
				return
			}
			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("flattenAlicloudCSCertificate(%v) want %v got %v", tt.certificate, tt.want, got)
			}
		})
	}
}
