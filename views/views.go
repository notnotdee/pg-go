package views

import (
	"github.com/gofiber/fiber"
)

func Default(ctx *fiber.Ctx) {
	ctx.Send("hello world")
}