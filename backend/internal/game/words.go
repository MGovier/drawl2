package game

import "math/rand"

var wordList = []string{
	// UK expressions & slang
	"dodgy kebab", "full English", "cheeky Nando's", "bog roll",
	"wobbly trolley", "jammie dodger", "ploughman's lunch",
	"beer garden", "chip butty", "mushy peas", "spotted dick",
	"sticky toffee pudding", "pie and mash", "Sunday roast", "cream tea",
	"lost the plot", "taking the mick", "having a kip", "throwing a wobbly",
	"off your trolley", "knees up", "bee's knees", "dog's dinner",
	"cat's pyjamas", "a right muppet", "gobsmacked pigeon", "chuffed to bits",

	// Cheeky double meanings (safe for image generation)
	"spotted woodpecker", "booby trap", "releasing the kraken",
	"love handles", "hot cross buns", "peacock", "a handful",
	"winner's trophy", "blowing a raspberry", "stuffed turkey",
	"plucking a chicken", "tossing a pancake", "great tits on a branch",
	"a cracking pair of eggs", "the family jewels", "a wobbly sausage roll",
	"a moist cake", "battered cod", "shaking the coconut tree", "banana split",

	// Funny situations
	"man vs seagull", "uncle dancing at a wedding", "stepping on Lego",
	"pigeons fighting over chips", "cat on a Roomba", "sunburnt tourist",
	"dad at a barbecue", "falling off a chair", "angry goose chase",
	"slipping on a banana", "dog stealing a sausage", "bad hair day",
	"nana on the dance floor", "crying at the pub", "kebab at 3am",
	"seagull stealing ice cream", "wasp at a picnic", "toilet out of paper",
	"missing the bus", "wrong Zoom call", "burnt toast", "flat tyre",

	// Pop culture & absurd
	"Bigfoot on holiday", "ghost driving a bus", "alien at Tesco",
	"vampire at the dentist", "the Loch Ness Monster", "wizard in Primark",
	"dinosaur at a rave", "mermaid with legs", "unicorn traffic jam",
	"robot doing yoga", "Elvis at the chippy", "ninja librarian",
	"pirate accountant", "zombie wedding", "cowboy plumber",

	// Animals doing things
	"posh cat", "angry swan", "hungover owl", "suspicious hedgehog",
	"flamingo in wellies", "bear in a hot tub", "penguin commute",
	"squirrel heist", "otter holding hands", "fox in a bin",
	"duck with attitude", "moose in a lift", "tortoise racing",
	"hamster on a wheel", "llama drama",

	// British culture
	"queue jumper", "motorway services", "rainy barbecue",
	"wonky shopping trolley", "warm pub pint", "dodgy tan lines",
	"caravan holiday", "garden gnome", "passive-aggressive note",
	"soggy camping trip", "roundabout rage", "charity shop bargain",
	"school disco", "ice cream van", "village fete",
	"car boot sale", "double decker bus", "phone box",
	"red post box", "fish and chips",

	// Easy-to-draw classics (keeping some for balance)
	"volcano", "rocket", "treehouse", "haunted house",
	"desert island", "snowman", "treasure map", "roller coaster",
	"hot air balloon", "pirate ship", "dragon", "castle",
	"UFO", "campfire", "lighthouse", "rainbow",
	"palm tree", "igloo", "submarine", "catapult",
}

func RandomWords(n int) []string {
	perm := rand.Perm(len(wordList))
	words := make([]string, n)
	for i := 0; i < n; i++ {
		words[i] = wordList[perm[i%len(perm)]]
	}
	return words
}
