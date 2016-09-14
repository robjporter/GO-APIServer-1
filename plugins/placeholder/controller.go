package placeholder

import (
	"io"
	"os"
	"fmt"
	"path"
	"time"
	"bytes"
	"image"
	"math"
	"runtime"
	"io/ioutil"
	"net/http"
	"image/color"
	"image/draw"
	"image/png"
	"strconv"
	"github.com/kataras/iris"
	"github.com/golang/freetype"
	//"github.com/golang/freetype/truetype"
	xfont "golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"github.com/ajstarks/svgo"
)

type myjson struct {
	Name string `json:"name"`
}

var (
	dpi = 72.00
	fontSize = 1.00
	maxTextBoundsToImageRatioY = 0.60
	maxTextBoundsToImageRatioX = 0.95
	backColour = "rgb(204, 204, 204)"
	frontColour = "rgb(150, 150, 150)"
)

func Index(ctx *iris.Context) {
	ctx.HTML(iris.StatusOK, "<b> Hello </b>")
}

func Draw(ctx *iris.Context) {
	width := ctx.Param("width")
	height := ctx.Param("height")
	text := ctx.Param("text")

	ctx.SetHeader("Content-Type","image/svg+xml")
	ctx.SetHeader("Cache-Control", "no-cache, no-store, must-revalidate")
	ctx.SetStatusCode(iris.StatusOK)
	var out io.Writer
	out = ctx.RequestCtx.Response.BodyWriter()
	
	intWidth,_ := strconv.Atoi(width)
	intHeight,_ := strconv.Atoi(height)
	if text == "" {text = width + " x " + height}
	intfontSize := int(fontSize)

	_, filename, _, _ := runtime.Caller(0)
	ttfPath := path.Join(path.Dir(filename),"fonts/Calibri.ttf")
	if _, err := os.Stat(ttfPath); err != nil {
		ctx.Write(err.Error())
	}

	fontBytes, err := ioutil.ReadFile(ttfPath)
	if err != nil {
		ctx.Write(err.Error())
	}
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		ctx.Write(err.Error())
	}
	
	c := freetype.NewContext()
	c.SetDPI(dpi)
	c.SetFont(font)
	c.SetFontSize(float64(intfontSize))
	c.SetHinting(xfont.HintingNone)

	var textExtent fixed.Point26_6
	drawPoint := freetype.Pt(0, int(PointToInt26_6(float64(intfontSize), dpi)>>6))
	textExtent, err = c.DrawString(text, drawPoint)
	if err != nil {
		ctx.Write(err.Error())
	}
	
	scaleX := float64(PointToInt26_6(float64(intWidth)*maxTextBoundsToImageRatioX, dpi)) / float64(textExtent.X)
	scaleY := float64(PointToInt26_6(float64(intHeight)*maxTextBoundsToImageRatioY, dpi)) / float64(textExtent.Y)
	fontSize2 := float64(intfontSize) * math.Min(scaleX, scaleY)
	
	canvas := svg.New(out)
	canvas.Start(intWidth, intHeight)
	canvas.Rect(0, 0, intWidth, intHeight, "fill:"+backColour)
	offset := 4
	canvas.Text(intWidth/2, (intHeight/2)+(int(fontSize2)/offset), text, "text-anchor:middle;font-family:Calibri;font-size:"+strconv.Itoa(int(fontSize2))+"px;font-weight:bold;fill:"+frontColour)
	canvas.End()
}

func Blank(ctx *iris.Context) {
	width := ctx.Param("width")
	height := ctx.Param("height")
	
	date := time.Now().Format(http.TimeFormat)
	ctx.SetHeader("Content-Type","image/svg+xml")
	ctx.SetHeader("Cache-Control", "no-cache, no-store, must-revalidate")
	ctx.SetHeader("Date", date)
	ctx.SetHeader("Expires", date)
	var out io.Writer
	out = ctx.RequestCtx.Response.BodyWriter()
	
	intWidth,_ := strconv.Atoi(width)
	intHeight,_ := strconv.Atoi(height)
	
	if intWidth < 1 {intWidth = 1}
	if intHeight == 0 {intHeight = intWidth}
	
	canvas := svg.New(out)
	canvas.Start(intWidth, intHeight)
	canvas.Rect(0, 0, intWidth, intHeight, "fill:"+backColour)
	canvas.End()
}


