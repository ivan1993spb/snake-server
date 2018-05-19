package snake

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ivan1993spb/snake-server/engine"
	"github.com/ivan1993spb/snake-server/game"
	"github.com/ivan1993spb/snake-server/playground"
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
	area, err := engine.NewArea(100, 100)
	require.Nil(t, err, "cannot create area")
	require.NotNil(t, area, "cannot create area")

	scene, err := engine.NewScene(area)
	require.Nil(t, err, "cannot create scene")
	require.NotNil(t, scene, "cannot create scene")

	pg := playground.NewPlayground(scene)
	require.NotNil(t, pg, "cannot initialize playground")
	world := game.NewWorld(pg)
	require.NotNil(t, world, "cannot initialize world")

	snake := &Snake{
		world:  world,
		length: 4,
		dots: []*engine.Dot{
			engine.NewDot(10, 0),
			engine.NewDot(9, 0),
			engine.NewDot(8, 0),
			engine.NewDot(7, 0),
		},
		direction: engine.DirectionEast,
		mux:       &sync.RWMutex{},
	}

	err = world.CreateObject(snake, engine.Location{
		engine.NewDot(10, 0),
		engine.NewDot(9, 0),
		engine.NewDot(8, 0),
		engine.NewDot(7, 0),
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
	area, err := engine.NewArea(100, 100)
	require.Nil(t, err, "cannot create area")
	require.NotNil(t, area, "cannot create area")

	scene, err := engine.NewScene(area)
	require.Nil(t, err, "cannot create scene")
	require.NotNil(t, scene, "cannot create scene")

	pg := playground.NewPlayground(scene)
	require.NotNil(t, pg, "cannot initialize playground")
	world := game.NewWorld(pg)
	require.NotNil(t, world, "cannot initialize world")

	// First case east

	snake := &Snake{
		world:  world,
		length: 4,
		dots: []*engine.Dot{
			engine.NewDot(10, 0),
			engine.NewDot(9, 0),
			engine.NewDot(8, 0),
			engine.NewDot(7, 0),
		},
		direction: engine.DirectionEast,
		mux:       &sync.RWMutex{},
	}

	err = world.CreateObject(snake, engine.Location{
		engine.NewDot(10, 0),
		engine.NewDot(9, 0),
		engine.NewDot(8, 0),
		engine.NewDot(7, 0),
	})
	require.Nil(t, err, "cannot create object")

	dot, err := snake.getNextHeadDot()
	require.Nil(t, err)
	require.Equal(t, engine.NewDot(11, 0), dot)

	// Second case west

	snake = &Snake{
		world:  world,
		length: 4,
		dots: []*engine.Dot{
			engine.NewDot(2, 5),
			engine.NewDot(3, 5),
			engine.NewDot(4, 5),
			engine.NewDot(5, 5),
		},
		direction: engine.DirectionWest,
	}

	err = world.CreateObject(snake, engine.Location{
		engine.NewDot(2, 5),
		engine.NewDot(3, 5),
		engine.NewDot(4, 5),
		engine.NewDot(5, 5),
	})
	require.Nil(t, err, "cannot create object")

	dot, err = snake.getNextHeadDot()
	require.Nil(t, err)
	require.Equal(t, engine.NewDot(1, 5), dot)

	// Third case north

	snake = &Snake{
		world:  world,
		length: 4,
		dots: []*engine.Dot{
			engine.NewDot(10, 10),
			engine.NewDot(10, 11),
			engine.NewDot(10, 12),
			engine.NewDot(10, 13),
		},
		direction: engine.DirectionNorth,
		mux:       &sync.RWMutex{},
	}

	err = world.CreateObject(snake, engine.Location{
		engine.NewDot(10, 10),
		engine.NewDot(10, 11),
		engine.NewDot(10, 12),
		engine.NewDot(10, 13),
	})
	require.Nil(t, err, "cannot create object")

	dot, err = snake.getNextHeadDot()
	require.Nil(t, err)
	require.Equal(t, engine.NewDot(10, 9), dot)

	// Fourth case south

	snake = &Snake{
		world:  world,
		length: 4,
		dots: []*engine.Dot{
			engine.NewDot(20, 24),
			engine.NewDot(20, 21),
			engine.NewDot(20, 20),
			engine.NewDot(20, 19),
		},
		direction: engine.DirectionSouth,
		mux:       &sync.RWMutex{},
	}

	err = world.CreateObject(snake, engine.Location{
		engine.NewDot(20, 24),
		engine.NewDot(20, 21),
		engine.NewDot(20, 20),
		engine.NewDot(20, 19),
	})
	require.Nil(t, err, "cannot create object")

	dot, err = snake.getNextHeadDot()
	require.Nil(t, err)
	require.Equal(t, engine.NewDot(20, 25), dot)
}

func Test_Snake_move_validLocation(t *testing.T) {
	area, err := engine.NewArea(100, 100)
	require.Nil(t, err, "cannot create area")
	require.NotNil(t, area, "cannot create area")

	scene, err := engine.NewScene(area)
	require.Nil(t, err, "cannot create scene")
	require.NotNil(t, scene, "cannot create scene")

	pg := playground.NewPlayground(scene)
	require.NotNil(t, pg, "cannot initialize playground")
	world := game.NewWorld(pg)
	require.NotNil(t, world, "cannot initialize world")

	snake := &Snake{
		world:  world,
		length: 4,
		dots: []*engine.Dot{
			engine.NewDot(10, 0),
			engine.NewDot(9, 0),
			engine.NewDot(8, 0),
			engine.NewDot(7, 0),
		},
		direction: engine.DirectionEast,
		mux:       &sync.RWMutex{},
	}

	err = world.CreateObject(snake, engine.Location{
		engine.NewDot(10, 0),
		engine.NewDot(9, 0),
		engine.NewDot(8, 0),
		engine.NewDot(7, 0),
	})
	require.Nil(t, err, "cannot create object")

	require.Nil(t, snake.move())
	require.Equal(t, []*engine.Dot{
		engine.NewDot(11, 0),
		engine.NewDot(10, 0),
		engine.NewDot(9, 0),
		engine.NewDot(8, 0),
	}, snake.dots)

	require.Nil(t, snake.move())
	require.Equal(t, []*engine.Dot{
		engine.NewDot(12, 0),
		engine.NewDot(11, 0),
		engine.NewDot(10, 0),
		engine.NewDot(9, 0),
	}, snake.dots)

	require.Nil(t, snake.move())
	require.Equal(t, []*engine.Dot{
		engine.NewDot(13, 0),
		engine.NewDot(12, 0),
		engine.NewDot(11, 0),
		engine.NewDot(10, 0),
	}, snake.dots)

	require.Nil(t, snake.move())
	require.Equal(t, []*engine.Dot{
		engine.NewDot(14, 0),
		engine.NewDot(13, 0),
		engine.NewDot(12, 0),
		engine.NewDot(11, 0),
	}, snake.dots)
}
