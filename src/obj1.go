package game

import (
	"game/BorderMethod"
	"game/gamehandler"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

func init(){
	gamehandler.InitObject(func(game *gamehandler.Game) {
		object := game.Add("object", "obj1", -game.Size.Width - 5, -game.Size.Height - 5, 5, 5, func(game *gamehandler.Game) fyne.CanvasObject {
			res := canvas.NewRectangle(color.RGBA{255, 0, 0, 255})
			return res
		})

		object.PreferredFPS = 30

		object.VelX = 3
		object.VelY = 3

		object.BorderMethod = BorderMethod.Bounce

		/* object.UpdateBasic = func(game *gamehandler.Game, thread *gamehandler.ThreadInfo) {
			// fmt.Println(object.X, object.Y)
		} */
	})
}
