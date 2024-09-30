package routers

import (
	"emma/dbs"
	middlewares "emma/http/middlewares"
	eventsHandlers "emma/internal/adapter/handlers/events"
	usersHandlers "emma/internal/adapter/handlers/users"
	eventsAdapter "emma/internal/adapter/mysql/events"
	usersAdapter "emma/internal/adapter/mysql/users"
	"emma/internal/domain"

	"github.com/gin-gonic/gin"
)


func InitRouters(r *gin.Engine) {
	db := dbs.InitDB()

	usersMySqlAdapter := usersAdapter.NewMySQLAdapter(db)
	usersRepository := domain.UserRepository(usersMySqlAdapter)
	usersAdapter := usersHandlers.NewUserRESTAdapter(usersRepository)
	usersRoutes := r.Group("/users")
	{
		usersRoutes.POST("/", usersAdapter.CreateUser)
		usersRoutes.GET("/:id", usersAdapter.FetchUser)
		usersRoutes.POST("/login", usersAdapter.Login)
		usersRoutes.PUT("/:eventId", middlewares.Authorize("admin", "user"), usersAdapter.UpdateUser)
		usersRoutes.GET("/:id/events", middlewares.Authorize("admin", "user"), usersAdapter.GetUsersEvents)
	}

	eventsMySqlAdapter := eventsAdapter.NewMySQLAdapter(db)
	eventsRepository := domain.EventRepository(eventsMySqlAdapter)
	eventsAdapter := eventsHandlers.NewEventRESTAdapter(eventsRepository)
	eventsRoutes := r.Group("/events")
	{
		eventsRoutes.POST("/", middlewares.Authorize("admin"), eventsAdapter.CreateEvent)
		eventsRoutes.GET("/:id", middlewares.Authorize("admin","user"), eventsAdapter.FetchEvent)
		eventsRoutes.GET("/", middlewares.Authorize("admin","user"), eventsAdapter.GetAllEvents)
		eventsRoutes.PUT("/:id", middlewares.Authorize("admin"), eventsAdapter.UpdateEvent)
		eventsRoutes.DELETE("/:id", middlewares.Authorize("admin"), eventsAdapter.DeleteEvent)
	}
}
