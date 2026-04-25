# barracudawaf_header_acl

`barracudawaf_header_acl` manages `Header ACL` on the Barracuda Web Application Firewall.

## Example Usage

```hcl
resource "barracudawaf_header_acl" "example" {
  parent      = [barracudawaf_services.example.name]
  name        = "example-header-acl"
  header_name = "X-Restricted-Header"
  status      = "On"
}
```

## Argument Reference

* `parent` (Required) - The parent service name.
* `name` (Required) - Header ACL Name.
* `header_name` (Required) - Header Name.
* `blocked_attack_types` (Optional) - Blocked Attack Types.
* `comments` (Optional) - Comments.
* `custom_blocked_attack_types` (Optional) - Custom Blocked Attack Types.
* `denied_metachars` (Optional) - Denied Metacharacters.
* `exception_patterns` (Optional) - Exception Patterns.
* `max_header_value_length` (Optional) - Max Header Value Length.
* `mode` (Optional) - Mode.
* `status` (Optional) - Status.
