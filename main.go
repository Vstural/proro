package main

import (
	"log"
	"net"
	"proro/httpserver"
	"proro/pkg/rtmp"
	"proro/pkg/streammgr"
	"strconv"
)

func init() {
	log.SetFlags(log.Ldate | log.Lshortfile)
}

func main() {
	rtmpstream := startRMTPServer(1935)
	// startHTTPServer(rtmpstream)

	httpserver, err := httpserver.NewHTTPServer(&httpserver.HTTPOption{
		RTMPstream:   rtmpstream,
		Lnaddr:       "0.0.0.0:8080",
		StreamingMgr: streammgr.NewStreamMgr(),
	})
	if err != nil {
		panic(err)
	}
	httpserver.Serve()
}

func startRMTPServer(port int) *rtmp.RtmpStream {
	if port < 0 {
		port = 1935 // rtmp default port
	}
	// stream := rtmp.NewRtmpStream()

	rtmpListen, err := net.Listen("tcp", "127.0.0.1:"+strconv.Itoa(port))
	if err != nil {
		log.Fatal(err)
	}

	rtmpServer := rtmp.NewRtmpServer()

	defer func() {
		if r := recover(); r != nil {
			log.Println("RTMP server panic: ", r)
		}
	}()

	go rtmpServer.Serve(rtmpListen)
	return rtmpServer.GetStream()
}
