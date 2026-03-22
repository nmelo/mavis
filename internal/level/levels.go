package level

type Finger string

const (
	LPinky Finger = "L.pinky"
	LRing  Finger = "L.ring"
	LMid   Finger = "L.mid"
	LIndex Finger = "L.index"
	RIndex Finger = "R.index"
	RMid   Finger = "R.mid"
	RRing  Finger = "R.ring"
	RPinky Finger = "R.pinky"
	LThumb Finger = "L.thumb"
	RThumb Finger = "R.thumb"
)

type Level struct {
	Number        int
	Name          string
	NewKeys       []rune
	HasWordDrills bool
	HasCodeDrills bool
}

var fingerMap = map[rune]Finger{
	'a': LPinky, 's': LRing, 'd': LMid, 'f': LIndex,
	'j': RIndex, 'k': RMid, 'l': RRing, ';': RPinky,
	'g': LIndex, 'h': RIndex,
	'q': LPinky, 'w': LRing, 'e': LMid, 'r': LIndex, 't': LIndex,
	'y': RIndex, 'u': RIndex, 'i': RMid, 'o': RRing, 'p': RPinky,
	'z': LPinky, 'x': LRing, 'c': LMid, 'v': LIndex, 'b': LIndex,
	'n': RIndex, 'm': RIndex, ',': RMid, '.': RRing, '/': RPinky,
	' ': RThumb,
	'1': LPinky, '2': LRing, '3': LMid, '4': LIndex, '5': LIndex,
	'6': RIndex, '7': RIndex, '8': RMid, '9': RRing, '0': RPinky,
	'-': RPinky, '=': RPinky, '[': RPinky, ']': RPinky,
	'\\': RPinky, '\'': RPinky,
}

var levels = []Level{
	{1, "Home Row: f j", []rune{'f', 'j'}, false, false},
	{2, "Home Row: d k", []rune{'d', 'k'}, false, false},
	{3, "Home Row: s l", []rune{'s', 'l'}, true, false},
	{4, "Home Row: a ;", []rune{'a', ';'}, true, false},
	{5, "Home Row: g h", []rune{'g', 'h'}, true, false},
	{6, "Top Row: t y", []rune{'t', 'y'}, true, false},
	{7, "Top Row: r u", []rune{'r', 'u'}, true, false},
	{8, "Top Row: e i", []rune{'e', 'i'}, true, false},
	{9, "Top Row: w o", []rune{'w', 'o'}, true, false},
	{10, "Top Row: q p", []rune{'q', 'p'}, true, true},
	{11, "Bottom Row: v m", []rune{'v', 'm'}, true, true},
	{12, "Bottom Row: c ,", []rune{'c', ','}, true, true},
	{13, "Bottom Row: x .", []rune{'x', '.'}, true, true},
	{14, "Bottom Row: z /", []rune{'z', '/'}, true, true},
	{15, "Space", []rune{' '}, true, true},
	{16, "Shift + Keys", []rune{}, true, true},
	{17, "Number Row", []rune{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}, true, true},
	{18, "Symbols", []rune{'-', '=', '[', ']', '\\', '\''}, true, true},
}

func All() []Level {
	return levels
}

func Get(n int) Level {
	return levels[n-1]
}

func UnlockedKeys(levelNum int) []rune {
	var keys []rune
	for i := 0; i < levelNum && i < len(levels); i++ {
		keys = append(keys, levels[i].NewKeys...)
	}
	return keys
}

func FingerForKey(k rune) Finger {
	return fingerMap[k]
}
