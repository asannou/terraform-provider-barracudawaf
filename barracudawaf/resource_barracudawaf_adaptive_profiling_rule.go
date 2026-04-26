package barracudawaf

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCudaWAFAdaptiveProfilingRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceCudaWAFAdaptiveProfilingRuleCreate,
		Read:   resourceCudaWAFAdaptiveProfilingRuleRead,
		Update: resourceCudaWAFAdaptiveProfilingRuleUpdate,
		Delete: resourceCudaWAFAdaptiveProfilingRuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				parts := strings.Split(d.Id(), "/")
				if len(parts) != 2 {
					return nil, fmt.Errorf("invalid ID specified. Supposed to be <service_name>/<rule_name>")
				}
				d.Set("parent", []string{parts[0]})
				d.SetId(parts[1])
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"name":                {Type: schema.TypeString, Required: true, Description: "Learn Rule Name"},
			"host":                {Type: schema.TypeString, Required: true, Description: "Host Match"},
			"url":                 {Type: schema.TypeString, Required: true, Description: "URL Match"},
			"status":              {Type: schema.TypeString, Optional: true, Description: "Status"},
			"learn_from_request":  {Type: schema.TypeString, Optional: true, Description: "Learn From Request"},
			"learn_from_response": {Type: schema.TypeString, Optional: true, Description: "Learn From Response"},
			"parent": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
		},

		Description: "`barracudawaf_adaptive_profiling_rule` manages `Adaptive Profiling Rule` on the Barracuda Web Application Firewall.",
	}
}

func resourceCudaWAFAdaptiveProfilingRuleCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Get("name").(string)

	log.Println("[INFO] Creating Barracuda WAF resource " + name)

	parent := d.Get("parent").([]interface{})
	serviceName := parent[0].(string)

	resourceEndpoint := "/services/" + serviceName + "/adaptive-profiling-rules"
	err := client.CreateBarracudaWAFResource(name, hydrateBarracudaWAFAdaptiveProfilingRuleResource(d, "post", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF resource (%s) (%v) ", name, err)
		return err
	}

	d.SetId(name)
	return resourceCudaWAFAdaptiveProfilingRuleRead(d, m)
}

func resourceCudaWAFAdaptiveProfilingRuleRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()
	log.Println("[INFO] Fetching Barracuda WAF resource " + name)

	parent := d.Get("parent").([]interface{})
	serviceName := parent[0].(string)

	resourceEndpoint := "/services/" + serviceName + "/adaptive-profiling-rules"
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
		if dataItems["name"] == name {
			break
		}
	}

	if dataItems["name"] != name {
		return fmt.Errorf("Barracuda WAF resource (%s) not found on the system", name)
	}

	d.Set("name", name)
	d.Set("host", dataItems["host"])
	d.Set("url", dataItems["url"])
	d.Set("status", dataItems["status"])
	d.Set("learn_from_request", dataItems["learn-from-request"])
	d.Set("learn_from_response", dataItems["learn-from-response"])

	return nil
}

func resourceCudaWAFAdaptiveProfilingRuleUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Updating Barracuda WAF resource " + name)

	parent := d.Get("parent").([]interface{})
	serviceName := parent[0].(string)

	resourceEndpoint := "/services/" + serviceName + "/adaptive-profiling-rules"
	err := client.UpdateBarracudaWAFResource(name, hydrateBarracudaWAFAdaptiveProfilingRuleResource(d, "put", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF resource (%s) (%v)", name, err)
		return err
	}

	return resourceCudaWAFAdaptiveProfilingRuleRead(d, m)
}

func resourceCudaWAFAdaptiveProfilingRuleDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Deleting Barracuda WAF resource " + name)

	parent := d.Get("parent").([]interface{})
	serviceName := parent[0].(string)

	resourceEndpoint := "/services/" + serviceName + "/adaptive-profiling-rules"
	request := &APIRequest{
		Method: "delete",
		URL:    resourceEndpoint,
	}

	err := client.DeleteBarracudaWAFResource(name, request)

	if err != nil {
		return fmt.Errorf("Unable to delete the Barracuda WAF resource (%s) (%v)", name, err)
	}

	return nil
}

func hydrateBarracudaWAFAdaptiveProfilingRuleResource(d *schema.ResourceData, method string, endpoint string) *APIRequest {

	//resourcePayload : payload for the resource
	resourcePayload := map[string]string{
		"name":                d.Get("name").(string),
		"host":                d.Get("host").(string),
		"url":                 d.Get("url").(string),
		"status":              d.Get("status").(string),
		"learn-from-request":  d.Get("learn_from_request").(string),
		"learn-from-response": d.Get("learn_from_response").(string),
	}

	// parameters not supported for updates
	if method == "put" {
		updatePayloadExceptions := [...]string{"name"}
		for _, param := range updatePayloadExceptions {
			delete(resourcePayload, param)
		}
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
