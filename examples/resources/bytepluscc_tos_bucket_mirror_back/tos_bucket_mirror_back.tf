resource "bytepluscc_tos_bucket_mirror_back" "Example" {
  bucket = "example-bucket"
  rules = [{
    condition = {
      http_code   = 404
      http_method = ["GET", "HEAD"]
      key_prefix  = "object-key-prefix"
      key_suffix  = "object-key-suffix"
      allow_host  = ["assets.example.com"]
    }
    id = "example-rule"
    redirect = {
      fetch_header_to_meta_data_rules = [{
        meta_data_suffix = "-meta"
        source_header    = "X-Custom-Header"
      }]
      fetch_source_on_redirect            = false
      fetch_source_on_redirect_with_query = false
      follow_redirect                     = true
      mirror_header = {
        pass_all = true
        pass     = ["aaa", "bbb"]
        remove   = ["xxx", "yyy"]
        set = [{
          key   = "X-Set-Header"
          value = "set-value"
        }]
      }
      pass_header_from_source      = ["X-Source-Header"]
      pass_query                   = true
      pass_status_code_from_source = [404, 500]
      private_source = {
        source_endpoint = {
          # Replace the endpoints and credentials with those of your source buckets.
          primary = [{
            bucket_name = "primary-bucket"
            endpoint    = "https://primary.example.com"
            credential_provider = {
              static_credential = {
                storage_vendor  = "S3"
                sk_encrypt_type = ""
                sk              = "replace-with-primary-secret-key"
                ak              = "replace-with-primary-access-key"
              }
              role   = ""
              region = "ap-southeast-1"
            }
          }]
          follower = [{
            bucket_name = "follower-bucket"
            endpoint    = "https://follower.example.com"
            credential_provider = {
              static_credential = {
                storage_vendor  = "S3"
                sk_encrypt_type = ""
                sk              = "replace-with-follower-secret-key"
                ak              = "replace-with-follower-access-key"
              }
              role   = ""
              region = "ap-southeast-1"
            }
          }]
        }
      }
      redirect_type = "Mirror"
      transform = {
        with_key_prefix = "additional-key-prefix"
        with_key_suffix = "additional-key-suffix"
        replace_key_prefix = {
          key_prefix   = "key-prefix"
          replace_with = "replace-with"
        }
      }
    }
  }]
}
