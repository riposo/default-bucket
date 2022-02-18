package internal_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	. "github.com/bsm/ginkgo/v2"
	. "github.com/bsm/gomega"
	"github.com/riposo/default-bucket/internal"
	"github.com/riposo/riposo/pkg/api"
	"github.com/riposo/riposo/pkg/mock"
	"github.com/riposo/riposo/pkg/riposo"
)

var _ = Describe("Config", func() {
	var txn *api.Txn
	var rts *api.Routes

	serve := func(r *http.Request) *httptest.ResponseRecorder {
		w := httptest.NewRecorder()
		rts.Mux().ServeHTTP(w, r)
		return w
	}

	BeforeEach(func() {
		rts = api.NewRoutes(&api.Config{
			Authz: api.Authz{
				"write": {"account:alice"},
				"read":  {"account:bob"},
			},
		})
		rts.Resource("/buckets", nil)
		rts.Resource("/buckets/{bucket_id}/collections", nil)
		rts.Resource("/buckets/{bucket_id}/collections/{collection_id}/records", nil)

		cfg := new(internal.Config)
		cfg.DefaultBucket.Secret = internal.HashSecret("foo")
		cfg.Mount(rts)

		txn = mock.Txn()
		txn.User = mock.User("account:alice")
	})

	AfterEach(func() {
		Expect(txn.Rollback()).To(Succeed())
	})

	It("rejects unauthenticated requests", func() {
		txn.User = mock.User(riposo.Everyone)

		w := serve(mock.Request(txn, http.MethodGet, "/buckets/default", nil))
		Expect(w.Code).To(Equal(http.StatusForbidden))
	})

	It("re-routes requests", func() {
		w := serve(mock.Request(txn, http.MethodPatch, "/buckets/default", strings.NewReader(`{"data":{"meta":"data"}}`)))
		Expect(w.Code).To(Equal(http.StatusOK))
		Expect(w.Body.Bytes()).To(MatchJSON(`{
			"data": {
				"id": "a53aa5f020d80439829eda9f6e3a4502",
				"last_modified": 1515151515678,
				"meta": "data"
			},
			"permissions": {
				"write": ["account:alice"]
			}
		}`))
	})

	It("creates default buckets via PUT", func() {
		w := serve(mock.Request(txn, http.MethodPut, "/buckets/default", strings.NewReader(`{"data":{"id":"default"}}`)))
		Expect(w.Code).To(Equal(http.StatusCreated))
		Expect(w.Body.Bytes()).To(MatchJSON(`{
			"data": {
				"id": "a53aa5f020d80439829eda9f6e3a4502",
				"last_modified": 1515151515677
			},
			"permissions": {
				"write": ["account:alice"]
			}
		}`))
	})

	It("auto-provisions buckets", func() {
		w := serve(mock.Request(txn, http.MethodGet, "/buckets", nil))
		Expect(w.Code).To(Equal(http.StatusOK))
		Expect(w.Body.Bytes()).To(MatchJSON(`{"data": []}`))

		w = serve(mock.Request(txn, http.MethodGet, "/buckets/default/collections", nil))
		Expect(w.Code).To(Equal(http.StatusOK))
		Expect(w.Body.Bytes()).To(MatchJSON(`{"data": []}`))

		w = serve(mock.Request(txn, http.MethodGet, "/buckets", nil))
		Expect(w.Code).To(Equal(http.StatusOK))
		Expect(w.Body.Bytes()).To(MatchJSON(`{
			"data": [
				{
					"id": "a53aa5f020d80439829eda9f6e3a4502",
					"last_modified": 1515151515677
				}
			]
		}`))
	})

	It("creates collections via PUT", func() {
		w := serve(mock.Request(txn, http.MethodPut, "/buckets/default/collections/foo", strings.NewReader(`{}`)))
		Expect(w.Code).To(Equal(http.StatusCreated))
		Expect(w.Body.Bytes()).To(MatchJSON(`{
			"data": {
				"id": "foo",
				"last_modified": 1515151515677
			},
			"permissions": {
				"write": ["account:alice"]
			}
		}`))

		w = serve(mock.Request(txn, http.MethodGet, "/buckets", nil))
		Expect(w.Code).To(Equal(http.StatusOK))
		Expect(w.Body.Bytes()).To(MatchJSON(`{
			"data": [
				{
					"id": "a53aa5f020d80439829eda9f6e3a4502",
					"last_modified": 1515151515677
				}
			]
		}`))
	})

	It("auto-provisions collections", func() {
		w := serve(mock.Request(txn, http.MethodGet, "/buckets/default/collections", nil))
		Expect(w.Code).To(Equal(http.StatusOK))
		Expect(w.Body.Bytes()).To(MatchJSON(`{"data": []}`))

		w = serve(mock.Request(txn, http.MethodGet, "/buckets", nil))
		Expect(w.Code).To(Equal(http.StatusOK))
		Expect(w.Body.Bytes()).To(MatchJSON(`{
			"data": [
				{
					"id": "a53aa5f020d80439829eda9f6e3a4502",
					"last_modified": 1515151515677
				}
			]
		}`))

		w = serve(mock.Request(txn, http.MethodGet, "/buckets/default/collections/foo/records", nil))
		Expect(w.Code).To(Equal(http.StatusOK))
		Expect(w.Body.Bytes()).To(MatchJSON(`{"data": []}`))

		w = serve(mock.Request(txn, http.MethodGet, "/buckets/default/collections", nil))
		Expect(w.Code).To(Equal(http.StatusOK))
		Expect(w.Body.Bytes()).To(MatchJSON(`{
			"data": [
				{
					"id": "foo",
					"last_modified": 1515151515677
				}
			]
		}`))
	})

	It("allows default buckets without explicit permission", func() {
		// bob cannot create buckets
		txn.User = mock.User("account:bob")

		w := serve(mock.Request(txn, http.MethodGet, "/buckets/default", nil))
		Expect(w.Code).To(Equal(http.StatusOK))
		Expect(w.Body.Bytes()).To(MatchJSON(`{
			"data": {
				"id": "b40605613e14af3f78da6c99efe88224",
				"last_modified": 1515151515677
			},
			"permissions": {
				"write": ["account:bob"]
			}
		}`))
	})
})

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "internal")
}
