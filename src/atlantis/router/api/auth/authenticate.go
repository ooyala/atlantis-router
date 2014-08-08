package auth


type Authorizer interface {
	SimpleAuth(user, secret string) error 
	IsSuperUser(user, secret string) error 
}


func IsAllowed(user, secret string)  error {
	return GetAuthorizer().SimpleAuth(user, secret)
}
