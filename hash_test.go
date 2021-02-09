package glowlib

import(
		"log"
		"math/rand"
		"testing"
		)

func InternalHashTest( t *testing.T, h1 int64, h2 int64 ) {
	hash_src 	 := Hash128{ h1, h2 }
	hash_str :=  hash_src.String()

	if ( len( hash_str ) != 32 ) {
		t.Error( "Hash %s was not 32 characters long", hash_str )
	}

	hash := Hash128FromHex( hash_str )
	if ( hash.H1 != h1 || hash.H2 != h2 ) {
		t.Error( "Hash128FromString conversion failed")
	}
}

func CheckHashIsValid( t *testing.T, hashname string, expected bool ) {
	if expected != HashIsValid( hashname ) {
		t.Error( "HashIsValid %s", hashname)
	}
}

func TestHash( t *testing.T) {
	InternalHashTest( t, 123419281928712, 98172938701286 )
	InternalHashTest( t, 0, 0 )
	InternalHashTest( t, 1, 1 )

	CheckHashIsValid( t, "8dfc46a60e3e996a26bd1107c01ff5ed", true )
	CheckHashIsValid( t, "00112233445566778899aabbccddeefg", false )
	CheckHashIsValid( t, "0", false )
	CheckHashIsValid( t, "00112233445566778899aabbccddeef", false )

	r := rand.New(rand.NewSource(99))
	for i := 0; i < 100; i++ {
		var hash Hash128 
		hash.H1 = r.Int63()
		hash.H2 = r.Int63()

		//enc := encodeUint64ToString( input )
		//log.Printf( "Input = %d Encode = %s Decode = %d\n", input, enc, decodeStringToUint64( enc ) )
		hres := HashDecodeFromUrl( HashEncodeToUrl( hash ) )
		if hres.H1 != hash.H1 || hres.H2 != hash.H2 {
			log.Printf( "Bad hash %v", hash)
		}
	}
}

