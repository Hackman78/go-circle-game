package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"time"
	"math/rand"
)

var spawnTimer float32
var spawnInterval float32 = 1.2 // seconds between spawns
var elapsedTime float32 = 0

var screenWidth = int32(800)
var screenHeight = int32(450)

type EnemyBall struct {
	Position rl.Vector2
	Radius   float32
	Speed    rl.Vector2
}

var enemies []EnemyBall

type Hero struct {
	Position rl.Vector2
	Radius   float32
	Speed    rl.Vector2
}

func detectCollision(radius *float32, ballPosition rl.Vector2) bool {
	newEnemies := enemies[:0] // reuse underlying array
	for _, e := range enemies {
		dx := ballPosition.X - e.Position.X
		dy := ballPosition.Y - e.Position.Y
		distance := rl.Vector2Length(rl.NewVector2(dx, dy))

		if distance < *radius+e.Radius {
			if *radius > e.Radius {
				*radius += e.Radius * 0.1
				continue // Skip adding this enemy (it’s eaten)
			} else {
				return true // Game over
			}
		}
		newEnemies = append(newEnemies, e)
	}
	enemies = newEnemies
	return false
}


// func randomSpawnEnemies() {
//
// 	deltaTime := rl.GetFrameTime()
// 	spawnTimer += deltaTime
//
// 	if spawnTimer >= spawnInterval {
// 		spawnTimer = 0
//
// 		enemies = append(enemies, EnemyBall{
// 			Position: rl.NewVector2(
// 				float32(screenWidth)+float32(rand.Intn(300)),
// 				float32(rand.Intn(int(screenHeight-40))+20),
// 				),
// 			Radius: float32(rand.Intn(10) + 10),
// 			Speed:  rl.NewVector2(-float32(rand.Intn(3)+1), 0),
// 		})
// 	}
// }

func spawnEnemyFromEdge(heroRadius float32) {
	side := rand.Intn(4) // 0=left, 1=right, 2=top, 3=bottom
	var pos rl.Vector2
	var speed rl.Vector2

	switch side {
	case 0: // Left
		pos = rl.NewVector2(-20, float32(rand.Intn(int(screenHeight))))
		speed = rl.NewVector2(rand.Float32()*2+1, rand.Float32()*2-1)
	case 1: // Right
		pos = rl.NewVector2(float32(screenWidth)+20, float32(rand.Intn(int(screenHeight))))
		speed = rl.NewVector2(-rand.Float32()*2-1, rand.Float32()*2-1)
	case 2: // Top
		pos = rl.NewVector2(float32(rand.Intn(int(screenWidth))), -20)
		speed = rl.NewVector2(rand.Float32()*2-1, rand.Float32()*2+1)
	case 3: // Bottom
		pos = rl.NewVector2(float32(rand.Intn(int(screenWidth))), float32(screenHeight)+20)
		speed = rl.NewVector2(rand.Float32()*2-1, -rand.Float32()*2-1)
	}

	speedMultiplier := 1.0 + elapsedTime/60.0
	if speedMultiplier > 3 {
		speedMultiplier = 3
	}
	speed.X *= float32(speedMultiplier)
	speed.Y *= float32(speedMultiplier)


	var radius float32
	if rand.Float32() < 0.5 {
		radius = heroRadius + rand.Float32()*10 // slightly larger
	} else {
		radius = rand.Float32()*heroRadius*0.8 // definitely smaller
	}
	enemies = append(enemies, EnemyBall{
		Position: pos,
		Radius:   radius,
		Speed:    speed,
	})
}

func detectKeys(ballPosition *rl.Vector2,speed float32, radius float32) {
	if rl.IsKeyDown(rl.KeyRight) {
		ballPosition.X += speed
	}
	if rl.IsKeyDown(rl.KeyLeft) {
		ballPosition.X -= speed
	}
	if rl.IsKeyDown(rl.KeyUp) {
		ballPosition.Y -= speed
	}
	if rl.IsKeyDown(rl.KeyDown) {
		ballPosition.Y += speed
	}

	if ballPosition.X < radius {
		ballPosition.X = radius
	}
	if ballPosition.X > float32(screenWidth)-radius {
		ballPosition.X = float32(screenWidth) - radius
	}
	if ballPosition.Y < radius {
		ballPosition.Y = radius
	}
	if ballPosition.Y > float32(screenHeight)-radius {
		ballPosition.Y = float32(screenHeight) - radius
	}

}


func main() {
	rl.InitWindow(screenWidth, screenHeight, "Platformer")
	defer rl.CloseWindow()

	ballPosition := rl.NewVector2(float32(screenWidth)/2, float32(screenHeight)/2)
	rand.Seed(time.Now().UnixNano())



	rl.SetTargetFPS(144)

	radius := float32(20)

	for i := 0; i < 10; i++ {
		spawnEnemyFromEdge(radius)
	}

	for !rl.WindowShouldClose() {
		speed := 100 / radius

		spawnInterval = 2.0 - (elapsedTime / 30.0) // over 60 seconds from 2s to 0.5s
		if spawnInterval < 0.5 {
			spawnInterval = 0.5
		}
		deltaTime := rl.GetFrameTime()
		elapsedTime += deltaTime
		spawnTimer += deltaTime

		if spawnTimer >= spawnInterval {
			spawnTimer = 0
			spawnEnemyFromEdge(radius)
		}

		detectKeys(&ballPosition, speed, radius)

		rl.BeginDrawing()

		if detectCollision(&radius, ballPosition) {
			rl.DrawText("Game Over", 300, 200, 40, rl.Red)
			rl.EndDrawing()
			time.Sleep(2 * time.Second)
			return
		}

		rl.ClearBackground(rl.RayWhite)

		rl.DrawText("move the ball with arrow keys", 10, 10, 20, rl.DarkGray)
		rl.DrawText("eat balls smaller than you", 10, 25, 20, rl.DarkGray)
		rl.DrawText("avoid balls bigger than you", 10, 40, 20, rl.DarkGray)
		rl.DrawCircleV(ballPosition, radius, rl.Maroon)

		for i := range enemies {
			enemies[i].Position.X += enemies[i].Speed.X
			enemies[i].Position.Y += enemies[i].Speed.Y
		}

		for _, e := range enemies {
			rl.DrawCircleV(e.Position, e.Radius, rl.DarkBlue)
		}


		rl.EndDrawing()
	}
}
