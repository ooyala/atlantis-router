package logger

import (
	"io"
	"log"
	"log/syslog"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	HAProxyFmtStr         = "haproxy[%d]: %s:%s [%s] %s %s/%s %d/%d/%d/%d/%d %d %d %s %s %s %d/%d/%d/%d/%d %d/%d {%s} {%s} \"%s\"\n"
	BadGatewayMsg         = "Bad Gateway"
	GatewayTimeoutMsg     = "Gateway Timeout"
	ServiceUnavailableMsg = "Service Unavailable"
)

var copier = NewCopier()

type HAProxyLogRecord struct {
	http.ResponseWriter
	*http.Request
	out                                       io.Writer
	pid                                       int
	clientIp, clientPort                      string
	acceptDate                                time.Time
	enterPoolTime                             time.Time
	enterServerTime                           time.Time
	serverResTime                             time.Time
	frontendPort                              string
	backendName                               string
	serverName                                string
	tq, tw, tc, tr, tt                        int64
	statusCode                                int
	bytesRead                                 int64
	capturedReqCookie                         string
	capturedResCookie                         string
	terminationState                          string
	actConn, feConn, beConn, srvConn, retries uint32
	srvQueue, backendQueue                    uint64
	capturedRequestHeaders                    string
	capturedResponseHeaders                   string
	httpRequest                               string
	sLog                                      *log.Logger
}

func NewShallowHAProxyLogRecord(out io.Writer, w http.ResponseWriter, r *http.Request) *HAProxyLogRecord {
	tlog, err := syslog.NewLogger(syslog.LOG_LOCAL5|syslog.LOG_INFO, 0)
	//if cannot connect to syslog just get basic logger to stdout
	if err != nil {
		tlog = log.New(os.Stdout, "", 0)
	}
	return &HAProxyLogRecord{
		ResponseWriter: w,
		Request:        r,
		sLog:           tlog,
	}
}

func NewHAProxyLogRecord(w http.ResponseWriter, r *http.Request, frontendPort string, feConn uint32, acceptDate time.Time) HAProxyLogRecord {
	var headStr, fullReq string
	for key, value := range r.Header {
		headStr += " " + key + ":" + value[0] + " |"
	}
	sz := len(headStr)
	if sz > 0 && headStr[sz-1] == '|' {
		headStr = headStr[:sz-1]
	}
	fullReq = r.Method + " " + r.RequestURI + " " + r.Proto
	colon := strings.LastIndex(r.RemoteAddr, ":")

	tlog, err := syslog.NewLogger(syslog.LOG_LOCAL5|syslog.LOG_INFO, 0)
	//if cannot connect to syslog just get basic logger to stdout
	if err != nil {
		tlog = log.New(os.Stdout, "", 0)
	}

	return HAProxyLogRecord{
		ResponseWriter:         w,
		Request:                r,
		pid:                    os.Getpid(),
		clientIp:               r.RemoteAddr[:colon],
		clientPort:             r.RemoteAddr[colon+1:],
		acceptDate:             acceptDate,
		frontendPort:           frontendPort,
		feConn:                 feConn,
		capturedRequestHeaders: headStr,
		httpRequest:            fullReq,
		actConn:                0,
		tq:                     0,
		tw:                     0,
		tc:                     0,
		tr:                     0,
		tt:                     0,
		backendName:            "-",
		serverName:             "-",
		capturedReqCookie:      "-",
		capturedResCookie:      "-",
		terminationState:       "--",
		sLog:                   tlog,
	}
}

