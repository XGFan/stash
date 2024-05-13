package transcoder

import (
	"reflect"
	"testing"

	"github.com/stashapp/stash/pkg/ffmpeg"
)

func TestScreenshotTimeDefaultUsesFastSeek(t *testing.T) {
	options := ScreenshotOptions{
		OutputPath: "out.jpg",
		OutputType: ScreenshotOutputTypeImage2,
	}

	got := ScreenshotTime("input.webm", 12.5, options)
	want := []string{
		"-v", "error",
		"-y",
		"-ss", "12.5",
		"-i", "input.webm",
		"-frames:v", "1",
		"-f", "image2",
		"out.jpg",
	}

	if !reflect.DeepEqual([]string(got), want) {
		t.Fatalf("ScreenshotTime() = %#v, want %#v", []string(got), want)
	}
}

func TestScreenshotTimeSlowSeek(t *testing.T) {
	options := ScreenshotOptions{
		OutputPath: "out.jpg",
		OutputType: ScreenshotOutputTypeImage2,
		SlowSeek:   true,
	}

	got := ScreenshotTime("input.webm", 12.5, options)
	want := []string{
		"-v", "error",
		"-y",
		"-i", "input.webm",
		"-ss", "12.5",
		"-frames:v", "1",
		"-f", "image2",
		"out.jpg",
	}

	if !reflect.DeepEqual([]string(got), want) {
		t.Fatalf("ScreenshotTime() = %#v, want %#v", []string(got), want)
	}
}

func TestScreenshotTimeHW(t *testing.T) {
	options := ScreenshotOptions{
		OutputPath: "-",
		OutputType: ScreenshotOutputTypeBMP,
		Width:      160,
		HW:         ffmpeg.NewHWGenerate(ffmpeg.VideoCodecRK264),
	}

	got := ScreenshotTime("input.mp4", 12.5, options)
	want := []string{
		"-v", "error",
		"-y",
		"-init_hw_device", "rkmpp=rk",
		"-filter_hw_device", "rk",
		"-hwaccel", "rkmpp",
		"-hwaccel_output_format", "drm_prime",
		"-ss", "12.5",
		"-i", "input.mp4",
		"-frames:v", "1",
		"-vf", "scale_rkrga=160:-2:format=nv12,hwdownload,format=nv12",
		"-c:v", "bmp",
		"-f", "rawvideo",
		"-",
	}

	if !reflect.DeepEqual([]string(got), want) {
		t.Fatalf("ScreenshotTime() = %#v, want %#v", []string(got), want)
	}
}
