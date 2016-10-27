package badge

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"unicode/utf8"
)

const badgeTemplate = `
<svg xmlns="http://www.w3.org/2000/svg" width="{{ .TotalWidth }}" height="18">
  <linearGradient id="smooth" x2="0" y2="100%">
    <stop offset="0"  stop-color="#fff" stop-opacity=".7"/>
    <stop offset=".1" stop-color="#aaa" stop-opacity=".1"/>
    <stop offset=".9" stop-color="#000" stop-opacity=".3"/>
    <stop offset="1"  stop-color="#000" stop-opacity=".5"/>
  </linearGradient>

  <mask id="round">
    <rect width="{{ .TotalWidth }}" height="18" rx="4" fill="#fff"/>
  </mask>

  <g mask="url(#round)">
    <rect width="{{ .LeftWidth }}" height="18" fill="#555"/>
    <rect x="{{ .LeftWidth }}" width="{{ .RightWidth }}" height="18" fill="{{ .Color }}"/>
    <rect width="{{ .TotalWidth }}" height="18" fill="url(#smooth)"/>
  </g>

  <g fill="#fff" text-anchor="middle" font-family="DejaVu Sans,Verdana,Geneva,sans-serif" font-size="12">
    <text x="{{ .LeftOffset }}" y="14" fill="#010101" fill-opacity=".3">{{ .Left }}</text>
    <text x="{{ .LeftOffset }}" y="13">{{ .Left }}</text>
    <text x="{{ .RightOffset }}" y="14" fill="#010101" fill-opacity=".3">{{ .Right }}</text>
    <text x="{{ .RightOffset }}" y="13">{{ .Right }}</text>
  </g>
</svg>
`
const downStream = `
<svg xmlns="http://www.w3.org/2000/svg" width="152.1" height="18">
  <linearGradient id="smooth" x2="0" y2="100%">
    <stop offset="0"  stop-color="#fff" stop-opacity=".7"/>
    <stop offset=".1" stop-color="#aaa" stop-opacity=".1"/>
    <stop offset=".9" stop-color="#000" stop-opacity=".3"/>
    <stop offset="1"  stop-color="#000" stop-opacity=".5"/>
  </linearGradient>

  <mask id="round">
    <rect width="152.1" height="18" rx="4" fill="#fff"/>
  </mask>

  <g mask="url(#round)">
    <rect width="90" height="18" fill="#555"/>
    <rect x="90" width="67.5" height="18" fill="red"/>
    <rect width="152.1" height="18" fill="url(#smooth)"/>
  </g>

  <g fill="#fff" text-anchor="middle" font-family="DejaVu Sans,Verdana,Geneva,sans-serif" font-size="12">
    <text x="46" y="14" fill="#010101" fill-opacity=".3">Downstream</text>
    <text x="46" y="13">Downstream</text>
    <text x="122.75" y="14" fill="#010101" fill-opacity=".3">Error</text>
    <text x="122.75" y="13">Error</text>
  </g>
</svg>
`

var tmplBadge = template.Must(template.New("").Parse(badgeTemplate))

type badgeParams struct {
	Left        string
	LeftWidth   float64
	LeftOffset  float64
	Right       string
	RightWidth  float64
	RightOffset float64
	TotalWidth  float64
	Color       string
}

func createBadgeParams(l, r, c string) badgeParams {
	b := badgeParams{Left: l, Right: r, Color: c}
	b.LeftWidth = float64(utf8.RuneCountInString(b.Left)) * 9
	b.LeftOffset = b.LeftWidth/2 + 1
	b.RightWidth = float64(utf8.RuneCountInString(b.Right)) * 13.5
	b.RightOffset = b.LeftWidth + b.RightWidth/2 - 1
	b.TotalWidth = b.LeftWidth + b.RightWidth*0.92
	return b
}

func CreateBadge(l, r, c string) (error, []byte) {
	bp := createBadgeParams(l, r, c)
	var buf bytes.Buffer
	if err := tmplBadge.Execute(&buf, bp); err != nil {
		return err, nil
	}
	return nil, buf.Bytes()
}

func DownstreamError(w http.ResponseWriter, err error) {
	fmt.Println(err)
	w.Header().Add("Content-Type", "image/svg+xml")
	w.Write([]byte(downStream))
}
