resource "bytepluscc_rdsmysql_endpoint_public_address" "EndpointPublicAddressDemo" {
  instance_id   = "mysql-f9fe535xxxxx"
  endpoint_id   = "mysql-f9fe53xxxxxa-cluster"
  eip_id        = "eip-37vr7k7r9t3b44etmwxxxxx"
  domain_prefix = "ccapi-terraform"
  port          = 23306
}