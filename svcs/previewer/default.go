package previewer

import (
	"github.com/foxfoxio/codelabs-preview-go/internal/bootstrap"
	"github.com/foxfoxio/codelabs-preview-go/internal/gdoc"
	"github.com/foxfoxio/codelabs-preview-go/internal/gdrive"
	"github.com/foxfoxio/codelabs-preview-go/internal/gstorage"
	"github.com/foxfoxio/codelabs-preview-go/internal/xfirebase"
	"github.com/foxfoxio/codelabs-preview-go/svcs/previewer/endpoints"
	"github.com/foxfoxio/codelabs-preview-go/svcs/previewer/transports"
	"github.com/foxfoxio/codelabs-preview-go/svcs/previewer/usecases"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"os"
	"strings"
)

func New() *bootstrap.Server {
	templateId := os.Getenv("CP_TEMPLATE_ID")
	driveRootId := os.Getenv("CP_DRIVE_ROOT_ID")
	driveTempId := os.Getenv("CP_DRIVE_TEMP_ID")
	adminEmail := os.Getenv("CP_ADMIN_EMAIL")
	bucketName := os.Getenv("CP_BUCKET_NAME")
	storagePath := os.Getenv("CP_STORAGE_PATH")
	allowedOriginsStr := os.Getenv("CP_ALLOWED_ORIGIN")

	allowedOrigins := make([]string, 0)
	if x := strings.TrimSpace(allowedOriginsStr); x != "" {
		allowedOrigins = append(allowedOrigins, strings.Split(x, ",")...)
	}

	if templateId == "" {
		templateId = "1X3kriKmznxdBrJ1U4NLVtM_kLHRJBXEjn92iZI9XcW4"
	}

	if driveRootId == "" {
		driveRootId = "1uH1lq__vo-PTusArFsOduKfHk6ZhW1gX"
	}

	if driveTempId == "" {
		driveTempId = "1uH1lq__vo-PTusArFsOduKfHk6ZhW1gX"
	}

	if bucketName == "" {
		bucketName = "codelabs-preview"
	}

	if storagePath == "" {
		storagePath = "files-dev"
	}

	driveClient := gdrive.NewClient()
	gdocClient := gdoc.NewClient()
	gStorageClient := gstorage.NewClient(bucketName)
	firebaseClient := xfirebase.NewDefaultClient()

	authUsecase := usecases.NewAuth(firebaseClient)
	viewerUsecase := usecases.NewViewer(driveClient, gdocClient, gStorageClient, templateId, driveRootId, adminEmail, storagePath, driveTempId)

	viewerEp := endpoints.NewViewer(viewerUsecase)
	router := mux.NewRouter()
	transports.RegisterHttpRouter(router, viewerEp)
	withAuth := authUsecase.Handler(router)

	options := cors.Options{
		AllowCredentials: true,
		Debug:            true,
	}
	if len(allowedOrigins) > 0 {
		options.AllowedOrigins = allowedOrigins
	}
	c := cors.New(options)
	server := &bootstrap.Server{
		HttpHandler: c.Handler(withAuth),
	}

	return server
}
