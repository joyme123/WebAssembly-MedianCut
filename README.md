中位切分法的go语言实现

**程序接收5个参数**

`-num`: int 输出的颜色数
`-img`: string 输入的图片
`-out`: string 输出的图片
`-debug`: bool 是否开启调试模式。调试模式下会输出调色板图片
`-fastMap`: bool 是否开启fastMap。开启fastMap效率会高一点，但是效果会差。


比如以下指令运行:

```
./main -num=256 -img=asset/girl.png -out=asset/girl-fastmap-out.png --fastMap=true
```

效果对比: [preview](https://github.com/joyme123/MedianCut/blob/master/asset/preview.md)

