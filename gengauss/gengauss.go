// Copyright © 2014 Lawrence E. Bakst. All rights reserved.
// THIS SOURCE CODE IS THE PROPRIETARY INTELLECTUAL PROPERTY AND CONFIDENTIAL
// INFORMATION OF LAWRENCE E. BAKST AND IS PROTECTED UNDER U.S. AND
// INTERNATIONAL LAW. ANY USE OF THIS SOURCE CODE WITHOUT THE PRIOR WRITTEN
// AUTHORIZATION OF LAWRENCE E. BAKST IS STRICTLY PROHIBITED.

package main

import "flag"
import "fmt"
import "math"
import "time"
import "math/rand"
import "leb/sp/stats"
import "leb/sp/point"

var a = flag.Int("a", 1, "a")
var b = flag.Int("b", 1000, "b")
var q = flag.Int("q", 100, "qualtiles")
var n = flag.Int("n", 10, "number of samples")
var t = flag.Int("t", 10, "number of trials")
var V = flag.Bool("v", false, "print samples")
var m = flag.Int("m", 10000000, "maximum number of samples")
var avg = flag.Float64("avg", 0.0, "mean of guassian distribution")
var sdv = flag.Float64("sdv", 1.0, "standard deviation of guassian distribution")
var r = rand.Float64

func rbetween(a int, b int) int {
	rf := r()
	diff := float64(b - a + 1)
	r2 := rf * diff
	r3 := r2 + float64(a)
	//	fmt.Printf("rbetween: a=%d, b=%d, rf=%f, diff=%f, r2=%f, r3=%f\n", a, b, rf, diff, r2, r3)
	ret := int(r3)
	return ret
}

// The polar form of the Box-Muller transformation is both faster and more robust numerically.
// http://www.design.caltech.edu/erik/Misc/Gaussian.html
// return a function that returns a random number corresponding to the guassian distribution
// with the specified mean and standard deviation
func GenGuassianFunc(mean, sd float64) (func() float64) {
var spare bool
var sv float64
var f = func() float64 {
	var u, v, s float64

		if spare {
			spare = false
			return mean + sd * sv
		}
		for {
			u = 2.0 * r() - 1.0
			v = 2.0 * r() - 1.0
			s = u * u + v * v
	        if s < 1.0 && s != 0.0 {
				break
	        }
	    }
		tmp := math.Sqrt( (-2.0 * math.Log( s ) ) / s )
		spare, sv = true, v * tmp
		return mean + sd * u * tmp
	}
	return f
}

func main() {
var values point.PointDatumSlice
var q99, q98 int

	flag.Parse()
	switch {
	case false:
	    f := GenGuassianFunc(*avg, *sdv)
		for i := 0; i < *n; i++ {
			v := f()
			if *V {
				fmt.Printf("%f\n", v)
			}
			values.AddPoint(&point.PointDatum{T: time.Now(), V: v}, *m)
		}
		s := stats.New(values, 10)
		fmt.Printf("s=%#v\n", s)
	    //t.Error("Sort out of order")
	// generate *n pages loads in *t trials; track how many page load times are in the 98th and 99th centile.
	case true:
		for i := 0; i < *t; i++ {
			values = nil
			for j := 0; j < *n; j++ {
				v := rbetween(*a, *b)
				if *V {
					fmt.Printf("%d\n", v)
				}
				values.AddPoint(&point.PointDatum{T: time.Now(), V: float64(v)}, *m)
			}
			s := stats.New(values, int64(*q), 1.0, 1000.0)
			tot := float64(0)
			for _, v := range s.Dec {
				tot += v
			}
			//fmt.Printf("tot=%f\n", tot)
			//fmt.Printf("s=%#v\n", s)
			if s.Dec[*q-1] > 0 {
				q99++
			}
			if s.Dec[*q-2] > 0 {
				q98++
			}
		}
		fmt.Printf("q98=%d, q99=%d, q98%%=%0.2f, q99%%=%0.2f, q98+q99%%=%0.2f\n", q98, q99, float64(q98)/float64(*t), float64(q99)/float64(*t),
			(float64(q98)+float64(q99))/float64(*t))
	}
}


