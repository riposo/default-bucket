module github.com/riposo/default-bucket

go 1.14

replace github.com/riposo/riposo => ../riposo

require (
	github.com/go-chi/chi v1.5.4
	github.com/onsi/ginkgo v1.15.2
	github.com/onsi/gomega v1.11.0
	github.com/riposo/riposo v0.0.0-00010101000000-000000000000
	github.com/valyala/bytebufferpool v1.0.0
	golang.org/x/crypto v0.0.0-20210314154223-e6e6c4f2bb5b
)
