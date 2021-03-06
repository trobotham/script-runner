// 2019 Daniel Oaks <daniel@danieloaks.net>
// released under the MIT license
package lib

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

const (
	htmlTemplate = `<html>
<head>
<title>%[1]s - IRC Foundation Test Framework</title>
<meta charset="utf-8">
<style>

:root {
	--sans-font-family: Trebuchet MS, Lucida Grande, Lucida Sans Unicode, Lucida Sans, Tahoma, sans-serif;
	--mono-font-family: monaco, Consolas, Lucida Console, monospace;

	--side-indent: 0.65em;
	--content-indent: 0.3em;
}

html {
	box-sizing: border-box;
}
*, *:before, *:after {
	box-sizing: inherit;
}

html {
	margin: 0;
	padding: 0;
	min-height: 100%%;
}
body {
	font-family: var(--sans-font-family);
	display: flex;
	flex-direction: column;
	min-height: 100%%;
	margin: 0;
	padding: 0;
}
header {
	flex: 0 0 content;
	padding: 0 var(--side-indent);
	display: flex;
	flex-direction: column;
}
.content {
	flex: 1 1 auto;
	background: #d4e3ed;
	min-height: 100%%;
	min-width: 100%%;
	display: flex;
	flex-direction: column;
}
footer {
	flex: 0 0 content;
	display: flex;
	padding: 0.3em var(--side-indent);
}

.tab-bar {
	display: flex;
	width: 100%%;
	padding: 0 1.5em 0 var(--content-indent);
	flex: 0 0 content;
	top: -1px;
	position: sticky;
	background: #fff;
}
.tabs {
	flex: 1 1 auto;
	display: flex;
}
.tab {
	padding: 0.3em 0.6em;
	border-top: 1px solid #ddf;
	border-right: 1px solid #ddf;
}
.tab:hover, .tab.active {
	background: #eef;
}
.tab:first-child {
	border-left: 1px solid #ddf;
}
.tab-options {
	flex: 0 0 content;
}
.tab-content {
	flex: 1 1 auto;
	padding: 0;
}
.tab-content .section:nth-child(2n) {
	background: #eee;
}
.options {
	display: flex;
}

.section {
	display: flex;
}
.lines {
	flex: 1 1 auto;
	font-size: 1.3em;
}
pre {
	margin: 0;
	white-space: pre-wrap;
}
pre.c {
	font-weight: bold;
}
pre.err {
	font-weight: bold;
	color: #822;
}
%[5]s

a {
	color: #217de4;
	font-style: italic;
	text-decoration: none;
}
h1 {
	font-size: 3em;
	color: #243847;
	padding: 0.5em 0;
	margin: 0;
}
.desc {
	display: block;
    color: #455e6e;
    margin-top: -0.8em;
    font-size: 1.05em;
    padding-bottom: 1.7em;
}

</style>
</head>
<body>

<header>
	<h1>%[1]s</h1>
	<span class="desc">%[2]s</span>
</header>

<div class="content">
	<div class="tab-bar">
		<div class="tabs">
%[3]s
		</div>
		<div class="options">
			<a id="sanitised-toggle" class="tab emoji" href="#" title="Toggle Sanitised/Raw">🎨</a>
		</div>
	</div>
	<div class="tab-content">
		<div class="section">
			<div class="lines active">
<pre>
Content here
</pre>
			</div>
			<div class="collapse-button">

			</div>
		</div>
	</div>
</div>

<footer>
	<a href="https://github.com/irccom/script-runner">IRC Foundation's script-runner</a>
</footer>

<script>

var serverInfo = %[4]s;

var server = serverInfo['default-server']
var sanitised = true

// print default server to console
console.log('default server is', serverInfo['default-server'])

// setup listeners for ircd buttons
var ircdButtons = document.querySelectorAll('[data-tab-button]')
for (var i = 0, len = ircdButtons.length; i < len; i++) {
	ircdButtons[i].addEventListener('click', (event) => {
		event.preventDefault()
		// console.log('btn ' + event.currentTarget.dataset.serverid)
		server = event.currentTarget.dataset.serverid
		showLogFor(server, sanitised)
	})
}

// setup listener for sanitised toggle
var sanitisedToggle = document.getElementById('sanitised-toggle')
sanitisedToggle.addEventListener('click', (event) => {
	event.preventDefault()
	sanitised = !sanitised
	showLogFor(server, sanitised)
	sanitisedToggle.classList.toggle('active')
})
sanitisedToggle.classList.add('active')

function showLogFor(ircd, sanitised) {
	console.log('showing log for', ircd)

	// 'press' button in the gui
	var ircdButtons = document.querySelectorAll('[data-tab-button]')
	for (var i = 0, len = ircdButtons.length; i < len; i++) {
		if (ircdButtons[i].dataset.serverid == ircd) {
			ircdButtons[i].classList.add('active')
		} else {
			ircdButtons[i].classList.remove('active')
		}
	}

	// construct and populate lines content
	var lines = document.createElement("div");
	lines.classList.add('lines')
	var logs
	if (sanitised) {
		logs = serverInfo["server-logs"][ircd]["sanitised"]
	} else {
		logs = serverInfo["server-logs"][ircd]["raw"]
	}
	for (var i = 0, len = logs.length; i < len; i++) {
		var raw = logs[i]
		// console.log(raw['c'], raw['s'], raw['l']);

		var content = ''
		if (0 < raw['e'].length) {
			content = raw['e']
		} else {
			content = raw['c'] + ' '
			if (raw['s'] == 'c') {
				content += ' ->'
			} else {
				content += '<- '
			}
			content += ' ' + raw['l']
		}

		var line = document.createElement("pre")
		line.classList.add('c-' + raw['c'])
		if (0 < raw['e'].length) {
			line.classList.add('err')
		} else {
			line.classList.add(raw['s'])
		}
		line.innerText = content

		lines.appendChild(line)
	}

	// replace it
	var linesElements = document.querySelectorAll('.lines.active')
	for (var i = 0, len = linesElements.length; i < len; i++) {
		var parent = linesElements[i].parentElement

		// out with the old
		linesElements[i].classList.remove('active')
		parent.removeChild(linesElements[i])

		// in with the new
		lines.classList.add('active')
		parent.appendChild(lines)
	}
}

showLogFor(server, sanitised);

</script>

</body>
</html>`

	tabText = `<a class="tab" href="" data-tab-button data-serverid="%[1]s">%[2]s</a>`

	cssPreTemplate = `pre.c-%[1]s {
	padding: 2px 0 2px calc(%[2]dch + var(--content-indent));
	text-indent: -%[2]dch;
	background: hsl(%[3]d, 90%%, 90%%);
}
`
)

