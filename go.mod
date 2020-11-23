module github.com/foxfoxio/codelabs-preview-go

go 1.13

require (
	cloud.google.com/go/storage v1.10.0
	github.com/googlecodelabs/tools/claat v0.0.0-20200918190358-3cc6629c4d3d
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/sessions v1.2.1
	github.com/rs/xid v1.2.1
	github.com/sirupsen/logrus v1.7.0
	github.com/stretchr/testify v1.4.0
	golang.org/x/net v0.0.0-20201016165138-7b1cca2348c0
	golang.org/x/oauth2 v0.0.0-20200902213428-5d25da1a8d43
	google.golang.org/api v0.30.0
)

replace github.com/googlecodelabs/tools/claat v0.0.0-20200918190358-3cc6629c4d3d => github.com/foxfoxio/tools/claat v0.0.0-20201121171015-125228c13156

replace gopkg.in/russross/blackfriday.v2 v2.0.1 => github.com/russross/blackfriday/v2 v2.0.1
