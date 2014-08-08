
package api

import (
	"net/http"
	"encoding/json"
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


func GetStatusJson(status string) string {

	m := map[string]interface{}
	m["Status"] = status
	b, err := json.Marshal(m)

	return string(b)	

}


