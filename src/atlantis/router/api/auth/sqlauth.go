
package auth



type SqlAuth struct {
	Name string
}

func (*SqlAuth) SimpleAuth(user, secret string) error{

	return nil

}

func (*SqlAuth) IsSuperUser(user, secret string) error {

	return nil
}

func GetSqlAuthorizer() *SqlAuth {

	return &SqlAuth{ Name : "Default" }
}
