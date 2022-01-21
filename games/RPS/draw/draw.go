//https://go.dev/blog/image-draw
// rock paper and scissor pictures: https://www.freepnglogos.com

package draw

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"strconv"
	"time"

	st "example.com/games/RPS/state"
	term "example.com/terminal"
)

// Consts and vars for image processing
const (
	width  = 400
	height = 200
	tmpDir = "games/RPS/draw/tmp/"
	srcDir = "games/RPS/draw/src/"
)

var bgSrc color.RGBA
var circleSrc image.Image
var yellowSrc image.Image
var whiteSrc image.Image
var redSrc image.Image

func init() {
	defaultColors()
	err := drawBackground()
	if err != nil {
		term.Print(term.ERROR, err.Error())
	}
}

// Resets image colors to default
func defaultColors() {
	bgSrc = color.RGBA{2, 48, 71, 255}
	circleSrc = &image.Uniform{color.RGBA{142, 202, 230, 255}}
	yellowSrc = &image.Uniform{color.RGBA{225, 183, 3, 255}}
	whiteSrc = &image.Uniform{color.RGBA{225, 225, 225, 255}}
	redSrc = &image.Uniform{color.RGBA{174, 32, 18, 255}}
}

// The following struct/functions are used to draw circles
type circle struct {
	p image.Point
	r int
}

func (c *circle) ColorModel() color.Model {
	return color.AlphaModel
}

func (c *circle) Bounds() image.Rectangle {
	return image.Rect(c.p.X-c.r, c.p.Y-c.r, c.p.X+c.r, c.p.Y+c.r)
}

func (c *circle) At(x, y int) color.Color {
	xx, yy, rr := float64(x-c.p.X)+0.5, float64(y-c.p.Y)+0.5, float64(c.r)
	if xx*xx+yy*yy < rr*rr {
		return color.Alpha{255}
	}
	return color.Alpha{0}
}

// Attempts to retrieve the file in /src/ and decode it as PNG
func getImg(file string) (image.Image, error) {
	openedFile, err := os.Open(srcDir + file)
	if err != nil {
		return nil, err
	}

	img, err := png.Decode(openedFile)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func newImg() *image.RGBA {
	topLeft := image.Point{0, 0}
	btmRight := image.Point{width, height}

	img := image.NewRGBA(image.Rectangle{topLeft, btmRight})

	return img
}

// Draws and saves a background /src/bg.png
func drawBackground() error {

	img := newImg()

	// Set Background color
	background := bgSrc
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img.Set(x, y, background)
		}
	}

	// Draw on PNG (destination)
	dst := draw.Image(img)

	// Draw circles
	draw.DrawMask(dst, dst.Bounds(), circleSrc, image.Point{}, &circle{image.Pt(80, 110), 50}, image.Point{}, draw.Over)
	draw.DrawMask(dst, dst.Bounds(), circleSrc, image.Point{}, &circle{image.Pt(320, 110), 50}, image.Point{}, draw.Over)

	// Draw Static Texts (using masks)
	mask, err := getImg("mask.png")
	if err != nil {
		return err
	}

	maskB := mask.Bounds()

	/*
		Reference: maskB.min     => top left coord of glyph (on mask.png)
				   Intersect     => rectangle from maskB.min to bottom right of glyph (relative)
				   Intersect.Add => Location on destination (dst)
	*/

	// "Rock Paper Scissors"
	rpsB := maskB.Intersect(image.Rect(0, 0, 164, 18)).Add(image.Point{115, 5})
	draw.DrawMask(dst, rpsB, whiteSrc, image.Point{}, mask, maskB.Min.Add(image.Point{60, 0}), draw.Over)

	// "Player"
	playerB := maskB.Intersect(image.Rect(0, 0, 58, 21)).Add(image.Point{50, 35})
	draw.DrawMask(dst, playerB, yellowSrc, image.Point{}, mask, maskB.Min, draw.Over)

	// "Bot"
	botB := maskB.Intersect(image.Rect(0, 0, 33, 17)).Add(image.Point{302, 35})
	draw.DrawMask(dst, botB, yellowSrc, image.Point{}, mask, maskB.Min.Add(image.Point{0, 21}), draw.Over)

	// "VS"
	vsB := maskB.Intersect(image.Rect(0, 0, 54, 32)).Add(image.Point{168, 92})
	draw.DrawMask(dst, vsB, whiteSrc, image.Point{}, mask, maskB.Min.Add(image.Point{0, 37}), draw.Over)

	// Encode background as PNG.
	bg_file, err := os.Create(srcDir + "bg.png")

	if err != nil {
		return err
	}

	png.Encode(bg_file, img)

	return nil
}

