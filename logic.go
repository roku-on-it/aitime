package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"log"
	"slices"
	"strconv"
	"strings"
	"time"
)

var client = resty.New().SetTimeout(time.Second * 15)

type TVSeriesParams struct {
	Intent  string
	Show    string
	Season  string
	Episode string
}

const (
	AddToWatched      = "ADD_TO_WATCHED"
	RemoveFromWatched = "REMOVE_FROM_WATCHED"
	WhereWasI         = "WHERE_WAS_I"

	TvTimeUserID = "XXXXXXXX"
	TVTimeToken  = "XXXXXXXX"
)

var retryCount = 0

func GetParametersFromAI(prompt string) (*TVSeriesParams, error) {
	if retryCount >= 3 {
		retryCount = 0

		log.Println("Reached maximum retry count. Killing process")

		return nil, errors.New("could not get a proper answer from AI")
	}

	var result ChatResponse

	r, err := client.R().SetBody(
		ChatReq{Model: "codellama",
			Messages: []RequestMsg{
				{Role: "user", Content: "You are an agent that can parse the TV Show names and understand the intent. Extract the parameters from sentence and answer with only one of the possible answers. Possible intents: ADD_TO_WATCHED (when user wants to add some episode to watched before list), REMOVE_FROM_WATCHED (when user wants to remove some episode to watched before list), WHERE_WAS_I (when user wants to remember the last episide they watched. which episode they're currently at) Possible answer formats: ADD_TO_WATCHED:<seriesName>:<season(just number)>:<episode(just number)> REMOVE_FROM_WATCHED:<seriesName>:<season(just number)>:<episode(just number)> WHERE_WAS_I:<seriesName>"},
				{Role: "assistant", Content: "Sure, I can do that! Please provide the sentence you would like me to parse."},
				{Role: "user", Content: prompt},
			}},
	).SetResult(&result).Post("http://localhost:11434/api/chat")

	if err != nil {
		if strings.Contains(err.Error(), "deadline exceeded") {
			fmt.Println("Request timeout after 15 seconds")
			return nil, context.DeadlineExceeded
		}

		panic(err)
	}

	if r.StatusCode() != 200 {
		fmt.Println("error trying to get parameters from ai:")
		fmt.Println(string(r.Body()))
	}

	splt := strings.Split(result.Message.Content, ":")
	intent := splt[0]
	params := &TVSeriesParams{
		Intent: intent,
	}

	retry := func() (*TVSeriesParams, error) {
		retryCount++

		return GetParametersFromAI(prompt)
	}

	switch intent {
	case AddToWatched:
		if len(splt) != 4 {
			return retry()
		}
		params.Show = splt[1]
		params.Season = splt[2]
		params.Episode = splt[3]
	case RemoveFromWatched:
		if len(splt) != 4 {
			return retry()
		}
		params.Show = splt[1]
		params.Season = splt[2]
		params.Episode = splt[3]
	case WhereWasI:
		if len(splt) != 2 {
			return retry()
		}
		params.Show = splt[1]
	default:
		return retry()
	}

	return params, nil
}

// GetNextEpisodeForShow returns the next episode to watch for the given show.
// Return values represent Season and Episode number respectively
func GetNextEpisodeForShow(show string) (int, int, error) {
	var resArr []WatchNextResponse

	r, err := client.R().SetResult(&resArr).Get("https://api2.tozelabs.com/v2/user/" + TvTimeUserID + "/to_watch?limit=500")
	if err != nil {
		if strings.Contains(err.Error(), "deadline exceeded") {
			fmt.Println("Request timeout after 15 seconds")
			return 0, 0, context.DeadlineExceeded
		}

		panic(err)
	}

	if r.StatusCode() != 200 {
		fmt.Println("error trying to get next episode:")
		fmt.Println(string(r.Body()))
	}

	indexOfShow := slices.IndexFunc(resArr, func(w WatchNextResponse) bool {
		fmt.Println(w.Show.Name)
		return strings.Contains(StripStr(w.Show.Name), StripStr(show))
	})

	fmt.Println(indexOfShow)

	if indexOfShow == -1 {
		return 0, 0, errors.New("show could not be found")
	}

	found := resArr[indexOfShow]

	return found.Season, found.Episode, nil
}

func MarkEpisode(show string, season, episode int, watched bool) error {
	var res FavoriteShowsResponse

	r, err := client.R().SetResult(&res).Get("https://api2.tozelabs.com/v2/user/" + TvTimeUserID + "?fields=shows.fields(id,name)")
	if err != nil {
		if strings.Contains(err.Error(), "deadline exceeded") {
			fmt.Println("Request timeout after 15 seconds")
			return context.DeadlineExceeded
		}

		panic(err)
	}

	if r.StatusCode() != 200 {
		fmt.Println("error trying to mark episode:")
		fmt.Println(string(r.Body()))
	}

	indexOfShow := slices.IndexFunc(res.Shows, func(s Show) bool {
		return strings.Contains(StripStr(s.Name), StripStr(show))
	})

	if indexOfShow == -1 {
		return errors.New("show could not be found")
	}

	found := res.Shows[indexOfShow]

	var episodesRes EpisodesResponse

	r, err = client.R().
		SetResult(&episodesRes).
		Get("https://msapi.tvtime.com/v1/series/" + strconv.Itoa(found.ID) + "/episodes")
	if err != nil {
		if strings.Contains(err.Error(), "deadline exceeded") {
			fmt.Println("Request timeout after 15 seconds")
			return context.DeadlineExceeded
		}

		panic(err)
	}

	if r.StatusCode() != 200 {
		fmt.Println("error trying to mark episode:")
		fmt.Println(string(r.Body()))
	}

	toMarkIndex := slices.IndexFunc(episodesRes.Data, func(data Data) bool {
		return data.Season.Number == season && data.Number == episode
	})

	foundEp := episodesRes.Data[toMarkIndex]

	if watched {
		r, err = client.R().SetAuthToken(TVTimeToken).Post("https://api2.tozelabs.com/v2/watched_episodes/episode/" + strconv.Itoa(foundEp.ID))
	} else {
		r, err = client.R().SetAuthToken(TVTimeToken).Delete("https://api2.tozelabs.com/v2/watched_episodes/episode/" + strconv.Itoa(foundEp.ID))
	}

	if err != nil {
		if strings.Contains(err.Error(), "deadline exceeded") {
			fmt.Println("Request timeout after 15 seconds")
			return context.DeadlineExceeded
		}

		panic(err)
	}

	if r.StatusCode() != 200 {
		fmt.Println("error trying to mark episode:")
		fmt.Println(string(r.Body()))
	}

	return nil
}
