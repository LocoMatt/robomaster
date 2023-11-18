package sdk2

import (
	"fmt"
	"sync"

	"github.com/brunoga/robomaster/sdk2/module/camera"
	"github.com/brunoga/robomaster/sdk2/module/connection"
	"github.com/brunoga/robomaster/sdk2/module/controller"
	"github.com/brunoga/robomaster/sdk2/module/robot"
	"github.com/brunoga/unitybridge"
	"github.com/brunoga/unitybridge/support/logger"
	"github.com/brunoga/unitybridge/wrapper"
)

type Client struct {
	l *logger.Logger

	ub unitybridge.UnityBridge

	cn *connection.Connection
	cm *camera.Camera
	ct *controller.Controller
	rb *robot.Robot

	m       sync.RWMutex
	started bool
}

// New creates a new Client instance with the given logger and appID.
func New(l *logger.Logger, appID uint64) (*Client, error) {
	ub := unitybridge.Get(wrapper.Get(l), true, l)

	cn, err := connection.New(ub, l, appID)
	if err != nil {
		return nil, err
	}

	rb, err := robot.New(ub, l)
	if err != nil {
		return nil, err
	}

	cm, err := camera.New(ub, l)
	if err != nil {
		return nil, err
	}

	ct, err := controller.New(rb, ub, l)
	if err != nil {
		return nil, err
	}

	return &Client{
		ub: ub,
		l:  l,
		cn: cn,
		rb: rb,
		cm: cm,
		ct: ct,
	}, nil
}

// Start starts the client and all associated modules.
func (c *Client) Start() error {
	c.m.Lock()
	defer c.m.Unlock()

	if c.started {
		return fmt.Errorf("client already started")
	}

	err := c.ub.Start()
	if err != nil {
		return err
	}

	// Start modules.

	// Connection.
	err = c.cn.Start()
	if err != nil {
		return err
	}

	if !c.cn.WaitForConnection() {
		return fmt.Errorf("network connection unexpectedly not established")
	}

	// Robot.
	err = c.rb.Start()
	if err != nil {
		return err
	}

	if !c.rb.WaitForConnection() {
		return fmt.Errorf("robot connection unexpectedly not established")
	}

	// Camera.
	err = c.cm.Start()
	if err != nil {
		return err
	}

	if !c.cm.WaitForConnection() {
		return fmt.Errorf("camera connection unexpectedly not established")
	}

	// Controller.
	err = c.ct.Start()
	if err != nil {
		return err
	}

	if !c.ct.WaitForConnection() {
		return fmt.Errorf("controller connection unexpectedly not established")
	}

	c.started = true

	return nil
}

// Connection returns the Connection module.
func (c *Client) Connection() *connection.Connection {
	return c.cn
}

// Camera returns the Camera module.
func (c *Client) Camera() *camera.Camera {
	return c.cm
}

// Controller returns the Controller module.
func (c *Client) Controller() *controller.Controller {
	return c.ct
}

// Robot returns the Robot module.
func (c *Client) Robot() *robot.Robot {
	return c.rb
}

// Stop stops the client and all associated modules.
func (c *Client) Stop() error {
	c.m.Lock()
	defer c.m.Unlock()

	if !c.started {
		return fmt.Errorf("client not started")
	}

	// Stop modules.

	// Controller.
	err := c.ct.Stop()
	if err != nil {
		return err
	}

	// Camera.
	err = c.cm.Stop()
	if err != nil {
		return err
	}

	// Robot.
	err = c.rb.Stop()
	if err != nil {
		return err
	}

	// Connection.
	err = c.cn.Stop()
	if err != nil {
		return err
	}

	// Stop Unity Bridge.
	err = c.ub.Stop()
	if err != nil {
		return err
	}

	return nil
}
