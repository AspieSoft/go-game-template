package gamehandler

import (
	"game/enum/BorderMethod"
	"game/enum/CollisionMethod"
	"game/enum/TypeCollisionMethod"
	"math"
	"sync"

	"fyne.io/fyne/v2"
	"github.com/AspieSoft/goutil/v5"
)

type CanvasSize struct {
	// RealWidth is the actual window width in pixels
	RealWidth float32

	// RealHeight is the actual window height in pixels
	RealHeight float32

	// Scale is a number based on the smallest dimension used to calculate how objects should be resized based on the screen size
	Scale float32

	// Width is calculated after scaling down
	//
	// for the actual pixel width, use 'RealWidth'
	Width float32

	// Height is calculated after scaling down
	//
	// for the actual pixel height, use 'RealHeight'
	Height float32
}

type ThreadInfo struct {
	FPS uint16
	Frame uint16

	// SpeedDelta can be multiplied by a consistently updating number that needs to recalculate if the max FPS is decreased in config.yml
	SpeedDelta float32
}

type Game struct {
	Canvas *fyne.Container
	CanvasList map[string]*fyne.Container
	CanvasListKeys []string

	Window fyne.Window
	Size CanvasSize

	MaxFPS uint16
	InconsistentRand bool

	MU sync.Mutex
}


var GameObjectInit []func(game *Game) = []func(game *Game){}

func InitObject(cb func(game *Game)){
	GameObjectInit = append(GameObjectInit, cb)
}

type GameObject struct{
	id string
	objType string
	name string
	Object fyne.CanvasObject
	MU sync.Mutex

	X float32
	Y float32

	Width float32
	Height float32

	VelX float32
	VelY float32

	// OnBorderX returns -1 or 1 if this object if touching a border in the x axis
	//
	// it will also return -3 or 3 if it is past the border and fully hidden
	//
	// -2 or 2 means we are pushing on a border (moving in that direction)
	//
	// -4 or 4 means we are pushing past a border (moving in that direction)
	OnBorderX int8

	// OnBorderY returns -1 or 1 if this object if touching a border in the y axis
	//
	// it will also return -3 or 3 if it is past the border and fully hidden
	//
	// -2 or 2 means we are pushing on a border (moving in that direction)
	//
	// -4 or 4 means we are pushing past a border (moving in that direction)
	OnBorderY int8

	// BorderMethod allows you to use a preset method on how to handle objects when they touch a border
	//
	// default: Ignore
	BorderMethod uint8

	// CollisionMethod determines how an objects collision hitbox should be calculated
	//
	// default: Ghost
	CollisionMethod uint8

	// Store is a basic map for storing extra data attached to an object if needed
	Store map[string]any

	// PreferredFPS is an optional FPS preference for detection updates for an object
	//
	// example: border detection
	//
	// default: 120
	PreferredFPS uint16

	// Update is an optional method that runs on a normal 60 fps game loop
	//
	// example: updating the player stats
	Update func(game *Game, thread *ThreadInfo)

	// Draw is an optional method that runs on a graphical 120 fps game loop
	//
	// example: drawing data to canvas
	Draw func(game *Game, thread *ThreadInfo)

	// UpdateBasic is an optional method that runs on an extra slow 15 fps game loop
	//
	// this method can be useful for long math operations or things that do not take priority
	//
	// example: updating particals (there could be many of them taking up a lot of the cpu)
	UpdateBasic func(game *Game, thread *ThreadInfo)

	// UpdateSlow is an optional method that runs on a slow 30 fps game loop
	//
	// this methid can be useful when many objects need to update frequently, and the cpu usage needs to be reduced
	//
	// example: updating entity stats (seperating this from the player can prevent input lag)
	UpdateSlow func(game *Game, thread *ThreadInfo)
}

type Direction struct {
	dist float32
	dirX float32
	dirY float32
}

var gameObjects map[string][]*GameObject = map[string][]*GameObject{}
var gameObjectsMU sync.Mutex

var gameColissionType map[string]uint8 = map[string]uint8{}

