package ConsulAdapter

import (
	"fmt"
	"strconv"

	"github.com/inconshreveable/log15"

	"github.com/casbin/casbin/model"
	"github.com/hashicorp/consul/api"
)

// DBAdapter represents the database adapter for policy persistence, can load policy from database or save policy to database.
// For now, only MySQL is tested, but it should work for other RDBMS.
type DBAdapter struct {
	kv *api.KV
}

// NewDBAdapter is the constructor for DBAdapter.
func NewDBAdapter() *DBAdapter {
	a := DBAdapter{}

	return &a
}

func (a *DBAdapter) init() {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		log15.Error("Could not return consul client", "Error", err)
	}
	a.kv = client.KV()

}

//Put function creates a key-value pair in consul
func (a *DBAdapter) Put(p *api.KVPair) (*api.WriteMeta, error) {
	return a.kv.Put(p, nil)
}

//Get Function gets the key value from consul
func (a *DBAdapter) Get(key string) (*api.KVPair, *api.QueryMeta, error) {
	return a.kv.Get(key, nil)
}

// func loadPolicyLine(line string, model model.Model) {
// 	if line == "" {
// 		return
// 	}

// 	tokens := strings.Split(line, ", ")

// 	key := tokens[0]
// 	fmt.Println("key: ", key)
// 	sec := key[:1]
// 	fmt.Println("sec: ", sec)
// 	model[sec][key].Policy = append(model[sec][key].Policy, tokens[1:])
// 	fmt.Println(model)
// }

// LoadPolicy loads policy from database.
func (a *DBAdapter) LoadPolicy(model model.Model) {
	a.init()

	for ptype, ast := range model["p"] {
		for _, rule := range ast.Policy {
			fmt.Println(ptype, " ", rule)
		}
	}

	for ptype, ast := range model["g"] {
		for _, rule := range ast.Policy {
			fmt.Println(ptype, "  ", rule)
		}
	}
}

func (a *DBAdapter) writePolicyLine(ptype string, rule []string, idx int) {
	line := "'" + ptype + "'"
	fmt.Println(line)
	for i := range rule {
		line += ",'" + rule[i] + "'"
	}
	for i := 0; i < 4-len(rule); i++ {
		line += ",''"
	}
	fmt.Println(idx)
	p := &api.KVPair{Key: ptype + strconv.Itoa(idx), Value: []byte(line)}
	success, _, err := a.kv.CAS(p, nil)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("CAS operation: ", success)

}

// SavePolicy saves policy to database.
func (a *DBAdapter) SavePolicy(model model.Model) {
	a.init()

	for ptype, ast := range model["p"] {
		for idx, rule := range ast.Policy {
			a.writePolicyLine(ptype, rule, idx)
		}
	}

	for ptype, ast := range model["g"] {
		for idx, rule := range ast.Policy {
			a.writePolicyLine(ptype, rule, idx)
		}
	}
}
