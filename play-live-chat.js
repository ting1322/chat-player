const body = document.getElementsByTagName('body');
const video1 = document.getElementById('video1');
const chat_div = document.getElementById('live-chat');
const timestamp_div = document.getElementById('timestamp');
const chat_templ = document.getElementById('live-chat-item-template');
const sc_templ = document.getElementById('live-chat-sc-template');
const timestamp_templ = document.getElementById('timestamp-template');
const chat_array = [];

function resizeChatDiv() {
    const height = video1.getClientRects()[0].height
    chat_div.style.height = (height-5) + "px";
}

window.addEventListener("resize", resizeChatDiv);
video1.addEventListener("resize", resizeChatDiv);

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
    timeInMs = Math.max(timeInMs, 0)

    var o;
    if (json.replayChatItemAction.actions.length > 0) {
        const action = json.replayChatItemAction.actions[0];
        if ('addChatItemAction' in action &&
            'item' in action.addChatItemAction) {

            const actionItem = action.addChatItemAction.item;
            if ('liveChatTextMessageRenderer' in actionItem) {
                o = render_liveChatTextMessage(actionItem.liveChatTextMessageRenderer, timeInMs);
            } else if ('liveChatPaidMessageRenderer' in actionItem) {
                o = render_liveChatPaidMessage(actionItem.liveChatPaidMessageRenderer, timeInMs);
            } else if ('liveChatMembershipItemRenderer' in actionItem) {
                o = render_liveChatPaidMessage(actionItem.liveChatMembershipItemRenderer, timeInMs);
            } else if ('liveChatPaidStickerRenderer' in actionItem) {
                o = render_liveChatSticker(actionItem.liveChatPaidStickerRenderer, timeInMs);
            } else if ('liveChatSponsorshipsGiftPurchaseAnnouncementRenderer' in actionItem) {
                o = render_liveChatGift(actionItem.liveChatSponsorshipsGiftPurchaseAnnouncementRenderer, timeInMs);
            }
        }
    }

    if (o != null) {
        chat_array.push(timeInMs);
        chat_div.appendChild(o.node);
    }
}

function render_liveChatTextMessage(liveChatTextMessageRenderer, timeInMs)
{
    var hasText = false;
    var o = newChatTextNode(timeInMs);

    if (!('message' in liveChatTextMessageRenderer)) {
        return;
    }
    if (!('runs' in liveChatTextMessageRenderer.message)) {
        return;
    }

    const runs = liveChatTextMessageRenderer.message.runs;
    o.c_content.innerHTML = "";
    for (const run of runs) {
        if ('text' in run) {
            span = document.createElement('span');
            span.innerHTML = run.text;
            o.c_content.appendChild(span);
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
                o.c_content.appendChild(img);
                hasText = true;
            }
        } else {
            console.log('unknown message');
        }
    }
    if ('authorName' in liveChatTextMessageRenderer) {
        const authorName = liveChatTextMessageRenderer.authorName;
        if ('simpleText' in authorName) {
            o.c_name.innerHTML = authorName.simpleText;
        } else {
            console.log('unknown authorName');
        }
    }
    if (hasText)
    {
        return o;
    }
    return null;
}

function newChatTextNode(timeInMs)
{
    const o = {};
    o.node = chat_templ.cloneNode(true);
    o.node.removeAttribute('id');
    o.c_content = o.node.getElementsByClassName('c_content')[0];
    o.c_name = o.node.getElementsByClassName('c_name')[0];
    o.c_time = o.node.getElementsByClassName('c_time')[0];
    o.c_time.setAttribute('time_in_ms', timeInMs);
    o.c_time.innerHTML = prettyFormatTime(timeInMs);
    o.c_time.onclick = comment_time_click;
    return o;
}

class SuperChatNode {
    constructor(timeInMs) {
        this.node = sc_templ.cloneNode(true);
        this.node.removeAttribute('id');
        this.c_header = this.node.getElementsByClassName('header')[0];
        this.c_text = this.node.getElementsByClassName('text')[0];
        this.c_name = this.node.getElementsByClassName('name')[0];
        this.c_paid = this.node.getElementsByClassName('paid')[0];
        this.c_time = this.node.getElementsByClassName('c_time')[0];
        this.c_time.setAttribute('time_in_ms', timeInMs);
        this.c_time.innerHTML = prettyFormatTime(timeInMs);
        this.c_time.onclick = comment_time_click;
    }

    setStickerMode() {
        this.node.setAttribute("class", "live-chat-sticker");
    }
}

