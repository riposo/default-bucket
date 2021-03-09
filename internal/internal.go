package internal

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/riposo/riposo/pkg/api"
	"github.com/riposo/riposo/pkg/riposo"
	"github.com/riposo/riposo/pkg/schema"
	"github.com/valyala/bytebufferpool"
)

// Config supports custom configuration.
type Config struct {
	DefaultBucket struct {
		Secret HashSecret
	} `split_words:"true"`
}

// Mount mounts the config onto the routes.
func (c *Config) Mount(rts *api.Routes) {
	h := http.HandlerFunc(c.serveHTTP)
	rts.Handle("/buckets/default", h)
	rts.Handle("/buckets/default/*", h)
}

func (c *Config) serveHTTP(w http.ResponseWriter, r *http.Request) {
	req := newRequest(r)
	if !req.txn.User.IsAuthenticated() {
		api.Render(w, schema.Forbidden)
		return
	}

	bucketID := hex.EncodeToString(c.DefaultBucket.Secret.Encode(req.txn.User.ID)[:16])
	if done, err := createBucket(w, req, bucketID); err != nil {
		api.Render(w, err)
		return
	} else if done {
		return
	}

	if done, err := createCollection(w, req, bucketID); err != nil {
		api.Render(w, err)
	} else if done {
		return
	}

	var body io.ReadCloser
	var contentLen int
	if req.Method == http.MethodPost || req.Method == http.MethodPut || req.Method == http.MethodPatch {
		// parse payload
		var payload schema.Resource
		if err := api.Parse(req.Request, &payload); err != nil {
			api.Render(w, err)
			return
		}

		// rewrite ID
		if payload.Data != nil && payload.Data.ID == "default" {
			payload.Data.ID = bucketID
		}

		// re-encode payload
		buf := bytebufferpool.Get()
		defer bytebufferpool.Put(buf)

		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(true)
		if err := enc.Encode(payload); err != nil {
			api.Render(w, err)
			return
		}

		body = io.NopCloser(bytes.NewReader(buf.Bytes()))
		contentLen = buf.Len()
	}

	r2 := req.Clone(req.Context())
	r2.URL.Path = strings.Replace(req.URL.Path, "/buckets/default", "/buckets/"+bucketID, 1)
	r2.Body = body
	normHeaders(r2.Header, contentLen)

	req.mux.ServeHTTP(w, r2)
}

func createBucket(parent http.ResponseWriter, req *request, bucketID string) (bool, error) {
	// skip if bucket-create request is in progress
	if strings.HasSuffix(req.URL.Path, "/buckets/default") && req.Method == http.MethodPut {
		return false, nil
	}

	return createResource(parent, req, riposo.Path("/buckets/"+bucketID))
}

func createCollection(parent http.ResponseWriter, req *request, bucketID string) (bool, error) {
	// determine relevant collection path
	var relevant riposo.Path
	req.path.Traverse(func(sub riposo.Path) {
		if !sub.IsNode() && sub.ResourceName() == "collection" {
			relevant = sub
		}
	})

	// skip if request doesn't involve a collection
	if relevant == "" {
		return false, nil
	}

	// skip if collection-create request is in progress
	if req.path == relevant && req.Method == http.MethodPut {
		return false, nil
	}

	return createResource(parent, req, relevant)
}

func createResource(parent http.ResponseWriter, req *request, path riposo.Path) (bool, error) {
	// extract objID and resKey
	objID := path.ObjectID()
	resKey := path.ResourceName() + ".created"

	// skip if resource was already created (e.g. as part of a batch request)
	var created []string
	if val, ok := req.txn.Data[resKey]; ok {
		created = val.([]string)
	}
	if containsString(created, objID) {
		return false, nil
	}

	// setup sub-request with clean context
	ctx := context.WithValue(req.Context(), chi.RouteCtxKey, nil)
	r, err := http.NewRequestWithContext(ctx, http.MethodPut, path.String(), strings.NewReader(`{}`))
	if err != nil {
		return false, err
	}

	// perform request, delegate failures
	w := &responseDelegator{parent: parent}
	req.mux.ServeHTTP(w, r)
	if w.Failed() {
		return true, nil
	}

	// remember resource as created
	req.txn.Data[resKey] = append(created, objID)
	return false, nil
}

type request struct {
	*http.Request
	ctx  context.Context
	mux  http.Handler
	txn  *api.Txn
	path riposo.Path
}

func newRequest(r *http.Request) *request {
	ctx := context.WithValue(r.Context(), chi.RouteCtxKey, nil)
	return &request{
		Request: r,
		ctx:     ctx,
		mux:     api.GetMux(r),
		txn:     api.GetTxn(r),
		path:    api.GetPath(r),
	}
}

func (r *request) Context() context.Context {
	return r.ctx
}

func containsString(vv []string, s string) bool {
	for _, v := range vv {
		if s == v {
			return true
		}
	}
	return false
}

// --------------------------------------------------------------------

type responseDelegator struct {
	Code   int
	header http.Header
	parent http.ResponseWriter
}

func (rw *responseDelegator) Header() http.Header {
	m := rw.header
	if m == nil {
		m = make(http.Header)
		rw.header = m
	}
	return m
}

func (rw *responseDelegator) Write(buf []byte) (int, error) {
	rw.WriteHeader(http.StatusOK)
	// propagate on non-2xx
	if rw.Failed() {
		return rw.parent.Write(buf)
	}
	return len(buf), nil
}

func (rw *responseDelegator) WriteHeader(code int) {
	if rw.Code == 0 {
		rw.Code = code

		// propagate on non-2xx
		if rw.Failed() {
			for k := range rw.header {
				rw.parent.Header().Set(k, rw.header.Get(k))
			}
			rw.parent.WriteHeader(rw.Code)
		}
	}
}

func (rw *responseDelegator) Failed() bool {
	return rw.Code >= 300
}
