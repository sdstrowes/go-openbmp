package obmp

import (
	"fmt"
	"encoding/hex"
	"testing"
)

func TestParse(t *testing.T) {
	input, err := hex.DecodeString("4f424d500107006400000033800c653f029d0002a952fef880938d19e9d632c815d1e95a87e1000a69732d61682d626d7031de6006b9151b7fa194e5c5aa87f2896480df3366000000000000000000000000000c726f7574652d76696577733200000001030000003302000000000000000000000000000000000000000000005be497010000792b000000005ec6f61b000389b9020000")
	if err != nil {
		t.Error("Failed decoding hex string (??)")
	}
	fmt.Println("IN: ", input)
	out, err := ParseHeader(input)
	fmt.Printf("%+v\n", out)

	if err != nil {
		t.Errorf("Bad parse from example message header: "+err.Error())
	}
}