function render_liveChatPaidMessage(liveChatPaidMessageRenderer, timeInMs)
{
    var o = new SuperChatNode(timeInMs)

    if ('headerBackgroundColor' in liveChatPaidMessageRenderer) {
        o.c_header.style.backgroundColor = toColor(liveChatPaidMessageRenderer.headerBackgroundColor);
        o.c_header.style.color = toColor(liveChatPaidMessageRenderer.headerTextColor);
        o.c_text.style.backgroundColor = toColor(liveChatPaidMessageRenderer.bodyBackgroundColor);
        o.c_text.style.color = toColor(liveChatPaidMessageRenderer.bodyTextColor);
    } else {
        o.c_header.style.backgroundColor = "rgb(10, 128, 67)";
        o.c_header.style.color = "rgb(0,0,0)";
        o.c_text.style.backgroundColor = "rgb(15, 157, 88)";
        o.c_text.style.color = "rgb(0,0,0)"
    }

    if ('purchaseAmountText' in liveChatPaidMessageRenderer &&
        'simpleText' in liveChatPaidMessageRenderer.purchaseAmountText) {
        o.c_paid.innerHTML = liveChatPaidMessageRenderer.purchaseAmountText.simpleText;
    } else if ('headerSubtext' in liveChatPaidMessageRenderer &&
               'simpleText' in liveChatPaidMessageRenderer.headerSubtext) {
        o.c_paid.innerHTML = liveChatPaidMessageRenderer.headerSubtext.simpleText;
    }

    if ('message' in liveChatPaidMessageRenderer &&
        'runs' in liveChatPaidMessageRenderer.message) {
        const runs = liveChatPaidMessageRenderer.message.runs;
        o.c_text.innerHTML = "";
        for (const run of runs) {
            if ('text' in run) {
                span = document.createElement('span');
                span.innerHTML = run.text;
                o.c_text.appendChild(span);
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
                    o.c_text.appendChild(img);
                }
            } else {
                console.log('unknown message');
            }
        }
    }
    if ('authorName' in liveChatPaidMessageRenderer) {
        const authorName = liveChatPaidMessageRenderer.authorName;
        if ('simpleText' in authorName) {
            o.c_name.innerHTML = authorName.simpleText;
        } else {
            console.log('unknown authorName');
        }
    }
    return o;
}

function render_liveChatSticker(liveChatPaidStickerRenderer, timeInMs)
{
    var o = new SuperChatNode(timeInMs);
    o.c_header.style.backgroundColor = toColor(liveChatPaidStickerRenderer.moneyChipBackgroundColor);
    o.c_header.style.color = toColor(liveChatPaidStickerRenderer.moneyChipTextColor);
    o.c_text.style.backgroundColor = toColor(liveChatPaidStickerRenderer.moneyChipBackgroundColor);
    o.c_text.style.color = toColor(liveChatPaidStickerRenderer.moneyChipTextColor);
    o.setStickerMode();

    if ('purchaseAmountText' in liveChatPaidStickerRenderer &&
        'simpleText' in liveChatPaidStickerRenderer.purchaseAmountText) {
        o.c_paid.innerHTML = liveChatPaidStickerRenderer.purchaseAmountText.simpleText;
    }

    if ('sticker' in liveChatPaidStickerRenderer &&
        'thumbnails' in liveChatPaidStickerRenderer.sticker) {
        const thumbnails = liveChatPaidStickerRenderer.sticker.thumbnails;
        image_url = thumbnails[thumbnails.length-1].url;
        img = document.createElement('img');
        img.setAttribute('src', image_url);
        img.setAttribute('class', 'sticker');
        //img.setAttribute('width', '1em');
        o.c_text.appendChild(img);
    }
    if ('authorName' in liveChatPaidStickerRenderer) {
        const authorName = liveChatPaidStickerRenderer.authorName;
        if ('simpleText' in authorName) {
            o.c_name.innerHTML = authorName.simpleText;
        } else {
            console.log('unknown authorName');
        }
    }
    return o;
}

function render_liveChatGift(liveChatSponsorshipsGiftPurchaseAnnouncementRenderer, timeInMs)
{
    if (!('header' in liveChatSponsorshipsGiftPurchaseAnnouncementRenderer))
        return;
    const header = liveChatSponsorshipsGiftPurchaseAnnouncementRenderer.header;
    if (!('liveChatSponsorshipsHeaderRenderer' in header))
        return;
    const renderer = header.liveChatSponsorshipsHeaderRenderer;

    var o = new SuperChatNode(timeInMs);
    o.c_header.style.backgroundColor = "rgb(10, 128, 67)";
    o.c_header.style.color = "rgb(0,0,0)";
    o.c_text.innerHTML = "";

    if ('authorName' in renderer) {
        const authorName = renderer.authorName;
        if ('simpleText' in authorName) {
            o.c_name.innerHTML = authorName.simpleText;
        } else {
            console.log('unknown authorName');
        }
    }

    if ('primaryText' in renderer) {
        const runs = renderer.primaryText.runs;
        o.c_paid.innerHTML = "";
        for (const run of runs) {
            if ('text' in run) {
                span = document.createElement('span');
                span.innerHTML = run.text;
                o.c_paid.appendChild(span);
                hasText = true;
            } else {
                console.log('unknown message');
            }
        }
    }
    return o
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
    video1.currentTime = timeInMs / 1000;
}

function timestamp_click(eventArg) {
    const element = eventArg.currentTarget;
    timeInMs = parseInt(element.getAttribute('time_in_ms'));
    video1.currentTime = timeInMs / 1000;
}

function sync_live_chat() {
    const current_time_in_ms = video1.currentTime * 1000;
    var i = 0;
    for(i=0; i<chat_array.length; i++) {
        if (chat_array[i] >= current_time_in_ms) {
            break;
        }
    }

    if (i == chat_array.length) {
        i = chat_array.length - 1;
    }

    const div_y = chat_div.getBoundingClientRect().bottom - 100;
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
