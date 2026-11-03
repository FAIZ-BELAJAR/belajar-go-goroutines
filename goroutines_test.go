package belajar_go_goroutines

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"
)

// Membuat Goroutines:
// Untuk Membuat goroutine di Golah sangatlah sederhana
// Kita hanya cukup menambahkan perintah go sebelum memanggil function yang akan kita jalankan dalam goroutine
// saat sebuah function kita jalankan dalam goroutine, function tersebut akan berjalan secara asynchronous, artinya tidak akan ditunggu sampai function tersebut selesai
// Aplikasi akan lanjut berjalan ke kode program selanjutnya tanpa menunggu goroutine yang kita buat selesai

func RunHelloWorld() {
	fmt.Println("Hello World")
}

func TestCreateGoroutine(t *testing.T) {
	go RunHelloWorld()
	fmt.Println("Ups")

	time.Sleep(1 * time.Second)
}

// Goroutine Sangat Ringan
// Seperti yang sebelumnya dijelaskan, bahwa goroutine itu sangat ringan
// Kita bisa membuat ribuan, bahkan sampi jutaan goroutine tanpa takut boros memory
// tidak seperti thread yang ukuran nya berat, goroutine sangatlah ringan

func DisplayNumber(number int32) {
	fmt.Println("Display", number)
}

func TestManyGoroutine(t *testing.T) {
	for i := 0; i < 10000; i++ {
		go DisplayNumber(int32(i))
	}

	time.Sleep(5 * time.Second)
}

// Pengenalan Channel
// Channel adalah tempat komunikasi secara synchronous yang bisa dilakukan oleh goroutine
// Di Channel terdapat perngirim dan penerima, biasanya pengirim dan penerima adalah gorutine yang berbeda
// saat melakukan pengiriman data ke channel, goroutine akan ter block, sampai ada yang menerima data tersebut
// maka dari itu, channel disebut sebagai alat komunikasi synchronous (blocking)
// Channel cocok sekali sebagai alternatif mekanisme seperti async await ang terdapat di beberapa bahasa pemrograman lain

// Karakteristik Channel
// Secara Default Channel hanya bisa menampung satu data, jika kita ingin menambahkan data lagi, harus menunggu data yang ada diambil channel
// Channel hanya bisa menerima satu jenis data
// Channel bisa diambil dari lebih dari satu goroutine
// Channel harus di close jika tidk digunakan, atau bisa menyebabkan memory leak

// Membuat Channel
// Channel di go direpresentasikan dalam tipe data chan /////////
// untuk membuat channel sangat mudah, kita bisa menggunakan make(), mirip seperti membuat map
// Namun saat pembuatan channel, kita harus tentukan tipe data apa yang bisa dimasukan ke dalam channel tersebut
// Contoh: make(chan tipe_data)

// Mengirim dan menerima data dari channel
// Seperti yang sudah dibahas sebelumnya, channel bisa digunakan untuk mengirim dan menerima data
// Untuk mengirim data kita bisa gunakan kode: channel <- data //////////////////
// Sedangkan untuk menerima data, kita bisa gunakan kode: data <- channel /////////////
// Jika selesai, jangan lupa untuk menurup channel menggunakan close() //////////
// Contoh: close(channel) ////////

func TestChannel(t *testing.T) {
	var channel = make(chan string)
	var hasil string = ""

	go func() {
		channel <- "Faiz"

	}()

	time.Sleep(3 * time.Second)

	go func(h *string) {
		*h = <-channel
	}(&hasil)

	time.Sleep(3 * time.Second)
	fmt.Println("ini datanya: ", hasil)

	close(channel)
}

// Channel Sebagai Parameter
// Dalam kenyataan pembuatan aplikasi, seringnya kita akan mengirim channel ke function lain via parameter
// Sebelumnya kita tahu di golang by default, parameter adalah pass by value, artinya value akan diduplikasi lalu dikirim ke function parameter, sehingga jika kita ingin mengirim data asli kita biasa gunakan pointer(agar pass by reference)
// Berbeda dengan channel, kita tidak perlu melakukan hal tersebut
// Karna channel akan otomatis pass by reference

