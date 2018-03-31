package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

type File struct {
	Filetype  string
	Url       string
	Extension string
}

type Result struct {
	Score  int64
	Title  string
	Author string
	Prize  string
	Files  []*File
}

type Competition struct {
	Competition string
	Results     []*Result
}

type GetResults struct {
	Succsess bool
	Action   string
	Data     *Competition
}

// The competitions
var contests = []int{
	13, //Fast Themed Graphics
	58, //Fast concept Visualization
	27, // Gamedev
	55, // Google Blocks
	56, // VR Demo
	25, //Demo
	7,  // Themed Music
	5,  // Freestyle Music
}

// We don't want this content
var contentBlacklist = map[int]bool{
	25: true,
}

func DownloadFile(location string, filename string, url string) (err error) {
	// Create the file
	os.MkdirAll(location, os.ModePerm)

	out, err := os.Create(location + "/" + filename)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func ReqGetResults(n int) *GetResults {
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest("GET", "https://unicorn.gathering.org/api/beamer/results/"+strconv.Itoa(n), nil)

	resp, err := client.Do(req)
	if err != nil {
		return nil
	}
	// Don't forget to close the response body
	defer resp.Body.Close()

	res := &GetResults{}

	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return res
}

func GetSortedWinners(entries []*Result) []*Result {
	var wrongEntries []*Result
	for _, entry := range entries {
		if entry.Prize == "" {
			break
		}
		wrongEntries = append(wrongEntries, entry)
	}
	ret := make([]*Result, len(wrongEntries))
	for j, n := len(wrongEntries), 0; j != 0; j, n = j-1, n+1 {
		ret[n] = wrongEntries[j-1]
	}
	return ret
}

var (
	dpi      = flag.Float64("dpi", 72, "screen resolution in Dots Per Inch")
	fontfile = flag.String("fontfile", "Helvetica-Condensed-Black.ttf", "filename of the ttf font")
	hinting  = flag.String("hinting", "full", "none | full")
	size     = flag.Float64("size", 115, "font size in points")
	spacing  = flag.Float64("spacing", 1.5, "line spacing (e.g. 2 means double spaced)")

	font1 = string("Helvetica-Black.ttf")
	font2 = string("Helvetica-Condensed-Black.ttf")
	text  = string("SEND PENGER")

	t_compo  = flag.String("c", "That was:", "Compo name/That was:")
	t_num    = flag.String("n", "03", "The number as a string")
	t_group  = flag.String("g", "TryHARDs", "Group/Author")
	t_msg    = flag.String("m", "…", "Message for the Crew")
	t_foot_l = string("Fast Themes Graphics")

	t_foot_r = string("Vote at competitions.gathering.org")

	res_w  = 1920
	res_h  = 1080
	border = 8
	res    = res_w * res_h

	// Font sizes. Should be a derived from res_w/res_h
	fsize_rect_compo = float64(111)
	fsize_rect_num   = float64(300)
	fsize_rect_title = float64(71)
	fsize_rect_group = float64(86)
	fsize_rect_msg   = float64(60)
	fsize_rect_footl = float64(100)
	fsize_rect_footr = float64(36)
)

func GenerateIntroSlide(comp *Competition) {
	// fmt.Printf("Loading fontfile %q\n", *fontfile)
	// fmt.Printf("Loading fontfile %q\n", test)

	// Hardcoded font 1
	// b1, err := ioutil.ReadFile(font1)
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }
	// f1, err := truetype.Parse(b1)
	// if err != nil {
	// 	log.Println("b1")
	// 	log.Println(err)
	// 	return
	// }

	// Hardcoded font 2
	b2, err := ioutil.ReadFile(font2)
	if err != nil {
		log.Println("font2")
		log.Println(err)
		return
	}
	f2, err := truetype.Parse(b2)
	if err != nil {
		log.Println("b2")
		log.Println(err)
		return
	}

	// // Generic font
	b, err := ioutil.ReadFile(*fontfile)
	if err != nil {
		log.Println(err)
		return
	}
	f, err := truetype.Parse(b)
	if err != nil {
		log.Println(err)
		return
	}

	// Freetype context
	// size_1 := float64((4*res_h/18 - 4*border) / 2)
	size_1 := float64(280)
	fmt.Printf("%+v\n", size_1)
	black, white, bg := image.Black, image.White, image.Transparent
	// white_a4 := color.RGBA{255, 255, 255, 172}
	// white_a5 := color.RGBA{255, 255, 255, 255}
	black_a4 := color.RGBA{100, 100, 100, 172}
	// black_a5 := color.RGBA{0, 0, 0, 214}
	yellow_1 := color.RGBA{255, 222, 0, 255}

	rgba := image.NewRGBA(image.Rect(0, 0, res_w, res_h))
	draw.Draw(rgba, rgba.Bounds(), bg, image.ZP, draw.Src)
	c := freetype.NewContext()
	c.SetDPI(*dpi)
	c.SetFont(f)
	c.SetFontSize(*size)
	c.SetFontSize(size_1)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(black)
	c.SetSrc(white)

	switch *hinting {
	default:
		c.SetHinting(font.HintingNone)
	case "full":
		c.SetHinting(font.HintingFull)
	}
	rect_msg := image.Rect(border, (res_h/3)*2, res_w-border, res_h/3)
	draw.Draw(rgba, rect_msg, &image.Uniform{black_a4}, image.ZP, draw.Src)

	// Fonts

	opts := truetype.Options{size_1, *dpi, font.HintingFull, 0, 64, 1}
	opts.Size = size_1

	opts.Size = fsize_rect_title
	face_num := truetype.NewFace(f2, &opts)
	face_footr := truetype.NewFace(f2, &opts)

	// sizes
	width_num_f := font.MeasureString(face_num, *t_num)
	width_num := width_num_f.Round()
	fmt.Printf("%+v\n", width_num)

	width_footr_f := font.MeasureString(face_footr, t_foot_l)
	width_footr := width_footr_f.Round()

	pt_foot_l := freetype.Pt((res_w/2)-(width_footr/2), (res_h / 2))
	c.SetSrc(&image.Uniform{yellow_1})
	c.SetFont(f2)
	c.SetFontSize(fsize_rect_footl)
	c.DrawString(t_foot_l, pt_foot_l)

	fmt.Printf("%+v\n", width_footr)
	// Save that RGBA image to disk.
	outFile, err := os.Create("out.png")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer outFile.Close()
	bf := bufio.NewWriter(outFile)
	err = png.Encode(bf, rgba)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	err = bf.Flush()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	fmt.Println("Wrote out.png OK.")
}

func GenerateWinnerSlide(entry *Result) {
}

func main() {
	flag.Parse()
	for _, v := range contests {
		l := ReqGetResults(v)
		if l != nil && l.Data.Results != nil {
			sortedEntries := GetSortedWinners(l.Data.Results)
			GenerateIntroSlide(l.Data)
			os.Exit(0)
			for _, entry := range sortedEntries {
				fmt.Println(entry.Prize)
			}
			// for _, entry := range l.Data.Results {

			// generate slide here

			// If we are in a blacklisted contest, move along
			// if val, ok := contentBlacklist[v]; ok {
			// 	continue
			// }

			// for _, file := range entry.Files {
			// 	DownloadFile("./output", entry.Title+"."+file.Extension, file.Url)
			// 	os.Exit(0)
			// }

			// }
		}
	}
}
