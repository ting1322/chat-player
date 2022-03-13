#! /usr/bin/python3
# -*- coding:utf-8; -*-

from hashlib import blake2b
import argparse
import json
import os
import pathlib
import re
import requests
import shutil
import sys
import urllib.parse

def main():
    cmd_parser = argparse.ArgumentParser(description='generate a html to play video with live-chat.json')
    cmd_parser.add_argument('video', type=pathlib.Path, help='video file (webm or mp4)')
    cmd_parser.add_argument('-c', '--chat-json', type=pathlib.Path, help='live chat json file (download by yt-dlp)')
    cmd_parser.add_argument('-s', '--set-list', type=pathlib.Path, help='時間軸 txt 檔')
    cmd_parser.add_argument('-o', '--output', type=pathlib.Path, help='output html file')
    cmd_parser.add_argument('--no-download-pic', action='store_true', help='不要把聊天室貼圖抓下來 (每次開網頁使用youtube檔案)')
    cmd = cmd_parser.parse_args()

    video_filename = cmd.video
    chat_json_filename = cmd.chat_json
    output_filename = cmd.output
    setlist_filename = cmd.set_list
    title = os.path.basename(os.path.splitext(video_filename)[0])

    if chat_json_filename is None:
        chat_json_filename = os.path.splitext(cmd.video)[0] + '.live_chat.json'
    if output_filename is None:
        output_filename = os.path.basename(cmd.video) + '.htm'

    if not os.path.exists(video_filename):
        print ('not found input file: ' + video_filename)
        sys.exit(-1)
    if not os.path.exists(chat_json_filename):
        print ('not found input file: ' + chat_json_filename)
        sys.exit(-1)

    app_dir = os.path.normpath(os.path.split(__file__)[0])
    video_dir = os.path.normpath(os.path.split(video_filename)[0])
    out_dir = os.path.normpath(os.path.split(os.path.abspath(output_filename))[0])

    video_filename = os.path.relpath(video_filename, start=out_dir)

    if setlist_filename is None:
        setlist_filename = os.path.join(video_dir, 'set-list.txt')
        if not os.path.exists(setlist_filename):
            setlist_filename = None

    if not setlist_filename is None \
       and not os.path.exists(setlist_filename):
        print ('not found input file: ' + setlist_filename)
        sys.exit(-1)

    video_type = str()
    if str(video_filename).endswith('mp4'):
        video_type = 'video/webm'
    elif str(video_filename).endswith('webm'):
        video_type = 'video/webm'
    else:
        print ('not support video type, only support mp4 and webm')

    setlist_json_text = convert_setlist_to_json(setlist_filename)

    template_htm = os.path.join(app_dir, 'template.htm.in')
    with open(template_htm, 'r') as in_file, \
         open(chat_json_filename, 'r', encoding='utf-8') as inline_file, \
         open(output_filename, 'w') as out_file:
        live_chat_json_lines = '\n'
        for line in inline_file.readlines():
            live_chat_json_lines += preprocess_json(line, cmd, out_dir)
        for line in in_file.readlines():
            text = line
            text = text.replace('{{video}}', urllib.parse.quote(str(video_filename)))
            text = text.replace('{{title}}', title)
            text = text.replace('{{video-type}}', video_type)
            text = text.replace('{{live-chat-json}}', live_chat_json_lines)
            text = text.replace('{{setlist-json}}', setlist_json_text)
            out_file.write(text)
    if out_dir != app_dir:
        shutil.copy(os.path.join(app_dir, 'play-live-chat.js'),
                    os.path.join(out_dir, 'play-live-chat.js'))

def convert_setlist_to_json(setlist_filename):
    if setlist_filename is None:
        return ""
    arr = []
    with open(setlist_filename, 'r') as f:
        for line in f.readlines():
            x = re.match(r"(\d+\. )?(\d{1,2}):(\d{2}):(\d{2}) (.+)", line)
            if x is None:
                continue
            hour = int(x.group(2))
            minute = int(x.group(3))
            second = int(x.group(4))
            title = x.group(5)
            totla_ms = ((((hour * 60) + minute) * 60) + second) * 1000
            obj = { 'time_in_ms': totla_ms,
                    'title': title }
            arr.append(obj)
    return json.dumps(arr)

def preprocess_json(line, cmd, out_dir):
    if cmd.no_download_pic:
        return line

    if len(line) < 10:
        return line

    if not os.path.exists(os.path.join(out_dir, 'images')):
        os.mkdir(os.path.join(out_dir, 'images'))

    js_root = json.loads(line)
    replayChat = js_root['replayChatItemAction']
    if not 'actions' in replayChat:
        return line
    for action in replayChat['actions']:
        if 'addChatItemAction' in action:
            addChatItem = action['addChatItemAction']['item']
            render = None
            if 'liveChatTextMessageRenderer' in addChatItem:
                render = addChatItem['liveChatTextMessageRenderer']
            elif 'liveChatPaidMessageRenderer' in addChatItem:
                render = addChatItem['liveChatPaidMessageRenderer']
            if render is None:
                continue
            if not 'message' in render:
                continue
            runs = render['message']['runs']
            for run in runs:
                if 'emoji' in run:
                    emoji = run['emoji']
                    thumbnails = emoji['image']['thumbnails']
                    for thumbnail in thumbnails:
                        image_url = thumbnail['url']
                        h = blake2b(digest_size=20)
                        h.update(image_url.encode('utf-8'))
                        filename = 'images/emoji-{0}'.format(h.hexdigest())
                        if str(image_url).endswith('.svg'):
                            filename += '.svg'
                        else:
                            filename += '.png'
                        http_download_image(os.path.join(out_dir, filename), image_url)
                        thumbnail['url'] = filename

    text = json.dumps(js_root) + '\n'
    text = text.replace('</', '\\u003C/')
    return text

def http_download_image(filename, image_url):
    if os.path.exists(filename):
        return
    try:
        print ('download emoji: ' + image_url)
        with open(filename, "wb") as f:
            r = requests.get(image_url)
            f.write(r.content)
            print ('save to: ' + filename)
    except:
        print ('靠北，下載表情符號失敗了: ' + image_url)

main()
