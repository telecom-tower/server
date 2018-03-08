//go:generate protoc -I telecomtower telecomtower/telecomtower.proto --go_out=plugins=grpc:telecomtower

package server

import (
	"io"
	"net"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	pb "github.com/telecom-tower/server/telecomtower"
	"google.golang.org/grpc"
)

const (
	displayHeight = 8
	displayWidth  = 8
)

type wsEngine interface {
	Init() error
	Render() error
	Wait() error
	Fini()
	Leds(channel int) []uint32
}

type towerServer struct {
	ws wsEngine
}

func coordinatesToIndex(x int32, y int32) int {
	if x%2 == 0 {
		y = displayHeight - 1 - y
	}
	return int(x*displayHeight + y)
}

func colorToint(c *pb.Color) uint32 {
	return ((c.Red>>8)&0xff)<<16 + ((c.Green>>8)&0xff)<<8 + ((c.Blue >> 8) & 0xff)
}

func (s *towerServer) fill(fill *pb.Fill) error {
	log.Debugf("fill")
	color := colorToint(fill.Color)
	for i := 0; i < len(s.ws.Leds(0)); i++ {
		s.ws.Leds(0)[i] = color
	}
	return nil
}

func (s *towerServer) setPixels(pixels *pb.SetPixels) error {
	log.Debugf("set pixels")
	for _, pix := range pixels.Pixels {
		index := coordinatesToIndex(pix.Point.Column, pix.Point.Row)
		if index < 0 || index >= len(s.ws.Leds(0)) {
			return errors.New("Index out of bounds")
		}
		color := colorToint(pix.Color)
		s.ws.Leds(0)[index] = color
	}
	return nil
}

func (s *towerServer) drawLine(*pb.DrawLine) error {
	log.Debug("draw line (NYI)")
	// TODO: Not yet implemented
	return nil
}

func (s *towerServer) drawRectangle(rect *pb.DrawRectangle) error {
	log.Debug("draw rectangle")
	color := colorToint(rect.Color)
	x0 := rect.Point0.Column
	x1 := rect.Point1.Column
	if x1 < x0 {
		x0, x1 = x1, x0
	}
	y0 := rect.Point0.Column
	y1 := rect.Point1.Column
	if y1 < y0 {
		y0, y1 = y1, y0
	}
	if x0 < 0 || x1 >= displayWidth || y0 < 0 || y1 >= displayHeight {
		return errors.New("Index out of bounds")
	}
	for x := x0; x <= x1; x++ {
		for y := y0; y < y1; y++ {
			index := coordinatesToIndex(x, y)
			s.ws.Leds(0)[index] = color
		}
	}
	return nil
}

func (s *towerServer) drawBitmap(*pb.DrawBitmap) error {
	log.Debug("draw bitmap (NYI)")
	// TODO: Not yet implemented
	return nil
}

func (s *towerServer) writeText(*pb.WriteText) error {
	log.Debug("write text (NYI)")
	// TODO: Not yet implemented
	return nil
}

func (s *towerServer) hScroll(*pb.HScroll) error {
	log.Debug("horizontal scroll (NYI)")
	// TODO: Not yet implemented
	return nil
}

func (s *towerServer) vScroll(*pb.VScroll) error {
	log.Debug("vertical scroll (NYI)")
	// TODO: Not yet implemented
	return nil
}

// Draw implements the main task of the server, namely drawing on the display
func (s *towerServer) Draw(stream pb.TowerDisplay_DrawServer) error {
	var status error = nil
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			if status == nil {
				status = s.ws.Render()
			}
			msg := ""
			if status != nil {
				msg = status.Error()
			}
			return stream.SendAndClose(&pb.DrawResponse{
				Message: msg,
			})
		}
		if err != nil {
			return err
		}

		if status == nil {
			switch t := in.Type.(type) {
			case *pb.DrawRequest_Fill:
				status = s.fill(t.Fill)
			case *pb.DrawRequest_SetPixels:
				status = s.setPixels(t.SetPixels)
			case *pb.DrawRequest_DrawLine:
				status = s.drawLine(t.DrawLine)
			case *pb.DrawRequest_DrawRectangle:
				status = s.drawRectangle(t.DrawRectangle)
			case *pb.DrawRequest_DrawBitmap:
				status = s.drawBitmap(t.DrawBitmap)
			case *pb.DrawRequest_WriteText:
				status = s.writeText(t.WriteText)
			case *pb.DrawRequest_HScroll:
				status = s.hScroll(t.HScroll)
			case *pb.DrawRequest_VScroll:
				status = s.vScroll(t.VScroll)
			}
		}
	}
}

// Serve starts a grpc server and handles the requests
func Serve(listener net.Listener, ws2811 wsEngine, opts ...grpc.ServerOption) error {
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterTowerDisplayServer(grpcServer, &towerServer{
		ws: ws2811,
	})
	log.Infof("Telecom Tower Server running at %v\n", listener.Addr().String())
	err := grpcServer.Serve(listener)
	if err != nil {
		return errors.WithMessage(err, "failed to serve")
	}
	return nil
}
