package log_test

import (
	"github.com/One-com/gonelog/log"
	"github.com/One-com/gonelog/syslog"
	"os"
	"time"
	"github.com/One-com/gonelog/context"
)

func output(l *log.Logger, msg string) {
	l.Output(2,msg)
}

func ExampleOutput() {
	l := log.New(os.Stdout,"",log.Lshortfile)
	l.DoCodeInfo(true)
	output(l, "output")
	// Output:
	// api_test.go:18: output
}

//------------------------------------

func h(l *log.Logger) {
	l.NOTICE("H1")
	log.NOTICE("H2")
}

func g(ctx context.Context) {
	ctx.NOTICE("G","a","b")
}

func f(ctx context.Context, ok chan<- struct{} ) {

	select {
	case <-ctx.Done():
		ctx.NOTICE("cancelled")
	}
	
	close(ok)
}


func ExampleContext() {
	log.SetPrefix("PFX:")
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LminFlags|log.Lshortfile)
	
	log.ERROR("hi",log.KV{"key":"val"})
	log.WARN("lazy", "hip", log.Lazy(func() interface{} {return "slow"}))
	
	cancel_context, cancel_func := context.WithCancel(context.Background())
	logger := log.Clone(log.KV{"tra":"la"})
	logger.DoCodeInfo(true)
	logging_context := context.WithLogging(cancel_context, logger)
	
	ok := make (chan struct{})
	
	go f(logging_context,ok)

	h(logger)
	g(logging_context)
	
	time.Sleep(time.Duration(1)*time.Second)

	cancel_func() // cancel the go-routine

	// wait for its exit
	select {
	case <- ok:
	}
	// Output:
	//[ERR]PFX:api_test.go:50: hi key=val
	//[WRN]PFX:api_test.go:51: lazy hip=slow
	//[NOT]PFX:api_test.go:26: H1 tra=la
	//[NOT]PFX:api_test.go:27: H2
	//[NOT]PFX:api_test.go:31: G tra=la a=b
	//[NOT]PFX:api_test.go:38: cancelled tra=la
}

//----------------------------------------------------------------

func ExampleGoneLogger() {
	h := log.NewStdFormatter(os.Stdout, "", log.Llevel)
	l := log.NewLogger(syslog.LOG_WARNING, h)
	l.SetDefaultLevel(syslog.LOG_NOTICE, false)

	// Traditional.
	// Evaluates arguments unless Lazy is used, but doesn't generate
	// Events above log level
	l.DEBUG("hej")
	l.INFO("hej")
	l.NOTICE("hej")
	l.WARN("hej")
	l.ERROR("hej")
	l.CRIT("hej")
	l.ALERT("hej")

	// Optimal
	// Doesn't do anything but checking the log level unless
	// something should be logged
	// A filtering handler would still be able to discard log events
	// based on level. Use Lazy to only evaluate just before formatting
	// Even by doing so a filtering writer might still discard the log line
	if f, ok := l.DEBUGok(); ok {
		f("dav")
	}
	if f, ok := l.INFOok(); ok {
		f("dav")
	}
	if f, ok := l.NOTICEok(); ok {
		f("dav")
	}
	if f, ok := l.WARNok(); ok {
		f("dav")
	}
	if f, ok := l.ERRORok(); ok {
		f("dav")
	}
	if f, ok := l.CRITok(); ok {
		f("dav")
	}
	if f, ok := l.ALERTok(); ok {
		f("dav")
	}

	// Primitive ... Allows for dynamically choosing log level.
	// Otherwise behaves like Traditional
	l.Log(syslog.LOG_DEBUG, "hop")
	l.Log(syslog.LOG_INFO, "hop")
	l.Log(syslog.LOG_NOTICE, "hop")
	l.Log(syslog.LOG_WARN, "hop")
	l.Log(syslog.LOG_ERROR, "hop")
	l.Log(syslog.LOG_CRIT, "hop")
	l.Log(syslog.LOG_ALERT, "hop")

	// Std logger compatible.
	// Will log with the default-level (default "INFO") - if that log-level is enabled.
	l.Print("default")
	// Fatal and Panic logs with level "ALERT"
	l.Fatal("fatal")
}

func ExamplePrintIgnores() {
	l := log.GetLogger("my/lib")
	h := log.NewFlxFormatter(log.SyncWriter(os.Stdout), "", log.Llevel|log.Lname)
	l.SetHandler(h)
	l.AutoColoring()
	l.SetLevel(syslog.LOG_ERROR)
	l.SetDefaultLevel(syslog.LOG_NOTICE,false)

	l.Print("ignoring level")
	// Output:
	// <5> (my/lib) ignoring level

}

func ExampleSubLogger() {
	l := log.GetLogger("my/lib")
	h := log.NewFlxFormatter(log.SyncWriter(os.Stdout), "", log.Llevel|log.Lname)
	l.SetHandler(h)
	l.SetLevel(syslog.LOG_ERROR)

	l2 := l.With("key","value")

	l3 := l2.With("more", "data")

	l3.ERROR("message")
	// Output:
	// <3> (my/lib) message more=data key=value

}

func ExampleNamedLogger() {
	l := log.GetLogger("my/lib")
	h := log.NewFlxFormatter(log.SyncWriter(os.Stdout), "", log.Llevel|log.Lname)
	l.SetHandler(h)
	l2 := log.GetLogger("my/lib/module")

	l3 := l2.With("k","v")
		
	l3.NOTICE("notice")
	// Output:
	// <5> (my/lib/module) notice k=v
}

func ExampleNamedClone() {
	l := log.Default()
	l.SetOutput(os.Stdout)
	l.SetFlags(log.Lname)
	l.SetPrefix("PFX")
	l2 := l.With("a","b")
	l3 := l2.NamedClone("myname")

	l3.ERROR("message")
	// Output:
	// PFXmessage a=b
}
