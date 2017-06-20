package consuladapter

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/casbin/casbin/model"
	"github.com/hashicorp/consul/api"
	"github.com/inconshreveable/log15"
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

//CAS - Check and set function returns true or false if the operation is successful
func (a *DBAdapter) CAS(kvpair *api.KVPair) (bool, *api.WriteMeta, error) {
	return a.kv.CAS(kvpair, nil)
}

// List is used to lookup all keys under a prefix
func (a *DBAdapter) List(prefix string) (api.KVPairs, *api.QueryMeta, error) {
	return a.kv.List(prefix, nil)
}

func loadPolicyLine(line string, model model.Model) {
	if line == "" {
		return
	}
	tokens := strings.Split(line, ", ")
	key := tokens[0]
	sec := key[:1]
	model[sec][key].Policy = append(model[sec][key].Policy, tokens[1:])
}

// LoadPolicy loads policy from consul.
func (a *DBAdapter) LoadPolicy(model model.Model) {
	a.init()

	pairs, meta, err := a.List("")
	for _, v := range pairs {
		line := string(v.Value)
		loadPolicyLine(line, model)
	}
	fmt.Println(meta.LastIndex)
	if err != nil {
		fmt.Println("List error API: ", err)
	}

}

func (a *DBAdapter) writePolicyLine(ptype string, rule []string) {
	line := ptype
	fmt.Println(line)
	for i := range rule {
		line += ", " + rule[i]
	}
	// for i := 0; i < 4-len(rule); i++ {
	// 	line += ","
	// }
	_, meta, err := a.List("")
	if err != nil {
		log15.Error("Could not retrieve key-value pair", "Error", err)
	}
	p := &api.KVPair{Key: ptype + strconv.FormatUint(meta.LastIndex, 10), Value: []byte(line)}
	if success, _, err := a.CAS(p); success {
		if err != nil {
			log15.Error("Check and Set failed for Consul KV", "Error", err)
		}
	} else {
		log15.Error("Check and Set returned false for Consul KV", "Error", err)
	}

}

// SavePolicy saves policy to consul.
func (a *DBAdapter) SavePolicy(model model.Model) {
	a.init()

	for ptype, ast := range model["p"] {
		for _, rule := range ast.Policy {
			a.writePolicyLine(ptype, rule)
		}
	}

	for ptype, ast := range model["g"] {
		for _, rule := range ast.Policy {
			a.writePolicyLine(ptype, rule)
		}
	}
}
