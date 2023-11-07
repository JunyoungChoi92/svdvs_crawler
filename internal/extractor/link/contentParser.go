package link

import (
	"fmt"
	"regexp"
	"strings"
)

type ContentParser struct {
	content string
}

var tagList = []string{
	"!DOCTYPE", "a", "abbr", "acronym", "address", "applet", "area", "article", "aside", "audio", "b", "base", "basefont", "bdi", "bdo", "big", "blockquote", "body", "br", "button", "canvas", "caption", "center", "cite", "code", "col",
	"colgroup", "data", "datalist", "dd", "del", "details", "dfn", "dialog", "dir", "div", "dl", "dt", "em", "embed", "fieldset", "figcaption", "figure", "font", "footer", "form", "frame", "frameset", "h1", "h2", "h3", "h4", "h5", "h6",
	"head", "header", "hr", "html", "i", "iframe", "img", "input", "ins", "kbd", "label", "legend", "li", "link", "main", "map", "mark", "meta", "meter", "nav", "noframes", "noscript", "object", "ol", "optgroup", "option", "output", "p",
	"param", "picture", "pre", "progress", "q", "rp", "rt", "ruby", "s", "samp", "section", "select", "small", "source", "span", "strike", "strong", "sub", "summary", "sup", "svg", "table", "tbody", "td", "template", "textarea", "tfoot",
	"th", "thead", "time", "title", "tr", "track", "tt", "u", "ul", "var", "video", "wbr",
}

func (contentParser *ContentParser) Content() string {
	return contentParser.content
}

func (contentParser *ContentParser) SetContent(content string) {
	contentParser.content = content
}

func (contentParser *ContentParser) extractURLsFromContent(unEscapeHtml string) []string {
	contentParser.SetContent(unEscapeHtml)
	contentParser.whiteReg().styleReg().highlightReg().scriptReg().emojiReg()

	tagStr := strings.Builder{}
	for _, d := range tagList {
		tagStr.WriteString(fmt.Sprintf("<%s(.*?)>|</%s>|", d, d))
	}
	tagStr.WriteString(`[0-9]|!-~|[\x{1F600}-\x{1F6FF}|[\x{2600}-\x{26FF}]`)
	reg := regexp.MustCompile(tagStr.String())

	contentParser.SetContent(reg.ReplaceAllString(contentParser.Content(), ""))

	linkReg := regexp.MustCompile(`(http|https)?://(www.)?([0-9a-zA-Z]?)+.([0-9a-zA-Z]?)+`)
	contentParser.SetContent(linkReg.ReplaceAllString(contentParser.Content(), ""))

	return linkReg.FindAllString(contentParser.Content(), -1)
}

// content replace all whitespace
func (contentParser *ContentParser) whiteReg() *ContentParser {
	whiteReg := regexp.MustCompile(`\s|&nbsp;`)
	contentParser.SetContent(whiteReg.ReplaceAllString(contentParser.Content(), ""))
	return contentParser
}

// replace all style tag inner css data
func (contentParser *ContentParser) styleReg() *ContentParser {
	styleReg := regexp.MustCompile("<style(.*?)>(.*?)</style>")
	contentParser.SetContent(styleReg.ReplaceAllString(contentParser.Content(), ""))
	return contentParser
}

// replace all highlight
func (contentParser *ContentParser) highlightReg() *ContentParser {
	highlightReg := regexp.MustCompile("<!--(.*?)-->")
	contentParser.SetContent(highlightReg.ReplaceAllString(contentParser.Content(), ""))
	return contentParser
}

// replace all script tag inner script data
func (contentParser *ContentParser) scriptReg() *ContentParser {
	scriptReg := regexp.MustCompile("<script(.*?)>(.*?)</script>")
	contentParser.SetContent(scriptReg.ReplaceAllString(contentParser.Content(), ""))
	return contentParser
}

