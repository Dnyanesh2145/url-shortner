package main

import (
	"fiber-url-shortner/config"
	"fiber-url-shortner/database"
	// "fiber-url-shortner/helpers"
	"fiber-url-shortner/routes"
	"fiber-url-shortner/utils"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

var requestChan = make(chan fiber.Ctx, 100)

func setupRoutes(app *fiber.App) {
	app.Get("/health", routes.Health)
	app.Get("/u/:url", routes.Resolve)
	app.Post("/v1/shorten", routes.Shorten)
}

func main() {
	app := fiber.New()
	database.DBConnect()
	database.RedisConnect()

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("reqid", c.Get("X-Request-Id", config.CONFIG_DEFAULT_REQID))
		c.Locals("device", c.Get("X-Device-INFO", "device;os;platform;platformversion;appversion;"))
		c.Locals("latlong", c.Get("X-LatLong", "lat;long"))
		c.Locals("userId", config.CONFIG_DEFAULT_UID)
		c.Locals("msg")
		return c.Next() // Proceed with processing the request
	})
	// Initialize default config
	// app.Use(helpers.CheckIPAdress)
	// Or extend your config for customization
	intvalue, _ := strconv.Atoi(config.EnvDBURI("RATELIMITIG_PER_MINUTE"))

	app.Use(limiter.New(limiter.Config{
		Max:        intvalue,
		Expiration: 60 * time.Second,
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(429).JSON(fiber.Map{"error": "Rate limit exceeded"})
		},
	}))

	utils.LogStringChannel = make(utils.LogStringChannelWriter, 1000)
	go utils.LogStringChannelConsumer()

	app.Use(logger.New(logger.Config{
		Format:     config.CONFIG_LOGFORMAT,
		TimeFormat: config.CONFIG_LOGTIME_FORMAT,
		Output:     utils.LogStringChannel,
	}))
	setupRoutes(app)

	// Set up a channel to capture interrupt signals
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	serverShutdown := make(chan struct{})

	// Start a separate goroutine that will listen for the interrupt signal
	go func() {
		<-c //wait for interrupt signal
		fmt.Println("Gracefully shutting down, Running cleanup tasks..")
		_ = app.Shutdown()
		close(requestChan) // Close the request channel after graceful shutdown signal is received
		utils.CloseLogStringChannel()
		// Allow some time to process pending requests before the application exits
		time.Sleep(5 * time.Second)
		serverShutdown <- struct{}{}
		database.RedisClose()
		database.DBClose()
		fmt.Println("Shutdown complete, bye.")
	}()

	if err := app.Listen(":" + config.EnvDBURI("HTTP_LISTEN_PORT")); err != nil {
		panic(err)
	}
	fmt.Println("server has been started")
	<-serverShutdown
}
