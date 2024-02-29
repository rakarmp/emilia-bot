<p align="center">
<img src="picture.jpg" src="emilia"/>

<h3 align="center">Emilia AI - ChatGPT Bot Telegram</h3>

<p align="center">Ini adalah bot Telegram yang dibuat dengan ChatGPT dan Golang. Bot ini menggunakan model bahasa GPT-3 Turbo OpenAI untuk menghasilkan respons pesan pengguna secara real-time.</p>
</p>

## Prerequisites
Untuk menjalankan bot ini, Anda harus memiliki perangkat lunak berikut yang terinstal di sistem Anda:

- Go language (versi terbaru)

## Features
- Menghasilkan respons seperti manusia terhadap pesan pengguna menggunakan ChatGPT API
- Menyimpan pesan pengguna dengan sqlite
- Dukungan Telegram
- Dibangun dengan Go untuk kinerja yang cepat dan efisien?

Sebelum Anda dapat menggunakan bot, Anda harus membuat bot Telegram menggunakan [kerangka kerja BotFather](https://t.me/botfather). Setelah Anda membuat bot dan mendapatkan token API, Anda juga memerlukan [API key dari OpenAI](https://platform.openai.com/account/api-keys)

Copy .env.example dengan perintah berikut
```sh
mv .env.example .env

# Atau

cp .env.example .env
```

Ini adalah contoh file `.env`
```.env
TELEGRAM_API_KEY=""
OPENAI_TOKEN=""
RETAIN_HISTORY="false"
```

`RETAIN_HISTORY="true"` mengirimkan percakapan sebelumnya dengan teks saat ini, [lihat di sini](https://platform.openai.com/docs/guides/chat/introduction), tetapi jika false, ini hanya mengirimkan prompt + teks pengguna saat ini, hal ini mengurangi jumlah token yang dikirim per permintaan.

membuat `prompt.txt` atau mengganti nama file contoh

```sh
$ mv prompt.example.txt prompt.txt

# Atau

$ cp prompt.example.txt prompt.txt
```
prompt membantu Anda menyesuaikan bagaimana bot akan bereaksi terhadap pesan

## Installing
Pertama, clone repositori ini:

```sh
$ git clone https://github.com/rakarmp/emilia-bot.git
```

Lalu, arahkan ke direktori project:

```sh
$ cd emilia-bot
```

Terakhir Bangun Project Dan Run:

```sh
$ go build -o file_name

$ ./file_name
```

## License
This project is licensed under the MIT License. See the [LICENSE](https://github.com/rakarmp/emilia-bot/LICENSE) file for details

## Resources
- [Go Documentation](https://golang.org/doc/)
- [Telegram Bot API](https://core.telegram.org/bots/api)                                                                                            