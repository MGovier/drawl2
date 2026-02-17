package game

import "math/rand"

var wordList = []string{
	// Food & drink
	"fish and chips", "cup of tea", "Sunday roast", "full English",
	"cream tea", "sausage roll", "pork pie", "Cornish pasty",
	"bacon sandwich", "beans on toast", "crumpet", "scone",
	"trifle", "mince pie", "Yorkshire pudding", "shepherd's pie",
	"bangers and mash", "ice cream van", "fish finger sandwich",
	"jam doughnut", "birthday cake", "cheese toastie",

	// British landmarks & places
	"Big Ben", "London Eye", "Tower Bridge", "Buckingham Palace",
	"Stonehenge", "red phone box", "double decker bus", "black cab",
	"the Shard", "Brighton pier", "castle", "lighthouse",
	"church spire", "village green", "cricket pitch", "canal boat",
	"post box", "pub", "chippy", "corner shop",

	// Weather & outdoors
	"rainbow", "rainy barbecue", "umbrella", "puddle jumping",
	"snowman", "muddy wellies", "thunderstorm", "deckchair",
	"rockpool", "sandcastle", "bonfire", "conker",
	"kite flying", "picnic blanket", "camping tent", "caravan",

	// Animals
	"angry swan", "fox", "hedgehog", "robin",
	"corgi", "sheepdog", "Highland cow", "badger",
	"seagull stealing chips", "squirrel", "puffin", "otter",
	"cat in a box", "horse jumping", "duck pond", "spider in the bath",

	// Everyday life
	"queue", "shopping trolley", "school uniform", "lollipop lady",
	"garden gnome", "wheelie bin", "washing line", "doorbell",
	"roundabout", "zebra crossing", "parking meter", "speed camera",
	"flat tyre", "broken umbrella", "lost keys", "missed bus",
	"alarm clock", "packed lunch", "ironing", "hoovering",

	// Sports & games
	"football", "cricket bat", "tennis racket", "rugby tackle",
	"darts", "snooker", "bowling", "swimming pool",
	"golf", "skateboard", "fishing rod", "bicycle",
	"rowing boat", "trophy", "medal", "referee",

	// Celebrations & culture
	"Christmas tree", "fireworks", "Halloween pumpkin", "Easter egg",
	"wedding cake", "party hat", "balloon animal", "disco ball",
	"carol singers", "nativity", "maypole", "pancake day",
	"fancy dress", "school disco", "village fete", "car boot sale",

	// Jobs & people
	"doctor", "firefighter", "astronaut", "pirate",
	"wizard", "knight", "chef", "builder",
	"postman", "farmer", "detective", "clown",
	"lifeguard", "pilot", "dentist", "hairdresser",

	// Things & objects
	"treasure chest", "hot air balloon", "rocket", "pirate ship",
	"treehouse", "robot", "telescope", "anchor",
	"crown", "sword", "shield", "catapult",
	"candle", "cuckoo clock", "jigsaw puzzle", "lava lamp",
	"roller coaster", "haunted house", "UFO", "volcano",

	// Funny scenes
	"dad at the barbecue", "stepping on Lego", "dog stealing a sausage",
	"wasp at a picnic", "burnt toast", "man vs seagull",
	"slipping on a banana", "cat on a keyboard", "falling off a chair",
	"pigeon on your head", "stuck in traffic", "tangled headphones",
	"walking into a glass door", "sleeping on the train",
}

func RandomWords(n int) []string {
	perm := rand.Perm(len(wordList))
	words := make([]string, n)
	for i := 0; i < n; i++ {
		words[i] = wordList[perm[i%len(perm)]]
	}
	return words
}
