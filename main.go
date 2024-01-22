package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"log"
	"manodarpanNewproject/authentication"
	"manodarpanNewproject/pkg/dbase"
	"manodarpanNewproject/pkg/userdao"
	"manodarpanNewproject/pkg/utilities"
	"manodarpanNewproject/server"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var Env server.Envs

func main() {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "db-username",
			Value:       "postgres",
			Usage:       "database username",
			EnvVar:      "DB_USERNAME",
			Destination: &Env.Database.Username,
		},
		cli.StringFlag{
			Name:        "db-password",
			Value:       "Sai@996361",
			Usage:       "database password",
			EnvVar:      "DB_PASSWORD",
			Destination: &Env.Database.Password,
		},
		cli.StringFlag{
			Name:        "db-hostname",
			Value:       "localhost",
			Usage:       "database hostname",
			EnvVar:      "DB_HOSTNAME",
			Destination: &Env.Database.Hostname,
		},
		cli.StringFlag{
			Name:        "db-port",
			Value:       "5432",
			Usage:       "database port",
			EnvVar:      "DB_PORT",
			Destination: &Env.Database.Port,
		},
		cli.StringFlag{
			Name:        "db-database",
			Value:       "testnew",
			Usage:       "database name",
			EnvVar:      "DB_DATABASE",
			Destination: &Env.Database.Database,
		}, cli.StringFlag{
			Name:        "port",
			Value:       "9000",
			Usage:       "server port",
			EnvVar:      "PORT",
			Destination: &Env.Port,
		},
	}

	app.Action = Run

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

}

func Run(_ *cli.Context) error {
	startServer := time.Now()

	fmt.Println("entered into run function")

	// wait for the pg database to be available
	if err := waitForHost(Env.Database.Hostname, Env.Database.Port); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Connecting to database, %v\n", Env.Database.Hostname)
	log.Println("database", Env.Database.Database)
	db, err := Env.Database.Open()
	if err != nil {
		log.Fatal("Error While Connecting to Database ", err)
	}
	fmt.Println("result for db", db)
	fmt.Println("Creating data access object handle")
	err = dbase.RunMigrations(db)
	if err != nil {
		log.Println("ERROR", err)
		log.Fatal(err)
	}
	dao := userdao.New(db)

	svr := server.Server{
		Env: Env,
		DAO: dao,
	}

	fmt.Println("Done!")

	fmt.Println("Creating data access object handle")

	ctx := context.Background()

	gin.SetMode(gin.ReleaseMode)
	routes := gin.New()
	routes.Use(
		gin.Recovery(),
		Authentication.CORS(),
		gin.Logger(),
	)

	v1 := routes.Group("/v1")
	v1.POST("/registration", svr.UserRegistration)
	v1.POST("/login", svr.LoginUser)
	v1.GET("/user/profile", Authentication.Auth(), svr.UserProfile)

	// Health APIs
	v1.GET("/readiness", func(c *gin.Context) {
		fmt.Println("called readiness")
		c.JSON(http.StatusOK, gin.H{
			"status": "UP",
		})
		return
	})

	srv := &http.Server{
		Addr:    ":" + Env.Port,
		Handler: routes,
	}

	fmt.Printf("serverStartSpan: %v\n", utilities.GetDuration(startServer))
	// launch our web server
	fmt.Printf("webserver listening on port %v\n", Env.Port)
	if err := srv.ListenAndServe(); err != nil {
		log.Printf("listen: %s\n", err)
		panic("at starting server " + err.Error())
	}
	//}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")

	return nil
}

func waitForHost(host, port string) error {
	timeOut := time.Second

	if host == "" {
		return errors.Errorf("unable to connect to %v:%v", host, port)
	}

	for i := 0; i < 60; i++ {
		fmt.Printf("waiting for %v:%v ...\n", host, port)
		conn, err := net.DialTimeout("tcp", host+":"+port, timeOut)
		if err == nil {
			fmt.Println("done!")
			conn.Close()
			return nil
		}

		time.Sleep(time.Second)
	}

	return errors.Errorf("timeout attempting to connect to %v:%v", host, port)
}
