package speakmanager

import (
	"github.com/cwdot/stdlib-go/wood"
	"hass/internal/config"

	"hass/internal/hassclient"
)

func New(targets map[string]config.SpeakerTarget) *SpeakManager {
	return &SpeakManager{targets: targets}
}

type SpeakManager struct {
	targets map[string]config.SpeakerTarget
}

func (c *SpeakManager) Speak(client *hassclient.Client, target string, message string) error {
	if st, ok := c.targets[target]; ok {
		entities := st.Players
		args := map[string]any{
			"entity_id":              "tts.piper",
			"message":                message,
			"media_player_entity_id": entities,
		}
		return client.Service("tts", "speak", args)
	}
	wood.Warnf("Failed to find target: %s", target)
	return nil
}
