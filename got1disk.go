package got1disk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"go.uber.org/zap"
)

type T1DiskConstructor struct {
	permanentAuth bool
	logger        *zap.Logger
	BaseUrlT1     string
}

func NewT1Disk(permanentAuth bool, logger *zap.Logger, BaseUrlT1 string) *T1DiskConstructor {
	return &T1DiskConstructor{
		permanentAuth: permanentAuth,
		logger:        logger,
		BaseUrlT1:     BaseUrlT1,
	}
}

func (t *T1DiskConstructor) Login(login, password string) (string, error) {
	fmt.Println(login)
	fmt.Println(password)

	url := fmt.Sprintf("%s/accounts/login/?login=email:%s&password=%s&permanent_auth=%t", t.BaseUrlT1, login, password, t.permanentAuth)
	fmt.Println(url)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return "", err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to login: %s", resp.Status)
	}

	var result struct {
		Token             string `json:"token"`
		DeviceDescription string `json:"device_description"`
		Created           int64  `json:"created"`
		Expires           int64  `json:"expires"`
		RemoteWipe        bool   `json:"remote_wipe"`
		ID                string `json:"id"`
		UserID            int    `json:"userid"`
		UserEID           int    `json:"user_eid"`
		Login             string `json:"login"`
		Domain            string `json:"domain"`
		Name              string `json:"name"`
		OfferURL          string `json:"offer_url"`
		CompanyID         int    `json:"company_id"`
		PreviousLoginDate int64  `json:"previous_login_date"`
		PreviousLoginIP   string `json:"previous_login_ip"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Token, nil
}

func (t *T1DiskConstructor) GetUploadURL(path, baseURL, token string, multipart bool) (string, string, string, error) {
	url := fmt.Sprintf("%s/files/create/?path=%s&multipart=%t", t.BaseUrlT1, path, multipart)
	fmt.Println(path)
	req, _ := http.NewRequest("POST", url, nil)
	req.Header.Set("Mountbit-Auth", token)
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

	var result struct {
		UploadURL  string            `json:"url"`
		Headers    map[string]string `json:"headers"`
		ConfirmURL string            `json:"confirm_url"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", "", "", err
	}

	contentType, ok := result.Headers["Content-Type"]
	if !ok {
		contentType = "application/octet-stream"
	}

	return result.UploadURL, contentType, result.ConfirmURL, nil
}

func (t *T1DiskConstructor) UploadFile(uploadURL, contentType, filePath, token string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %v", err)
	}

	req, err := http.NewRequest("PUT", uploadURL, file)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Mountbit-Auth", token)
	req.ContentLength = fileInfo.Size()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to upload file: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to upload file: %s", resp.Status)
	}

	return nil
}

func (t *T1DiskConstructor) ConfirmUpload(confirmURL, token string) error {
	req, err := http.NewRequest("POST", confirmURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Mountbit-Auth", token)

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

func (t *T1DiskConstructor) UploadToT1Disk(path, baseURL, filePath, token string, multipart bool) error {
	uploadURL, contentType, confirmURL, err := t.GetUploadURL(path, baseURL, token, multipart)
	if err != nil {
		return err
	}

	if err := t.UploadFile(uploadURL, contentType, filePath, token); err != nil {
		return err
	}

	if err := t.ConfirmUpload(confirmURL, token); err != nil {
		return err
	}

	return nil
}