func Save(ctx *iris.Context) {
	width := ctx.Param("width")
	height := ctx.Param("height")
	text := ctx.Param("text")
	foreground := color.RGBA{150, 150, 150, 255}
	background := color.RGBA{204, 204, 204, 255}
	
	intWidth,_ := strconv.Atoi(width)
	intHeight,_ := strconv.Atoi(height)
	if intWidth < 1 {intWidth = 1}
	if intHeight == 0 {intHeight = intWidth}
	
	if text == "" {
		text = strconv.Itoa(intWidth) + " x " + strconv.Itoa(intHeight)
	}

	_, filename, _, _ := runtime.Caller(0)
	ttfPath := path.Join(path.Dir(filename),"fonts/DejaVuSans-Bold.ttf")
	if _, err := os.Stat(ttfPath); err != nil {
		ctx.Write(err.Error())
	}

	fontBytes, err := ioutil.ReadFile(ttfPath)
	if err != nil {
		ctx.Write(err.Error())
	}
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		ctx.Write(err.Error())
	}
	
	fg_img := image.NewUniform(foreground)
	bg_img := image.NewUniform(background)
	test_img := image.NewRGBA(image.Rect(0, 0, 0, 0))

	c := freetype.NewContext()
	c.SetDPI(dpi)
	c.SetFont(font)
	c.SetFontSize(fontSize)
	c.SetSrc(fg_img)
	c.SetDst(test_img)
	c.SetClip(test_img.Bounds())
	c.SetHinting(xfont.HintingNone)
	
	var textExtent fixed.Point26_6
	drawPoint := freetype.Pt(0, int(PointToInt26_6(fontSize, dpi)>>6))
	textExtent, err = c.DrawString(text, drawPoint)
	if err != nil {
		ctx.Write(err.Error())
	}

	scaleX := float64(PointToInt26_6(float64(intWidth)*maxTextBoundsToImageRatioX, dpi)) / float64(textExtent.X)
	scaleY := float64(PointToInt26_6(float64(intHeight)*maxTextBoundsToImageRatioY, dpi)) / float64(textExtent.Y)
	fontSize = fontSize * math.Min(scaleX, scaleY)

	c.SetFontSize(fontSize)
	drawPoint = freetype.Pt(0, 0)
	textExtent, err = c.DrawString(text, drawPoint)
	if err != nil {
		ctx.Write(err.Error())	
	}
	
	drawPoint = freetype.Pt(
		int(PointToInt26_6(float64(intWidth)/2.0, dpi)-textExtent.X/2)>>6,
		int(PointToInt26_6(float64(intHeight)/2.0+fontSize/2.6, dpi))>>6)

	img := image.NewRGBA(image.Rect(0, 0, intWidth, intHeight))
	draw.Draw(img, img.Bounds(), bg_img, image.ZP, draw.Src)

	c.SetDst(img)
	c.SetClip(img.Bounds())
	
	_, err = c.DrawString(text, drawPoint)
	if err != nil {
		ctx.Write(err.Error())
	}
	
	outputFilename := strconv.FormatInt(makeTimestamp(), 10)
	outputFilename = outputFilename + ".png"
	
	buffer := new(bytes.Buffer)
	png.Encode(buffer,img)
	fso, err := os.Create(path.Join(path.Dir(filename),"output/" + outputFilename))
	if err == nil {
		png.Encode(fso,img)
	}
	fso.Close()
	
	ctx.HTML(iris.StatusOK, "Saved - <a href='http://" + ctx.HostString() + "/placeholder/output/" + outputFilename + "'>View</a>")
}

func Base64(ctx *iris.Context) {
	width := ctx.Param("width")
	height := ctx.Param("height")
	text := ctx.Param("text")
	fmt.Println(width)
	fmt.Println(height)
	fmt.Println(text)
	ctx.Write("BASE64")
}

