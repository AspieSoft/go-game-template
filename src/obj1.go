package game

import (
	"game/BorderMethod"
	"game/CollisionMethod"
	"game/gamehandler"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

func init(){
	gamehandler.InitObject(func(game *gamehandler.Game) {
		size := float32(4)

		r := GameRandSeed.Get(0, 12)
		border := uint8(0)
		if r % 2 == 0 {
			border = 2
		}
		if r % 3 == 0 || r % 4 == 0 {
			border++
		}

		x := float32(0)
		y := float32(0)
		if border == 0 {
			y = -game.Size.Height - size
		}else if border == 1 {
			x = game.Size.Width + size
		}else if border == 2 {
			y = game.Size.Height + size
		}else if border == 3 {
			x = -game.Size.Width - size
		}

		if x == 0 {
			x = float32(GameRandSeed.Get(int(-game.Size.Width + size), int(game.Size.Width - size)))
		}else if y == 0 {
			y = float32(GameRandSeed.Get(int(-game.Size.Height + size), int(game.Size.Height - size)))
		}

		object := game.Add("object", "obj1", x, y, size, size, func(game *gamehandler.Game) fyne.CanvasObject {
			res := canvas.NewRectangle(color.RGBA{255, 0, 0, 255})
			return res
		})

		object.PreferredFPS = 30

		object.VelX = 3
		object.VelY = 3

		object.BorderMethod = BorderMethod.Bounce
		object.CollisionMethod = CollisionMethod.Box

		/* object.UpdateBasic = func(game *gamehandler.Game, thread *gamehandler.ThreadInfo) {
			// fmt.Println(object.X, object.Y)
		} */
	})
}
