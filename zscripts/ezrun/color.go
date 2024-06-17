package main

import (
	"fmt"
	"slices"

	"github.com/lucasb-eyer/go-colorful"

	"ezpkg.io/bytez"
)

// https://github.com/lucasb-eyer/go-colorful/blob/master/doc/gradientgen/gradientgen.go

var C50_900 = []int{50, 100, 200, 300, 400, 500, 600, 700, 800, 900}

type Color struct {
	colorful.Color
}

func (c Color) String() string {
	return c.Hex()
}

func (c Color) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%q", c.Hex())), nil
}

func (c Color) Darken(p float64) Color {
	l, a, b := c.Lab()
	l *= (1 - p)
	return Color{Color: colorful.Lab(l, a, b)}
}

type MapColors []struct {
	Code  int
	Color Color
}

func (cs MapColors) MarshalJSON() ([]byte, error) {
	var b bytez.Buffer
	b.WriteByteZ('{')
	for i, c := range cs {
		if i > 0 {
			b.WriteByteZ(',')
		}
		b.Printf(`"%d":%q`, c.Code, c.Color.Hex())
	}
	b.WriteByteZ('}')
	return b.Bytes(), nil
}

func GradientTable_0() *GradientTable {
	return NewGradientTable(0).
		Add(0.0, MustParseColors("#9E0142")).
		Add(0.0, MustParseColors("#9E0142")).
		Add(0.1, MustParseColors("#D53E4F")).
		Add(0.2, MustParseColors("#F46D43")).
		// Add( {0.3, MustParseColors("#FDAE61")).
		// Add( {0.4, MustParseColors("#FEE090")).
		Add(0.3, MustParseColors("#D5A600")).
		// Add( {0.6, MustParseColors("#E6F598")).
		Add(0.5, MustParseColors("#65C97A")).
		// Add( {0.8, MustParseColors("#66C2A5")).
		Add(0.75, MustParseColors("#3288BD")).
		Add(0.85, MustParseColors("#5E4FA2")).
		Add(1.0, MustParseColors("#925CB1"))
}

func GradientTable_1() *GradientTable {
	return NewGradientTable().
		Add(0.0, MustParseColors("#9e0142")).
		Add(0.1, MustParseColors("#d53e4f")).
		Add(0.2, MustParseColors("#f46d43")).
		Add(0.3, MustParseColors("#fdae61")).
		Add(0.4, MustParseColors("#fee090")).
		Add(0.5, MustParseColors("#ffffbf")).
		Add(0.6, MustParseColors("#e6f598")).
		Add(0.7, MustParseColors("#abdda4")).
		Add(0.8, MustParseColors("#66c2a5")).
		Add(0.9, MustParseColors("#3288bd")).
		Add(1.0, MustParseColors("#5e4fa2"))
}

