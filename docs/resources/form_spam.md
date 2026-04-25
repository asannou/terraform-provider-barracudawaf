# barracudawaf_form_spam

`barracudawaf_form_spam` manages `Form Spam` on the Barracuda Web Application Firewall.

## Example Usage

```hcl
resource "barracudawaf_form_spam" "example" {
  parent     = [barracudawaf_services.example.name]
  name       = "example-form"
  action_url = "/login"
  status     = "On"
}
```

## Argument Reference

* `parent` (Required) - The parent service name.
* `name` (Required) - Form Name.
* `action_url` (Required) - Action URL.
* `minimum_form_fill_time` (Optional) - Minimum Form Fill Time.
* `mode` (Optional) - Mode.
* `parameter_class` (Optional) - Parameter Class.
* `parameter_name` (Optional) - Parameter Name.
* `status` (Optional) - Status.
