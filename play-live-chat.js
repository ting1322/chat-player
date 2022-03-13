const body = document.getElementsByTagName('body');
const video1 = document.getElementById('video1');
const chat_div = document.getElementById('live-chat');
const timestamp_div = document.getElementById('timestamp');
const chat_templ = document.getElementById('live-chat-item-template');
const sc_templ = document.getElementById('live-chat-sc-template');
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
    
    var node;
    if (json.replayChatItemAction.actions.length > 0) {
        const action = json.replayChatItemAction.actions[0];
        if ('addChatItemAction' in action &&
            'item' in action.addChatItemAction) {

            const actionItem = action.addChatItemAction.item;
            if ('liveChatTextMessageRenderer' in actionItem) {
                node = render_liveChatTextMessage(actionItem.liveChatTextMessageRenderer, timeInMs);
            } else if ('liveChatPaidMessageRenderer' in actionItem) {
                node = render_liveChatPaidMessage(actionItem.liveChatPaidMessageRenderer, timeInMs);
            }
        }
    }
    
    if (node != null) {
        chat_array.push(timeInMs);
        chat_div.appendChild(node);
    }
}

function render_liveChatTextMessage(liveChatTextMessageRenderer, timeInMs)
{
    var hasText = false;
    var node = chat_templ.cloneNode(true);
    node.removeAttribute('id');
    c_content = node.getElementsByClassName('c_content')[0];
    c_name = node.getElementsByClassName('c_name')[0];
    c_time = node.getElementsByClassName('c_time')[0];
    c_time.setAttribute('time_in_ms', timeInMs);

    if (!('message' in liveChatTextMessageRenderer)) {
        return;
    }
    if (!('runs' in liveChatTextMessageRenderer.message)) {
        return;
    }

    const runs = liveChatTextMessageRenderer.message.runs;
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
    if ('authorName' in liveChatTextMessageRenderer) {
        const authorName = liveChatTextMessageRenderer.authorName;
        if ('simpleText' in authorName) {
            c_name.innerHTML = authorName.simpleText;
        } else {
            console.log('unknown authorName');
        }
    }
    if (hasText)
    {
        c_time.innerHTML = prettyFormatTime(timeInMs);
        c_time.onclick = comment_time_click;
        return node;
    }
    return null;
}

function render_liveChatPaidMessage(liveChatPaidMessageRenderer, timeInMs)
{
    var hasText = false;
    var node = sc_templ.cloneNode(true);
    node.removeAttribute('id');
    c_header = node.getElementsByClassName('header')[0];
    c_text = node.getElementsByClassName('text')[0];
    c_name = node.getElementsByClassName('name')[0];
    c_paid = node.getElementsByClassName('paid')[0];
    c_time = node.getElementsByClassName('c_time')[0];
    c_time.setAttribute('time_in_ms', timeInMs);

    if (!('message' in liveChatPaidMessageRenderer)) {
        return;
    }
    if (!('runs' in liveChatPaidMessageRenderer.message)) {
        return;
    }

    c_header.style.backgroundColor = toColor(liveChatPaidMessageRenderer.headerBackgroundColor);
    c_header.style.color = toColor(liveChatPaidMessageRenderer.headerTextColor);
    c_text.style.backgroundColor = toColor(liveChatPaidMessageRenderer.bodyBackgroundColor);
    c_text.style.color = toColor(liveChatPaidMessageRenderer.bodyTextColor);

    if ('purchaseAmountText' in liveChatPaidMessageRenderer &&
        'simpleText' in liveChatPaidMessageRenderer.purchaseAmountText) {
        c_paid.innerHTML = liveChatPaidMessageRenderer.purchaseAmountText.simpleText;
    }

    const runs = liveChatPaidMessageRenderer.message.runs;
    c_text.innerHTML = "";
    for (const run of runs) {
        if ('text' in run) {
            span = document.createElement('span');
            span.innerHTML = run.text;
            c_text.appendChild(span);
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
                c_text.appendChild(img);
                hasText = true;
            }
        } else {
            console.log('unknown message');
        }
    }
    if ('authorName' in liveChatPaidMessageRenderer) {
        const authorName = liveChatPaidMessageRenderer.authorName;
        if ('simpleText' in authorName) {
            c_name.innerHTML = authorName.simpleText;
        } else {
            console.log('unknown authorName');
        }
    }
    if (hasText)
    {
        c_time.innerHTML = prettyFormatTime(timeInMs);
        c_time.onclick = comment_time_click;
        return node;
    }
    return null;
}

function toColor(num) {
    num >>>= 0;
    var b = num & 0xFF,
        g = (num & 0xFF00) >>> 8,
        r = (num & 0xFF0000) >>> 16,
        a = ( (num & 0xFF000000) >>> 24 ) / 255 ;
    return "rgba(" + [r, g, b, a].join(",") + ")";
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
