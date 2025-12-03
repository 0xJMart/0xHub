package main

import (
	"log"

	"0xhub/backend/internal/handlers"
	"0xhub/backend/internal/models"
	"0xhub/backend/internal/store"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize store and seed with sample data
	store := store.NewStore()
	seedProjects(store)

	// Initialize handlers
	projectsHandler := handlers.NewProjectsHandler(store)

	// Setup router
	router := gin.Default()

	// CORS configuration
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173", "http://localhost:3000"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	router.Use(cors.New(config))

	// Health check endpoint
	router.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// API routes
	api := router.Group("/api")
	{
		api.GET("/projects", projectsHandler.GetProjects)
		api.GET("/projects/:id", projectsHandler.GetProject)
		api.POST("/projects", projectsHandler.CreateProject)
		api.PUT("/projects/:id", projectsHandler.UpdateProject)
		api.DELETE("/projects/:id", projectsHandler.DeleteProject)
	}

	// Start server
	log.Println("Server starting on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// seedProjects adds sample projects for testing
func seedProjects(store *store.Store) {
	projects := []*models.Project{
		{
			ID:          "1",
			Name:        "Kubernetes",
			Description: "Production-Grade Container Orchestration",
			URL:         "https://kubernetes.io",
			Icon:        "https://kubernetes.io/images/favicon.png",
			Category:    "Infrastructure",
			Status:      "active",
		},
		{
			ID:          "2",
			Name:        "Docker",
			Description: "Empowering App Development for Developers",
			URL:         "https://www.docker.com",
			Icon:        "https://www.docker.com/favicon.ico",
			Category:    "Containerization",
			Status:      "active",
		},
		{
			ID:          "3",
			Name:        "Prometheus",
			Description: "Monitoring and alerting toolkit",
			URL:         "https://prometheus.io",
			Icon:        "https://prometheus.io/assets/favicon.ico",
			Category:    "Monitoring",
			Status:      "active",
		},
		{
			ID:          "4",
			Name:        "Grafana",
			Description: "The open observability platform",
			URL:         "https://grafana.com",
			Icon:        "https://grafana.com/favicon.ico",
			Category:    "Visualization",
			Status:      "active",
		},
		{
			ID:          "5",
			Name:        "Helm",
			Description: "The package manager for Kubernetes",
			URL:         "https://helm.sh",
			Icon:        "https://helm.sh/img/favicon-32x32.png",
			Category:    "DevOps",
			Status:      "active",
		},
		{
			ID:          "6",
			Name:        "Istio",
			Description: "Connect, secure, control, and observe services",
			URL:         "https://istio.io",
			Icon:        "https://istio.io/favicon.ico",
			Category:    "Service Mesh",
			Status:      "active",
		},
	}

	for _, project := range projects {
		store.Create(project)
	}
	log.Println("Seeded", len(projects), "sample projects")
}

