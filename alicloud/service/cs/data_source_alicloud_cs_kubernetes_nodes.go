package cs

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"time"

	roaCS "github.com/alibabacloud-go/cs-20151215/v8/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	sdkretry "github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

const kubernetesNodesDataSourceTypeName = "alicloud_cs_kubernetes_nodes"

var (
	_ datasource.DataSource              = &kubernetesNodesDataSource{}
	_ datasource.DataSourceWithConfigure = &kubernetesNodesDataSource{}
)

func NewKubernetesNodesDataSource() datasource.DataSource {
	return &kubernetesNodesDataSource{}
}

type kubernetesNodesDataSource struct {
	fwadapt.DataSourceBase
}

// kubernetesNodesClient is the read-only slice of the CS client this data source
// needs, so the filtering, pagination and ordering logic can be unit tested
// without a cluster.
type kubernetesNodesClient interface {
	DescribeClusterNodePools(clusterId *string, request *roaCS.DescribeClusterNodePoolsRequest) (*roaCS.DescribeClusterNodePoolsResponse, error)
	DescribeClusterNodes(clusterId *string, request *roaCS.DescribeClusterNodesRequest) (*roaCS.DescribeClusterNodesResponse, error)
}

type kubernetesNodeModel struct {
	Id           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	NodepoolId   types.String `tfsdk:"nodepool_id"`
	HostName     types.String `tfsdk:"host_name"`
	InstanceRole types.String `tfsdk:"instance_role"`
	State        types.String `tfsdk:"state"`
}

type kubernetesNodesModel struct {
	Id           types.String          `tfsdk:"id"`
	ClusterId    types.String          `tfsdk:"cluster_id"`
	NodepoolId   types.String          `tfsdk:"nodepool_id"`
	NodepoolName types.String          `tfsdk:"nodepool_name"`
	Ids          []string              `tfsdk:"ids"`
	NodeNames    []string              `tfsdk:"node_names"`
	Nodes        []kubernetesNodeModel `tfsdk:"nodes"`
}

// kubernetesNodeObjectType describes the element shape of the nodes list. A
// plain list attribute with an object element type is used instead of the
// nested attribute variants: this provider serves every half behind one
// protocol 5 server, which cannot express nested attributes.
var kubernetesNodeObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"id":            types.StringType,
	"name":          types.StringType,
	"nodepool_id":   types.StringType,
	"host_name":     types.StringType,
	"instance_role": types.StringType,
	"state":         types.StringType,
}}

func (d *kubernetesNodesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = kubernetesNodesDataSourceTypeName
}

func (d *kubernetesNodesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides the nodes of an ACK cluster, optionally limited to a single node pool.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the data source. It is the cluster ID, followed by the node pool ID separated by a colon when a node pool is selected.",
			},
			"cluster_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the ACK cluster.",
			},
			"nodepool_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The ID of the node pool whose nodes are queried. A configured value takes precedence over `nodepool_name`. When only `nodepool_name` is configured, the resolved node pool ID is exported here.",
			},
			"nodepool_name": schema.StringAttribute{
				Optional:    true,
				Description: "The name of the node pool whose nodes are queried. It is resolved to the node pool ID on every read, so a node pool created in the same configuration can be referenced without creating a dependency cycle. If no node pool of the cluster has this name, the data source returns empty lists.",
			},
			"ids": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "A list of instance IDs of the nodes.",
			},
			"node_names": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "A list of Kubernetes node names, sorted in ascending order.",
			},
			"nodes": schema.ListAttribute{
				Computed:    true,
				ElementType: kubernetesNodeObjectType,
				Description: "A list of nodes. The nodes are sorted by node name.",
			},
		},
	}
}

