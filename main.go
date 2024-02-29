package main

import (
	"context"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/joho/godotenv"
	"github.com/rakyll/openai-go"
	"github.com/rakyll/openai-go/chat"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	retainHistory bool
	promptName    = "prompt.txt"
)

func main() {
	// setup logger
	log.Logger = log.With().Caller().Logger()
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Environment files
	err := godotenv.Load()
	if err != nil {
		log.Debug().Msg(err.Error())
	}

	retainHistory = os.Getenv("RETAIN_HISTORY") == "true"

	if err := ConnectDB(); err != nil {
		log.Fatal().Msg(err.Error())
	}

	// start server
	StartServer()
}

// Males Nulis Comment YTTA Aja
func StartServer() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(handler),
	}

	b, err := bot.New(os.Getenv("TELEGRAM_API_KEY"), opts...)
	if err != nil {
		panic(err)
	}

	log.Debug().Msg("Telegram bot started!")
	b.Start(ctx)
}

func SendToChatGPT(chatId, textMsg string) []*chat.Choice {
	var (
		ctx = context.Background()
		s   = openai.NewSession(os.Getenv("OPENAI_TOKEN"))

		gptMsgs = make([]*chat.Message, 0)
	)

	prevMessages, err := FindMessages(chatId)
	if err != nil {
		log.Err(err)
	}

	prmptB, _ := os.ReadFile(promptName)

	if len(prevMessages) == 0 {
		log.Debug().Msg("added system prompt because its a first time user")
		gptMsgs = append(gptMsgs, &chat.Message{
			Role:    "user", // "system"
			Content: string(prmptB),
		})

	} else {
		if retainHistory {
			for _, prevMsg := range prevMessages {
				gptMsgs = append(gptMsgs, &chat.Message{
					Role:    prevMsg.Role,
					Content: prevMsg.Content,
				})
			}
		} else {
			gptMsgs = append(gptMsgs, &chat.Message{
				Role:    "user", // "system"
				Content: string(prmptB),
			})
		}
	}

	// add this current message
	gptMsgs = append(gptMsgs, &chat.Message{
		Role:    "user",
		Content: textMsg,
	})

	// process request
	client := chat.NewClient(s, "gpt-3.5-turbo-0301")
	resp, err := client.CreateCompletion(ctx, &chat.CreateCompletionParams{
		Messages: gptMsgs,
	})
	if err != nil {
		log.Error().Msgf("Failed to complete: %v", err)
		return nil
	}

	if len(prevMessages) == 0 {
		for _, gptMsg := range gptMsgs {
			_, err := CreateMessage(Message{
				ChatID:  chatId,
				Content: gptMsg.Content,
				Role:    gptMsg.Role,

				PromptTokens:     resp.Usage.PromptTokens,
				CompletionTokens: resp.Usage.CompletionTokens,
				TotalTokens:      resp.Usage.TotalTokens,
			})
			if err != nil {
				log.Error().Msgf("unable to save message: %v", err)
			}
		}
	} else {
		_, err := CreateMessage(Message{
			ChatID:  chatId,
			Role:    "user",
			Content: textMsg,

			PromptTokens:     resp.Usage.PromptTokens,
			CompletionTokens: resp.Usage.CompletionTokens,
			TotalTokens:      resp.Usage.TotalTokens,
		})
		if err != nil {
			log.Error().Msgf("unable to current message: %v", err)
		}
	}

	for _, choice := range resp.Choices {
		_, err := CreateMessage(Message{
			ChatID:  chatId,
			Role:    choice.Message.Role,
			Content: choice.Message.Content,

			PromptTokens:     resp.Usage.PromptTokens,
			CompletionTokens: resp.Usage.CompletionTokens,
			TotalTokens:      resp.Usage.TotalTokens,
		})
		if err != nil {
			log.Error().Msgf("unable save chat response: %v", err)
		}
	}

	log.Info().
		Int("TotalTokens", resp.Usage.TotalTokens).
		Int("CompletionTokens", resp.Usage.CompletionTokens).
		Int("PromptTokens", resp.Usage.PromptTokens).
		Msg("usage")

	return resp.Choices
}

// handler
func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	// Initialize random seed
	rand.Seed(time.Now().UnixNano())

	outgoingMsg := update.Message.Text
	chatId := update.Message.Chat.ID
	log.Debug().Msg(outgoingMsg)

	chatIdStr := strconv.Itoa(int(chatId))
	chatResp := SendToChatGPT(chatIdStr, outgoingMsg)
	if chatResp == nil {

		// Define an array of responses
		responses := []string{
			"Maaf, ada masalah sebentar. Aku bakal terus coba dan kabarin kamu begitu udah online lagi.",
			"Hmm, ada yang kurang beres nih. Aku lagi cek dan bakal update kamu kalo udah berfungsi lagi.",
			"Kayaknya aku lagi ngalamin sedikit masalah. Aku lagi pantauin dan bakal kabarin kamu kalo udah balik normal.",
			"Waduh, aku lagi down nih. Aku bakal usahain buat nyambung lagi dan kabarin kamu.",
			"Nggak bisa nih, kayaknya aku nggak bisa nyampe ke tujuan akhir. Tapi aku bakal kabarin kamu kalo udah online lagi.",
			"Oh no, aku lagi down. Aku bakal terus coba dan kabarin kamu kalo aku udah online lagi.",
		}
		randIndex := rand.Intn(len(responses))

		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   responses[randIndex],
		})
		return
	}

	for _, choice := range chatResp {
		incomingMsg := choice.Message
		log.Printf("role=%q, content=%q", incomingMsg.Role, incomingMsg.Content)

		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   incomingMsg.Content,
		})
	}
}
