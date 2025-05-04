# Time Server

[![Go Report Card](https://goreportcard.com/badge/github.com/yourusername/timeserver)](https://goreportcard.com/report/github.com/yourusername/timeserver)
[![GoDoc](https://godoc.org/github.com/yourusername/timeserver?status.svg)](https://godoc.org/github.com/yourusername/timeserver)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Time Server adalah implementasi server TCP konkuren yang ditulis dalam Go (Golang). Server ini menyediakan layanan waktu sederhana yang dapat diakses melalui koneksi TCP.

## Fitur Utama

- 🚀 **Server TCP Konkuren** - Menangani banyak koneksi klien secara bersamaan menggunakan goroutines
- 🕰️ **Layanan Waktu** - Menyediakan waktu server saat ini saat diminta
- 📝 **Logging Komprehensif** - Mencatat semua aktivitas server dan koneksi klien
- 🛑 **Graceful Shutdown** - Menangani penghentian dengan bersih saat menerima sinyal sistem

## Instalasi

### Prasyarat

- Go 1.16 atau lebih tinggi

### Langkah Instalasi

1. Clone repository:
   ```bash
   git clone https://github.com/yourusername/timeserver.git
   cd timeserver
   ```

2. Build project:
   ```bash
   go build -o timeserver
   ```

## Penggunaan

### Menjalankan Server

```bash
./timeserver
```

Server secara default akan berjalan pada `0.0.0.0:45000`.

### Opsi Konfigurasi

Anda dapat menyesuaikan konfigurasi dengan flag command line:

```bash
./timeserver -host=localhost -port=8080 -loglevel=debug
```

Flag yang tersedia:
- `-host`: Host address untuk binding (default: "0.0.0.0")
- `-port`: Port untuk listening (default: 45000)
- `-loglevel`: Level log (debug, info, warn, error) (default: "info")

## Protokol Komunikasi

Time Server mendukung protokol berbasis teks sederhana:

| Perintah | Deskripsi | Respons |
|----------|-----------|---------|
| `TIME`   | Mendapatkan waktu server saat ini | `JAM HH:MM:SS` |
| `QUIT`   | Menutup koneksi | *Koneksi ditutup* |

### Contoh Interaksi

```
$ nc localhost 45000
Welcome to Time Server. Send 'TIME' for current time, 'QUIT' to disconnect.
TIME
JAM 15:04:05
QUIT
```

## Struktur Project

```
timeserver/
├── main.go       # Entry point aplikasi
├── config.go     # Konfigurasi server
├── logger.go     # Sistem logging
├── server.go     # Implementasi server
├── go.mod
└── README.md
```

## Arsitektur

Time Server dibangun dengan arsitektur berikut:

1. **Main Application** - Menginisialisasi server, mengelola sinyal shutdown
2. **Server Component** - Menangani koneksi jaringan dan protokol
3. **Configuration** - Mengelola pengaturan server
4. **Logger** - Menyediakan fasilitas logging terstruktur

### Diagram Alur

```mermaid
graph TD
    MainApp["Main App"] --> Server["Server"]
    Logger["Logger"] <--> Server
    Server --> Config["Config"]
    Server --> ClientHandlers["Client Handlers"]
    ClientHandlers --> Goroutines["(Goroutines)"]
```

## Konkurensi dan Keamanan Thread

Time Server menggunakan model konkurensi Go dengan:

- **Goroutines** untuk menangani koneksi klien secara bersamaan
- **Mutex/RWMutex** untuk melindungi variabel bersama
- **WaitGroup** untuk penghentian yang terkoordinasi
- **Context** untuk pengelolaan lifecycle

## Penanganan Error

Server mengimplementasikan strategi penanganan error yang komprehensif:

- Error wrapping dengan `fmt.Errorf` dan `%w`
- Sentinel errors untuk error yang umum
- Timeout untuk operasi I/O
- Graceful degradation saat terjadi error

## Pengujian

Untuk menjalankan test:

```bash
go test ./...
```

Project ini mencakup:
- Unit tests untuk komponen individual
- Integration tests untuk fungsionalitas end-to-end
- Benchmarks untuk komponen kritis

## Performa

Time Server dirancang untuk performa tinggi dengan penggunaan resource yang efisien:

- Konkurensi goroutine yang ringan
- Penggunaan memory yang dioptimalkan dengan buffer tetap
- Penggunaan RWMutex untuk pembacaan simultan

## Kasus Penggunaan

Time Server cocok untuk:

- Demo/pembelajaran untuk pemrograman jaringan dengan Go
- Basis untuk membangun server TCP yang lebih kompleks
- Pengujian aplikasi klien sederhana
- Tutorial untuk konkurensi dan pengelolaan resource dalam Go

## Kontribusi

Kontribusi selalu disambut! Jika Anda ingin berkontribusi:

1. Fork repository
2. Buat branch fitur (`git checkout -b feature/amazing-feature`)
3. Commit perubahan Anda (`git commit -m 'Add some amazing feature'`)
4. Push ke branch (`git push origin feature/amazing-feature`)
5. Buka Pull Request

## Lisensi

Project ini dilisensikan di bawah Lisensi MIT - lihat file [LICENSE](LICENSE) untuk detail.

## Kontak

Nama Anda - [@twitter_handle](https://twitter.com/twitter_handle) - email@example.com

Project Link: [https://github.com/yourusername/timeserver](https://github.com/yourusername/timeserver)

## Acknowledgments

- [Go Documentation](https://golang.org/doc/)
- [Effective Go](https://golang.org/doc/effective_go.html)
- [Go Concurrency Patterns](https://blog.golang.org/concurrency-patterns)