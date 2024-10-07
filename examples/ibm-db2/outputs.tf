output "db2_hostname" {
  value     = ibm_resource_key.db2.credentials["connection.db2.host_ros.0"]
  sensitive = true
}
output "db2_user" {
  value     = ibm_resource_key.db2.credentials["connection.db2.authentication.username"]
  sensitive = true
}
output "db2_password" {
  value     = ibm_resource_key.db2.credentials["connection.db2.authentication.password"]
  sensitive = true
}
output "db2_cert_base64" {
  value     = ibm_resource_key.db2.credentials["connection.db2.certificate.certificate_base64"]
  sensitive = true
}