# barracudawaf_system_export_log_settings

`barracudawaf_system_export_log_settings` manages `System Export Log Settings` on the Barracuda Web Application Firewall.

## Example Usage

```hcl
resource "barracudawaf_system_export_log_settings" "example" {
  export_access_logs           = "Enable"
  export_audit_logs            = "Enable"
  export_web_firewall_logs     = "Enable"
  export_network_firewall_logs = "Enable"
  export_system_logs           = "Enable"
}
```

## Argument Reference

* `export_access_logs` (Optional) - Enable export of Access Logs. Valid values are `Enable` or `Disable`. Default is `Disable`.
* `export_audit_logs` (Optional) - Enable export of Audit Logs. Valid values are `Enable` or `Disable`. Default is `Disable`.
* `export_web_firewall_logs` (Optional) - Enable export of Web Firewall Logs. Valid values are `Enable` or `Disable`. Default is `Disable`.
* `export_network_firewall_logs` (Optional) - Enable export of Network Firewall Logs. Valid values are `Enable` or `Disable`. Default is `Disable`.
* `export_system_logs` (Optional) - Enable export of System Logs. Valid values are `Enable` or `Disable`. Default is `Disable`.
