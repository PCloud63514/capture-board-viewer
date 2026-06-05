package captureboard

// SentryDSN은 빌드 시 -ldflags로 주입됩니다.
// 예: go build -ldflags "-X 'capture-board-selector/internal/captureboard.SentryDSN=https://...'"
var SentryDSN string
