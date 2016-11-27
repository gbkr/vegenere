package vegenerelib

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"sort"
)

const (
	NormalizingCoefficient    = 26
	MaxKeyLength              = 10
	IndexOfCoincidenceEnglish = 1.73
)

var alphabet = setAlphabet()
var key_position int = 0

func Encrypt(source string, key string) string {
	if len(key) == 0 {
		exitWithMessage("No key specified")
	}

	content := openFile(source)
	msg := processFile(content, key, "encrypt")

	return msg
}

func Decrypt(source string, key string) string {
	content := openFile(source)

	// If no key is given, attemp to determine key
	if len(key) == 0 {
		keySize := calcKeyLength(content)
		key = identifyKey(content, keySize)
	}

	result := processFile(content, key, "decrypt")
	return result
}

func DecryptKey(source string) string {
	content := openFile(source)
	keySize := calcKeyLength(content)
	key := identifyKey(content, keySize)
	return key
}

func chiSquaredStatistic(content []byte) float64 {
	var chiSquared float64 = 0

	for i := 0; i < len(alphabet); i++ {
		letter := alphabet[i]
		actualLetterIncidence := float64(letterCount(letter, content))
		if actualLetterIncidence > 0 {
			expectedLetterIncidence := (letterFrequency(letter) / float64(100)) * float64(len(content))
			chiSquared += math.Pow((expectedLetterIncidence-actualLetterIncidence), 2) / expectedLetterIncidence
		}
	}
	return chiSquared
}

// Find the shift value for a Caesar cipher
func findKeyChar(caesar []byte) byte {
	results := make(map[int]float64)
	for step := 0; step < len(alphabet); step++ {
		testCipher := []byte{}
		for _, value := range caesar {
			v := upperByte(value) - byte(step)
			if v < alphabetStart(v) {
				v += 26
			}
			testCipher = append(testCipher, v)
		}

		results[step] = chiSquaredStatistic(testCipher)
	}

	// find the key of the lowest value
	values := make([]float64, len(results))
	for i, v := range results {
		values[i] = v
	}
	sort.Float64s(values)

	var key byte
	for i, v := range results {
		if v == values[0] {
			key = byte(i)
		}
	}

	return key
}

func identifyKey(content []byte, keySize int) string {
	key := []byte{}
	sanitizedStr := sanitizeContent(content)

	for offset := 0; offset < keySize; offset++ {
		caesar := []byte{}
		for i := offset; i < len(sanitizedStr); i += keySize {
			caesar = append(caesar, sanitizedStr[i])
		}
		key = append(key, findKeyChar(caesar))
	}

	for i, v := range key {
		key[i] = alphabet[v]
	}
	return string(key)
}

func letterCount(letter byte, content []byte) int {
	count := 0
	for _, v := range content {
		if upperByte(v) == letter {
			count++
		}
	}
	return count
}

func encryptableContentLength(content []byte) int {
	count := 0
	for _, letter := range content {
		if isEncryptable(letter) {
			count++
		}
	}
	return count
}

func probabilityOfDrawingLetterTwice(letter byte, content []byte) float64 {
	if !isEncryptable(letter) {
		return 0
	}
	count := float64(letterCount(letter, content))
	contentLength := float64(encryptableContentLength(content))
	return (count / contentLength) * ((count - 1) / (contentLength - 1))
}

// Probability of drawing two matching letters by randomly selecting two letters from given text
// https://en.wikipedia.org/wiki/Index_of_coincidence
func indexOfCoincidence(content []byte) float64 {
	var probabilitySum float64 = 0
	for _, letter := range alphabet {
		probability := probabilityOfDrawingLetterTwice(letter, content)
		probabilitySum += probability
	}

	return NormalizingCoefficient * probabilitySum
}

func ciForKeyLength(keyLength int, offset int, content []byte) float64 {
	column := []byte{}
	for i := offset; i < len(content); i += keyLength {
		column = append(column, content[i])
	}
	return indexOfCoincidence(column)
}

