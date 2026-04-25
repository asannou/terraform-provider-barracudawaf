# barracudawaf_bond

`barracudawaf_bond` manages `Bond` on the Barracuda Web Application Firewall.

## Example Usage

```hcl
resource "barracudawaf_bond" "example" {
  name       = "bond0"
  bond_ports = "WAN,LAN"
  mode       = "Active-Backup"
}
```

## Argument Reference

* `name` (Required) - Name.
* `bond_ports` (Required) - Ports.
* `min_link` (Optional) - Minimum Links.
* `mode` (Optional) - Bonding Mode.
* `duplexity` (Optional) - Duplexity.
* `mtu` (Optional) - MTU.
* `speed` (Optional) - Speed.
