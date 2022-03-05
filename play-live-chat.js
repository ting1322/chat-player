const body = document.getElementsByTagName('body');
const video1 = document.getElementById('video1');
const chat_div = document.getElementById('live-chat');
const timestamp_div = document.getElementById('timestamp');
const chat_templ = document.getElementById('live-chat-item-template');
const timestamp_templ = document.getElementById('timestamp-template');
const chat_array = [];

function prettyFormatTime(timeInMs) {
    const timeInSec = Math.floor(timeInMs / 1000);
    const timeSec = timeInSec % 60;
    const timeInMinute = Math.floor(timeInSec / 60);
    const timeMinute = timeInMinute % 60;
    const timeInHour = Math.floor(timeInMinute / 60);
    var str = "";
    if (timeInHour < 10)
        str += "0";
    str += timeInHour + ":";
    if (timeMinute < 10)
        str += "0";
    str += timeMinute + ":";
    if (timeSec < 10)
        str += "0";
    str += timeSec;
    return str;
}

async function wait(timeInMs) {
    await new Promise(resolve => setTimeout(resolve, timeInMs));
}

async function create_chat_item(json_text) {
    if (json_text.length < 10)
        return;
    var json;
    try {
        json = JSON.parse(json_text);
    } catch (e) {
        console.log('parse json error: ' + json_text);
        return;
    }
    var node = chat_templ.cloneNode(true);
    node.removeAttribute('id');
    c_content = node.getElementsByClassName('c_content')[0];
    c_name = node.getElementsByClassName('c_name')[0];
    c_time = node.getElementsByClassName('c_time')[0];
    if (!('replayChatItemAction' in json))
        return;
    var timeInMs = -1;
    if ('videoOffsetTimeMsec' in json.replayChatItemAction)
        timeInMs = json.replayChatItemAction.videoOffsetTimeMsec;
    else if ('videoOffsetTimeMsec' in json)
        timeInMs = json.videoOffsetTimeMsec;
    else {
        console.log('error: not found videoOffsetTimeMsec');
        return;
    }
    c_time.setAttribute('time_in_ms', timeInMs);
    
    var hasText = false;
    for (const action of json.replayChatItemAction.actions) {
        if (!('addChatItemAction' in action)) {
            continue;
        }
        if (!('item' in action.addChatItemAction)) {
            continue;
        }

        const actionItem = action.addChatItemAction.item;
        if (!('liveChatTextMessageRenderer' in actionItem)) {
            continue;
        }
        if (!('message' in actionItem.liveChatTextMessageRenderer)) {
            continue;
        }
        if (!('runs' in actionItem.liveChatTextMessageRenderer.message)) {
            continue;
        }

        const runs = actionItem.liveChatTextMessageRenderer.message.runs;
        c_content.innerHTML = "";
        for (const run of runs) {
            if ('text' in run) {
                span = document.createElement('span');
                span.innerHTML = run.text;
                c_content.appendChild(span);
                hasText = true;
            }
            else if ('emoji' in run) {
                const emoji = run.emoji;
                var image_url = "";
                if ('image' in emoji && 'thumbnails' in emoji.image) {
                    thumbnails = emoji.image.thumbnails;
                    image_url = thumbnails[thumbnails.length-1].url;
                    img = document.createElement('img');
                    img.setAttribute('src', image_url);
                    img.setAttribute('class', 'emoji');
                    //img.setAttribute('width', '1em');
                    c_content.appendChild(img);
                    hasText = true;
                }
            } else {
                console.log('unknown message');
            }
        }
        if ('authorName' in actionItem.liveChatTextMessageRenderer) {
            const authorName = actionItem.liveChatTextMessageRenderer.authorName;
            if ('simpleText' in authorName) {
                c_name.innerHTML = authorName.simpleText;
            } else {
                console.log('unknown authorName');
            }
        }
    }
    
    if (hasText) {
        chat_array.push(timeInMs);
        c_time.innerHTML = prettyFormatTime(timeInMs);
        chat_div.appendChild(node);
        c_time.onclick = comment_time_click;
    }
}

async function init_js_from_embedded() {
    const pre = document.getElementById('live-chat-json-text');
    const text = pre.innerHTML;
    var json_lines = text.split(/\r?\n/);

    var idx = 0;
    for (const json_line of json_lines) {
        create_chat_item(json_line);
        idx += 1;
        // 由於 JSON.parse 需要時間，每 256 次就放開 thread 50ms
        // 讓其他UI能夠更新
        if ((idx & 0xFF) == 0) {
            await wait(50);
        }
    }
}

function init_setlist_from_embedded() {
    pre = document.getElementById('setlist-json-text');
    json_text = pre.innerHTML;
    if (json_text.length == 0)
        return;
    var json;
    try {
        json = JSON.parse(json_text);
    } catch (e) {
        console.log('parse json error: ' + json_text);
        return;
    }
    for (const chapter of json) {
        const node = timestamp_templ.cloneNode(true);
        node.removeAttribute('id');
        t_time = node.getElementsByClassName('t_time')[0];
        t_title = node.getElementsByClassName('t_title')[0];
        t_time.innerHTML = prettyFormatTime(chapter.time_in_ms);
        t_title.innerHTML = chapter.title;
        t_time.setAttribute('time_in_ms', chapter.time_in_ms);
        timestamp_div.appendChild(node)
        t_time.onclick = timestamp_click;
    }
}

function comment_time_click(eventArg) {
    const element = eventArg.currentTarget;
    const timeInMs = parseInt(element.getAttribute('time_in_ms'));
    video1.fastSeek(timeInMs / 1000);
}

function timestamp_click(eventArg) {
    const element = eventArg.currentTarget;
    timeInMs = parseInt(element.getAttribute('time_in_ms'));
    video1.fastSeek(timeInMs / 1000);
}

function sync_live_chat() {
    const current_time_in_ms = video1.currentTime * 1000;
    var i = 0;
    for(i=0; i<chat_array.length; i++) {
        if (chat_array[i] >= current_time_in_ms) {
            break;
        }
    }
    i = i - 5;
    if (i < 0) i=0;

    const div_y = chat_div.getBoundingClientRect().y;
    const child_y = chat_div.children[i].getBoundingClientRect().y;
    const scrollTop = chat_div.scrollTop;
    
    const scroll_target = child_y - div_y + scrollTop;
    chat_div.scrollTo({
        top: scroll_target,
        behavior: "smooth"
    });
}

function sync_live_chat_work() {
    if (!video1.paused) {
        sync_live_chat();
    }
    setTimeout(sync_live_chat_work, 500);
}

video1.onseeked = function() {
    if (video1.paused) {
        sync_live_chat();
    }
}

init_js_from_embedded();
init_setlist_from_embedded();

setTimeout(sync_live_chat_work, 500);
