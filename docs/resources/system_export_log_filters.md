# barracudawaf_system_export_log_filters

`barracudawaf_system_export_log_filters` manages `System Export Log Filters` on the Barracuda Web Application Firewall.

## Example Usage

```hcl
resource "barracudawaf_system_export_log_filters" "example" {
  web_firewall_log_severity = "5-Notice"
  system_log_severity       = "5-Notice"
}
```

## Argument Reference

* `web_firewall_log_severity` (Optional) - Web Firewall Log Severity. Valid values are `0-Emergency`, `1-Alert`, `2-Critical`, `3-Error`, `4-Warning`, `5-Notice`, `6-Information`, `7-Debug`. Default is `7-Debug`.
* `system_log_severity` (Optional) - System Log Severity. Valid values are `0-Emergency`, `1-Alert`, `2-Critical`, `3-Error`, `4-Warning`, `5-Notice`, `6-Information`, `7-Debug`. Default is `7-Debug`.
