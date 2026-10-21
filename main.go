package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Dhorq/job-tracker-lab/internal/database"
	"github.com/Dhorq/job-tracker-lab/internal/models"
	"github.com/Dhorq/job-tracker-lab/internal/rabbitmq"
	"github.com/Dhorq/job-tracker-lab/internal/repository"
	"github.com/gin-gonic/gin"
)

func main() {
	var err error
	db, err := database.Connect()
	if err != nil {
		fmt.Printf("WARNING: Failed to connect to database: %v\n", err)
	}
	defer db.Close()

	if err := database.Migrate(db); err != nil {
		fmt.Printf("WARNING: Failed to migrate database: %v\n", err)
	}
	fmt.Println("Table is ready to use")

	rq, err := rabbitmq.New("amqp://guest:guest@rabbitmq:5672/", "job_applications")
	if err != nil {
		log.Fatalf("Failed to connect Rabbitmq: %v", err)
	}
	defer rq.Close()
	fmt.Println("RabbitMQ Connected")

	go rq.Consume()

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello World")
	})

	r.POST("/applications", func(c *gin.Context) {
		var app models.JobApplication
		if err := c.ShouldBindJSON(&app); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		id, err := repository.Create(db, app)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		msg := fmt.Sprintf("Pelamar baru: ID %d, Perusahaan %s, Posisi %s", id, app.CompanyName, app.Position)
		if err := rq.Publish(msg); err != nil {
			log.Printf("Failed to publish message: %v", err)
		}
		fmt.Println("Message published to RabbitMQ")

		c.JSON(http.StatusCreated, gin.H{"id": id, "message": "Data berhasil disimpan"})
	})

	r.GET("/applications", func(c *gin.Context) {
		apps, err := repository.GetAll(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, apps)
	})

	r.Run(":8080")
}
