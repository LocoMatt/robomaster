package chassis

import (
	"log/slog"
	"os"
	"testing"

	robomaster "github.com/brunoga/robomaster"
	"github.com/brunoga/robomaster/module"
	"github.com/brunoga/robomaster/module/chassis"
	"github.com/brunoga/robomaster/module/robot"
	"github.com/brunoga/robomaster/unitybridge/support"
	"github.com/brunoga/robomaster/unitybridge/support/logger"
)

var chassisModule *chassis.Chassis

func TestMain(m *testing.M) {
	c, err := robomaster.NewWithModules(logger.New(slog.LevelDebug), support.AnyAppID,
		module.TypeConnection|module.TypeRobot|module.TypeController|module.TypeChassis)
	if err != nil {
		panic(err)
	}

	if err := c.Start(); err != nil {
		panic(err)
	}
	defer func() {
		if err := c.Stop(); err != nil {
			panic(err)
		}
	}()

	chassisModule = c.Chassis()

	// Set controller mode to SDK for the tests here.
	//err = c.Controller().SetMode(controller.ModeSDK)
	//if err != nil {
	//	panic(err)
	//}

	// Enable robot movement.
	c.Robot().EnableFunction(robot.FunctionTypeMovementControl, true)

	os.Exit(m.Run())
}
