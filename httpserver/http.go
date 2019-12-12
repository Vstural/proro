package httpserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"proro/pkg/rtmp"
	"proro/pkg/streammgr"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"
)

type stream struct {
	Key string `json:"key"`
	Id  string `json:"id"`
}
type streamsInfo struct {
	Publishers []stream `json:"publishers"`
	Players    []stream `json:"players"`
}
type HTTPServer struct {
	rtmpstream *rtmp.RtmpStream
	ln         net.Listener
	e          *gin.Engine
	lnaddr     string

	streamingMgr *streammgr.Mgr
}

type HTTPOption struct {
	RTMPstream   *rtmp.RtmpStream
	Lnaddr       string
	StreamingMgr *streammgr.Mgr
}

func NewHTTPServer(option *HTTPOption) (*HTTPServer, error) {
	if option == nil {
		return nil, errors.New("nil option get")
	}

	if option.RTMPstream == nil {
		return nil, errors.New("rtmp streams should not be bil")
	}
	if option.StreamingMgr == nil {
		return nil, errors.New("streaming mgr should not be nil")
	}

	e := gin.Default()

	server := &HTTPServer{
		rtmpstream:   option.RTMPstream,
		e:            e,
		lnaddr:       option.Lnaddr,
		streamingMgr: option.StreamingMgr,
	}
	server.initRouters()
	return server, nil
}

func (s *HTTPServer) initRouters() {
	const prefix = "./web/src"

	// static files
	s.e.LoadHTMLGlob("htmls/*")
	s.e.Static("/js", prefix+"/js")
	s.e.Static("/css", prefix+"/css")
	s.e.Static("/img", prefix+"/img")

	// html pages
	s.e.GET("/index", s.index)

	// apis
	s.e.GET("/apis/hello", s.hello)
	s.e.Any("/streams", s.streamsInfo)
	s.e.Any("/push/*rtmpurl", s.push)
}

func (s *HTTPServer) index(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "index.html", gin.H{})
}

func (s *HTTPServer) hello(ctx *gin.Context) {
	ctx.Writer.WriteString("Hello_proro")
}

func (s *HTTPServer) push(ctx *gin.Context) {
	// note: do something here
	defer func() {
		if r := recover(); r != nil {
			debug.PrintStack()
			log.Println("http flv handleConn panic: ", r)
		}
	}()

	// url := ctx.Request.URL.String()
	u := ctx.Request.URL.Path
	path := strings.TrimSuffix(strings.TrimLeft(u, "push/"), ".flv")

	// finalPath, err := s.pathMapping(path)

	// if err != nil {
	// 	ctx.Writer.WriteString(err.Error())
	// 	return
	// }

	// paths := strings.SplitN(path, "/", 2)
	// paths := strings.SplitN(finalPath, "/", 2)
	// log.Println("url:", u, "path:", path, "paths:", paths)

	// todo: fix this
	// ctx.Header("Access-Control-Allow-Origin", "*")

	// writer := flvwriter.NewFLVWriter(paths[0], paths[1], url, ctx.Writer)
	// fmt.Println(paths[0], paths[1], url)
	// fmt.Println("start serve!")
	// s.rtmpstream.HandleWriter(writer)
	// writer.Wait()
	stream, err := s.rtmpstream.GetStream(path)
	if err != nil {
		ctx.Writer.WriteString(err.Error())
		ctx.Status(http.StatusNotFound)
		return
	}
	stream.HandleHTTPFLVWriter(ctx.Writer)

}

func (s *HTTPServer) streamsInfo(ctx *gin.Context) {
	msgs := &streamsInfo{}

	s.rtmpstream.GetStreams().Range(func(key interface{}, value interface{}) bool {
		if s, ok := value.(*rtmp.Stream); ok {
			if s.GetReader() != nil {
				msg := stream{key.(string), s.GetReader().Info().UID}
				msgs.Publishers = append(msgs.Publishers, msg)
			}
		}
		return true
	})

	s.rtmpstream.GetStreams().Range(func(key interface{}, value interface{}) bool {
		ws := value.(*rtmp.Stream).GetWs()

		ws.Range(func(key interface{}, value interface{}) bool {
			if pw, ok := value.(*rtmp.PackWriterCloser); ok {
				if pw.GetWriter() != nil {
					msg := stream{key.(string), pw.GetWriter().Info().UID}
					msgs.Players = append(msgs.Players, msg)
				}
			}
			return true
		})
		return true
	})

	resp, _ := json.Marshal(msgs)
	ctx.Header("Content-Type", "application/json")
	ctx.Writer.Write(resp)
}

func (s *HTTPServer) authRequire(ctx *gin.Context) {}

func (s *HTTPServer) Serve() error {
	ln, err := net.Listen("tcp", s.lnaddr)
	if err != nil {
		return fmt.Errorf("listen %s fail:%s", s.lnaddr, err)
	}
	s.ln = ln
	httpserver := http.Server{Handler: s.e}
	httpserver.Serve(s.ln)
	return nil
}

func (s *HTTPServer) Close() error {
	return s.ln.Close()
}

func (s *HTTPServer) pathMapping(ss string) (string, error) {
	return "", nil
}
