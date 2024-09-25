// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package db2

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/IBM-Cloud/bluemix-go/models"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/service/resourcecontroller"
	rc "github.com/IBM/platform-services-go-sdk/resourcecontrollerv2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	RsInstanceSuccessStatus       = "active"
	RsInstanceProgressStatus      = "in progress"
	RsInstanceProvisioningStatus  = "provisioning"
	RsInstanceInactiveStatus      = "inactive"
	RsInstanceFailStatus          = "failed"
	RsInstanceRemovedStatus       = "removed"
	RsInstanceReclamation         = "pending_reclamation"
	RsInstanceUpdateSuccessStatus = "succeeded"
)

func ResourceIBMDb2() *schema.Resource {
	riSchema := resourcecontroller.ResourceIBMResourceInstance().Schema

	riSchema["autoscaling_config"] = &schema.Schema{
		Description: "The db2 auto scaling config",
		Optional:    true,
		Type:        schema.TypeList,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"db2_auto_scaling_threshold": {
					Description: "The db2_auto_scaling_threshold of the instance",
					Required:    true,
					Type:        schema.TypeString,
				},
				"db2_auto_scaling_over_time_period": {
					Description: "The db2_auto_scaling_threshold of the instance",
					Required:    true,
					Type:        schema.TypeString,
				},
			},
		},
	}

	riSchema["whitelist_config"] = &schema.Schema{
		Description: "The db2 whitelist config",
		Optional:    true,
		Type:        schema.TypeList,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"db2_ip_whitelist": {
					Description: "The db2_ip_whitelist of the instance",
					Required:    true,
					Type:        schema.TypeString,
				},
				"db2_whitelist_description": {
					Description: "The db2_whitelist_description of the instance",
					Required:    true,
					Type:        schema.TypeString,
				},
			},
		},
	}

	riSchema["db2_userdetails"] = &schema.Schema{
		Description: "The db2 user details",
		Optional:    true,
		Type:        schema.TypeList,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"db2_userid": {
					Description: "The db2_userid of the instance",
					Required:    true,
					Type:        schema.TypeString,
				},
				"db2_mailid": {
					Description: "The db2_mailid of the instance",
					Required:    true,
					Type:        schema.TypeString,
				},
				"db2_username": {
					Description: "The db2_username of the instance",
					Required:    true,
					Type:        schema.TypeString,
				},
				"db2_role": {
					Description: "The db2_role of the instance",
					Required:    true,
					Type:        schema.TypeString,
				},
				"db2_password": {
					Description: "The db2_password of the instance",
					Required:    true,
					Type:        schema.TypeString,
				},
			},
		},
	}

	return &schema.Resource{
		Create:   resourceIBMDb2Create,
		Read:     resourcecontroller.ResourceIBMResourceInstanceRead,
		Update:   resourcecontroller.ResourceIBMResourceInstanceUpdate,
		Delete:   resourcecontroller.ResourceIBMResourceInstanceDelete,
		Exists:   resourcecontroller.ResourceIBMResourceInstanceExists,
		Importer: &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		CustomizeDiff: customdiff.Sequence(
			func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
				return flex.ResourceTagsCustomizeDiff(diff)
			},
		),

		Schema: riSchema,
	}
}

func addWhitelistIP(crn string, endpoint string, oauthtoken string, ip interface{}, ip_description string) {
	url := endpoint + "/dbapi/v4/dbsettings/whitelistips"
	payload := fmt.Sprintf("{\"ip_addresses\":[{\"address\":\"%s\",\"description\":\"%s\"}]}", ip, ip_description)
	req, _ := http.NewRequest("POST", url, strings.NewReader(payload))
	req.Header.Add("x-deployment-id", crn)
	req.Header.Add("Authorization", "Bearer "+oauthtoken)
	req.Header.Add("Content-Type", "application/json")

	res, _ := http.DefaultClient.Do(req)
	if res.StatusCode != http.StatusOK {
		fmt.Printf("Error Creating user %s", res.Status)
	}
	log.Printf("Status Code : %d", res.StatusCode)
	fmt.Printf("Status Code : %d", res.StatusCode)
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	log.Println(string(body))
	fmt.Println(string(body))
}

