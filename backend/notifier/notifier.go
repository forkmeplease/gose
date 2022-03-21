package notifier

import (
	"bytes"
	"fmt"
	"math"
	"net"
	"path/filepath"
	"text/template"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/containrrr/shoutrrr"
	"github.com/containrrr/shoutrrr/pkg/router"
	"github.com/containrrr/shoutrrr/pkg/types"
	"github.com/stv0g/gose/backend/config"
	"github.com/stv0g/gose/backend/utils"
)

type NotifierArgs struct {
	URL              string
	FileSize         int64
	FileSizeHuman    string
	FileName         string
	FileType         string
	UploaderIP       string
	UploaderHostname string
	Env              map[string]string
}

type Notifier struct {
	*router.ServiceRouter

	template *template.Template
}

func NewNotifier(cfg *config.NotificationConfig) (*Notifier, error) {
	sender, err := shoutrrr.CreateSender(cfg.URLs...)
	if err != nil {
		return nil, err
	}

	t := template.New("action")

	t, err = t.Parse(cfg.Template)
	if err != nil {
		return nil, fmt.Errorf("failed to parse notification template: %w", err)
	}

	return &Notifier{
		ServiceRouter: sender,
		template:      t,
	}, nil
}

func (n *Notifier) Notify(svc *s3.S3, cfg *config.Config, key string) error {
	obj, err := svc.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(cfg.S3.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return err
	}

	env, err := utils.EnvToMap()
	if err != nil {
		return fmt.Errorf("failed to get env: %w", err)
	}

	data := NotifierArgs{
		FileName:      filepath.Base(key),
		FileSize:      *obj.ContentLength,
		FileSizeHuman: humanizeBytes(*obj.ContentLength),
		FileType:      *obj.ContentType,
		Env:           env,
	}

	if u, ok := obj.Metadata["Url"]; ok {
		data.URL = *u
	}

	if upl, ok := obj.Metadata["Uploaded-By"]; ok {
		data.UploaderIP = *upl

		if addrs, err := net.LookupAddr(data.UploaderIP); err != nil && len(addrs) > 0 {
			data.UploaderHostname = addrs[0]
		}
	}

	var tpl bytes.Buffer
	if err := n.template.Execute(&tpl, data); err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	msg := tpl.String()

	if errs := n.Send(msg, &types.Params{
		"title": "New upload",
	}); errs != nil {
		for _, err := range errs {
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func logn(n, b float64) float64 {
	return math.Log(n) / math.Log(b)
}

func humanizeBytes(s int64) string {
	sizes := []string{"B", "KiB", "MiB", "GiB", "TiB", "PiB", "EiB"}
	base := 1024.0

	if s < 10 {
		return fmt.Sprintf("%d B", s)
	}
	e := math.Floor(logn(float64(s), base))
	suffix := sizes[int(e)]
	val := math.Floor(float64(s)/math.Pow(base, e)*10+0.5) / 10
	f := "%.0f %s"
	if val < 10 {
		f = "%.1f %s"
	}

	return fmt.Sprintf(f, val, suffix)
}
