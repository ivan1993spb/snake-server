package snake

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ivan1993spb/snake-server/engine"
	"github.com/ivan1993spb/snake-server/world"
)

func Test_NewSnake(t *testing.T) {

}

func Test_Snake_calculateDelay(t *testing.T) {
	firstSnake := &Snake{
		length: 10,
	}
	require.NotZero(t, firstSnake.calculateDelay())

	secondSnake := &Snake{
		length: 11,
	}
	require.NotZero(t, secondSnake.calculateDelay())

	require.True(t, firstSnake.calculateDelay() < secondSnake.calculateDelay())
}

func Test_Snake_setMovementDirection(t *testing.T) {
	world, err := world.NewWorld(100, 100)
	require.Nil(t, err, "cannot initialize world")
	require.NotNil(t, world, "cannot initialize world")

	snake := &Snake{
		world:  world,
		length: 4,
		dots: []engine.Dot{
			{10, 0},
			{9, 0},
			{8, 0},
			{7, 0},
		},
		direction: engine.DirectionEast,
		mux:       &sync.RWMutex{},
	}

	err = world.CreateObject(snake, engine.Location{
		engine.Dot{10, 0},
		engine.Dot{9, 0},
		engine.Dot{8, 0},
		engine.Dot{7, 0},
	})
	require.Nil(t, err, "cannot create object")

	require.Nil(t, snake.setMovementDirection(engine.DirectionNorth))
	require.Equal(t, engine.DirectionNorth, snake.direction)

	require.Nil(t, snake.setMovementDirection(engine.DirectionSouth))
	require.Equal(t, engine.DirectionSouth, snake.direction)

	require.NotNil(t, snake.setMovementDirection(engine.DirectionWest))
	require.Equal(t, engine.DirectionSouth, snake.direction)
}

func Test_Snake_getNextHeadDot(t *testing.T) {
	world, err := world.NewWorld(100, 100)
	require.Nil(t, err, "cannot initialize world")
	require.NotNil(t, world, "cannot initialize world")

	// First case east

	snake := &Snake{
		world:  world,
		length: 4,
		dots: []engine.Dot{
			{10, 0},
			{9, 0},
			{8, 0},
			{7, 0},
		},
		direction: engine.DirectionEast,
		mux:       &sync.RWMutex{},
	}

	err = world.CreateObject(snake, engine.Location{
		engine.Dot{10, 0},
		engine.Dot{9, 0},
		engine.Dot{8, 0},
		engine.Dot{7, 0},
	})
	require.Nil(t, err, "cannot create object")

	dot, err := snake.getNextHeadDot()
	require.Nil(t, err)
	require.Equal(t, engine.Dot{11, 0}, dot)

	// Second case west

	snake = &Snake{
		world:  world,
		length: 4,
		dots: []engine.Dot{
			{2, 5},
			{3, 5},
			{4, 5},
			{5, 5},
		},
		direction: engine.DirectionWest,
	}

	err = world.CreateObject(snake, engine.Location{
		engine.Dot{2, 5},
		engine.Dot{3, 5},
		engine.Dot{4, 5},
		engine.Dot{5, 5},
	})
	require.Nil(t, err, "cannot create object")

	dot, err = snake.getNextHeadDot()
	require.Nil(t, err)
	require.Equal(t, engine.Dot{1, 5}, dot)

	// Third case north

	snake = &Snake{
		world:  world,
		length: 4,
		dots: []engine.Dot{
			{10, 10},
			{10, 11},
			{10, 12},
			{10, 13},
		},
		direction: engine.DirectionNorth,
		mux:       &sync.RWMutex{},
	}

	err = world.CreateObject(snake, engine.Location{
		engine.Dot{10, 10},
		engine.Dot{10, 11},
		engine.Dot{10, 12},
		engine.Dot{10, 13},
	})
	require.Nil(t, err, "cannot create object")

	dot, err = snake.getNextHeadDot()
	require.Nil(t, err)
	require.Equal(t, engine.Dot{10, 9}, dot)

	// Fourth case south

	snake = &Snake{
		world:  world,
		length: 4,
		dots: []engine.Dot{
			{20, 24},
			{20, 21},
			{20, 20},
			{20, 19},
		},
		direction: engine.DirectionSouth,
		mux:       &sync.RWMutex{},
	}

	err = world.CreateObject(snake, engine.Location{
		engine.Dot{20, 24},
		engine.Dot{20, 21},
		engine.Dot{20, 20},
		engine.Dot{20, 19},
	})
	require.Nil(t, err, "cannot create object")

	dot, err = snake.getNextHeadDot()
	require.Nil(t, err)
	require.Equal(t, engine.Dot{20, 25}, dot)
}

func Test_Snake_move_validLocation(t *testing.T) {
	world, err := world.NewWorld(100, 100)
	require.Nil(t, err, "cannot initialize world")
	require.NotNil(t, world, "cannot initialize world")

	snake := &Snake{
		world:  world,
		length: 4,
		dots: []engine.Dot{
			{10, 0},
			{9, 0},
			{8, 0},
			{7, 0},
		},
		direction: engine.DirectionEast,
		mux:       &sync.RWMutex{},
	}

	err = world.CreateObject(snake, engine.Location{
		engine.Dot{10, 0},
		engine.Dot{9, 0},
		engine.Dot{8, 0},
		engine.Dot{7, 0},
	})
	require.Nil(t, err, "cannot create object")

	require.Nil(t, snake.move())
	require.Equal(t, []engine.Dot{
		{11, 0},
		{10, 0},
		{9, 0},
		{8, 0},
	}, snake.dots)

	require.Nil(t, snake.move())
	require.Equal(t, []engine.Dot{
		{12, 0},
		{11, 0},
		{10, 0},
		{9, 0},
	}, snake.dots)

	require.Nil(t, snake.move())
	require.Equal(t, []engine.Dot{
		{13, 0},
		{12, 0},
		{11, 0},
		{10, 0},
	}, snake.dots)

	require.Nil(t, snake.move())
	require.Equal(t, []engine.Dot{
		{14, 0},
		{13, 0},
		{12, 0},
		{11, 0},
	}, snake.dots)
}
