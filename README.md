# consul_adapter
casbin KV store adapter for casbin

## Usage
```go
package main

import (
	"github.com/casbin/casbin"
	"github.com/ankitm123/consul_adapter"
)

func main() {
	//This is how it should ideally work:
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

    //This is how it works now. We have a csv file which is read by the adpater using savepolicy
    e := casbin.NewEnforcer("./rbac/rbac_model.conf", "./rbac/rbac_policy.csv")
	//a := ConsulAdapter.NewDBAdapter()
	a.SavePolicy(e.GetModel())
	//a.LoadPolicy(e.GetModel())
	e = casbin.NewEnforcer("./rbac/rbac_model.conf", a)

	e.LoadPolicy()
}
```
