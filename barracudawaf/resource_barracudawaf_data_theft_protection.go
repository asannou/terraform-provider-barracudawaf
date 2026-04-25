package barracudawaf

import (
	"fmt"
	"log"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCudaWAFDataTheftProtection() *schema.Resource {
	return &schema.Resource{
		Create: resourceCudaWAFDataTheftProtectionCreate,
		Read:   resourceCudaWAFDataTheftProtectionRead,
		Update: resourceCudaWAFDataTheftProtectionUpdate,
		Delete: resourceCudaWAFDataTheftProtectionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Data Theft Element Name",
			},
			"parent": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
				Description: "Parent format: [policy_name]",
			},
			"enable": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Enabled",
			},
			"identity_theft_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Identity Theft Type",
			},
			"custom_identity_theft_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Custom Identity Theft Type",
			},
			"action": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Action",
			},
			"initial_characters_to_keep": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Initial Characters to Keep",
			},
			"trailing_characters_to_keep": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Trailing Characters to Keep",
			},
		},

		Description: "`barracudawaf_data_theft_protection` manages `Data Theft Protection` on the Barracuda Web Application Firewall.",
	}
}

func resourceCudaWAFDataTheftProtectionCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Get("name").(string)

	log.Println("[INFO] Creating Barracuda WAF resource " + name)

	parent := d.Get("parent").([]interface{})
	policyName := parent[0].(string)

	resourceEndpoint := fmt.Sprintf("/security-policies/%s/protected-data-types", policyName)
	err := client.CreateBarracudaWAFResource(name, hydrateBarracudaWAFDataTheftProtectionResource(d, "post", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF resource (%s) (%v) ", name, err)
		return err
	}


	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF sub resource (%s) (%v) ", name, err)
		return err
	}

	d.SetId(name)
	return resourceCudaWAFDataTheftProtectionRead(d, m)
}

func resourceCudaWAFDataTheftProtectionRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()
	log.Println("[INFO] Fetching Barracuda WAF resource " + name)

	parent := d.Get("parent").([]interface{})
	policyName := parent[0].(string)

	resourceEndpoint := fmt.Sprintf("/security-policies/%s/protected-data-types", policyName)
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
	d.Set("enable", dataItems["enable"])
	d.Set("identity_theft_type", dataItems["identity-theft-type"])
	d.Set("custom_identity_theft_type", dataItems["custom-identity-theft-type"])
	d.Set("action", dataItems["action"])
	d.Set("initial_characters_to_keep", dataItems["initial-characters-to-keep"])
	d.Set("trailing_characters_to_keep", dataItems["trailing-characters-to-keep"])
	return nil
}

func resourceCudaWAFDataTheftProtectionUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Updating Barracuda WAF resource " + name)

	parent := d.Get("parent").([]interface{})
	policyName := parent[0].(string)

	resourceEndpoint := fmt.Sprintf("/security-policies/%s/protected-data-types", policyName)
	err := client.UpdateBarracudaWAFResource(name, hydrateBarracudaWAFDataTheftProtectionResource(d, "put", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF resource (%s) (%v)", name, err)
		return err
	}


	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF sub resource (%s) (%v)", name, err)
		return err
	}

	return resourceCudaWAFDataTheftProtectionRead(d, m)
}

func resourceCudaWAFDataTheftProtectionDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Deleting Barracuda WAF resource " + name)

	parent := d.Get("parent").([]interface{})
	policyName := parent[0].(string)

	resourceEndpoint := fmt.Sprintf("/security-policies/%s/protected-data-types", policyName)
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

func hydrateBarracudaWAFDataTheftProtectionResource(d *schema.ResourceData, method string, endpoint string) *APIRequest {

	//resourcePayload : payload for the resource
	resourcePayload := map[string]interface{}{
		"name":                        d.Get("name").(string),
		"enable":                      d.Get("enable").(string),
		"identity-theft-type":        d.Get("identity_theft_type").(string),
		"custom-identity-theft-type": d.Get("custom_identity_theft_type").(string),
		"action":                      d.Get("action").(string),
		"initial-characters-to-keep": d.Get("initial_characters_to_keep").(string),
		"trailing-characters-to-keep": d.Get("trailing_characters_to_keep").(string),
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
