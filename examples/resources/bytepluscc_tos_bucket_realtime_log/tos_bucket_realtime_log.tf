resource "bytepluscc_tos_bucket_realtime_log" "Example" {
  bucket_name = "example-bucket"
  real_time_log = {
    role              = "TOSLogArchiveTLSRole"
    use_service_topic = true
    tls_project_id    = "replace-with-tls-project-id"
    tls_topic_id      = "replace-with-tls-topic-id"
  }
}