func GradientTable_Tailwind() *GradientTable {
	return NewGradientTable(C50_900...).
		/* Slate */
		// AddX("#f8fafc", "#f1f5f9", "#e2e8f0", "#cbd5e1", "#94a3b8", "#64748b", "#475569", "#334155", "#1e293b", "#0f172a").
		/* Gray */
		// AddX("#f9fafb", "#f3f4f6", "#e5e7eb", "#d1d5db", "#9ca3af", "#6b7280", "#4b5563", "#374151", "#1f2937", "#111827").
		/* Zinc */
		// AddX("#fafafa", "#f4f4f5", "#e4e4e7", "#d4d4d8", "#a1a1aa", "#71717a", "#52525b", "#3f3f46", "#27272a", "#18181b").
		/* Neutral */
		// AddX("#fafafa", "#f5f5f5", "#e5e5e5", "#d4d4d4", "#a3a3a3", "#737373", "#525252", "#404040", "#262626", "#171717").
		/* Stone */
		// AddX("#fafaf9", "#f5f5f4", "#e7e5e4", "#d6d3d1", "#a8a29e", "#78716c", "#57534e", "#44403c", "#292524", "#1c1917").
		/* Red */
		AddX("#fef2f2", "#fee2e2", "#fecaca", "#fca5a5", "#f87171", "#ef4444", "#dc2626", "#b91c1c", "#991b1b", "#7f1d1d").
		/* Orange */
		AddX("#fff7ed", "#ffedd5", "#fed7aa", "#fdba74", "#fb923c", "#f97316", "#ea580c", "#c2410c", "#9a3412", "#7c2d12").
		/* Amber */
		AddX("#fffbeb", "#fef3c7", "#fde68a", "#fcd34d", "#fbbf24", "#f59e0b", "#d97706", "#b45309", "#92400e", "#78350f").
		/* Yellow */
		AddX("#fefce8", "#fef9c3", "#fef08a", "#fde047", "#facc15", "#eab308", "#ca8a04", "#a16207", "#854d0e", "#713f12").
		/* Lime */
		AddX("#f7fee7", "#ecfccb", "#d9f99d", "#bef264", "#a3e635", "#84cc16", "#65a30d", "#4d7c0f", "#3f6212", "#365314").
		/* Green */
		AddX("#f0fdf4", "#dcfce7", "#bbf7d0", "#86efac", "#4ade80", "#22c55e", "#16a34a", "#15803d", "#166534", "#14532d").
		/* Emerald */
		AddX("#ecfdf5", "#d1fae5", "#a7f3d0", "#6ee7b7", "#34d399", "#10b981", "#059669", "#047857", "#065f46", "#064e3b").
		/* Teal */
		AddX("#f0fdfa", "#ccfbf1", "#99f6e4", "#5eead4", "#2dd4bf", "#14b8a6", "#0d9488", "#0f766e", "#115e59", "#134e4a").
		/* Cyan */
		AddX("#ecfeff", "#cffafe", "#a5f3fc", "#67e8f9", "#22d3ee", "#06b6d4", "#0891b2", "#0e7490", "#155e75", "#164e63").
		/* Sky */
		AddX("#f0f9ff", "#e0f2fe", "#bae6fd", "#7dd3fc", "#38bdf8", "#0ea5e9", "#0284c7", "#0369a1", "#075985", "#0c4a6e").
		/* Blue */
		AddX("#eff6ff", "#dbeafe", "#bfdbfe", "#93c5fd", "#60a5fa", "#3b82f6", "#2563eb", "#1d4ed8", "#1e40af", "#1e3a8a").
		/* Indigo */
		AddX("#eef2ff", "#e0e7ff", "#c7d2fe", "#a5b4fc", "#818cf8", "#6366f1", "#4f46e5", "#4338ca", "#3730a3", "#312e81").
		/* Violet */
		AddX("#f5f3ff", "#ede9fe", "#ddd6fe", "#c4b5fd", "#a78bfa", "#8b5cf6", "#7c3aed", "#6d28d9", "#5b21b6", "#4c1d95").
		/* Purple */
		AddX("#faf5ff", "#f3e8ff", "#e9d5ff", "#d8b4fe", "#c084fc", "#a855f7", "#9333ea", "#7e22ce", "#6b21a8", "#581c87").
		/* Fuchsia */
		AddX("#fdf4ff", "#fae8ff", "#f5d0fe", "#f0abfc", "#e879f9", "#d946ef", "#c026d3", "#a21caf", "#86198f", "#701a75").
		/* Pink */
		AddX("#fdf2f8", "#fce7f3", "#fbcfe8", "#f9a8d4", "#f472b6", "#ec4899", "#db2777", "#be185d", "#9d174d", "#831843").
		/* Rose */
		AddX("#fff1f2", "#ffe4e6", "#fecdd3", "#fda4af", "#fb7185", "#f43f5e", "#e11d48", "#be123c", "#9f1239", "#881337").
		RecalculatePositions()
}

// This table contains the "keypoints" of the colorgradient you want to generate.
// The position of each keypoint has to live in the range [0,1]
type GradientTable struct {
	Codes []int
	Items []GradientTableItem
}

type GradientTableItem struct {
	Pos    float64
	Colors []Color
}

func NewGradientTable(codes ...int) *GradientTable {
	return &GradientTable{Codes: codes}
}

func (gt *GradientTable) Add(pos float64, colors []Color) *GradientTable {
	if len(colors) != len(gt.Codes) {
		panic(fmt.Sprintf("wrong number of colors (expected %v)", len(gt.Codes)))
	}
	gt.Items = append(gt.Items, GradientTableItem{Pos: pos, Colors: colors})
	return gt
}

func (gt *GradientTable) AddX(colors ...string) *GradientTable {
	return gt.Add(0, MustParseColors(colors...))
}

func (gt *GradientTable) RecalculatePositions() *GradientTable {
	for i := range gt.Items {
		gt.Items[i].Pos = float64(i) / float64(len(gt.Items)-1)
	}
	return gt
}

// This is the meat of the gradient computation. It returns a HCL-blend between
// the two colors around `t`.
// Note: It relies heavily on the fact that the gradient keypoints are sorted.
func (gt GradientTable) GetInterpolatedColorFor(code int, t float64) Color {
	idx := slices.Index(gt.Codes, code)
	if idx < 0 {
		panic(fmt.Sprintf("invalid code: %v", code))
	}
	N := len(gt.Items)
	for i := 0; i < N-1; i++ {
		p1, p2 := gt.Items[i], gt.Items[i+1]
		if p1.Pos <= t && t <= p2.Pos {
			// We are in between p1 and p2. Go blend them!
			t := (t - p1.Pos) / (p2.Pos - p1.Pos)
			c1, c2 := p1.Colors[idx], p2.Colors[idx]
			return Color{Color: c1.BlendHcl(c2.Color, t).Clamped()}
		}
	}

	// Nothing found? Means we're at (or past) the last gradient keypoint.
	return gt.Items[N-1].Colors[idx]
}

func (gt GradientTable) GetInterpolatedPaletteFor(t float64) []Color {
	colors := make([]Color, len(gt.Codes))
	for i, code := range gt.Codes {
		colors[i] = gt.GetInterpolatedColorFor(code, t)
	}
	return colors
}

func MustParseHex(s string) Color {
	c, err := colorful.Hex(s)
	if err != nil {
		panic("MustParseHex: " + err.Error())
	}
	return Color{c}
}

func MustParseColors(ss ...string) []Color {
	colors := make([]Color, len(ss))
	for i, s := range ss {
		colors[i] = MustParseHex(s)
	}
	return colors
}
