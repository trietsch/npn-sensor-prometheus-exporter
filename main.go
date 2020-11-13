package main

import (
	"flag"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stianeikeland/go-rpio"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var listenAddr string
var verboseLogging bool

var currentWaterLevelFile string
var currentWaterLevel uint

var (
	registry        = prometheus.NewRegistry()
	waterGaugeLevel = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "water_usage",
		Help: "Water usage in liters",
	})
)

func init() {
	registry.MustRegister(waterGaugeLevel)
}

func main() {
	flag.StringVar(&listenAddr, "listen", "127.0.0.1:8787", "Listen address for HTTP metrics")
	flag.StringVar(&currentWaterLevelFile, "file", "water_gauge_level.txt", "File to read/write the state")
	flag.BoolVar(&verboseLogging, "verbose", false, "Verbose output logging")
	flag.Parse()

	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	if verboseLogging {
		logrus.SetLevel(logrus.DebugLevel)
	}

	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Get the current value
	currentWaterLevel = getWaterLevelFromFile()

	// Get the pin on which the NPN sensor is connected
	npnPin := rpio.Pin(21)

	// Set the pin to input mode
	npnPin.Input()
	npnPin.Detect(rpio.FallEdge)

	// Clean up on ctrl-c
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		logrus.Infoln("Shutting down NPN exporter. Bye!")
		writeWaterLevelToFile()
		os.Exit(0)
	}()

	defer rpio.Close()

	ticker := time.NewTicker(5 * time.Second)
	go func() {
		for {
			select {
			case _ = <-ticker.C:
				writeWaterLevelToFile()
			}
		}
	}()

	go func() {
		// Infinite loop for listening for changes
		for {
			if npnPin.EdgeDetected() {
				currentWaterLevel++
				waterGaugeLevel.Set(float64(currentWaterLevel))
			}

			time.Sleep(500 * time.Millisecond)
			currentWaterLevel++
			waterGaugeLevel.Set(float64(currentWaterLevel))
		}
	}()

	logrus.Infoln("Start listening at", listenAddr)
	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	logrus.Fatalln(http.ListenAndServe(listenAddr, nil))
}

func writeWaterLevelToFile() {
	logrus.Infoln("Writing water level: ", currentWaterLevel)

	file, _ := os.Create(currentWaterLevelFile)

	defer file.Close()

	_, err := file.WriteString(fmt.Sprintf("%d", currentWaterLevel))
	if err != nil {
		logrus.Errorln("error writing string: ", err)
	}
}

func getWaterLevelFromFile() uint {
	logrus.Infoln("Reading water level")
	f, _ := os.Open(currentWaterLevelFile)
	var waterLevel uint
	_, _ = fmt.Fscanln(f, &waterLevel)
	return waterLevel
}
