package user

import "github.com/open-policy-agent/opa/rego"

func TestOpa() {
	r := rego.New(
		rego.Query("x = data.example.allow"),
		rego.Load([]string{"./example.rego"}, nil))
}
