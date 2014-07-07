package logger

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	HAProxyFmtStr = "haproxy[%d]: %s:%s [%s] %s %s %s/%s/%s/%s/%s %d %d %s %s %s %d/%d/%d/%d/%d %d/%d {%s} {%s} \"%s\""
	BadGatewayMsg = "Bad Gateway"
	GatewayTimeoutMsg = "Gateway Timeout"
	ServiceUnavailableMsg = "Service Unavailable"
	
)

var copier = NewCopier()

type HAProxyLogRecord struct {
	http.ResponseWriter
	*http.Request
	pid                                       int
	clientIp, clientPort                      string
	acceptDate                                time.Time
	frontendPort                              string
	backendName                               string
	serverName                                string
	tq, tw, tc, tr, tt                        string
	statusCode                                int
	bytesRead                                 int
	capturedReqCookie                         string
	capturedResCookie                         string
	terminationState                          string
	actConn, feConn, beConn, srvConn, retries uint32
	srvQueue, backendQueue                    uint64 
	capturedRequestHeaders                    string
	capturedResponseHeaders                   string
	httpRequest                               string
}

func NewHAProxyLogRecord(wr http.ResponseWriter, r *http.Request, frontendPort string, feConn uint32, acceptDate time.Time) HAProxyLogRecord {
	var headStr, fullReq string
	for key, value := range r.Header {
		headStr += "| " + key + ":" + value[0] + " "
	}
	fullReq = r.Method + " " + r.RequestURI + " " + r.Proto
	colon := strings.LastIndex(r.RemoteAddr, ":")
	return HAProxyLogRecord{
		ResponseWriter:         wr,
		Request:                r,
		pid:                    100,
		clientIp:               r.RemoteAddr[:colon],
		clientPort:             r.RemoteAddr[colon+1:],
		acceptDate:             acceptDate,
		frontendPort:           frontendPort,
		feConn:                 feConn,
		capturedRequestHeaders: "{" + headStr + "}",
		httpRequest:            fullReq,
	}
}

func (r *HAProxyLogRecord) Log(out io.Writer) {
	var resHeadStr string
	for key, value := range r.ResponseWriter.Header() {
		resHeadStr += " " + key + ":" + value[0] + " |"
	}
	resHeadStr = "{" + resHeadStr + "}"
	timeFormatted := r.acceptDate.Format("02/Jan/2006:03:04:05.555")
	fmt.Fprintf(out, HAProxyFmtStr, r.pid, r.clientIp, r.clientPort,
		timeFormatted, r.frontendPort, r.backendName, r.serverName,
		r.tq, r.tw, r.tc, r.tr, r.tt, r.statusCode, r.bytesRead, r.capturedReqCookie,
		r.capturedResCookie, r.terminationState, r.actConn, r.feConn, r.beConn,
		r.srvConn, r.retries, r.srvQueue, r.backendQueue, r.capturedRequestHeaders,
		resHeadStr, r.httpRequest)

}

func (r *HAProxyLogRecord) Write(p []byte) (int, error) {
	written, err := r.ResponseWriter.Write(p)
	//r.responseBytes += int64(written)
	return written, err
}

func (r *HAProxyLogRecord) PoolUpdateRecord(name string, count uint32, beQueue uint64) {
	r.backendName = name
	r.beConn = count
	r.backendQueue = beQueue
}
func (r *HAProxyLogRecord) ServerUpdateRecord(name string, sQueue uint64, sConn uint32) {
	r.serverName = name
	r.srvQueue = sQueue
	r.srvConn = sConn
}
func (r *HAProxyLogRecord) AddResponseHeaderMap(hdrs http.Header) {
	for hdr, vals := range hdrs {
		for _, val := range vals {
			r.AddResponseHeader(hdr, val)
		}
	}
}
func (r *HAProxyLogRecord) AddResponseHeader(hdr, val string) {
	r.ResponseWriter.Header().Add(hdr, val)
}
func (r *HAProxyLogRecord) Copy(src io.Reader) (err error){
	_, err = copier.Copy(r.ResponseWriter, src)
	return err
		
}
func (r *HAProxyLogRecord) SetResponseStatusCode(code int){
	if code >= 100 && code <= 505{ 
		r.statusCode = code
		r.ResponseWriter.WriteHeader(code)
 	}	
	

}
func (r *HAProxyLogRecord) GetResponseStatusCode() int{
	//if status code has not been set
	if r.statusCode < 100 || r.statusCode > 505 {
		r.statusCode = http.StatusOK
		r.SetResponseStatusCode(http.StatusOK)
	}
	return r.statusCode
}
//Responds to the specified request with the provided error. 
func (r *HAProxyLogRecord) Error(error string, code int){	
	r.statusCode = code
	http.Error(r.ResponseWriter, error, code)
}
func (r *HAProxyLogRecord) GetResponseHeaders() http.Header {
	return r.ResponseWriter.Header()
}