// game methods

// Add adds a new object to the game
func (game *Game) Add(objType string, name string, x, y, width, height float32, cb func(game *Game) fyne.CanvasObject) *GameObject {
	object := GameObject{
		id: string(goutil.Crypt.RandBytes(64)),
		objType: objType,
		name: name,
		Object: cb(game),

		X: x,
		Y: y,
		Width: width,
		Height: height,
	}

	gameObjectsMU.Lock()
	if _, ok := gameObjects[objType]; !ok {
		gameObjects[objType] = []*GameObject{}
	}
	gameObjects[objType] = append(gameObjects[objType], &object)
	if _, ok := game.CanvasList[objType]; ok {
		game.CanvasList[objType].Add(object.Object)
		game.CanvasList[objType].Refresh()
	}
	gameObjectsMU.Unlock()

	return &object
}

// RemoveType clears all objects of the same type
func (game *Game) RemoveType(objType string){
	gameObjectsMU.Lock()
	defer gameObjectsMU.Unlock()

	if _, ok := gameObjects[objType]; !ok {
		return
	}

	gameObjects[objType] = []*GameObject{}
	game.CanvasList[objType].RemoveAll()
	game.CanvasList[objType].Refresh()
}

// Get returns a list of objects by type and name
func (game *Game) Get(objType string, name string) []*GameObject {
	gameObjectsMU.Lock()
	defer gameObjectsMU.Unlock()

	if _, ok := gameObjects[objType]; !ok {
		return []*GameObject{}
	}

	list := []*GameObject{}

	for _, object := range gameObjects[objType] {
		if object.name == name {
			list = append(list, object)
		}
	}

	return list
}

// GetType returns a list of objects by type
func (game *Game) GetType(objType string) []*GameObject {
	gameObjectsMU.Lock()
	defer gameObjectsMU.Unlock()

	if _, ok := gameObjects[objType]; !ok {
		return []*GameObject{}
	}

	return gameObjects[objType]
}

// GetID returns an object by its ID
func (game *Game) GetID(objType string, id string) *GameObject {
	gameObjectsMU.Lock()
	defer gameObjectsMU.Unlock()

	if _, ok := gameObjects[objType]; !ok {
		return nil
	}

	for _, object := range gameObjects[objType] {
		if object.id == id {
			return object
		}
	}

	return nil
}

func (game *Game) eachObject(cb func(object *GameObject)){
	for i := 0; i < len(game.CanvasListKeys); i++ {
		for _, object := range gameObjects[game.CanvasListKeys[i]] {
			cb(object)
		}
	}
}

// SetTypeCollision sets a specific collision type to an object
//
// default: Any
//
// example: you can set a GUI to type ghost if objects should not interact with it
func (game *Game) SetTypeCollision(objType string, typeCollisionMethod uint8){
	if typeCollisionMethod > 3 {
		typeCollisionMethod = 0
	}

	gameObjectsMU.Lock()
	gameColissionType[objType] = typeCollisionMethod
	gameObjectsMU.Unlock()
}


// object methods

// Remove removes this object from the game
func (object *GameObject) Remove(game *Game, thread *ThreadInfo){
	gameObjectsMU.Lock()
	defer gameObjectsMU.Unlock()

	game.Canvas.Remove(object.Object)
	for i, obj := range gameObjects[object.objType] {
		if obj.id == object.id {
			gameObjects[object.objType] = append(gameObjects[object.objType][:i], gameObjects[object.objType][i+1:]...)
		}
	}
}

