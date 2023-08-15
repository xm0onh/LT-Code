package Cryptography

import (
	"go.dedis.ch/kyber/pairing/bn256"
	"go.dedis.ch/kyber/sign/bls"
	"go.dedis.ch/kyber/util/random"

	"go.dedis.ch/kyber"
	//	"os/exec"

	"os"
	//"go.dedis.ch/kyber"
	//	"os/exec"
	"fmt"
	"strconv"
)

func Load_Own_keys(priv string, pub string) (kyber.Scalar, kyber.Point) {
	privatekeyfile, err := os.Open(priv)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Priv Key name is", priv)
	//	decoder := gob.NewDecoder(privatekeyfile)
	privatekey, publickey := bls.NewKeyPair(bn256.NewSuite(), random.New())
	fmt.Print("Privatekey before assignment", privatekey)
	//err = decoder.Decode(&privatekey)

	_, err1 := privatekey.UnmarshalFrom(privatekeyfile)
	if err1 != nil {
		fmt.Println("reading priv key fails", err.Error())
		os.Exit(1)
	}

	privatekeyfile.Close()
	fmt.Print("Privatekey after assignment", privatekey)

	publickeyfile, err := os.Open(pub)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//	decoder = gob.NewDecoder(publickeyfile)

	_, err2 := publickey.UnmarshalFrom(publickeyfile)
	//err = decoder.Decode(&publickey)

	if err2 != nil {
		fmt.Println("reading pub key fails", err2.Error())
		os.Exit(1)
	}

	publickeyfile.Close()

	return privatekey, publickey

}
func Load_PubKeys(NodeSliceLen int) []kyber.Point {
	PubArr := make([]kyber.Point, 0, NodeSliceLen)
	for i := 0; i < NodeSliceLen; i++ {
		//Pubkey := fmt.Sprintf("Pub%s", ArrInstanceID[i])
		//	var    ToBreadedPublickKey kyber.Point
		_, publickey := bls.NewKeyPair(bn256.NewSuite(), random.New())
		Pubkey := fmt.Sprintf("Pub%s", strconv.Itoa(i))
		fmt.Println("Pubkey filename is", Pubkey)
		publickeyfile, err := os.Open(Pubkey)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		//	decoder := gob.NewDecoder(publickeyfile)

		_, err1 := publickey.UnmarshalFrom(publickeyfile)
		//	err = decoder.Decode(&publickey)

		if err1 != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		publickeyfile.Close()
		PubArr = append(PubArr, publickey)

	}
	return PubArr

}

func KeySetup(NodeSliceLen int) {
	//suite := bn256.NewSuite()
	fmt.Println("Setting up keys")
	for i := 0; i < NodeSliceLen+150; i++ {
		//bs := []byte(strconv.Itoa(i+1000000))
		//	reader := rand.New(rand.NewSource(int64(i)))
		//      reader := rand.Reader
		//reader:=i
		//	Pubkey := fmt.Sprintf("Pub%s", strArr[i])
		//	Privkey := fmt.Sprintf("Priv%s", strArr[i])
		Pubkey := fmt.Sprintf("Pub%s", strconv.Itoa(i))
		Privkey := fmt.Sprintf("Priv%s", strconv.Itoa(i))
		key, publicKey := bls.NewKeyPair(bn256.NewSuite(), random.New())
		outFilePrv, err := os.Create(Privkey)
		checkError(err)
		defer outFilePrv.Close()
		outFilePub, err := os.Create(Pubkey)
		checkError(err)
		defer outFilePub.Close()
		key.MarshalTo(outFilePrv)
		publicKey.MarshalTo(outFilePub)
		fmt.Println("Private and pubkey  are", key, publicKey)

	}

}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}
