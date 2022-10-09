package data

import (
	"encoding/json"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/names"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/validate"
	"github.com/ethereum/go-ethereum/common"
	"github.com/prysmaticlabs/prysm/v3/crypto/hash"
	"github.com/prysmaticlabs/prysm/v3/crypto/rand"
)

func GetProjects() (projects []Project) {
	filepath.Walk(DataFolder()+"raw/", func(path string, info fs.FileInfo, err error) error {
		if strings.HasSuffix(path, ".json") {
			fp, err := os.Open(path)
			if err != nil {
				return err
			}
			defer fp.Close()
			b, _ := io.ReadAll(fp)
			w := Wrapper{}
			json.Unmarshal(b, &w)
			p := w.Data.Project
			if len(p.Id) > 0 && validate.IsValidAddress(p.WalletAddress) {
				p.WalletAddress = strings.ToLower(p.WalletAddress)
				p.Title = strings.Trim(strings.Replace(p.Title, "\n", " ", -1), " ")
				for i := 0; i < len(p.Addresses); i++ {
					p.Addresses[i].Address = strings.ToLower(p.Addresses[i].Address)
				}
				projects = append(projects, p)
			}
		}
		return nil
	})
	return projects
}

type SimpleProject struct {
	Address    string `json:"address"`
	GrantId    string `json:"grantId"`
	Name       string `json:"name"`
	Tag        string `json:"tag"`
	Active     bool   `json:"active"`
	Core       bool   `json:"core"`
	Categories string `json:"categories"`
}

