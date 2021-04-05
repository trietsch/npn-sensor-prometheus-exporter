clean:
	rm -rf dist
	rm npn_sensor_prometheus_exporter

build: clean
	env GOOS=linux GOARCH=arm GOARM=7 go build -o npn_sensor_prometheus_exporter

copy: build
	scp -i ~/.ssh/id_iron npn_sensor_prometheus_exporter pi@meterkast.local:/home/pi/
