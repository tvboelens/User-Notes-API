package routes

import (
	"net/http"
	"user-notes-api/auth"
	"user-notes-api/controllers"
	"user-notes-api/middleware"
	"user-notes-api/repositories"
	"user-notes-api/services"
	"user-notes-api/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"runtime"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB, jwt_secret string) {
	user_repo := repositories.NewUserRepository(db)
	note_repo := repositories.NewNoteRepository(db)

	threads := uint8(runtime.GOMAXPROCS(0))
	pwd_hasher := utils.Argon2IdHasher{Time: 1, SaltLen: 32, Memory: 64 * 1024, Threads: threads, KeyLen: 256}

	login_manager := auth.LoginManager{UserReader: user_repo, PwdComparer: &pwd_hasher}
	registration_manager := auth.RegistrationManager{UserCreator: user_repo, PwdHasher: &pwd_hasher}

	login_service := services.NewLoginService(&login_manager, jwt_secret)
	registration_service := services.NewRegistrationService(&registration_manager, jwt_secret)

	note_service := services.NewNoteService(note_repo, note_repo, user_repo)
	note_controller := controllers.NewNoteController(note_service, note_service)

	r.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	auth_controller := controllers.NewAuthController(login_service, registration_service)
	r.POST("/register", auth_controller.Register)
	r.POST("/login", auth_controller.Login)

	auth := r.Group("/")
	auth.Use(middleware.JwtMiddleware(jwt_secret))
	auth.POST("/notes", note_controller.Create)
	auth.GET("/notes", note_controller.GetNotes)
	auth.GET("/notes/:id", note_controller.GetSingleNote)
	// Todo: GET /notes, PUT /notes/:id, DELETE /notes/:id
}