func autoScaling(crn string, endpoint string, oauthtoken string, auto_scaling_threshold int, auto_scaling_over_time_period int) {
	url := endpoint + "/dbapi/v4/manage/scaling/auto"
	payload := fmt.Sprintf("{\"auto_scaling_enabled\":\"YES\",\"auto_scaling_allow_plan_limit\":\"YES\",\"auto_scaling_threshold\":%d,\"auto_scaling_over_time_period\":%d}", auto_scaling_threshold, auto_scaling_over_time_period)
	req, _ := http.NewRequest("PUT", url, strings.NewReader(payload))
	req.Header.Add("x-deployment-id", crn)
	req.Header.Add("Authorization", "Bearer "+oauthtoken)
	req.Header.Add("Content-Type", "application/json")

	res, _ := http.DefaultClient.Do(req)
	if res.StatusCode != http.StatusOK {
		fmt.Printf("Error Creating user %s", res.Status)
	}
	log.Printf("Status Code : %d", res.StatusCode)
	fmt.Printf("Status Code : %d", res.StatusCode)
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	log.Println(string(body))
	fmt.Println(string(body))
}

func createUser(crn string, endpoint string, oauthtoken string, db2_userid string, db2_username string, db2_role string, db2_mailid string, db2_password string) {
	url := endpoint + "/dbapi/v4/users"
	payload := fmt.Sprintf("{\"id\":\"%s\",\"name\":\"%s\",\"role\":\"%s\",\"password\":\"%s\",\"email\":\"%s\",\"authentication\":{\"method\":\"internal\",\"policy_id\":\"Default\"}}", db2_userid, db2_username, db2_role, db2_password, db2_mailid)
	log.Println(string(payload))
	fmt.Println(string(payload))
	fmt.Println(string(oauthtoken))
	req, _ := http.NewRequest("POST", url, strings.NewReader(payload))
	req.Header.Add("x-deployment-id", crn)
	req.Header.Add("Authorization", "Bearer "+oauthtoken)
	req.Header.Add("Content-Type", "application/json")

	res, _ := http.DefaultClient.Do(req)
	if res.StatusCode != http.StatusOK {
		fmt.Printf("Error Creating user %s", res.Status)
	}
	log.Printf("Status Code : %d", res.StatusCode)
	fmt.Printf("Status Code : %d", res.StatusCode)
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	log.Println(string(body))
	fmt.Println(string(body))
}

