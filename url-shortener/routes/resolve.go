package routes

import (
	"fiber-url-shortner/config"
	"fiber-url-shortner/database"

	"github.com/gofiber/fiber/v2"
)

func Resolve(ctx *fiber.Ctx) error {
	url := ctx.Params("url")

	value, err := database.GetHashValue(config.SHORT_URL_KEY, url)
	if value != "" && err == nil {
		ctx.Locals("msg", "Redirected succesfully "+url+" "+value)
		ctx.Next()
		return ctx.Redirect(value, 301)
	}

	if res, err := database.GetURL(url); err != nil {
		ctx.Locals("msg", "Invalid URL "+url)
		ctx.Next()
		return ctx.Status(404).JSON(fiber.Map{"message": "Invalid URL"})
	} else if result := database.SetHashkey(config.SHORT_URL_KEY, url, res); result != nil {
		ctx.Locals("msg", "sethashkey "+result.Error())
		ctx.Next()
		return ctx.Status(253).JSON(fiber.Map{"error": "sethashkey " + result.Error()})
	} else {
		ctx.Locals("msg", "Redirected succesfully "+url+" "+res)
		ctx.Next()
		return ctx.Redirect(res, 301)
	}
}
