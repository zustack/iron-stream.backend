#!/bin/bash

if [ "$#" -ne 1 ]; then
    echo "Uso: $0 <ruta-al-archivo-de-video>"
    exit 1
fi

VIDEO_PATH="$1"

if [ ! -f "$VIDEO_PATH" ]; then
    echo "El archivo $VIDEO_PATH no existe."
    exit 1
fi

ffprobe -v error -select_streams v:0 -show_entries stream=duration -of csv=p=0 "$VIDEO_PATH" | awk '{print int($1)}'