func oldDraw(ctx *iris.Context) {
	width := ctx.Param("width")
	height := ctx.Param("height")
	text := ctx.Param("text")
	foreground := color.RGBA{150, 150, 150, 255}
	background := color.RGBA{204, 204, 204, 255}
	
	intWidth,_ := strconv.Atoi(width)
	intHeight,_ := strconv.Atoi(height)
	if intWidth < 1 {intWidth = 1}
	if intHeight == 0 {intHeight = intWidth}
	
	if text == "" {
		text = strconv.Itoa(intWidth) + " x " + strconv.Itoa(intHeight)
	}

	_, filename, _, _ := runtime.Caller(0)
	ttfPath := path.Join(path.Dir(filename),"fonts/DejaVuSans-Bold.ttf")
	if _, err := os.Stat(ttfPath); err != nil {
		ctx.Write(err.Error())
	}

	fontBytes, err := ioutil.ReadFile(ttfPath)
	if err != nil {
		ctx.Write(err.Error())
	}
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		ctx.Write(err.Error())
	}
	
	fg_img := image.NewUniform(foreground)
	bg_img := image.NewUniform(background)
	test_img := image.NewRGBA(image.Rect(0, 0, 0, 0))

	c := freetype.NewContext()
	c.SetDPI(dpi)
	c.SetFont(font)
	c.SetFontSize(fontSize)
	c.SetSrc(fg_img)
	c.SetDst(test_img)
	c.SetClip(test_img.Bounds())
	c.SetHinting(xfont.HintingNone)
	
	var textExtent fixed.Point26_6
	drawPoint := freetype.Pt(0, int(PointToInt26_6(fontSize, dpi)>>6))
	textExtent, err = c.DrawString(text, drawPoint)
	if err != nil {
		ctx.Write(err.Error())
	}

	scaleX := float64(PointToInt26_6(float64(intWidth)*maxTextBoundsToImageRatioX, dpi)) / float64(textExtent.X)
	scaleY := float64(PointToInt26_6(float64(intHeight)*maxTextBoundsToImageRatioY, dpi)) / float64(textExtent.Y)
	fontSize = fontSize * math.Min(scaleX, scaleY)

	c.SetFontSize(fontSize)
	drawPoint = freetype.Pt(0, 0)
	textExtent, err = c.DrawString(text, drawPoint)
	if err != nil {
		ctx.Write(err.Error())	
	}
	
	drawPoint = freetype.Pt(
		int(PointToInt26_6(float64(intWidth)/2.0, dpi)-textExtent.X/2)>>6,
		int(PointToInt26_6(float64(intHeight)/2.0+fontSize/2.6, dpi))>>6)

	img := image.NewRGBA(image.Rect(0, 0, intWidth, intHeight))
	draw.Draw(img, img.Bounds(), bg_img, image.ZP, draw.Src)

	c.SetDst(img)
	c.SetClip(img.Bounds())
	
	_, err = c.DrawString(text, drawPoint)
	if err != nil {
		ctx.Write(err.Error())
	}
	
	ctx.SetHeader("Content-Type","image/png")
	ctx.SetHeader("Cache-Control", "no-cache, no-store, must-revalidate")
	ctx.SetStatusCode(iris.StatusOK)
	var out io.Writer
	out = ctx.RequestCtx.Response.BodyWriter()
	png.Encode(out,img)
}

func oldBlank(ctx *iris.Context) {
	width := ctx.Param("width")
	height := ctx.Param("height")
	
	date := time.Now().Format(http.TimeFormat)
	ctx.SetHeader("Content-Type","image/png")
	ctx.SetHeader("Cache-Control", "no-cache, no-store, must-revalidate")
	ctx.SetHeader("Date", date)
	ctx.SetHeader("Expires", date)
		
	intWidth,_ := strconv.Atoi(width)
	intHeight,_ := strconv.Atoi(height)
	
	if intWidth < 1 {intWidth = 1}
	if intHeight == 0 {intHeight = intWidth}
	
	m := image.NewRGBA(image.Rect(0, 0, intWidth, intHeight))
	grey := color.RGBA{0, 0, 0, 50}
	draw.Draw(m, m.Bounds(), &image.Uniform{grey}, image.ZP, draw.Src)
	buffer := new(bytes.Buffer)
	png.Encode(buffer, m)
	ctx.Write(string(buffer.Bytes()))
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / (int64(time.Millisecond)/int64(time.Nanosecond))
}

func PointToInt26_6(x, dpi float64) fixed.Int26_6 {
	return fixed.Int26_6(x * dpi * (64.0 / 72.0))
}