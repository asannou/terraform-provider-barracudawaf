# barracudawaf_http_response_rewrite

`barracudawaf_http_response_rewrite` manages `HTTP Response Rewrite` on the Barracuda Web Application Firewall.

## Example Usage

```hcl
resource "barracudawaf_http_response_rewrite" "example" {
  parent          = [barracudawaf_services.example.name]
  name            = "example-rule"
  action          = "Insert Header"
  header          = "X-Example-Response-Header"
  rewrite_value   = "ExampleValue"
  sequence_number = "100"
}
```

## Argument Reference

* `parent` (Required) - The parent service name.
* `name` (Required) - Rule Name.
* `action` (Required) - Action.
* `comments` (Optional) - Comments.
* `condition` (Optional) - Rewrite Condition.
* `continue_processing` (Optional) - Continue Processing.
* `header` (Optional) - Header Name.
* `old_value` (Optional) - Old Value.
* `rewrite_value` (Optional) - Rewrite Value.
* `sequence_number` (Required) - Sequence Number.
