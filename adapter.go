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

	if success, _, err := a.kv.CAS(p,nil); success {
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

	}
	if len(model["g"]["g"].Policy) != 0 {
		rule = append(model["g"]["g"].Policy, rule...)
	}

	err := a.writePolicyKey(rule)
	if err != nil {
		return err
	}

	return nil
}