// handleBorder handles border detection math
//
// this method will be called by the different update methods depending on an objects PreferredFPS
func (object *GameObject) handleBorder(game *Game, thread *ThreadInfo){
	{ // check if object in on or past border
		if object.X + object.Width < -game.Size.Width {
			object.OnBorderX = -3
		} else if object.X - object.Width > game.Size.Width {
			object.OnBorderX = 3
		}else if object.X - object.Width <= -game.Size.Width {
			object.OnBorderX = -1
		}else if object.X + object.Width >= game.Size.Width {
			object.OnBorderX = 1
		}else{
			object.OnBorderX = 0
		}
	
		if object.Y + object.Height < -game.Size.Height {
			object.OnBorderY = -3
		} else if object.Y - object.Height > game.Size.Height {
			object.OnBorderY = 3
		}else if object.Y - object.Height <= -game.Size.Height {
			object.OnBorderY = -1
		}else if object.Y + object.Height >= game.Size.Height {
			object.OnBorderY = 1
		}else{
			object.OnBorderY = 0
		}
	
		if object.VelX < 0 && object.OnBorderX <= -1 {
			object.OnBorderX--
		}else if object.VelX > 0 && object.OnBorderX >= 1 {
			object.OnBorderX++
		}

		if object.VelY < 0 && object.OnBorderY <= -1 {
			object.OnBorderY--
		}else if object.VelY > 0 && object.OnBorderY >= 1 {
			object.OnBorderY++
		}
	}

	// handle object border method
	switch object.BorderMethod {
	case BorderMethod.PushLimit:
		if object.X - object.Width < -game.Size.Width {
			object.X = -game.Size.Width + object.Width
		}else if object.X + object.Width > game.Size.Width {
			object.X = game.Size.Width - object.Width
		}

		if object.Y - object.Height < -game.Size.Height {
			object.Y = -game.Size.Height + object.Height
		}else if object.Y + object.Height > game.Size.Height {
			object.Y = game.Size.Height - object.Height
		}

	case BorderMethod.PushHide:
		if object.X + object.Width < -game.Size.Width {
			object.X = -game.Size.Width - object.Width - 0.25
		}else if object.X - object.Width > game.Size.Width {
			object.X = game.Size.Width + object.Width + 0.25
		}

		if object.Y + object.Height < -game.Size.Height {
			object.Y = -game.Size.Height - object.Height - 0.25
		}else if object.Y - object.Height > game.Size.Height {
			object.Y = game.Size.Height + object.Height + 0.25
		}

	case BorderMethod.Bounce:
		if object.OnBorderX != 0 && object.OnBorderX % 2 == 0 {
			object.VelX *= -1
		}
		if object.OnBorderY != 0 && object.OnBorderY % 2 == 0 {
			object.VelY *= -1
		}

	case BorderMethod.Teleport:
		if object.OnBorderX <= -4 {
			object.X = game.Size.Width + object.Width
		}else if object.OnBorderX >= 4 {
			object.X = -game.Size.Width - object.Width
		}

		if object.OnBorderY <= -4 {
			object.Y = game.Size.Height + object.Height
		}else if object.OnBorderY >= 4 {
			object.Y = -game.Size.Height - object.Height
		}

	case BorderMethod.RemoveObject:
		if object.OnBorderX <= -4 || object.OnBorderX >= 4 || object.OnBorderY <= -4 || object.OnBorderY >= 4 {
			object.Remove(game, thread)
		}
	}
}


// GetDistance calculates the distance between 2 objects
func (obj1 *GameObject) GetDistance(obj2 *GameObject) float32 {
	diffX := obj1.X - obj2.X
	diffY := obj1.Y - obj2.Y
	return float32(math.Sqrt(math.Pow(float64(diffX), 2) + math.Pow(float64(diffY), 2)))
}

// GetDirection is similar to the 'GetDistance' method,
// but this method aso returns the direction of an object and the difference between the x and y distance
//
// example use: velX += Direction.dirX * speed; velY += Direction.dirY * speed (to move twards an object)
func (obj1 *GameObject) GetDirection(obj2 *GameObject) Direction {
	diffX := obj1.X - obj2.X
	diffY := obj1.Y - obj2.Y
	dist := float32(math.Sqrt(math.Pow(float64(diffX), 2) + math.Pow(float64(diffY), 2)))

	dirX := (1/dist * diffX)
	dirY := (1/dist * diffY)

	return Direction{
		dist,
		dirX,
		dirY,
	}
}

