# SwiftX Express - Go

[API 文档](https://open.swiftx-express.com/)

## SwiftX邮编覆盖总表_251215_海外仓版

https://www.kdocs.cn/l/ck9uqAf4OVTk

## 测试数据说明

为了方便开发者进行API测试和集成调试，我们提供了一组预设的测试运单号。这些测试运单号可以直接用于轨迹查询接口，无需真实的订单数据：

- 已送达的运单号（可查询 POD 签收图片）：

```text
SWX784390000000365027
SWX295610000000373749
SWX847260000000377341
SWX531820000000384745
SWX672950000000387197
SWX418630000000378294
```

- 投递失败的运单号（无 POD 签收图片）：

```text
SWX847260000000377348
SWX923510000000485672
SWX756340000000629183
SWX682170000000794825
``` 

这些测试运单号包含完整的物流轨迹信息，其中状态为"已送达"的运单号还支持签收证明图片的下载。您可以使用这些运单号来：

测试 getTrackingInfo 和 batchGetTrackingInfo 接口

测试 downloadPodImages 和 batchDownloadPodImages 接口（仅限已送达的运单）

验证您的集成代码的正确性