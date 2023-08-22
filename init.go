package main

import (
	"game/gamehandler"
	game "game/src"
	"time"
)

func Init(gameData *gamehandler.Game){
	time.Sleep(1 * time.Second)

	game.Init(gameData)

	go GameLoop(gameData, 60, gamehandler.Update)
	go GameLoop(gameData, 120, gamehandler.Draw)

	go GameLoop(gameData, 15, gamehandler.UpdateBasic)
	go GameLoop(gameData, 30, gamehandler.UpdateSlow)

	for _, objInit := range gamehandler.GameObjectInit {
		objInit(gameData)
	}
}
