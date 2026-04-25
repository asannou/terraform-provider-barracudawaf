package barracudawaf

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCudaWAFCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceCudaWAFClusterUpdate,
		Read:   resourceCudaWAFClusterRead,
		Update: resourceCudaWAFClusterUpdate,
		Delete: schema.Noop,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"cluster_name":                  {Type: schema.TypeString, Required: true, Description: "Cluster Name"},
			"cluster_shared_secret":         {Type: schema.TypeString, Optional: true, Description: "Cluster Shared Secret", Sensitive: true},
			"data_path_failure_action":      {Type: schema.TypeString, Optional: true, Description: "Data-path failure action"},
			"failback_mode":                 {Type: schema.TypeString, Optional: true, Description: "Failback Mode"},
			"heartbeat_count_per_interface": {Type: schema.TypeString, Optional: true, Description: "Heartbeat Count Per Interface"},
			"heartbeat_frequency":           {Type: schema.TypeString, Optional: true, Description: "Heartbeat Frequency"},
			"monitor_link":                  {Type: schema.TypeString, Optional: true, Description: "Monitor Link"},
			"transmit_heartbeat_on":         {Type: schema.TypeString, Optional: true, Description: "Transmit Heartbeat On"},
			"vx_aa_enable":                  {Type: schema.TypeString, Optional: true, Description: "Active-Active HA Clustering"},
		},

		Description: "`barracudawaf_cluster` manages `Cluster` on the Barracuda Web Application Firewall.",
	}
}

func resourceCudaWAFClusterRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()
	log.Println("[INFO] Fetching Barracuda WAF resource " + name)

	resourceEndpoint := "/cluster"
	request := &APIRequest{
		Method: "get",
		URL:    resourceEndpoint,
	}

	resources, err := client.GetBarracudaWAFResource(name, request)

	if err != nil {
		log.Printf("[ERROR] Unable to Retrieve Barracuda WAF resource (%s) (%v) ", name, err)
		return err
	}

	if resources.Data == nil {
		log.Printf("[WARN] Barracuda WAF resource (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	var dataItems map[string]interface{}
	dataItems = resources.Data["Cluster"]
	if dataItems == nil {
		// Fallback for different API structure if needed
		for _, v := range resources.Data {
			dataItems = v
			break
		}
	}

	if dataItems == nil {
		return fmt.Errorf("Barracuda WAF resource (%s) not found on the system", name)
	}

	d.Set("cluster_name", dataItems["cluster-name"])
	d.Set("data_path_failure_action", dataItems["data-path-failure-action"])
	d.Set("failback_mode", dataItems["failback-mode"])
	d.Set("heartbeat_count_per_interface", fmt.Sprintf("%v", dataItems["heartbeat-count-per-interface"]))
	d.Set("heartbeat_frequency", fmt.Sprintf("%v", dataItems["heartbeat-frequency"]))
	d.Set("monitor_link", dataItems["monitor-link"])
	d.Set("transmit_heartbeat_on", dataItems["transmit-heartbeat-on"])
	d.Set("vx_aa_enable", dataItems["vx-aa-enable"])

	return nil
}

func resourceCudaWAFClusterUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Get("cluster_name").(string)
	if d.Id() == "" {
		d.SetId(name)
	}

	log.Println("[INFO] Updating Barracuda WAF resource " + name)

	resourceEndpoint := "/cluster"
	request := hydrateBarracudaWAFClusterResource(d, "put", resourceEndpoint)
	_, err := client.putReq(request.Body, resourceEndpoint)

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF resource (%s) (%v)", name, err)
		return err
	}

	return resourceCudaWAFClusterRead(d, m)
}

func hydrateBarracudaWAFClusterResource(d *schema.ResourceData, method string, endpoint string) *APIRequest {

	//resourcePayload : payload for the resource
	resourcePayload := map[string]string{
		"cluster-name":                  d.Get("cluster_name").(string),
		"cluster-shared-secret":         d.Get("cluster_shared_secret").(string),
		"data-path-failure-action":      d.Get("data_path_failure_action").(string),
		"failback-mode":                 d.Get("failback_mode").(string),
		"heartbeat-count-per-interface": d.Get("heartbeat_count_per_interface").(string),
		"heartbeat-frequency":           d.Get("heartbeat_frequency").(string),
		"monitor-link":                  d.Get("monitor_link").(string),
		"transmit-heartbeat-on":         d.Get("transmit_heartbeat_on").(string),
		"vx-aa-enable":                  d.Get("vx_aa_enable").(string),
	}

	// remove empty parameters from resource payload
	for key, val := range resourcePayload {
		if len(val) == 0 {
			delete(resourcePayload, key)
		}
	}

	return &APIRequest{
		URL:  endpoint,
		Body: resourcePayload,
	}
}