// replace all emoji in content data
func (contentParser *ContentParser) emojiReg() *ContentParser {
	reg := "[0-9#*]️?⃣|[©®‼⁉™ℹ↔-↙↩↪⌨⏏⏭-⏯⏱⏲⏸-⏺Ⓜ▪▫▶◀◻◼☀-☄☎☑☘☠☢☣☦☪☮☯☸-☺♀♂♟♠♣♥♦♨♻♾⚒⚔-⚗⚙⚛⚜⚠⚧⚰⚱⛈⛏⛑⛓⛩⛰⛱⛴⛷⛸✂✈✉✏✒✔✖✝✡✳✴❄❇❣➡⤴⤵⬅-⬇〰〽㊗㊙🅰🅱🅾🅿🈂🈷🌡🌤-🌬🌶🍽🎖🎗🎙-🎛🎞🎟🏍🏎🏔-🏟🏵🏷🐿📽🕉🕊🕯🕰🕳🕶-🕹🖇🖊-🖍🖥🖨🖱🖲🖼🗂-🗄🗑-🗓🗜-🗞🗡🗣🗨🗯🗳🗺🛋🛍-🛏🛠-🛥🛩🛰🛳]️?|[☝✌✍🕴🖐][️🏻-🏿]?|[⛹🏋🏌🕵](?:‍[♀♂]️?|[️🏻-🏿](?:‍[♀♂]️?)?)?|[✊✋🎅🏂🏇👂👃👆-👐👦👧👫-👭👲👴-👶👸👼💃💅💏💑💪🕺🖕🖖🙌🙏🛀🛌🤌🤏🤘-🤟🤰-🤴🤶🥷🦵🦶🦻🧒🧓🧕🫃-🫅🫰🫲-🫶][🏻-🏿]?|❤(?:‍[🔥🩹]|️(?:‍[🔥🩹])?)?|🇦[🇨-🇬🇮🇱🇲🇴🇶-🇺🇼🇽🇿]|🇧[🇦🇧🇩-🇯🇱-🇴🇶-🇹🇻🇼🇾🇿]|🇨[🇦🇨🇩🇫-🇮🇰-🇵🇷🇺-🇿]|🇩[🇪🇬🇯🇰🇲🇴🇿]|🇪[🇦🇨🇪🇬🇭🇷-🇺]|🇫[🇮-🇰🇲🇴🇷]|🇬[🇦🇧🇩-🇮🇱-🇳🇵-🇺🇼🇾]|🇭[🇰🇲🇳🇷🇹🇺]|🇮[🇨-🇪🇱-🇴🇶-🇹]|🇯[🇪🇲🇴🇵]|🇰[🇪🇬-🇮🇲🇳🇵🇷🇼🇾🇿]|🇱[🇦-🇨🇮🇰🇷-🇻🇾]|🇲[🇦🇨-🇭🇰-🇿]|🇳[🇦🇨🇪-🇬🇮🇱🇴🇵🇷🇺🇿]|🇴🇲|🇵[🇦🇪-🇭🇰-🇳🇷-🇹🇼🇾]|🇶🇦|🇷[🇪🇴🇸🇺🇼]|🇸[🇦-🇪🇬-🇴🇷-🇹🇻🇽-🇿]|🇹[🇦🇨🇩🇫-🇭🇯-🇴🇷🇹🇻🇼🇿]|🇺[🇦🇬🇲🇳🇸🇾🇿]|🇻[🇦🇨🇪🇬🇮🇳🇺]|🇼[🇫🇸]|🇽🇰|🇾[🇪🇹]|🇿[🇦🇲🇼]|[🏃🏄🏊👮👰👱👳👷💁💂💆💇🙅-🙇🙋🙍🙎🚣🚴-🚶🤦🤵🤷-🤹🤽🤾🦸🦹🧍-🧏🧔🧖-🧝](?:‍[♀♂]️?|[🏻-🏿](?:‍[♀♂]️?)?)?|🏳(?:‍(?:⚧️?|🌈)|️(?:‍(?:⚧️?|🌈))?)?|🏴(?:‍☠️?|󠁧󠁢(?:󠁥󠁮󠁧|󠁳󠁣󠁴|󠁷󠁬󠁳)󠁿)?|🐈(?:‍⬛)?|🐕(?:‍🦺)?|🐻(?:‍❄️?)?|👁(?:‍🗨️?|️(?:‍🗨️?)?)?|👨(?:‍(?:[⚕⚖✈]️?|❤️?‍(?:💋‍)?👨|👦(?:‍👦)?|👧(?:‍[👦👧])?|[👨👩]‍(?:👦(?:‍👦)?|👧(?:‍[👦👧])?)|[🌾🍳🍼🎓🎤🎨🏫🏭💻💼🔧🔬🚀🚒🦯-🦳🦼🦽])|🏻(?:‍(?:[⚕⚖✈]️?|❤️?‍(?:💋‍)?👨[🏻-🏿]|🤝‍👨[🏼-🏿]|[🌾🍳🍼🎓🎤🎨🏫🏭💻💼🔧🔬🚀🚒🦯-🦳🦼🦽]))?|🏼(?:‍(?:[⚕⚖✈]️?|❤️?‍(?:💋‍)?👨[🏻-🏿]|🤝‍👨[🏻🏽-🏿]|[🌾🍳🍼🎓🎤🎨🏫🏭💻💼🔧🔬🚀🚒🦯-🦳🦼🦽]))?|🏽(?:‍(?:[⚕⚖✈]️?|❤️?‍(?:💋‍)?👨[🏻-🏿]|🤝‍👨[🏻🏼🏾🏿]|[🌾🍳🍼🎓🎤🎨🏫🏭💻💼🔧🔬🚀🚒🦯-🦳🦼🦽]))?|🏾(?:‍(?:[⚕⚖✈]️?|❤️?‍(?:💋‍)?👨[🏻-🏿]|🤝‍👨[🏻-🏽🏿]|[🌾🍳🍼🎓🎤🎨🏫🏭💻💼🔧🔬🚀🚒🦯-🦳🦼🦽]))?|🏿(?:‍(?:[⚕⚖✈]️?|❤️?‍(?:💋‍)?👨[🏻-🏿]|🤝‍👨[🏻-🏾]|[🌾🍳🍼🎓🎤🎨🏫🏭💻💼🔧🔬🚀🚒🦯-🦳🦼🦽]))?)?|👩(?:‍(?:[⚕⚖✈]️?|❤️?‍(?:[👨👩]|💋‍[👨👩])|👦(?:‍👦)?|👧(?:‍[👦👧])?|👩‍(?:👦(?:‍👦)?|👧(?:‍[👦👧])?)|[🌾🍳🍼🎓🎤🎨🏫🏭💻💼🔧🔬🚀🚒🦯-🦳🦼🦽])|🏻(?:‍(?:[⚕⚖✈]️?|❤️?‍(?:💋‍)?[👨👩][🏻-🏿]|🤝‍[👨👩][🏼-🏿]|[🌾🍳🍼🎓🎤🎨🏫🏭💻💼🔧🔬🚀🚒🦯-🦳🦼🦽]))?|🏼(?:‍(?:[⚕⚖✈]️?|❤️?‍(?:💋‍)?[👨👩][🏻-🏿]|🤝‍[👨👩][🏻🏽-🏿]|[🌾🍳🍼🎓🎤🎨🏫🏭💻💼🔧🔬🚀🚒🦯-🦳🦼🦽]))?|🏽(?:‍(?:[⚕⚖✈]️?|❤️?‍(?:💋‍)?[👨👩][🏻-🏿]|🤝‍[👨👩][🏻🏼🏾🏿]|[🌾🍳🍼🎓🎤🎨🏫🏭💻💼🔧🔬🚀🚒🦯-🦳🦼🦽]))?|🏾(?:‍(?:[⚕⚖✈]️?|❤️?‍(?:💋‍)?[👨👩][🏻-🏿]|🤝‍[👨👩][🏻-🏽🏿]|[🌾🍳🍼🎓🎤🎨🏫🏭💻💼🔧🔬🚀🚒🦯-🦳🦼🦽]))?|🏿(?:‍(?:[⚕⚖✈]️?|❤️?‍(?:💋‍)?[👨👩][🏻-🏿]|🤝‍[👨👩][🏻-🏾]|[🌾🍳🍼🎓🎤🎨🏫🏭💻💼🔧🔬🚀🚒🦯-🦳🦼🦽]))?)?|[👯🤼🧞🧟](?:‍[♀♂]️?)?|😮(?:‍💨)?|😵(?:‍💫)?|😶(?:‍🌫️?)?|🧑(?:‍(?:[⚕⚖✈]️?|🤝‍🧑|[🌾🍳🍼🎄🎓🎤🎨🏫🏭💻💼🔧🔬🚀🚒🦯-🦳🦼🦽])|🏻(?:‍(?:[⚕⚖✈]️?|❤️?‍(?:💋‍)?🧑[🏼-🏿]|🤝‍🧑[🏻-🏿]|[🌾🍳🍼🎄🎓🎤🎨🏫🏭💻💼🔧🔬🚀🚒🦯-🦳🦼🦽]))?|🏼(?:‍(?:[⚕⚖✈]️?|❤️?‍(?:💋‍)?🧑[🏻🏽-🏿]|🤝‍🧑[🏻-🏿]|[🌾🍳🍼🎄🎓🎤🎨🏫🏭💻💼🔧🔬🚀🚒🦯-🦳🦼🦽]))?|🏽(?:‍(?:[⚕⚖✈]️?|❤️?‍(?:💋‍)?🧑[🏻🏼🏾🏿]|🤝‍🧑[🏻-🏿]|[🌾🍳🍼🎄🎓🎤🎨🏫🏭💻💼🔧🔬🚀🚒🦯-🦳🦼🦽]))?|🏾(?:‍(?:[⚕⚖✈]️?|❤️?‍(?:💋‍)?🧑[🏻-🏽🏿]|🤝‍🧑[🏻-🏿]|[🌾🍳🍼🎄🎓🎤🎨🏫🏭💻💼🔧🔬🚀🚒🦯-🦳🦼🦽]))?|🏿(?:‍(?:[⚕⚖✈]️?|❤️?‍(?:💋‍)?🧑[🏻-🏾]|🤝‍🧑[🏻-🏿]|[🌾🍳🍼🎄🎓🎤🎨🏫🏭💻💼🔧🔬🚀🚒🦯-🦳🦼🦽]))?)?|[⌚⌛⏩-⏬⏰⏳◽◾☔☕♈-♓♿⚓⚡⚪⚫⚽⚾⛄⛅⛎⛔⛪⛲⛳⛵⛺⛽✅✨❌❎❓-❕❗➕-➗➰➿⬛⬜⭐⭕🀄🃏🆎🆑-🆚🈁🈚🈯🈲-🈶🈸-🈺🉐🉑🌀-🌠🌭-🌵🌷-🍼🍾-🎄🎆-🎓🎠-🏁🏅🏆🏈🏉🏏-🏓🏠-🏰🏸-🐇🐉-🐔🐖-🐺🐼-🐾👀👄👅👑-👥👪👹-👻👽-💀💄💈-💎💐💒-💩💫-📼📿-🔽🕋-🕎🕐-🕧🖤🗻-😭😯-😴😷-🙄🙈-🙊🚀-🚢🚤-🚳🚷-🚿🛁-🛅🛐-🛒🛕-🛗🛝-🛟🛫🛬🛴-🛼🟠-🟫🟰🤍🤎🤐-🤗🤠-🤥🤧-🤯🤺🤿-🥅🥇-🥶🥸-🦴🦷🦺🦼-🧌🧐🧠-🧿🩰-🩴🩸-🩼🪀-🪆🪐-🪬🪰-🪺🫀-🫂🫐-🫙🫠-🫧]|🫱(?:🏻(?:‍🫲[🏼-🏿])?|🏼(?:‍🫲[🏻🏽-🏿])?|🏽(?:‍🫲[🏻🏼🏾🏿])?|🏾(?:‍🫲[🏻-🏽🏿])?|🏿(?:‍🫲[🏻-🏾])?)?"
	emojiReg := regexp.MustCompile(reg)
	contentParser.SetContent(emojiReg.ReplaceAllString(contentParser.Content(), ""))
	return contentParser
}
