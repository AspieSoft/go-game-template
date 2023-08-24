package TypeCollisionMethod

// Any allows object types to freely detect collision
const Any uint8 = 0

// Ghost sets an object to ignore any collision
const Ghost uint8 = 1

// Self limits a type to only detecting collisions of the same type
const Self uint8 = 2

// Other limits a type to only detecting collisions of a different type
const Other uint8 = 3
