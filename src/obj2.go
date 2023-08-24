package game

import (
	"fmt"
	"game/BorderMethod"
	"game/CollisionMethod"
	"game/gamehandler"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

func init(){
	gamehandler.InitObject(func(game *gamehandler.Game) {
		object := game.Add("object", "obj2", -10, -10, 4, 4, func(game *gamehandler.Game) fyne.CanvasObject {
			rect := canvas.NewRectangle(color.RGBA{35, 190, 15, 255})

			return container.NewMax(rect)
		})

		object.PreferredFPS = 30

		// object.VelX = 3
		// object.VelY = 3

		object.BorderMethod = BorderMethod.Bounce
		object.CollisionMethod = CollisionMethod.Box

		/* object.UpdateBasic = func(game *gamehandler.Game, thread *gamehandler.ThreadInfo) {
			// fmt.Println(object.X, object.Y)
		} */

		var player *gamehandler.GameObject

		object.UpdateBasic = func(game *gamehandler.Game, thread *gamehandler.ThreadInfo) {
			if player == nil {
				playerList := game.Get("player", "player")
				if len(playerList) != 0 {
					player = playerList[0]
				}
			}

			if player != nil {
				fmt.Println("colliding:", object.IsColideing(player))
				// object.IsColideing(player)
			}
		}
	})
}
