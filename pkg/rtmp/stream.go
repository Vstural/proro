package rtmp

import (
	"errors"
	"fmt"
	"io"
	"log"
	"proro/pkg/flvwriter"
	"proro/pkg/rtmp/av"
	"proro/pkg/rtmp/cache"
	"sync"
	"time"
)

var (
	EmptyID = ""
)

type RtmpStream struct {
	streams *sync.Map //key
}

func NewRtmpStream() *RtmpStream {
	ret := &RtmpStream{
		streams: &sync.Map{},
	}
	go ret.CheckAlive()
	return ret
}

func (rs *RtmpStream) HandleReader(r av.ReadCloser) {
	info := r.Info()
	log.Printf("HandleReader: info[%v]", info)

	var stream *Stream
	i, ok := rs.streams.Load(info.Key)
	if stream, ok = i.(*Stream); ok {
		stream.TransStop()
		id := stream.ID()
		if id != EmptyID && id != info.UID {
			ns := NewStream()
			stream.Copy(ns)
			stream = ns
			rs.streams.Store(info.Key, ns)
		}
	} else {
		stream = NewStream()
		rs.streams.Store(info.Key, stream)
		stream.info = info
	}

	stream.AddReader(r)
}

func (rs *RtmpStream) GetStream(key string) (*Stream, error) {
	item, ok := rs.streams.Load(key)
	if !ok {
		return nil, errors.New("stream not found")
	}
	return item.(*Stream), nil
}

func (rs *RtmpStream) HandleWriter(w av.WriteCloser) {
	info := w.Info()
	log.Printf("HandleWriter: info[%v]", info)

	var s *Stream
	_, ok := rs.streams.Load(info.Key)
	if !ok {
		fmt.Printf("try get stream with key %s not found\n", info.Key)
		s = NewStream()
		rs.streams.Store(info.Key, s)
		s.info = info
	} else {
		item, ok := rs.streams.Load(info.Key)
		if ok {
			s = item.(*Stream)
			s.AddWriter(w)
		}
	}
}

func (rs *RtmpStream) GetStreams() *sync.Map {
	return rs.streams
}

func (rs *RtmpStream) CheckAlive() {
	// todo: check alive?
	for {
		<-time.After(5 * time.Second)
		rs.streams.Range(func(key interface{}, item interface{}) bool {
			v := item.(*Stream)
			if v.CheckAlive() == 0 {
				fmt.Printf("check alive %s fail,remove\n", key)
				rs.streams.Delete(key)
			}
			return true
		})
	}
}

type Stream struct {
	isStart bool
	cache   *cache.Cache
	r       av.ReadCloser
	ws      *sync.Map
	info    av.Info
}

type PackWriterCloser struct {
	init bool
	w    av.WriteCloser
}

func (p *PackWriterCloser) GetWriter() av.WriteCloser {
	return p.w
}

func NewStream() *Stream {
	return &Stream{
		cache: cache.NewCache(),
		ws:    &sync.Map{},
	}
}

func NewStreamInfo(app, title, url string) *Stream {
	return &Stream{
		cache: cache.NewCache(),
		ws:    &sync.Map{},
	}
}

func (s *Stream) HandleHTTPFLVWriter(w io.Writer) {
	flvwriter := flvwriter.NewFLVWriter(paths[0], paths[1], url, ctx.Writer)
	s.AddWriter(flvwriter)
	flvwriter.Wait()
}

func (s *Stream) ID() string {
	if s.r != nil {
		return s.r.Info().UID
	}
	return EmptyID
}

func (s *Stream) GetReader() av.ReadCloser {
	return s.r
}

func (s *Stream) GetWs() *sync.Map {
	return s.ws
}

func (s *Stream) Copy(dst *Stream) {
	s.ws.Range(func(key interface{}, value interface{}) bool {
		v := value.(*PackWriterCloser)
		s.ws.Delete(key)
		v.w.CalcBaseTimestamp()
		dst.AddWriter(v.w)
		return true
	})
}

func (s *Stream) AddReader(r av.ReadCloser) {
	s.r = r
	go s.TransStart()
}

func (s *Stream) AddWriter(w av.WriteCloser) {
	info := w.Info()
	pw := &PackWriterCloser{w: w}
	s.ws.Store(info.UID, pw)
}

func (s *Stream) TransStart() {
	s.isStart = true
	var p av.Packet
	log.Printf("TransStart:%v", s.info)
	for {
		if !s.isStart {
			s.closeInter()
			return
		}
		err := s.r.Read(&p)
		if err != nil {
			s.closeInter()
			s.isStart = false
			return
		}

		s.cache.Write(p)

		s.ws.Range(func(key interface{}, value interface{}) bool {
			v := value.(*PackWriterCloser)
			if !v.init {
				//log.Printf("cache.send: %v", v.w.Info())
				if err = s.cache.Send(v.w); err != nil {
					log.Printf("[%s] send cache packet error: %v, remove", v.w.Info(), err)
					s.ws.Delete(key)
					return true
				}
				v.init = true
			} else {
				new_packet := p
				//writeType := reflect.TypeOf(v.w)
				//log.Printf("w.Write: type=%v, %v", writeType, v.w.Info())
				if err = v.w.Write(&new_packet); err != nil {
					log.Printf("[%s] write packet error: %v, remove", v.w.Info(), err)
					s.ws.Delete(key)
				}
			}
			return true
		})
	}
}

func (s *Stream) TransStop() {
	log.Printf("TransStop: %s", s.info.Key)
	if s.isStart && s.r != nil {
		s.r.Close(errors.New("stop old"))
	}
	s.isStart = false
}

func (s *Stream) CheckAlive() (n int) {
	if s.r != nil && s.isStart {
		if s.r.Alive() {
			n++
		} else {
			s.r.Close(errors.New("read timeout"))
		}
	}

	s.ws.Range(func(key interface{}, value interface{}) bool {
		v := value.(*PackWriterCloser)
		if v.w != nil {
			if !v.w.Alive() && s.isStart {
				s.ws.Delete(key)
				v.w.Close(errors.New("write timeout"))
				return true
			}
			n++
		}
		return true
	})
	return
}

func (s *Stream) closeInter() {
	s.ws.Range(func(key interface{}, value interface{}) bool {
		v := value.(*PackWriterCloser)
		if v.w != nil {
			if v.w.Info().IsInterval() {
				v.w.Close(errors.New("closed"))
				s.ws.Delete(key)
				log.Printf("[%v] player closed and remove\n", v.w.Info())
			}
		}
		return true
	})
}
