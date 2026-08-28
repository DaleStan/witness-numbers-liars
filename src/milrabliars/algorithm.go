package milrabliars

const MUST_BE_64_BIT = int(1 << 48)
const MONTGOMERY_SHIFT = 32
const MONTGOMERY_R = 1 << MONTGOMERY_SHIFT
const MONTGOMERY_MASK = MONTGOMERY_R - 1

// DressedUpInt is an int represented as 2^(Pow2) * Simplified + 1.
type DressedUpInt struct {
	Pow2       int
	Simplified int
}

// DressUp the given number by representing it as 2^(Pow2) * Simplified + 1.
func DressUp(n int) DressedUpInt {
	var pow2 int = 1
	var simplified int = (n - 1) / 2

	for simplified%2 == 0 {
		pow2++
		simplified /= 2
	}

	return DressedUpInt{
		Pow2:       pow2,
		Simplified: simplified,
	}
}

type MontgomeryInfo struct {
	// The magic constant k = (r * r⁻¹ - 1) / n, used in the reduction
	k           int
	// The Montgomery forms of 1 and -1, so comparisons can be done without an
	// expensive convert-out operation.
	One         int
	MinusOne    int
}

// Initializes a MontgomeryInfo that can be used for operations modulo n
func InitializeMontgomery(n int) MontgomeryInfo {
	inverse := ModularInverse(MONTGOMERY_R, n)
	return MontgomeryInfo{
		k:        MONTGOMERY_R * inverse / n,
		One:      ConvertToMontgomery(1, n),
		MinusOne: ConvertToMontgomery(n - 1, n),
	}
}

// Computes a⁻¹ such that a * a⁻¹ ≡ 1 (mod n)
func ModularInverse(a, n int) int {
	var t, newt int = 0, 1
	r, newr := n, a

	for newr != 0 {
		quotient := r / newr
		t, newt = newt, t - quotient * newt
		r, newr = newr, r - quotient * newr
	}

	if (t < 0) {
		t += n;
	}
	return t;
}

// Converts the input number into Montgomery form, for computation modulo n.
func ConvertToMontgomery(x, n int) int {
	return (x << MONTGOMERY_SHIFT) % n
}

// For testing purposes only, convert a number back from Montgomery form.
// Production code compares Montgomery products to the pre-computed Montgomery
// forms of 1 and -1.
func ConvertFromMontgomery(x, n int) int {
	inverse := ModularInverse(MONTGOMERY_R, n)
	return x * inverse % n
}

// CalculateAPowXModN calculates a^x (mod n), provided a < n and x is odd.
// a and the output are in Montgomery form, x and n are in normal form, and
// the precomputed Montgomery values for n are in info
func CalculateAPowXModN(a, x, n int, info *MontgomeryInfo) int {
	if x == 1 {
		return a // a < n
	}

	// The exponent is always odd, so handle its low-order bit unconditionally
	// (and skip a 1 * <something> multiply)
	result := a
	x >>= 1
	squared := MontgomeryMultiply(a, a, n, info)

	for x > 1 {
		if x&1 != 0 {
			result = MontgomeryMultiply(result, squared, n, info)
		}
		x >>= 1
		squared = MontgomeryMultiply(squared, squared, n, info)
	}

	// x (exponent) is always 1 here
	result = MontgomeryMultiply(result, squared, n, info)
	return result
}

// Compute the Montgomery form of a*b%n, given the Montgomery forms of a and
// b, the normal form of n, and the precomputed values in info
func MontgomeryMultiply(a, b, n int, info *MontgomeryInfo) int {
	x := a * b
	// `x * info.k` may overflow if n > 2^21, but any overflowed bits won't
	// affect the value of s
	s := (x * info.k) & MONTGOMERY_MASK
	t := x + s * n
	u := t >> MONTGOMERY_SHIFT
	if u >= n { u -= n }
	return u
}

// IsStrongProbablePrimeToBaseADirect checks if the given number `n` is a strong
// probably prime to base `a`, i.e., if the witness a testifies that the number
// n is prime. This must be given `s` and `d` such that `n = 2^s * d + 1`.
//
// A number n is a strong probably prime to a iff any of the following hold:
// - a^d ≡ 1 (mod n)
// - there exists some r integer where 0 ≤ r < s and a^(2^r * d) ≡ -1 (mod n)
func IsStrongProbablePrimeToBaseADirect(n, a, s, d int, info *MontgomeryInfo) bool {
	res := CalculateAPowXModN(ConvertToMontgomery(a, n), d, n, info)
	if res == info.One || res == info.MinusOne {
		return true
	}

	for r := 1; r < s; r++ {
		res = MontgomeryMultiply(res, res, n, info)
		if res == info.MinusOne {
			return true
		}
	}

	return false
}

// IsStrongProbablePrimeToBaseA checks if the given number is a strong probable
// prime to base A, ie., if the witness a testifies that the number n is prime.
// If this returns true, then n is probably prime. Otherwise, n is composite.
func IsStrongProbablePrimeToBaseA(n, a int, info *MontgomeryInfo) bool {
	dressedUp := DressUp(n)
	return IsStrongProbablePrimeToBaseADirect(n, a, dressedUp.Pow2, dressedUp.Simplified, info)
}

var StarWitnessesInt32 = [...]int{2, 3, 5, 7}

// CalculateIsPrimeAndLiars both determines if n is prime deterministically and
// any lying witnesses. n must be odd and >= 3. If the number is prime then it's
// not possible for a witness to lie and so the slice is nil rather than empty.
func CalculateIsPrimeAndLiars(n int) (bool, []int) {
	info := InitializeMontgomery(n)

	var firstTruthfulStarWitnessIndex int = -1
	for idx, witness := range StarWitnessesInt32 {
		if !IsStrongProbablePrimeToBaseA(n, witness, &info) {
			firstTruthfulStarWitnessIndex = idx
			break
		}

		if n < 2047 {
			// checking just 2 is enough for n < 2047
			return true, nil
		}

		if idx == 1 && n < 1373653 {
			return true, nil
		}

		if idx == 2 && n < 25326001 {
			return true, nil
		}
	}

	if firstTruthfulStarWitnessIndex == -1 {
		return true, nil
	}

	liars := make([]int, 0, 64)
	var idxHint int = 0
	for a := 2; a < StarWitnessesInt32[firstTruthfulStarWitnessIndex]; a++ {
		if a == StarWitnessesInt32[idxHint] {
			// already checked
			liars = append(liars, a)
			liars = append(liars, n-a)
			idxHint++
		} else if IsStrongProbablePrimeToBaseA(n, a, &info) {
			liars = append(liars, a)
			liars = append(liars, n-a)
		}
	}

	halfN := n / 2
	for a := StarWitnessesInt32[firstTruthfulStarWitnessIndex] + 1; a <= halfN; a++ {
		if IsStrongProbablePrimeToBaseA(n, a, &info) {
			liars = append(liars, a)
			liars = append(liars, n-a)
		}
	}

	return false, liars
}
