
// Copyright 2018 The godag Authors
// This file is part of the godag library.
//
// The godag library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The godag library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the godag library. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"testing"
	"math"
	"math/big"
)
// Calculate K according to formula: k(D_{max},δ). For fraction λ
//
func TestCalK(t *testing.T) {

	Dmax := 15.0
	Lambdaλ := 0.2	// 0.1 ~ 0.9

	targetδ := big.NewFloat(0.001)

	t1 := 2.0 * Dmax * Lambdaλ
	t2 := math.Exp(float64(-t1))
	t3 := new(big.Float).Quo(big.NewFloat(t2), big.NewFloat(1.0-t2))

	println("Dmax=", Dmax, " λ=", Lambdaλ, " δ=", targetδ.String())
	println("t1='2*Dmax*λ'=", t1, " t2=exp(-2 * Dmax * λ)=", t2, " t3=t2/(1-t2)=", t3.String())
	println("")

	if t1-float64(int64(t1)) > 0 {
		t.Errorf("error! '2*Dmax*λ' must be an integer (easy to calculate). but now it's %f", t1)
		return
	}

	for k := 0; k <= 50; k++ {
		sum := big.NewInt(0)
		for j := k + 1; j <= k+1000; j++ { // instead of using ∞, we only loop 1000 to get rough sum.

			numerator := new(big.Int).Exp(big.NewInt(int64(t1)), big.NewInt(int64(j)), nil)
			jFactorial := new(big.Int).MulRange(1, int64(j)) // calculate j!

			t4 := new(big.Int).Div(numerator.Lsh(numerator, 32), jFactorial)
			sum = sum.Add(sum, t4)
			if t4.Cmp(big.NewInt(0)) == 0 {
				//println("k'=",k, " break at j=", j, " (2*Dmax*λ)^j=", numerator.String(), " (j!)=", jFactorial.String(), " t4=(2*Dmax*λ)^j/(j!)=", t4.String(), " sum=", sum.String())
				break
			}
		}

		sum.Rsh(sum, 32)

		sumFloat := new(big.Float)
		sumFloat.SetInt(sum)
		sumFloat.Mul(sumFloat, t3)

		result := sumFloat.Cmp(targetδ)
		// -1 if x < y
		if result < 0 {
			println("if k'=", k, " probability=", sumFloat.String(), " found it! this is the minimum k'.")
			break
		}else {
			println("if k'=", k, " probability=", sumFloat.String(), " sum=", sum.String())
		}
	}
}


// Calculate K according to formula: k(D_{max},δ).  For integer λ
//
func TestCalKBigRate(t *testing.T) {

	Dmax := 15.0
	Lambdaλ := 1.0

	targetδ := big.NewFloat(0.01 )

	t1 := 2.0 * Dmax * Lambdaλ
	t2 := math.Exp( float64(-t1) )
	t3 := new(big.Float).Quo(big.NewFloat(1.0-t2), big.NewFloat(t2*10000.0))

	println("Dmax=", Dmax, " λ=", Lambdaλ, " δ=", targetδ.String())
	println("t1='2*Dmax*λ'=", t1, " t2=exp(-2 * Dmax * λ)=", t2, " t3=(1-t2)/t2=", t3.String())
	println("")

	targetδ.Mul(targetδ, big.NewFloat(10000.0))

	if Lambdaλ < 1.0 {
		t.Errorf("error! λ must be integer (use 'TestCalK' for λ<1.0 ). but now it's %f", Lambdaλ)
		return
	}

	if t1 - float64(int64(t1)) > 0 {
		t.Errorf("error! '2*Dmax*λ' must be an integer (easy to calculate). but now it's %f", t1)
		return
	}

	t3BigInt := new(big.Int)
	t3.Int(t3BigInt)

	for k:=0; k<=1000; k++ {
		sum := big.NewInt(0)
		for j := k+1; j<=k+1000; j++ {	// instead of using ∞, we only loop 1000 to get rough sum.

			nominator := new(big.Int).Exp(big.NewInt(int64(t1)), big.NewInt(int64(j)), nil)
			jFactorial := new(big.Int).MulRange(1, int64(j))		// calculate (j!)

			t4 := new(big.Int).Div(nominator, jFactorial)
			sum = sum.Add(sum, t4)
			if t4.Cmp(big.NewInt(0)) == 0 {
				//println("k=",k, " j=", j, " (2*Dmax*λ)^j=", nominator.String(), " (j!)=", jFactorial.String(), " t4=", t4.String(), " sum=", sum.String())
				break
			}
		}

		dv := new(big.Int).Div(sum, t3BigInt)
		targetδInt := new(big.Int)
		targetδ.Int(targetδInt)
		result := dv.Cmp(targetδInt)
		if result <= 0 {
			println("if k'=", k, " probability=", dv.String(), "/10000. found it! this is the minimum k'.",  " sum=", sum.String())
			break
		}else{
			println("if k'=", k, " probability=", dv.String(), "/10000. sum=", sum.String())
		}
	}

}
