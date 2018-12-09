// Package mediancut 包封装了中位切分法相关的算法
package mediancut

import (
	"image"
	"image/color"
	"image/png"
	"io"
	"log"
	"os"
	"sort"

	"github.com/joyme123/WebAssembly-MedianCut/util"
)

// MedianCut 是中位切分法的封装
type MedianCut struct {
	HSIZE     uint16
	MAXCOLOR  int
	HistPtr   []uint16
	CubeList  []util.ColorCube
	LongDim   int
	Hist      []uint16  // 图片的颜色统计直方图
	ColorMap  [][3]byte //切割后的立方体对应的颜色
	MaxCubes  int
	m         image.Image // 切割的图片
	imgWidth  int
	imgHeight int
}

// NewMedianCut 负责初始化MedianCut中一些必要参数
func NewMedianCut() *MedianCut {

	mc := &MedianCut{}
	mc.HSIZE = 32768
	mc.MAXCOLOR = 256

	mc.HistPtr = make([]uint16, mc.HSIZE)
	mc.CubeList = make([]util.ColorCube, mc.MAXCOLOR)

	return mc
}

func (mc *MedianCut) Read(filepath string, maxCubes int) {

	reader, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	mc.Decode(reader, maxCubes)

}

// Decode 接收reader
func (mc *MedianCut) Decode(reader io.Reader, maxCubes int) {

	mc.MaxCubes = maxCubes

	mc.ColorMap = make([][3]byte, maxCubes)

	var imgErr error
	mc.m, _, imgErr = image.Decode(reader)
	if imgErr != nil {
		log.Fatal(imgErr)
	}

	bounds := mc.m.Bounds()

	mc.imgWidth = bounds.Size().X
	mc.imgHeight = bounds.Size().Y

	colorCnt := mc.imgWidth * mc.imgHeight

	mc.Hist = make([]uint16, colorCnt)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := mc.m.At(x, y).RGBA()

			r = r >> 8
			g = g >> 8
			b = b >> 8

			color := util.RGB((byte)(r), (byte)(g), (byte)(b))

			mc.Hist[color]++
		}
	}
}

// MedianCut 用来切割ColorCude, hist是图片中颜色的直方图
func (mc *MedianCut) MedianCut(fastMap bool) int {
	var lr, lg, lb byte
	var i, median, color uint16
	var count int
	var k, level, ncubes, splitpos int

	var cube util.ColorCube

	ncubes = 0
	cube.Count = 0

	i = 0
	color = 0

	for ; i < mc.HSIZE-1; i++ {
		if mc.Hist[i] != 0 {
			mc.HistPtr[color] = i
			color++
			cube.Count = cube.Count + (int)(mc.Hist[i])
		}
	}

	cube.Lower = 0
	cube.Upper = color - 1
	cube.Level = 0

	cube.HistPtr = mc.HistPtr[:]
	cube.Shrink()
	mc.CubeList[ncubes] = cube
	ncubes++

	for ncubes < mc.MaxCubes {
		level = 255
		splitpos = -1

		for k = 0; k <= ncubes-1; k++ {
			if mc.CubeList[k].Lower == mc.CubeList[k].Upper {

			} else if mc.CubeList[k].Level < level {
				level = mc.CubeList[k].Level
				splitpos = k
			}
		}

		if splitpos == -1 {
			break
		}

		cube = mc.CubeList[splitpos]
		lr = cube.Rmax - cube.Rmin
		lg = cube.Gmax - cube.Gmin
		lb = cube.Bmax - cube.Bmin

		if lr >= lg && lr >= lb {
			mc.LongDim = 0
		}

		if lg >= lr && lg >= lb {
			mc.LongDim = 1
		}

		if lb >= lr && lb >= lg {
			mc.LongDim = 2
		}

		histList := &(util.HistList{HistPtr: mc.HistPtr[cube.Lower:cube.Upper], Longdim: mc.LongDim})

		sort.Sort(histList)

		count = 0
		for i = cube.Lower; i <= cube.Upper-1; i++ {
			if count >= cube.Count/2 {
				break
			}
			color = mc.HistPtr[i]
			count = count + (int)(mc.Hist[color])
		}

		median = i

		cubeA := cube
		cubeA.Upper = median - 1
		cubeA.Count = count
		cubeA.Level = cube.Level + 1
		cubeA.Shrink()
		mc.CubeList[splitpos] = cubeA

		cubeB := cube
		cubeB.Lower = median
		cubeB.Count = cube.Count - count
		cubeB.Level = cube.Level + 1
		cubeB.Shrink()
		mc.CubeList[ncubes] = cubeB
		ncubes++
	}

	// 得到了足够的切割后的cube，现在计算所有方块的颜色,做颜色的映射
	mc.invMap(ncubes, fastMap)

	return ncubes
}

