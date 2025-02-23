package main

import "fmt"
import "math/rand"

func abs(x int) int {
	if x >= 0 { return x; }
	return -x;
}

func encode(utf32 []rune) []byte {
	utf8 := make([]byte, 0);
	var numberOfOctets int;
	for _, codePoint := range utf32 {
		if codePoint <= 0x7F {
			numberOfOctets = 1;
		} else if codePoint <= 0x7FF {
			numberOfOctets = 2;
		} else if codePoint <= 0xFFFF {
			numberOfOctets = 3;
		} else {
			numberOfOctets = 4;
		}
		if numberOfOctets == 1 {
			utf8 = append(utf8, (byte)(codePoint));
		} else {
			utf8 = append(utf8, 0xF0 << abs(numberOfOctets - 4) |
					    (byte)(codePoint >> (6 * (numberOfOctets - 1))));
			for j := numberOfOctets - 1; j > 0; j-- {
				utf8 = append(utf8, 0x80 | (byte)(codePoint >> (6 * (j - 1)) & 0x3F));
			}
		}
	}
	return utf8;
}

func decode(utf8 []byte) []rune {
	var utf32 []rune = make([]rune, 0);
	var codePoint rune = 0;
	var numberOfOctets int;
	for i := 0; i < len(utf8); i++ {
		if utf8[i] & 0x80 == 0 {
			numberOfOctets = 1;
		} else if utf8[i] & 0x20 == 0 {
			numberOfOctets = 2;
		} else if utf8[i] & 0x10 == 0 {
			numberOfOctets = 3;
		} else {
			numberOfOctets = 4;
		}
		if numberOfOctets == 1 {
			utf32 = append(utf32, (rune)(utf8[i]));
		} else {
			codePoint = (rune)(utf8[i] & ((1 << (7 - numberOfOctets)) - 1)) <<
				    (6 * (numberOfOctets - 1));
			for j := 1; j < numberOfOctets; j++ {
				codePoint = codePoint |
				((rune)(utf8[i + j] & 0x3F)) << (6 * (numberOfOctets - 1 - j));
			}
			i = i + numberOfOctets - 1;
			utf32 = append(utf32, codePoint);
		}
	}
	return utf32;
}

func randSeq(n int) string {
    b := make([]rune, n)
    for i := 0; i < n; i++ {
	    b[i] = rand.Int31n(1 << 20);
    }
    return string(b)
}

func main() {
	var s string = randSeq(99);
	fmt.Printf("%s\n%s\n", s, (string)(encode(([]rune)(s))));
	fmt.Println(s == (string)(encode(([]rune)(s))));
	fmt.Printf("%d, %d\n\n", len(s), len((string)(encode(([]rune)(s)))));
	fmt.Printf("%s\n%s\n", s, (string)(decode(([]byte)(s))));
	fmt.Println(s == (string)(decode(([]byte)(s))));
	fmt.Printf("%d, %d\n", len(s), len((string)(decode(([]byte)(s)))));
}

