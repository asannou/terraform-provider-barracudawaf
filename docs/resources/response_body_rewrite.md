# barracudawaf_response_body_rewrite

`barracudawaf_response_body_rewrite` manages `Response Body Rewrite` on the Barracuda Web Application Firewall.

## Example Usage

```hcl
resource "barracudawaf_response_body_rewrite" "example" {
  parent          = [barracudawaf_services.example.name]
  name            = "example-rule"
  host            = "*"
  url             = "/*"
  search_string   = "foo"
  replace_string  = "bar"
  sequence_number = "100"
}
```

## Argument Reference

* `parent` (Required) - The parent service name.
* `name` (Required) - Rule Name.
* `comments` (Optional) - Comments.
* `host` (Required) - Host Match.
* `replace_string` (Optional) - Replace String.
* `search_string` (Required) - Search String.
* `sequence_number` (Required) - Sequence number.
* `url` (Required) - URL Match.
