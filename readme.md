# Go Game Template

[![donation link](https://img.shields.io/badge/buy%20me%20a%20coffee-paypal-blue)](https://paypal.me/shaynejrtaylor?country.x=US&locale.x=en_US)

A multi threaded game template for modifying to your needs.

## Note: This module is still in Alpha, and it may likely see many breaking changes until it is released under v1.0

## Installation

```shell script
git clone https://github.com/AspieSoft/go-game-template
```

## Setup

Most of the files you will need to configure are in the src directory.

The init.go file is called with its `Init` method once the game starts (a capital `Init`, not the default lowercase `init`).

### Creating New Game Objects

```go

import (
  "game/BorderMethod"
  "game/gamehandler"
  "image/color"

  "fyne.io/fyne/v2"
  "fyne.io/fyne/v2/canvas"
)

func init(){
  gamehandler.InitObject(func(game *gamehandler.Game) {
    object := game.Add("object", "MyObject", -game.Size.Width - 5, -game.Size.Height - 5, 5, 5, func(game *gamehandler.Game) fyne.CanvasObject {
      res := canvas.NewRectangle(color.RGBA{255, 0, 0, 255})
      return res
    })

    // by default, border detection and other builtin features run on the 120 fps draw method
    // setting an objects Preferred FPS allows you to override this to use a slower thread if needed
    // recommended: use 60 fps for a player, to prevent input lag (30 fps works better for entities)
    object.PreferredFPS = 30

    object.VelX = 3
    object.VelY = 3

    // there are optional builtin methods for handling basic common tasks
    object.BorderMethod = BorderMethod.Bounce

    object.Update := func(game *gamehandler.Game, thread *gamehandler.ThreadInfo) {
      // run updates on a 60 fps thread
      // useful for most common game updates, or player related updates with less lag
    }

    object.UpdateSlow := func(game *gamehandler.Game, thread *gamehandler.ThreadInfo) {
      // run updates on a slower 30 fps thread
      // useful for a large number of extra objects to share a slower thread
    }

    object.UpdateBasic := func(game *gamehandler.Game, thread *gamehandler.ThreadInfo) {
      // run updates on a slow, cpu saving 15 fps thread
      // useful for larger math operations
    }

    object.Draw := func(game *gamehandler.Game, thread *gamehandler.ThreadInfo) {
      // run updates on a fast 120 fps thread
      // useful for graphics
    }
  }
}

```
