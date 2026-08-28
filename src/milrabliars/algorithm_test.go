package milrabliars_test

import (
	"fmt"
	"milrabliars/milrabliars"
	"sort"
)

func ExampleDressUp() {
	dressedUp := milrabliars.DressUp(747)
	fmt.Printf("747 = 2^%d * %d + 1", dressedUp.Pow2, dressedUp.Simplified)
	// Output:
	// 747 = 2^1 * 373 + 1
}

func ExampleDressUp_pow2Plus1() {
	dressedUp := milrabliars.DressUp(9)
	fmt.Printf("9 = 2^%d * %d + 1", dressedUp.Pow2, dressedUp.Simplified)
	// Output:
	// 9 = 2^3 * 1 + 1
}

func ExampleCalculateAPowXModN() {
	info := milrabliars.InitializeMontgomery(7)
	fmt.Printf("3^5 ≡ %d (mod 7)", milrabliars.ConvertFromMontgomery(milrabliars.CalculateAPowXModN(milrabliars.ConvertToMontgomery(3, 7, info), 5, 7, info), 7))
	// Output:
	// 3^5 ≡ 5 (mod 7)
}

func ExampleIsStrongProbablePrimeToBaseA_pass91and10() {
	fmt.Printf("91 is a strong probable prime to base 10: %t", milrabliars.IsStrongProbablePrimeToBaseA(91, 10, milrabliars.InitializeMontgomery(91)))
	// Output:
	// 91 is a strong probable prime to base 10: true
}

func ExampleIsStrongProbablePrimeToBaseA_pass97and2() {
	fmt.Printf("97 is a strong probable prime to base 2: %t", milrabliars.IsStrongProbablePrimeToBaseA(97, 2, milrabliars.InitializeMontgomery(97)))
	// Output:
	// 97 is a strong probable prime to base 2: true
}

func ExampleIsStrongProbablePrimeToBaseA_pass747and23() {
	fmt.Printf("747 is a strong probable prime to base 23: %t", milrabliars.IsStrongProbablePrimeToBaseA(747, 23, milrabliars.InitializeMontgomery(747)))
	// Output:
	// 747 is a strong probable prime to base 23: false
}

func ExampleIsStrongProbablePrimeToBaseA_pass65and8() {
	fmt.Printf("65 is a strong probable prime to base 8: %t", milrabliars.IsStrongProbablePrimeToBaseA(65, 8, milrabliars.InitializeMontgomery(65)))
	// Output:
	// 65 is a strong probable prime to base 8: true
}

func ExampleIsStrongProbablePrimeToBaseA_pass21and8() {
	fmt.Printf("21 is a strong probable prime to base 8: %t", milrabliars.IsStrongProbablePrimeToBaseA(21, 8, milrabliars.InitializeMontgomery(21)))
	// Output:
	// 21 is a strong probable prime to base 8: false
}

func ExampleIsStrongProbablePrimeToBaseA_compositesFor8() {
	for i := 11; i < 100; i += 2 {
		info := milrabliars.InitializeMontgomery(i)
		// 2 does not lie for these numbers
		if !milrabliars.IsStrongProbablePrimeToBaseA(i, 2, info) && milrabliars.IsStrongProbablePrimeToBaseA(i, 8, info) {
			fmt.Printf("%d is a strong probable prime to base 8\n", i)
		}
	}
	// Output:
	// 65 is a strong probable prime to base 8
}

func ExampleCalculateIsPrimeAndLiars_pass9() {
	prime, arr := milrabliars.CalculateIsPrimeAndLiars(9)
	sort.Slice(arr, func(i, j int) bool { return arr[i] < arr[j] })
	fmt.Printf("%t: %v", prime, arr)
	// Output:
	// false: []
}

func ExampleCalculateIsPrimeAndLiars_pass91() {
	prime, arr := milrabliars.CalculateIsPrimeAndLiars(91)
	sort.Slice(arr, func(i, j int) bool { return arr[i] < arr[j] })
	fmt.Printf("%t: %v", prime, arr)
	// Output:
	// false: [9 10 12 16 17 22 29 38 53 62 69 74 75 79 81 82]
}

func ExampleCalculateIsPrimeAndLiars_pass97() {
	prime, arr := milrabliars.CalculateIsPrimeAndLiars(97)
	sort.Slice(arr, func(i, j int) bool { return arr[i] < arr[j] })
	fmt.Printf("%t: %v", prime, arr)
	// Output:
	// true: []
}
