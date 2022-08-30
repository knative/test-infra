package fixtures

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git-fixtures/v4/internal/tgz"
	"gopkg.in/check.v1"
)

//go:generate esc -o data.go -pkg=fixtures data

var (
	files      = make(map[string]string)
	Filesystem = osfs.New(os.TempDir())
)

var fixtures = Fixtures{{
	Tags:         []string{"packfile", "ofs-delta", ".git", "root-reference"},
	URL:          "https://github.com/git-fixtures/root-references.git",
	Head:         "6ecf0ef2c2dffb796033e5a02219af86ec6584e5",
	PackfileHash: "135fe3d1ad828afe68706f1d481aedbcfa7a86d2",
	DotGitHash:   "78c5fb882e76286d8201016cffee63ea7060a0c2",
	ObjectsCount: 68,
}, {
	Tags:         []string{"packfile", "ofs-delta", ".git"},
	URL:          "https://github.com/git-fixtures/basic.git",
	Head:         "6ecf0ef2c2dffb796033e5a02219af86ec6584e5",
	PackfileHash: "a3fed42da1e8189a077c0e6846c040dcf73fc9dd",
	DotGitHash:   "7a725350b88b05ca03541b59dd0649fda7f521f2",
	ObjectsCount: 31,
}, {
	Tags:         []string{"packfile", "ref-delta", ".git"},
	URL:          "https://github.com/git-fixtures/basic.git",
	Head:         "6ecf0ef2c2dffb796033e5a02219af86ec6584e5",
	PackfileHash: "c544593473465e6315ad4182d04d366c4592b829",
	DotGitHash:   "7cbde0ca02f13aedd5ec8b358ca17b1c0bf5ee64",
	ObjectsCount: 31,
}, {
	Tags:         []string{"packfile", "ofs-delta", ".git", "single-branch"},
	URL:          "https://github.com/git-fixtures/basic.git",
	Head:         "6ecf0ef2c2dffb796033e5a02219af86ec6584e5",
	PackfileHash: "61f0ee9c75af1f9678e6f76ff39fbe372b6f1c45",
	DotGitHash:   "21504f6d2cc2ef0c9d6ebb8802c7b49abae40c1a",
	ObjectsCount: 28,
}, {
	Tags:       []string{".git", "merge-conflict"},
	URL:        "https://github.com/git-fixtures/basic.git",
	DotGitHash: "4870d54b5b04e43da8cf99ceec179d9675494af8",
}, {
	Tags:       []string{".git", "resolve-undo"},
	URL:        "https://github.com/git-fixtures/basic.git",
	DotGitHash: "df6781fd40b8f4911d70ce71f8387b991615cd6d",
}, {
	Tags:       []string{".git", "intent-to-add"},
	URL:        "https://github.com/git-fixtures/basic.git",
	DotGitHash: "4e7600af05c3356e8b142263e127b76f010facfc",
}, {
	Tags:       []string{".git", "index-v4"},
	URL:        "https://github.com/git-fixtures/basic.git",
	DotGitHash: "935e5ac17c41c309c356639816ea0694a568c484",
}, {
	Tags:       []string{".git", "end-of-index-entry"},
	URL:        "https://github.com/git-fixtures/basic.git",
	DotGitHash: "ab06771a67110b976953d34400d4dbc465ccd2d9",
}, {
	Tags:         []string{"worktree"},
	URL:          "https://github.com/git-fixtures/basic.git",
	WorktreeHash: "d2e42ddd68eacbb6034e7724e0dd4117ff1f01ee",
}, {
	Tags:         []string{"worktree", "submodule"},
	URL:          "https://github.com/git-fixtures/submodule.git",
	WorktreeHash: "8b4d55c85677b6b94bef2e46832ed2174ed6ecaf",
}, {
	Tags:         []string{"packfile", ".git", "unpacked", "multi-packfile"},
	URL:          "https://github.com/src-d/go-git.git",
	Head:         "e8788ad9165781196e917292d6055cba1d78664e",
	PackfileHash: "3559b3b47e695b33b0913237a4df3357e739831c",
	DotGitHash:   "174be6bd4292c18160542ae6dc6704b877b8a01a",
	ObjectsCount: 2133,
}, {
	Tags:         []string{"packfile", ".git", "tags"},
	URL:          "https://github.com/git-fixtures/tags.git",
	Head:         "f7b877701fbf855b44c0a9e86f3fdce2c298b07f",
	DotGitHash:   "c0c7c57ab1753ddbd26cc45322299ddd12842794",
	PackfileHash: "b68617dd8637fe6409d9842825a843a1d9a6e484",
	ObjectsCount: 7,
}, {
	Tags:         []string{"packfile"},
	URL:          "https://github.com/spinnaker/spinnaker.git",
	Head:         "06ce06d0fc49646c4de733c45b7788aabad98a6f",
	PackfileHash: "f2e0a8889a746f7600e07d2246a2e29a72f696be",
}, {
	Tags:         []string{"packfile"},
	URL:          "https://github.com/jamesob/desk.git",
	Head:         "d2313db6e7ca7bac79b819d767b2a1449abb0a5d",
	PackfileHash: "4ec6344877f494690fc800aceaf2ca0e86786acb",
}, {
	Tags:         []string{"packfile", "empty-folder"},
	URL:          "https://github.com/cpcs499/Final_Pres_P.git",
	Head:         "70bade703ce556c2c7391a8065c45c943e8b6bc3",
	PackfileHash: "29f304662fd64f102d94722cf5bd8802d9a9472c",
	DotGitHash:   "e1580a78f7d36791249df76df8a2a2613d629902",
}, {
	Tags:         []string{"packfile", "diff-tree"},
	URL:          "https://github.com/github/gem-builder.git",
	PackfileHash: "1ea0b3971fd64fdcdf3282bfb58e8cf10095e4e6",
}, {
	Tags:         []string{"packfile", "diff-tree"},
	URL:          "https://github.com/githubtraining/example-branches.git",
	PackfileHash: "bb8ee94710d3fa39379a630f76812c187217b312",
}, {
	Tags:         []string{"packfile", "diff-tree"},
	URL:          "https://github.com/rumpkernel/rumprun-xen.git",
	PackfileHash: "7861f2632868833a35fe5e4ab94f99638ec5129b",
}, {
	Tags:         []string{"packfile", "diff-tree"},
	URL:          "https://github.com/mcuadros/skeetr.git",
	PackfileHash: "36ef7a2296bfd526020340d27c5e1faa805d8d38",
}, {
	Tags:         []string{"packfile", "diff-tree"},
	URL:          "https://github.com/dezfowler/LiteMock.git",
	PackfileHash: "0d9b6cfc261785837939aaede5986d7a7c212518",
}, {
	Tags:         []string{"packfile", "diff-tree"},
	URL:          "https://github.com/tyba/storable.git",
	PackfileHash: "0d3d824fb5c930e7e7e1f0f399f2976847d31fd3",
}, {
	Tags:         []string{"packfile", "diff-tree"},
	URL:          "https://github.com/toqueteos/ts3.git",
	PackfileHash: "21b33a26eb7ffbd35261149fe5d886b9debab7cb",
}, {
	Tags:         []string{"empty", ".git"},
	URL:          "https://github.com/git-fixtures/empty.git",
	DotGitHash:   "bf3fedcc8e20fd0dec9172987ceea0038d17b516",
	ObjectsCount: 0,
}, {
	Tags:         []string{"worktree", "alternates"},
	WorktreeHash: "a6b6ff89c593f042347113203ead1c14ab5733ce",
}, {
	Tags:         []string{"worktree", "dirty"},
	WorktreeHash: "7203669c66103305e56b9dcdf940a7fbeb515f28",
}, {
	Tags:         []string{"packfile", "standalone"}, // standalone packfile that does not have any dependencies nor is part of any other fixture repo
	PackfileHash: "3638209d310e10ea8d90c362d568be65dd5e03a6",
}, {
	Tags:         []string{"thinpack"}, // adds commit on top of spinnaker fixture 06ce06d0fc49646c4de733c45b7788aabad98a6f via a thin pack
	PackfileHash: "ee4fef0ef8be5053ebae4ce75acf062ddf3031fb",
	Head:         "ee372bb08322c1e6e7c6c4f953cc6bf72784e7fb", // the thin pack adds this commit
}, {
	Tags:       []string{"merge-base"},
	DotGitHash: "26baa505b9f6fb2024b9999c140b75514718c988",
}, {
	Tags:         []string{"commit-graph"},
	Head:         "b9d69064b190e7aedccf84731ca1d917871f8a1c",
	PackfileHash: "769137af7784db501bca677fbd56fef8b52515b7",
	DotGitHash:   "cf717ccadce761d60bb4a8557a7b9a2efd23816a",
	ObjectsCount: 31,
}, {
	Tags:         []string{"worktree", "linked-worktree"},
	WorktreeHash: "363d996b02d9c3b598f0176619f5c6a44a82480a",
}, {
	Tags:         []string{"worktree", "main-branch", "no-master-head"},
	WorktreeHash: "e3b91f99d8d050cac81d84fbef89172f58eeb745",
}}

