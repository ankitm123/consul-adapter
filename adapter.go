package consuladapter

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/casbin/casbin/model"
	"github.com/hashicorp/consul/api"
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

func loadPolicyKey(line string, model model.Model) {
	if line == "" {
		return
	}

	tokens := strings.Split(line, ";")
	key := tokens[0]
	sec := key[:1]
	model[sec][key].Policy = append(model[sec][key].Policy, tokens[1:])

}

// LoadPolicy loads policy from consul.
func (a *KVAdapter) LoadPolicy(model model.Model) error {
	line := [][]string{}
	//rule := ""

	//TODO: Write a get function
	pair, _, err := a.kv.Get("rp", nil)
	if err != nil {
		return err
	}
	if pair != nil {
		json.Unmarshal(pair.Value, &line)

		for _, v := range line {
			if len(v) > 2 {
				v = append([]string{"p"}, v...)
			} else {
				v = append([]string{"g"}, v...)
			}

			rule := strings.Join(v, ";")
			loadPolicyKey(rule, model)

		}

	}

	return nil
}

func (a *KVAdapter) writePolicyKey(rule [][]string) error {
	pair, _, err := a.kv.Get("rp", nil)
	if err != nil {
		return err
	}

	value, _ := json.Marshal(rule)

	p := &api.KVPair{Key: "rp", Value: []byte(value)}

	//If not set, the default value is 0, and CAS will fail
	if pair != nil {
		p.ModifyIndex = pair.ModifyIndex
	}

	if success, _, err := a.cas(p); success {
		if err != nil {
			return err
		}
	} else {
		err = errors.New("Check and set returned false for Consul KV")
		return err
	}
	return nil
}

// SavePolicy saves policy to consul.
func (a *KVAdapter) SavePolicy(model model.Model) error {

	var rule [][]string
	if len(model["p"]["p"].Policy) != 0 {

		rule = append(model["p"]["p"].Policy, rule...)
		a.writePolicyKey(rule)
	}
	if len(model["g"]["g"].Policy) != 0 {

		rule = append(model["g"]["g"].Policy, rule...)
		a.writePolicyKey(rule)
	}
	if len(rule) == 0 {
		return errors.New("Invalid policy (policy cannot be empty)")
	}

	return nil
}
