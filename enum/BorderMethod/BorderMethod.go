package BorderMethod

// Ignore does nothing if an object touches a border
const Ignore uint8 = 0

// Limit prevents an object from moving off screen
const Limit uint8 = 1

// Hide allows an object to go fully off screen, then prevents it from wondering too far
const Hide uint8 = 2

// PushLimit is just like 'Limit', but it also pushes back an element if the canvas size changes
const PushLimit uint8 = 3

// PushHide is just like 'Hide', but it also pushes back an element if the canvas size changes
const PushHide uint8 = 4

// Bounce reverses an objects velocity when it touches a border
const Bounce uint8 = 5

// Teleport moves an object to the other side of the screen if once it becomes hidden
const Teleport uint8 = 6

// RemoveObject deletes an object from the game if it goes off screen
const RemoveObject uint8 = 7
