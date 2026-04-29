# barracudawaf_system

`barracudawaf_system` manages `System` settings on the Barracuda Web Application Firewall.

## Example Usage

```hcl
resource "barracudawaf_system" "test" {
  domain   = "example.com"
  hostname = "waf"
  locale   = "English"
}
```

## Argument Reference

* `domain` - (Required) Default Domain.
* `locale` - (Optional) Default Language and Encoding. Valid values are `한국어`, `Dansk`, `日本語`, `Italiano`, `Español`, `Deutsch`, `Magyar`, `Français`, `Русский`, `English`, `Česky`, `繁體中文`, `Polski`, `Catalan`, `íslenska`, `简体中文`, `Português (BR)`, `Nederlands`.
* `interface_for_system_services` - (Optional) Interface for System Services. Valid values are `WAN`, `Management`.
* `time_zone` - (Optional) Time Zone.
* `hostname` - (Optional) Default Hostname.
* `enable_ipv6` - (Optional) Enable IPv6. Valid values are `Yes`, `No`.
* `device_name` - (Optional) Firewall Name.
* `operation_mode` - (Optional) Mode of Operation. Valid values are `Bridge All Traffic`, `Proxy`.
