package barracudawaf

import (
	"fmt"
	"log"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCudaWAFJSONKeyProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceCudaWAFJSONKeyProfileCreate,
		Read:   resourceCudaWAFJSONKeyProfileRead,
		Update: resourceCudaWAFJSONKeyProfileUpdate,
		Delete: resourceCudaWAFJSONKeyProfileDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Profile Name",
			},
			"parent": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
				Description: "Parent format: [service_name, json_profile_name]",
			},
			"key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Key",
			},
			"value_class": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Custom Parameter Class",
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Status",
			},
			"value_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Value Type",
			},
			"max_length": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Max Length",
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
			"max_keys": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Max Number Of Keys",
			},
			"max_array_elements": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Max Array Elements",
			},
			"allow_null": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Allow NULL",
			},
			"base64_decode_parameter_value": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Base64 Decode",
			},
			"validate_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Validate Key",
			},
			"allowed_metachars": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Allowed Metacharacters",
			},
			"exception_patterns": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: "Exception Patterns",
			},
			"comments": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Comments",
			},
		},

		Description: "`barracudawaf_json_key_profile` manages `JSON Key Profile` on the Barracuda Web Application Firewall.",
	}
}

func resourceCudaWAFJSONKeyProfileCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Get("name").(string)

	log.Println("[INFO] Creating Barracuda WAF resource " + name)

	parent := d.Get("parent").([]interface{})
	serviceName := parent[0].(string)
	jsonProfileName := parent[1].(string)

	resourceEndpoint := fmt.Sprintf("/services/%s/json-profiles/%s/json-key-profiles", serviceName, jsonProfileName)
	err := client.CreateBarracudaWAFResource(name, hydrateBarracudaWAFJSONKeyProfileResource(d, "post", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF resource (%s) (%v) ", name, err)
		return err
	}


	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF sub resource (%s) (%v) ", name, err)
		return err
	}

	d.SetId(name)
	return resourceCudaWAFJSONKeyProfileRead(d, m)
}

func resourceCudaWAFJSONKeyProfileRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()
	log.Println("[INFO] Fetching Barracuda WAF resource " + name)

	parent := d.Get("parent").([]interface{})
	serviceName := parent[0].(string)
	jsonProfileName := parent[1].(string)

	resourceEndpoint := fmt.Sprintf("/services/%s/json-profiles/%s/json-key-profiles", serviceName, jsonProfileName)
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
	d.Set("key", dataItems["key"])
	d.Set("value_class", dataItems["value-class"])
	d.Set("status", dataItems["status"])
	d.Set("value_type", dataItems["value-type"])
	d.Set("max_length", dataItems["max-length"])
	d.Set("max_number_value", dataItems["max-number-value"])
	d.Set("min_number_value", dataItems["min-number-value"])
	d.Set("max_keys", dataItems["max-keys"])
	d.Set("max_array_elements", dataItems["max-array-elements"])
	d.Set("allow_null", dataItems["allow-null"])
	d.Set("base64_decode_parameter_value", dataItems["base64-decode-parameter-value"])
	d.Set("validate_key", dataItems["validate-key"])
	d.Set("allowed_metachars", dataItems["allowed-metachars"])
	d.Set("exception_patterns", dataItems["exception-patterns"])
	d.Set("comments", dataItems["comments"])
	return nil
}

func resourceCudaWAFJSONKeyProfileUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Updating Barracuda WAF resource " + name)

	parent := d.Get("parent").([]interface{})
	serviceName := parent[0].(string)
	jsonProfileName := parent[1].(string)

	resourceEndpoint := fmt.Sprintf("/services/%s/json-profiles/%s/json-key-profiles", serviceName, jsonProfileName)
	err := client.UpdateBarracudaWAFResource(name, hydrateBarracudaWAFJSONKeyProfileResource(d, "put", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF resource (%s) (%v)", name, err)
		return err
	}


	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF sub resource (%s) (%v)", name, err)
		return err
	}

	return resourceCudaWAFJSONKeyProfileRead(d, m)
}

func resourceCudaWAFJSONKeyProfileDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Deleting Barracuda WAF resource " + name)

	parent := d.Get("parent").([]interface{})
	serviceName := parent[0].(string)
	jsonProfileName := parent[1].(string)

	resourceEndpoint := fmt.Sprintf("/services/%s/json-profiles/%s/json-key-profiles", serviceName, jsonProfileName)
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

func hydrateBarracudaWAFJSONKeyProfileResource(d *schema.ResourceData, method string, endpoint string) *APIRequest {

	//resourcePayload : payload for the resource
	resourcePayload := map[string]interface{}{
		"name":                          d.Get("name").(string),
		"key":                           d.Get("key").(string),
		"value-class":                   d.Get("value_class").(string),
		"status":                        d.Get("status").(string),
		"value-type":                    d.Get("value_type").(string),
		"max-length":                    d.Get("max_length").(string),
		"max-number-value":              d.Get("max_number_value").(string),
		"min-number-value":              d.Get("min_number_value").(string),
		"max-keys":                      d.Get("max_keys").(string),
		"max-array-elements":            d.Get("max_array_elements").(string),
		"allow-null":                    d.Get("allow_null").(string),
		"base64-decode-parameter-value": d.Get("base64_decode_parameter_value").(string),
		"validate-key":                  d.Get("validate_key").(string),
		"allowed-metachars":             d.Get("allowed_metachars").(string),
		"exception-patterns":            d.Get("exception_patterns"),
		"comments":                      d.Get("comments").(string),
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
