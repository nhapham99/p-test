package googleplay

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"payment-module/configs"
	"payment-module/logger"

	"golang.org/x/oauth2/google"
)

func CheckGoogleAuth(purchaseToken string) (*interface{}, bool, error) {
	data, err := os.ReadFile("55300f006edc.json")
	if err != nil {
		logger.Error("ReadFile error 17: ", err)
		return nil, false, err
	}

	conf, err := google.JWTConfigFromJSON(data, "https://www.googleapis.com/auth/androidpublisher")
	if err != nil {
		logger.Error("JWTConfigFromJSON error 23: ", err)
		return nil, false, err
	}

	tokenSource := conf.TokenSource(context.Background())
	_, err = tokenSource.Token()
	if err != nil {
		logger.Error("TokenSource error 30: ", err)
		return nil, false, err
	}

	url := "https://androidpublisher.googleapis.com/androidpublisher/v3/applications/" + configs.GetEnv("GOOGLE_PACKAGE_NAME") + "/purchases/subscriptionsv2/tokens/" + purchaseToken
	client := conf.Client(context.Background())
	resp, err := client.Get(url)
	if err != nil {
		logger.Error("client.Get error 38: ", err)
		return nil, false, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("ReadAll error 45: ", err)
		return nil, false, err
	}

	var result interface{}
	json.Unmarshal(body, &result)

	if resp.StatusCode == 200 {
		logger.Info("CheckGoogleAuth success")
		return &result, true, nil
	}

	logger.Error("CheckGoogleAuth error 56: ", result)

	return nil, false, nil
}

func SubscriptionAcknowledge(token string, subscriptionId string) error {
	data, err := os.ReadFile("55300f006edc.json")
	if err != nil {
		return err
	}

	conf, err := google.JWTConfigFromJSON(data, "https://www.googleapis.com/auth/androidpublisher")
	if err != nil {
		return err
	}
	url := "https://androidpublisher.googleapis.com/androidpublisher/v3/applications/" + configs.GetEnv("GOOGLE_PACKAGE_NAME") + "/purchases/subscriptions/" + subscriptionId + "/tokens/" + token + ":acknowledge"
	client := conf.Client(context.Background())

	resp, err := client.Post(url, "application/json", nil)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return nil
}
