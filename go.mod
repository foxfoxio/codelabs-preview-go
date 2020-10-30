module github.com/foxfoxio/codelabs-preview-go

go 1.13

require (
	github.com/googlecodelabs/tools/claat v0.0.0-20200918190358-3cc6629c4d3d
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/sessions v1.2.1
	github.com/rs/xid v1.2.1
	github.com/sirupsen/logrus v1.7.0
	golang.org/x/oauth2 v0.0.0-20200902213428-5d25da1a8d43
)

replace github.com/googlecodelabs/tools/claat v0.0.0-20200918190358-3cc6629c4d3d => github.com/foxfoxio/tools/claat v0.0.0-20201030163433-0f71522d5224

replace gopkg.in/russross/blackfriday.v2 v2.0.1 => github.com/russross/blackfriday/v2 v2.0.1
