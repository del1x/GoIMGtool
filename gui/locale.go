package gui

type Locale struct {
	LanguageLabel          string
	WatermarkLabel         string
	ImageDirLabel          string
	FormatLabel            string
	QualityLabel           string
	WebSizeHint            string
	WidthLabel             string
	HeightLabel            string
	TargetSizeLabel        string
	WatermarkModeLabel     string
	WatermarkPlaceholder   string
	ImageDirPlaceholder    string
	ProcessButton          string
	CurrentFileLabel       string
	ProcessingDone         string
	BrowseButton           string
	BrowseFolderButton     string
	ErrorTitle             string
	NoWatermarkOrFolder    string
	FailedSelectWatermark  string
	InvalidFile            string
	FailedSelectFolder     string
	InvalidFolder          string
	FailedInitProcessor    string
	ProcessingFailed       string
	WidthExceedsWatermark  string
	HeightExceedsWatermark string
}

var locales = map[string]Locale{
	"en": {
		LanguageLabel:          "Language:",
		WatermarkLabel:         "Watermark file:",
		ImageDirLabel:          "Image folder:",
		FormatLabel:            "Output format:",
		QualityLabel:           "Quality (1-100):",
		WebSizeHint:            "Note: For web, target size ≤100 KB is optimal",
		WidthLabel:             "Max Width (100-4096):",
		HeightLabel:            "Max Height (100-4096):",
		TargetSizeLabel:        "Target Size (KB, 50-5000):",
		WatermarkModeLabel:     "Watermark Mode:",
		WatermarkPlaceholder:   "Select watermark.png",
		ImageDirPlaceholder:    "Select image folder",
		ProcessButton:          "Process",
		CurrentFileLabel:       "Processing: None",
		ProcessingDone:         "Processing: Done",
		BrowseButton:           "Browse...",
		BrowseFolderButton:     "Browse Folder...",
		ErrorTitle:             "Error",
		NoWatermarkOrFolder:    "Please select a watermark file and image folder!",
		FailedSelectWatermark:  "Failed to select watermark file!",
		InvalidFile:            "Invalid file: %v",
		FailedSelectFolder:     "Failed to select folder!",
		InvalidFolder:          "Invalid folder: %v",
		FailedInitProcessor:    "Failed to initialize processor: %v",
		ProcessingFailed:       "Processing failed: %v",
		WidthExceedsWatermark:  "Width exceeds watermark width (%d px)",
		HeightExceedsWatermark: "Height exceeds watermark height (%d px)",
	},
	"ru": {
		LanguageLabel:          "Язык:",
		WatermarkLabel:         "Файл водяного знака:",
		ImageDirLabel:          "Папка с изображениями:",
		FormatLabel:            "Формат вывода:",
		QualityLabel:           "Качество (1-100):",
		WebSizeHint:            "Примечание: Для веба оптимальный размер ≤100 КБ",
		WidthLabel:             "Макс. ширина (100-4096):",
		HeightLabel:            "Макс. высота (100-4096):",
		TargetSizeLabel:        "Целевой размер (КБ, 50-5000):",
		WatermarkModeLabel:     "Режим водяного знака:",
		WatermarkPlaceholder:   "Выберите watermark.png",
		ImageDirPlaceholder:    "Выберите папку с изображениями",
		ProcessButton:          "Обработать",
		CurrentFileLabel:       "Обработка: Нет",
		ProcessingDone:         "Обработка: Завершено",
		BrowseButton:           "Выбрать watermark...",
		BrowseFolderButton:     "Выбрать папку...",
		ErrorTitle:             "Ошибка",
		NoWatermarkOrFolder:    "Пожалуйста, выберите файл водяного знака и папку с изображениями!",
		FailedSelectWatermark:  "Не удалось выбрать файл водяного знака!",
		InvalidFile:            "Недопустимый файл: %v",
		FailedSelectFolder:     "Не удалось выбрать папку!",
		InvalidFolder:          "Недопустимая папка: %v",
		FailedInitProcessor:    "Не удалось инициализировать процессор: %v",
		ProcessingFailed:       "Ошибка обработки: %v",
		WidthExceedsWatermark:  "Ширина превышает ширину водяного знака (%d пикс.)",
		HeightExceedsWatermark: "Высота превышает высоту водяного знака (%d пикс.)",
	},
}
