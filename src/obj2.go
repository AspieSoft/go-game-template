package game

import (
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
		rect := canvas.NewRectangle(color.RGBA{35, 190, 15, 255})
		
		object := game.Add("object", "obj2", 30, 20, 4, 4, func(game *gamehandler.Game) fyne.CanvasObject {
			// rect := canvas.NewRectangle(rectColor)

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

			if player != nil && object.IsColideing(player) {
				rect.FillColor = color.RGBA{255, 0, 0, 255}
			}else if len(object.IsColideingType("object")) != 0 {
				rect.FillColor = color.RGBA{50, 90, 200, 255}
			}else{
				rect.FillColor = color.RGBA{35, 190, 15, 255}
			}
		}
	})
}
