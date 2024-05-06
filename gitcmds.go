package gitcmds

import (
	"strings"

	"github.com/astrica1/cmdutils"
)

type Repository interface {
	Clone() error
	Fetch() error
	BranchList() ([]string, error)
	RemoveBranch(branchName string) error
	RemoveTag(tagName string) error
	GetTagBranch(tagID string) ([]string, error)
	GetBranchTags(branchName string) ([]string, error)
	GetAllTagsList() ([]string, error)
	GetAllTags() ([]string, error)
}

type repository struct {
	url      string
	executer cmdutils.Executer
}

func NewRepository(repositoryURL string, executer cmdutils.Executer) Repository {
	return &repository{
		url:      repositoryURL,
		executer: executer,
	}
}

func (r *repository) Clone() error {
	cmd := "git clone " + r.url
	_, err := r.executer.Execute(cmd)

	return err
}

func (r *repository) Fetch() error {
	cmd := "git fetch --prune origin"
	_, err := r.executer.Execute(cmd)

	return err
}

func (r *repository) BranchList() ([]string, error) {
	cmd := "for branch in `git branch -r | grep -v HEAD`;do echo -e `git show --format=\"%ci %cr\" $branch | head -n 1` \t$branch; done | sort -r"

	branches, err := r.executer.Execute(cmd)
	if err != nil {
		return nil, err
	}

	if branches == "" {
		return nil, nil
	}

	branchesList := strings.Split(branches, "\n")

	return branchesList, nil
}

func (r *repository) RemoveBranch(branchName string) error {
	cmd := "git branch --delete " + branchName
	if _, err := r.executer.Execute(cmd); err != nil {
		return err
	}

	cmd = "git push origin --delete " + branchName
	if _, err := r.executer.Execute(cmd); err != nil {
		return err
	}

	return nil
}

func (r *repository) RemoveTag(tagName string) error {
	cmd := "git push origin :" + tagName
	if _, err := r.executer.Execute(cmd); err != nil {
		return err
	}

	cmd = "git push --delete origin " + tagName
	if _, err := r.executer.Execute(cmd); err != nil {
		return err
	}

	return nil
}

func (r *repository) GetTagBranch(tagID string) ([]string, error) {
	cmd := "git branch -a --contains " + tagID

	branches, err := r.executer.Execute(cmd)
	if err != nil {
		return nil, err
	}

	if branches == "" {
		return nil, nil
	}

	branchesList := strings.Split(branches, "\n")
	branchesList = branchesList[:len(branchesList)-1]

	for i, b := range branchesList {
		println(b)

		if strings.HasSuffix(b, "main") && len(b) < 8 {
			branchesList[i] = "main"
			continue
		}

		if strings.Contains(b, "HEAD ->") {
			branchesList = append(branchesList[:i], branchesList[i+1:]...)
		}
		// tmp := strings.Split(b, "/")
		// var idx int = -1
		// for j, t := range tmp {
		// 	if t == "origin" {
		// 		idx = j
		// 		break
		// 	}
		// }
		// if idx == -1 {
		// 	break
		// }
		// branchesList[i] = strings.Join(tmp[idx+1:], "/")
		tmp := strings.Split(b, "origin/")
		if len(tmp) > 1 {
			branchesList[i] = strings.Join(tmp[1:], "origin/")
		}
	}

	return branchesList, nil
}

func (r *repository) GetBranchTags(branchName string) ([]string, error) {
	cmd := "git describe --tags $(git rev-list --tags --max-count=150) " + branchName

	tags, err := r.executer.Execute(cmd)
	if err != nil {
		return nil, err
	}

	tagsList := strings.Split(tags, "\n")

	return tagsList[:len(tagsList)-1], nil
}

func (r *repository) GetAllTagsList() ([]string, error) {
	cmd := "git tag --list | xargs -n1 echo"

	tags, err := r.executer.Execute(cmd)
	if err != nil {
		return nil, err
	}

	tagsList := strings.Split(tags, "\n")

	return tagsList[:len(tagsList)-1], nil
}

func (r *repository) GetAllTags() ([]string, error) {
	cmd := "git for-each-ref --sort=-creatordate --format '%(creatordate:iso8601) %(refname:short)' refs/tags"

	tags, err := r.executer.Execute(cmd)
	if err != nil {
		return nil, err
	}

	tagsList := strings.Split(tags, "\n")

	return tagsList[:len(tagsList)-1], nil
}
