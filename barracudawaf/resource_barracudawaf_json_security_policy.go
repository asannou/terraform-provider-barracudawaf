package barracudawaf

import (
	"fmt"
	"log"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCudaWAFJSONSecurityPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceCudaWAFJSONSecurityPolicyCreate,
		Read:   resourceCudaWAFJSONSecurityPolicyRead,
		Update: resourceCudaWAFJSONSecurityPolicyUpdate,
		Delete: resourceCudaWAFJSONSecurityPolicyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Policy Name",
			},
			"max_keys": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Max Keys",
			},
			"max_key_length": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Max Key Length",
			},
			"max_object_depth": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Max Object Depth",
			},
			"max_value_length": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Max Value Length",
			},
			"max_siblings": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Max Siblings",
			},
			"max_array_elements": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Max Array Elements",
			},
			"max_number_value": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Max Number Value",
			},
			"min_number_value": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Min Number Value",
			},
		},

		Description: "`barracudawaf_json_security_policy` manages `JSON Security Policy` on the Barracuda Web Application Firewall.",
	}
}

func resourceCudaWAFJSONSecurityPolicyCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Get("name").(string)

	log.Println("[INFO] Creating Barracuda WAF resource " + name)

	resourceEndpoint := "/json-security-policies"
	err := client.CreateBarracudaWAFResource(name, hydrateBarracudaWAFJSONSecurityPolicyResource(d, "post", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF resource (%s) (%v) ", name, err)
		return err
	}


	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF sub resource (%s) (%v) ", name, err)
		return err
	}

	d.SetId(name)
	return resourceCudaWAFJSONSecurityPolicyRead(d, m)
}

func resourceCudaWAFJSONSecurityPolicyRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()
	log.Println("[INFO] Fetching Barracuda WAF resource " + name)

	resourceEndpoint := "/json-security-policies"
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
		log.Printf("[WARN] Barracuda WAF resource (%s) not found on the system", name)
		d.SetId("")
		return nil
	}

	d.Set("name", name)
	d.Set("max_keys", dataItems["max-keys"])
	d.Set("max_key_length", dataItems["max-key-length"])
	d.Set("max_object_depth", dataItems["max-object-depth"])
	d.Set("max_value_length", dataItems["max-value-length"])
	d.Set("max_siblings", dataItems["max-siblings"])
	d.Set("max_array_elements", dataItems["max-array-elements"])
	d.Set("max_number_value", dataItems["max-number-value"])
	d.Set("min_number_value", dataItems["min-number-value"])
	return nil
}

func resourceCudaWAFJSONSecurityPolicyUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Updating Barracuda WAF resource " + name)

	resourceEndpoint := "/json-security-policies"
	err := client.UpdateBarracudaWAFResource(name, hydrateBarracudaWAFJSONSecurityPolicyResource(d, "put", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF resource (%s) (%v)", name, err)
		return err
	}


	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF sub resource (%s) (%v)", name, err)
		return err
	}

	return resourceCudaWAFJSONSecurityPolicyRead(d, m)
}

func resourceCudaWAFJSONSecurityPolicyDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Deleting Barracuda WAF resource " + name)

	resourceEndpoint := "/json-security-policies"
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

func hydrateBarracudaWAFJSONSecurityPolicyResource(d *schema.ResourceData, method string, endpoint string) *APIRequest {

	//resourcePayload : payload for the resource
	resourcePayload := map[string]interface{}{
		"name":               d.Get("name").(string),
		"max-keys":           d.Get("max_keys").(string),
		"max-key-length":     d.Get("max_key_length").(string),
		"max-object-depth":   d.Get("max_object_depth").(string),
		"max-value-length":   d.Get("max_value_length").(string),
		"max-siblings":       d.Get("max_siblings").(string),
		"max-array-elements": d.Get("max_array_elements").(string),
		"max-number-value":   d.Get("max_number_value").(string),
		"min-number-value":   d.Get("min_number_value").(string),
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
		if reflect.ValueOf(val).Kind() == reflect.String && reflect.ValueOf(val).Len() == 0 {
			delete(resourcePayload, key)
		}
	}

	return &APIRequest{
		URL:  endpoint,
		Body: resourcePayload,
	}
}
