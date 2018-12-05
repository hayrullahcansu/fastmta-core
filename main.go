package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"time"

	"./boss"
	OS "./cross"
	"./logger"

	"./iris/web/controllers"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
)

var ops int64

const (
	Timeout   = time.Second * time.Duration(30)
	KeepAlive = time.Second * time.Duration(30)
	//MtaName       = "ZetaMail"
	//MaxErrorLimit = 10
)

func main() {
	// load command line arguments
	start := time.Now()
	//m := &sync.Mutex{}

	//runtime.Goexit()
	elapsed := time.Since(start)
	log.Printf("Binomial took %s", elapsed)

	name := flag.String("name", "FastMta", "name to print")
	flag.Parse()
	log.Printf("Starting service for %s%s", *name, OS.NewLine)
	// setup signal catching
	sigs := make(chan os.Signal, 1)
	// catch all signals since not explicitly listing
	signal.Notify(sigs)
	//signal.Notify(sigs,syscall.SIGQUIT)
	// method invoked upon seeing signal

	go func() {
		s := <-sigs
		logger.Info.Printf("RECEIVED SIGNAL: %s%s", s, OS.NewLine)
		AppCleanup()
		os.Exit(1)
	}()

	boss.InitSystem()
	boss := boss.New()
	// rabbitClient := queue.New()
	// rabbitClient.Connect(true)
	// _, _ = rabbitClient.Consume(queue.InboundQueueName, "", false, false, true, nil)

	// rabbitClient2 := queue.New()
	// rabbitClient2.Connect(true)
	// _, _ = rabbitClient2.Consume(queue.InboundQueueName, "", false, false, true, nil)

	boss.Run()
	// infinite print loop
	app := iris.New()
	// You got full debug messages, useful when using MVC and you want to make
	// sure that your code is aligned with the Iris' MVC Architecture.
	app.Logger().SetLevel("debug")
	// Load the template files.
	tmpl := iris.HTML("./views", ".html").
		Layout("shared/layout.html").
		Reload(true)
	app.RegisterView(tmpl)

	app.StaticWeb("/", "./public")

	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewData("Message", ctx.Values().
			GetStringDefault("message", "The page you're looking for doesn't exist"))
		ctx.View("shared/error.html")
	})

	// ---- Serve our controllers. ----

	// "/users" based mvc application.

	defaultDashboard := mvc.New(app.Party("/"))
	dashboard := mvc.New(app.Party("/dashboard"))
	// Add the basic authentication(admin:password) middleware
	// for the /users based requests.
	//users.Router.Use(middleware.BasicAuth)
	// Bind the "userService" to the UserController's Service (interface) field.

	defaultDashboard.Handle(new(controllers.DashboardController))
	dashboard.Handle(new(controllers.DashboardController))

	// "/user" based mvc application.
	// sessManager := sessions.New(sessions.Config{
	// 	Cookie:  "sessioncookiename",
	// 	Expires: 24 * time.Hour,
	// })

	// http://localhost:8080/noexist
	// and all controller's methods like
	// http://localhost:8080/users/1
	// http://localhost:8080/user/register
	// http://localhost:8080/user/login
	// http://localhost:8080/user/me
	// http://localhost:8080/user/logout
	// basic auth: "admin", "password", see "./middleware/basicauth.go" source file.
	app.Run(
		// Starts the web server at localhost:8080
		iris.Addr("localhost:8080"),
		// Ignores err server closed log when CTRL/CMD+C pressed.
		iris.WithoutServerError(iris.ErrServerClosed),
		// Enables faster json serialization and more.
		iris.WithOptimizations,
	)
	select {}

}
func AppCleanup() {
	time.Sleep(time.Millisecond * time.Duration(1000))
	logger.Info.Println("CLEANUP APP BEFORE EXIT!!!")
}
