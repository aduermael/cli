package docker

import "testing"

func TestImagesFilter(t *testing.T) {
	eng := NewTestEngine(t)
	defer nuke(mkDaemonFromEngine(eng, t))

	if err := eng.Job("tag", unitTestImageName, "utest", "tag1").Run(); err != nil {
		t.Fatal(err)
	}

	if err := eng.Job("tag", unitTestImageName, "utest/docker", "tag2").Run(); err != nil {
		t.Fatal(err)
	}

	if err := eng.Job("tag", unitTestImageName, "utest:5000/docker", "tag3").Run(); err != nil {
		t.Fatal(err)
	}

	images := getImages(eng, t, false, "utest*/*")

	if len(images[0].RepoTags) != 2 {
		t.Fatal("incorrect number of matches returned")
	}

	images = getImages(eng, t, false, "utest")

	if len(images[0].RepoTags) != 1 {
		t.Fatal("incorrect number of matches returned")
	}

	images = getImages(eng, t, false, "utest*")

	if len(images[0].RepoTags) != 1 {
		t.Fatal("incorrect number of matches returned")
	}

	images = getImages(eng, t, false, "*5000*/*")

	if len(images[0].RepoTags) != 1 {
		t.Fatal("incorrect number of matches returned")
	}
}
