#!/bin/bash
filename=$1
extension="${filename##*.}"
filename="${filename%.*}"
( ffmpeg -i "$1" -f s16le -ar 48000 -ac 2 pipe:1 | dca > "$filename.dca" )
# if [ $extension == 'opus' ]; then
#     echo fuck
#     ( cat "$1" | dca > "$filename.dca" )
#     echo shit
# else
#     echo ing shit
#     ffmpeg -i "$1" -f s16le -ar 48000 -ac 2 pipe:1 | dca > "$filename.dca"
# fi

( rm "$1" )