func All() Fixtures {
	return fixtures
}

func Basic() Fixtures {
	return ByURL("https://github.com/git-fixtures/basic.git").
		Exclude("single-branch")
}

func ByURL(url string) Fixtures {
	return fixtures.ByURL(url)
}

func ByTag(tag string) Fixtures {
	return fixtures.ByTag(tag)
}

type Fixture struct {
	URL          string
	Tags         []string
	Head         string
	PackfileHash string
	DotGitHash   string
	WorktreeHash string
	ObjectsCount int32
}

func (f *Fixture) Is(tag string) bool {
	for _, t := range f.Tags {
		if t == tag {
			return true
		}
	}

	return false
}

func (f *Fixture) file(path string) (billy.File, error) {
	if fpath, ok := files[path]; ok {
		return Filesystem.Open(fpath)
	}

	bytes, err := FSByte(false, "/data/"+path)
	if err != nil {
		return nil, err
	}

	file, err := Filesystem.TempFile("", "go-git-fixtures")
	if err != nil {
		return nil, err
	}

	if _, err := file.Write(bytes); err != nil {
		return nil, err
	}

	if err := file.Close(); err != nil {
		return nil, err
	}

	files[path] = file.Name()
	return Filesystem.Open(file.Name())
}

func (f *Fixture) Packfile() billy.File {
	file, err := f.file(fmt.Sprintf("pack-%s.pack", f.PackfileHash))
	if err != nil {
		panic(err)
	}

	return file
}

