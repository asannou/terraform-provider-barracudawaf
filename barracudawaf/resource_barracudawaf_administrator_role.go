package barracudawaf

import (
	"fmt"
	"log"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCudaWAFAdministratorRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceCudaWAFAdministratorRoleCreate,
		Read:   resourceCudaWAFAdministratorRoleRead,
		Update: resourceCudaWAFAdministratorRoleUpdate,
		Delete: resourceCudaWAFAdministratorRoleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Role Name",
			},
			"api_privilege": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "API Privilege",
			},
			"authentication_services": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: "Auth Services",
			},
			"objects": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: "Object access permissions",
			},
			"operations": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: "Specify Allowed Operations",
			},
			"role_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Role Type",
			},
			"security_policies": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: "Security Policies",
			},
			"service_groups": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: "Service Group",
			},
			"services": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: "Services",
			},
			"vsites": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: "Vsites",
			},
		},

		Description: "`barracudawaf_administrator_role` manages `Administrator Role` on the Barracuda Web Application Firewall.",
	}
}

func resourceCudaWAFAdministratorRoleCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Get("name").(string)

	log.Println("[INFO] Creating Barracuda WAF resource " + name)

	resourceEndpoint := "/administrator-roles"
	err := client.CreateBarracudaWAFResource(name, hydrateBarracudaWAFAdministratorRoleResource(d, "post", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF resource (%s) (%v) ", name, err)
		return err
	}


	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF sub resource (%s) (%v) ", name, err)
		return err
	}

	d.SetId(name)
	return resourceCudaWAFAdministratorRoleRead(d, m)
}

func resourceCudaWAFAdministratorRoleRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()
	log.Println("[INFO] Fetching Barracuda WAF resource " + name)

	resourceEndpoint := "/administrator-roles"
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
	d.Set("api_privilege", dataItems["api-privilege"])
	d.Set("authentication_services", dataItems["authentication-services"])
	d.Set("objects", dataItems["objects"])
	d.Set("operations", dataItems["operations"])
	d.Set("role_type", dataItems["role-type"])
	d.Set("security_policies", dataItems["security-policies"])
	d.Set("service_groups", dataItems["service-groups"])
	d.Set("services", dataItems["services"])
	d.Set("vsites", dataItems["vsites"])
	return nil
}

func resourceCudaWAFAdministratorRoleUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Updating Barracuda WAF resource " + name)

	resourceEndpoint := "/administrator-roles"
	err := client.UpdateBarracudaWAFResource(name, hydrateBarracudaWAFAdministratorRoleResource(d, "put", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF resource (%s) (%v)", name, err)
		return err
	}


	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF sub resource (%s) (%v)", name, err)
		return err
	}

	return resourceCudaWAFAdministratorRoleRead(d, m)
}

func resourceCudaWAFAdministratorRoleDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Deleting Barracuda WAF resource " + name)

	resourceEndpoint := "/administrator-roles"
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

func hydrateBarracudaWAFAdministratorRoleResource(d *schema.ResourceData, method string, endpoint string) *APIRequest {

	//resourcePayload : payload for the resource
	resourcePayload := map[string]interface{}{
		"name":                    d.Get("name").(string),
		"api-privilege":           d.Get("api_privilege").(string),
		"authentication-services": d.Get("authentication_services"),
		"objects":                 d.Get("objects"),
		"operations":              d.Get("operations"),
		"role-type":               d.Get("role_type").(string),
		"security-policies":       d.Get("security_policies"),
		"service-groups":          d.Get("service_groups"),
		"services":                d.Get("services"),
		"vsites":                  d.Get("vsites"),
	}

	// parameters not supported for updates
	if method == "put" {
		updatePayloadExceptions := [...]string{}
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
