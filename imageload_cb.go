package mediancut

import (
	"bytes"
	"fmt"
	"reflect"
	"syscall/js"
	"unsafe"
)

// setInitMemCb 是设置初始化内存时的回调函数
func (api *WebAPI) setInitMemCb() {
	api.onMemInitCb = js.NewCallback(func(args []js.Value) {
		length := args[0].Int()
		api.console.Call("log", "length", length) // 调用js的console.log("length", length)
		api.inBuf = make([]uint8, length)
		// 拿到这个slice的SliceHeader
		hdr := (*reflect.SliceHeader)(unsafe.Pointer(&api.inBuf))
		ptr := uintptr(unsafe.Pointer(hdr.Data))

		api.console.Call("log", "ptr:", ptr)
		js.Global().Call("gotMem", ptr)

		fmt.Println("初始化Mem成功")
	})

}

func (api *WebAPI) setOnImgLoadCb() {
	api.onImgLoadCb = js.NewCallback(func(args []js.Value) {

		fmt.Println("开始回调图片上传")

		reader := bytes.NewReader(api.inBuf)
		api.mc.Decode(reader, 256)

		js.Global().Get("document").Call("getElementById", "maxCubes").Set("value", 256)

		fmt.Println("开始设置maxCubes")

		api.mc.MedianCut(true)
		paletted := api.mc.Out(256, true)
		api.updateImage(paletted)
	})
}