// Draws the player and bot moves onto destination image
func drawMoves(m1 st.Move, m2 st.Move, dst draw.Image) error {
	// Points to overlay on left or right circles
	leftPt := image.Pt(36, 70)
	rightPt := image.Pt(276, 70)

	// Images
	rock, err := getImg("rock.png")
	if err != nil {
		return err
	}
	paper, err := getImg("paper.png")
	if err != nil {
		return err
	}
	scissors, err := getImg("scissors.png")
	if err != nil {
		return err
	}

	// All three images are 90x90
	imgBounds := rock.Bounds()

	// Player move
	player := image.Rectangle{leftPt, leftPt.Add(imgBounds.Size())}
	switch m1 {
	case st.ROCK:
		draw.Draw(dst, player, rock, image.Point{}, draw.Over)
	case st.PAPER:
		draw.Draw(dst, player, paper, image.Point{}, draw.Over)
	case st.SCISSORS:
		draw.Draw(dst, player, scissors, image.Point{}, draw.Over)
	}

	// Bot move
	bot := image.Rectangle{rightPt, rightPt.Add(imgBounds.Size())}
	switch m2 {
	case st.ROCK:
		draw.Draw(dst, bot, rock, image.Point{}, draw.Over)
	case st.PAPER:
		draw.Draw(dst, bot, paper, image.Point{}, draw.Over)
	case st.SCISSORS:
		draw.Draw(dst, bot, scissors, image.Point{}, draw.Over)
	}

	return nil
}

// Generates a unique filename to store a RPS_GenerateImage request
// to 'tmp' folder. Only one game daily per user is saved.
func generateFilename(player string) string {
	const layout = "01-02-2006"
	t := time.Now()
	return tmpDir + player + t.Format(layout) + ".png"
}

var leftStatus image.Point = image.Point{50, 174}
var centerStatus image.Point = image.Point{166, 174}
var rightStatus image.Point = image.Point{290, 174}

// Generates a image detailing the results of an RPS game and returns a
// string denoting the file generated
// Input: User name, move, bot move, and result
func RPS_GenerateImage(playerName string, playerM st.Move, botM st.Move, res st.Result) (string, error) {
	// Create new image and set background/moves
	bg, err := getImg("bg.png")
	if err != nil {
		return "", err
	}
	dst := draw.Image(newImg())
	draw.Draw(dst, bg.Bounds(), bg, image.Point{}, draw.Over)
	drawMoves(playerM, botM, dst)

	mask, err := getImg("mask.png")
	if err != nil {
		return "", err
	}
	maskB := mask.Bounds()
	leftPt := image.Point{50, 174}
	centerPt := image.Point{166, 174}
	rightPt := image.Point{290, 174}

	switch res {
	case st.WIN:
		winB := maskB.Intersect(image.Rect(0, 0, 65, 13)).Add(leftPt)
		draw.DrawMask(dst, winB, redSrc, image.Point{}, mask, maskB.Min.Add(image.Point{60, 17}), draw.Over)
		loseB := maskB.Intersect(image.Rect(0, 0, 65, 13)).Add(rightPt)
		draw.DrawMask(dst, loseB, redSrc, image.Point{}, mask, maskB.Min.Add(image.Point{60, 31}), draw.Over)
	case st.LOSE:
		loseB := maskB.Intersect(image.Rect(0, 0, 65, 13)).Add(leftPt)
		draw.DrawMask(dst, loseB, redSrc, image.Point{}, mask, maskB.Min.Add(image.Point{60, 31}), draw.Over)
		winB := maskB.Intersect(image.Rect(0, 0, 65, 13)).Add(rightPt)
		draw.DrawMask(dst, winB, redSrc, image.Point{}, mask, maskB.Min.Add(image.Point{60, 17}), draw.Over)
	case st.TIE:
		tieB := maskB.Intersect(image.Rect(0, 0, 65, 13)).Add(centerPt)
		draw.DrawMask(dst, tieB, redSrc, image.Point{}, mask, maskB.Min.Add(image.Point{60, 45}), draw.Over)
	default:
		break
	}

	// Encode background as PNG.
	file := generateFilename(playerName)
	bg_file, _ := os.Create(file)
	png.Encode(bg_file, dst)

	return file, nil
}

// Sets a user-defined color for bg.png
// Returns True if successful
// Returns false with a string if there is a user-error
// Returns an error otherwise
func RPS_SetColor(args []string) (string, error) {

	cmd := args[0]

	if cmd == "set-default" {
		defaultColors()
		err := drawBackground()
		if err != nil {
			return "", err
		}
		return "RPS Game: Default colors Set!", nil
	}

	if len(args) < 4 {
		return "RPS Game Error: Insufficient number of arguments", nil
	}

	R, err1 := strconv.ParseInt(args[1], 10, 8)
	G, err2 := strconv.ParseInt(args[2], 10, 8)
	B, err3 := strconv.ParseInt(args[3], 10, 8)

	if err1 != nil || err2 != nil || err3 != nil {
		return "RPS Game Error: Invalid RGB values entered", nil
	}

	if R < 0 || R > 255 ||
		G < 0 || G > 255 ||
		B < 0 || B > 255 {
		return "RPS Game Error: RGB values must be between 0 and 255", nil
	}

	switch cmd {
	case "set-bg":
		bgSrc = color.RGBA{uint8(R), uint8(G), uint8(B), 255}
	case "set-circle":
		circleSrc = &image.Uniform{color.RGBA{uint8(R), uint8(G), uint8(B), 255}}
	default:
	}

	err := drawBackground()
	if err != nil {
		return "", err
	}
	return "RPS Game: Color set!", nil
}

// Cleans up file i/o related to this game
func RPS_DeleteTmpImage(dir string) {
	os.Remove(dir)
}
