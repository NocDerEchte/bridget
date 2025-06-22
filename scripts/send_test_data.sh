#!/usr/bin/env bash

BROKER='192.168.178.61'
PORT='1883'
CLIENT_ID='test'
TOPIC='test-db/test'
MESSAGE='{"tags": {"room": "living_room"},"fields": {"temp": "24.2", "hum": 40, "cp": "15i"}}'


docker run -it --rm efrecon/mqtt-client mosquitto_pub -d \
        -h "${BROKER}" \
        -p "${PORT}" \
        -t "${TOPIC}" \
        -m "${MESSAGE}" \
        -i "${CLIENT_ID}"
