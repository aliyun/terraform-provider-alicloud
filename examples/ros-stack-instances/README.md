# ROS Stack Instances Example

This example demonstrates how to use the `alicloud_ros_stack_instances` resource to deploy infrastructure across multiple regions and accounts using ROS (Resource Orchestration Service) Stack Groups.

## Overview

This example creates:
1. A ROS Stack Group with a template that provisions VPC and VSwitch resources
2. Stack Instances deployed across multiple target regions and accounts
3. Demonstrates both self-managed and service-managed permission models

## Architecture

The example deploys a simple VPC and VSwitch infrastructure to multiple regions simultaneously, showcasing the power of ROS Stack Instances for multi-region deployments.

## Prerequisites

Before running this example, ensure you have:

1. **Alibaba Cloud Account**: An active Alibaba Cloud account with appropriate permissions
2. **Terraform**: Terraform >= 0.12 installed
3. **ROS Permissions**: RAM permissions for ROS operations including:
   - `ros:CreateStackGroup`
   - `ros:CreateStackInstances`
   - `ros:ListStackInstances`
   - `ros:UpdateStackInstances`
   - `ros:DeleteStackInstances`
4. **Multi-Account Setup** (optional): For multi-account deployments, ensure proper RAM role configuration

## Usage

### Step 1: Initialize Terraform

```bash
terraform init
```

### Step 2: Review Configuration

Review the `variables.tf` file and adjust the parameters according to your requirements:

- `target_regions`: List of regions where you want to deploy (default: cn-beijing, cn-shanghai)
- `target_accounts`: List of account IDs for deployment (update with your actual account IDs)
- `vpc_cidr_block`: CIDR block for VPCs
- `max_concurrent_count`: Number of concurrent deployments

### Step 3: Plan the Deployment

```bash
terraform plan
```

You will be prompted to enter variable values, or they will use the defaults from `variables.tf`.

### Step 4: Apply the Configuration

```bash
terraform apply
```

Confirm the action when prompted. The deployment will:
1. Create a ROS Stack Group
2. Deploy stack instances to all specified regions and accounts
3. Wait for all operations to complete (up to 60 minutes by default)

### Step 5: View Outputs

After successful deployment, Terraform will display outputs including:
- Stack group information
- List of created stack instances
- Deployment status for each region/account combination

```bash
terraform output
```

### Step 6: Verify Resources

You can verify the deployed resources in the Alibaba Cloud Console:
1. Navigate to ROS Console
2. Go to "Stack Groups" section
3. Find your stack group by name
4. View stack instances and their statuses

## Configuration Examples

### Self-Managed Permissions (Default)

The example uses self-managed permissions by default, where you explicitly specify account IDs:

```hcl
resource "alicloud_ros_stack_instances" "self_managed" {
  stack_group_name = "my-stack-group"
  region_ids       = ["cn-beijing", "cn-shanghai"]
  account_ids      = ["123456789012****", "098765432109****"]
  
  parameter_overrides {
     parameter_key   = "VpcCidrBlock"
     parameter_value = "172.16.0.0/12"
  }
}
```

### Service-Managed Permissions

To use service-managed permissions with Resource Directory folders, uncomment the `service_managed` block in `main.tf` and comment out the `self_managed` block:

```hcl
resource "alicloud_ros_stack_instances" "service_managed" {
  stack_group_name = "my-stack-group"
  region_ids       = ["cn-beijing", "cn-shanghai"]
  
  deployment_targets {
    accounts      = ["111111111111****"]
    rd_folder_ids = ["fd-abc123****"]
  }
  
  deployment_options = ["IgnoreExisting"]
}
```

## Important Notes

### Asynchronous Operations
- Stack instance operations are asynchronous batch operations
- The provider waits for operation completion but partial failures may occur
- If some instances fail, successful instances are NOT automatically rolled back

### Operation Conflicts
- Only one stack group operation can run at a time
- Concurrent operations will result in `StackGroupOperationInProgress` errors
- The provider automatically retries in such cases

### State Management
- If create operation has partial failures, the resource ID is NOT set in state
- This allows Terraform to retry the entire operation on next apply
- This prevents orphaned resources in your environment

### ForceNew Parameters
The following parameters require resource recreation if modified:
- `stack_group_name`
- `region_ids`
- `account_ids` / `deployment_targets`
- `disable_rollback`
- `deployment_options`

Only these parameters support in-place updates:
- `parameter_overrides`
- `operation_preferences`
- `timeout_in_minutes`
- `operation_description`

## Cleanup

To destroy all resources created by this example:

```bash
terraform destroy
```

This will:
1. Delete all stack instances from target regions and accounts
2. Delete the stack group
3. Remove all associated resources

**Warning**: This action cannot be undone. All deployed infrastructure will be permanently deleted.

## Troubleshooting

### Common Issues

1. **StackGroupOperationInProgress Error**
   - Wait for the current operation to complete
   - Check the ROS console for ongoing operations
   - The provider will automatically retry

2. **Partial Failures**
   - Check the `stack_instances` output for individual instance statuses
   - Review failed instances in the ROS console
   - Manually retry or fix issues and reapply

3. **Permission Errors**
   - Ensure RAM roles are properly configured
   - Verify account has necessary ROS permissions
   - Check cross-account trust relationships for multi-account setups

4. **Timeout Issues**
   - Increase `timeout_in_minutes` if deployments take longer
   - Reduce `max_concurrent_count` to avoid resource contention
   - Check network connectivity to all target regions

### Viewing Operation Logs

Check operation details in the Alibaba Cloud ROS Console:
1. Navigate to Stack Groups
2. Select your stack group
3. View "Operations" tab for detailed logs
4. Check individual stack instance statuses

## Additional Resources

- [ROS Stack Instances Documentation](https://www.alibabacloud.com/help/en/resource-orchestration-service/latest/api-ros-2019-09-10-createstackinstances)
- [ROS Stack Groups Best Practices](https://www.alibabacloud.com/help/en/resource-orchestration-service/latest/stack-groups)
- [Terraform Alicloud Provider Documentation](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs)

## Support

For issues related to:
- **Terraform Provider**: Open an issue on [GitHub](https://github.com/aliyun/terraform-provider-alicloud/issues)
- **ROS Service**: Contact [Alibaba Cloud Support](https://www.alibabacloud.com/support)
- **This Example**: Open an issue with the example tag
