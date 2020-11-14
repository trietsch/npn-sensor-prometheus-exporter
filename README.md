# NPN Meter Prometheus Exporter

Based on this [blogpost](https://ehoco.nl/watermeter-uitlezen-in-domoticz-python-script/), I wanted to get the NPN data (and thus the water meter gauge) into my monitoring stack (Prometheus + Grafana).
Therefore, I decided to write my own Prometheus exporter to read the data from GPIO and expose the data through a metrics endpoint.

## Contributions

You're welcome to clone, fork, adjust the code to your needs. If you think that others may benefit as well, please make a pull request.
