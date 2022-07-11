package websocket

import (
	"fmt"
	"net/http"
	"order-system/utils"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func CreateWebsocketHandler(hub *Hub) func(echo.Context) error {
	return func(c echo.Context) error {
		user, err := utils.VerifyAndParseJWT(c.QueryParam("jwt"))

		if err != nil {
			return echo.ErrUnauthorized
		}

		ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			return err
		}
		hub.RegisterNewClient(user.ID, ws)
		fmt.Println("registered")
		return nil
	}
}