func (mc *MedianCut) invMap(ncubes int, fastMap bool) {
	var r, g, b byte
	var i, k, color uint16
	var rsum, gsum, bsum float32
	var cube util.ColorCube

	for k = 0; k <= (uint16)(ncubes)-1; k++ {
		cube = mc.CubeList[k]
		rsum = 0.0
		gsum = 0.0
		bsum = 0.0

		// fmt.Printf("upper是%d, lower是%d\n", cube.Upper, cube.Lower)
		for i = cube.Lower; i <= cube.Upper; i++ {
			// fmt.Printf("i是%d\n", i)
			color = mc.HistPtr[i]
			r = util.RED(color)
			rsum += (float32)(r) * (float32)(mc.Hist[color])
			g = util.GREEN(color)
			gsum += (float32)(g) * (float32)(mc.Hist[color])
			b = util.BLUE(color)
			bsum += (float32)(b) * (float32)(mc.Hist[color])
		}

		mc.ColorMap[k][0] = (byte)(rsum / (float32)(cube.Count))
		mc.ColorMap[k][1] = (byte)(gsum / (float32)(cube.Count))
		mc.ColorMap[k][2] = (byte)(bsum / (float32)(cube.Count))
	}

	if fastMap {

		for k = 0; k < (uint16)(ncubes); k++ {
			cube = mc.CubeList[k]
			for i = cube.Lower; i <= cube.Upper; i++ {
				color = mc.HistPtr[i]
				mc.Hist[color] = k
			}
		}
	} else {
		var dmin, dr, dg, db, d float32
		var index uint16

		for k = 0; k < (uint16)(ncubes); k++ {
			cube = mc.CubeList[k]
			for i = cube.Lower; i <= cube.Upper; i++ {
				color = mc.HistPtr[i]
				r = util.RED(color)
				g = util.GREEN(color)
				b = util.BLUE(color)

				/* Search for closest entry in "ColMap" */
				dmin = 99999999999999
				for j := 0; j < ncubes; j++ {
					dr = (float32)(mc.ColorMap[j][0]) - (float32)(r)
					dg = (float32)(mc.ColorMap[j][1]) - (float32)(g)
					db = (float32)(mc.ColorMap[j][2]) - (float32)(b)
					d = dr*dr + dg*dg + db*db
					if d == 0.0 {
						index = (uint16)(j)
						break
					} else if d < dmin {
						dmin = d
						index = (uint16)(j)
					}
				}
				mc.Hist[color] = index
			}
		}
	}

}

func (mc *MedianCut) Write(ncubes int, out string, debug bool) {

	paletted := mc.Out(ncubes, debug)

	file, err := os.Create(out)

	if err != nil {
		log.Fatal(err)
	}

	err2 := png.Encode(file, paletted)

	if err2 != nil {
		log.Fatal(err2)
	}
}

// Out 向writer中写入
func (mc *MedianCut) Out(ncubes int, debug bool) *image.Paletted {

	var palette []color.Color

	palette = make([]color.Color, ncubes)

	pIndex := 0

	for colorMapIndex, rgb := range mc.ColorMap {

		if colorMapIndex >= ncubes {
			break
		}

		rc := rgb[0]
		gc := rgb[1]
		bc := rgb[2]

		rgba := color.RGBA{rc, gc, bc, 255}
		palette[pIndex] = rgba
		pIndex++
	}

	// 输出调色板
	if debug {
		// debugPaletted := image.NewPaletted(image.Rect(0, 0, 640, 640), palette)
		// for y := 0; y < mc.imgHeight; y++ {
		// 	for x := 0; x < mc.imgWidth; x++ {
		// 		colorIndex := (x / 40) + 16*(y/40)

		// 		debugPaletted.SetColorIndex(x, y, (uint8)(colorIndex))
		// 	}
		// }

		// file, err := os.Create("debug.png")

		// if err != nil {
		// 	log.Fatal(err)
		// }

		// err2 := png.Encode(file, debugPaletted)

		// if err2 != nil {
		// 	log.Fatal(err2)
		// }
	}

	// 要生成改变颜色后的图片,这里使用调色板模式
	paletted := image.NewPaletted(image.Rect(0, 0, mc.imgWidth, mc.imgHeight), palette)

	for y := 0; y < mc.imgHeight; y++ {
		for x := 0; x < mc.imgWidth; x++ {

			r, g, b, _ := mc.m.At(x, y).RGBA()

			r = r >> 8
			g = g >> 8
			b = b >> 8

			color := util.RGB((byte)(r), (byte)(g), (byte)(b))
			paletted.SetColorIndex(x, y, (uint8)(mc.Hist[color]))
		}
	}

	return paletted
}
