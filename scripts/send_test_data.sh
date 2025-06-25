#!/usr/bin/env bash

BROKER="$1"
PORT='1883'
CLIENT_ID='test-client'
TOPIC="$2"
MESSAGE='{"tags": {"room": "living_room"},"fields": {"temp": "24.2", "hum": 40, "cp": "15i"}}'


docker run -it --rm efrecon/mqtt-client mosquitto_pub -d \
        -h "${BROKER}" \
        -p "${PORT}" \
        -t "${TOPIC}" \
        -m "${MESSAGE}" \
        -i "${CLIENT_ID}"
