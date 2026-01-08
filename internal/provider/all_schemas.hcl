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

resource_schema "byteplus_vke_permission" {
  cloudcontrol_type_name = "Byteplus::VKE::Permission"
}

resource_schema "byteplus_vefaas_sandbox" {
  cloudcontrol_type_name = "Byteplus::VEFAAS::Sandbox"
}

resource_schema "byteplus_vefaas_kafka_trigger" {
  cloudcontrol_type_name = "Byteplus::VEFAAS::KafkaTrigger"
}

resource_schema "byteplus_cdn_share_config" {
  cloudcontrol_type_name = "Byteplus::CDN::ShareConfig"
}

resource_schema "byteplus_apig_upstream" {
  cloudcontrol_type_name = "Byteplus::APIG::Upstream"
}

resource_schema "byteplus_apig_upstream_source" {
  cloudcontrol_type_name = "Byteplus::APIG::UpstreamSource"
}

resource_schema "byteplus_apig_gateway_service" {
  cloudcontrol_type_name = "Byteplus::APIG::GatewayService"
}

resource_schema "byteplus_iam_policy" {
  cloudcontrol_type_name = "Byteplus::IAM::Policy"
}

resource_schema "byteplus_iam_role" {
  cloudcontrol_type_name = "Byteplus::IAM::Role"
}

resource_schema "byteplus_iam_user" {
  cloudcontrol_type_name = "Byteplus::IAM::User"
}

resource_schema "byteplus_clb_listener" {
  cloudcontrol_type_name = "Byteplus::CLB::Listener"
}

resource_schema "byteplus_natgateway_ngw" {
  cloudcontrol_type_name = "Byteplus::NatGateway::NGW"
}

resource_schema "byteplus_organization_unit" {
  cloudcontrol_type_name = "Byteplus::Organization::Unit"
}

resource_schema "byteplus_rdsmysql_db_account" {
  cloudcontrol_type_name = "Byteplus::RDSMySQL::DBAccount"
}

resource_schema "byteplus_autoscaling_scaling_group" {
  cloudcontrol_type_name = "Byteplus::AutoScaling::ScalingGroup"
}

resource_schema "byteplus_cen_cen" {
  cloudcontrol_type_name = "Byteplus::CEN::CEN"
}

resource_schema "byteplus_natgateway_snatentry" {
  cloudcontrol_type_name = "Byteplus::NatGateway::Snatentry"
}

resource_schema "byteplus_natgateway_dnatentry" {
  cloudcontrol_type_name = "Byteplus::NatGateway::Dnatentry"
}

resource_schema "byteplus_natgateway_nat_ip" {
  cloudcontrol_type_name = "Byteplus::NatGateway::NatIp"
}

resource_schema "byteplus_vke_cluster" {
  cloudcontrol_type_name = "Byteplus::VKE::Cluster"
}

resource_schema "byteplus_rdsmysql_instance" {
  cloudcontrol_type_name = "Byteplus::RDSMySQL::Instance"
}

resource_schema "byteplus_escloud_instance" {
  cloudcontrol_type_name = "Byteplus::ESCloud::Instance"
}

resource_schema "byteplus_kms_key" {
  cloudcontrol_type_name = "Byteplus::KMS::Key"
}

resource_schema "byteplus_apig_gateway" {
  cloudcontrol_type_name = "Byteplus::APIG::Gateway"
}

resource_schema "byteplus_apig_custom_domain" {
  cloudcontrol_type_name = "Byteplus::APIG::CustomDomain"
}

resource_schema "byteplus_hbase_instance" {
  cloudcontrol_type_name = "Byteplus::HBase::Instance"
}

resource_schema "byteplus_alb_certificate" {
  cloudcontrol_type_name = "Byteplus::ALB::Certificate"
}


resource_schema "byteplus_vke_addon" {
  cloudcontrol_type_name = "Byteplus::VKE::Addon"
}

resource_schema "byteplus_vke_node_pool" {
  cloudcontrol_type_name = "Byteplus::VKE::NodePool"
}

resource_schema "byteplus_vpc_security_group" {
  cloudcontrol_type_name = "Byteplus::VPC::SecurityGroup"
}

resource_schema "byteplus_vpc_eni" {
  cloudcontrol_type_name = "Byteplus::VPC::ENI"
}

resource_schema "byteplus_redis_instance" {
  cloudcontrol_type_name = "Byteplus::Redis::Instance"
}

resource_schema "byteplus_vepfs_instance" {
  cloudcontrol_type_name = "Byteplus::VEPFS::Instance"
}

resource_schema "byteplus_storageebs_snapshot" {
  cloudcontrol_type_name = "Byteplus::StorageEBS::Snapshot"
}

resource_schema "byteplus_kafka_topic" {
  cloudcontrol_type_name = "Byteplus::Kafka::Topic"
}

resource_schema "byteplus_dns_zone" {
  cloudcontrol_type_name = "Byteplus::DNS::Zone"
}

resource_schema "byteplus_iam_project" {
  cloudcontrol_type_name = "Byteplus::IAM::Project"
}

resource_schema "byteplus_cloudidentity_user" {
  cloudcontrol_type_name = "Byteplus::CloudIdentity::User"
}

resource_schema "byteplus_bmq_instance" {
  cloudcontrol_type_name = "Byteplus::BMQ::Instance"
}

resource_schema "byteplus_privatelink_endpoint_service" {
  cloudcontrol_type_name = "Byteplus::PrivateLink::EndpointService"
}

resource_schema "byteplus_fwcenter_dns_control_policy" {
  cloudcontrol_type_name = "Byteplus::FWCenter::DnsControlPolicy"
}


resource_schema "byteplus_clb_rule" {
  cloudcontrol_type_name = "Byteplus::CLB::Rule"
}

resource_schema "byteplus_vpc_subnet" {
  cloudcontrol_type_name = "Byteplus::VPC::Subnet"
}

resource_schema "byteplus_bmq_group" {
  cloudcontrol_type_name = "Byteplus::BMQ::Group"
}

resource_schema "byteplus_cloudmonitor_rule" {
  cloudcontrol_type_name = "Byteplus::CloudMonitor::Rule"
}

resource_schema "byteplus_clb_nlb_server_group" {
  cloudcontrol_type_name = "Byteplus::CLB::NLBServerGroup"
}

resource_schema "byteplus_clb_nlb_listener" {
  cloudcontrol_type_name = "Byteplus::CLB::NLBListener"
}

resource_schema "byteplus_clb_server_group" {
  cloudcontrol_type_name = "Byteplus::CLB::ServerGroup"
}

resource_schema "byteplus_directconnect_direct_connect_gateway" {
  cloudcontrol_type_name = "Byteplus::DirectConnect::DirectConnectGateway"
}

resource_schema "byteplus_vpc_eip" {
  cloudcontrol_type_name = "Byteplus::VPC::EIP"
}

resource_schema "byteplus_vpc_route_table" {
  cloudcontrol_type_name = "Byteplus::VPC::RouteTable"
}

resource_schema "byteplus_filenas_instance" {
  cloudcontrol_type_name = "Byteplus::FileNAS::Instance"
}


