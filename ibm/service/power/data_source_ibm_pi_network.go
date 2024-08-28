// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/helpers"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceIBMPINetwork() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPINetworkRead,
		Schema: map[string]*schema.Schema{
			// Arguments
			Arg_CloudInstanceID: {
				Description:  "The GUID of the service instance associated with an account.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
			Arg_NetworkName: {
				Description:  "The unique identifier or name of a network.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},

			// Attributes
			Attr_AccessConfig: {
				Computed:    true,
				Description: "The network communication configuration option of the network (for satellite locations only).",
				Type:        schema.TypeString,
			},
			Attr_AvailableIPCount: {
				Computed:    true,
				Description: "The total number of IP addresses that you have in your network.",
				Type:        schema.TypeFloat,
			},
			Attr_CIDR: {
				Computed:    true,
				Description: "The CIDR of the network.",
				Type:        schema.TypeString,
			},
			Attr_CRN: {
				Computed:    true,
				Description: "The CRN of this resource.",
				Type:        schema.TypeString,
			},
			Attr_DNS: {
				Computed:    true,
				Description: "The DNS Servers for the network.",
				Elem:        &schema.Schema{Type: schema.TypeString},
				Type:        schema.TypeSet,
			},
			Attr_Gateway: {
				Computed:    true,
				Description: "The network gateway that is attached to your network.",
				Type:        schema.TypeString,
			},
			Attr_Jumbo: {
				Computed:    true,
				Deprecated:  "This field is deprecated, use mtu instead.",
				Description: "MTU Jumbo option of the network (for multi-zone locations only).",
				Type:        schema.TypeBool,
			},
			Attr_MTU: {
				Computed:    true,
				Description: "Maximum Transmission Unit option of the network.",
				Type:        schema.TypeInt,
			},
			Attr_Name: {
				Computed:    true,
				Deprecated:  "This field is deprecated, use pi_network_name instead.",
				Description: "The unique identifier or name of a network.",
				Type:        schema.TypeString,
			},
			Attr_Type: {
				Computed:    true,
				Description: "The type of network.",
				Type:        schema.TypeString,
			},
			Attr_UsedIPCount: {
				Computed:    true,
				Description: "The number of used IP addresses.",
				Type:        schema.TypeFloat,
			},
			Attr_UsedIPPercent: {
				Computed:    true,
				Description: "The percentage of IP addresses used.",
				Type:        schema.TypeFloat,
			},
			Attr_UserTags: {
				Computed:    true,
				Description: "The user tags attached to this resource.",
				Elem:        &schema.Schema{Type: schema.TypeString},
				Type:        schema.TypeList,
			},
			Attr_VLanID: {
				Computed:    true,
				Description: "The VLAN ID that the network is connected to.",
				Type:        schema.TypeInt,
			},
		},
	}
}

func dataSourceIBMPINetworkRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)

	networkC := instance.NewIBMPINetworkClient(ctx, sess, cloudInstanceID)
	networkdata, err := networkC.Get(d.Get(helpers.PINetworkName).(string))
	if err != nil || networkdata == nil {
		return diag.FromErr(err)
	}

	d.SetId(*networkdata.NetworkID)
	d.Set(Attr_AccessConfig, networkdata.AccessConfig)
	if networkdata.IPAddressMetrics.Available != nil {
		d.Set(Attr_AvailableIPCount, networkdata.IPAddressMetrics.Available)
	}
	if networkdata.Cidr != nil {
		d.Set(Attr_CIDR, networkdata.Cidr)
	}
	if networkdata.Crn != "" {
		d.Set(Attr_CRN, networkdata.Crn)
	}
	if len(networkdata.DNSServers) > 0 {
		d.Set(Attr_DNS, networkdata.DNSServers)
	}
	d.Set(Attr_Gateway, networkdata.Gateway)
	d.Set(Attr_Jumbo, networkdata.Jumbo)
	d.Set(Attr_MTU, networkdata.Mtu)
	if networkdata.Name != nil {
		d.Set(Attr_Name, networkdata.Name)
	}
	if networkdata.Type != nil {
		d.Set(Attr_Type, networkdata.Type)
	}
	if networkdata.IPAddressMetrics.Used != nil {
		d.Set(Attr_UsedIPCount, networkdata.IPAddressMetrics.Used)
	}
	if networkdata.IPAddressMetrics.Utilization != nil {
		d.Set(Attr_UsedIPPercent, networkdata.IPAddressMetrics.Utilization)
	}
	if len(networkdata.UserTags) > 0 {
		d.Set(Attr_UserTags, networkdata.UserTags)
	}
	if networkdata.VlanID != nil {
		d.Set(Attr_VLanID, networkdata.VlanID)
	}

	return nil
}
