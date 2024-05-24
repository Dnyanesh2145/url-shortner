package routes

import (
	"fiber-url-shortner/database"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func Health(c *fiber.Ctx) error {
	strStatus := "SDB,MDB,REDIS: OK"
	if _, err := database.RedisPing(); err != nil {
		strStatus = "Redis: " + err.Error()
		fmt.Println("HealthCheck failed due to cache/redis server connectivity: " + err.Error())
	}
	if _, err := database.DBPing(); err != nil {
		strStatus = "DB: " + err.Error()
		fmt.Println("HealthCheck failed due to cache/redis server connectivity: " + err.Error())
	}
	if strStatus == "SDB,MDB,REDIS: OK" {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"msg": strStatus})
	}
	return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"msg": strStatus})
}
