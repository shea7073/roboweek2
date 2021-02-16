package main

import (
	"fmt"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/aio"
	"gobot.io/x/gobot/drivers/i2c"
	g "gobot.io/x/gobot/platforms/dexter/gopigo3"
	"gobot.io/x/gobot/platforms/raspi"
	"time"
)

//robotRunLoop is the main function for the robot, the gobot framework
//will spawn a new thread in the NewRobot factory functin and run this
//function in that new thread. Do all of your work in this function and
//in other functions that this function calls. don't read from sensors or
//use actuators frmo main or you will get a panic.
//add
func robotRunLoop(lightSensor *aio.GroveLightSensorDriver, soundSensor *aio.GroveSoundSensorDriver, gpg *g.Driver, lidarSensor *i2c.LIDARLiteDriver, m map[int]int64) {
	gpg.SetMotorDps(g.MOTOR_RIGHT, 30)
	err := lidarSensor.Start()
	if err != nil {
		fmt.Errorf("Error starting lidar %+v", err)
	}
	for {
		sensorVal, err := lightSensor.Read()
		if err != nil {
			fmt.Errorf("Error reading light sensor %+v", err)
		}
		soundSensorVal, err := soundSensor.Read()
		if err != nil {
			fmt.Errorf("Error reading from Sound Sensor %+v", err)
		}
		val, err := gpg.GetMotorEncoder(g.MOTOR_RIGHT)
		if err != nil {
			fmt.Errorf("Error reading from encoder %+v", err)
		}
		lidarVal, err := lidarSensor.Distance()
		if err != nil {
			fmt.Errorf("Error reading from lidar %+v", err)
		}

		fmt.Println("Light Value is ", sensorVal)
		fmt.Println("Sound Value is ", soundSensorVal)
		fmt.Println("encoder value: ", val)
		m[sensorVal] = val
		time.Sleep(time.Second)

		gpg.Start()

	}
}

func maxNumber(m map[int]int64) int {
	var max int
	for n := range m {
		max = n
		break
	}
	for n := range m {
		if n > max {
			max = n
		}
	}
	return max
}
func main() {
	//We create the adaptors to connect the GoPiGo3 board with the Raspberry Pi 3
	//also create any sensor drivers here
	raspiAdaptor := raspi.NewAdaptor()
	gopigo3 := g.NewDriver(raspiAdaptor)
	lightSensor := aio.NewGroveLightSensorDriver(gopigo3, "AD_2_1") //AnalogDigital Port 1 is "AD_1_1" this is port 2
	soundSensor := aio.NewGroveSoundSensorDriver(gopigo3, "AD_1_1")
	lidarSensor := i2c.NewLIDARLiteDriver(raspiAdaptor)
	//end create hardware drivers

	//here we create an anonymous function assigned to a local variable
	//the robot framework will create a new thread and run this function
	//I'm calling my robot main loop here. Pass any of the variables we created
	//above to that function if you need them
	m := make(map[int]int64)
	mainRobotFunc := func() {
		robotRunLoop(lightSensor, soundSensor, gopigo3, lidarSensor, m)
	}

	//this is the crux of the gobot framework. The factory function to create a new robot
	//struct (go uses structs and not objects) It takes four parameters
	robot := gobot.NewRobot("gopigo3sensorChecker", //first a name
		[]gobot.Connection{raspiAdaptor},                               //next a slice of connections to one or more robot controllers
		[]gobot.Device{gopigo3, lightSensor, soundSensor, lidarSensor}, //next a slice of one or more sensors and actuators for the robots
		mainRobotFunc, //the variable holding the function to run in a new thread as the main function
	)

	robot.Start() //actually run the function
}
