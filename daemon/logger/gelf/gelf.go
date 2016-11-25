// +build linux

// Package gelf provides the log driver for forwarding server logs to
// endpoints that support the Graylog Extended Log Format.
package gelf

import (
	"compress/flate"
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"time"

	"github.com/Graylog2/go-gelf/gelf"
	"github.com/Sirupsen/logrus"
	"github.com/docker/docker/daemon/logger"
	"github.com/docker/docker/daemon/logger/loggerutils"
	"github.com/docker/docker/pkg/urlutil"
)

const name = "gelf"

type gelfLogger struct {
	writer   *gelf.Writer
	ctx      logger.Context
	hostname string
	rawExtra json.RawMessage
}

func init() {
	if err := logger.RegisterLogDriver(name, New); err != nil {
		logrus.Fatal(err)
	}
	if err := logger.RegisterLogOptValidator(name, ValidateLogOpt); err != nil {
		logrus.Fatal(err)
	}
}

// New creates a gelf logger using the configuration passed in on the
// context. The supported context configuration variable is gelf-address.
func New(ctx logger.Context) (logger.Logger, error) {
	// parse gelf address
	address, err := parseAddress(ctx.Config["gelf-address"])
	if err != nil {
		return nil, err
	}

	// collect extra data for GELF message
	hostname, err := ctx.Hostname()
	if err != nil {
		return nil, fmt.Errorf("gelf: cannot access hostname to set source field")
	}

	// parse log tag
	tag, err := loggerutils.ParseLogTag(ctx, loggerutils.DefaultTemplate)
	if err != nil {
		return nil, err
	}

	extra := map[string]interface{}{
		"_container_id":   ctx.ContainerID,
		"_container_name": ctx.Name(),
		"_image_id":       ctx.ContainerImageID,
		"_image_name":     ctx.ContainerImageName,
		"_command":        ctx.Command(),
		"_tag":            tag,
		"_created":        ctx.ContainerCreated,
	}

	extraAttrs := ctx.ExtraAttributes(func(key string) string {
		if key[0] == '_' {
			return key
		}
		return "_" + key
	})
	for k, v := range extraAttrs {
		extra[k] = v
	}

	rawExtra, err := json.Marshal(extra)
	if err != nil {
		return nil, err
	}

	// create new gelfWriter
	gelfWriter, err := gelf.NewWriter(address)
	if err != nil {
		return nil, fmt.Errorf("gelf: cannot connect to GELF endpoint: %s %v", address, err)
	}

	if v, ok := ctx.Config["gelf-compression-type"]; ok {
		switch v {
		case "gzip":
			gelfWriter.CompressionType = gelf.CompressGzip
		case "zlib":
			gelfWriter.CompressionType = gelf.CompressZlib
		case "none":
			gelfWriter.CompressionType = gelf.CompressNone
		default:
			return nil, fmt.Errorf("gelf: invalid compression type %q", v)
		}
	}

	if v, ok := ctx.Config["gelf-compression-level"]; ok {
		val, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("gelf: invalid compression level %s, err %v", v, err)
		}
		gelfWriter.CompressionLevel = val
	}

	return &gelfLogger{
		writer:   gelfWriter,
		ctx:      ctx,
		hostname: hostname,
		rawExtra: rawExtra,
	}, nil
}

func (s *gelfLogger) Log(msg *logger.Message) error {
	level := gelf.LOG_INFO
	if msg.Source == "stderr" {
		level = gelf.LOG_ERR
	}

	m := gelf.Message{
		Version:  "1.1",
		Host:     s.hostname,
		Short:    string(msg.Line),
		TimeUnix: float64(msg.Timestamp.UnixNano()/int64(time.Millisecond)) / 1000.0,
		Level:    level,
		RawExtra: s.rawExtra,
	}

	if err := s.writer.WriteMessage(&m); err != nil {
		return fmt.Errorf("gelf: cannot send GELF message: %v", err)
	}
	return nil
}

func (s *gelfLogger) Close() error {
	return s.writer.Close()
}

func (s *gelfLogger) Name() string {
	return name
}

// ValidateLogOpt looks for gelf specific log option gelf-address.
func ValidateLogOpt(cfg map[string]string) error {
	for key, val := range cfg {
		switch key {
		case "gelf-address":
		case "tag":
		case "labels":
		case "env":
		case "gelf-compression-level":
			i, err := strconv.Atoi(val)
			if err != nil || i < flate.DefaultCompression || i > flate.BestCompression {
				return fmt.Errorf("unknown value %q for log opt %q for gelf log driver", val, key)
			}
		case "gelf-compression-type":
			switch val {
			case "gzip", "zlib", "none":
			default:
				return fmt.Errorf("unknown value %q for log opt %q for gelf log driver", val, key)
			}
		default:
			return fmt.Errorf("unknown log opt %q for gelf log driver", key)
		}
	}

	if _, err := parseAddress(cfg["gelf-address"]); err != nil {
		return err
	}

	return nil
}

func parseAddress(address string) (string, error) {
	if address == "" {
		return "", nil
	}
	if !urlutil.IsTransportURL(address) {
		return "", fmt.Errorf("gelf-address should be in form proto://address, got %v", address)
	}
	url, err := url.Parse(address)
	if err != nil {
		return "", err
	}

	// we support only udp
	if url.Scheme != "udp" {
		return "", fmt.Errorf("gelf: endpoint needs to be UDP")
	}

	// get host and port
	if _, _, err = net.SplitHostPort(url.Host); err != nil {
		return "", fmt.Errorf("gelf: please provide gelf-address as udp://host:port")
	}

	return url.Host, nil
}
