package proxy

import (
	"encoding/binary"
	"errors"
	gox "game/utils"
	"io"
	"log"
)

const (
	Magic        = 1995
	Version      = 1
	HeaderLength = 16
)

type CMD uint32

const (
	CMD_OK    CMD = 0x00010000 //服务和客户端通用
	CMD_ERROR CMD = 0x00010001 //服务和客户端通用

	CMD_CLIENT_HOST     CMD = 0x00010002 //主机命令 用于握手交换主机信息
	CMD_CLIENT_EXCHANGE CMD = 0x00010003 //交换数据命令 用于传输数据
)

var (
	ErrNilReader      = errors.New("nil reader")
	ErrNilBuffer      = errors.New("nil byte buffer")
	ErrParseBuffer    = errors.New("error parse byte to Bat")
	ErrLengthMisMatch = errors.New("error length not Match")
)

// Bat 总长度 2 + 2 + 4 + 8 +len(buffer)
type Bat struct {
	Magic   uint16 //魔法数字 2
	Version uint16 //版本号 2
	Cmd     CMD    //
	Length  uint64 // 数命令 4据包总长度 8
	Packet  []byte //数据包
}

func (bat *Bat) Write(w io.Writer) (n int, err error) {
	return w.Write(bat.toByte())
}

func (bat *Bat) WriteEncrypt(w io.Writer) (n int, err error) {
	return w.Write(bat.Encrypt())
}

// 转为字节码
func (bat *Bat) toByte() []byte {
	buffer := make([]byte, 0)
	buffer = binary.BigEndian.AppendUint16(buffer, bat.Magic)
	buffer = binary.BigEndian.AppendUint16(buffer, bat.Version)
	buffer = binary.BigEndian.AppendUint32(buffer, uint32(bat.Cmd))
	buffer = binary.BigEndian.AppendUint64(buffer, bat.Length)
	buffer = append(buffer, bat.Packet...)
	return buffer
}

func (bat *Bat) Encrypt() []byte {
	encrypt, err := gox.EncryptByAes(bat.toByte())
	if err != nil {
		encrypt = ""
		log.Printf("encrypt err %v", err)
	}
	buffer := make([]byte, 0)
	buffer = binary.BigEndian.AppendUint16(buffer, Magic)
	buffer = binary.BigEndian.AppendUint32(buffer, uint32(len(encrypt)))
	buffer = append(buffer, []byte(encrypt)...)
	return buffer
}

func DecryptBat(r io.Reader) (*Bat, error) {
	length := make([]byte, 6)
	if _, err := io.ReadFull(r, length); err != nil {
		return nil, err
	}
	magic := binary.BigEndian.Uint16(length[:2])
	if magic != Magic {
		return nil, ErrParseBuffer
	}
	len := binary.BigEndian.Uint32(length[2:])
	encryptBuf := make([]byte, len)
	if _, err := io.ReadFull(r, encryptBuf); err != nil {
		return nil, err
	}
	raw, err := gox.DecryptByAes(string(encryptBuf))
	if err != nil {
		log.Printf("encrypt err %v", err)
		return nil, err
	}
	return FillBat(raw)
}

func newBat(version uint16, Cmd CMD, Packet []byte) *Bat {
	if Packet == nil {
		Packet = make([]byte, 0)
	}
	bat := &Bat{
		Magic:   Magic,
		Version: version,
		Cmd:     Cmd,
		Length:  uint64(HeaderLength + len(Packet)),
		Packet:  Packet,
	}
	return bat
}

func ZeroBat(Cmd CMD, Packet []byte) *Bat {
	return newBat(Version, Cmd, Packet)
}

func FillBat(buffer []byte) (*Bat, error) {
	if buffer == nil || len(buffer) < HeaderLength {
		return nil, ErrNilBuffer
	}
	magic := binary.BigEndian.Uint16(buffer[0:2])
	if magic != Magic {
		return nil, ErrParseBuffer
	}
	version := binary.BigEndian.Uint16(buffer[2:4])
	cmd := binary.BigEndian.Uint32(buffer[4:8])
	length := binary.BigEndian.Uint64(buffer[8:HeaderLength])
	packet := buffer[16:length]
	bat := newBat(version, CMD(cmd), packet)
	return bat, nil
}

func ReadBat(r io.Reader) (*Bat, error) {
	if r == nil {
		return nil, ErrNilReader
	}
	header := make([]byte, HeaderLength)
	if rd, err := r.Read(header); err != nil || rd != HeaderLength {
		if err == nil {
			err = ErrLengthMisMatch
		}
		return nil, err
	}
	magic := binary.BigEndian.Uint16(header[0:2])
	if magic != Magic {
		return nil, ErrParseBuffer
	}
	version := binary.BigEndian.Uint16(header[2:4])
	cmd := binary.BigEndian.Uint32(header[4:8])
	length := binary.BigEndian.Uint64(header[8:HeaderLength])
	pLen := length - HeaderLength
	packet := make([]byte, pLen)
	if pLen > 0 {
		if rd, err := r.Read(packet); err != nil || uint64(rd) != (pLen) {
			log.Printf("read err: %v", err)
			return nil, err
		}
	}
	bat := newBat(version, CMD(cmd), packet)
	return bat, nil
}
