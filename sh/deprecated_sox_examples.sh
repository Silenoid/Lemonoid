#! /usr/bin/env bash

sox -v 0.2 "./res/music/$musicdir/$namenoext.mp3" -c 1 -r 48000 "./res/music_baked/$musicdir/$namenoext.wav"
sox "temp/elevenlabs_result.mp3" -c 1 -r 48000 "temp/elevenlabs_result.wav"
sox  -v 0.3 ./temp/elevenlabs_music.wav -c 1 -r 48000 ./temp/elevenlabs_music_p.wav
soxi ./temp/elevenlabs_music_p.wav
sox -m ./temp/elevenlabs_result.wav ./temp/elevenlabs_music_p.wav ./temp/elevenlabs_voice_p.wav trim 0 "$(soxi -D ./temp/elevenlabs_result.wav)"
opusenc ./temp/elevenlabs_voice_p.wav ./temp/elevenlabs_voice.ogg
sox ./temp/downloaded.wav ./temp/downloaded_p.wav norm 4 highpass 1000
soxi ./temp/downloaded_p.wav
sox  -v 0.2 ./temp/music.wav -c 1 -r 48000 ./temp/music_p.wav
soxi ./temp/music_p.wav
sox -m ./temp/downloaded_p.wav ./temp/music_p.wav ./temp/output_rec.wav trim 0 "$(soxi -D ./temp/downloaded_p.wav)"
opusenc ./temp/output_rec.wav ./temp/output_rec.ogg