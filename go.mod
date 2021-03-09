module github.com/riposo/default-bucket

go 1.14

replace github.com/riposo/riposo => ../riposo

require (
	github.com/go-chi/chi v1.5.3
	github.com/onsi/ginkgo v1.14.2
	github.com/onsi/gomega v1.10.4
	github.com/riposo/riposo v0.0.0-20210226155134-b4e129732a1c
	github.com/valyala/bytebufferpool v1.0.0
	golang.org/x/crypto v0.0.0-20210220033148-5ea612d1eb83
)