func GiveMeResponse(channel chan string) {
	time.Sleep(time.Second * 2)
	channel <- "Faiz Annan Cuy"
}

func TestChannelAsParam(t *testing.T) {
	var channel chan string = make(chan string)

	go GiveMeResponse(channel)

	data := <-channel
	fmt.Println(data)
	close(channel)
}

// Channel In dan Out
// Saat kita mengirim channel sebagai parameter, isi function tersebut bisa mengirim dan menerima data dari channel tersebut
// kadang ktia ingin memberi tahu terhadap sebuah funciton, misal bahwa channel tersebut hanya digunakan untuk mengirim data, atau hanya dapat digunakan untuk menerima data
// Hal ini bisa dilakukan di parameter dengan cara menandai apakah channel ini digunakan untuk in (mengirim data masuk) atau out (mengirim data keluar)

func OnlyInChannel(c chan<- string, d string) {
	time.Sleep(time.Second * 1)
	c <- d
}

func OnlyOutChannel(c <-chan string) {
	time.Sleep(time.Second * 1)
	data := <-c
	fmt.Println(data)
}

func TestInOutChannel(t *testing.T) {
	var cenel chan string = make(chan string)

	go OnlyInChannel(cenel, "Faizz")
	go OnlyOutChannel(cenel)

	time.Sleep(time.Second * 6)

	close(cenel)
}

// Buffered Channel
// Seperti yang dijelaskan sebelumnya bahwa secara default, channel itu hanya bisa menerima 1 data
// Artinya jika kita menambah data ke-2, maka kita akan diminta menunggu sampai data ke 1 ada yang mengambil
// kadang kadang ada kasus dimana pengirim lebih cepat dibanding penerima, dalam hal ini jika kita menggunakan channel, maka ototmatis pengirim akan ikut lambat juga
// Untungnya ada Buffered Channel, yaitu buffer yang bisa digunakan untuk menampung data antrian di channel

// Buffer Capacity
// Kita bisa bebad memasukan berapa jumlah kapasitas antrian di dalam buffer
// Jika kita set misal 5, artinya kita bisa menerima 5 data di buffer
// Jika kita mengirim data ke 6, maka kita diminta untuk menunggu sampai buffer nya ada yang kosong
// ini cocok sekali ketika memang gorroutine yang menerima data lebih lambat dari yang mengirim data

// cara membuat nya adalah (channel, tipe_data, buffer)///////////

func TestBufferedChannel(t *testing.T) {
	channel := make(chan string, 3)
	defer close(channel)

	channel <- "Faiz"

	fmt.Println(<-channel)

	fmt.Println("Selesai")
}

// Range Channel
// Kadang kadang ada kasus sebuah channel dikirim data secara terus menerus oleh pengirim
// dan kadang tidak jelas kapan channel tersebut akan berhenti menerima data
// salah satu yang bisa kita lakukan adalah dengan menggunakan perulangan range ketika menerima data dari channel
// ketika sebuah channel di close(), maka secara otomatis perulangan tersebut akan berhenti
// ini lebih sederhana daripada kita melakukan engecekan channel secara manual

func TestRangeChannel(t *testing.T) {
	const diulang = 10
	publisher := func(c chan<- string, k string, sebanyak int8) {

		for i := 0; i < int(sebanyak); i++ {
			c <- k + strconv.Itoa(i)
		}

	}

	consumer := func(c <-chan string) {
		for k := range c {
			fmt.Println("ini adalah pesan nya: ", k)
		}
	}

	var cenel chan string = make(chan string)
	defer close(cenel)

	go publisher(cenel, "anjay", int8(diulang))

	go consumer(cenel)

	time.Sleep(time.Second*diulang + 3)

	fmt.Println("selesai")
}