var (
	adjectives = []string{"able", "above", "absolute", "accepted", "accurate", "ace", "active", "actual", "adapted", "adapting", "adequate", "adjusted", "advanced", "alert", "alive", "allowed", "allowing", "amazed", "amazing", "ample", "amused", "amusing", "apparent", "apt", "arriving", "artistic", "assured", "assuring", "awaited", "awake", "aware", "balanced", "becoming", "beloved", "better", "big", "blessed", "bold", "boss", "brave", "brief", "bright", "bursting", "busy", "calm", "capable", "capital", "careful", "caring", "casual", "causal", "central", "certain", "champion", "charmed", "charming", "cheerful", "chief", "choice", "civil", "classic", "clean", "clear", "clever", "climbing", "close", "closing", "coherent", "comic", "communal", "complete", "composed", "concise", "concrete", "content", "cool", "correct", "cosmic", "crack", "creative", "credible", "crisp", "crucial", "cuddly", "cunning", "curious", "current", "cute", "daring", "darling", "dashing", "dear", "decent", "deciding", "deep", "definite", "delicate", "desired", "destined", "devoted", "direct", "discrete", "distinct", "diverse", "divine", "dominant", "driven", "driving", "dynamic", "eager", "easy", "electric", "elegant", "emerging", "eminent", "enabled", "enabling", "endless", "engaged", "engaging", "enhanced", "enjoyed", "enormous", "enough", "epic", "equal", "equipped", "eternal", "ethical", "evident", "evolved", "evolving", "exact", "excited", "exciting", "exotic", "expert", "factual", "fair", "faithful", "famous", "fancy", "fast", "feasible", "fine", "finer", "firm", "first", "fit", "fitting", "fleet", "flexible", "flowing", "fluent", "flying", "fond", "frank", "free", "fresh", "full", "fun", "funky", "funny", "game", "generous", "gentle", "genuine", "giving", "glad", "glorious", "glowing", "golden", "good", "gorgeous", "grand", "grateful", "great", "growing", "grown", "guided", "guiding", "handy", "happy", "hardy", "harmless", "healthy", "helped", "helpful", "helping", "heroic", "hip", "holy", "honest", "hopeful", "hot", "huge", "humane", "humble", "humorous", "ideal", "immense", "immortal", "immune", "improved", "in", "included", "infinite", "informed", "innocent", "inspired", "integral", "intense", "intent", "internal", "intimate", "inviting", "joint", "just", "keen", "key", "kind", "knowing", "known", "large", "lasting", "leading", "learning", "legal", "legible", "lenient", "liberal", "light", "liked", "literate", "live", "living", "logical", "loved", "loving", "loyal", "lucky", "magical", "magnetic", "main", "major", "many", "massive", "master", "mature", "maximum", "measured", "meet", "merry", "mighty", "mint", "model", "modern", "modest", "moral", "more", "moved", "moving", "musical", "mutual", "national", "native", "natural", "nearby", "neat", "needed", "neutral", "new", "next", "nice", "noble", "normal", "notable", "noted", "novel", "obliging", "on", "one", "open", "optimal", "optimum", "organic", "oriented", "outgoing", "patient", "peaceful", "perfect", "pet", "picked", "pleasant", "pleased", "pleasing", "poetic", "polished", "polite", "popular", "positive", "possible", "powerful", "precious", "precise", "premium", "prepared", "present", "pretty", "primary", "prime", "pro", "probable", "profound", "promoted", "userprompt", "proper", "proud", "proven", "pumped", "pure", "quality", "quick", "quiet", "rapid", "rare", "rational", "ready", "real", "refined", "regular", "related", "relative", "relaxed", "relaxing", "relevant", "relieved", "renewed", "renewing", "resolved", "rested", "rich", "right", "robust", "romantic", "ruling", "sacred", "safe", "saved", "saving", "secure", "select", "selected", "sensible", "set", "settled", "settling", "sharing", "sharp", "shining", "simple", "sincere", "singular", "skilled", "smart", "smashing", "smiling", "smooth", "social", "solid", "sought", "sound", "special", "splendid", "square", "stable", "star", "steady", "sterling", "still", "stirred", "stirring", "striking", "strong", "stunning", "subtle", "suitable", "suited", "summary", "sunny", "super", "superb", "supreme", "sure", "sweeping", "sweet", "talented", "teaching", "tender", "thankful", "thorough", "tidy", "tight", "together", "tolerant", "top", "topical", "tops", "touched", "touching", "tough", "true", "trusted", "trusting", "trusty", "ultimate", "unbiased", "uncommon", "unified", "unique", "united", "up", "upright", "upward", "usable", "useful", "valid", "valued", "vast", "verified", "viable", "vital", "vocal", "wanted", "warm", "wealthy", "welcome", "welcomed", "well", "whole", "willing", "winning", "wired", "wise", "witty", "wondrous", "workable", "working", "worthy"}
	adverbs    = []string{"abnormally", "absolutely", "accurately", "actively", "actually", "adequately", "admittedly", "adversely", "allegedly", "amazingly", "annually", "apparently", "arguably", "awfully", "badly", "barely", "basically", "blatantly", "blindly", "briefly", "brightly", "broadly", "carefully", "centrally", "certainly", "cheaply", "cleanly", "clearly", "closely", "commonly", "completely", "constantly", "conversely", "correctly", "curiously", "currently", "daily", "deadly", "deeply", "definitely", "directly", "distinctly", "duly", "eagerly", "early", "easily", "eminently", "endlessly", "enormously", "entirely", "equally", "especially", "evenly", "evidently", "exactly", "explicitly", "externally", "extremely", "factually", "fairly", "finally", "firmly", "firstly", "forcibly", "formally", "formerly", "frankly", "freely", "frequently", "friendly", "fully", "generally", "gently", "genuinely", "ghastly", "gladly", "globally", "gradually", "gratefully", "greatly", "grossly", "happily", "hardly", "heartily", "heavily", "hideously", "highly", "honestly", "hopefully", "hopelessly", "horribly", "hugely", "humbly", "ideally", "illegally", "immensely", "implicitly", "incredibly", "indirectly", "infinitely", "informally", "inherently", "initially", "instantly", "intensely", "internally", "jointly", "jolly", "kindly", "largely", "lately", "legally", "lightly", "likely", "literally", "lively", "locally", "logically", "loosely", "loudly", "lovely", "luckily", "mainly", "manually", "marginally", "mentally", "merely", "mildly", "miserably", "mistakenly", "moderately", "monthly", "morally", "mostly", "multiply", "mutually", "namely", "nationally", "naturally", "nearly", "neatly", "needlessly", "newly", "nicely", "nominally", "normally", "notably", "noticeably", "obviously", "oddly", "officially", "only", "openly", "optionally", "overly", "painfully", "partially", "partly", "perfectly", "personally", "physically", "plainly", "pleasantly", "poorly", "positively", "possibly", "precisely", "preferably", "presently", "presumably", "previously", "primarily", "privately", "probably", "promptly", "properly", "publicly", "purely", "quickly", "quietly", "radically", "randomly", "rapidly", "rarely", "rationally", "readily", "really", "reasonably", "recently", "regularly", "reliably", "remarkably", "remotely", "repeatedly", "rightly", "roughly", "routinely", "sadly", "safely", "scarcely", "secondly", "secretly", "seemingly", "sensibly", "separately", "seriously", "severely", "sharply", "shortly", "similarly", "simply", "sincerely", "singularly", "slightly", "slowly", "smoothly", "socially", "solely", "specially", "steadily", "strangely", "strictly", "strongly", "subtly", "suddenly", "suitably", "supposedly", "surely", "terminally", "terribly", "thankfully", "thoroughly", "tightly", "totally", "trivially", "truly", "typically", "ultimately", "unduly", "uniformly", "uniquely", "unlikely", "urgently", "usefully", "usually", "utterly", "vaguely", "vastly", "verbally", "vertically", "vigorously", "violently", "virtually", "visually", "weekly", "wholly", "widely", "wildly", "willingly", "wrongly", "yearly"}
	nouns      = []string{"ox", "ant", "ape", "asp", "bat", "bee", "boa", "bug", "cat", "cod", "cow", "cub", "doe", "dog", "eel", "eft", "elf", "elk", "emu", "ewe", "fly", "fox", "gar", "gnu", "hen", "hog", "imp", "jay", "kid", "kit", "koi", "lab", "man", "owl", "pig", "pug", "pup", "ram", "rat", "ray", "yak", "bass", "bear", "bird", "boar", "buck", "bull", "calf", "chow", "clam", "colt", "crab", "crow", "dane", "deer", "dodo", "dory", "dove", "drum", "duck", "fawn", "fish", "flea", "foal", "fowl", "frog", "gnat", "goat", "grub", "gull", "hare", "hawk", "ibex", "joey", "kite", "kiwi", "lamb", "lark", "lion", "loon", "lynx", "mako", "mink", "mite", "mole", "moth", "mule", "mutt", "newt", "orca", "oryx", "pika", "pony", "puma", "seal", "shad", "slug", "sole", "stag", "stud", "swan", "tahr", "teal", "tick", "toad", "tuna", "wasp", "wolf", "worm", "wren", "yeti", "adder", "akita", "alien", "aphid", "bison", "boxer", "bream", "bunny", "burro", "camel", "chimp", "civet", "cobra", "coral", "corgi", "crane", "dingo", "drake", "eagle", "egret", "filly", "finch", "gator", "gecko", "ghost", "ghoul", "goose", "guppy", "heron", "hippo", "horse", "hound", "husky", "hyena", "koala", "krill", "leech", "lemur", "liger", "llama", "louse", "macaw", "midge", "molly", "moose", "moray", "mouse", "panda", "perch", "prawn", "quail", "racer", "raven", "rhino", "robin", "satyr", "shark", "sheep", "shrew", "skink", "skunk", "sloth", "snail", "snake", "snipe", "squid", "stork", "swift", "swine", "tapir", "tetra", "tiger", "troll", "trout", "viper", "wahoo", "whale", "zebra", "alpaca", "amoeba", "baboon", "badger", "beagle", "bedbug", "beetle", "bengal", "bobcat", "caiman", "cattle", "cicada", "collie", "condor", "cougar", "coyote", "dassie", "donkey", "dragon", "earwig", "falcon", "feline", "ferret", "gannet", "gibbon", "glider", "goblin", "gopher", "grouse", "guinea", "hermit", "hornet", "iguana", "impala", "insect", "jackal", "jaguar", "jennet", "kitten", "kodiak", "lizard", "locust", "maggot", "magpie", "mammal", "mantis", "marlin", "marmot", "marten", "martin", "mayfly", "minnow", "monkey", "mullet", "muskox", "ocelot", "oriole", "osprey", "oyster", "parrot", "pigeon", "piglet", "poodle", "possum", "python", "quagga", "rabbit", "raptor", "rodent", "roughy", "salmon", "sawfly", "serval", "shiner", "shrimp", "spider", "sponge", "tarpon", "thrush", "tomcat", "toucan", "turkey", "turtle", "urchin", "vervet", "walrus", "weasel", "weevil", "wombat", "anchovy", "anemone", "bluejay", "buffalo", "bulldog", "buzzard", "caribou", "catfish", "chamois", "cheetah", "chicken", "chigger", "cowbird", "crappie", "crawdad", "cricket", "dogfish", "dolphin", "firefly", "garfish", "gazelle", "gelding", "giraffe", "gobbler", "gorilla", "goshawk", "grackle", "griffon", "grizzly", "grouper", "haddock", "hagfish", "halibut", "hamster", "herring", "jackass", "javelin", "jawfish", "jaybird", "katydid", "ladybug", "lamprey", "lemming", "leopard", "lioness", "lobster", "macaque", "mallard", "mammoth", "manatee", "mastiff", "meerkat", "mollusk", "monarch", "mongrel", "monitor", "monster", "mudfish", "muskrat", "mustang", "narwhal", "oarfish", "octopus", "opossum", "ostrich", "panther", "peacock", "pegasus", "pelican", "penguin", "phoenix", "piranha", "polecat", "primate", "quetzal", "raccoon", "rattler", "redbird", "redfish", "reptile", "rooster", "sawfish", "sculpin", "seagull", "skylark", "snapper", "spaniel", "sparrow", "sunbeam", "sunbird", "sunfish", "tadpole", "termite", "terrier", "unicorn", "vulture", "wallaby", "walleye", "warthog", "whippet", "wildcat", "aardvark", "airedale", "albacore", "anteater", "antelope", "arachnid", "barnacle", "basilisk", "blowfish", "bluebird", "bluegill", "bonefish", "bullfrog", "cardinal", "chipmunk", "cockatoo", "crayfish", "dinosaur", "doberman", "duckling", "elephant", "escargot", "flamingo", "flounder", "foxhound", "glowworm", "goldfish", "grubworm", "hedgehog", "honeybee", "hookworm", "humpback", "kangaroo", "killdeer", "kingfish", "labrador", "lacewing", "ladybird", "lionfish", "longhorn", "mackerel", "malamute", "marmoset", "mamoswine", "moccasin", "mongoose", "monkfish", "mosquito", "pangolin", "parakeet", "pheasant", "pipefish", "platypus", "polliwog", "porpoise", "reindeer", "ringtail", "sailfish", "scorpion", "seahorse", "seasnail", "sheepdog", "shepherd", "silkworm", "squirrel", "stallion", "starfish", "starling", "stingray", "stinkbug", "sturgeon", "terrapin", "titmouse", "tortoise", "treefrog", "werewolf", "woodcock"}
)