// IsColideing returns true if this object is colliding with the target object
//
// this method calculates differently depending on an objects CollisionMethod
func (obj1 *GameObject) IsColideing(obj2 *GameObject) bool {
	if obj1.CollisionMethod == CollisionMethod.Ghost || obj2.CollisionMethod == CollisionMethod.Ghost {
		return false
	}

	// prevent accidental self collision
	if obj1.id == obj2.id {
		return false
	}

	if obj1.CollisionMethod == CollisionMethod.Box && obj2.CollisionMethod == CollisionMethod.Box {
		if (obj1.X + obj1.Width > obj2.X - obj2.Width && obj1.X - obj1.Width < obj2.X + obj2.Width) && 
		(obj1.Y + obj1.Height > obj2.Y - obj2.Height && obj1.Y - obj1.Height < obj2.Y + obj2.Height) {
			return true
		}
	}else if obj1.CollisionMethod == CollisionMethod.Radius && obj2.CollisionMethod == CollisionMethod.Radius {
		dir := obj1.GetDirection(obj2)
		return dir.dist <= float32(math.Sqrt(math.Pow(float64(obj1.Width + obj2.Width), 2) + math.Pow(float64(obj1.Height + obj2.Height), 2))) / (math.Pi / 2.25)
	}else if obj1.CollisionMethod == CollisionMethod.Box && obj2.CollisionMethod == CollisionMethod.Radius {
		size := float32(math.Sqrt(math.Pow(float64(obj2.Width), 2) + math.Pow(float64(obj2.Height), 2))) / (math.Pi / 2.25)

		// skip math loop if object is too far away
		if dist := obj1.GetDistance(obj2); dist > size + (obj1.Width * 2) && dist > size + (obj1.Height * 2) {
			return false
		}

		for w := -obj1.Width; w <= obj1.Width; w += obj2.Width / math.Pi {
			for h := -obj1.Height; h <= obj1.Height; h += obj2.Height / math.Pi {
				diffX := (obj1.X + w) - obj2.X
				diffY := (obj1.Y + h) - obj2.Y
				dist := float32(math.Sqrt(math.Pow(float64(diffX), 2) + math.Pow(float64(diffY), 2)))

				if dist <= size {
					return true
				}
			}

			// cover final height check
			diffX := (obj1.X + w) - obj2.X
			diffY := (obj1.Y + obj1.Height) - obj2.Y
			dist := float32(math.Sqrt(math.Pow(float64(diffX), 2) + math.Pow(float64(diffY), 2)))
			if dist <= size {
				return true
			}
		}

		// cover final width checks
		for h := -obj1.Height; h <= obj1.Height; h += obj2.Height / math.Pi {
			diffX := (obj1.X + obj1.Width) - obj2.X
			diffY := (obj1.Y + h) - obj2.Y
			dist := float32(math.Sqrt(math.Pow(float64(diffX), 2) + math.Pow(float64(diffY), 2)))
			if dist <= size {
				return true
			}
		}

		// cover final width and height check
		diffX := (obj1.X + obj1.Width) - obj2.X
		diffY := (obj1.Y + obj1.Height) - obj2.Y
		dist := float32(math.Sqrt(math.Pow(float64(diffX), 2) + math.Pow(float64(diffY), 2)))
		if dist <= size {
			return true
		}

		return false
	}else if obj1.CollisionMethod == CollisionMethod.Radius && obj2.CollisionMethod == CollisionMethod.Box {
		size := float32(math.Sqrt(math.Pow(float64(obj1.Width), 2) + math.Pow(float64(obj1.Height), 2))) / (math.Pi / 2.25)

		// skip math loop if object is too far away
		if dist := obj1.GetDistance(obj2); dist > size + (obj1.Width * 2) && dist > size + (obj1.Height * 2) {
			return false
		}

		for w := -obj2.Width; w <= obj2.Width; w += obj1.Width / math.Pi {
			for h := -obj2.Height; h <= obj2.Height; h += obj1.Height / math.Pi {
				diffX := (obj2.X + w) - obj1.X
				diffY := (obj2.Y + h) - obj1.Y
				dist := float32(math.Sqrt(math.Pow(float64(diffX), 2) + math.Pow(float64(diffY), 2)))

				if dist <= size {
					return true
				}
			}

			// cover final height check
			diffX := (obj2.X + w) - obj1.X
			diffY := (obj2.Y + obj2.Height) - obj1.Y
			dist := float32(math.Sqrt(math.Pow(float64(diffX), 2) + math.Pow(float64(diffY), 2)))
			if dist <= size {
				return true
			}
		}

		// cover final width checks
		for h := -obj2.Height; h <= obj2.Height; h += obj1.Height / math.Pi {
			diffX := (obj2.X + obj2.Width) - obj1.X
			diffY := (obj2.Y + h) - obj1.Y
			dist := float32(math.Sqrt(math.Pow(float64(diffX), 2) + math.Pow(float64(diffY), 2)))
			if dist <= size {
				return true
			}
		}

		// cover final width and height check
		diffX := (obj2.X + obj2.Width) - obj1.X
		diffY := (obj2.Y + obj2.Height) - obj1.Y
		dist := float32(math.Sqrt(math.Pow(float64(diffX), 2) + math.Pow(float64(diffY), 2)))
		if dist <= size {
			return true
		}

		return false
	}

	return false
}

