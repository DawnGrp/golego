package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"testing"

	"github.com/open-policy-agent/opa/rego"
)

func TestOpa(t *testing.T) {

	ctx := context.Background()

	// Construct a Rego object that can be prepared or evaluated.
	r := rego.New(
		rego.Query("data.auth"), //data is rego root name, auth is the package name of rego file
		rego.Load([]string{
			"/Users/zeta/workspace/golego/modules/auth/opa_test.rego",
			"/Users/zeta/workspace/golego/modules/auth/opa_test.json",
		},
			nil))
	// Create a prepared query that can be evaluated.
	query, err := r.PrepareForEval(ctx)
	if err != nil {
		log.Fatal(err)
	}

	input := map[string]interface{}{
		"age": 20,
	}

	// Execute the prepared query.
	rs, err := query.Eval(ctx, rego.EvalInput(input))
	if err != nil {
		log.Fatal("err", err)
	}

	// Do something with the result.
	rsJson, _ := json.Marshal(rs)
	rsJson2, _ := json.Marshal(rs[0].Expressions[0].Value)

	fmt.Println("rs", string(rsJson))
	fmt.Println("value", string(rsJson2))

}

func TestGenRego(t *testing.T) {

}
