package external

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type External struct {
	bot *tgbotapi.BotAPI
}

type botGetter interface {
	GetBot() *tgbotapi.BotAPI
}

func New(tgclient botGetter) *External {
	return &External{
		bot: tgclient.GetBot(),
	}
}

func (e *External) IsUserInChat(chatID, userID int64) (bool, error) {
	member, err := e.bot.GetChatMember(tgbotapi.GetChatMemberConfig{
		ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
			ChatID: chatID,
			UserID: userID,
		},
	})
	if err != nil {
		return false, err
	}

	if member.Status == "member" || member.IsCreator() || member.IsAdministrator() {
		return true, nil
	}

	return false, nil
}
