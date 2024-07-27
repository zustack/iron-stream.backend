#!/bin/bash

if [ "$#" -ne 2 ]; then
    echo "Uso: $0 <input_path> <output_directory>"
    exit 1
fi

input_path=$1
output_dir=$2

ffmpeg -i "$input_path" \
-filter_complex \
"[0:v]split=4[v1][v2][v3][v4]; \
[v1]scale=w=640:h=360:force_original_aspect_ratio=decrease,pad=640:360:(ow-iw)/2:(oh-ih)/2[v1out]; \
[v2]scale=w=854:h=480:force_original_aspect_ratio=decrease,pad=854:480:(ow-iw)/2:(oh-ih)/2[v2out]; \
[v3]scale=w=1280:h=720:force_original_aspect_ratio=decrease,pad=1280:720:(ow-iw)/2:(oh-ih)/2[v3out]; \
[v4]scale=w=1920:h=1080:force_original_aspect_ratio=decrease,pad=1920:1080:(ow-iw)/2:(oh-ih)/2[v4out]" \
-map "[v1out]" -map 0:a -map "[v2out]" -map 0:a -map "[v3out]" -map 0:a -map "[v4out]" -map 0:a \
-c:v libx264 -crf 23 -preset slower -c:a aac -ar 48000 \
-b:v:0 800k -maxrate:v:0 856k -bufsize:v:0 1200k -b:a:0 96k \
-b:v:1 1400k -maxrate:v:1 1498k -bufsize:v:1 2100k -b:a:1 128k \
-b:v:2 2800k -maxrate:v:2 2996k -bufsize:v:2 4200k -b:a:2 128k \
-b:v:3 5000k -maxrate:v:3 5350k -bufsize:v:3 7500k -b:a:3 192k \
-var_stream_map "v:0,a:0,name:360p v:1,a:1,name:480p v:2,a:2,name:720p v:3,a:3,name:1080p" \
-keyint_min 48 -g 48 -sc_threshold 0 -hls_time 4 \
-hls_playlist_type event -hls_list_size 0 -f hls \
-hls_flags independent_segments+delete_segments \
-master_pl_name "master.m3u8" \
"$output_dir/master-%v.m3u8"
