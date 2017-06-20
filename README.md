# consul_adapter
casbin KV store adapter for casbin

## Usage
```go
package main

import (
	"github.com/casbin/casbin"
	"github.com/casbin/mysql_adapter"
)

func main() {
	// Initialize a consul adapter and use it in a Casbin enforcer:
	a := consuladapter.NewDBAdapter()
	e := casbin.NewEnforcer("examples/rbac_model.conf", a)
	
	// Load the policy from consul KV store.
	e.LoadPolicy()
	
	// Check the permission.
	e.Enforce("alice", "data1", "read")
	
	// Modify the policy.
	// e.AddPolicy(...)
	// e.RemovePolicy(...)
	
	// Save the policy back to consul KV store.
	e.SavePolicy()
}
```
