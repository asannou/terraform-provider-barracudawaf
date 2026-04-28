package barracudawaf

import (
	"fmt"
	"log"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCudaWAFJSONProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceCudaWAFJSONProfileCreate,
		Read:   resourceCudaWAFJSONProfileRead,
		Update: resourceCudaWAFJSONProfileUpdate,
		Delete: resourceCudaWAFJSONProfileDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "JSON Profile Name",
			},
			"parent": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
			"host_match": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Host Match",
			},
			"url_match": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "URL Match",
			},
			"method": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: "Methods",
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Status",
			},
			"mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Mode",
			},
			"json_policy": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "JSON Policy",
			},
			"validate_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Validate Key",
			},
			"ignore_keys": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: "Ignore Keys",
			},
			"allowed_content_types": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: "Inspect Mime Types",
			},
			"exception_patterns": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: "Exception Patterns",
			},
			"comment": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Comments",
			},
		},

		Description: "`barracudawaf_json_profile` manages `JSON Profile` on the Barracuda Web Application Firewall.",
	}
}

func resourceCudaWAFJSONProfileCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Get("name").(string)

	log.Println("[INFO] Creating Barracuda WAF resource " + name)

	parent := d.Get("parent").([]interface{})
	serviceName := parent[0].(string)

	resourceEndpoint := "/services/" + serviceName + "/json-profiles"
	err := client.CreateBarracudaWAFResource(name, hydrateBarracudaWAFJSONProfileResource(d, "post", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF resource (%s) (%v) ", name, err)
		return err
	}


	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF sub resource (%s) (%v) ", name, err)
		return err
	}

	d.SetId(name)
	return resourceCudaWAFJSONProfileRead(d, m)
}

func resourceCudaWAFJSONProfileRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()
	log.Println("[INFO] Fetching Barracuda WAF resource " + name)

	parent := d.Get("parent").([]interface{})
	serviceName := parent[0].(string)

	resourceEndpoint := "/services/" + serviceName + "/json-profiles"
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
	d.Set("host_match", dataItems["host-match"])
	d.Set("url_match", dataItems["url-match"])
	if val, ok := dataItems["method"]; ok && val != nil {
		d.Set("method", sortFileList(val.([]interface{}), ""))
	} else {
		d.Set("method", nil)
	}
	d.Set("status", dataItems["status"])
	d.Set("mode", dataItems["mode"])
	d.Set("json_policy", dataItems["json-policy"])
	d.Set("validate_key", dataItems["validate-key"])
	if val, ok := dataItems["ignore-keys"]; ok && val != nil {
		d.Set("ignore_keys", sortFileList(val.([]interface{}), ""))
	} else {
		d.Set("ignore_keys", nil)
	}
	if val, ok := dataItems["allowed-content-types"]; ok && val != nil {
		d.Set("allowed_content_types", sortFileList(val.([]interface{}), ""))
	} else {
		d.Set("allowed_content_types", nil)
	}
	if val, ok := dataItems["exception-patterns"]; ok && val != nil {
		d.Set("exception_patterns", sortFileList(val.([]interface{}), ""))
	} else {
		d.Set("exception_patterns", nil)
	}
	d.Set("comment", dataItems["comment"])
	return nil
}

func resourceCudaWAFJSONProfileUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Updating Barracuda WAF resource " + name)

	parent := d.Get("parent").([]interface{})
	serviceName := parent[0].(string)

	resourceEndpoint := "/services/" + serviceName + "/json-profiles"
	err := client.UpdateBarracudaWAFResource(name, hydrateBarracudaWAFJSONProfileResource(d, "put", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF resource (%s) (%v)", name, err)
		return err
	}


	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF sub resource (%s) (%v)", name, err)
		return err
	}

	return resourceCudaWAFJSONProfileRead(d, m)
}

func resourceCudaWAFJSONProfileDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Deleting Barracuda WAF resource " + name)

	parent := d.Get("parent").([]interface{})
	serviceName := parent[0].(string)

	resourceEndpoint := "/services/" + serviceName + "/json-profiles"
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

func hydrateBarracudaWAFJSONProfileResource(d *schema.ResourceData, method string, endpoint string) *APIRequest {

	//resourcePayload : payload for the resource
	resourcePayload := map[string]interface{}{
		"name":                  d.Get("name").(string),
		"host-match":            d.Get("host_match").(string),
		"url-match":             d.Get("url_match").(string),
		"method":                d.Get("method"),
		"status":                d.Get("status").(string),
		"mode":                  d.Get("mode").(string),
		"json-policy":           d.Get("json_policy").(string),
		"validate-key":          d.Get("validate_key").(string),
		"ignore-keys":           d.Get("ignore_keys"),
		"allowed-content-types": d.Get("allowed_content_types"),
		"exception-patterns":    d.Get("exception_patterns"),
		"comment":               d.Get("comment").(string),
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
