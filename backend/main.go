package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"

	"genai.ai/automonitor/handlers"
	"genai.ai/automonitor/service"
)


func ReverseProxy(frontURL string) gin.HandlerFunc {
	return func(c *gin.Context) {
		remote, _ := url.Parse(frontURL)
		proxy := httputil.NewSingleHostReverseProxy(remote)
		proxy.Director = func(req *http.Request) {
			req.Header = c.Request.Header
			req.Host = remote.Host
			req.URL.Scheme = remote.Scheme
			req.URL.Host = remote.Host
			req.URL.Path = c.Request.URL.Path
		}

		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
  
func main() {

	logger := log.New(os.Stdout, "", log.LstdFlags)

	cfg, err := service.LoadConfig()
	if err != nil {
		log.Fatalf("error loading configuration: %v", err)
	}

	//Connect MySQL database
	err = service.ConnectDB(cfg.DBName, cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort)
	if err != nil {
		logger.Println("[MySQL connect failed]: please check MySQL dbHost, dbUser,dbPasword, dbName .... ")
		return
	}

	//Start Job Scheduler
	service.StartScheduler()


	//setup router
	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	authMiddleware := handlers.AuthMiddleware()
	router.Use(authMiddleware)

	router.GET("/", ReverseProxy(cfg.Front))

	

	v1:=router.Group("/api/v1")

	v1.POST("/login", handlers.SignInHandler)
	
	v1.GET("/projects", handlers.GetProjects)
	v1.GET("/catalogs", handlers.GetCatalogs)

	v1.GET("/metrics", handlers.GetMetrics)
	v1.POST("/metric", handlers.CreateMonitorMetric)
	v1.GET("/metric/:id", handlers.GetMonitorMetricByID)
	v1.DELETE("/metric/:id", handlers.DeleteMonitorMetric)
	v1.PUT("/metric/:id", handlers.UpdateMonitorMetric)
	

	v1.GET("/connections", handlers.GetMonitorConnections)
	v1.GET("/connections/name",handlers.GetMonitorConnectionName)
	v1.POST("/connection", handlers.CreateMonitorConnection)
	v1.GET("/connection/:id", handlers.GetMonitorConnection)
	v1.PUT("/connection/:id", handlers.UpdateMonitorConnection)
	v1.DELETE("/connection/:id", handlers.DeleteMonitorConnection)
	
	

	v1.GET("/jobs", handlers.GetMonitorJobs)
	v1.POST("/job", handlers.CreateMonitorJob)
	v1.GET("/job/:id", handlers.GetMonitorJob)
	v1.PUT("/job/:id", handlers.UpdateMonitorJob)
	v1.DELETE("/job/:id", handlers.DeleteMonitorJob)

	v1.GET("/status/:project", handlers.GetJobStatus)

	v1.POST("/run/:id", handlers.RunMetricByID)

	v1.GET("/config", handlers.Config)

	// Serve a staic endpoint
	router.Static("/static", "./static")

	
	router.NoRoute(ReverseProxy(cfg.Front))


	router.Run(":8080")
}
