# barracudawaf_ntp_server

`barracudawaf_ntp_server` manages `NTP Server` on the Barracuda Web Application Firewall.

## Example Usage

```hcl
resource "barracudawaf_ntp_server" "example" {
  name       = "ntp-1"
  ip_address = "pool.ntp.org"
}
```

## Argument Reference

* `name` (Required) - System NTP Server Name.
* `ip_address` (Required) - System NTP Server IP/Address.
* `description` (Optional) - Description.
