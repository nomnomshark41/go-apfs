package cmd

import (
        "os"
        "io"

        "github.com/blacktop/go-apfs/pkg/disk"
)

func isRawAPFS(path string) bool {
        f, err := os.Open(path)
        if err != nil {
                return false
        }
        defer f.Close()
        fi, err := f.Stat()
        if err != nil || fi.Size() < 512 {
                return false
        }
        buf := make([]byte, 4)
        f.ReadAt(buf, fi.Size()-512)
        return string(buf) != "koly"
}

func openRawAPFS(path string) (disk.Device, error) {
        f, err := os.Open(path)
        if err != nil {
                return nil, err
        }
        fi, err := f.Stat()
        if err != nil {
                f.Close()
                return nil, err
        }
        // find NXSB magic offset
        buf := make([]byte, 4096)
        f.ReadAt(buf, 0)
        offset := int64(0)
        for i := 0; i < len(buf)-4; i++ {
                if buf[i] == 'N' && buf[i+1] == 'X' && buf[i+2] == 'S' && buf[i+3] == 'B' {
                        offset = int64(i)
                        break
                }
        }
        sr := io.NewSectionReader(f, offset, fi.Size()-offset)
        g := disk.NewGeneric(sr)
        return g, nil
}
