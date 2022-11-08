package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image/png"
	"os"
	"strings"
)

func saveFavicon(input string, ip string, port uint16) {
	decoded, err := base64.StdEncoding.DecodeString(strings.Split(input, ",")[1])
	if err != nil {
		panic("Cannot decode b64")
	}
	
	filename := strings.Replace(ip, ".", "_", 3) + "_" + fmt.Sprint(port)
	
	r := bytes.NewReader(decoded)
	im, err := png.Decode(r)
	if err != nil {
		fmt.Errorf("Invalid png")
		return 
	}
	f, err := os.OpenFile("out/"+filename+".png", os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		fmt.Errorf("Cannot open file")
	}
	png.Encode(f, im)
	}