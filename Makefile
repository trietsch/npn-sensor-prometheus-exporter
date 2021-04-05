clean:
	rm -rf dist
	rm -f npn_sensor_prometheus_exporter

build: clean
	env GOOS=linux GOARCH=arm GOARM=7 go build -o npn_sensor_prometheus_exporter

copy: build
	scp -i ~/.ssh/id_iron npn_sensor_prometheus_exporter pi@meterkast.local:/home/pi/

release: clean
	@read -p "Enter new tag name:" tagName; \
	git tag -a "$${tagName}" -m "Tagging $${tagName}"; \
	git push origin $${tagName}; \
	goreleaser release --rm-dist