func (d *kubernetesNodesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var model kubernetesNodesModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	csClient, err := d.Meta().NewRoaCsClient()
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: creating the CS client", kubernetesNodesDataSourceTypeName),
			err.Error(),
		)
		return
	}

	clusterId := model.ClusterId.ValueString()
	// Only the configuration decides which node pool is queried. Reading the
	// computed nodepool_id back from state would keep pointing at the previous
	// pool when the configured name changes, or when a pool is replaced by a new
	// one carrying the same name.
	configuredNodepoolId := ""
	if !model.NodepoolId.IsNull() && !model.NodepoolId.IsUnknown() {
		configuredNodepoolId = model.NodepoolId.ValueString()
	}
	nodepoolName := ""
	if !model.NodepoolName.IsNull() && !model.NodepoolName.IsUnknown() {
		nodepoolName = model.NodepoolName.ValueString()
	}

	result, err := fetchKubernetesNodes(csClient, clusterId, configuredNodepoolId, nodepoolName)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling DescribeClusterNodes", kubernetesNodesDataSourceTypeName),
			err.Error(),
		)
		return
	}

	model.Id = types.StringValue(kubernetesNodesId(clusterId, result.NodepoolId))
	if result.NodepoolId == "" {
		model.NodepoolId = types.StringNull()
	} else {
		model.NodepoolId = types.StringValue(result.NodepoolId)
	}
	model.Ids = make([]string, 0, len(result.Nodes))
	model.NodeNames = make([]string, 0, len(result.Nodes))
	model.Nodes = make([]kubernetesNodeModel, 0, len(result.Nodes))
	for _, node := range result.Nodes {
		instanceId := tea.StringValue(node.InstanceId)
		nodeName := tea.StringValue(node.NodeName)
		model.Ids = append(model.Ids, instanceId)
		model.NodeNames = append(model.NodeNames, nodeName)
		model.Nodes = append(model.Nodes, kubernetesNodeModel{
			Id:           types.StringValue(instanceId),
			Name:         types.StringValue(nodeName),
			NodepoolId:   types.StringValue(tea.StringValue(node.NodepoolId)),
			HostName:     types.StringValue(tea.StringValue(node.HostName)),
			InstanceRole: types.StringValue(tea.StringValue(node.InstanceRole)),
			State:        types.StringValue(tea.StringValue(node.State)),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

type kubernetesNodesResult struct {
	NodepoolId string
	Nodes      []*roaCS.DescribeClusterNodesResponseBodyNodes
}

// fetchKubernetesNodes resolves the node pool to query and returns its nodes in
// a stable node-name order.
func fetchKubernetesNodes(api kubernetesNodesClient, clusterId, configuredNodepoolId, nodepoolName string) (kubernetesNodesResult, error) {
	nodepoolId := configuredNodepoolId
	if nodepoolId == "" && nodepoolName != "" {
		resolved, err := resolveNodepoolIdByName(api, clusterId, nodepoolName)
		if err != nil {
			return kubernetesNodesResult{}, err
		}
		if resolved == "" {
			// A named pool that does not exist yet yields an empty result. Falling
			// back to every node of the cluster would answer a different question
			// than the one asked.
			return kubernetesNodesResult{}, nil
		}
		nodepoolId = resolved
	}

	nodes, err := describeAllClusterNodes(api, clusterId, nodepoolId)
	if err != nil {
		return kubernetesNodesResult{}, err
	}
	// A stable order keeps references such as node_names[0] from flipping between reads.
	sort.Slice(nodes, func(i, j int) bool {
		return tea.StringValue(nodes[i].NodeName) < tea.StringValue(nodes[j].NodeName)
	})
	return kubernetesNodesResult{NodepoolId: nodepoolId, Nodes: nodes}, nil
}

func kubernetesNodesId(clusterId, nodepoolId string) string {
	if nodepoolId == "" {
		return clusterId
	}
	return fmt.Sprintf("%s:%s", clusterId, nodepoolId)
}

func resolveNodepoolIdByName(csClient kubernetesNodesClient, clusterId, name string) (string, error) {
	var response *roaCS.DescribeClusterNodePoolsResponse
	request := &roaCS.DescribeClusterNodePoolsRequest{
		NodepoolName: tea.String(name),
	}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err := sdkretry.Retry(5*time.Minute, func() *sdkretry.RetryError {
		var callErr error
		response, callErr = csClient.DescribeClusterNodePools(tea.String(clusterId), request)
		if callErr != nil {
			if helper.NeedRetry(callErr) {
				wait()
				return sdkretry.RetryableError(callErr)
			}
			return sdkretry.NonRetryableError(callErr)
		}
		return nil
	})
	logs.AddDebug("DescribeClusterNodePools", response, request)
	if err != nil {
		return "", errs.WrapErrorf(err, errs.DataDefaultErrorMsg, kubernetesNodesDataSourceTypeName, "DescribeClusterNodePools", errs.AlibabaCloudSdkGoERROR)
	}
	if response == nil || response.Body == nil {
		return "", nil
	}

	for _, nodepool := range response.Body.Nodepools {
		if nodepool == nil || nodepool.NodepoolInfo == nil {
			continue
		}
		if tea.StringValue(nodepool.NodepoolInfo.Name) == name {
			return tea.StringValue(nodepool.NodepoolInfo.NodepoolId), nil
		}
	}
	return "", nil
}

func describeAllClusterNodes(csClient kubernetesNodesClient, clusterId, nodepoolId string) ([]*roaCS.DescribeClusterNodesResponseBodyNodes, error) {
	var nodes []*roaCS.DescribeClusterNodesResponseBodyNodes
	pageNumber := 1
	pageSize := helper.PageSizeLarge
	for {
		request := &roaCS.DescribeClusterNodesRequest{
			PageNumber: tea.String(strconv.Itoa(pageNumber)),
			PageSize:   tea.String(strconv.Itoa(pageSize)),
		}
		if nodepoolId != "" {
			request.NodepoolId = tea.String(nodepoolId)
		}

		var response *roaCS.DescribeClusterNodesResponse
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err := sdkretry.Retry(5*time.Minute, func() *sdkretry.RetryError {
			var callErr error
			response, callErr = csClient.DescribeClusterNodes(tea.String(clusterId), request)
			if callErr != nil {
				if helper.NeedRetry(callErr) {
					wait()
					return sdkretry.RetryableError(callErr)
				}
				return sdkretry.NonRetryableError(callErr)
			}
			return nil
		})
		logs.AddDebug("DescribeClusterNodes", response, request)
		if err != nil {
			return nil, errs.WrapErrorf(err, errs.DataDefaultErrorMsg, kubernetesNodesDataSourceTypeName, "DescribeClusterNodes", errs.AlibabaCloudSdkGoERROR)
		}
		if response == nil || response.Body == nil || len(response.Body.Nodes) == 0 {
			break
		}

		nodes = append(nodes, response.Body.Nodes...)
		returnedPageSize := pageSize
		if response.Body.Page != nil {
			if response.Body.Page.TotalCount != nil && len(nodes) >= int(tea.Int32Value(response.Body.Page.TotalCount)) {
				break
			}
			// The server may cap the page size below the requested one, so the short-page
			// test uses the size it reports rather than the size that was asked for.
			if response.Body.Page.PageSize != nil && int(tea.Int32Value(response.Body.Page.PageSize)) > 0 {
				returnedPageSize = int(tea.Int32Value(response.Body.Page.PageSize))
			}
		}
		if len(response.Body.Nodes) < returnedPageSize {
			break
		}
		pageNumber += 1
	}
	return nodes, nil
}
