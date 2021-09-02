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
		rego.Query("data.user"),
		rego.Load([]string{"/Users/zeta/workspace/golego/modules/user/example.rego"}, nil))
	// Create a prepared query that can be evaluated.
	query, err := r.PrepareForEval(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Load the input document from stdin.
	var input interface{}

	inputString := `{
			"user": "bob",
			"age": 19,
			"method": "post",
			"path":"doc"
		}`

	err = json.Unmarshal([]byte(inputString), &input)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(inputString)
	// Execute the prepared query.
	rs, err := query.Eval(ctx, rego.EvalInput(input))
	if err != nil {
		log.Fatal("err", err)
	}

	// Do something with the result.
	rsJson, _ := json.Marshal(rs)

	fmt.Println("rs", string(rsJson))

}
