//Benchmak
package main

import "testing"

import "github.com/forsoki/gohash/str2hash"
import "github.com/forsoki/gohash/mutation"

func BenchmarkToLeet(b *testing.B) {
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		mutations.ToLeet("This is a long string!")
	}
	b.StopTimer()
}

func BenchmarkStr2hashMD5(b *testing.B) {
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		str2hash.MD5("asdf")
	}
	b.StopTimer()
}

func BenchmarkStr2hashSHA1(b *testing.B) {
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		str2hash.SHA1("asdf")
	}
	b.StopTimer()
}
