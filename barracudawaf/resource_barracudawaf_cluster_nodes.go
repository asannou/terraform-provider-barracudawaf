package barracudawaf

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCudaWAFClusterNodes() *schema.Resource {
	return &schema.Resource{
		Create: resourceCudaWAFClusterNodesUpdate,
		Read:   resourceCudaWAFClusterNodesRead,
		Update: resourceCudaWAFClusterNodesUpdate,
		Delete: schema.Noop,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"ip_address": {Type: schema.TypeString, Required: true, Description: "Cluster System IP"},
		},

		Description: "`barracudawaf_cluster_nodes` manages `Cluster nodes` on the Barracuda Web Application Firewall.",
	}
}

func resourceCudaWAFClusterNodesRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()
	log.Println("[INFO] Fetching Barracuda WAF resource " + name)

	resourceEndpoint := "/cluster/nodes"
	request := &APIRequest{
		Method: "get",
		URL:    resourceEndpoint,
	}

	var dataItems map[string]interface{}
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

	for _, dataItems = range resources.Data {
		if dataItems["ip-address"] == name {
			break
		}
	}

	if dataItems["ip-address"] != name {
		return fmt.Errorf("Barracuda WAF resource (%s) not found on the system", name)
	}

	d.Set("ip_address", name)

	return nil
}

func resourceCudaWAFClusterNodesUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Get("ip_address").(string)
	if d.Id() == "" {
		d.SetId(name)
	}

	log.Println("[INFO] Updating Barracuda WAF resource " + name)

	resourceEndpoint := "/cluster/nodes"
	request := hydrateBarracudaWAFClusterNodesResource(d, "put", resourceEndpoint)
	_, err := client.putReq(request.Body, resourceEndpoint)

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF resource (%s) (%v)", name, err)
		return err
	}

	return resourceCudaWAFClusterNodesRead(d, m)
}

func hydrateBarracudaWAFClusterNodesResource(d *schema.ResourceData, method string, endpoint string) *APIRequest {

	//resourcePayload : payload for the resource
	resourcePayload := map[string]string{
		"ip-address": d.Get("ip_address").(string),
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
