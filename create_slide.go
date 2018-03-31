//usr/bin/env test -x $0 && (([ ! -x "${0}c" ] || [ "$0" -nt "${0}c" ]) && go build -o "${0}c" "$0"; "${0}c" $@); exit "$?"

package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"log"
	"os"
)

var (
	dpi      = flag.Float64("dpi", 72, "screen resolution in Dots Per Inch")
	fontfile = flag.String("fontfile", "/home/andy/TG/TG16_Beamer/Fonts/Helvetica-Condensed-Black.ttf", "filename of the ttf font")
	hinting  = flag.String("hinting", "full", "none | full")
	size     = flag.Float64("size", 115, "font size in points")
	spacing  = flag.Float64("spacing", 1.5, "line spacing (e.g. 2 means double spaced)")

	font1    = string("/home/andy/TG/TG16_Beamer/Fonts/Helvetica-Black.ttf")
	font2    = string("/home/andy/TG/TG16_Beamer/Fonts/Helvetica-Condensed-Black.ttf")
	text     = string("SEND PENGER")

	t_compo  = flag.String("c", "That was:", "Compo name/That was:")
	t_num    = flag.String("n", "03", "The number as a string")
	t_group  = flag.String("g", "TryHARDs", "Group/Author")
	t_msg    = flag.String("m", "…", "Message for the Crew")
	t_foot_l = flag.String("l", "Hei Mamma!", "Left footer")

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
	fsize_rect_footl = float64(36)
	fsize_rect_footr = float64(36)
)

var t_title  = flag.String("t", "YoloSWAG 4 Lyf", "Title of artwork")

