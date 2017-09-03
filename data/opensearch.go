// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

import (
	"github.com/caixw/typing/data/xmlwriter"
	"github.com/caixw/typing/vars"
)

// Opensearch 相关内容
type Opensearch struct {
	URL     string
	Type    string
	Title   string
	Content []byte
}

type opensearchConfig struct {
	URL   string `yaml:"url"`             // opensearch 的地址，不能包含域名
	Type  string `yaml:"type,omitempty"`  // mimeType 默认取 vars.ContentTypeOpensearch
	Title string `yaml:"title,omitempty"` // 出现于 html>head>link.title 属性中

	ShortName   string `yaml:"shortName"`
	Description string `yaml:"description"`
	LongName    string `yaml:"longName,omitempty"`
	Image       *Icon  `yaml:"image,omitempty"`
}

// 用于生成一个符合 atom 规范的 XML 文本。
func (d *Data) buildOpensearch() error {
	w := xmlwriter.New()
	o := d.Config.Opensearch

	w.WriteStartElement("OpenSearchDescription", map[string]string{
		"xmlns": "http://a9.com/-/spec/opensearch/1.1/",
	})

	w.WriteElement("InputEncoding", "UTF-8", nil)
	w.WriteElement("OutputEncoding", "UTF-8", nil)
	w.WriteElement("ShortName", o.ShortName, nil)
	w.WriteElement("Description", o.Description, nil)

	if len(o.LongName) > 0 {
		w.WriteElement("LongName", o.LongName, nil)
	}

	if o.Image != nil {
		w.WriteElement("Image", o.Image.URL, map[string]string{
			"type": o.Image.Type,
		})
	}

	w.WriteCloseElement("Url", map[string]string{
		"type":     d.Config.Type,
		"template": vars.SearchURL("{searchTerms}", 0),
	})

	w.WriteElement("Developer", vars.AppName, nil)
	w.WriteElement("Language", d.Config.Language, nil)

	w.WriteEndElement("OpenSearchDescription")

	bs, err := w.Bytes()
	if err != nil {
		return err
	}
	d.Opensearch = &Opensearch{
		URL:     d.Config.Opensearch.URL,
		Type:    d.Config.Opensearch.Type,
		Title:   d.Config.Opensearch.Title,
		Content: bs,
	}

	return nil
}

// 检测 opensearch 取值是否正确
func (s *opensearchConfig) sanitize(conf *Config) *FieldError {
	switch {
	case len(s.URL) == 0:
		return &FieldError{Message: "不能为空", Field: "Opensearch.URL"}
	case len(s.ShortName) == 0:
		return &FieldError{Message: "不能为空", Field: "Opensearch.ShortName"}
	case len(s.Description) == 0:
		return &FieldError{Message: "不能为空", Field: "Opensearch.Description"}
	}

	if len(s.Type) == 0 {
		s.Type = vars.ContentTypeOpensearch
	}

	if s.Image == nil && conf.Icon != nil {
		s.Image = conf.Icon
	}

	return nil
}
