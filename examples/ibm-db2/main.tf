
provider "ibm" {
  ibmcloud_api_key = var.ibmcloud_api_key
}

resource "ibm_db2" "db2_instance" {
  name              = "demo-db2-prasad22"
  service           = "dashdb-for-transactions"
  plan              = "performance-dev" 
  location          = "us-east"
  resource_group_id = "0f39969ff2da4ec986cd89e4684bb181"
  service_endpoints = "public"
  parameters_json   = <<EOF
    {
        "version": "12",
        "node_type": null
    }
  EOF

  timeouts {
    create = "720m"
    update = "30m"
    delete = "30m"
  }
}


resource "ibm_resource_key" "db2" {
  name                 = "default-db2-creds"
  role                 = "Manager"
  resource_instance_id = ibm_db2.db2_instance.id
}

/*
    autoscaling_config {
      db2_auto_scaling_threshold = "60"
      db2_auto_scaling_over_time_period = "15"
    }
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