func resourceIBMDb2Create(d *schema.ResourceData, meta interface{}) error {
	serviceName := d.Get("service").(string)
	plan := d.Get("plan").(string)
	name := d.Get("name").(string)
	location := d.Get("location").(string)

	rsConClient, err := meta.(conns.ClientSession).ResourceControllerV2API()
	if err != nil {
		return err
	}

	rsInst := rc.CreateResourceInstanceOptions{
		Name: &name,
	}

	rsCatClient, err := meta.(conns.ClientSession).ResourceCatalogAPI()
	if err != nil {
		return err
	}
	rsCatRepo := rsCatClient.ResourceCatalog()

	serviceOff, err := rsCatRepo.FindByName(serviceName, true)
	if err != nil {
		return fmt.Errorf("[ERROR] Error retrieving service offering: %s", err)
	}

	if metadata, ok := serviceOff[0].Metadata.(*models.ServiceResourceMetadata); ok {
		if !metadata.Service.RCProvisionable {
			return fmt.Errorf("%s cannot be provisioned by resource controller", serviceName)
		}
	} else {
		return fmt.Errorf("[ERROR] Cannot create instance of resource %s\nUse 'ibm_service_instance' if the resource is a Cloud Foundry service", serviceName)
	}

	servicePlan, err := rsCatRepo.GetServicePlanID(serviceOff[0], plan)
	if err != nil {
		return fmt.Errorf("[ERROR] Error retrieving plan: %s", err)
	}
	rsInst.ResourcePlanID = &servicePlan

	deployments, err := rsCatRepo.ListDeployments(servicePlan)
	if err != nil {
		return fmt.Errorf("[ERROR] Error retrieving deployment for plan %s : %s", plan, err)
	}
	if len(deployments) == 0 {
		return fmt.Errorf("[ERROR] No deployment found for service plan : %s", plan)
	}
	deployments, supportedLocations := resourcecontroller.FilterDeployments(deployments, location)

	if len(deployments) == 0 {
		locationList := make([]string, 0, len(supportedLocations))
		for l := range supportedLocations {
			locationList = append(locationList, l)
		}
		return fmt.Errorf("[ERROR] No deployment found for service plan %s at location %s.\nValid location(s) are: %q.\nUse 'ibm_service_instance' if the service is a Cloud Foundry service", plan, location, locationList)

	}
	fmt.Print("Service Off:")
	fmt.Println(serviceOff)
	fmt.Printf("Service plan %s", servicePlan)
	log.Printf("Service plan %s", servicePlan)
	fmt.Print("Deployments :")
	fmt.Println(deployments)

	rsInst.Target = &deployments[0].CatalogCRN

	if rsGrpID, ok := d.GetOk("resource_group_id"); ok {
		rg := rsGrpID.(string)
		rsInst.ResourceGroup = &rg
	} else {
		defaultRg, err := flex.DefaultResourceGroup(meta)
		if err != nil {
			return err
		}
		rsInst.ResourceGroup = &defaultRg
	}

	params := map[string]interface{}{}

	if serviceEndpoints, ok := d.GetOk("service_endpoints"); ok {
		params["service-endpoints"] = serviceEndpoints.(string)
	}

	if parameters, ok := d.GetOk("parameters"); ok {
		temp := parameters.(map[string]interface{})
		for k, v := range temp {
			if v == "true" || v == "false" {
				b, _ := strconv.ParseBool(v.(string))
				params[k] = b
			} else if strings.HasPrefix(v.(string), "[") && strings.HasSuffix(v.(string), "]") {
				//transform v.(string) to be []string
				arrayString := v.(string)
				result := []string{}
				trimLeft := strings.TrimLeft(arrayString, "[")
				trimRight := strings.TrimRight(trimLeft, "]")
				if len(trimRight) == 0 {
					params[k] = result
				} else {
					array := strings.Split(trimRight, ",")
					for _, a := range array {
						result = append(result, strings.Trim(a, "\""))
					}
					params[k] = result
				}
			} else {
				params[k] = v
			}
		}

	}
	if s, ok := d.GetOk("parameters_json"); ok {
		json.Unmarshal([]byte(s.(string)), &params)
	}

	rsInst.Parameters = params

	//Start to create resource instance
	instance, resp, err := rsConClient.CreateResourceInstance(&rsInst)
	if err != nil {
		log.Printf(
			"Error when creating resource instance: %s, Instance info  NAME->%s, LOCATION->%s, GROUP_ID->%s, PLAN_ID->%s",
			err, *rsInst.Name, *rsInst.Target, *rsInst.ResourceGroup, *rsInst.ResourcePlanID)
		return fmt.Errorf("[ERROR] Error when creating resource instance: %s with resp code: %s", err, resp)
	}

	d.SetId(*instance.ID)

	_, err = waitForResourceInstanceCreate(d, meta)
	if err != nil {
		return fmt.Errorf("[ERROR] Error waiting for create resource instance (%s) to be succeeded: %s", d.Id(), err)
	}

	log.Printf("Instance ID %s", *instance.ID)
	log.Printf("Instance CRN %s", *instance.CRN)
	log.Printf("Instance URL %s", *instance.URL)
	log.Printf("Instance DashboardURL %s", *instance.DashboardURL)

	endpoint := strings.Split(*instance.DashboardURL, "databases.appdomain.cloud")[0] + "cloud.ibm.com"

	sess, err := meta.(conns.ClientSession).BluemixSession()
	if err != nil {
		return err
	}
	oauthtoken := sess.Config.IAMAccessToken
	oauthtoken = strings.Replace(oauthtoken, "Bearer ", "", -1)

	if whitelistConfigRaw, ok := d.GetOk("whitelist_config"); ok {
		if whitelistConfigRaw == nil || reflect.ValueOf(whitelistConfigRaw).IsNil() {
			fmt.Print("No whitelistConfig paramas provided; skipping")
		} else {
			whitelistConfig := whitelistConfigRaw.([]interface{})[0].(map[string]interface{})
			fmt.Print(whitelistConfig)
			fmt.Print(whitelistConfig["db2_ip_whitelist"].(string))
			addWhitelistIP(*instance.ID, endpoint, oauthtoken, whitelistConfig["db2_ip_whitelist"].(string), whitelistConfig["db2_whitelist_description"].(string))
		}
	}

	if autoscalingConfigRaw, ok := d.GetOk("autoscaling_config"); ok {
		if autoscalingConfigRaw == nil || reflect.ValueOf(autoscalingConfigRaw).IsNil() {
			fmt.Print("No Autoscaling paramas provided; skipping")
		} else {
			autoscalingConfig := autoscalingConfigRaw.([]interface{})[0].(map[string]interface{})
			fmt.Print(autoscalingConfig["db2_auto_scaling_over_time_period"].(string))
			db2_auto_scaling_threshold, err := strconv.Atoi(autoscalingConfig["db2_auto_scaling_threshold"].(string))
			if err != nil {
				return err
			}
			db2_auto_scaling_over_time_period, err := strconv.Atoi(autoscalingConfig["db2_auto_scaling_over_time_period"].(string))
			if err != nil {
				return err
			}
			autoScaling(*instance.ID, endpoint, oauthtoken, db2_auto_scaling_threshold, db2_auto_scaling_over_time_period)
		}
	}

	if db2userDetailsRaw, ok := d.GetOk("db2_userdetails"); ok {
		if db2userDetailsRaw == nil || reflect.ValueOf(db2userDetailsRaw).IsNil() {
			fmt.Print("No db2 paramas provided; skipping")
		} else {
			db2userDetais := db2userDetailsRaw.([]interface{})[0].(map[string]interface{})
			fmt.Print(db2userDetailsRaw)
			createUser(*instance.ID, endpoint, oauthtoken, db2userDetais["db2_userid"].(string), db2userDetais["db2_username"].(string), db2userDetais["db2_role"].(string), db2userDetais["db2_mailid"].(string), db2userDetais["db2_password"].(string))
		}
	}

	v := os.Getenv("IC_ENV_TAGS")
	if _, ok := d.GetOk("tags"); ok || v != "" {
		oldList, newList := d.GetChange("tags")
		err = flex.UpdateTagsUsingCRN(oldList, newList, meta, *instance.CRN)
		if err != nil {
			log.Printf(
				"Error on create of resource instance (%s) tags: %s", d.Id(), err)
		}
	}

	return resourcecontroller.ResourceIBMResourceInstanceRead(d, meta)

}

