package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"sync"
	"time"

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

	runAll := func(cmd string, w io.WriteCloser) {
		l.Lock()
		if len(players) == 0 {
			go func() {
				w.Write([]byte("No Client Connected"))
				w.Close()
			}()
		} else {
			for _, p := range players {
				t := time.NewTimer(3 * time.Second)
				has_run := false
				go func() {
					<-t.C
					if !has_run {
						w.Write([]byte("Timed Out"))
						w.Close()
						has_run = true
					}
				}()

				p.Exec(cmd, func(data map[string]any) {
					logrus.Infof("Command %s", cmd)
					if message, ok := data["statusMessage"]; ok {
						fmt.Printf("Result:\n%s\n\n", message)
					} else {
						fmt.Printf("Result:\n%#+v\n\n", data)
					}
					if !has_run {
						body, _ := json.MarshalIndent(data, "", "\t")
						w.Write(body)
						w.Close()
						has_run = true
					}
				})
			}
		}
		l.Unlock()
	}

	app.Post("/Exec", func(c *fiber.Ctx) error {
		cmd := string(c.Body())
		if len(cmd) == 0 {
			return c.Status(400).SendString("No Command Specified")
		}
		r, w := io.Pipe()
		runAll(cmd, w)

		return c.Status(200).SendStream(r)
	})

	app.Get("/Exec/+", func(c *fiber.Ctx) error {
		cmd, err := url.QueryUnescape(c.Params("+"))
		if err != nil {
			return err
		}

		r, w := io.Pipe()
		runAll(cmd, w)

		return c.Status(200).SendStream(r)
	})

	app.Use("/", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return websocket.New(func(c *websocket.Conn) {
				server.HandleConnection(c)
			})(c)
		}
		return c.Next()
	})

	err := app.Listen(":8080")
	if err != nil {
		logrus.Fatal(err)
	}
}
