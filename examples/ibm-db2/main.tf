
provider "ibm" {
  ibmcloud_api_key = var.ibmcloud_api_key
}

resource "ibm_db2" "db2_instance" {
  name              = "demo-db2-prasad44"
  service           = "dashdb-for-transactions"
  plan              = "dashdbpreprod" 
  location          = "us-south"
  resource_group_id = "6084241f97f74bc1bf99fdf4f8eb4001"
  service_endpoints = "public"
  autoscaling_config {
      db2_auto_scaling_threshold = "60"
      db2_auto_scaling_over_time_period = "15"
    }
  parameters_json   = <<EOF
    {
        "disk_encryption_instance_crn": "none",
        "disk_encryption_key_crn": "none",
        "high_availability": "no",
        "oracle_compatibility": "no"
    }
  EOF

  timeouts {
    create = "720m"
    update = "30m"
    delete = "30m"
  }
}
/*
  whitelist_config {
      db2_ip_whitelist = "192.168.3.42"
      db2_whitelist_description = "test1"
    }
  db2_userdetails {
    db2_userid = "akash"
    db2_mailid = "akash.david@ibm.com"
    db2_password = var.db2_password
    db2_username = "akashdavid"
    db2_role = "bluuser"
  }
*/
