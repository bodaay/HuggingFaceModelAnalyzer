package hfanalyzer

type ModelFramework string

const (
	HUGGINGFACE ModelFramework = "hf"
	PYTORCH     ModelFramework = "pytorch"
	TENSORFLOW  ModelFramework = "tensorflow"
	ONNX        ModelFramework = "onnx"
)

type LLMModelInfoV1 struct {
	Version                string
	LocalPath              string
	RemoteURL              string
	ModelName              string
	ModelParametersSize    string
	ModelNumberOfTokens    string
	ModelConfigFile        string
	ModelTokenerConfigFile string
	ModelTokenerModelFile  string
	ModelReadMeFile        string
	PromptTypeFormat       string
	PromptFormatUser       string
	PromtFormatAssistant   string

	PreTrainedModels []LLMPreTrainedModel
}

type LLMPreTrainedModel struct {
	ModelFrameWork  ModelFramework
	IsCompressed    bool
	CompressionType string
	IsQuantaized    bool
	SafeTensors     bool
	Files           []LLMPreTrainedModelFiles
}

type LLMPreTrainedModelQuantConfig struct {
	QuantEngine  string
	Bits         int
	GroupSize    int
	WithActOrder bool
	IsSequencial bool
}

type LLMPreTrainedModelFiles struct {
	FileName string
	FileSize int64
	SHA256   string
}