// DeterministicName returns a deterministic triple of adverb-adjective-name
// given a random seed for initialization.
func DeterministicName(seed []byte, separator string) string {
	rng := rand.NewDeterministicGenerator()
	hashedValue := hash.FastSum64(seed)
	rng.Seed(int64(hashedValue)) // lint:ignore uintcast -- It's safe to do this for deterministic naming.
	adverb := adverbs[rng.Intn(len(adverbs)-1)]
	adjective := adjectives[rng.Intn(len(adjectives)-1)]
	name := nouns[rng.Intn(len(nouns)-1)]
	petname := []string{adverb, adjective, name}
	return strings.Join(petname, separator)
}

var namesMap *names.NamesMap

func GetChifraName(addrIn string) (*names.Name, error) {
	if namesMap == nil {
		nM, err := names.LoadNamesMap("mainnet")
		if err != nil {
			return nil, err
		}
		namesMap = &nM
	}
	a := common.HexToAddress(addrIn)
	name := (*namesMap)[a]
	if len(name.Name) == 0 {
		name.Name = DeterministicName([]byte(addrIn), " ")
	}
	name.Name = "[" + addrIn[:5] + "..." + addrIn[len(addrIn)-3:] + "] " + strings.Replace(strings.ToLower(name.Name), " ", "-", -1)
	return &name, nil
}

