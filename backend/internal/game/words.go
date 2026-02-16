package game

import "math/rand"

var wordList = []string{
	"airplane", "alarm clock", "alien", "anchor", "angel",
	"ant", "apple", "arrow", "astronaut", "axe",
	"baby", "backpack", "balloon", "banana", "barn",
	"baseball", "basket", "bat", "beach", "bear",
	"bed", "bee", "bicycle", "bird", "birthday cake",
	"boat", "bomb", "book", "bowtie", "bridge",
	"broccoli", "broom", "bucket", "bug", "bus",
	"butterfly", "cactus", "cake", "camel", "camera",
	"campfire", "candle", "cannon", "car", "castle",
	"cat", "caterpillar", "chair", "cheese", "cherry",
	"chicken", "church", "clock", "cloud", "clown",
	"coffee", "comb", "computer", "cookie", "couch",
	"cow", "crab", "crayon", "crocodile", "crown",
	"cup", "cupcake", "diamond", "dice", "dinosaur",
	"dog", "dolphin", "donut", "door", "dragon",
	"drum", "duck", "eagle", "ear", "egg",
	"elephant", "envelope", "eye", "feather", "fence",
	"fire", "fish", "flag", "flashlight", "flower",
	"football", "fork", "fountain", "fox", "frog",
	"garden", "giraffe", "glasses", "globe", "glove",
	"goal", "gorilla", "grapes", "grass", "guitar",
	"hamburger", "hammer", "hand", "hat", "headphones",
	"heart", "helicopter", "hippo", "horse", "hot dog",
	"house", "ice cream", "igloo", "island", "jacket",
	"jellyfish", "kangaroo", "key", "kite", "knife",
	"ladder", "lamp", "leaf", "lemon", "lightbulb",
	"lighthouse", "lion", "lizard", "lock", "lollipop",
	"mailbox", "map", "megaphone", "mermaid", "microphone",
	"monkey", "moon", "mosquito", "mountain", "mouse",
	"mushroom", "nail", "necklace", "nest", "nose",
	"ocean", "octopus", "owl", "paintbrush", "palm tree",
	"panda", "parachute", "parrot", "peanut", "pencil",
	"penguin", "piano", "pie", "pig", "pillow",
	"pineapple", "pirate", "pizza", "planet", "present",
	"pumpkin", "rabbit", "raccoon", "radio", "rainbow",
	"robot", "rocket", "roller coaster", "rose", "sailboat",
	"sandwich", "saw", "scarecrow", "scissors", "scorpion",
	"shark", "sheep", "shoe", "skateboard", "skeleton",
	"skull", "snail", "snake", "snowflake", "snowman",
	"soccer ball", "spaceship", "spider", "spoon", "star",
	"starfish", "stop sign", "strawberry", "sun", "sunflower",
	"surfboard", "swan", "sword", "table", "teapot",
	"telephone", "telescope", "tent", "tiger", "toilet",
	"tomato", "tooth", "tornado", "tractor", "train",
	"treasure chest", "tree", "trophy", "truck", "trumpet",
	"turtle", "umbrella", "unicorn", "vacuum", "vampire",
	"vase", "volcano", "watermelon", "whale", "wheel",
	"windmill", "witch", "wizard", "worm", "zebra",
}

func RandomWords(n int) []string {
	perm := rand.Perm(len(wordList))
	words := make([]string, n)
	for i := 0; i < n; i++ {
		words[i] = wordList[perm[i%len(perm)]]
	}
	return words
}
