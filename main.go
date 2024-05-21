package main

import (
	"log"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

type Chatroom struct {
	CID   string
	Users map[string]*websocket.Conn
}

type Message struct {
	UserId     string `json:"user_id"`
	ChatroomId string `json:"cid"`
	Content    string `json:"content"`
}

type Request struct {
	UserId     string `json:"user_id"`
	ChatroomId string `json:"cid"`
}

type Response struct {
	Data string
}

func main() {
	app := fiber.New()

	hub := make(map[string]Chatroom)

	app.Use("/ws", func(c *fiber.Ctx) error {

		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws/:guid/:cid", websocket.New(func(c *websocket.Conn) {
		guid := c.Params("guid")
		cid := c.Params("cid")

		hub[cid].Users[guid] = c

		log.Println("hub", hub)

		for {
			clientMessage := &Message{}
			if err := c.ReadJSON(&clientMessage); err != nil {
				log.Println("read:", err)
				delete(hub[cid].Users, guid)
				return
			}

			log.Println(clientMessage.Content)

			chatroom := clientMessage.ChatroomId
			chatroomUsers := hub[chatroom]
			log.Println("chatroom Users", chatroomUsers)
			for userId, connection := range chatroomUsers.Users {
				if userId != clientMessage.UserId {
					if err := connection.WriteJSON(clientMessage); err != nil {
						log.Println(err)
						break
					}
				}
			}

		}

	}))

	app.Post("/create", func(c *fiber.Ctx) error {
		var request Request

		if err := c.BodyParser(&request); err != nil {
			log.Println(err)
			return c.Status(400).JSON(Response{
				Data: "Bad request",
			})
		}

		if _, ok := hub[request.ChatroomId]; !ok {
			hub[request.ChatroomId] = Chatroom{
				CID:   request.ChatroomId,
				Users: make(map[string]*websocket.Conn),
			}
			return c.Status(200).JSON(Response{
				Data: "Chatroom created",
			})
		}

		return c.Status(400).JSON(Response{
			Data: "Chatroom already exists",
		})

	})

	log.Fatal(app.Listen(":3000"))

}

func ListenUser(c *websocket.Conn) {

}
