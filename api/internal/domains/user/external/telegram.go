package external

import (
	"io"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type External struct {
	bot *tgbotapi.BotAPI
}

func New(bot *tgbotapi.BotAPI) *External {
	return &External{
		bot: bot,
	}
}

// GetAvatar returns the avatar of the user with the given ID as a byte slice.
func (t *External) GetAvatar(userID int64) ([]byte, error) {
	photos, err := t.bot.GetUserProfilePhotos(tgbotapi.UserProfilePhotosConfig{
		UserID: userID,
	})

	if err != nil {
		return nil, err
	}

	if len(photos.Photos) > 0 {
		// Get the first photo (first one is the smallest)
		photo := photos.Photos[0][0]

		// Get the file
		fileURL, err := t.bot.GetFileDirectURL(photo.FileID)

		if err != nil {
			return nil, err
		}

		response, err := http.Get(fileURL)
		if err != nil {
			return nil, err
		}
		defer response.Body.Close()

		data, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}

		return data, nil

	}
	return nil, nil
}
