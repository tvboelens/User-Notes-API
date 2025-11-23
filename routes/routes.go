package routes

import (
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

	login_service := services.NewLoginService(&pwd_hasher, user_repo, jwt_secret)
	registration_service := services.NewRegistrationService(&pwd_hasher, user_repo, jwt_secret)

	note_service := services.NewNoteService(note_repo, note_repo, user_repo)
	note_controller := controllers.NewNoteController(note_service, note_service)

	auth_controller := controllers.NewAuthController(login_service, registration_service)
	r.POST("/register", auth_controller.Register)
	r.POST("/login", auth_controller.Login)

	auth := r.Group("/")
	auth.Use(middleware.JwtMiddleware(jwt_secret))
	auth.POST("/notes", note_controller.Create)
	// Todo: GET /notes, PUT /notes/:id, DELETE /notes/:id
}
