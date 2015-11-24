package log

import (
	"io/ioutil"
	stdlog "log"
	"testing"
)

func BenchmarkGoStdPrintln(b *testing.B) {
	const testString = "test"
	l := stdlog.New(ioutil.Discard, "", LstdFlags)
	for i := 0; i < b.N; i++ {
		l.Println(testString)
	}
}

func BenchmarkStdPrintln(b *testing.B) {
	const testString = "test"
	l := New(ioutil.Discard, "", LstdFlags)
	for i := 0; i < b.N; i++ {
		l.Println(testString)
	}
}

func BenchmarkFlxPrintln(b *testing.B) {
	const testString = "test"
	h := NewFlxFormatter(ioutil.Discard, "", LstdFlags)
	l := NewLogger(LvlDEFAULT, h)
	l.DoTime(true)
	for i := 0; i < b.N; i++ {
		l.Println(testString)
	}
}

func BenchmarkParallelGoStdPrintln(b *testing.B) {
	const testString = "test"
	l := stdlog.New(ioutil.Discard, "", LstdFlags)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l.Println(testString)
		}
	})
}

func BenchmarkParallelStdPrintln(b *testing.B) {
	const testString = "test"
	l := New(ioutil.Discard, "", LstdFlags)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l.Println(testString)
		}
	})
}

func BenchmarkParallelFlxPrintln(b *testing.B) {
	const testString = "test"
	h := NewFlxFormatter(ioutil.Discard, "", LstdFlags)
	l := NewLogger(LvlDEFAULT, h)
	l.DoTime(true)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l.Println(testString)
		}
	})
}

func BenchmarkParallelMinPrintln(b *testing.B) {
	const testString = "test"
	h := NewMinFormatter(ioutil.Discard)
	l := NewLogger(LvlDEFAULT, h)
	l.DoTime(false)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l.Println(testString)
		}
	})
}

func BenchmarkParallelMinERRORok(b *testing.B) {
	const testString = "test"
	h := NewMinFormatter(ioutil.Discard)
	l := NewLogger(LvlDEFAULT, h)
	l.DoTime(false)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if f, ok := l.ERRORok(); ok {
				f(testString)
			}
		}
	})
}