func ToSimpleProject(p *Project) SimpleProject {
	isCore := strings.HasPrefix(p.Id, "core-")
	cats := ""
	for _, c := range p.Categories {
		if len(cats) > 0 {
			cats += " "
		}
		cats += c.Name
	}

	return SimpleProject{
		Address:    p.WalletAddress,
		GrantId:    p.Id,
		Name:       p.Title,
		Tag:        "31-Projects:Giveth",
		Active:     p.Listed,
		Core:       isCore,
		Categories: cats,
	}
}

type Project struct {
	Id                  string           `json:"id"`
	Title               string           `json:"title"`
	Balance             float64          `json:"balance"`
	Image               string           `json:"image"`
	Slug                string           `json:"slug"`
	SlugHistory         []string         `json:"slugHistory"`
	CreationDate        string           `json:"creationDate"`
	UpdatedAt           string           `json:"updatedAt"`
	Admin               string           `json:"admin"`
	Description         string           `json:"description"`
	WalletAddress       string           `json:"walletAddress"`
	ImpactLocation      string           `json:"impactLocation"`
	QualityScore        int              `json:"qualityScore"`
	Verified            bool             `json:"verified"`
	TraceCampaignId     *string          `json:"traceCampaignId"`
	Listed              bool             `json:"listed"`
	GivingBlocksId      *string          `json:"givingBlocksId"`
	Status              Status           `json:"status"`
	Categories          []Category       `json:"categories"`
	Reaction            *Reaction        `json:"reaction"`
	AdminUser           User             `json:"adminUser"`
	Organization        Organization     `json:"organization"`
	Addresses           []NetworkAddress `json:"addresses"`
	TotalReactions      int              `json:"totalReactions"`
	TotalDonations      int              `json:"totalDonations"`
	TotalTraceDonations int              `json:"totalTraceDonations"`
}

func (p Project) String() string {
	b, _ := json.MarshalIndent(p, "", "  ")
	return string(b)
}

type NetworkAddress struct {
	Address     string `json:"address"`
	IsRecipient bool   `json:"isRecipient"`
	NetworkId   int    `json:"networkId"`
}

type Organization struct {
	Name                string `json:"name"`
	Label               string `json:"label"`
	SupportCustomTokens bool   `json:"supportCustomTokens"`
}

type User struct {
	Id            int     `json:"id"`
	Email         *string `json:"email"`
	FirstName     string  `json:"firstName"`
	WalletAddress string  `json:"walletAddress"`
}

type Reaction struct {
	Id int `json:"id"`
}

type Category struct {
	Name string `json:"name"`
}

type Status struct {
	Id          int    `json:"id"`
	Symbol      string `json:"symbol"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type DataWrapper struct {
	Project Project `json:"projectById"`
}

type Wrapper struct {
	Data DataWrapper `json:"data"`
}
