# Go Convert
## Script ini untuk menangani sebuah case mengupload file excell dengan jumlah yang sangat besar lalu dijadikan file JSON

1. Streaming File Excel:
excelize.OpenFile() digunakan untuk membuka file Excel tanpa memuat seluruh file ke dalam memori.
Fungsi f.Rows("Sheet1") digunakan untuk membaca baris dari file secara bertahap.

2. Manajemen Memori:
Data Excel tidak ditampung seluruhnya di memori. Baris dibaca satu per satu dan segera diproses untuk diubah menjadi JSON.
Menggunakan buffered channel (rowsChan) untuk meminimalkan penggunaan memori dengan membaca data Excel secara batch.

3. Worker Goroutine:
Fungsi worker() digunakan untuk memproses setiap baris Excel menjadi JSON dalam worker goroutine.
Worker pool dibuat dengan jumlah worker sebanyak jumlah CPU yang tersedia (runtime.NumCPU()).

4. JSON Streaming:
writeJSONToFile() menggunakan json.NewEncoder() untuk menulis baris JSON secara streaming ke file output, sehingga tidak ada akumulasi data di memori.
JSON array dibuka ([) di awal dan ditutup (]) di akhir setelah semua baris diproses.

5. Sinkronisasi:
sync.WaitGroup digunakan untuk menunggu semua worker goroutine menyelesaikan tugas mereka.
doneChan digunakan untuk memastikan penulisan ke file JSON selesai sebelum program berhenti.



# Optimasi Tambahan:
- Buffering: Dengan channel yang di-buffer (rowsChan dan jsonLinesChan), kita bisa menjaga aliran data tanpa menunggu terlalu lama antara membaca dan menulis.
- Concurrency: Dengan memanfaatkan beberapa core CPU menggunakan goroutine, proses konversi dapat dilakukan lebih cepat.