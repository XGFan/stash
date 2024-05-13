package ffmpeg

import (
	"reflect"
	"testing"
)

func TestHWGenerateRK264(t *testing.T) {
	hw := NewHWGenerate(VideoCodecRK264)

	if got := hw.Codec(); got != VideoCodecRK264 {
		t.Errorf("Codec() = %v, want %v", got, VideoCodecRK264)
	}

	wantInput := []string{
		"-init_hw_device", "rkmpp=rk",
		"-filter_hw_device", "rk",
		"-hwaccel", "rkmpp",
		"-hwaccel_output_format", "drm_prime",
	}
	if got := []string(hw.InputArgs()); !reflect.DeepEqual(got, wantInput) {
		t.Errorf("InputArgs() = %#v, want %#v", got, wantInput)
	}

	if got, want := hw.ScaleFilter(640, -2), VideoFilter("scale_rkrga=640:-2:format=nv12"); got != want {
		t.Errorf("ScaleFilter() = %q, want %q", got, want)
	}

	if got, want := hw.ScaleDownloadFilter(-2, 160), VideoFilter("scale_rkrga=-2:160:format=nv12,hwdownload,format=nv12"); got != want {
		t.Errorf("ScaleDownloadFilter() = %q, want %q", got, want)
	}

	if got, want := []string(hw.EncoderArgs()), []string{"-rc_mode", "VBR"}; !reflect.DeepEqual(got, want) {
		t.Errorf("EncoderArgs() = %#v, want %#v", got, want)
	}
}

func TestHWGenerateNotSupported(t *testing.T) {
	f := &FFMpeg{hwCodecSupport: []VideoCodec{VideoCodecN264, VideoCodecV264}}
	if got := f.HWGenerate(); got != nil {
		t.Errorf("HWGenerate() = %v, want nil", got)
	}

	f = &FFMpeg{hwCodecSupport: []VideoCodec{VideoCodecRK264}}
	if got := f.HWGenerate(); got == nil || got.Codec() != VideoCodecRK264 {
		t.Errorf("HWGenerate() = %v, want RK264", got)
	}
}
