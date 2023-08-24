package main

import (
	"game/gamehandler"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"github.com/AspieSoft/goutil/v5"
	"gopkg.in/yaml.v3"
)

func main(){
	// set default config
	maxFPS := uint16(120)
	inconsistentRand := true
	objectTypes := []string{"object"}

	// get game config file
	if gameConfigFile, err := os.ReadFile("./src/config.yml"); err == nil {
		gameConfig := map[string]interface{}{}
		if err := yaml.Unmarshal(gameConfigFile, &gameConfig); err == nil {
			if val, ok := gameConfig["MaxFPS"]; ok {
				if v := goutil.ToType[uint16](val); v != 0 {
					maxFPS = v
				}
			}

			if val, ok := gameConfig["MaxFPS"]; ok {
				inconsistentRand = goutil.ToType[bool](val)
			}
			
			if val, ok := gameConfig["ObjectTypes"]; ok {
				if v := goutil.ToType[[]string](val); len(v) != 0 {
					objectTypes = v
				}
			}
		}
	}


	// create app and window
	a := app.New()
	defer a.Quit()

	w := a.NewWindow("Test")
	defer w.Close()

	// get background image
	var img *canvas.Image
	if bg, err := os.ReadFile("./assets/background.jpg"); err == nil {
		img = canvas.NewImageFromResource(fyne.NewStaticResource("background", bg))
	}else if bg, err := os.ReadFile("./assets/background.png"); err == nil {
		img = canvas.NewImageFromResource(fyne.NewStaticResource("background", bg))
	}

	canvasList := map[string]*fyne.Container{}
	canvasListKeys := []string{}
	canvasListArr := []fyne.CanvasObject{}
	for _, objType := range objectTypes {
		canvasList[objType] = container.NewWithoutLayout()
		canvasListKeys = append(canvasListKeys, objType)
		canvasListArr = append(canvasListArr, canvasList[objType])
	}
	canvasBox := container.NewWithoutLayout(canvasListArr...)

	var box *fyne.Container
	if img != nil {
		box = container.NewMax(img, canvasBox)
	}else{
		box = container.NewMax(canvasBox)
	}
	w.SetContent(box)

	w.Resize(fyne.NewSize(720, 480))
	w.CenterOnScreen()
	w.SetFixedSize(false)
	w.SetPadded(false)
	// w.SetFullScreen(true)
	w.SetMaster()

	if icon, err := os.ReadFile("./assets/icon.png"); err == nil {
		res := fyne.NewStaticResource("icon", icon)
		a.SetIcon(res)
		w.SetIcon(res)
	}

	canvasWidth := canvasBox.Size().Width
	canvasHeight := canvasBox.Size().Height
	var canvasScale float32
	if canvasWidth > canvasHeight {
		canvasScale = canvasHeight
	}else{
		canvasScale = canvasWidth
	}
	canvasScale /= 100

	gameData := gamehandler.Game{
		Canvas: canvasBox,
		CanvasList: canvasList,
		CanvasListKeys: canvasListKeys,
		Window: w,

		Size: gamehandler.CanvasSize{
			RealWidth: canvasWidth + 0.000025,
			RealHeight: canvasHeight + 0.000025,
			Scale: canvasScale,

			Width: (canvasWidth + 0.000025) / canvasScale / 2,
			Height: (canvasHeight + 0.000025) / canvasScale / 2,
		},

		MaxFPS: maxFPS,
		InconsistentRand: inconsistentRand,
	}

	go func(){
		for {
			time.Sleep(300 * time.Millisecond)

			canvasWidth := canvasBox.Size().Width
			canvasHeight := canvasBox.Size().Height
			var canvasScale float32
			if canvasWidth > canvasHeight {
				canvasScale = canvasHeight
			}else{
				canvasScale = canvasWidth
			}

			canvasScale /= 100

			gameData.MU.Lock()
			gameData.Size = gamehandler.CanvasSize{
				RealWidth: canvasWidth + 0.000025,
				RealHeight: canvasHeight + 0.000025,
				Scale: canvasScale,

				Width: (canvasWidth + 0.000025) / canvasScale / 2,
				Height: (canvasHeight + 0.000025) / canvasScale / 2,
			}
			gameData.MU.Unlock()
		}
	}()

	go Init(&gameData)

	w.ShowAndRun()
}

// GameLoop creates a new game loop
//
// call this with `go GameLoop(...)` to run this on a new thread/goroutine
//
// this method has been modified from the notch game loop
func GameLoop(gameData *gamehandler.Game, fps uint16, cb func(game *gamehandler.Game, thread *gamehandler.ThreadInfo)){
	time.Sleep(100 * time.Millisecond)

	speedDelta := float32(1)
	if fps > gameData.MaxFPS {
		speedDelta = float32(fps) / float32(gameData.MaxFPS)
		fps = gameData.MaxFPS
	}

	lastTime := float64(time.Now().UnixNano())
	ns := float64(time.Second.Nanoseconds()) / float64(fps)
	delta := float64(0)
	frames := uint16(0)
	timeMS := float64(time.Now().UnixMilli())

	currentFPS := fps

	for {
		now := float64(time.Now().UnixNano())
		delta += (now - lastTime) / ns
		lastTime = now

		if delta >= 1 {
			gameData.MU.Lock()
			cb(gameData, &gamehandler.ThreadInfo{
				FPS: currentFPS,
				Frame: frames,
				SpeedDelta: speedDelta,
			})

			gameData.MU.Unlock()

			frames++
			delta--
			if float64(time.Now().UnixMilli()) - timeMS >= float64(time.Second.Milliseconds()) {
				currentFPS = frames
				timeMS += float64(time.Second.Milliseconds())
				frames = 0
			}
		}
	}
}
