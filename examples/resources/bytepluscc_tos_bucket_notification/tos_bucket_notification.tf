resource "bytepluscc_tos_bucket_notification" "Example" {
  bucket_name = "example-bucket"
  notification_rules = [{
    # Replace the destination identifiers and roles with your own values.
    destination = {
      kafka = [{
        instance_id = "replace-with-kafka-instance-id"
        region      = "ap-southeast-1"
        role        = "replace-with-kafka-role-resource-id"
        topic       = "example-topic"
        user        = "example-user"
      }]
      rocket_mq = [{
        access_key_id = "replace-with-rocketmq-access-key-id"
        instance_id   = "replace-with-rocketmq-instance-id"
        role          = "replace-with-rocketmq-role-resource-id"
        topic         = "example-topic"
      }]
      ve_faa_s = [{
        function_id = "replace-with-function-id"
      }]
    }
    events = ["tos:ObjectCreated:Put", "tos:LifecycleExpiration:Delete"]
    filter = {
      tos_key = {
        filter_rules = [{
          name  = "prefix"
          value = "prefix"
          }, {
          name  = "suffix"
          value = "suffix"
        }]
      }
    }
    rule_id = "example-rule"
  }]
}
