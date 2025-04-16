package main

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
)

var mu sync.Mutex
var count int

func main() {
    http.HandleFunc("/", handler)
    http.HandleFunc("/count", counter)
    http.HandleFunc("/lissajous", lissajousHandler)
    log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "%s %s %s\n", r.Method, r.URL, r.Proto)
    for k, v := range r.Header {
        fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
    }
    fmt.Fprintf(w, "Host = %q\n", r.Host)
    fmt.Fprintf(w, "RemoteAddr = %q\n", r.RemoteAddr)
    if err := r.ParseForm(); err != nil {
        log.Print(err)
    }
    for k, v := range r.Form {
        fmt.Fprintf(w, "Form[%q] = %q\n", k, v)
    }
}

func counter(w http.ResponseWriter, r *http.Request) {
    mu.Lock()
    fmt.Fprintf(w, "Count %d\n", count)
    mu.Unlock()
}

func lissajousHandler(w http.ResponseWriter, r *http.Request) {
    // Парсим параметры запроса
    r.ParseForm()

    // Считываем параметр cycles из URL
    cyclesStr := r.Form.Get("cycles")
    cycles := 5 // Значение по умолчанию
    if cyclesStr != "" {
        var err error
        cycles, err = strconv.Atoi(cyclesStr)
        if err != nil {
            http.Error(w, "Invalid 'cycles' parameter", http.StatusBadRequest)
            return
        }
    }

    // Генерируем GIF с учетом параметра cycles
    lissajous(w, cycles)
}

func lissajous(out io.Writer, cycles int) {
    const (
        res    = 0.001 // Угловое разрешение
        size   = 100   // Размер изображения [-size..+size]
        nframes = 64   // Количество кадров
        delay   = 8    // Задержка между кадрами (в 10 мс)
    )

    freq := rand.Float64() * 3.0
    anim := gif.GIF{LoopCount: nframes}
    phase := 0.0

    for i := 0; i < nframes; i++ {
        rect := image.Rect(0, 0, 2*size+1, 2*size+1)
        img := image.NewPaletted(rect, []color.Color{color.White, color.Black})
        for t := 0.0; t < float64(cycles)*2*math.Pi; t += res {
            x := math.Sin(t)
            y := math.Sin(t*freq + phase)
            img.SetColorIndex(size+int(x*size+0.5), size+int(y*size+0.5), 1)
        }

        phase += 0.1
        anim.Delay = append(anim.Delay, delay)
        anim.Image = append(anim.Image, img)
    }
    gif.EncodeAll(out, &anim)
}