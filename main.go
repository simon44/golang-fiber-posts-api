package main

import (
	"strconv"
	"time"
	"github.com/gofiber/fiber/v2"
	"github.com/go-playground/validator/v10"
)

type Post struct {
	Id int `json:"id"`
	Content string `json:"content"`
	CreatedAt time.Time `json:"createAt"`
}

var posts = []Post{
	{Id: 1, Content: "My first post", CreatedAt: time.Now()},
	{Id: 2, Content: "My second post", CreatedAt: time.Now()},
}

func main() {
    app := fiber.New()

    app.Get("/", func(c *fiber.Ctx) error {
        return c.SendString("Hello, World 👋!")
    })

		app.Get("/posts", GetPosts)
		app.Post("/posts", CreatePost)
		app.Get("/posts/:id", GetPost)

    app.Listen(":3000")
}

func GetPosts(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(posts)
}

func GetPost(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": "Cannot parse id",
		})
	}
	for _, p := range posts {
		if p.Id == id {
			return c.Status(fiber.StatusOK).JSON(p)
		}
	}
	return c.Status(fiber.StatusNotFound).JSON(&fiber.Map{
		"error": "Post not found",
	})
}

func CreatePost(c *fiber.Ctx) error {
	type request struct {
		Content string `json:"content" validate:"required,min=5"`
	}

	r := new(request)
	if err := c.BodyParser(r); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": "Cannot parse body",
		})	
	}

	err := validator.New().Struct(r)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
				"error": "Field is not valid " + err.StructNamespace() + " : " + err.Tag() + " : " + err.Param(),
			})
		}
	}

	nextId := posts[len(posts)-1].Id+1
	post := Post{Id: nextId, Content: r.Content, CreatedAt: time.Now()}

	posts = append(posts, post)

	return c.Status(fiber.StatusOK).JSON(post)
}