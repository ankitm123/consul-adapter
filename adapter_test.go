//Write some tests
package consuladapter

import (
	"log"
	"testing"

	"github.com/casbin/casbin"
	"github.com/casbin/casbin/util"
)

func testGetPolicy(t *testing.T, e *casbin.Enforcer, res [][]string) {
	myRes := e.GetPolicy()
	log.Print("Policy: ", myRes)

	if !util.Array2DEquals(res, myRes) {
		t.Error("Policy: ", myRes, ", supposed to be ", res)
	}
}

func TestAdapter(t *testing.T) {
	e := casbin.NewEnforcer("../rbac_model.conf", "../rbac_policy.csv")

	a, err := NewKVAdapter("127.0.0.1:8500")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	err = a.SavePolicy(e.GetModel())
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	e.ClearPolicy()
	testGetPolicy(t, e, [][]string{})

	err = a.LoadPolicy(e.GetModel())
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	testGetPolicy(t, e, [][]string{{"admin", "*", "*"}, {"user", "dashboard.html", "GET"}, {"user", "settings.html", "POST"}})

	e = casbin.NewEnforcer("../rbac_model.conf", a)
	testGetPolicy(t, e, [][]string{{"admin", "*", "*"}, {"user", "dashboard.html", "GET"}, {"user", "settings.html", "POST"}})
}

func TestBadPort(t *testing.T) {
	e := casbin.NewEnforcer("../rbac_model.conf", "../rbac_policy.csv")
	a, _ := NewKVAdapter("127.0.0.1:9800")
	err := a.SavePolicy(e.GetModel())
	if err == nil {
		t.Fatalf("err: %v", err)
	}

}
