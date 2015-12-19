package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/docker/distribution"
	"github.com/docker/distribution/digest"
	"github.com/docker/distribution/manifest"
	"github.com/docker/distribution/manifest/manifestlist"
	"github.com/docker/distribution/manifest/schema2"
	"github.com/docker/docker/pkg/integration/checker"
	"github.com/go-check/check"
)

// testPullImageWithAliases pulls a specific image tag and verifies that any aliases (i.e., other
// tags for the same image) are not also pulled down.
//
// Ref: docker/docker#8141
func testPullImageWithAliases(c *check.C) {
	repoName := fmt.Sprintf("%v/dockercli/busybox", privateRegistryURL)

	repos := []string{}
	for _, tag := range []string{"recent", "fresh"} {
		repos = append(repos, fmt.Sprintf("%v:%v", repoName, tag))
	}

	// Tag and push the same image multiple times.
	for _, repo := range repos {
		dockerCmd(c, "tag", "busybox", repo)
		dockerCmd(c, "push", repo)
	}

	// Clear local images store.
	args := append([]string{"rmi"}, repos...)
	dockerCmd(c, args...)

	// Pull a single tag and verify it doesn't bring down all aliases.
	dockerCmd(c, "pull", repos[0])
	dockerCmd(c, "inspect", repos[0])
	for _, repo := range repos[1:] {
		_, _, err := dockerCmdWithError("inspect", repo)
		c.Assert(err, checker.NotNil, check.Commentf("Image %v shouldn't have been pulled down", repo))
	}
}

func (s *DockerRegistrySuite) TestPullImageWithAliases(c *check.C) {
	testPullImageWithAliases(c)
}

func (s *DockerSchema1RegistrySuite) TestPullImageWithAliases(c *check.C) {
	testPullImageWithAliases(c)
}

// testConcurrentPullWholeRepo pulls the same repo concurrently.
func testConcurrentPullWholeRepo(c *check.C) {
	repoName := fmt.Sprintf("%v/dockercli/busybox", privateRegistryURL)

	repos := []string{}
	for _, tag := range []string{"recent", "fresh", "todays"} {
		repo := fmt.Sprintf("%v:%v", repoName, tag)
		_, err := buildImage(repo, fmt.Sprintf(`
		    FROM busybox
		    ENTRYPOINT ["/bin/echo"]
		    ENV FOO foo
		    ENV BAR bar
		    CMD echo %s
		`, repo), true)
		c.Assert(err, checker.IsNil)
		dockerCmd(c, "push", repo)
		repos = append(repos, repo)
	}

	// Clear local images store.
	args := append([]string{"rmi"}, repos...)
	dockerCmd(c, args...)

	// Run multiple re-pulls concurrently
	results := make(chan error)
	numPulls := 3

	for i := 0; i != numPulls; i++ {
		go func() {
			_, _, err := runCommandWithOutput(exec.Command(dockerBinary, "pull", "-a", repoName))
			results <- err
		}()
	}

	// These checks are separate from the loop above because the check
	// package is not goroutine-safe.
	for i := 0; i != numPulls; i++ {
		err := <-results
		c.Assert(err, checker.IsNil, check.Commentf("concurrent pull failed with error: %v", err))
	}

	// Ensure all tags were pulled successfully
	for _, repo := range repos {
		dockerCmd(c, "inspect", repo)
		out, _ := dockerCmd(c, "run", "--rm", repo)
		c.Assert(strings.TrimSpace(out), checker.Equals, "/bin/sh -c echo "+repo)
	}
}

func (s *DockerRegistrySuite) testConcurrentPullWholeRepo(c *check.C) {
	testConcurrentPullWholeRepo(c)
}

func (s *DockerSchema1RegistrySuite) testConcurrentPullWholeRepo(c *check.C) {
	testConcurrentPullWholeRepo(c)
}

// testConcurrentFailingPull tries a concurrent pull that doesn't succeed.
func testConcurrentFailingPull(c *check.C) {
	repoName := fmt.Sprintf("%v/dockercli/busybox", privateRegistryURL)

	// Run multiple pulls concurrently
	results := make(chan error)
	numPulls := 3

	for i := 0; i != numPulls; i++ {
		go func() {
			_, _, err := runCommandWithOutput(exec.Command(dockerBinary, "pull", repoName+":asdfasdf"))
			results <- err
		}()
	}

	// These checks are separate from the loop above because the check
	// package is not goroutine-safe.
	for i := 0; i != numPulls; i++ {
		err := <-results
		c.Assert(err, checker.NotNil, check.Commentf("expected pull to fail"))
	}
}

