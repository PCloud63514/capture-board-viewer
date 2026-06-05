package ghrelease

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"capture-board-selector/internal/captureboard/domain"
)

const (
	apiURL    = "https://api.github.com/repos/PCloud63514/capture-board-viewer/releases/latest"
	assetName = "capture-board-selector.exe"
)

type githubReleaseResponse struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

type GitHubChecker struct{}

func NewGitHubChecker() domain.UpdateChecker {
	return &GitHubChecker{}
}

func (c *GitHubChecker) Check(currentVersion string) (*domain.UpdateInfo, error) {
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("업데이트 확인 실패: %w", err)
	}
	defer resp.Body.Close()

	var release githubReleaseResponse
	if err = json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("응답 파싱 실패: %w", err)
	}

	if !isNewer(currentVersion, release.TagName) {
		return nil, nil
	}

	for _, asset := range release.Assets {
		if asset.Name == assetName {
			return &domain.UpdateInfo{
				Version:     release.TagName,
				DownloadURL: asset.BrowserDownloadURL,
			}, nil
		}
	}
	return nil, nil
}

// isNewer reports whether latest is a higher semver than current.
func isNewer(current, latest string) bool {
	if current == "" || latest == "" {
		return false
	}
	cv := parseVersion(strings.TrimPrefix(current, "v"))
	lv := parseVersion(strings.TrimPrefix(latest, "v"))
	for i := range lv {
		if i >= len(cv) {
			return true
		}
		if lv[i] > cv[i] {
			return true
		}
		if lv[i] < cv[i] {
			return false
		}
	}
	return false
}

func parseVersion(v string) []int {
	parts := strings.Split(v, ".")
	nums := make([]int, len(parts))
	for i, p := range parts {
		nums[i], _ = strconv.Atoi(p)
	}
	return nums
}
