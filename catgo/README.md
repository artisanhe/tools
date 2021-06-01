# Catgo CAT 客户端集成

在 tools 的 **v1.2.0-paystd** 开始集成cat 的客户端，基于 github.com/Meituan-Dianping/cat-go 修改封装

Transaction

* transaction适合记录跨越系统边界的程序访问行为，比如远程调用，数据库调用，也适合执行时间较长的业务逻辑监控
* 某些运行期单元要花费一定时间完成工作, 内部需要其他处理逻辑协助, 我们定义为Transaction.
* Transaction可以嵌套(如http请求过程中嵌套了sql处理).
* 大部分的Transaction可能会失败, 因此需要一个结果状态码.
* 如果Transaction开始和结束之间没有其他消息产生, 那它就是Atomic Transaction(合并了起始标记).

Event

Event用来记录次数，表名单位时间内消息发生次数，比如记录系统异常，它和transaction相比缺少了时间的统计，开销比transaction要小

Metric
一共有三个API，分别用来记录次数、平均、总和，统一粒度为一分钟

* logMetricForCount用于记录一个指标值出现的次数
* logMetricForDuration用于记录一个指标出现的平均值
* logMetricForSum用于记录一个指标出现的总和


Heartbeat
这个是系统CAT客户端使用，应用程序不使用此API.
Heartbeta表示程序内定期产生的统计信息, 如CPU%, MEM%, 连接池状态, 系统负载等。

## 使用方式


### 配置

||||
|:---:|:---:|:----:|
|Domain|域|g7pay|
|Debug|是否输出debug 日志| true|
|HttpServers|接入服务器信息||

### 初始化

```
(&Catgo{}).Init()
```

在tools 中只需要配置在 config 中，默认会调用 Init 初始化

### api

**首先先确认cat 的引用**

```
import "git.chinawayltd.com/golib/tools/catgo/cat-go/cat"
```

示例

```
// send transaction
func case1() {
	t := cat.NewTransaction(TestType, "test")
	defer t.Complete()

	if rand.Int31n(100) == 0 {
		t.SetStatus(cat.FAIL)
	}

	t.AddData("foo", "bar")

	t.NewEvent(TestType, "event-1")
	t.Complete()

	if rand.Int31n(100) == 0 {
		t.LogEvent(TestType, "event-2", cat.FAIL)
	} else {
		t.LogEvent(TestType, "event-2")
	}
	t.LogEvent(TestType, "event-3", cat.SUCCESS, "k=v")

	t.SetDurationStart(time.Now().Add(-5 * time.Second))
	t.SetTime(time.Now().Add(-5 * time.Second))
	t.SetDuration(time.Millisecond * 500)
}

// send completed transaction with duration
func case2() {
	cat.NewCompletedTransactionWithDuration(TestType, "completed", time.Second*24)
	cat.NewCompletedTransactionWithDuration(TestType, "completed-over-60s", time.Second*65)
}

// send event
func case3() {
	// way 1
	e := cat.NewEvent(TestType, "event-4")
	e.Complete()
	// way 2

	if rand.Int31n(100) == 0 {
		cat.LogEvent(TestType, "event-5", cat.FAIL)
	} else {
		cat.LogEvent(TestType, "event-5")
	}
	cat.LogEvent(TestType, "event-6", cat.SUCCESS, "foobar")
}

// send error with backtrace
func case4() {
	if rand.Int31n(100) == 0 {
		err := errors.New("error")
		cat.LogError(err)
	}
}

// send metric
func case5() {
	cat.LogMetricForCount("metric-1")
	cat.LogMetricForCount("metric-2", 3)
	cat.LogMetricForDuration("metric-3", 150*time.Millisecond)
	cat.NewMetricHelper("metric-4").Count(7)
	cat.NewMetricHelper("metric-5").Duration(time.Second)
}

```


## 测试环境

Http服务地址： test.dubhe.chinawayltd.com:2280

查看地址： http://test.dubhe.chinawayltd.com/cat/r/t?domain=g7pay&ip=All&reportType=day&op=view