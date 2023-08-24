package CollisionMethod

// Ghost sets an object to ignore any collision
const Ghost uint8 = 0

// Box compares the X and Y position of an object, allowing for a square hitbox
const Box uint8 = 1

// Radius compares the distance of an object, allowing for a round hitbox
const Radius uint8 = 2
