package scheduler

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	anlz "webcrawler/analyzer"
	base "webcrawler/base"
	dl "webcrawler/downloader"
	ipl "webcrawler/itempipeline"
	mdw "webcrawler/middleware"
)

func generateChannelManager(channelArgs base.ChannelArgs) mdw.ChannelManager {
	return mdw.NewChannelManager(channelArgs)
}

func generatePageDownloaderPool(
	poolSize uint32,
	httpClientGenerator GenHttpClient) (dl.PageDownloaderPool, error) {
	dlPool, err := dl.NewPageDownloaderPool(
		poolSize,
		func() dl.PageDownloader {
			return dl.NewPageDownloader(httpClientGenerator())
		},
	)
	if err != nil {
		return nil, err
	}
	return dlPool, nil
}

func generateAnalyzerPool(poolSize uint32) (anlz.AnalyzerPool, error) {
	analyzerPool, err := anlz.NewAnalyzerPool(
		poolSize,
		func() anlz.Analyzer {
			return anlz.NewAnalyzer()
		},
	)
	if err != nil {
		return nil, err
	}
	return analyzerPool, nil
}

func generateItemPipeline(itemProcessors []ipl.ProcessItem) ipl.ItemPipeline {
	return ipl.NewItemPipeline(itemProcessors)
}

// 生成组件实例代号。
func generateCode(prefix string, id uint32) string {
	return fmt.Sprintf("%s-%d", prefix, id)
}

// 解析组件实例代号。
func parseCode(code string) []string {
	result := make([]string, 2)
	var codePrefix string
	var id string
	index := strings.Index(code, "-")
	if index > 0 {
		codePrefix = code[:index]
		id = code[index+1:]
	} else {
		codePrefix = code
	}
	result[0] = codePrefix
	result[1] = id
	return result
}

var regexpForIp = regexp.MustCompile(`((?:(?:25[0-5]|2[0-4]\d|[01]?\d?\d)\.){3}(?:25[0-5]|2[0-4]\d|[01]?\d?\d))`)

var regexpForDomains = []*regexp.Regexp{
	// *.xx or *.xxx.xx
	regexp.MustCompile(`\.(com|com\.\w{2})$`),
	regexp.MustCompile(`\.(gov|gov\.\w{2})$`),
	regexp.MustCompile(`\.(net|net\.\w{2})$`),
	regexp.MustCompile(`\.(org|org\.\w{2})$`),
	// *.xx
	regexp.MustCompile(`\.me$`),
	regexp.MustCompile(`\.biz$`),
	regexp.MustCompile(`\.info$`),
	regexp.MustCompile(`\.name$`),
	regexp.MustCompile(`\.mobi$`),
	regexp.MustCompile(`\.so$`),
	regexp.MustCompile(`\.asia$`),
	regexp.MustCompile(`\.tel$`),
	regexp.MustCompile(`\.tv$`),
	regexp.MustCompile(`\.cc$`),
	regexp.MustCompile(`\.co$`),
	regexp.MustCompile(`\.\w{2}$`),
}

func getPrimaryDomain(host string) (string, error) {
	host = strings.TrimSpace(host)
	if host == "" {
		return "", errors.New("The host is empty!")
	}
	if regexpForIp.MatchString(host) {
		return host, nil
	}
	var suffixIndex int
	for _, re := range regexpForDomains {
		pos := re.FindStringIndex(host)
		if pos != nil {
			suffixIndex = pos[0]
			break
		}
	}
	if suffixIndex > 0 {
		var pdIndex int
		firstPart := host[:suffixIndex]
		index := strings.LastIndex(firstPart, ".")
		if index < 0 {
			pdIndex = 0
		} else {
			pdIndex = index + 1
		}
		return host[pdIndex:], nil
	} else {
		return "", errors.New("Unrecognized host!")
	}
}
