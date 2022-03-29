package taps

/*
import (
	"crypto/tls"
	"fmt"
)

func ExampleSecurityParameters_Set() {
	sp := NewSecurityParameters()

	var suite *tls.CipherSuite
	// find CipherSuite
	for _, suite = range tls.CipherSuites() {
		if suite.ID == tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256 {
			break
		}
	}

	// Calling Set() is possible, because the TAPS (draft)
	// specifies a Set function on the SecurityParameters Object
	// (i.e., struct).
	err := sp.Set("ciphersuite", suite)
	if err != nil {
		panic(err)
	}

	// Idiomatic Go would be to instead directly access the Field:
	sp.CipherSuite = suite

	fmt.Printf("%s", sp.CipherSuite.Name)
	// Output: TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256

}
*/