// Select Channel
// Kadang ada kasus dimana kita membuat beberap channel dan menjalankan beberapa goroutine
// lalu kita ingin mendapatkan data dari semua channel tersebut
// untuk melakukan hal tersebut, kita bisa menggunakan select channel di go lang
// dengan select channel, kita bisa memilih data tercepat dari beberapa channel jika data datang secara bersamaan di beberapa channel, maka akan dipilih secara random

func TestSelectMultipleChannel(t *testing.T) {
	channel1 := make(chan string)
	channel2 := make(chan string)
	defer close(channel1)
	defer close(channel2)

	go GiveMeResponse(channel1)
	go GiveMeResponse(channel2)

	counter := 0

	for {
		select {
		case data := <-channel1:
			fmt.Println("Data dari channel 1", data)
			counter++
		case data := <-channel2:
			fmt.Println("Data dari channel 1", data)
			counter++
		}

		if counter == 2 {
			break
		}
	}

}

// Default Select
// Apa yang akan terjadi jika kita melakukan select terhadap channel yang ternyata tidak ada datanya?
// Maka kita akan menunggu sampai data ada
// Kadang mungkin kita ingin melakukan sesuatu jika misal semua channel tidak ada datanya ketika kita melakukan select channel
// Dalam select, kita bisa menambahkan defaukt, ddimana ini akan dieksekusi jika memang di semua channel yang kita select tidak ada datanya

func TestDefaultSelect(t *testing.T) {
	channel1 := make(chan string)
	channel2 := make(chan string)
	defer close(channel1)
	defer close(channel2)
	counter := 0
	for {
		select {
		case data := <-channel1:
			fmt.Println("Data dari channel 1: ", data)
			counter++
		case data := <-channel2:
			fmt.Println("Data dari channel 2: ", data)
			counter++
		default:
			fmt.Println("Menunggu data")
		}

		if counter == 2 {
			break
		}
	}
}

// Race Condition
// Masalah dengan goroutine
// Saat kita menggunakan goroutine, dia tidak hanya berjalan secara conccurent, tapi bisa parallel juga, karena bisa ada beberapa thread yang berjalan secara parallel
// Hal ini sangat berbahaya ketika kita melakukan manipulasi data variable yang sama oleh beberapa goroutine secara bersamaan
// Hal ini bisa menyebabkan masalah yang namanya Race condition

func TestRaceCondition(t *testing.T) {
	var x = 0
	for i := 0; i <= 100; i++ {
		go func() {
			for j := i; j <= 100; j++ {
				x = x + 1
			}
		}()
	}

	time.Sleep(5 * time.Second)
	fmt.Println("Counter: ", x)
}

// sync.Mutex
// Mutex adalah singkatan dari Mutual Exclusion
// Untuk mengatasi masalah race condition tersebut, di Go Lang terdapat sebuah struct bernama sync.Mutex
// Mutex bisa digunakan untuk melakukan locking dan unlocking, dimana ketika kita melakukan locking terhadap mutex, maka tidak ada yang bisa melakukan locking lagi sampai kita melakukan unlock
// Dengan demikian, jika ada beberapa goroutine melakukan lock terhadap Mutex, maka hanya ada 1 goroutine yang diperbolehkan, setelah goroutine tersebut melakukan unlock, baru goroutine selanjutnya diperbolehkan melakukan lock lagi
// Ini sangat cocok sebagai solusi ketika ada masalah race condition yang sebelumnya kita hadapi

func TestRaceConditionMutex(t *testing.T) {
	var m sync.Mutex
	var x = 0
	for i := 0; i <= 100; i++ {
		go func() {
			for j := i; j <= 100; j++ {
				m.Lock()
				x = x + 1
				m.Unlock()
			}
		}()
	}

	time.Sleep(5 * time.Second)
	fmt.Println("Counter: ", x)
}

