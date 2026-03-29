package cmd

import (
        "io"
        "os"

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
        // read from offset 0 - ObjPhysT header is at 0, NXSB magic at 32 is correct
        sr := io.NewSectionReader(f, 0, fi.Size())
        g := disk.NewGeneric(sr)
        return g, nil
}
