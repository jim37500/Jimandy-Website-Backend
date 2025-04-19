package data

import "github.com/gofiber/fiber/v2"

type AccessControlList struct {
	Path        string
	Method      string
	Handler     fiber.Handler
	Authorization string
}