// IsColideingAny returns a list of colliding objects
func (object *GameObject) IsColideingAny() []*GameObject {
	colType := uint8(0)
	if v, ok := gameColissionType[object.objType]; ok {
		colType = v
	}

	list := []*GameObject{}
	if colType == TypeCollisionMethod.Ghost {
		return list
	}

	if colType == TypeCollisionMethod.Self {
		if objList, ok := gameObjects[object.objType]; ok {
			for _, obj := range objList {
				if object.IsColideing(obj) {
					list = append(list, obj)
				}
			}
		}

		return list
	}

	for key, objList := range gameObjects {
		cType := uint8(0)
		if v, ok := gameColissionType[key]; ok {
			cType = v
		}

		if cType == TypeCollisionMethod.Ghost ||
		(object.objType != key && cType == TypeCollisionMethod.Self) ||
		(object.objType == key && colType == TypeCollisionMethod.Other) {
			continue
		}

		for _, obj := range objList {
			if object.IsColideing(obj) {
				list = append(list, obj)
			}
		}
	}

	return list
}

// IsColideingAny returns a list of colliding objects of a specific type
func (object *GameObject) IsColideingType(objType string) []*GameObject {
	colType := uint8(0)
	if v, ok := gameColissionType[object.objType]; ok {
		colType = v
	}

	list := []*GameObject{}
	if colType == TypeCollisionMethod.Ghost ||
	(object.objType != objType && colType == TypeCollisionMethod.Self) ||
	(object.objType == objType && colType == TypeCollisionMethod.Other) {
		return list
	}

	if objList, ok := gameObjects[objType]; ok {
		cType := uint8(0)
		if v, ok := gameColissionType[objType]; ok {
			cType = v
		}

		if cType == TypeCollisionMethod.Ghost ||
		(object.objType != objType && cType == TypeCollisionMethod.Self) {
			return list
		}
		
		for _, obj := range objList {
			if object.IsColideing(obj) {
				list = append(list, obj)
			}
		}
	}

	return list
}

// IsColideingAny returns a list of colliding objects of a specific type and name
func (object *GameObject) IsColideingName(objType string, name string) []*GameObject {
	colType := uint8(0)
	if v, ok := gameColissionType[object.objType]; ok {
		colType = v
	}

	list := []*GameObject{}
	if colType == TypeCollisionMethod.Ghost ||
	(object.objType != objType && colType == TypeCollisionMethod.Self) ||
	(object.objType == objType && colType == TypeCollisionMethod.Other) {
		return list
	}

	if objList, ok := gameObjects[objType]; ok {
		cType := uint8(0)
		if v, ok := gameColissionType[objType]; ok {
			cType = v
		}

		if cType == TypeCollisionMethod.Ghost ||
		(object.objType != objType && cType == TypeCollisionMethod.Self) {
			return list
		}
		
		for _, obj := range objList {
			if obj.name == name && object.IsColideing(obj) {
				list = append(list, obj)
			}
		}
	}

	return list
}


