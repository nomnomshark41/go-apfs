package cmd

import "os"

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
