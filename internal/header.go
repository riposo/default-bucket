package internal

import (
	"net/http"
	"net/textproto"
	"strconv"
	"strings"
)

// Hop-by-hop headers. These are removed when sent to the backend.
// As of RFC 7230, hop-by-hop headers are required to appear in the
// Connection header field. These are the headers defined by the
// obsoleted RFC 2616 (section 13.5.1) and are used for backward
// compatibility.
var hopHeaders = []string{
	"Connection",
	"Proxy-Connection", // non-standard but still sent by libcurl and rejected by e.g. google
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"Te",      // canonicalized version of "TE"
	"Trailer", // not Trailers per URL above; https://www.rfc-editor.org/errata_search.php?eid=4522
	"Transfer-Encoding",
	"Upgrade",
}

func normHeaders(h http.Header, contentLen int) {
	// remove hop-by-hop headers listed in the "Connection" header, see RFC 7230, section 6.1
	for _, f := range h["Connection"] {
		for _, sf := range strings.Split(f, ",") {
			if sf = textproto.TrimString(sf); sf != "" {
				h.Del(sf)
			}
		}
	}

	// clear Hop-by-hop headers
	for _, k := range hopHeaders {
		h.Del(k)
	}

	// clear Content-Encoding as we decode the body
	h.Del("Content-Encoding")
	h.Del("Content-Length")
	if contentLen != 0 {
		h.Set("Content-Length", strconv.Itoa(contentLen))
	}
}