func main() {
	flag.Parse()
	test := *t_title
	fmt.Printf("Loading fontfile %q\n", *fontfile)
	fmt.Printf("Loading fontfile %q\n", test)

	// Hardcoded font 1
	b1, err := ioutil.ReadFile(font1)
	if err != nil {
		log.Println(err)
		return
	}
	f1, err := truetype.Parse(b1)
	if err != nil {
		log.Println(err)
		return
	}

	// Hardcoded font 2
	b2, err := ioutil.ReadFile(font2)
	if err != nil {
		log.Println(err)
		return
	}
	f2, err := truetype.Parse(b2)
	if err != nil {
		log.Println(err)
		return
	}

	// Generic font
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
	//size_1 := float64((4*res_h/18 - 4 * border)/2)
	size_1 := float64(180)
	fmt.Printf("%+v\n", size_1)
	black, white, bg := image.Black, image.White, image.Transparent
	//white_a4 := color.RGBA{255, 255, 255, 172}
	white_a5 := color.RGBA{255, 255, 255, 255}
	black_a4 := color.RGBA{0, 0, 0, 172}
	black_a5 := color.RGBA{0, 0, 0, 214}
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

	// Define the boxes
	rect_compo := image.Rect(border, border, res_w - border, 3*res_h/18 - border)
	rect_num   := image.Rect(border, 3*res_h/18 + border, res_w/4 - border, 9*res_h/18 - border)
	rect_title := image.Rect(res_w/4 + border, 3*res_h/18 + border, res_w - border, 7*res_h/18 - border)
	rect_group := image.Rect(res_w/4 + border, 7*res_h/18 + border, res_w - border, 9*res_h/18 - border)
	rect_msg   := image.Rect(border, 9*res_h/18 + border, res_w - border, 17*res_h/18 - border)
	rect_foot  := image.Rect(border, 17*res_h/18 + border, res_w - border, 18*res_h/18 - border)

	// Draw the rectangles
	draw.Draw(rgba, rect_compo, &image.Uniform{black_a5}, image.ZP, draw.Src)
	draw.Draw(rgba, rect_num,   &image.Uniform{yellow_1}, image.ZP, draw.Src)
	draw.Draw(rgba, rect_title, &image.Uniform{white_a5}, image.ZP, draw.Src)
	draw.Draw(rgba, rect_group, &image.Uniform{black_a5}, image.ZP, draw.Src)
	draw.Draw(rgba, rect_msg,   &image.Uniform{black_a4}, image.ZP, draw.Src)
	draw.Draw(rgba, rect_foot,  &image.Uniform{black_a5}, image.ZP, draw.Src)

	// Truetype stuff
	opts := truetype.Options{size_1, *dpi, font.HintingFull, 0, 64, 1}
	opts.Size = size_1
	//face := truetype.NewFace(f, &opts)

	/*opts.Size = fsize_rect_compo
	face_compo := truetype.NewFace(f, &opts)*/

	opts.Size = fsize_rect_num
	face_num := truetype.NewFace(f2, &opts)

	/*
	opts.Size = fsize_rect_title
	face_title := truetype.NewFace(f, &opts)

	opts.Size = fsize_rect_group
	face_group := truetype.NewFace(f, &opts)

	opts.Size = fsize_rect_msg
	face_msg := truetype.NewFace(f, &opts)*/

	opts.Size = fsize_rect_footr
	face_footr := truetype.NewFace(f2, &opts)

	// String is how big?
	width_num_f := font.MeasureString(face_num, *t_num)
	width_num   := width_num_f.Round()
	fmt.Printf("%+v\n", width_num)

	width_footr_f := font.MeasureString(face_footr, t_foot_r)
	width_footr   := width_footr_f.Round()
	fmt.Printf("%+v\n", width_footr)

	// Compo
	pt_compo := freetype.Pt(int(0.5 * fsize_rect_compo), int(3*res_h/36) + int(0.377 * fsize_rect_compo))
	c.SetSrc(&image.Uniform{yellow_1})
	c.SetFont(f1)
	c.SetFontSize(fsize_rect_compo)
	c.DrawString(*t_compo, pt_compo)

	// Number
	pt_num := freetype.Pt((res_w/4 - width_num)/2, 6*res_h/18 + int(0.377 * fsize_rect_num))
	c.SetSrc(black)
	c.SetFont(f2)
	c.SetFontSize(fsize_rect_num)
	c.DrawString(*t_num, pt_num)

	// Title
	pt_title := freetype.Pt(res_w/4 + int(0.5 * fsize_rect_title), int(5*res_h/18) + int(0.377 * fsize_rect_title))
	c.SetSrc(black)
	c.SetFont(f2)
	c.SetFontSize(fsize_rect_title)
	c.DrawString(test, pt_title)

	// Group
	pt_group := freetype.Pt(res_w/4 + int(0.5 * fsize_rect_group), int(8*res_h/18) + int(0.377 * fsize_rect_group))
	c.SetSrc(&image.Uniform{yellow_1})
	c.SetFont(f2)
	c.SetFontSize(fsize_rect_group)
	c.DrawString("By: " + *t_group, pt_group)

	// Message
	pt_msg := freetype.Pt(int(0.5 * fsize_rect_msg), int(10*res_h/18) + int(0.377 * fsize_rect_msg))
	c.SetSrc(white)
	c.SetFont(f2)
	c.SetFontSize(fsize_rect_msg)
	c.DrawString(*t_msg, pt_msg)

	// Footer right
	pt_foot_l := freetype.Pt(2*border, int(35*res_h/36) + int(0.377 * fsize_rect_footl))
	c.SetSrc(&image.Uniform{yellow_1})
	c.SetFont(f2)
	c.SetFontSize(fsize_rect_footl)
	c.DrawString(*t_foot_l, pt_foot_l)

	// Footer right
	pt_foot_r := freetype.Pt(res_w - width_footr - 2*border, int(35*res_h/36) + int(0.377 * fsize_rect_footr))
	c.SetSrc(&image.Uniform{yellow_1})
	c.SetFont(f2)
	c.SetFontSize(fsize_rect_footr)
	c.DrawString(t_foot_r, pt_foot_r)



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