// sync.RWMutex
// RWMutex adalah singkatan dari Read Write Mutex
// Kadang ada kasus dimana kita ingin melakukan locking tidak hanya pada proses mengubah data, tapi juga membaca data
// Kita sebenernya bisa menggunakan Mutex aja, namun masalahnya nanti akan rebutan antara proses membaca dan mengubah
// di Go Lang telah di sediakan struct RWMutex (Read Write Mutex) untuk menangani hal ini, dimana Mutex jenis ini memiliki dua lock, yakni Lock untuk Read dan Lock untuk Write

type BankAccount struct {
	RWMutex sync.RWMutex
	Balance int
}

func (account *BankAccount) AddBalance(amount int) {
	account.RWMutex.Lock()
	account.Balance = account.Balance + amount
	account.RWMutex.Unlock()
}

func (account *BankAccount) GetBalance() int {
	account.RWMutex.RLock()
	balance := account.Balance
	account.RWMutex.RUnlock()
	return balance
}

func (BankAccount) Baru(duit int) BankAccount {
	return BankAccount{
		RWMutex: sync.RWMutex{},
		Balance: duit,
	}
}

func TestReadWriteMutex(t *testing.T) {
	var akun BankAccount = BankAccount{}.Baru(20000)

	for i := 0; i < 100; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				akun.AddBalance(1)
				fmt.Println(akun.GetBalance())
			}
		}()
	}

	time.Sleep(20 * time.Second)
	fmt.Println("duit: ", akun.GetBalance())
}

// DeadLock
// Hati Hati saat kita membuat aplikasi yang parallel atau conccurent, masalah yang sering kita hadapi adalah deadlock
// Deadlock adalah keadaan dimana sebuah proses goroutine saling menunggu lock sehingga tidak ada satupun goroutine yang bisa jalan
// Sekarang kita coba simulasikan proses deadlock

type UserBalance struct {
	sync.Mutex
	Name    string
	Balance int
}

func (u *UserBalance) Lock() {
	u.Mutex.Lock()
}

func (u *UserBalance) Unlock() {
	u.Mutex.Unlock()
}

func (u *UserBalance) Change(amount int) {
	u.Balance = u.Balance + amount
}

func Transfer(u1, u2 *UserBalance, amount int) {
	u1.Lock()
	fmt.Println("Lock", u1.Name)
	u1.Change(-amount)

	time.Sleep(1 * time.Second)

	u2.Lock()
	fmt.Println("Lock", u2.Name)
	u2.Change(amount)

	time.Sleep(1 * time.Second)

	u1.Unlock()
	u2.Unlock()
}

func TestDeadlock(t *testing.T) {
	user1 := UserBalance{
		Name:    "Faiz",
		Balance: 1000000,
	}

	user2 := UserBalance{
		Name:    "Juki",
		Balance: 1000000,
	}

	go Transfer(&user1, &user2, 10000)
	go Transfer(&user2, &user1, 10000)

	time.Sleep(10 * time.Second)

	fmt.Println("User ", user1.Name, "Balance ", user1.Balance)
	fmt.Println("User ", user2.Name, "Balance ", user2.Balance)
}

func TestD(t *testing.T) {
	var data chan string = make(chan string)

	go func(d chan string) {
		fmt.Println("Anjay", <-data)
	}(data)
	data <- "MAKAN"

	close(data)
}

// Jadi begini channel harus diterima oleh thread / proses asinkron lain, ketika kita membuat unbuffered channel maka
// akan terus menunggu samapai proses asinkron lain mengambil datanya

// sync.WaitGroup
// WaitGroup adalah fitur yang bisa digunakan untuk menunggu sebuah proses selesai dilakukan
// Hal ini kadang diperlukan, misal ktia ingin menjalankan beberapa proses menggunakan goroutine,
// tapi kita ingin semua proses selesai terlebih dahulu sebelum aplikasi kita selesai
// Kasus seperti ini bisa menggunakan WaitGroup
// Untuk menandai bahwa ada proses goroutine, kita bisa menggunakan method Add(int), setelah proses goroutine selesai, kita bisa gunakan method Done()
// Untuk menunggu semua proses selesai, kita bisa menggunakan method Wait()

