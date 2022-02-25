package resolvers_test

import (
	"fmt"
	"testing"

	"github.com/brice-74/golang-base-api/internal/testutils"
	"github.com/graph-gophers/graphql-go/gqltesting"
)

func TestRegisterUserAccount(t *testing.T) {
	var (
		db     = testutils.PrepareDB(t)
		schema = testutils.ParseTestSchema(db)
	)

	const (
		email      string = "test@test.com"
		password   string = "passWORD123!"
		profilName string = "name"
	)

	test := &gqltesting.Test{
		Schema: schema,
		Query: fmt.Sprintf(`
			mutation {
				registerUserAccount(input: {
					email: "%s",
					password: "%s",
					profilName: "%s"
				}) {
					email
					profilName
					roles
				}
			}`, email, password, profilName),
		ExpectedResult: fmt.Sprintf(`
			{
				"registerUserAccount": {
					"email": "%s",
					"profilName": "%s",
					"roles": [
						"ROLE_USER"
					]
				}
			}`, email, profilName),
	}

	gqltesting.RunTest(t, test)
}
