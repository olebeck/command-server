package main

import (
	"sync"

	"github.com/sandertv/mcwss"
	"github.com/sandertv/mcwss/protocol/event"
	"github.com/sirupsen/logrus"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func main() {
	app := fiber.New()
	server := mcwss.NewServer(nil)
	l := sync.Mutex{}
	players := make(map[string]*mcwss.Player)

	server.OnConnection(func(player *mcwss.Player) {
		l.Lock()
		players[player.Name()] = player
		l.Unlock()
		log := logrus.WithField("Player", player.Name())
		log.Info("Connected")
		player.SendMessage("Connected to Command Server")
		player.OnScriptLoaded(func(event *event.ScriptLoaded) {
			log.Infof("Loaded Script: %s", event.ScriptName)
		})
	})

	server.OnDisconnection(func(player *mcwss.Player) {
		log := logrus.WithField("Player", player.Name())
		l.Lock()
		delete(players, player.Name())
		l.Unlock()
		log.Info("Disconnected")
	})

	app.Post("/Exec", func(c *fiber.Ctx) error {
		cmd := string(c.Body())
		if len(cmd) == 0 {
			return c.Status(400).SendString("No Command Specified")
		}

		l.Lock()
		for _, p := range players {
			p.Exec(cmd, func(data map[string]any) {
				logrus.Debugf("Command %s\n%+#v", cmd, data)
			})
		}
		l.Unlock()

		return c.SendStatus(200)
	})

	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		server.HandleConnection(c)
	}))

	app.Listen(":8080")
}
