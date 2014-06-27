package testutils

import (

	"atlantis/router/logger"
	"net/http"
	"net/http/httptest"
)



func NewTestHAProxyLogRecord(getUrl string) (*logger.HAProxyLogRecord, *httptest.ResponseRecorder) {	
	r, _ := http.NewRequest("GET", getUrl, nil)
	rr   := httptest.NewRecorder()
	return &logger.HAProxyLogRecord{
		ResponseWriter:		rr,
		Request:		r,
	} , rr
}