func (s *DockerRegistrySuite) testConcurrentFailingPull(c *check.C) {
	testConcurrentFailingPull(c)
}

func (s *DockerSchema1RegistrySuite) testConcurrentFailingPull(c *check.C) {
	testConcurrentFailingPull(c)
}

// testConcurrentPullMultipleTags pulls multiple tags from the same repo
// concurrently.
func testConcurrentPullMultipleTags(c *check.C) {
	repoName := fmt.Sprintf("%v/dockercli/busybox", privateRegistryURL)

	repos := []string{}
	for _, tag := range []string{"recent", "fresh", "todays"} {
		repo := fmt.Sprintf("%v:%v", repoName, tag)
		_, err := buildImage(repo, fmt.Sprintf(`
		    FROM busybox
		    ENTRYPOINT ["/bin/echo"]
		    ENV FOO foo
		    ENV BAR bar
		    CMD echo %s
		`, repo), true)
		c.Assert(err, checker.IsNil)
		dockerCmd(c, "push", repo)
		repos = append(repos, repo)
	}

	// Clear local images store.
	args := append([]string{"rmi"}, repos...)
	dockerCmd(c, args...)

	// Re-pull individual tags, in parallel
	results := make(chan error)

	for _, repo := range repos {
		go func(repo string) {
			_, _, err := runCommandWithOutput(exec.Command(dockerBinary, "pull", repo))
			results <- err
		}(repo)
	}

	// These checks are separate from the loop above because the check
	// package is not goroutine-safe.
	for range repos {
		err := <-results
		c.Assert(err, checker.IsNil, check.Commentf("concurrent pull failed with error: %v", err))
	}

	// Ensure all tags were pulled successfully
	for _, repo := range repos {
		dockerCmd(c, "inspect", repo)
		out, _ := dockerCmd(c, "run", "--rm", repo)
		c.Assert(strings.TrimSpace(out), checker.Equals, "/bin/sh -c echo "+repo)
	}
}

func (s *DockerRegistrySuite) TestConcurrentPullMultipleTags(c *check.C) {
	testConcurrentPullMultipleTags(c)
}

func (s *DockerSchema1RegistrySuite) TestConcurrentPullMultipleTags(c *check.C) {
	testConcurrentPullMultipleTags(c)
}

// testPullIDStability verifies that pushing an image and pulling it back
// preserves the image ID.
func testPullIDStability(c *check.C) {
	derivedImage := privateRegistryURL + "/dockercli/id-stability"
	baseImage := "busybox"

	_, err := buildImage(derivedImage, fmt.Sprintf(`
	    FROM %s
	    ENV derived true
	    ENV asdf true
	    RUN dd if=/dev/zero of=/file bs=1024 count=1024
	    CMD echo %s
	`, baseImage, derivedImage), true)
	if err != nil {
		c.Fatal(err)
	}

	originalID, err := getIDByName(derivedImage)
	if err != nil {
		c.Fatalf("error inspecting: %v", err)
	}
	dockerCmd(c, "push", derivedImage)

	// Pull
	out, _ := dockerCmd(c, "pull", derivedImage)
	if strings.Contains(out, "Pull complete") {
		c.Fatalf("repull redownloaded a layer: %s", out)
	}

	derivedIDAfterPull, err := getIDByName(derivedImage)
	if err != nil {
		c.Fatalf("error inspecting: %v", err)
	}

	if derivedIDAfterPull != originalID {
		c.Fatal("image's ID unexpectedly changed after a repush/repull")
	}

	// Make sure the image runs correctly
	out, _ = dockerCmd(c, "run", "--rm", derivedImage)
	if strings.TrimSpace(out) != derivedImage {
		c.Fatalf("expected %s; got %s", derivedImage, out)
	}

	// Confirm that repushing and repulling does not change the computed ID
	dockerCmd(c, "push", derivedImage)
	dockerCmd(c, "rmi", derivedImage)
	dockerCmd(c, "pull", derivedImage)

	derivedIDAfterPull, err = getIDByName(derivedImage)
	if err != nil {
		c.Fatalf("error inspecting: %v", err)
	}

	if derivedIDAfterPull != originalID {
		c.Fatal("image's ID unexpectedly changed after a repush/repull")
	}
	if err != nil {
		c.Fatalf("error inspecting: %v", err)
	}

	// Make sure the image still runs
	out, _ = dockerCmd(c, "run", "--rm", derivedImage)
	if strings.TrimSpace(out) != derivedImage {
		c.Fatalf("expected %s; got %s", derivedImage, out)
	}
}