func RunAsynchronous(group *sync.WaitGroup, count int, d chan<- string) {
	defer group.Done()

	group.Add(1)

	d <- "Halo Ke: " + strconv.Itoa(count) + "dari channel"
	time.Sleep(1 * time.Second)
}

func TestWaitGroup(t *testing.T) {
	w := sync.WaitGroup{}
	var data chan string = make(chan string)
	go func(d chan string) {
		for message := range d {
			fmt.Println(message)
		}
	}(data)

	for i := 0; i < 10; i++ {
		go RunAsynchronous(&w, i, data)
	}

	w.Wait()
	fmt.Println("Selesai")
}

// sync.Once
// Once adalah fitur di golang yang bisa kita gunakan untuk memastikan bahwa sebuah function di eksekusi hanya sekali
// Jadi berapa banyak pun goroutine yang mengakses, bisa dipastikan bahwa goroutine yang pertama yang bisa mengeksekusi function nya
// Goroutine yang lain akan di hiraukan, artinya function tidak akan dieksekusi Lagi

var counter = 0

func TestOnce(t *testing.T) {
	var once sync.Once
	var group sync.WaitGroup

	for i := 0; i < 100; i++ {
		go func() {
			group.Add(1)
			once.Do(func() {
				counter++
			})
			group.Done()
		}()
	}

	group.Wait()
	fmt.Print("Counter: ", counter)
}

// sync.Pool
// Pool adalah implementasi design pattern bernama object pool pattern
// Sederhananya, design pattern Pool ini digunakan untuk menyimpan data, selanjutnya untuk menggunakan datanya, kita bisa mengambil data dari pool, dan setelah selesai menggunakan datanya, kita bisa menyimpan kembali ke pool nya
// Implementasi Pool di Go lang ini sudah aman dari problem race condition
// Biasanya pool ini sering digunakan di kasus seperti koneksi ke database
// Pada defaultnya jika data kosong di TestPool maka akan diberikan value <nil>, kita bisa mengubah nya
// dengan cara menambahkan properti New pada object pool dengan fungsi yang meretrun interface kosong
// kita bebas boleh mengisi datanya dengan apapun

func TestPool(t *testing.T) {
	var pool sync.Pool = sync.Pool{
		New: func() interface{} {
			return "kosong"
		},
	}

	pool.Put("Faiz")
	pool.Put("Anancuy")
	pool.Put("Mantap")

	for i := 0; i < 10; i++ {
		go func() {
			data := pool.Get()
			switch v := any(data).(type) {
			case string:
				if v == "kosong" {
					fmt.Println("Kosong cuk")
				}
			}
			fmt.Println(data)
			time.Sleep(1)
			pool.Put(data)
		}()
	}

	time.Sleep(5 * time.Second)
}

// sync.Map
// Go lang memiliki sebuah struct bernama sync.Map
// Map ini miripGo Lang map, namun yang membedakan, Map ini aman untuk menggunakan conccurent menggunakan goroutine
// Ada beberapa method yan bisa kita gunakan di Map:
// Store(key): untuk menyimpan data ke dalam Map
// Load(key): untuk mengambil data dari Map menggunakan key
// Delete(key): untuk menghapus data di Map menggunaan key
// Range(function(key, value)) digunakan untuk melakukan iterasi seluruh data di Map

func TestSyncMap(t *testing.T) {
	var data sync.Map
	var addToMap = func(val int) {
		data.Store(val, val)
	}

	for i := 0; i < 10; i++ {
		go addToMap(i)
	}

	time.Sleep(3 * time.Second)
	data.Range(func(key, val interface{}) bool {
		fmt.Println(key, ":", val)
		return true
	})
}
