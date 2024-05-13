package ffmpeg

import "fmt"

// HWGenerate provides hardware acceleration arguments for offline generation
// tasks (preview videos, sprite screenshots, phash screenshots).
//
// Unlike live transcoding, generation tasks cannot probe every input in advance,
// so callers are expected to fall back to the software path when a hardware
// accelerated command fails.
//
// Currently only Rockchip MPP (rkmpp) is supported.
type HWGenerate struct {
	codec VideoCodec
}

// HWGenerate returns the hardware acceleration available for generation tasks,
// or nil if none is available.
func (f *FFMpeg) HWGenerate() *HWGenerate {
	for _, codec := range f.getHWCodecSupport() {
		switch codec {
		case VideoCodecRK264:
			return NewHWGenerate(codec)
		}
	}

	return nil
}

// NewHWGenerate returns a HWGenerate for the given hardware codec.
// Prefer FFMpeg.HWGenerate, which only returns codecs detected as supported.
func NewHWGenerate(codec VideoCodec) *HWGenerate {
	return &HWGenerate{codec: codec}
}

// Codec returns the hardware video encoder.
func (h *HWGenerate) Codec() VideoCodec {
	return h.codec
}

// InputArgs returns the arguments to be placed before the input in order to
// initialise the hardware device and enable hardware decoding.
func (h *HWGenerate) InputArgs() Args {
	var args Args
	switch h.codec {
	case VideoCodecRK264:
		args = append(args,
			"-init_hw_device", "rkmpp=rk",
			"-filter_hw_device", "rk",
			"-hwaccel", "rkmpp",
			"-hwaccel_output_format", "drm_prime",
		)
	}
	return args
}

// ScaleFilter returns a hardware scaling filter. Frames remain in device memory,
// suitable for feeding into Codec().
// Use -n for w or h to maintain aspect ratio as a multiple of n.
func (h *HWGenerate) ScaleFilter(w, h2 int) VideoFilter {
	var vf VideoFilter
	switch h.codec {
	case VideoCodecRK264:
		// force nv12: 10-bit sources decode to nv15 which h264_rkmpp does not accept
		vf = vf.Append(fmt.Sprintf("scale_rkrga=%v:%v:format=nv12", w, h2))
	}
	return vf
}

// ScaleDownloadFilter returns a hardware scaling filter whose output is
// downloaded back into system memory, suitable for software encoders such as
// bmp used for screenshots.
func (h *HWGenerate) ScaleDownloadFilter(w, h2 int) VideoFilter {
	vf := h.ScaleFilter(w, h2)
	if vf == "" {
		return vf
	}
	vf = vf.Append("hwdownload")
	vf = vf.Append("format=nv12")
	return vf
}

// EncoderArgs returns the encoder specific arguments for Codec().
func (h *HWGenerate) EncoderArgs() Args {
	var args Args
	switch h.codec {
	case VideoCodecRK264:
		args = append(args, "-rc_mode", "VBR")
	}
	return args
}
