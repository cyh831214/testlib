package glowlib

import ( 
        "crypto/aes"
        "crypto/cipher"
        "crypto/rand"
        "errors"
        "io"
        "log"
        "os"
        )

var AESLogger = log.New( os.Stdout, "[aes] ", log.LstdFlags )

func AESDecrypt( ciphertext []byte, secret string  ) ([]byte, error) {
 
 	block, err := aes.NewCipher( []byte(secret) )
 	if err != nil {
    	return nil, errors.New( "aes decrypt setup" )
 		
 	}

 	if len(ciphertext) < aes.BlockSize {
    	return nil, errors.New( "aes decrypt bad length" )
    }
  
	iv := ciphertext[:aes.BlockSize]
 	ciphertext = ciphertext[aes.BlockSize:]

 	// CBC mode always works in whole blocks.
 	if len(ciphertext)%aes.BlockSize != 0 {
    	return nil, errors.New( "aes decrypt bad length" ) 		
 	}
 	
 	mode := cipher.NewCBCDecrypter(block, iv)
 
 	// CryptBlocks can work in-place if the two arguments are the same.
 	mode.CryptBlocks(ciphertext, ciphertext )
 	
	// padding indicated by the last X bytes containing the number X
	padding_size := ciphertext[ len(ciphertext ) - 1 ]
	ciphertext = ciphertext[:len(ciphertext )- int(padding_size)]	
 	
 	return ciphertext, nil
}

func AESEncrypt( plaintext []byte, secret string ) ([]byte, error) {
	// pad the block with the length of the padding
	plaintext_len := len( plaintext )
	len_rem 	:= plaintext_len%aes.BlockSize

	// default pad an entire block
	pad_length 	:= aes.BlockSize
	if len_rem != 0 {
		pad_length = (aes.BlockSize-len_rem)
	}

	// create the padded version , pad with length of padding according to RFC????
	padded_plaintext := make( []byte, len( plaintext ) + pad_length )
    copy( padded_plaintext, plaintext )

    for i := 0; i < pad_length; i++ {
    	padded_plaintext[ len( padded_plaintext ) - pad_length + i ] = byte(pad_length)
	}

    plaintext = padded_plaintext
    if ( len( plaintext) % aes.BlockSize != 0 ) {
    	return nil, errors.New( "AES encrypt bad plaintext length")
    }

    block, err := aes.NewCipher( []byte( secret ) )
    if err != nil {
    	return nil, err
    }
    
    // The IV needs to be unique, but not secure. Therefore it's common to
    // include it at the beginning of the ciphertext.
    ciphertext := make([]byte, aes.BlockSize+len(plaintext) )
    iv := ciphertext[:aes.BlockSize]
    if _, err := io.ReadFull(rand.Reader, iv); err != nil {
    	return nil, err
    }
    
    mode := cipher.NewCBCEncrypter(block, iv)
    mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)

    return ciphertext, nil
}