func (r *HAProxyLogRecord) Log() {

	//build response header string
	var resHeadStr string
	for key, value := range r.ResponseWriter.Header() {
		resHeadStr += " " + key + ":" + value[0] + " |"
	}
	sz := len(resHeadStr)
	if sz > 0 && resHeadStr[sz-1] == '|' {
		resHeadStr = resHeadStr[:sz-1]
	}

	//build cookie strings
	r.capturedReqCookie = getCookiesString(r.Request.Cookies())

	timeFormatted := r.acceptDate.Format("02/Jan/2006:03:04:05.555")

	//calculate the queue/wait times
	r.tt = int64((r.serverResTime.UnixNano() - r.acceptDate.UnixNano()) / int64(time.Millisecond))      // total time from accepted to final response
	r.tw = int64((r.enterServerTime.UnixNano() - r.enterPoolTime.UnixNano()) / int64(time.Millisecond)) //total time spent waiting in queues

	r.sLog.Printf(HAProxyFmtStr, r.pid, r.clientIp, r.clientPort,
		timeFormatted, r.frontendPort, r.backendName, r.serverName,
		r.tq, r.tw, r.tc, r.tr, r.tt, r.statusCode, r.bytesRead, r.capturedReqCookie,
		r.capturedResCookie, r.terminationState, r.actConn, r.feConn, r.beConn,
		r.srvConn, r.retries, r.srvQueue, r.backendQueue, r.capturedRequestHeaders,
		resHeadStr, r.httpRequest)
}

func getCookiesString(cookies []*http.Cookie) string {
	var cookieStr string
	for _, c := range cookies {
		cookieStr += " " + c.Name + ":" + c.Value + " |"
	}
	sz := len(cookieStr)
	if sz > 0 && cookieStr[sz-1] == '|' {
		cookieStr = cookieStr[:sz-1]
		return "{ " + cookieStr + "}"
	}
	return "-"
}

func (r *HAProxyLogRecord) Write(p []byte) (int, error) {
	written, err := r.ResponseWriter.Write(p)
	r.bytesRead += int64(written)
	return written, err
}

func (r *HAProxyLogRecord) PoolUpdateRecord(name string, count uint32, beQueue uint64, pTime time.Time) {
	r.backendName = name
	r.beConn = count
	r.backendQueue = beQueue
	r.enterPoolTime = pTime
}
func (r *HAProxyLogRecord) ServerUpdateRecord(name string, sQueue uint64, sConn uint32, sTime time.Time) {
	r.serverName = name
	r.srvQueue = sQueue
	r.srvConn = sConn
	r.enterServerTime = sTime
}
func (r *HAProxyLogRecord) CopyHeaders(hdrs http.Header) {
	for hdr, vals := range hdrs {
		for _, val := range vals {
			r.AddResponseHeader(hdr, val)
		}
	}
}
func (r *HAProxyLogRecord) AddResponseHeader(hdr, val string) {
	r.bytesRead += int64(len(hdr))
	r.bytesRead += int64(len(val))
	r.ResponseWriter.Header().Add(hdr, val)
}
func (r *HAProxyLogRecord) Copy(src io.Reader) (err error) {
	bwritten, err := copier.Copy(r.ResponseWriter, src)
	r.bytesRead += bwritten
	return err

}
func (r *HAProxyLogRecord) WriteHeader(code int) {
	if code >= 100 && code <= 599 {
		r.statusCode = code
		r.ResponseWriter.WriteHeader(code)
	}

}
func (r *HAProxyLogRecord) GetResponseStatusCode() int {
	//if status code has not been set
	if r.statusCode < 100 || r.statusCode > 599 {
		r.statusCode = http.StatusOK
		r.WriteHeader(http.StatusOK)
	}
	return r.statusCode
}

//Responds to the specified request with the provided error.
func (r *HAProxyLogRecord) Error(error string, code int) {
	r.statusCode = code
	http.Error(r.ResponseWriter, error, code)
}
func (r *HAProxyLogRecord) GetResponseHeaders() http.Header {
	return r.ResponseWriter.Header()
}

func (r *HAProxyLogRecord) UpdateTr(resStartTime, resRetTime time.Time) {
	r.tr = int64((resRetTime.UnixNano() - resStartTime.UnixNano()) / int64(time.Millisecond))
	r.serverResTime = resRetTime
}

//Set's the termination state and log's the request
func (r *HAProxyLogRecord) Terminate(termState string) {
	r.terminationState = termState
	r.Log()
}
