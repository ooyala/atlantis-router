package routing

import (
	"net/http"
	"strings"
)

type HostMatcher struct {
	Host string
}

func (h *HostMatcher) Match(r *http.Request) bool {
	return r.Host == h.Host
}

func NewHostMatcher(r string) Matcher {
	return &HostMatcher{r}
}

type HeaderMatcher struct {
	Header     string
	Value      string
	ParseError bool
}

func (h *HeaderMatcher) Match(r *http.Request) bool {
	if h.ParseError {
		return false
	}
	return r.Header.Get(h.Header) == h.Value
}

func NewHeaderMatcher(r string) Matcher {
	hdrVal := strings.Split(r, ":")
	if len(hdrVal) != 2 {
		return &HeaderMatcher{"", "", true}
	}

	if hdrVal[0] == "" || hdrVal[1] == "" {
		return &HeaderMatcher{"", "", true}
	}

	return &HeaderMatcher{hdrVal[0], hdrVal[1], false}
}
