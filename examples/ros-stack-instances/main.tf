# ROS Stack Group for managing stack instances
resource "alicloud_ros_stack_group" "example" {
  stack_group_name = var.stack_group_name
  description      = "Stack group for multi-region deployment example"
  
  template_body = jsonencode({
    ROSTemplateFormatVersion = "2015-09-01"
    Description              = "Example VPC and VSwitch"
    Parameters = {
      VpcCidrBlock = {
        Type        = "String"
        Default     = "172.16.0.0/12"
        Description = "CIDR block for VPC"
      }
      VSwitchCidrBlock = {
        Type        = "String"
        Default     = "172.16.0.0/16"
        Description = "CIDR block for VSwitch"
      }
    }
    Resources = {
      Vpc = {
        Type       = "ALIYUN::ECS::VPC"
        Properties = {
          CidrBlock = { Ref = "VpcCidrBlock" }
          VpcName   = { "Fn::Join": ["-", ["example-vpc", { Ref = "ALIYUN::Region" }]] }
        }
      }
      VSwitch = {
        Type       = "ALIYUN::ECS::VSwitch"
        Properties = {
          VpcId     = { Ref = "Vpc" }
          CidrBlock = { Ref = "VSwitchCidrBlock" }
          ZoneId    = { "Fn::Select": ["0", { "Fn::GetAZs": { Ref = "ALIYUN::Region" } }] }
          VSwitchName = { "Fn::Join": ["-", ["example-vswitch", { Ref = "ALIYUN::Region" }]] }
        }
        DependsOn = ["Vpc"]
      }
    }
    Outputs = {
      VpcId = {
        Value = { Ref = "Vpc" }
        Description = "The ID of the created VPC"
      }
      VSwitchId = {
        Value = { Ref = "VSwitch" }
        Description = "The ID of the created VSwitch"
      }
    }
  })

  tags = {
    Environment = var.environment
    ManagedBy   = "Terraform"
    Project     = var.project_name
  }
}

# ROS Stack Instances - Self-Managed Permissions Model
resource "alicloud_ros_stack_instances" "self_managed" {
  stack_group_name = alicloud_ros_stack_group.example.stack_group_name
  
  # Target regions for deployment (1-20 regions)
  region_ids = var.target_regions
  
  # Target accounts for self-managed permissions (1-50 accounts)
  account_ids = var.target_accounts

  # Override template parameters for all stack instances
  parameter_overrides {
    parameter_key   = "VpcCidrBlock"
    parameter_value = var.vpc_cidr_block
  }

  parameter_overrides {
    parameter_key   = "VSwitchCidrBlock"
    parameter_value = var.vswitch_cidr_block
  }

  # Operation preferences for batch deployment
  operation_preferences {
    max_concurrent_count    = var.max_concurrent_count
    failure_tolerance_count = var.failure_tolerance_count
    region_concurrency_type = "PARALLEL"
  }

  # Timeout configuration
  timeout_in_minutes = var.operation_timeout

  # Operation description for tracking
  operation_description = "Deploy infrastructure to ${length(var.target_regions)} regions across ${length(var.target_accounts)} accounts"
  
  # Disable rollback on failure (optional)
  disable_rollback = false
}

# Alternative: Service-Managed Permissions Model (commented out)
# Uncomment this block and comment out the self_managed block above to use service-managed permissions
/*
resource "alicloud_ros_stack_instances" "service_managed" {
  stack_group_name = alicloud_ros_stack_group.example.stack_group_name
  
  # Target regions for deployment
  region_ids = var.target_regions
  
  # Deployment targets for service-managed permissions
  deployment_targets {
    rd_folder_ids = var.rd_folder_ids
  }

  # Deployment options
  deployment_options = ["IgnoreExisting"]

  # Override template parameters
  parameter_overrides {
    parameter_key   = "VpcCidrBlock"
    parameter_value = var.vpc_cidr_block
  }

  parameter_overrides {
    parameter_key   = "VSwitchCidrBlock"
    parameter_value = var.vswitch_cidr_block
  }

  # Sequential deployment across regions
  operation_preferences {
    max_concurrent_count    = 2
    failure_tolerance_count = 0
    region_concurrency_type = "SEQUENTIAL"
  }

  timeout_in_minutes    = var.operation_timeout
  operation_description = "Enterprise-wide infrastructure rollout"
}
*/