// basic methods

// Update should run on a GameLoop thread
//
// recommended: 60 fps
func Update(game *Game, thread *ThreadInfo){
	game.eachObject(func(object *GameObject) {
		if object.PreferredFPS >= 60 && object.PreferredFPS < 120 { 
			object.handleBorder(game, thread)
		}

		if object.Update != nil {
			object.Update(game, thread)
		}
	})
}

// Draw should run on a GameLoop thread
//
// recommended: 120 fps
func Draw(game *Game, thread *ThreadInfo){
	game.eachObject(func(object *GameObject) {
		if object.PreferredFPS == 0 || object.PreferredFPS > 120 { 
			object.handleBorder(game, thread)
		}

		// handle object border method
		switch object.BorderMethod {
		case BorderMethod.Ignore:
			object.X += object.VelX / 10 * thread.SpeedDelta
			object.Y += object.VelY / 10 * thread.SpeedDelta

		case BorderMethod.Limit:
			if object.OnBorderX == 0 || object.OnBorderX % 2 != 0 {
				object.X += object.VelX / 10 * thread.SpeedDelta
			}
			if object.OnBorderY == 0 || object.OnBorderY % 2 != 0 {
				object.Y += object.VelY / 10 * thread.SpeedDelta
			}

		case BorderMethod.Hide:
			if object.OnBorderX >= -3 && object.OnBorderX <= 3 {
				object.X += object.VelX / 10 * thread.SpeedDelta
			}
			if object.OnBorderY >= -3 && object.OnBorderY <= 3 {
				object.Y += object.VelY / 10 * thread.SpeedDelta
			}

		case BorderMethod.PushLimit:
			if object.OnBorderX == 0 || object.OnBorderX % 2 != 0 {
				object.X += object.VelX / 10 * thread.SpeedDelta
			}
			if object.OnBorderY == 0 || object.OnBorderY % 2 != 0 {
				object.Y += object.VelY / 10 * thread.SpeedDelta
			}

		case BorderMethod.PushHide:
			if object.OnBorderX >= -3 && object.OnBorderX <= 3 {
				object.X += object.VelX / 10 * thread.SpeedDelta
			}
			if object.OnBorderY >= -3 && object.OnBorderY <= 3 {
				object.Y += object.VelY / 10 * thread.SpeedDelta
			}

		default:
			object.X += object.VelX / 10 * thread.SpeedDelta
			object.Y += object.VelY / 10 * thread.SpeedDelta
		}

		if object.Draw != nil {
			object.Draw(game, thread)
		}

		object.Object.Move(fyne.NewPos(((object.X - object.Width) * game.Size.Scale) + (game.Size.RealWidth/2), ((object.Y - object.Height) * game.Size.Scale) + (game.Size.RealHeight/2)))
		object.Object.Resize(fyne.NewSize((object.Width * 2) * game.Size.Scale, (object.Height * 2) * game.Size.Scale))
		object.Object.Refresh()
	})
}

// UpdateBasic should run on a GameLoop thread
//
// recommended: 15 fps
func UpdateBasic(game *Game, thread *ThreadInfo){
	game.eachObject(func(object *GameObject) {
		if object.PreferredFPS >= 15 && object.PreferredFPS < 30 { 
			object.handleBorder(game, thread)
		}

		if object.UpdateBasic != nil {
			object.UpdateBasic(game, thread)
		}
	})
}

// UpdateSlow should run on a GameLoop thread
//
// recommended: 30 fps
func UpdateSlow(game *Game, thread *ThreadInfo){
	game.eachObject(func(object *GameObject) {
		if object.PreferredFPS >= 30 && object.PreferredFPS < 60 { 
			object.handleBorder(game, thread)
		}

		if object.UpdateSlow != nil {
			object.UpdateSlow(game, thread)
		}
	})
}
