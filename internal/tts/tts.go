package tts

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"log"
	"time"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"cloud.google.com/go/texttospeech/apiv1/texttospeechpb"
	"github.com/gopxl/beep"
	"github.com/gopxl/beep/speaker"
	"github.com/gopxl/beep/wav"
)

type TTS struct {
	Client *texttospeech.Client
	Ctx    context.Context
	Word   string
	audio  audioMessage
}

type audioMessage struct {
	AudioContent []byte
	Word         string
}

// SayWord takes a word and uses the Google Cloud Text-to-Speech API to generate and play the audio for that word.
// It checks if the audio for the word is already generated and stored in the audioMessage struct.
// If the audio is already generated, it plays the audio directly without calling the API again.
// If the audio is not generated, it calls the API to synthesize the speech and then plays the audio.
func (t *TTS) SayWord() error {
	// call tts api only if not already done
	if t.audio.Word != t.Word {
		audioContent, err := t.synthesizeSpeech()
		if err != nil {
			return errors.New(fmt.Sprintf("error synthesizing speech: %v", err))
		}
		t.audio.AudioContent = audioContent
		t.audio.Word = t.Word
	}
	// play the audio
	t.PlayAudio()
	return nil
}

func (t *TTS) PlayAudio() {
	r := bytes.NewReader(t.audio.AudioContent)
	// Play the audio
	streamer, format, err := wav.Decode(r)
	if err != nil {
		log.Fatal(err)
	}
	defer streamer.Close()

	sr := format.SampleRate
	speaker.Init(sr, sr.N(time.Second/10))

	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))
	<-done
}

// synthesizeSpeech uses the Google Cloud Text-to-Speech API to synthesize speech from the given word.
// It creates a request with the word, voice parameters, and audio configuration.
// It then calls the API and returns the synthesized audio content.
func (t *TTS) synthesizeSpeech() ([]byte, error) {
	// set up request
	req := texttospeechpb.SynthesizeSpeechRequest{
		Input: &texttospeechpb.SynthesisInput{
			InputSource: &texttospeechpb.SynthesisInput_Text{
				Text: t.Word,
			},
		},
		// configure voice
		Voice: &texttospeechpb.VoiceSelectionParams{
			LanguageCode: "en-US",
			Name:         "en-US-Chirp3-HD-Fenrir", // this one makes me laugh bc he's zesty
			// Name: "en-US-Wavenet-D", // this one is more natural
			// Name: "en-US-Neural2-J", // this is the og
			SsmlGender: texttospeechpb.SsmlVoiceGender_MALE,
		},
		// Configure the audio output
		AudioConfig: &texttospeechpb.AudioConfig{
			AudioEncoding: texttospeechpb.AudioEncoding_LINEAR16, // WAV format
		},
	}

	// Call the API
	resp, err := t.Client.SynthesizeSpeech(t.Ctx, &req)
	if err != nil {
		return nil, err
	}

	return resp.AudioContent, nil
}
