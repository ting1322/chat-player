package cplayer

import (
	"encoding/json"
	"log"
	"path/filepath"
	"strconv"
	"strings"
)

type jmap = map[string]any

func preprocessJson(option *Option, down ImgDownloader, jsonText, outDir string) (string, error) {
	if option.NoDownloadPic || len(jsonText) < 10 {
		return jsonText, nil
	}

	var jsonmap jmap
	json.Unmarshal([]byte(jsonText), &jsonmap)
	replayChat := jsonmap["replayChatItemAction"].(jmap)
	actions, exist := replayChat["actions"]
	if !exist {
		return jsonText, nil
	}
	for _, action := range actions.([]any) {
		_, exist = action.(jmap)["addChatItemAction"]
		if exist {
			var render jmap
			addChatItem := action.(jmap)["addChatItemAction"].(jmap)["item"].(jmap)
			if _, exist := addChatItem["liveChatTextMessageRenderer"]; exist {
				render = addChatItem["liveChatTextMessageRenderer"].(jmap)
			} else if _, exist := addChatItem["liveChatPaidMessageRenderer"]; exist {
				render = addChatItem["liveChatPaidMessageRenderer"].(jmap)
			} else if _, exist := addChatItem["liveChatViewerEngagementMessageRenderer"]; exist {
				// system message, don't care
			} else if _, exist := addChatItem["liveChatMembershipItemRenderer"]; exist {
				render = addChatItem["liveChatMembershipItemRenderer"].(jmap)
			} else if _, exist := addChatItem["liveChatPaidStickerRenderer"]; exist {
				// super stick
				render = addChatItem["liveChatPaidStickerRenderer"].(jmap)
			} else if _, exist := addChatItem["liveChatSponsorshipsGiftPurchaseAnnouncementRenderer"]; exist {
				// membership gift
			} else if _, exist := addChatItem["liveChatSponsorshipsGiftRedemptionAnnouncementRenderer"]; exist {
				// membership gift ??
				// unimplement
			} else {
				log.Printf("unknown addChatItemAction.item node: %v", jsonText)
				continue
			}

			if n, exist := render["authorBadges"]; exist {
				authorBadges := n.([]any)
				if len(authorBadges) > 0 {
					liveChatAuthorBadgeRenderer := authorBadges[0].(jmap)["liveChatAuthorBadgeRenderer"].(jmap)
					if n, exist := liveChatAuthorBadgeRenderer["customThumbnail"]; exist {
						customThumbnail := n.(jmap)
						if n, exist := customThumbnail["thumbnails"]; exist {
							thumbnails := n.([]any)
							for _, thumbnail := range thumbnails {
								image_url := thumbnail.(jmap)["url"].(string)
								filename := hashUrlFilename(image_url)
								filename = down.Download(filepath.Join(outDir, filename), image_url)
								filename, _ = filepath.Rel(outDir, filename)
								thumbnail.(jmap)["url"] = filename
							}
						}

					}
				}
			}

			if _, exist := render["message"]; exist {
				runs := render["message"].(jmap)["runs"].([]any)

				for _, run := range runs {
					if _, exist := run.(jmap)["emoji"]; exist {
						emoji := run.(jmap)["emoji"].(jmap)
						thumbnails := emoji["image"].(jmap)["thumbnails"].([]any)
						for _, thumbnail := range thumbnails {
							image_url := thumbnail.(jmap)["url"].(string)
							filename := hashUrlFilename(image_url)
							filename = down.Download(filepath.Join(outDir, filename), image_url)
							filename, _ = filepath.Rel(outDir, filename)
							thumbnail.(jmap)["url"] = filename
						}
					}
				}
			} else if _, exist := render["sticker"]; exist {
				thumbnails := render["sticker"].(jmap)["thumbnails"].([]any)
				for _, thumbnail := range thumbnails {
					image_url := thumbnail.(jmap)["url"].(string)
					filename := hashUrlFilename(image_url)
					filename = down.Download(filepath.Join(outDir, filename), image_url)
					filename, _ = filepath.Rel(outDir, filename)
					thumbnail.(jmap)["url"] = filename
				}
			}
		}
	}
	if _, exist := replayChat["videoOffsetTimeMsec"]; exist && option.TimeOffsetInSec > 0 {
		timeInMs, err := strconv.Atoi(replayChat["videoOffsetTimeMsec"].(string))
		if err == nil {
			timeInMs += option.TimeOffsetInSec * 1000
			replayChat["videoOffsetTimeMsec"] = strconv.Itoa(timeInMs)
		}
	} else if _, exist := jsonmap["videoOffsetTimeMsec"]; exist && option.TimeOffsetInSec > 0 {
		timeInMs, err := strconv.Atoi(jsonmap["videoOffsetTimeMsec"].(string))
		if err == nil {
			timeInMs += option.TimeOffsetInSec * 1000
			jsonmap["videoOffsetTimeMsec"] = strconv.Itoa(timeInMs)
		}
	}

	jsondata, err := json.Marshal(jsonmap)
	if err != nil {
		return jsonText, nil
	}
	text := string(jsondata)
	text = strings.ReplaceAll(text, "</", "\\u003C/")
	return text, nil
}