func (s *DockerRegistrySuite) TestPullIDStability(c *check.C) {
	testPullIDStability(c)
}

func (s *DockerSchema1RegistrySuite) TestPullIDStability(c *check.C) {
	testPullIDStability(c)
}

// TestPullFallbackOn404 tries to pull a nonexistent manifest and confirms that
// the pull falls back to the v1 protocol.
//
// Ref: docker/docker#18832
func (s *DockerRegistrySuite) TestPullFallbackOn404(c *check.C) {
	repoName := fmt.Sprintf("%v/does/not/exist", privateRegistryURL)

	out, _, _ := dockerCmdWithError("pull", repoName)

	c.Assert(out, checker.Contains, "v1 ping attempt")
}

func (s *DockerRegistrySuite) TestPullManifestList(c *check.C) {
	pushDigest, err := setupImage(c)
	c.Assert(err, checker.IsNil, check.Commentf("error setting up image"))

	// Inject a manifest list into the registry
	manifestList := &manifestlist.ManifestList{
		Versioned: manifest.Versioned{
			SchemaVersion: 2,
			MediaType:     manifestlist.MediaTypeManifestList,
		},
		Manifests: []manifestlist.ManifestDescriptor{
			{
				Descriptor: distribution.Descriptor{
					Digest:    "sha256:1a9ec845ee94c202b2d5da74a24f0ed2058318bfa9879fa541efaecba272e86b",
					Size:      3253,
					MediaType: schema2.MediaTypeManifest,
				},
				Platform: manifestlist.PlatformSpec{
					Architecture: "bogus_arch",
					OS:           "bogus_os",
				},
			},
			{
				Descriptor: distribution.Descriptor{
					Digest:    pushDigest,
					Size:      3253,
					MediaType: schema2.MediaTypeManifest,
				},
				Platform: manifestlist.PlatformSpec{
					Architecture: runtime.GOARCH,
					OS:           runtime.GOOS,
				},
			},
		},
	}

	manifestListJSON, err := json.MarshalIndent(manifestList, "", "   ")
	c.Assert(err, checker.IsNil, check.Commentf("error marshalling manifest list"))

	manifestListDigest := digest.FromBytes(manifestListJSON)
	hexDigest := manifestListDigest.Hex()

	registryV2Path := filepath.Join(s.reg.dir, "docker", "registry", "v2")

	// Write manifest list to blob store
	blobDir := filepath.Join(registryV2Path, "blobs", "sha256", hexDigest[:2], hexDigest)
	err = os.MkdirAll(blobDir, 0755)
	c.Assert(err, checker.IsNil, check.Commentf("error creating blob dir"))
	blobPath := filepath.Join(blobDir, "data")
	err = ioutil.WriteFile(blobPath, []byte(manifestListJSON), 0644)
	c.Assert(err, checker.IsNil, check.Commentf("error writing manifest list"))

	// Add to revision store
	revisionDir := filepath.Join(registryV2Path, "repositories", remoteRepoName, "_manifests", "revisions", "sha256", hexDigest)
	err = os.Mkdir(revisionDir, 0755)
	c.Assert(err, checker.IsNil, check.Commentf("error creating revision dir"))
	revisionPath := filepath.Join(revisionDir, "link")
	err = ioutil.WriteFile(revisionPath, []byte(manifestListDigest.String()), 0644)
	c.Assert(err, checker.IsNil, check.Commentf("error writing revision link"))

	// Update tag
	tagPath := filepath.Join(registryV2Path, "repositories", remoteRepoName, "_manifests", "tags", "latest", "current", "link")
	err = ioutil.WriteFile(tagPath, []byte(manifestListDigest.String()), 0644)
	c.Assert(err, checker.IsNil, check.Commentf("error writing tag link"))

	// Verify that the image can be pulled through the manifest list.
	out, _ := dockerCmd(c, "pull", repoName)

	// The pull output includes "Digest: <digest>", so find that
	matches := digestRegex.FindStringSubmatch(out)
	c.Assert(matches, checker.HasLen, 2, check.Commentf("unable to parse digest from pull output: %s", out))
	pullDigest := matches[1]

	// Make sure the pushed and pull digests match
	c.Assert(manifestListDigest.String(), checker.Equals, pullDigest)

	// Was the image actually created?
	dockerCmd(c, "inspect", repoName)

	dockerCmd(c, "rmi", repoName)
}
