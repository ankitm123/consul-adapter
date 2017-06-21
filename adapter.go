package consuladapter

import (
	"errors"
	"strconv"
	"strings"

	"github.com/casbin/casbin/model"
	"github.com/hashicorp/consul/api"
	"github.com/inconshreveable/log15"
)

// KVAdapter represents the consul adapter for policy persistence, can load policy from consul or save policy to consul.
type KVAdapter struct {
	kv *api.KV
}

// NewKVAdapter is the constructor for KVAdapter.
func NewKVAdapter(kv *api.KV) *KVAdapter {
	a := KVAdapter{}
	a.kv = kv

	return &a
}

//cas - Check and set function returns true or false if the operation is successful
func (a *KVAdapter) cas(kvpair *api.KVPair) (bool, *api.WriteMeta, error) {
	return a.kv.CAS(kvpair, nil)
}

// list is used to lookup all keys under a prefix
func (a *KVAdapter) list(prefix string) (api.KVPairs, *api.QueryMeta, error) {
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
func (a *KVAdapter) LoadPolicy(model model.Model) {

	pairs, _, err := a.list("rp")
	if err != nil {
		log15.Error("Could not retreive list of key-value pairs", "Error", err)
		//return err
	}
	for _, v := range pairs {
		line := string(v.Value)
		loadPolicyLine(line, model)
	}

	if err != nil {

	}
	//return nil
}

func (a *KVAdapter) writePolicyLine(ptype string, rule []string) error {
	line := ptype

	for i := range rule {
		line += ", " + rule[i]
	}
	// for i := 0; i < 4-len(rule); i++ {
	// 	line += ","
	// }
	_, meta, err := a.list("rp")
	if err != nil {
		log15.Error("Could not retrieve key-value pair", "Error", err)
		return err
	}
	p := &api.KVPair{Key: ptype + strconv.FormatUint(meta.LastIndex, 10), Value: []byte(line)}
	if success, _, err := a.cas(p); success {
		if err != nil {
			log15.Error("Check and Set failed for Consul KV", "Error", err)
			return err
		}
	} else {
		err = errors.New("Check and set returned false for Consul KV")
		log15.Error("Check and Set returned false for Consul KV", "Error", err)
		return err
	}
	return nil
}

// SavePolicy saves policy to consul.
func (a *KVAdapter) SavePolicy(model model.Model) {

	//Loop over the policies
	for ptype, ast := range model["p"] {
		for _, rule := range ast.Policy {
			err := a.writePolicyLine(ptype, rule)
			if err != nil {
				log15.Error("Error storing policy to consul KV store", "Error", err)
				//return err
			}
		}
	}

	//Loop over group to role mapping
	for ptype, ast := range model["g"] {
		for _, rule := range ast.Policy {
			err := a.writePolicyLine(ptype, rule)
			if err != nil {
				log15.Error("Error storing policy to consul KV store", "Error", err)
				//return err
			}
		}
	}
	//return nil
}
