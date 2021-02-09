package glowlib

import (
		"math/big"
		"math/rand"
		"strings"
)

// create an alphabetical uri from a hash that is shorter than a hex conversion
var alphabet 		= ShuffleString( "abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ-_" )
var alphabetBase 	= uint64( len( alphabet ) )
var alphabetBaseBig	= big.NewInt( int64( alphabetBase ) )
var bigZero 		= big.NewInt( 0 )

//http://math.stackexchange.com/questions/9508/looking-for-a-bijective-discrete-function-that-behaves-as-chaotically-as-possib
func ShuffleString( input string ) string {
	perm := rand.Perm( len( input ) )
	var output string
	for _, v := range perm {
		output += string(input[v])
	}

	return output
}

var linkMulPrime1 = new( big.Int ).SetUint64( 687724237 )
var linkMulPrime2 = new( big.Int ).SetUint64( 890234203 )
var linkMulAdd 	 = new( big.Int ).SetUint64( 6066941 )
var linkMulXor 	 = new( big.Int ).SetUint64( 9223372036854775837 )

var linkBigMod 	 = new( big.Int ).Lsh( new( big.Int ).SetUint64( 1 << 63 ), 64 )

// some kind of obfuscated link encoding. not super great!
func LinkIdEncode( publish_id uint64 ) string {
	var id big.Int
	// multiply, mod 1<< 63
	id.SetUint64( publish_id )
	id2 := new(big.Int)
	id2.Mul( &id, linkMulPrime1 )
	id2.Add( id2, linkMulAdd )
	id2.Lsh( id2, 7 )
	id2.Mul( id2, linkMulPrime2 )
	//id2.Xor( id2, linkMulXor )

	link_str := ""
	for id2.Cmp( bigZero ) > 0 {
		link_str = string( alphabet[ id2.Uint64() % alphabetBase ] ) + link_str
		id2.Rsh( id2, 6 )
	}

	return link_str
}

func LinkIdDecode( link_id string ) uint64 {

	id := big.NewInt( 0 )
	idx := big.NewInt( 0 )

	for i := 0; i < len( link_id ); i++ {
		id.Mul( id, alphabetBaseBig )
		idx.SetUint64( uint64( strings.IndexByte( alphabet, link_id[ i ] ) ) )
		id.Add( id, idx )
	}

	id2 := new(big.Int)
	//id2.Xor( id, linkMulXor )
	id2.Div( id, linkMulPrime2 )
	id2.Rsh( id2, 7 )
	id2.Sub( id2, linkMulAdd )
	id2.Div( id2, linkMulPrime1 )

	return id2.Uint64()
}


// i guess we;re base64 encoding it?
func EncodeUint64ToString( t uint64 ) string {
	if (t == 0) {
		return string( alphabet[0] )
	} 

	str := ""

	for  t > 0 {
		str = string( alphabet[ t % alphabetBase ] ) + str
		t = t / alphabetBase
	}
	return str
}

func DecodeStringToUint64( str string ) uint64 {
	var t uint64 = 0
	for i := 0; i < len( str ); i++ {
		t = ( t * alphabetBase ) + uint64( strings.IndexByte( alphabet, str[ i ] ) )
	}

	return t
}

func HashEncodeToUrl( t Hash128 ) string {
	str := ""
	t1 := uint64(t.H1)
	t2 := uint64(t.H2)

	for i := 0; i < 10; i++ {
		str = string ( alphabet[ t1 & 63 ] ) + str 
		t1 = t1 >> 6	
	}

	tt := t1 & 15 | (t2 & 3) << 4
	t2 = t2 >> 2
	str = string ( alphabet[ tt & 63 ] ) + str 

	for t2 > 0 {
		str = string ( alphabet[ t2 & 63 ] ) + str 
		t2 = t2 >> 6	
	}

	return str
}

func HashDecodeFromUrl( str string ) Hash128 {
	t := big.NewInt( 0 )

	idx := big.NewInt( 0 )
	for i := 0; i < len( str ); i++ {
		t.Mul( t, alphabetBaseBig )
		idx.SetUint64( uint64( strings.IndexByte( alphabet, str[ i ] ) ) )
		t.Add( t, idx )
	}

	var res Hash128
	res.H1 = int64( t.Uint64() )
	t.Rsh( t, 64 )
	res.H2 = int64( t.Uint64() )
	return res
}