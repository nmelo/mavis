package drill

import "github.com/nmelo/mavis/internal/level"

// StoryLines returns the curated story lines for a given level's word drills.
// The story is inspired by True Names by Vernor Vinge: hackers, the Other Plane,
// digital identities, and the tension between the net and the real world.
// Each line uses ONLY characters unlocked at that level (plus space).
var storyLines = map[int][]string{
	// Level 3: f j d k s l + space
	// Pure drill patterns — no real English words exist with just these letters.
	3: {
		"ff jj dd kk ss ll",
		"fjdk slkj fdsk ljfd",
		"jkl fds jkl fds",
		"skl dfs lkj fds",
		"fds jkl sdk fld",
		"slf dkf jsl fdk",
		"dfl jsk fld ksl",
		"jfd ksl fjd ksl",
		"lkj fds lkj fds",
		"sdk flj dks ljf",
	},

	// Level 4: + a ;
	// Words: ask, all, fall, flask, lads, lass, salad, sad, fad, dad, alas
	4: {
		"alas a sad flask",
		"all lads ask dad",
		"a lass adds a salad",
		"ask all; ask dad",
		"a lad falls; alas all falls",
		"dad asks all lads",
		"a sad fad; all flask salad",
		"alas dad; all lads fall",
		"flask salad; ask all",
		"a sad lass asks a lad",
	},

	// Level 5: + g h
	// Words: flash, glass, gash, hash, dash, glad, hall, shall, half, flag, gala
	5: {
		"flash ash gash hall",
		"a glad lad had a dash",
		"shall all flags flash",
		"half a glass has a hash",
		"glass slag ash flash gala",
		"a gash had flash; shall dash",
		"glad gals dash; flags flash",
		"shag hall gall lash slag",
		"hash flash half glass flag",
		"a glad hall shall flash all",
	},

	// Level 6: + t y
	// Words: that, stay, hasty, salty, ghastly, flatly, lastly, flashy, stag, shaft
	6: {
		"that flat gala stays shy",
		"flashy stag daft lastly",
		"stay at that tall shaft",
		"salty data flash ghastly",
		"a hasty stag halts flatly",
		"that last flag sadly falls",
		"flashy at last; tall flags stay",
		"salty flaky data; that sadly",
		"stag halts lastly; stay fast",
		"a ghastly shaft; that tally",
	},

	// Level 7: + r u
	// Words: dust, rust, trust, guard, dark, dusk, rush, hurt, skull, stray, lurk
	7: {
		"dusty guards rush at dusk",
		"trust that rusty dark skull",
		"a dark stray lurks far",
		"just shut dusty dusty halls",
		"rust hurts darkly at dusk",
		"dusty guards trust a sultry fury",
		"dark rust drags at trust",
		"hurry fast; guard that dusty gust",
		"dusty shards guard a harsh shaft",
		"a dusk rush starts ultra darkly",
	},

	// Level 8: + e i
	// Words: the, fire, shield, strike, digital, desire, field, elite, light, riddle
	8: {
		"the digital shield is hid",
		"fire strikes their dark field",
		"desire stirs the skilled elite",
		"the elite guard hid the data",
		"she figured the riddle first",
		"their real title stayed hid",
		"she strikes the dark keys right",
		"the shield sighs at firelight",
		"seek the digital field; delete fear",
		"a little desire stirs the fight",
	},

	// Level 9: + w o
	// Words: world, shadow, follow, other, tower, sword, flow, throw, old, wood
	9: {
		"who follows the other world",
		"shadows flow to the old tower",
		"the hollow sword shows its worth",
		"words flow softly toward the light",
		"she showed her the world outside",
		"slow shadows follow the door",
		"the world grows weirder to follow",
		"throw the old sword to the tower",
		"two words would show their worth",
		"wood doors glow throughout",
	},

	// Level 10: + q p
	// Words: speak, quest, equip, plot, query, portal, power, password, opaque
	10: {
		"speak the password to the portal",
		"quest through opaque pathways",
		"equip your thoughts with power",
		"the plot quietly spreads its roots",
		"speak of the split quest at top",
		"who plots through quiet steps",
		"equipped people stop their queries",
		"speak quietly through the portal",
		"a queer path spoke with power",
		"queued requests pulse with purpose",
	},

	// Level 11: + v m
	// Words: move, virtual, malware, vast, volume, morph, vivid, solve, doom, storm
	11: {
		"move through the virtual realm",
		"five vast volumes of malware",
		"the dim map seemed almost alive",
		"virtual avatars shimmer quietly",
		"solve the massive firewall doom",
		"malware moved past every limit",
		"the virtual map reforms its shape",
		"vivid memories from the vast depths",
		"some rapid morph swept the void",
		"a massive virtual storm looms with dread",
	},

	// Level 12: + c ,
	// Words: code, circle, cosmic, crack, cypher, crypt, cache, cloak, cult, cascade
	12: {
		"the cult circles, code is core",
		"crack the cypher, move with care",
		"cosmic data cascades from the rift",
		"each circle cloaks a secret voice",
		"code compiles, cryptic logic",
		"dark cache clears, the cult commits",
		"cascades fracture, cosmic cracks appear",
		"scope each circuit, complete the code",
		"the cult claims victory, a crucial choice",
		"compact circles of cosmic crypt data",
	},

	// Level 13: + x .
	// Words: matrix, complex, exploit, proxy, hex, vortex, text, exit, fix, extract
	13: {
		"the matrix exposes a complex code.",
		"complex exploits ripple outward.",
		"exit the proxy. fix the vortex.",
		"hex oxide exists as complex text.",
		"the apex exposed. exploit it exactly.",
		"mix the complex. extract the fix.",
		"exact pixel depths exceed limits.",
		"explore the matrix. expose the text.",
		"six exploits. execute with thought.",
		"the vortex expels complex flux.",
	},

	// Level 14: + z /
	// Words: wizard, freeze, maze, haze, zigzag, fizz, fuzz, zero, daze, quiz, seize
	14: {
		"the wizard gazed at a froze maze",
		"fuzzy data sleeps ooze of haze",
		"a zigzag riff sparks the freeze",
		"the wizard says seize the prize",
		"zigzag through the froze depths",
		"freeze the ooze to daze a foe",
		"lazy fizz drifts to the zero wire",
		"she seizes the quiz amidst haze",
		"the froze wizard gazes at zero",
		"fuzzy daze grips the dizzied quiz",
	},
}

// GetStoryLine returns a story line for the given level and drill number (0-indexed).
// Returns empty string if no story exists for that level.
func GetStoryLine(levelNum, drillIndex int) string {
	lines, ok := storyLines[levelNum]
	if !ok || len(lines) == 0 {
		return ""
	}
	return lines[drillIndex%len(lines)]
}

// HasStory returns true if curated story content exists for the given level.
func HasStory(levelNum int) bool {
	lines, ok := storyLines[levelNum]
	return ok && len(lines) > 0
}

// ValidateStoryLines checks that all story lines only use characters available
// at their level. Returns a map of level -> list of invalid lines.
func ValidateStoryLines() map[int][]string {
	invalid := map[int][]string{}
	for lvl, lines := range storyLines {
		allowed := make(map[rune]bool)
		for _, k := range level.UnlockedKeys(lvl) {
			allowed[k] = true
		}
		allowed[' '] = true // space always allowed in word drills

		for _, line := range lines {
			for _, ch := range line {
				if !allowed[ch] {
					invalid[lvl] = append(invalid[lvl], line)
					break
				}
			}
		}
	}
	return invalid
}
