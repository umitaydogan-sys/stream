package recording

import (
	"encoding/binary"
)

func recordingParseAVCConfigToAnnexB(data []byte) []byte {
	if len(data) < 11 {
		return nil
	}
	cfg := data
	if len(data) > 5 {
		cfg = data[5:]
	}
	if len(cfg) < 7 {
		return nil
	}

	pos := 6
	numSPS := int(cfg[5] & 0x1F)
	out := make([]byte, 0, 256)

	for i := 0; i < numSPS; i++ {
		if pos+2 > len(cfg) {
			return out
		}
		l := int(binary.BigEndian.Uint16(cfg[pos : pos+2]))
		pos += 2
		if l <= 0 || pos+l > len(cfg) {
			return out
		}
		out = append(out, 0x00, 0x00, 0x00, 0x01)
		out = append(out, cfg[pos:pos+l]...)
		pos += l
	}

	if pos >= len(cfg) {
		return out
	}

	numPPS := int(cfg[pos])
	pos++
	for i := 0; i < numPPS; i++ {
		if pos+2 > len(cfg) {
			return out
		}
		l := int(binary.BigEndian.Uint16(cfg[pos : pos+2]))
		pos += 2
		if l <= 0 || pos+l > len(cfg) {
			return out
		}
		out = append(out, 0x00, 0x00, 0x00, 0x01)
		out = append(out, cfg[pos:pos+l]...)
		pos += l
	}

	return out
}

func recordingAVCCToAnnexB(data []byte) []byte {
	if len(data) < 4 {
		return data
	}
	out := make([]byte, 0, len(data)+32)
	pos := 0

	for pos+4 <= len(data) {
		naluLen := int(binary.BigEndian.Uint32(data[pos : pos+4]))
		pos += 4
		if naluLen <= 0 || pos+naluLen > len(data) {
			return data
		}
		out = append(out, 0x00, 0x00, 0x00, 0x01)
		out = append(out, data[pos:pos+naluLen]...)
		pos += naluLen
	}

	if len(out) == 0 {
		return data
	}
	return out
}

func recordingParseAACAudioSpecificConfig(rec *Recording, data []byte) {
	if rec == nil || len(data) < 4 {
		return
	}
	asc := data[2:]
	if len(asc) < 2 {
		return
	}
	audioObjectType := int((asc[0] >> 3) & 0x1F)
	freqIdx := int(((asc[0] & 0x07) << 1) | ((asc[1] >> 7) & 0x01))
	chCfg := int((asc[1] >> 3) & 0x0F)

	if audioObjectType >= 2 {
		rec.aacProfile = audioObjectType - 1
	}
	if freqIdx >= 0 && freqIdx <= 12 {
		rec.aacFreqIndex = freqIdx
	}
	if chCfg > 0 && chCfg <= 7 {
		rec.aacChannelCfg = chCfg
	}
}

func recordingAddADTSHeader(raw []byte, profile, freqIdx, chCfg int) []byte {
	frameLen := len(raw) + 7
	adts := make([]byte, 7)
	adts[0] = 0xFF
	adts[1] = 0xF1
	adts[2] = byte((profile<<6)&0xC0 | (freqIdx<<2)&0x3C | (chCfg>>2)&0x01)
	adts[3] = byte((chCfg&0x03)<<6 | (frameLen>>11)&0x03)
	adts[4] = byte((frameLen >> 3) & 0xFF)
	adts[5] = byte((frameLen&0x07)<<5 | 0x1F)
	adts[6] = 0xFC
	return append(adts, raw...)
}
