package routers

import (
	"fmt"
	"log"

	"emma/cmd/configs"
	eventsHandlers "emma/internal/adapter/handlers/events"
	usersHandlers "emma/internal/adapter/handlers/users"
	eventsAdapter "emma/internal/adapter/mysql/events"
	usersAdapter "emma/internal/adapter/mysql/users"
	events "emma/internal/core/domain/events"
	users "emma/internal/core/domain/users"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)


func InitRouters(r *gin.Engine) {
	dbName := configs.GetConfig().MYSQL_DB
	dbUser := configs.GetConfig().MYSQL_USER
	dbPassword := configs.GetConfig().MYSQL_PASSWORD
	dbTcp := configs.GetConfig().MYSQL_TCP
	dbDriver := dbUser+":"+dbPassword+dbTcp+dbName+
	"?charset=utf8&parseTime=True"
  db, err := gorm.Open(mysql.Open(dbDriver), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to database", err)
	}

	usersMySqlAdapter := usersAdapter.NewMySQLAdapter(db)
	usersService := users.UserService(usersMySqlAdapter)
	usersAdapter := usersHandlers.NewUserRESTAdapter(usersService)
	usersRoutes := r.Group("/userss")
	{
		usersRoutes.POST("/", usersAdapter.CreateUser)
		usersRoutes.GET("/:username", usersAdapter.FetchUser)
	}

	eventsMySqlAdapter := eventsAdapter.NewMySQLAdapter(db)
	eventsService := events.EventService(eventsMySqlAdapter)
	eventsAdapter := eventsHandlers.NewEventRESTAdapter(eventsService)
	eventsRoutes := r.Group("/events")
	{
		eventsRoutes.POST("/", eventsAdapter.CreateEvent)
		eventsRoutes.GET("/:id", eventsAdapter.FetchEvent)
		eventsRoutes.GET("/", eventsAdapter.GetAllEvents)
		eventsRoutes.PUT("/:id", eventsAdapter.UpdateEvent)
		eventsRoutes.DELETE("/:id", eventsAdapter.DeleteEvent)
	}

	fmt.Println("InitRouters main rout")
}
