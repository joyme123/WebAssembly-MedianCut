package mediancut

import (
	"bytes"
	"image"
	"image/png"
	"log"
	"reflect"
	"syscall/js"
	"unsafe"
)

// WebAPI 是暴露给浏览器的接口
type WebAPI struct {
	onMaxCubeCb js.Callback // 最大颜色数的回调
	onImgLoadCb js.Callback // 图片加载回调
	onMemInitCb js.Callback // 内存初始化回调
	inBuf       []uint8     // reader
	outBuf      bytes.Buffer

	console js.Value
	done    chan struct{}

	mc *MedianCut
}

// Init 是初始化函数s
func (api *WebAPI) Init() {

	api.console = js.Global().Get("console")
	api.done = make(chan struct{})
	api.mc = NewMedianCut()

	api.setInitMemCb()
	js.Global().Set("initMem", api.onMemInitCb)

	api.setOnImgLoadCb()
	js.Global().Set("loadImage", api.onImgLoadCb)

	api.setMaxCubeCb()
	js.Global().Set("maxCubeChange", api.onMaxCubeCb)

	<-api.done
	api.onMemInitCb.Release()
	api.onImgLoadCb.Release()
	api.onMaxCubeCb.Release()
}

func (api *WebAPI) updateImage(img image.Image) {
	// api.mc.Out(&api.outBuf, img)

	err := png.Encode(&api.outBuf, img)

	if err != nil {
		log.Fatal(err)
	}

	out := api.outBuf.Bytes()

	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&out))
	ptr := uintptr(unsafe.Pointer(hdr.Data))

	js.Global().Call("displayImage", ptr, len(out))
	api.outBuf.Reset()
}
