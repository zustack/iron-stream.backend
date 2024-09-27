#!/bin/bash

if [ "$#" -ne 2 ]; then
    echo "Uso: $0 <input_path> <output_directory>"
    exit 1
fi

input_path=$1
output_dir=$2

# Mantener la resolución original del video
ffmpeg -i "$input_path" \
-map 0:v -map 0:a \
-c:v copy -c:a copy \
-hls_time 4 \
-hls_playlist_type event -hls_list_size 0 -f hls \
-hls_flags independent_segments+delete_segments \
-master_pl_name "master.m3u8" \
"$output_dir/master.m3u8"
