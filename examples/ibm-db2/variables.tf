variable "ibmcloud_api_key" {
  description = "IBM Cloud API key"
  type        = string
}

variable "db2_password" {
  description = "IBM Db2oC user password"
  type        = string
  sensitive = true
  default = ""
}
