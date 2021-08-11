package main

import (
	"bufio"
	"github.com/winlinvip/go-aresample/aresample"
	"github.com/winlinvip/go-fdkaac/fdkaac"
	"io"
	"log"
	"os"
)

var (
	aacData *os.File
	part    []byte
	err     error
	count   int
	pcm     []byte
)

func main() {
	aacFile := "./sample.aac"
	pcmFile := "./sample.pcm"

	aacData, err = os.Open(aacFile)
	if err != nil {
		log.Fatal(err)
	}
	defer aacData.Close()
	log.Println("Open aac file ", aacFile)

	pcmData, err := os.Create(pcmFile)
	if err != nil {
		log.Fatal(err)
	}
	defer pcmData.Close()
	log.Println("Create pcm file ", aacFile)

	d := fdkaac.NewAacDecoder()

	if err := d.InitAdts(); err != nil {
		log.Fatal("init decoder failed, err is", err)
		return
	}
	defer d.Close()

	var r aresample.ResampleSampleRate
	if r,err = aresample.NewPcmS16leResampler(2, 44100, 8000); err != nil {
		log.Println("aresample failed, err is", err)
		return
	}

	reader := bufio.NewReader(aacData)
	part := make([]byte, 128)

	var npcm []byte

	for {
		if count, err = reader.Read(part); err != nil {
			break
		}

		if pcm, err = d.Decode(part[:count]); err != nil {
			log.Fatal("decode failed, err is", err)
			return
		}

		if len(pcm) == 0 {
			continue
		}

		if npcm,err = r.Resample(pcm); err != nil {
			log.Println("aresample failed, err is", err)
			return
		}
		log.Println("44.1 KHZ PCM:", len(pcm))
		log.Println("8 KHZ NPCM:", len(npcm))

		pcmData.Write(npcm)
	}
	if err != io.EOF {
		log.Fatal("reading failed file: ", aacFile, " err is ", err)
	} else {
		err = nil
	}
}
