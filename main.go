package main

import rl "github.com/gen2brain/raylib-go/raylib"

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

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	rl.InitWindow(screenWidth, screenHeight, "Platformer")

	ballPosition := rl.NewVector2(float32(screenWidth)/2, float32(screenHeight)/2)

	defer rl.CloseWindow()

	rl.SetTargetFPS(144)

	radius := float32(20)

	for i := 0; i < 5; i++ {
		enemies = append(enemies, EnemyBall{
			Position: rl.NewVector2(float32(screenWidth)+float32(i*100), float32(50+i*70)),
			Radius:   10 + float32(i*3),
			Speed:    rl.NewVector2(-2, 0), // Move left
		})
	}


	for !rl.WindowShouldClose() {
		speed := 100 / radius

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
		rl.BeginDrawing()


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
