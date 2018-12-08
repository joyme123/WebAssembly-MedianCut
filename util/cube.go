package util

type ColorCube struct {
	Lower uint16 // 盒子的低位指针
	Upper uint16 // 盒子的高位指针
	Count int    // 盒子中的颜色数量
	Level int    // 盒子切割的深度

	Rmin, Rmax, Gmin, Gmax, Bmin, Bmax byte // RGB值的范围

	HistPtr []uint16
}

// Shrink 方法会将盒子收缩
func (cube *ColorCube) Shrink() {
	var r, g, b byte
	var i, color uint16

	cube.Rmin = 255
	cube.Rmax = 0
	cube.Gmin = 255
	cube.Gmax = 0
	cube.Bmin = 255
	cube.Bmax = 0

	for i = cube.Lower; i <= cube.Upper; i++ {
		color = cube.HistPtr[i]
		r = RED(color)
		if r > cube.Rmax {
			cube.Rmax = r
		}

		if r < cube.Rmin {
			cube.Rmin = r
		}

		g = GREEN(color)
		if g > cube.Gmax {
			cube.Gmax = g
		}

		if g < cube.Gmin {
			cube.Gmin = g
		}

		b = BLUE(color)
		if b > cube.Bmax {
			cube.Bmax = b
		}

		if b < cube.Bmin {
			cube.Bmin = b
		}

	}
}

// GetColor 用来获取盒子代表的颜色
func (cube *ColorCube) GetColor() {

}

// HistList 封装了HistPtr，用来作排序
type HistList struct {
	HistPtr []uint16
	Longdim int
}

func (hist *HistList) Len() int {
	return len(hist.HistPtr)
}

func (hist *HistList) Less(i, j int) bool {
	colorA := hist.HistPtr[i]
	colorB := hist.HistPtr[j]

	var c1 byte
	var c2 byte

	// fmt.Println("执行londim", hist.Longdim)

	switch hist.Longdim {

	case 0:
		c1 = RED(colorA)
		c2 = RED(colorB)
	case 1:
		c1 = GREEN(colorA)
		c2 = GREEN(colorB)
	case 2:
		c1 = BLUE(colorA)
		c2 = BLUE(colorB)
	}
	return c1-c2 < 0

}

func (hist *HistList) Swap(i, j int) {
	hist.HistPtr[i], hist.HistPtr[j] = hist.HistPtr[j], hist.HistPtr[i]
}
