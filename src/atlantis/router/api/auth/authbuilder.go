package auth


func GetAuthorizer() Authorizer {

	return GetSqlAuthorizer() 
}