func (f *Fixture) Idx() billy.File {
	file, err := f.file(fmt.Sprintf("pack-%s.idx", f.PackfileHash))
	if err != nil {
		panic(err)
	}

	return file
}

// DotGit creates a new temporary directory and unpacks the repository .git
// directory into it. Multiple calls to DotGit returns different directories.
func (f *Fixture) DotGit() billy.Filesystem {
	if f.DotGitHash == "" && f.WorktreeHash != "" {
		fs, _ := f.Worktree().Chroot(".git")
		return fs.(billy.Filesystem)
	}

	file, err := f.file(fmt.Sprintf("git-%s.tgz", f.DotGitHash))
	if err != nil {
		panic(err)
	}

	fs, err := tgz.Extract(Filesystem, file.Name())
	if err != nil {
		panic(err)
	}

	return fs
}

// EnsureIsBare overrides the config file with one where bare is true.
func EnsureIsBare(fs billy.Filesystem) error {
	if _, err := fs.Stat("config"); err != nil {
		fmt.Printf("not .git folder: %s\n", err)
	}

	cfg, err := fs.OpenFile("config", os.O_TRUNC|os.O_WRONLY, 0)
	if err != nil {
		return err
	}

	defer cfg.Close()

	content := strings.NewReader("" +
		"[core]\n" +
		"repositoryformatversion = 0\n" +
		"filemode = true\n" +
		"bare = true\n" +
		"[http]\n" +
		"receivepack = true\n",
	)

	_, err = io.Copy(cfg, content)
	return err
}

func (f *Fixture) Worktree() billy.Filesystem {
	file, err := f.file(fmt.Sprintf("worktree-%s.tgz", f.WorktreeHash))
	if err != nil {
		panic(err)
	}

	fs, err := tgz.Extract(Filesystem, file.Name())
	if err != nil {
		panic(err)
	}

	return fs
}

type Fixtures []*Fixture

func (g Fixtures) Test(c *check.C, test func(*Fixture)) {
	for _, f := range g {
		c.Logf("executing test at %s %s", f.URL, f.Tags)
		test(f)
	}
}

func (g Fixtures) One() *Fixture {
	return g[0]
}

func (g Fixtures) ByTag(tag string) Fixtures {
	r := make(Fixtures, 0)
	for _, f := range g {
		if f.Is(tag) {
			r = append(r, f)
		}
	}

	return r
}
func (g Fixtures) ByURL(url string) Fixtures {
	r := make(Fixtures, 0)
	for _, f := range g {
		if f.URL == url {
			r = append(r, f)
		}
	}

	return r
}

func (g Fixtures) Exclude(tag string) Fixtures {
	r := make(Fixtures, 0)
	for _, f := range g {
		if !f.Is(tag) {
			r = append(r, f)
		}
	}

	return r
}

// Clean cleans all the temporal files created
func Clean() error {
	for fname, f := range files {
		if err := Filesystem.Remove(f); err != nil {
			return err
		}
		delete(files, fname)
	}
	return nil
}

type Suite struct{}

func (s *Suite) TearDownSuite(c *check.C) {
	Clean()
}
