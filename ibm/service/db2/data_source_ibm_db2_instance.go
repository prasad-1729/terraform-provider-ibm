// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package db2

import (
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/service/resourcecontroller"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMDb2() *schema.Resource {
	riSchema := resourcecontroller.DataSourceIBMResourceInstance().Schema

	riSchema["autoscaling_config"] = &schema.Schema{
		Description: "The db2 auto scaling config",
		Computed:    true,
		Type:        schema.TypeList,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"db2_auto_scaling_threshold": {
					Description: "The db2_auto_scaling_threshold of the instance",
					Computed:    true,
					Type:        schema.TypeString,
				},
				"db2_auto_scaling_over_time_period": {
					Description: "The db2_auto_scaling_threshold of the instance",
					Computed:    true,
					Type:        schema.TypeString,
				},
			},
		},
	}

	riSchema["whitelist_config"] = &schema.Schema{
		Description: "The db2 whitelist config",
		Computed:    true,
		Type:        schema.TypeList,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"db2_ip_whitelist": {
					Description: "The db2_ip_whitelist of the instance",
					Computed:    true,
					Type:        schema.TypeString,
				},
				"db2_whitelist_description": {
					Description: "The db2_whitelist_description of the instance",
					Computed:    true,
					Type:        schema.TypeString,
				},
			},
		},
	}

	riSchema["db2_userdetails"] = &schema.Schema{
		Description: "The db2 user details",
		Computed:    true,
		Type:        schema.TypeList,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"db2_userid": {
					Description: "The db2_userid of the instance",
					Computed:    true,
					Type:        schema.TypeString,
				},
				"db2_mailid": {
					Description: "The db2_mailid of the instance",
					Computed:    true,
					Type:        schema.TypeString,
				},
				"db2_username": {
					Description: "The db2_username of the instance",
					Computed:    true,
					Type:        schema.TypeString,
				},
				"db2_role": {
					Description: "The db2_role of the instance",
					Computed:    true,
					Type:        schema.TypeString,
				},
				"db2_password": {
					Description: "The db2_password of the instance",
					Computed:    true,
					Type:        schema.TypeString,
				},
			},
		},
	}

	return &schema.Resource{
		Read:   resourcecontroller.DataSourceIBMResourceInstanceRead,
		Schema: riSchema,
	}
}
