package runner

import (
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"

	httpapi "github.com/authena-ru/courses-organization/internal/adapter/delivery/http"
	mongorepo "github.com/authena-ru/courses-organization/internal/adapter/repository/mongodb"
	"github.com/authena-ru/courses-organization/internal/app"
	"github.com/authena-ru/courses-organization/internal/app/command"
	"github.com/authena-ru/courses-organization/internal/app/command/mock"
	"github.com/authena-ru/courses-organization/internal/config"
	"github.com/authena-ru/courses-organization/internal/server"
	"github.com/authena-ru/courses-organization/pkg/database/mongodb"
)

func Start(configsDir string) {
	cfg := newConfig(configsDir)
	db := newMongoDatabase(cfg)
	application := newApplication(db)
	startServer(cfg, application)
}

func newConfig(configsDir string) *config.Config {
	cfg, err := config.New(configsDir)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to parse config")
	}
	return cfg
}

func newMongoDatabase(cfg *config.Config) *mongo.Database {
	client, err := mongodb.NewClient(cfg.Mongo.URI, cfg.Mongo.Username, cfg.Mongo.Password)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to connect to MongoDB")
	}
	return client.Database(cfg.Mongo.DatabaseName)
}

func newApplication(db *mongo.Database) app.Application {
	coursesRepository := mongorepo.NewCoursesRepository(db)
	academicsService := mock.NewAcademicsService(
		[]string{"d3e2490f-5944-4a87-b29a-94177d1caaed", "4edefb83-4b6b-479d-9ce2-60cd465630b6"},
		[]string{"798155cb-91b7-41d4-9f91-a1970339707e"},
		[]string{"95dca190-f307-4954-8700-f992f8c12a86"},
	)
	return app.Application{
		Commands: app.Commands{
			CreateCourse:       command.NewCreateCourseHandler(coursesRepository),
			ExtendCourse:       command.NewExtendCourseHandler(coursesRepository),
			AddCollaborator:    command.NewAddCollaboratorHandler(coursesRepository, academicsService),
			RemoveCollaborator: command.NewRemoveCollaboratorHandler(coursesRepository),
			AddStudent:         command.NewAddStudentHandler(coursesRepository, academicsService),
			RemoveStudent:      command.NewRemoveStudentHandler(coursesRepository),
			AddTask:            command.NewAddTaskHandler(coursesRepository),
		},
	}
}

func startServer(cfg *config.Config, application app.Application) {
	logrus.Info("Starting HTTP server on address :8080")
	httpServer := server.New(cfg, httpapi.NewHandler(application))
	err := httpServer.Run()
	logrus.WithError(err).Fatal("HTTP server stopped")
}
