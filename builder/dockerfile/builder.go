package dockerfile

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/docker/docker/builder"
	"github.com/docker/docker/builder/dockerfile/parser"
	"github.com/docker/docker/pkg/stringid"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/container"
)

var validCommitCommands = map[string]bool{
	"cmd":        true,
	"entrypoint": true,
	"env":        true,
	"expose":     true,
	"label":      true,
	"onbuild":    true,
	"user":       true,
	"volume":     true,
	"workdir":    true,
}

// BuiltinAllowedBuildArgs is list of built-in allowed build args
var BuiltinAllowedBuildArgs = map[string]bool{
	"HTTP_PROXY":  true,
	"http_proxy":  true,
	"HTTPS_PROXY": true,
	"https_proxy": true,
	"FTP_PROXY":   true,
	"ftp_proxy":   true,
	"NO_PROXY":    true,
	"no_proxy":    true,
}

// Builder is a Dockerfile builder
// It implements the builder.Backend interface.
type Builder struct {
	options *types.ImageBuildOptions

	Stdout io.Writer
	Stderr io.Writer

	docker  builder.Backend
	context builder.Context

	dockerfile       *parser.Node
	runConfig        *container.Config // runconfig for cmd, run, entrypoint etc.
	flags            *BFlags
	tmpContainers    map[string]struct{}
	image            string // imageID
	noBaseImage      bool
	maintainer       string
	cmdSet           bool
	disableCommit    bool
	cacheBusted      bool
	cancelled        chan struct{}
	cancelOnce       sync.Once
	allowedBuildArgs map[string]bool // list of build-time args that are allowed for expansion/substitution and passing to commands in 'run'.

	// TODO: remove once docker.Commit can receive a tag
	id string
}

// NewBuilder creates a new Dockerfile builder from an optional dockerfile and a Config.
// If dockerfile is nil, the Dockerfile specified by Config.DockerfileName,
// will be read from the Context passed to Build().
func NewBuilder(config *types.ImageBuildOptions, backend builder.Backend, context builder.Context, dockerfile io.ReadCloser) (b *Builder, err error) {
	if config == nil {
		config = new(types.ImageBuildOptions)
	}
	if config.BuildArgs == nil {
		config.BuildArgs = make(map[string]string)
	}
	b = &Builder{
		options:          config,
		Stdout:           os.Stdout,
		Stderr:           os.Stderr,
		docker:           backend,
		context:          context,
		runConfig:        new(container.Config),
		tmpContainers:    map[string]struct{}{},
		cancelled:        make(chan struct{}),
		id:               stringid.GenerateNonCryptoID(),
		allowedBuildArgs: make(map[string]bool),
	}
	if dockerfile != nil {
		b.dockerfile, err = parser.Parse(dockerfile)
		if err != nil {
			return nil, err
		}
	}

	return b, nil
}

// Build runs the Dockerfile builder from a context and a docker object that allows to make calls
// to Docker.
//
// This will (barring errors):
//
// * read the dockerfile from context
// * parse the dockerfile if not already parsed
// * walk the AST and execute it by dispatching to handlers. If Remove
//   or ForceRemove is set, additional cleanup around containers happens after
//   processing.
// * Print a happy message and return the image ID.
// * NOT tag the image, that is responsibility of the caller.
//
func (b *Builder) Build() (string, error) {
	// If Dockerfile was not parsed yet, extract it from the Context
	if b.dockerfile == nil {
		if err := b.readDockerfile(); err != nil {
			return "", err
		}
	}

	var shortImgID string
	for i, n := range b.dockerfile.Children {
		select {
		case <-b.cancelled:
			logrus.Debug("Builder: build cancelled!")
			fmt.Fprintf(b.Stdout, "Build cancelled")
			return "", fmt.Errorf("Build cancelled")
		default:
			// Not cancelled yet, keep going...
		}
		if err := b.dispatch(i, n); err != nil {
			if b.options.ForceRemove {
				b.clearTmp()
			}
			return "", err
		}
		shortImgID = stringid.TruncateID(b.image)
		fmt.Fprintf(b.Stdout, " ---> %s\n", shortImgID)
		if b.options.Remove {
			b.clearTmp()
		}
	}

	// check if there are any leftover build-args that were passed but not
	// consumed during build. Return an error, if there are any.
	leftoverArgs := []string{}
	for arg := range b.options.BuildArgs {
		if !b.isBuildArgAllowed(arg) {
			leftoverArgs = append(leftoverArgs, arg)
		}
	}
	if len(leftoverArgs) > 0 {
		return "", fmt.Errorf("One or more build-args %v were not consumed, failing build.", leftoverArgs)
	}

	if b.image == "" {
		return "", fmt.Errorf("No image was generated. Is your Dockerfile empty?")
	}

	fmt.Fprintf(b.Stdout, "Successfully built %s\n", shortImgID)
	return b.image, nil
}

// Cancel cancels an ongoing Dockerfile build.
func (b *Builder) Cancel() {
	b.cancelOnce.Do(func() {
		close(b.cancelled)
	})
}

// BuildFromConfig will do build directly from parameter 'changes',
// which is treated like it is coming from a Dockerfile. It will:
// - Call parse.Parse() to get an AST root for the concatenated Dockerfile entries.
// - Do build by calling builder.dispatch() to call all entries' handling routines
//
// BuildFromConfig is used by the /commit endpoint, with the changes
// coming from the query parameter of the same name.
//
// TODO: Remove?
func BuildFromConfig(config *container.Config, changes []string) (*container.Config, error) {
	ast, err := parser.Parse(bytes.NewBufferString(strings.Join(changes, "\n")))
	if err != nil {
		return nil, err
	}

	// ensure that the commands are valid
	for _, n := range ast.Children {
		if !validCommitCommands[n.Value] {
			return nil, fmt.Errorf("%s is not a valid change command", n.Value)
		}
	}

	b, err := NewBuilder(nil, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	b.runConfig = config
	b.Stdout = ioutil.Discard
	b.Stderr = ioutil.Discard
	b.disableCommit = true

	for i, n := range ast.Children {
		if err := b.dispatch(i, n); err != nil {
			return nil, err
		}
	}

	return b.runConfig, nil
}
