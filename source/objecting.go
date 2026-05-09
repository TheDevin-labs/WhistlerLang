package main

import (
	"encoding/binary"
	"errors"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

const objMagicSimple = "WHSTLOBJ"
const objVer = uint16(0x0102)

type ObjectInfo struct {
	Magic     string
	Version   uint16
	Arch      string
	Timestamp int64
	Payload   []byte
	Script    string
}

func archNameFromRuntime() string {
	switch runtime.GOARCH {
	case "amd64":
		return "amd64"
	case "arm64":
		return "arm64"
	case "386":
		return "386"
	case "arm":
		return "arm"
	default:
		return runtime.GOARCH
	}
}

func BuildScriptToObject(scriptPath, outPath, archOverride string) (string, error) {
	src, err := os.ReadFile(scriptPath)
	if err != nil {
		return "", err
	}
	payload := src
	ts := time.Now().Unix()
	arch := archOverride
	if arch == "" {
		arch = archNameFromRuntime()
	}
	if outPath == "" {
		base := filepath.Base(scriptPath)
		outPath = filepath.Join(filepath.Dir(scriptPath), "release", base+".o")
	}
	if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
		return "", err
	}
	f, err := os.Create(outPath)
	if err != nil {
		return "", err
	}
	defer f.Close()
	magic := make([]byte, 8)
	copy(magic, []byte(objMagicSimple))
	if _, err := f.Write(magic); err != nil {
		return "", err
	}
	if err := binary.Write(f, binary.LittleEndian, objVer); err != nil {
		return "", err
	}
	archB := []byte(arch)
	if err := binary.Write(f, binary.LittleEndian, uint16(len(archB))); err != nil {
		return "", err
	}
	if _, err := f.Write(archB); err != nil {
		return "", err
	}
	if err := binary.Write(f, binary.LittleEndian, ts); err != nil {
		return "", err
	}
	if err := binary.Write(f, binary.LittleEndian, uint32(len(payload))); err != nil {
		return "", err
	}
	if _, err := f.Write(payload); err != nil {
		return "", err
	}
	return outPath, nil
}

func GetObjectInfo(path string) (ObjectInfo, error) {
	var info ObjectInfo
	f, err := os.Open(path)
	if err != nil {
		return info, err
	}
	defer f.Close()
	var magic [8]byte
	if _, err := io.ReadFull(f, magic[:]); err != nil {
		return info, err
	}
	info.Magic = string(magic[:])
	if info.Magic != objMagicSimple {
		return info, errors.New("bad magic")
	}
	if err := binary.Read(f, binary.LittleEndian, &info.Version); err != nil {
		return info, err
	}
	var archLen uint16
	if err := binary.Read(f, binary.LittleEndian, &archLen); err != nil {
		return info, err
	}
	archB := make([]byte, archLen)
	if _, err := io.ReadFull(f, archB); err != nil {
		return info, err
	}
	info.Arch = string(archB)
	if err := binary.Read(f, binary.LittleEndian, &info.Timestamp); err != nil {
		return info, err
	}
	var payloadLen uint32
	if err := binary.Read(f, binary.LittleEndian, &payloadLen); err != nil {
		return info, err
	}
	payload := make([]byte, payloadLen)
	if _, err := io.ReadFull(f, payload); err != nil {
		return info, err
	}
	info.Payload = payload
	info.Script = "<embedded>"
	return info, nil
}
