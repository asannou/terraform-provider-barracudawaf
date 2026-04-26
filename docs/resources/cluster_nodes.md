# barracudawaf_cluster_nodes

`barracudawaf_cluster_nodes` manages `Cluster nodes` on the Barracuda Web Application Firewall.

## Example Usage

```hcl
resource "barracudawaf_cluster_nodes" "node2" {
  ip_address = "192.168.1.11"
}
```

## Argument Reference

* `ip_address` (Required) - Cluster System IP.
