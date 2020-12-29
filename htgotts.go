package htgotts

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/weiztd/htgo-tts/handlers"
)

/**
 * Required:
 * - mplayer
 *
 * Use:
 *
 * speech := htgotts.Speech{Folder: "audio", Language: "en", Handler: MPlayer}
 */

// Speech struct
type Speech struct {
	Folder   string
	Language string
	Handler  handlers.PlayerInterface
}

// Speak downloads speech and plays it using mplayer
func (speech *Speech) Speak(text string) error {

	fileName := speech.Folder + "/TTS.mp3"

	var err error
	if err = speech.createFolderIfNotExists(speech.Folder); err != nil {
		return err
	}
	if err = speech.download(fileName, text); err != nil {
		return err
	}

	if speech.Handler == nil {
		mplayer := handlers.MPlayer{}
		return mplayer.Play(fileName)
	}

	return speech.Handler.Play(fileName)
}

/**
 * Create the folder if does not exists.
 */
func (speech *Speech) createFolderIfNotExists(folder string) error {
	dir, err := os.Open(folder)
	if os.IsNotExist(err) {
		return os.MkdirAll(folder, 0700)
	}

	dir.Close()
	return nil
}

/**
 * Download the voice file if does not exists.
 */
func (speech *Speech) download(fileName string, text string) error {
	url := fmt.Sprintf("http://translate.google.com/translate_tts?ie=UTF-8&total=1&idx=0&textlen=32&client=tw-ob&q=%s&tl=%s", url.QueryEscape(text), speech.Language)
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	output, err := os.Create(fileName)
	if err != nil {
		return err
	}

	_, err = io.Copy(output, response.Body)
	if err != nil {
		return err
	}
	return nil

}