func waitForResourceInstanceCreate(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	rsConClient, err := meta.(conns.ClientSession).ResourceControllerV2API()
	if err != nil {
		return false, err
	}
	instanceID := d.Id()
	resourceInstanceGet := rc.GetResourceInstanceOptions{
		ID: &instanceID,
	}

	stateConf := &retry.StateChangeConf{
		Pending: []string{RsInstanceProgressStatus, RsInstanceInactiveStatus, RsInstanceProvisioningStatus},
		Target:  []string{RsInstanceSuccessStatus},
		Refresh: func() (interface{}, string, error) {
			instance, resp, err := rsConClient.GetResourceInstance(&resourceInstanceGet)
			if err != nil {
				if resp != nil && resp.StatusCode == 404 {
					return nil, "", fmt.Errorf("[ERROR] The resource instance %s does not exist anymore: %v", d.Id(), err)
				}
				return nil, "", fmt.Errorf("[ERROR] Get the resource instance %s failed with resp code: %s, err: %v", d.Id(), resp, err)
			}
			if *instance.State == RsInstanceFailStatus {
				return instance, *instance.State, fmt.Errorf("[ERROR] The resource instance '%s' creation failed: %v", d.Id(), err)
			}
			return instance, *instance.State, nil
		},
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForStateContext(context.Background())
}
