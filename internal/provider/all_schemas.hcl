# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

defaults {
  schema_cache_directory     = "../service/cloudcontrol/schemas"
  terraform_type_name_prefix = "bytepluscc"
}

meta_schema {
  path = "../service/cloudcontrol/meta-schemas/provider.definition.schema.v1.json"
}

# CloudControl resource types schemas are available for use with the Cloud Control API.
resource_schema "byteplus_ecs_image" {
  cloudcontrol_type_name = "Byteplus::ECS::Image"
}

resource_schema "byteplus_clb_listener" {
  cloudcontrol_type_name = "Byteplus::CLB::Listener"
}

resource_schema "byteplus_vpc_vpc" {
  cloudcontrol_type_name = "Byteplus::VPC::VPC"

}

resource_schema "byteplus_storageebs_volume" {
  cloudcontrol_type_name = "Byteplus::StorageEBS::Volume"
}

resource_schema "byteplus_alb_listener" {
  cloudcontrol_type_name = "Byteplus::ALB::Listener"
}

resource_schema "byteplus_cr_repository" {
  cloudcontrol_type_name = "Byteplus::CR::Repository"
}

resource_schema "byteplus_autoscaling_scaling_configuration" {
  cloudcontrol_type_name = "Byteplus::AutoScaling::ScalingConfiguration"
}

resource_schema "byteplus_ecs_instance" {
  cloudcontrol_type_name = "Byteplus::ECS::Instance"
}


