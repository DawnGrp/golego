package user

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
		rego.Query("data.user"), //data is rego root name, user is the package name of rego file
		rego.Load([]string{
			"/Users/zeta/workspace/golego/modules/user/opa_test.rego",
			"/Users/zeta/workspace/golego/modules/user/opa_test.json",
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

	fmt.Println("rs", string(rsJson))

}

func TestGenRego(t *testing.T) {

}
