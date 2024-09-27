package speech

import (
	"fmt"
)

type RecognizeRequest struct {
	// 음성파일 처리를 위한 Config를 정의합니다.
	// Config를 작성하지 않으면 Default 값으로 사용됩니다.
	Config RecognitionConfig
	// Content 와 FilePath 둘 중 하나만을 전달해야합니다.
	// 만약 두 개 동시에 제공한다면 에러가 발생합니다.
	AudioSource RecognitionAudio
}

type ResultId string

type RecognitionConfig struct {
	// 사용할 모델 이름을 정의합니다. Default 값으로 sommers가 사용됩니다.
	ModelName string `json:"model_name,omitempty"`
	// 모델 중 whisper가 사용되었을 시에 Language를 제공해야합니다.
	Language string `json:"language,omitempty"`
	// 화자 분리 여부를 정의합니다. Default 값으로 False 입니다.
	UseDiarization *bool `json:"use_diarization,omitempty"`
	// 화자 분리 사용 시 이미 화자 수를 알 경우 사용하는 파라미터입니다.
	Diarization *DiarizationConfig `json:"diarization,omitempty"`
	// 영어,단어,숫자 등의 표현 변환을 정의합니다. Default 값으로 True 입니다.
	UseItn *bool `json:"use_itn,omitempty"`
	// 간투어 필터 설정을 정의합니다. Default 값으로 True 입니다.
	UseDisfluencyFilter *bool `json:"use_disfluency_filter,omitempty"`
	// 비속어 필터 설정을 정의합니다. Default 값으로 False 입니다.
	UseProfanityFilter *bool `json:"use_profanity_filter,omitempty"`
	// 문단 나누기 설정을 정의합니다. Default 값으로 True입니다.
	UseParagraphSplitter *bool `json:"use_paragraph_splitter,omitempty"`
	// 문단 나누기 수준을 정의합니다. Default 값으로 50 입니다.
	ParagraphSpliter *ParagraphSplitterConfig `json:"paragraph_splitter,omitempty"`
	// 도메인 설정을 정의합니다. GENERAL, CALL 이 존재하며 Default 값으로 GENERAL입니다.
	Domain string `json:"domain,omitempty"`
	// 단어 수준의 타임스탬프 설정을 정의합니다. Default 값은 False입니다.
	// 전사된 텍스트를 원래 오디오와 정확히 일치시킬 필요가 있을 때 특히 유용합니다.
	UseWordTimestamp *bool `json:"use_word_timestamp,omitempty"`
	// 특정 키워드에 대한 전사 정확도를 높이기 위해 사용됩니다.
	// 키워드는 한글만 지원합니다.
	Keywords []string `json:"keywords,omitempty"`
}

// DiarizationConfig는 발화자 분리 설정을 포함하는 구조체입니다.
type DiarizationConfig struct {
	SpkCount int `json:"spk_count"`
}

// ParagraphSplitterConfig는 단락 분리 설정을 포함하는 구조체입니다.
type ParagraphSplitterConfig struct {
	Max int `json:"max"`
}

// Content 와 FilePath 둘 중 하나만을 전달해야합니다.
// 만약 두 개 동시에 제공한다면 에러가 발생합니다.
type RecognitionAudio struct {
	Content  []byte
	FilePath string
}

func (ra *RecognitionAudio) validate() error {
	if ra.Content != nil && ra.FilePath != "" {
		return fmt.Errorf("both Content and FilePath are provided; please provide only one")
	}
	if ra.Content == nil && ra.FilePath == "" {
		return fmt.Errorf("neither Content nor FilePath is provided; please provide one")
	}
	return nil
}

type RecognizeResponse struct {
	Id      ResultId `json:"id"`
	Status  string   `json:"status"`
	Results *Results `json:"results"`
}

type Results struct {
	Utterances []*Utterance `json:"utterances"`
	Verified   bool         `json:"verified"`
}
type Utterance struct {
	Duration int              `json:"duration"`
	Msg      string           `json:"msg"`
	Spk      int              `json:"spk"`
	SpkType  string           `json:"spk_type"`
	StartAt  int              `json:"start_at"`
	Words    []*TimeStampWord `json:"words"`
}

type TimeStampWord struct {
	StartAt  int    `json:"start_at"`
	Duration int    `json:"duration"`
	Text     string `json:"text"`
}
