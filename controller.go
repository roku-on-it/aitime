package main

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
	"strconv"
)

func MountController(router *fiber.App) {
	router.Post("/handle-prompt", HandlePrompt)
}

func HandlePrompt(ctx *fiber.Ctx) error {
	var body HandlePromptReqBody

	if err := json.Unmarshal(ctx.Body(), &body); err != nil {
		return ctx.Status(400).JSON(fiber.Map{"text_to_speech": err.Error()})
	}

	promptLen := len(body.Prompt)

	fmt.Println(body.Prompt)

	if promptLen < 3 || promptLen > 90 {
		return ctx.Status(400).JSON(fiber.Map{"text_to_speech": "prompt parameter must be between 3 and 90 characters"})
	}

	params, err := GetParametersFromAI(body.Prompt)

	if err != nil {
		log.Println(err)
		return ctx.Status(500).JSON(fiber.Map{"text_to_speech": err.Error()})
	}

	switch params.Intent {
	case WhereWasI:
		season, episode, err := GetNextEpisodeForShow(params.Show)
		if err != nil {
			return ctx.Status(500).JSON(fiber.Map{"text_to_speech": "something went wrong trying to get next episode for " + params.Show})
		}

		tts := fmt.Sprintf("You're on season %d, episode %d, in %s", season, episode, params.Show)
		return ctx.Status(200).JSON(fiber.Map{
			"text_to_speech": tts,
		})

	case AddToWatched:
		seasonInt, err := strconv.Atoi(params.Season)
		if err != nil {
			return ctx.Status(500).JSON(fiber.Map{"text_to_speech": "could not convert season string to integer for " + params.Show})
		}
		epInt, err := strconv.Atoi(params.Episode)
		if err != nil {
			return ctx.Status(500).JSON(fiber.Map{"text_to_speech": "could not convert episode string to integer for " + params.Show})
		}

		err = MarkEpisode(params.Show, seasonInt, epInt, true)
		if err != nil {
			return ctx.Status(500).JSON(fiber.Map{"text_to_speech": "something went wrong marking episode for " + params.Show})
		}
		tts := fmt.Sprintf("Marked season %s, episode %s, off %s, as watched", params.Season, params.Episode, params.Show)
		return ctx.Status(200).JSON(fiber.Map{
			"text_to_speech": tts,
		})

	case RemoveFromWatched:
		seasonInt, err := strconv.Atoi(params.Season)
		if err != nil {
			return ctx.Status(500).JSON(fiber.Map{"text_to_speech": "could not convert season string to integer for " + params.Show})
		}
		epInt, err := strconv.Atoi(params.Episode)
		if err != nil {
			return ctx.Status(500).JSON(fiber.Map{"text_to_speech": "could not convert episode string to integer for " + params.Show})
		}

		err = MarkEpisode(params.Show, seasonInt, epInt, false)
		if err != nil {
			return ctx.Status(500).JSON(fiber.Map{"text_to_speech": "something went wrong marking episode for " + params.Show})
		}

		tts := fmt.Sprintf("Successfully removed season %s, episode %s, off %s", params.Season, params.Episode, params.Show)
		return ctx.Status(200).JSON(fiber.Map{
			"text_to_speech": tts,
		})
	}

	fmt.Println("Intent:", params.Intent)
	fmt.Println("Show:", params.Show)
	fmt.Println("Season:", params.Season)
	fmt.Println("Episode:", params.Episode)

	return nil
}
