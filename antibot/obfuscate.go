package antibot

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"math/rand"
	"regexp"
)

// randomString produces a random alphabetic string of the given length.
// These random strings will be used as variable names to hide meaning.
func randomString(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// ObfuscateJS takes a JavaScript code string and returns an obfuscated version,
// wrapped in a <script> tag. The obfuscated code is heavily hidden using a two‑step process:
//   - The JS code is minified by removing extra whitespace
//   - It is then XOR‑encrypted with a random key and Base64‑encoded.
//
// At runtime, the generated script decodes and decrypts the code before passing it to eval.
func ObfuscateJS(jsCode string) string {
	// Step 1: Remove extra whitespace (a simple minification step).
	re := regexp.MustCompile(`\s+`)
	compactJS := re.ReplaceAllString(jsCode, " ")

	// Step 2: Generate a random key for XOR encryption.
	xorKey := randomString(10)

	// Step 3: Perform XOR encryption with the random key.
	var encrypted bytes.Buffer
	codeBytes := []byte(compactJS)
	for i, b := range codeBytes {
		encrypted.WriteByte(b ^ xorKey[i%len(xorKey)])
	}
	// Base64 encode the encrypted data.
	encodedData := base64.StdEncoding.EncodeToString(encrypted.Bytes())

	// Step 4: Create randomized variable and function names for the output snippet.
	keyVar := randomString(8)         // Holds the XOR key.
	dataVar := randomString(8)        // Holds the Base64 encoded encrypted data.
	decryptFuncVar := randomString(8) // The decryption function.

	// Step 5: Build the self-decoding JavaScript snippet.
	// Note: The double '%%' escapes a percent sign so that "% k.length" shows correctly in the output.
	obfuscatedJS := fmt.Sprintf(
		`<script>(function(){
    var %s = "%s";
    var %s = "%s";
    function %s(d, k){
        var r = "";
        var s = atob(d);
        for(var i = 0; i < s.length; i++){
            r += String.fromCharCode(s.charCodeAt(i) ^ k.charCodeAt(i %% k.length));
        }
        return r;
    }
    eval(%s(%s, %s));
})();</script>`,
		keyVar, xorKey,
		dataVar, encodedData,
		decryptFuncVar,
		decryptFuncVar, dataVar, keyVar)

	return obfuscatedJS
}