type htmlJSONBlob struct {
	DefaultServer string                `json:"default-server"`
	ServerNames   map[string]string     `json:"server-names"`
	ServerLogs    map[string]serverBlob `json:"server-logs"`
}

type serverBlob struct {
	Raw       []lineBlob `json:"raw"`
	Sanitised []lineBlob `json:"sanitised"`
}

type lineBlob struct {
	Client string `json:"c"`
	SentBy string `json:"s"`
	Line   string `json:"l"`
	Error  string `json:"e"`
}

// HTMLFromResults takes a set of results and outputs an HTML representation of those results.
func HTMLFromResults(script *Script, serverConfigs map[string]ServerConfig, scriptResults map[string]*ScriptResults) string {
	// sorted ID list
	var sortedIDs []string
	for id := range serverConfigs {
		sortedIDs = append(sortedIDs, id)
	}
	sort.Strings(sortedIDs)

	// tab buttons
	var tabs string
	for _, id := range sortedIDs {
		tabs += fmt.Sprintf(tabText, id, serverConfigs[id].DisplayName)
	}

	// css
	var css string
	hue := 150
	for id := range script.Clients {
		// + 5 for ' <-  ' or similar
		css += fmt.Sprintf(cssPreTemplate, id, len(id)+5, hue)
		hue += 40
	}

	// construct JSON blob used by the page
	blob := htmlJSONBlob{
		DefaultServer: sortedIDs[0],
		ServerNames:   make(map[string]string),
		ServerLogs:    make(map[string]serverBlob),
	}
	for id, info := range serverConfigs {
		blob.ServerNames[id] = info.DisplayName
	}
	for id, sr := range scriptResults {
		var sBlob serverBlob

		var actionIndex int
		for _, srl := range sr.Lines {
			switch srl.Type {
			case ResultIRCMessage:
				// raw line
				lineRaw := lineBlob{
					Client: srl.Client,
					SentBy: "s",
					Line:   strings.TrimSuffix(srl.RawLine, "\r\n"),
				}
				sBlob.Raw = append(sBlob.Raw, lineRaw)

				// sanitised line
				sanitisedLine := srl.RawLine
				for orig, new := range serverConfigs[id].SanitisedReplacements {
					sanitisedLine = strings.Replace(sanitisedLine, orig, new, -1)
				}
				lineSanitised := lineBlob{
					Client: srl.Client,
					SentBy: "s",
					Line:   strings.TrimSuffix(sanitisedLine, "\r\n"),
				}
				sBlob.Sanitised = append(sBlob.Sanitised, lineSanitised)
			case ResultActionSync:
				thisAction := script.Actions[actionIndex]
				if thisAction.LineToSend != "" {
					// sent line always stays the same
					line := lineBlob{
						Client: thisAction.Client,
						SentBy: "c",
						Line:   strings.TrimSuffix(thisAction.LineToSend, "\r\n"),
					}
					sBlob.Raw = append(sBlob.Raw, line)
					sBlob.Sanitised = append(sBlob.Sanitised, line)
				}
				actionIndex++
			case ResultDisconnected:
				line := lineBlob{
					Client: srl.Client,
					Error:  fmt.Sprintf("%s was disconnected unexpectedly", srl.Client),
				}
				sBlob.Raw = append(sBlob.Raw, line)
				sBlob.Sanitised = append(sBlob.Sanitised, line)
			case ResultDisconnectedExpected:
				line := lineBlob{
					Client: srl.Client,
					Error:  fmt.Sprintf("%s was disconnected", srl.Client),
				}
				sBlob.Raw = append(sBlob.Raw, line)
				sBlob.Sanitised = append(sBlob.Sanitised, line)
			}
		}

		blob.ServerLogs[id] = sBlob
	}

	// marshall json blob
	blobBytes, err := json.Marshal(blob)
	blobString := "{'error': 1}"
	if err == nil {
		blobString = string(blobBytes)
	}

	// assemble template
	return fmt.Sprintf(htmlTemplate, script.Name, script.ShortDescription, tabs, blobString, css)
}
