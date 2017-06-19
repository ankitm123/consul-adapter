// Copyright 2017 The casbin Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ConsulAdapter

import (
	"fmt"

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

// LoadPolicy loads policy from database.
func (a *DBAdapter) LoadPolicy(model model.Model) {
	a.init()

	for ptype, ast := range model["p"] {
		for _, rule := range ast.Policy {
			fmt.Println("Load Policy: ", ptype, "  ", rule)
		}
	}

	for ptype, ast := range model["g"] {
		for _, rule := range ast.Policy {
			fmt.Println(ptype, "  ", rule)
		}
	}
}

// SavePolicy saves policy to database.
func (a *DBAdapter) SavePolicy(model model.Model) {
	a.init()

	for ptype, ast := range model["p"] {
		for _, rule := range ast.Policy {
			fmt.Println(ptype, "  ", rule)
		}
	}

	for ptype, ast := range model["g"] {
		for _, rule := range ast.Policy {
			fmt.Println("Load Policy: ", ptype, "  ", rule)
		}
	}
}
