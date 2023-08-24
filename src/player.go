package game

import (
	"game/enum/BorderMethod"
	"game/enum/CollisionMethod"
	"game/gamehandler"
	"image/color"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
)


func init(){
	gamehandler.InitObject(func(game *gamehandler.Game) {
		// (game.Size.Width / game.Size.Scale / 2), (game.Size.Height / game.Size.Scale / 2)
		object := game.Add("player", "player", 0, 0, 5, 5, func(game *gamehandler.Game) fyne.CanvasObject {
			if icon, err := os.ReadFile("./assets/objects/player/white.png"); err == nil {
				res := canvas.NewImageFromResource(fyne.NewStaticResource("icon", icon))
				return res
			}

			res := canvas.NewCircle(color.White)
			return res
		})

		object.PreferredFPS = 60

		object.BorderMethod = BorderMethod.PushLimit
		object.CollisionMethod = CollisionMethod.Radius

		speed := float32(4)
		moveX := int8(0)
		moveY := int8(0)

		object.Update = func(game *gamehandler.Game, thread *gamehandler.ThreadInfo) {
			if moveX < 0 {
				object.VelX = -speed
			}else if moveX > 0 {
				object.VelX = speed
			}else{
				object.VelX = 0
			}

			if moveY < 0 {
				object.VelY = -speed
			}else if moveY > 0 {
				object.VelY = speed
			}else{
				object.VelY = 0
			}
		}

		if fyne.CurrentDevice().HasKeyboard() {
			if deskCanvas, ok := game.Window.Canvas().(desktop.Canvas); ok {
				deskCanvas.SetOnKeyDown(func(key *fyne.KeyEvent) {
					if key.Name == fyne.KeyW && object.VelY >= 0 {
						moveY--
					}else if key.Name == fyne.KeyS && object.VelY <= 0 {
						moveY++
					}else if key.Name == fyne.KeyA && object.VelX >= 0 {
						moveX--
					}else if key.Name == fyne.KeyD && object.VelX <= 0 {
						moveX++
					}
				})
				deskCanvas.SetOnKeyUp(func(key *fyne.KeyEvent) {
					if key.Name == fyne.KeyW && object.VelY <= 0 {
						moveY++
					}else if key.Name == fyne.KeyS && object.VelY >= 0 {
						moveY--
					}else if key.Name == fyne.KeyA && object.VelX <= 0 {
						moveX++
					}else if key.Name == fyne.KeyD && object.VelX >= 0 {
						moveX--
					}
				})
			}
		}

		//todo: add mobile key method
	})
}
