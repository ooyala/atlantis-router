
package api

import (
	"fmt"
	"net/http"
	"encoding/json"
	"atlantis/router/api/auth"
)



func GetRequestBody(r http.Request) ([]byte, err){

	cLen := r.ContentLength
	var b [cLen]byte
	n, err := r.Body.Read(b)
	if n < cLen {
		if err != nil {
			return err
		}
		
		return errors.New("Could not read full request body")
	}

	return b,nil

}


func GetMapFromReqJson(r http.Request) (map[string]interface{}, err) {

	body, err := GetRequestBody(r)
	if err != nil {
		return nil, err
	}

	var v interface{}
	err = json.Unmarshal(body, &v)
	if err != nil {
		return err
	}
	
	return v.(map[string]interface{})
	


}

GetErrorStatusJson(status string, err error){

	return GetStatusJson(fmt.Sprint(status, err))
}


func GetStatusJson(status string) string {

	m := map[string]interface{}
	m["Status"] = status
	b, err := json.Marshal(m)

	return string(b)	

}

func WriteResponse(w *http.ResponseWriter, code int, json string) {

	//TODO: check if code is valid
	w.WriteHeader(code)
	fmt.Fprintf(w, "%s", json)

}

func GetUserSecretAndAuth(r *http.Request) error{

	user := r.Header.Get("User")
	secret := r.Header.Get("Secret")

	if len(user) == 0 {
		return errors.New("Could not get user header val")
	} else if len(secret) == 0 {
		return errors.New("Could not get secret header val")
	}

	return auth.IsAllowed(user, secret)

}
