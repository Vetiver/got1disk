package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type T1DiskConstructor struct {
	login         string
	password      string
	token         string
	baseURL       string
	permanentAuth bool
}

func NewT1Disk(login, password string, permanentAuth bool, baseURL string) *T1DiskConstructor {
	return &T1DiskConstructor{
		login:         login,
		password:      password,
		baseURL:       baseURL,
		permanentAuth: permanentAuth,
	}
}

func (t *T1DiskConstructor) Login() error {
	url := fmt.Sprintf("%s/accounts/login/", t.baseURL)
	data := map[string]interface{}{
		"login":          t.login,
		"password":       t.password,
		"permanent_auth": t.permanentAuth,
	}

	body, _ := json.Marshal(data)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to login: %s", resp.Status)
	}

	var result map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	t.token = result["token"]
	return nil
}

func (t *T1DiskConstructor) GetUploadURL(path string, multipart bool) (string, string, string, error) {
	url := fmt.Sprintf("%s/files/create/", t.baseURL)
	data := map[string]interface{}{
		"path":      path,
		"multipart": multipart,
	}

	body, _ := json.Marshal(data)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	req.Header.Set("Mountbit-Auth", t.token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", "", fmt.Errorf("failed to get upload URL: %s", resp.Status)
	}

	var result map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", "", "", err
	}

	return result["url"], result["Content-Type"], result["confirm_url"], nil
}

func (t *T1DiskConstructor) UploadFile(uploadURL, contentType, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	req, err := http.NewRequest("PUT", uploadURL, file)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Mountbit-Auth", t.token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to upload file: %s", resp.Status)
	}

	return nil
}

func (t *T1DiskConstructor) ConfirmUpload(confirmURL string) error {
	req, err := http.NewRequest("POST", confirmURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Mountbit-Auth", t.token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to confirm upload: %s", resp.Status)
	}

	return nil
}

func (t *T1DiskConstructor) UploadToT1Disk(path, filePath string, multipart bool) error {
	uploadURL, contentType, confirmURL, err := t.GetUploadURL(path, multipart)
	if err != nil {
		return err
	}

	if err := t.UploadFile(uploadURL, contentType, filePath); err != nil {
		return err
	}

	if err := t.ConfirmUpload(confirmURL); err != nil {
		return err
	}

	fmt.Println("Файл успешно загружен!")
	return nil
}
