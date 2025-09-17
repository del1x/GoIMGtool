# GoIMGtool 🖼️
<a name="goimgtool"></a>

**GoIMGtool** is a lightweight image processing tool that allows resizing images and applying watermarks.
It supports batch processing, multiple formats, and optimal web image sizes.

[Перейти на русскую версию](#goimgtool-русская-версия)

### Features

* 🔍 **Resize** – Users can choose the output size. Images are only limited by the watermark size.
* 💧 **Watermark** – Semi-transparent watermark overlay with two modes: `crop` and `resize`.
* 🚀 **Format support** – PNG, JPG, JPEG, WebP.
* 🐳 **Docker-ready** – Can run via Docker or Docker Compose.
* ⚡ **Optimized for Web** – Watermarked images <= 100KB.
* 🖥️ **GUI-based** – Users select image folders and watermark via GUI (Fyne).

### Requirements

* Go 1.24.6
* MSYS2 with `mingw-w64-x86_64-libwebp` (for WebP on Windows)
* X Server (Windows) for GUI interaction (if using Go run version)

---

## Quick Start (Windows Executable)

If you just want to run the application on Windows, simply go to the `dist/` folder and launch the executable.

---

## Installation (Native)

1. Install Go 1.24.6: [https://golang.org/dl/](https://golang.org/dl/)
2. Install MSYS2: [https://www.msys2.org/](https://www.msys2.org/) and run:

   ```bash
   pacman -S mingw-w64-x86_64-libwebp
   ```
3. Set CGO environment variables (Windows CMD, adjust the path to your MSYS2 installation):

   ```bash
   set CGO_CFLAGS=-I<path_to_msys2>\mingw64\include
   set CGO_LDFLAGS=-L<path_to_msys2>\mingw64\lib -lwebp
   set CGO_ENABLED=1
   set PATH=<path_to_msys2>\mingw64\bin;%PATH%
   ```
4. Clone repository and initialize Go modules:

   ```bash
   git clone https://github.com/del1x/GoIMGtool.git
   cd GoIMGtool
   go mod tidy
   ```

---

## Usage (Native)

1. Run:

   ```bash
   go run -tags "desktop" main.go
   ```
2. Using the GUI, select the folder with images and the watermark file.
3. Choose output size, format (webp/png), and watermark mode (`crop` or `resize`).
   Processed files appear in `Images_watermarked/` and are optimized to <= 100KB.

---

## Docker Usage

### Build Docker Image

```bash
docker build -t goimgtool:latest .
```

### Run Docker Container

```bash
docker run --rm -e DISPLAY=host.docker.internal:0 -v /path/to/images:/app/Images goimgtool:latest
```

### Docker Compose (Development)

```bash
docker-compose -f docker-compose.yml up --build
```

> **Note:** X Server must be running on Windows for GUI.

---

## Notes

* Ensure the watermark file is valid.
* WebP support requires CGO settings on Windows.
* Users choose the image folder and watermark file via GUI.
* Two watermark modes available: `crop` and `resize`.
* Optimized for web: output images <= 100KB.

---

# GoIMGtool 🖼️ (Русская версия)

<a name="goimgtool-русская-версия"></a>
[Перейти на английскую версию](#goimgtool-🖼️)

**GoIMGtool** — инструмент для обработки изображений: изменение размера и наложение водяных знаков.
Поддерживает пакетную обработку, несколько форматов и оптимальные размеры изображений для веба.

### Возможности

* 🔍 **Изменение размера** — пользователь выбирает размер; ограничение только по размеру водяного знака.
* 💧 **Водяной знак** — полупрозрачный, два режима наложения: `crop` и `resize`.
* 🚀 **Поддержка форматов** — PNG, JPG, JPEG, WebP.
* 🐳 **Docker-ready** — можно запускать через Docker или Docker Compose.
* ⚡ **Оптимизация для веб** — водяные изображения <= 100KB.
* 🖥️ **GUI** — выбор папки с изображениями и водяного знака пользователем (Fyne).

### Требования

* Go 1.24.6
* MSYS2 с `mingw-w64-x86_64-libwebp` (для WebP на Windows)
* X Server (Windows) для GUI (если используется версия через Go run)

---

## Быстрый запуск (Windows Executable)

Если нужно просто запустить приложение на Windows, достаточно открыть папку `dist/` и запустить исполняемый файл.

---

## Установка (Native)

1. Установите Go 1.24.6: [https://golang.org/dl/](https://golang.org/dl/)
2. Установите MSYS2: [https://www.msys2.org/](https://www.msys2.org/) и выполните:

   ```bash
   pacman -S mingw-w64-x86_64-libwebp
   ```
3. Установите переменные окружения CGO (Windows CMD, путь к вашей установке MSYS2):

   ```bash
   set CGO_CFLAGS=-I<path_to_msys2>\mingw64\include
   set CGO_LDFLAGS=-L<path_to_msys2>\mingw64\lib -lwebp
   set CGO_ENABLED=1
   set PATH=<path_to_msys2>\mingw64\bin;%PATH%
   ```
4. Клонируйте репозиторий и установите зависимости Go:

   ```bash
   git clone https://github.com/del1x/GoIMGtool.git
   cd GoIMGtool
   go mod tidy
   ```

---

## Использование (Native)

1. Запустите:

   ```bash
   go run -tags "desktop" main.go
   ```
2. Через GUI выберите папку с изображениями и файл водяного знака.
3. Выберите размер вывода, формат (webp/png) и режим водяного знака (`crop` или `resize`).
   Обработанные файлы появятся в `Images_watermarked/`, размер <= 100KB.

---

## Docker Использование

### Сборка Docker Image

```bash
docker build -t goimgtool:latest .
```

### Запуск Docker Container

```bash
docker run --rm -e DISPLAY=host.docker.internal:0 -v /path/to/images:/app/Images goimgtool:latest
```

### Docker Compose (Разработка)

```bash
docker-compose -f docker-compose.yml up --build
```

> **Примечание:** На Windows должен быть запущен X Server для GUI.

---

## Заметки

* Убедитесь, что файл водяного знака корректный.
* Для поддержки WebP на Windows нужны настройки CGO.
* Пользователь выбирает папку с изображениями и файл водяного знака через GUI.
* Доступны два режима водяного знака: `crop` и `resize`.
* Оптимизация для веб: итоговые изображения <= 100KB.
