# T1Disk Go Library

T1Disk Go Library предоставляет удобный интерфейс для загрузки файлов в T1 Диск через API. Библиотека позволяет получить токен авторизации, загрузить файлы в объектное хранилище S3 и подтвердить успешную загрузку.


Использование

### Шаг 1: Авторизация

Перед тем как загружать файлы, вам нужно получить токен авторизации. Для этого выполните функцию `Login()`.

t1disk := t1disk.NewT1Disk("your-login", "your-password", true)
if err := t1disk.Login(); err != nil {
	fmt.Println("Ошибка при авторизации:", err)
	return
}
### Шаг 2: Загрузка файла

Чтобы загрузить файл, используйте функцию `UploadToT1Disk()`. Она выполняет все шаги, включая:

* Получение ссылки для загрузки в S3,
* Загрузку файла,
* Подтверждение успешной загрузки.

Пример загрузки файла:

filePath := "/path/to/your/file.mp4"      // Путь к вашему файлу на локальной машине
pathOnT1Disk := "T1 Диск/таймер/10 hour timer.mp4"  // Путь, куда будет загружен файл на T1 Диск

if err := t1disk.UploadToT1Disk(pathOnT1Disk, filePath, false); err != nil {
	fmt.Println("Ошибка при загрузке файла:", err)
}

### Описание API

1. **`NewT1Disk(login string, password string, permanentAuth bool) *T1Disk`** : Конструктор для создания экземпляра библиотеки. Принимает логин, пароль и флаг, указывающий на постоянную авторизацию.
2. **`Login() error`** : Выполняет запрос на получение токена авторизации. В случае успешного входа токен будет сохранен внутри структуры.
3. **`GetUploadURL(path string, multipart bool) (string, string, string, error)`** : Запрашивает ссылку для загрузки файла в объектное хранилище S3. Возвращает URL для загрузки, `Content-Type`, и `confirm_url`.
4. **`UploadFile(uploadURL string, contentType string, filePath string) error`** : Загружает файл по полученному `uploadURL`. Использует `Content-Type` для правильной загрузки и токен авторизации.
5. **`ConfirmUpload(confirmURL string) error`** : Подтверждает успешную загрузку файла с помощью запроса к `confirm_url`.
6. **`UploadToT1Disk(path string, filePath string, multipart bool) error`** : Полный процесс загрузки файла, который включает получение URL, загрузку файла и подтверждение загрузки.

### Пример использования

package main

import (
	"fmt"
	"t1disk-go-library/t1disk"
)

func main() {
	// Создаем новый экземпляр клиента
	t1disk := t1disk.NewT1Disk("your-login", "your-password", true)

	// Шаг 1: Авторизация
	if err := t1disk.Login(); err != nil {
		fmt.Println("Ошибка при авторизации:", err)
		return
	}

	// Шаг 2: Загрузка файла
	filePath := "/path/to/your/file.mp4"  // Локальный путь к файлу
	pathOnT1Disk := "T1 Диск/таймер/10 hour timer.mp4"  // Путь на T1 Диске

	if err := t1disk.UploadToT1Disk(pathOnT1Disk, filePath, false); err != nil {
		fmt.Println("Ошибка при загрузке файла:", err)
		return
	}

	fmt.Println("Файл успешно загружен!")
}
