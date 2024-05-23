package routes

import (
	// "fiber-url-shortner/database"
	"fiber-url-shortner/config"
	"fiber-url-shortner/database"
	"fiber-url-shortner/helpers"
	"log"
	"math/rand"
	"time"

	// "strconv"
	// "time"

	"github.com/gofiber/fiber/v2"
)

type request struct {
	URL string `json:"url"`
}

type response struct {
	ShortURL string `json:"shortURL"`
}

var randSource = rand.NewSource(time.Now().UnixNano())

func Shorten(ctx *fiber.Ctx) error {
	body := request{}

	if err := ctx.BodyParser(&body); err != nil {
		ctx.Locals("msg", "cannot parse JSON "+string(ctx.Body()))
		ctx.Next()
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "cannot parse JSON"})
	}

	// check for if it is domain error
	if helpers.DomainError(body.URL) {
		ctx.Locals("msg", "Cannot shorten same or empty domain")
		ctx.Next()
		return ctx.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"msg": "Cannot shorten same/empty domain"})
	}

	r := rand.New(randSource)
	id := helpers.Base62Encoder(uint64(r.Intn(1000000000)))
	for database.CheckIDExists(id) {
		log.Fatalf("RANDOM ID exists, golang rand error for %s, regenerating", id)
		id = helpers.Base62Encoder(uint64(r.Intn(1000000000)))
	}

	if resID, err := database.InsertData(id, body.URL); err != nil {
		ctx.Locals("msg", "Data Insertion error "+err.Error())
		ctx.Next()
		return ctx.Status(252).JSON(fiber.Map{"msg": "Data Insertion error " + err.Error()})
	} else {
		resp := response{
			ShortURL: config.EnvDBURI("DOMAIN_HTTP_PROTOCAL") + config.EnvDBURI("URL_SHORTENER_DOMAIN") + "/u/" + resID,
		}
		ctx.Locals("msg", "Request Succesfull "+id+" "+body.URL)
		ctx.Next()
		return ctx.Status(fiber.StatusOK).JSON(resp)
	}
}
