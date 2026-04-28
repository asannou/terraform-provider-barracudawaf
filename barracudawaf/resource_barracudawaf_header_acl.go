package barracudawaf

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCudaWAFHeaderACL() *schema.Resource {
	return &schema.Resource{
		Create: resourceCudaWAFHeaderACLCreate,
		Read:   resourceCudaWAFHeaderACLRead,
		Update: resourceCudaWAFHeaderACLUpdate,
		Delete: resourceCudaWAFHeaderACLDelete,
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
			"name":                        {Type: schema.TypeString, Required: true, Description: "Header ACL Name"},
			"header_name":                 {Type: schema.TypeString, Required: true, Description: "Header Name"},
			"blocked_attack_types":        {Type: schema.TypeList, Optional: true, Elem: &schema.Schema{Type: schema.TypeString}, Description: "Blocked Attack Types"},
			"comments":                    {Type: schema.TypeString, Optional: true, Description: "Comments"},
			"custom_blocked_attack_types": {Type: schema.TypeList, Optional: true, Elem: &schema.Schema{Type: schema.TypeString}, Description: "Custom Blocked Attack Types"},
			"denied_metachars":            {Type: schema.TypeString, Optional: true, Description: "Denied Metacharacters"},
			"exception_patterns":          {Type: schema.TypeList, Optional: true, Elem: &schema.Schema{Type: schema.TypeString}, Description: "Exception Patterns"},
			"max_header_value_length":     {Type: schema.TypeString, Optional: true, Description: "Max Header Value Length"},
			"mode":                        {Type: schema.TypeString, Optional: true, Description: "Mode"},
			"status":                      {Type: schema.TypeString, Optional: true, Description: "Status"},
			"parent": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
		},

		Description: "`barracudawaf_header_acl` manages `Header ACL` on the Barracuda Web Application Firewall.",
	}
}

func resourceCudaWAFHeaderACLCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Get("name").(string)

	log.Println("[INFO] Creating Barracuda WAF resource " + name)

	parent := d.Get("parent").([]interface{})
	serviceName := parent[0].(string)

	resourceEndpoint := "/services/" + serviceName + "/header-acls"
	err := client.CreateBarracudaWAFResource(name, hydrateBarracudaWAFHeaderACLResource(d, "post", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF resource (%s) (%v) ", name, err)
		return err
	}

	d.SetId(name)
	return resourceCudaWAFHeaderACLRead(d, m)
}

func resourceCudaWAFHeaderACLRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()
	log.Println("[INFO] Fetching Barracuda WAF resource " + name)

	parent := d.Get("parent").([]interface{})
	serviceName := parent[0].(string)

	resourceEndpoint := "/services/" + serviceName + "/header-acls"
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
	d.Set("header_name", dataItems["header-name"])
	if val, ok := dataItems["blocked-attack-types"]; ok && val != nil {
		d.Set("blocked_attack_types", sortFileList(val.([]interface{}), ""))
	} else {
		d.Set("blocked_attack_types", nil)
	}
	d.Set("comments", dataItems["comments"])
	if val, ok := dataItems["custom-blocked-attack-types"]; ok && val != nil {
		d.Set("custom_blocked_attack_types", sortFileList(val.([]interface{}), ""))
	} else {
		d.Set("custom_blocked_attack_types", nil)
	}
	d.Set("denied_metachars", dataItems["denied-metachars"])
	if val, ok := dataItems["exception-patterns"]; ok && val != nil {
		d.Set("exception_patterns", sortFileList(val.([]interface{}), ""))
	} else {
		d.Set("exception_patterns", nil)
	}
	d.Set("max_header_value_length", fmt.Sprintf("%v", dataItems["max-header-value-length"]))
	d.Set("mode", dataItems["mode"])
	d.Set("status", dataItems["status"])

	return nil
}

func resourceCudaWAFHeaderACLUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Updating Barracuda WAF resource " + name)

	parent := d.Get("parent").([]interface{})
	serviceName := parent[0].(string)

	resourceEndpoint := "/services/" + serviceName + "/header-acls"
	err := client.UpdateBarracudaWAFResource(name, hydrateBarracudaWAFHeaderACLResource(d, "put", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF resource (%s) (%v)", name, err)
		return err
	}

	return resourceCudaWAFHeaderACLRead(d, m)
}

func resourceCudaWAFHeaderACLDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Deleting Barracuda WAF resource " + name)

	parent := d.Get("parent").([]interface{})
	serviceName := parent[0].(string)

	resourceEndpoint := "/services/" + serviceName + "/header-acls"
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

func hydrateBarracudaWAFHeaderACLResource(d *schema.ResourceData, method string, endpoint string) *APIRequest {

	//resourcePayload : payload for the resource
	resourcePayload := map[string]interface{}{
		"name":                        d.Get("name").(string),
		"header-name":                 d.Get("header_name").(string),
		"blocked-attack-types":        d.Get("blocked_attack_types"),
		"comments":                    d.Get("comments").(string),
		"custom-blocked-attack-types": d.Get("custom_blocked_attack_types"),
		"denied-metachars":            d.Get("denied_metachars").(string),
		"exception-patterns":          d.Get("exception_patterns"),
		"max-header-value-length":     d.Get("max_header_value_length").(string),
		"mode":                        d.Get("mode").(string),
		"status":                      d.Get("status").(string),
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
		if v, ok := val.(string); ok && len(v) == 0 {
			delete(resourcePayload, key)
		}
	}

	return &APIRequest{
		URL:  endpoint,
		Body: resourcePayload,
	}
}
