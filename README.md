# common-tools
## 工具包
### 并发
在并发的世界，经常会遇到下面两种需求：  
1. 并发执行多个任务
2. 按照指定速率执行任务  

Marmot 和 BumbleBee两个工具类正好覆盖这两种场景：  
- Marmot是通过指定并发任务数来执行任务
- BumbleBee则是以固定速率的方式来并发执行任务

#### Marmot(土拨鼠)
Marmot 以每秒N个并发来同时执行任务。使用场景一般在同时并发执行的压力测试场景.
>Demo如下:
```
m := NewMarmot(queueLength,concurrentNum)
go m.StartWork()

m.AddProcessor(&MockProcessor{})
m.AddProcessor(&MockProcessor{})
...
// work done
m.WorkDone()
m.WaitForClose()

type MockProcessor struct {
}

func (m *MockProcessor) PreProcess()   {
  // do something before work
}
func (m *MockProcessor) DoProcess()    {
  // the main process
}
func (m *MockProcessor) AfterProcess() {
  // do something after work
}
```
#### BumbleBee(大黄蜂)
BumbleBee 通过生成令牌的方式来进行并发任务处理.
>Demo如下:
```
b := NewBumbleBee(3,1)
b.AddProcessor(&MockProcessor{})
b.AddProcessor(&MockProcessor{})
b.AddProcessor(&MockProcessor{})
go b.StartWork()

// work done
b.WorkDone()
b.WaitForClose()

```