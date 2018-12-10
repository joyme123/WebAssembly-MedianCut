package mediancut

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
	"syscall/js"
)

func (api *WebAPI) setMaxCubeCb() {
	api.onMaxCubeCb = js.NewCallback(func(args []js.Value) {
		maxCubes, err := strconv.Atoi(args[0].String()) // 更新最大的颜色数量

		if err != nil {
			log.Fatal(err)
		}

		api.mc = NewMedianCut()

		reader := bytes.NewReader(api.inBuf)
		api.mc.Decode(reader, maxCubes)

		fmt.Println("开始设置maxCubes")

		api.mc.MedianCut(true)
		paletted := api.mc.Out(maxCubes, true)
		api.updateImage(paletted)
	})
}

func (api *WebAPI) setShutdownCb() {
	api.onShutdownCb = js.NewEventCallback(js.PreventDefault, func(e js.Value) {
		api.done <- struct{}{}
	})
}