func keyLengthThatProducesEnglishLetterFreq(frequencies []float64) int {
	freqOffset := make(map[float64]int)

	for i, v := range frequencies {
		freqOffset[math.Abs(v-IndexOfCoincidenceEnglish)] = i
	}

	var offsets []float64
	for v := range freqOffset {
		offsets = append(offsets, v)
	}
	sort.Float64s(offsets)

	keyLengthIndex := offsets[0]
	return freqOffset[keyLengthIndex] + 1
}

func calcKeyLength(content []byte) int {
	sanitizedStr := sanitizeContent(content)
	icsForColumns := make([]float64, MaxKeyLength)

	for keyLength := 1; keyLength <= MaxKeyLength; keyLength++ {
		var sum float64 = 0
		for offset := 0; offset < keyLength; offset++ {
			sum += ciForKeyLength(keyLength, offset, sanitizedStr)
		}

		icsForColumns[keyLength-1] = (sum / float64(keyLength))
	}

	keyLength := keyLengthThatProducesEnglishLetterFreq(icsForColumns)
	return keyLength
}

func sanitizeContent(content []byte) []byte {
	sanitized := []byte{}
	for _, value := range content {
		if isEncryptable(value) {
			sanitized = append(sanitized, value)
		}
	}
	return sanitized
}

func setAlphabet() []byte {
	alphabet := []byte{}
	var i byte
	for i = 65; i <= 90; i++ {
		alphabet = append(alphabet, i)
	}
	return alphabet
}

func processFile(content []byte, key string, action string) string {
	output := make([]byte, len(content))
	for i, n := range content {
		output[i] = processByte(n, key, action)
	}
	return string(output)
}

func upperByte(b byte) byte {
	asciiStart := alphabetStart(b)
	if asciiStart == 97 {
		return b - 32
	} else {
		return b
	}
}

func charIndex(b byte) byte {
	return byte(bytes.IndexByte(alphabet, upperByte(b)))
}

func processByte(b byte, key string, action string) byte {
	if !isEncryptable(b) {
		return b
	}
	steps := charIndex(key[key_position])
	var transformed_byte byte
	switch action {
	case "decrypt":
		transformed_byte = decrypt(b, steps)
	case "encrypt":
		transformed_byte = encrypt(b, steps)
	}

	increment_key_position(key)
	return transformed_byte
}

func decrypt(b byte, steps byte) byte {
	enc_byte := b - steps
	if enc_byte < alphabetStart(b) {
		enc_byte += 26
	}
	return enc_byte
}

func encrypt(b byte, steps byte) byte {
	asciiStart := alphabetStart(b)
	i := byte(math.Mod((float64(b - asciiStart + steps)), 26))
	return i + asciiStart

}

// Returns the first letter of the alphabet as a byte in the same case as the input
func alphabetStart(b byte) byte {
	if b >= 97 {
		return 97
	} else {
		return 65
	}
}

func isEncryptable(b byte) bool {
	return bytes.IndexByte(alphabet, upperByte(b)) != -1
}

func increment_key_position(key string) {
	if key_position == (len(key) - 1) {
		key_position = -1
	}
	key_position++
}

func exitWithMessage(message string) {
	fmt.Println(message)
	fmt.Println("Run 'vegenere --help' for usage.")
	os.Exit(0)
}

func checkErr(err error, message string) {
	if err != nil {
		exitWithMessage(message)
	}
}

func openFile(fromFile string) []byte {
	if len(fromFile) == 0 {
		exitWithMessage("No source file specified")
	}

	content, err := ioutil.ReadFile(fromFile)
	checkErr(err, "Invalid source file")
	content = []byte(string(content))

	return content
}

func letterFrequency(char byte) float64 {
	frequencies := map[string]float64{
		"A": 8.167,
		"B": 1.492,
		"C": 2.782,
		"D": 4.253,
		"E": 12.702,
		"F": 2.228,
		"G": 2.015,
		"H": 6.094,
		"I": 6.966,
		"J": 0.153,
		"K": 0.772,
		"L": 4.025,
		"M": 2.406,
		"N": 6.749,
		"O": 7.507,
		"P": 1.929,
		"Q": 0.095,
		"R": 5.987,
		"S": 6.327,
		"T": 9.056,
		"U": 2.758,
		"V": 0.978,
		"W": 2.360,
		"X": 0.150,
		"Y": 1.974,
		"Z": 0.074,
	}
	return frequencies[string(char)]
}
