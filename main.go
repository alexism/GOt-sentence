package main

import (
	"bufio"
	"github.com/iris-contrib/middleware/logger"
	"github.com/iris-contrib/middleware/recovery"
	"github.com/kataras/iris"
	"log"
	"os"
	"regexp"
	"math/rand"
	"strings"
)

type Generator struct {
	firstPairs map[[2]string]bool
	words      []string
	data       map[string][]string
}

func NewGenerator(sentences []string) Generator {
	firstPairs := make(map[[2]string]bool, 0)
	var allWords []string
	data := make(map[string][]string)
	re := regexp.MustCompile("[[:space:]]")
	for _, s := range sentences {
		words := re.Split(s, -1)
		if len(words) < 2 {
			continue
		}

		var pair [2]string
		copy(pair[:], words[0:2])
		firstPairs[pair] = true
		allWords = append(allWords, words...)
		for i := 0; i < len(words) - 2; i++ {
			k := words[i] + words[i + 1]
			v := words[i + 2]
			// Appending to a nil slice just allocates a new slice
			data[k] = append(data[k], v)
		}
	}
	return Generator{firstPairs, allWords, data}
}
func (g *Generator) generate(size int) string {
	var sentence []string
	var w1, w2 string
	for i := 0; i < size; i++ {
		key := w1 + w2
		if _, exists := g.data[key]; !exists {
			p := g.pickSeedPair()
			w1 = p[0]
			w2 = p[1]
			key = w1 + w2
		}
		sentence = append(sentence, w1)
		pick := rand.Intn(len(g.data[key]))
		w3 := g.data[key][pick]
		w1 = w2
		w2 = w3
	}

	return strings.Join(sentence, " ")
}

func (g *Generator) pickSeedPair() [2]string {
	var pair [2]string
	var w1, w2 string

	for true {
		if _, exists := g.data[w1 + w2]; !exists {
			pair = g.pickKey()
			w1 = pair[0]
			w2 = pair[1]
		} else {
			break
		}

	}
	return pair
}

func (g *Generator) pickKey() [2]string {
	i := rand.Intn(len(g.words) - 3)
	pair := [2]string{g.words[i], g.words[i + 1]}
	log.Printf("%d %v", i, pair)
	return pair
}


func init() {
	log.SetFlags(log.Flags() | log.Lshortfile)
}

func main() {
	log.Println("starting...")

	iris.Use(recovery.New(os.Stderr))
	iris.Use(logger.New(iris.Logger))

	iris.Config.Render.Template.Directory = "./public"

	fileNames := os.Args[1:]
	sentences, err := Load(fileNames)
	if err != nil {
		log.Panicf("Failed to load files: %v", fileNames)
	}
	log.Printf("Number of sentences loaded: %d", len(sentences))

	iris.Get("/", func(ctx *iris.Context) {
		ctx.ServeFile("./public/index.html", false)
	})

	g := NewGenerator(sentences)

	iris.Get("/api/sentence", func(ctx *iris.Context) {
		s := g.generate(rand.Intn(20) + 10)
		log.Println(s)
		ctx.Text(iris.StatusOK, s)
	})

	iris.Static("/public", "./public", 1)
	iris.Listen(":8080")
}

func Load(files []string) ([]string, error) {
	lines := make([]string, 0, 20)
	for _, f := range files {
		log.Printf("loading %s", f)
		readLine(f, &lines)
	}
	return lines, nil
}

func readLine(path string, lines *[]string) {
	inFile, _ := os.Open(path)
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		*lines = append(*lines, scanner.Text())
	}
}

