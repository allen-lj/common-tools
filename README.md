# common-tools
## 工具包
### 并发
#### Marmot(土拨鼠)
Marmot实现方式是以指定并发数同时运行，使用场景一般以同时并发请求的压力测试场景.
>Demo如下:
```
m := NewMarmot(queueLength,concurrentNum)
go m.StartWork()

m.AddProcessor(&MockProcessor{})
m.AddProcessor(&MockProcessor{})
...

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
  // 
}
```