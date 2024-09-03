package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-gst/go-gst/gst"
)

func main() {
	gstInit()
	go bench()
	go bench()
	go bench()
	go bench()
	go bench()
	go bench()
	go bench()
	go bench()
	go bench()
	go bench()
	go bench()
	go bench()
	go bench()
	go bench()
	bench()
}

func bench() {
	for {
		if err := test(); err != nil {
			fmt.Println("ERROR: ", err)
			return
		}
		fmt.Print(".")
		// runtime.GC()
	}
}

func gstInit() {
	gst.Init(nil)
}

func test() error {
	var pipeline *gst.Pipeline
	var err error
	strPipeline := `rtpsession name=r sdes="application/x-rtp-source-sdes,cname=497328127001,tool=IVR" rtp-profile=avpf
 udpsrc name=udpsrc_rtp port=36258 caps=application/x-rtp,media=audio,clock-rate=8000,payload=96,encoding-name=TELEPHONE-EVENT timeout=30000000000
   ! r.recv_rtp_sink
 udpsrc name=udpsrc_rtcp port=36259 timeout=30000000000
   ! r.recv_rtcp_sink
 audiotestsrc name=background-tone-src
   ! audiomixer name=amixer force-live=true
 amixer.
   ! mulawenc
   ! rtppcmupay name=rtppay pt=0
   ! application/x-rtp,media=audio,clock-rate=8000,payload=0,encoding-name=PCMU
   ! r.send_rtp_sink
 r.send_rtp_src ! identity name=rtp-in-inspector
   ! udpsink name=udp_rtp_sink host=127.0.0.1 port=6000 async=false
 r.send_rtcp_src
   ! udpsink name=udp_rtcp_sink host=127.0.0.1 port=6001 async=false`

	if pipeline, err = gst.NewPipelineFromString(strPipeline); err != nil {
		return err
	}
	if !pipeline.GetPipelineBus().AddWatch(handleEvent) {
		return errors.New("can't add watch")
	}
	defer pipeline.GetPipelineBus().RemoveWatch()

	err = pipeline.BlockSetState(gst.StatePlaying)
	if err != nil {
		return err
	}

	err = pipeline.BlockSetState(gst.StateReady)
	if err != nil {
		return err
	}

	err = pipeline.BlockSetState(gst.StatePlaying)
	if err != nil {
		return err
	}

	inspRTP, err := pipeline.GetElementByName("rtp-in-inspector")
	if err != nil {
		return err
	}

	h, err := inspRTP.Connect("handoff", func(self *gst.Element, buff *gst.Buffer) {})
	if err != nil {
		return err
	}
	defer inspRTP.HandlerDisconnect(h)

	time.Sleep(500 * time.Millisecond)

	return pipeline.BlockSetState(gst.StateNull)
}

func handleEvent(msg *gst.Message) bool {
	// fmt.Println(msg)
	return true
}
