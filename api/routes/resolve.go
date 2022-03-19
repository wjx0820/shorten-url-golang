package routes

import (
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/wjx0820/shorten-url-golang/database"
)

func ResolveURL(c *fiber.Ctx) error {
	url := c.Params("url")

	r := database.CreateClient(0)
	defer r.Close()

	val, err := r.Get(database.Ctx, url).Result()
	if err == redis.Nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "short not found in the db"})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot connect to db"})
	}

	rInr := database.CreateClient(1)
	defer rInr.Close()

	_ = rInr.Incr(database.Ctx, "counter")

	return c.Redirect(val, 301)
}
