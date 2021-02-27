package apps

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/itsjamie/gin-cors"
	"so-crm/src/libs"
	"so-crm/src/models"
	"time"
)

type App struct {
	configFile string
	app        *gin.Engine
	config     models.Config
}

func NewApp() *App {
	return &App{
		configFile: "config.yml",
	}
}

func (a *App) run() {
	a.loadConfig()
	a.runDatabase()
}

func (a *App) GetApp() *gin.Engine {
	return a.app
}

func (a *App) loadConfig() {
	env := libs.Env{ConfigFile: a.configFile}
	env.ReadFile(&a.config)
	env.ReadEnv(&a.config)
}

func (a *App) runDatabase() {
	models.ConnectDataBase(a.config)
}

func (a *App) cors() {
	a.GetApp().Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     true,
		ValidateHeaders: false,
	}))
}

func (a *App) middleware() {
	//a.GetApp().Use(middleware.AuthorizeJWT())
}

func (a *App) Routes(f func(*App)) {
	a.app = gin.Default()
	a.middleware()
	a.cors()
	f(a)
}

func (a *App) Start() {
	a.run()
	_ = fmt.Sprintf("Application running on %s:%d", a.config.Server.Host, a.config.Server.Port)
	if err := a.app.Run(
		fmt.Sprintf(
			"%s:%d",
			a.config.Server.Host,
			a.config.Server.Port,
		),
	); err != nil {
		panic(err)
	}
}
