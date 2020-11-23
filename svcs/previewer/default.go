package previewer

import (
	"github.com/foxfoxio/codelabs-preview-go/internal/gdoc"
	"github.com/foxfoxio/codelabs-preview-go/internal/gdrive"
	"github.com/foxfoxio/codelabs-preview-go/svcs/previewer/endpoints"
	"github.com/foxfoxio/codelabs-preview-go/svcs/previewer/transports"
	"github.com/foxfoxio/codelabs-preview-go/svcs/previewer/usecases"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"os"
)

func New(rootRouter *mux.Router) {
	clientId := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	config := &oauth2.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		Endpoint:     google.Endpoint,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/drive.file",
			"https://www.googleapis.com/auth/drive.readonly",
			"openid",
		},
		RedirectURL: os.Getenv("GOOGLE_REDIRECT_URL"),
	}

	templateId := os.Getenv("CP_TEMPLATE_ID")
	driveRootId := os.Getenv("CP_DRIVE_ROOT_ID")
	adminEmail := os.Getenv("CP_ADMIN_EMAIL")

	if templateId == "" {
		templateId = "1X3kriKmznxdBrJ1U4NLVtM_kLHRJBXEjn92iZI9XcW4"
	}

	if driveRootId == "" {
		driveRootId = "1uH1lq__vo-PTusArFsOduKfHk6ZhW1gX"
	}

	store := sessions.NewCookieStore([]byte("t0p-secret"))
	driveClient := gdrive.NewClient()
	gdocClient := gdoc.NewClient()

	sessionUsecase := usecases.NewSession(store, "__session")
	viewerUsecase := usecases.NewViewer(driveClient, gdocClient, templateId, driveRootId, adminEmail)
	authUsecase := usecases.NewAuth(config)

	authEp := endpoints.NewAuth(sessionUsecase, authUsecase)
	viewerEp := endpoints.NewViewer(sessionUsecase, viewerUsecase, authUsecase)

	transports.RegisterHttpRouter(rootRouter, authEp, viewerEp)
}
