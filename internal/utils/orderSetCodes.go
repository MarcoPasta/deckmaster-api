package utils

import (
	"fmt"
	"strings"
)

var setCodes = map[string]int{
	"base":            1,
	"jungle":          2,
	"fossil":          3,
	"base2":           4,
	"teamrocket":      5, // Team Rocket set
	"gymheroes":       6,
	"gymchallenge":    7,
	"neogenesis":      8,
	"neodiscovery":    9,
	"southernislands": 10,
	"neorevelation":   11,
	"neodestiny":      12,
	"expedition":      13,
	"aquapolis":       14,
	"skyridge":        15,
	"ex1":             16,  // EX Ruby & Sapphire
	"ex2":             17,  // EX Sandstorm
	"ex3":             18,  // EX Dragon
	"ex4":             19,  // EX Team Magma vs Team Aqua
	"ex5":             20,  // EX Hidden Legends
	"ex6":             21,  // EX FireRed & LeafGreen
	"ex7":             22,  // EX Team Rocket Returns
	"ex8":             23,  // EX Deoxys
	"ex9":             24,  // EX Emerald
	"ex10":            25,  // EX Unseen Forces
	"ex11":            26,  // EX Unseen Forces Unown Collection
	"ex12":            27,  // EX Delta Species
	"ex13":            28,  // EX Legend Maker
	"ex14":            29,  // EX Holon Phantoms
	"ex15":            30,  // EX Crystal Guardians
	"ex16":            31,  // EX Dragon Frontiers
	"ex17":            32,  // EX Power Keepers
	"dpp":             33,  // Diamond & Pearl Black Star Promos
	"dp1":             34,  // Diamond & Pearl
	"dp2":             35,  // Mysterious Treasures
	"dp3":             36,  // Secret Wonders
	"dp4":             37,  // Great Encounters
	"dp5":             38,  // Majestic Dawn
	"dp6":             39,  // Legends Awakened
	"dp7":             40,  // Stormfront
	"plat1":           41,  // Platinum
	"plat2":           42,  // Rising Rivals
	"plat3":           43,  // Supreme Victors
	"plat4":           44,  // Arceus
	"plat5":           45,  // Pokemon Rumble
	"hgss1":           46,  // HeartGold & SoulSilver
	"hgss2":           47,  // HeartGold & SoulSilver Promos
	"hgss3":           48,  // HeartGold & SoulSilver Unleashed
	"hgss4":           49,  // HeartGold & SoulSilver Undaunted
	"hgss5":           50,  // HeartGold & SoulSilver Triumphant
	"bwp":             51,  // Black & White Promos
	"bw1":             52,  // Black & White
	"bw2":             53,  // Emerging Powers
	"bw3":             54,  // Noble Victories
	"bw4":             55,  // Next Destinies
	"bw5":             56,  // Dark Explorers
	"bw6":             57,  // Dragons Exalted
	"bw7":             58,  // Dragon Vault <- Promo Set
	"bw8":             59,  // Boundaries Crossed
	"bw9":             60,  // Plasma Storm
	"bw10":            61,  // Plasma Freeze
	"bw11":            62,  // Plasma Blast
	"bw12":            63,  // Legendary Treasures
	"bw13":            64,  // Radian Collection <- Promo Set
	"xy1":             65,  // X & Y
	"xy2":             66,  // Promos
	"xy3":             67,  // Kalos Starter Set
	"xy4":             68,  // Flashfire
	"xy5":             69,  // Furious Fists
	"xy6":             70,  // Phantom Forces
	"xy7":             71,  // Primal Clash
	"xy8":             72,  // Double Crisis
	"xy9":             73,  // Roaring Skies
	"xy10":            74,  // Ancient Origins
	"xy11":            75,  // BREAKthrough
	"xy12":            76,  // BREAKpoint
	"xy13":            77,  // Generations
	"xy14":            78,  // Fates Collide
	"xy15":            79,  // Steam Siege
	"xy16":            80,  // Evolutions
	"sm1":             81,  // Sun & Moon
	"sm2":             82,  // Promos
	"sm3":             83,  // Guardians Rising
	"sm4":             84,  // Burning Shadows
	"sm5":             85,  // Shining Legends
	"sm6":             86,  // Crimson Invasion
	"sm7":             87,  // Ultra Prism
	"sm8":             88,  // Forbidden Light
	"sm9":             89,  // Celestial Storm
	"sm10":            90,  // Dragon Majesty
	"sm11":            91,  // Lost Thunder
	"sm12":            92,  // Team Up
	"sm13":            93,  // Detective Pikachu
	"sm14":            94,  // Unbroken Bonds
	"sm15":            95,  // Unified Minds
	"sm16":            96,  // Hidden Fates
	"sm17":            97,  // Cosmic Eclipse
	"swsh1":           98,  // Sword & Shield
	"swsh2":           99,  // Rebel Clash
	"swsh3":           100, // Darkness Ablaze
	"swsh4":           101, // Champion's Path
	"swsh5":           102, // Vivid Voltage
	"swsh6":           103, // Shining Fates
	"swsh7":           104, // Battle Styles
	"swsh8":           105, // Chilling Reign
	"swsh9":           106, // Evolving Skies
	"swsh10":          107, // Celebrations
	"swsh11":          108, // Fusion Strike
	"swsh12":          109, // Brilliant Stars
	"swsh13":          110, // Astral Radiance
	"swsh14":          111, // Pokemon Go
	"swsh15":          112, // Lost Origin
	"swsh16":          113, // Crown Zenith
	"svp":             114, // Scarlet & Violet Promos
	"sv1":             115, // Scarlet & Violet
	"sve":             116, // Scarlet & Violet Energies
	"pal":             117, // Paldea Evolved
	"obf":             118, // Obsidian Flames
	"151":             119, // 151 Collection
	"par":             120,
	"paf":             121,
	"tef":             122,
	"twm":             123,
	"sfa":             124,
	"scr":             125,
	"ssp":             126,
	"pre":             127,
	"jtg":             128,
	"dri":             129,
	"blk":             130,
	"wht":             131,
	"meg":             132,
	"mep":             133,
	"pfl":             134,
	"asc":             135,
	"por":             136,
	"cri":             137,
	"pbl":             138,
}

func PokemonSetOrderCase() string {
	cases := make([]string, 0, len(setCodes))
	for name, order := range setCodes {
		cases = append(cases, fmt.Sprintf("WHEN '%s' THEN %d", name, order))
	}
	return "CASE code " + strings.Join(cases, " ") + " ELSE 999 END"
}
